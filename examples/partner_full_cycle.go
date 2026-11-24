package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func PartnerFullCycle() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cfg := cycleConfigForRun()
	yifi.UseTestnet = true
	client, auth := partnerTestClient()
	partnerID := auth.PartnerID
	fmt.Printf("running full cycle for existing partner: %s\n", partnerID)

	// Preserve the original policy so the temporary test fee can be removed
	// without changing the partner's production configuration.
	originalPolicy, err := client.GetPartnerPolicy(ctx, partnerID)
	must(err)
	printJSON("original policy", originalPolicy)

	// Exercise API-key creation, listing, and deferred revocation. The existing
	// credential remains in use for signed requests; this key is only for testing.
	created, err := client.CreateAPIKey(ctx, partnerID, fmt.Sprintf("full-cycle-%d", time.Now().Unix()))
	must(err)
	fmt.Printf("created test api key: %s\n", created.APIKey.APIKey)
	defer func() {
		if err := client.RevokeAPIKey(ctx, partnerID, created.APIKey.APIKey); err != nil {
			fmt.Printf("revoke test api key failed: %v\n", err)
		}
	}()

	keys, err := client.GetAPIKeys(ctx, partnerID)
	must(err)
	if len(keys) == 0 {
		must(errors.New("created api key is missing from key list"))
	}
	supportedChains, err := client.GetSupportedChains(ctx)
	must(err)
	chainSupported := false
	for _, chain := range supportedChains {
		if int64(chain.ChainID) == cfg.chainID && chain.ChainName != "" {
			chainSupported = true
			break
		}
	}
	if !chainSupported {
		must(fmt.Errorf("chain %d is missing from supported partner chains", cfg.chainID))
	}

	// Apply a temporary 0.3% fee. Restoration is deferred so even a failed cycle
	// attempts to return the policy to its original state.
	commonFee := decimal.RequireFromString("0.3")
	policy := originalPolicy
	policy.Fees = &models.FeePolicy{CommonFee: &commonFee}
	must(client.UpdatePartnerPolicy(ctx, partnerID, policy))
	defer func() {
		if err := client.UpdatePartnerPolicy(ctx, partnerID, originalPolicy); err != nil {
			fmt.Printf("restore policy failed: %v\n", err)
		}
	}()

	key, err := signerKey()
	must(err)
	signer := keyAddress(key)
	rpc := rpcClient(ctx, cfg.chain)
	defer rpc.Close()
	chainID, err := rpc.ChainID(ctx)
	must(err)
	if chainID.Int64() != cfg.chainID {
		must(fmt.Errorf("unexpected chain id: %s", chainID.String()))
	}

	tokenAddress := strings.TrimPrefix(cfg.inputToken, cfg.chain+":")
	requiredAmount, ok := newAmount(cfg.inputAmount)
	if !ok {
		must(errors.New("invalid input amount"))
	}
	balance, err := erc20Uint256(ctx, rpc, tokenAddress, "0x70a08231", signer)
	must(err)
	nativeBalance, err := rpc.BalanceAt(ctx, signer, nil)
	must(err)
	fmt.Printf("wallet: %s native(wei)=%s input(base units)=%s\n", signer.Hex(), nativeBalance.String(), balance.String())
	if (tokenAddress != "0x0000000000000000000000000000000000000000" && balance.Cmp(requiredAmount) < 0) || nativeBalance.Sign() == 0 {
		must(errors.New("insufficient wallet funds"))
	}

	// A resume order skips quote/build/broadcast so a transient timeout never
	// forces a duplicate zap. Without it, quote and build must share the order ID
	// and output amounts before the partner's HMAC-attributed transaction is signed.
	resumeOrderID := getenv("YIFI_RESUME_ORDER_ID", "")
	quote := models.ZapQuoteResponse{OrderID: resumeOrderID}
	txHash := getenv("YIFI_RESUME_TX_HASH", "")
	if resumeOrderID == "" {
		req := models.ZapQuoteRequest{
			StrategyID:                  cfg.strategyID,
			Sender:                      signer.Hex(),
			InputToken:                  cfg.inputToken,
			InputAmount:                 cfg.inputAmount,
			OutputToken:                 cfg.strategyID,
			SlippageTolerancePercentage: decimal.NewFromInt(1),
			Engine:                      cfg.engine,
		}
		quote, err = client.QuoteZap(ctx, req)
		must(err)
		printJSON("partner quote", quote)
		req.OrderID = quote.OrderID
		build, err := client.BuildZap(ctx, req)
		must(err)
		printJSON("partner build", build)
		if build.OrderID != quote.OrderID || build.MinOutputAmount != quote.MinOutputAmount || build.OutputAmount != quote.OutputAmount {
			must(errors.New("build differs from quote"))
		}
		// ERC-20 inputs require an allowance for build's spender. Native inputs
		// skip this check and carry value in the signed transaction instead.
		if tokenAddress != "0x0000000000000000000000000000000000000000" {
			allowance, err := erc20Allowance(ctx, rpc, tokenAddress, signer, common.HexToAddress(build.ApprovalAddress))
			must(err)
			if allowance.Cmp(requiredAmount) < 0 {
				must(fmt.Errorf("insufficient allowance for %s", build.ApprovalAddress))
			}
		}

		txHash, err = signAndSend(ctx, rpcURL(cfg.chain), privateKey(), build)
		must(err)
		fmt.Printf("partner zap sent: %s\n", txHash)
		receipt := waitForReceipt(ctx, rpc, common.HexToHash(txHash))
		printJSON("partner zap receipt", receipt)
		if receipt.Status != types.ReceiptStatusSuccessful {
			must(errors.New("partner zap reverted"))
		}
	}

	// Execution status and partner attribution arrive asynchronously. Poll both
	// independently to avoid assuming the mined receipt already updated the API.
	statusFound := false
	for attempt := 1; attempt <= 20; attempt++ {
		status, err := client.GetTxInfo(ctx, quote.OrderID)
		if err == nil {
			printJSON("partner execution status", status)
			statusFound = true
			break
		}
		time.Sleep(3 * time.Second)
	}
	if !statusFound {
		must(errors.New("partner execution status unavailable after polling"))
	}

	// Confirm that the partner feed contains this exact hash before checking fees.
	transactionFound := false
	for attempt := 1; attempt <= 20; attempt++ {
		transactions, err := client.GetPartnerTransactions(ctx, partnerID, models.PartnerTransactionsRequest{Hash: txHash, Limit: 10})
		if err == nil && len(transactions.Items) > 0 {
			printJSON("partner transactions", transactions)
			transactionFound = true
			break
		}
		time.Sleep(3 * time.Second)
	}
	if !transactionFound {
		must(errors.New("partner transaction unavailable after polling"))
	}

	// Fees are accounted in the token that actually settled on-chain. For native
	// inputs this is the strategy's wrapped asset, not the zero address.
	feeTokenID := cfg.inputToken
	if feeTokenID == cfg.chain+":0x0000000000000000000000000000000000000000" {
		feeTokenID = cfg.strategyID
	}
	balanceBeforeClaim, err := client.GetFeeBalance(ctx, partnerID)
	must(err)
	printJSON("fee balance before claim", balanceBeforeClaim)
	if !cycleFeeAmount(balanceBeforeClaim, feeTokenID).IsPositive() {
		must(errors.New("partner fee was not accrued"))
	}

	addresses, err := client.GetAllClaimAddresses(ctx, partnerID)
	must(err)
	claimAddress := ""
	for _, address := range addresses {
		if int64(address.ChainID) == cfg.chainID && address.Status {
			claimAddress = address.Address
			break
		}
	}
	// Reuse the active claim address for this chain or attach the signer when
	// none exists. The claim must be signed by this same address.
	if claimAddress == "" {
		must(client.AttachNewClaimAddress(ctx, models.ClaimAddressPublic{
			PartnerID: partnerID,
			ChainID:   int(cfg.chainID),
			Address:   signer.Hex(),
		}))
		claimAddress = signer.Hex()
	}
	fmt.Printf("claim address: %s\n", claimAddress)
	// Opt in to the TOTP-protected update with the same address. This verifies
	// the route without changing claim ownership or redirecting funds.
	if getenv("YIFI_UPDATE_CLAIM_ADDRESS", "") == "1" {
		claimName := "full-cycle"
		must(client.UpdateExistingClaimAddress(ctx, models.ClaimAddress{
			PartnerID: partnerID,
			ChainID:   int(cfg.chainID),
			Address:   claimAddress,
			Name:      &claimName,
			TOTPCode:  getenv("YIFI_TOTP_CODE", ""),
		}))
		fmt.Println("claim address updated with TOTP")
	}

	// Claim the accrued fee with either the single-token or all-tokens endpoint.
	// The API credits the distributor first and returns unsigned claim calldata.
	var claim models.TxInfo
	if getenv("YIFI_CLAIM_MODE", "") == "single" {
		claim, err = client.ClaimFees(ctx, models.PartnerFeeClaimRequest{
			PartnerID:    partnerID,
			TokenID:      feeTokenID,
			ChainID:      cfg.chain,
			ClaimAddress: claimAddress,
		})
	} else {
		claim, err = client.ClaimAllFees(ctx, models.PartnerFeeClaimAllRequest{
			PartnerID:    partnerID,
			ChainID:      cfg.chain,
			ClaimAddress: claimAddress,
		})
	}
	must(err)
	printJSON("claim transaction", claim)
	if claim.Data == "" {
		must(errors.New("claim returned no calldata despite a positive fee balance"))
	}
	claimHash, err := signAndSend(ctx, rpcURL(cfg.chain), privateKey(), models.Signature{
		From:  claim.From,
		To:    claim.To,
		Data:  claim.Data,
		Value: "0",
	})
	must(err)
	fmt.Printf("fee claim sent: %s\n", claimHash)
	claimReceipt := waitForReceipt(ctx, rpc, common.HexToHash(claimHash))
	printJSON("claim receipt", claimReceipt)
	if claimReceipt.Status != types.ReceiptStatusSuccessful {
		must(errors.New("fee claim reverted"))
	}

	// The claim is settled only after the outgoing distributor webhook updates
	// the partner ledger. Poll the actual fee token until it reaches zero.
	for attempt := 1; attempt <= 20; attempt++ {
		balanceAfterClaim, err := client.GetFeeBalance(ctx, partnerID)
		must(err)
		if cycleFeeAmount(balanceAfterClaim, feeTokenID).IsZero() {
			printJSON("fee balance after claim", balanceAfterClaim)
			break
		}
		if attempt == 20 {
			printJSON("fee balance after claim", balanceAfterClaim)
			must(errors.New("fee balance did not reach zero after claim webhook"))
		}
		time.Sleep(3 * time.Second)
	}

	if getenv("YIFI_CLAIM_MODE", "") == "single" {
		emptyClaim, err := client.ClaimAllFees(ctx, models.PartnerFeeClaimAllRequest{
			PartnerID:    partnerID,
			ChainID:      cfg.chain,
			ClaimAddress: claimAddress,
		})
		must(err)
		printJSON("claim-all after settlement", emptyClaim)
	} else {
		emptyClaim, err := client.ClaimFees(ctx, models.PartnerFeeClaimRequest{
			PartnerID:    partnerID,
			TokenID:      feeTokenID,
			ChainID:      cfg.chain,
			ClaimAddress: claimAddress,
		})
		must(err)
		printJSON("single-token claim after settlement", emptyClaim)
	}

	// Finish by checking aggregate and per-strategy reporting, not only the
	// execution status and fee balance.
	summary, err := client.GetPartnerSummary(ctx, partnerID, models.FromTo{})
	must(err)
	printJSON("partner summary", summary)
	metrics, err := client.GetPartnerMetrics(ctx, partnerID, models.FromTo{})
	must(err)
	printJSON("partner metrics", metrics)
	performance, err := client.GetStrategyPerformance(ctx, partnerID, models.FromTo{})
	must(err)
	printJSON("strategy performance", performance)
}

func cycleFeeAmount(balance models.Balance, tokenID string) decimal.Decimal {
	for _, token := range balance.Tokens {
		if token.TokenID == tokenID {
			return token.Amount.Total
		}
	}
	return decimal.Zero
}

func signerKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(privateKey(), "0x"))
}

func keyAddress(key *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(key.PublicKey)
}

func rpcClient(ctx context.Context, chain string) *ethclient.Client {
	client, err := ethclient.DialContext(ctx, rpcURL(chain))
	must(err)
	return client
}

func newAmount(raw string) (*big.Int, bool) {
	return new(big.Int).SetString(raw, 10)
}
