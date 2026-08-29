package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

func Zap() {
	privateKey := privateKey()
	rpcURL := rpcURL("arbitrum")

	auth := models.PartnerAuth{
		PartnerID: partnerID(),
		APIKey:    apiKey(),
		Secret:    apiSecret(),
	}

	testClient, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
	if err != nil {
		panic(fmt.Sprintf("yifi test client init failed: %s", err))
	}

	ctx := context.Background()

	fee := decimal.NewFromFloat(0.4)
	err = testClient.UpdatePartnerPolicy(ctx, auth.PartnerID, models.Policy{
		Fees: &models.FeePolicy{
			CommonFee: &fee,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("update policy: %s", err))
	}
	fmt.Println("fee updated to 0.4%")

	// name := "bsc-main"
	// err = testClient.UpdateExistingClaimAddress(ctx, models.ClaimAddress{
	// 	PartnerID: auth.PartnerID,
	// 	ChainID:   56,
	// 	Address:   "0x256DD8D76C0d258A529503cB4Ac2ccDa94c437F3",
	// 	Name:      &name,
	// })
	// if err != nil && !strings.Contains(err.Error(), "unmarshal") {
	// 	panic(fmt.Sprintf("update claim address: %s", err))
	// }
	// fmt.Println("claim address updated")

	// Arbitrum: USDT -> aArbUSDT
	sig, err := testClient.BuildZap(ctx, models.ZapQuoteRequest{
		StrategyID:                  "arbitrum:0x6ab707aca953edaefbc4fd23ba73294241490620",
		Sender:                      walletAddr(),
		InputToken:                  "arbitrum:0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9", // USDT (6 decimals)
		InputAmount:                 "50000",                                               // 0.05 USDT
		OutputToken:                 "arbitrum:0x6ab707aca953edaefbc4fd23ba73294241490620", // aArbUSDT
		SlippageTolerancePercentage: decimal.NewFromInt(3),
	})
	if err != nil {
		panic(fmt.Sprintf("build zap: %s", err))
	}
	printJSON("zap quote", sig)

	txHash, err := signAndSend(ctx, rpcURL, privateKey, sig)
	if err != nil {
		panic(fmt.Sprintf("sign and send: %s", err))
	}
	fmt.Printf("tx sent: https://arbiscan.io/tx/%s\n", txHash)
}

// Claim pulls a single token's accrued fees out of the TokenDistributor on
// optimism. It attaches the claim address, calls ClaimFees (which pushes the
// free balance to pending via addBalances and returns the unsigned claim(token)
// calldata), then signs and broadcasts it — firing the distributor's Claimed
// event that the "out" webhook consumes.
func Claim() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	// The claim address MUST be the signer of the claim tx (see broadcastClaim).
	rpcURL := rpcURL("arbitrum")
	privateKey := privateKey()
	claimAddr := walletAddr()

	name := "arb-main"
	if err := client.AttachNewClaimAddress(ctx, models.ClaimAddressPublic{
		PartnerID: auth.PartnerID,
		ChainID:   42161,
		Address:   claimAddr,
		Name:      &name,
	}); err != nil {
		fmt.Printf("attach claim address (may already exist): %s\n", err)
	}

	tx, err := client.ClaimFees(ctx, models.PartnerFeeClaimRequest{
		PartnerID: auth.PartnerID,
		TokenID:   "arbitrum:0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9", // USDT
		ChainID:   "arbitrum",
	})
	if err != nil {
		panic(fmt.Sprintf("claim fees: %s", err))
	}
	printJSON("claim tx", tx)

	if tx.Data == "" {
		fmt.Println("nothing to claim: no pending balance on-chain")
		return
	}

	txHash, err := signAndSend(ctx, rpcURL, privateKey, models.Signature{
		From:  tx.From,
		To:    tx.To,
		Data:  tx.Data,
		Value: "0",
	})
	if err != nil {
		panic(fmt.Sprintf("sign and send claim: %s", err))
	}
	fmt.Printf("claim tx sent: https://arbiscan.io/tx/%s\n", txHash)
}

// ClaimAll pulls every pending token's fees out of the TokenDistributor in one
// tx. It calls ClaimAllFees for the unsigned claimAll(tokens) calldata, then
// signs and broadcasts it from the partner's claim address.
func ClaimAll() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	tx, err := client.ClaimAllFees(ctx, models.PartnerFeeClaimAllRequest{
		PartnerID: auth.PartnerID,
		ChainID:   "bsc",
	})
	if err != nil {
		panic(fmt.Sprintf("claim all fees: %s", err))
	}
	broadcastClaim(ctx, tx)
}

// partnerTestClient builds the partner-authed client used by the claim examples.
func partnerTestClient() (yifi.SDK, models.PartnerAuth) {
	auth := models.PartnerAuth{
		PartnerID: partnerID(),
		APIKey:    apiKey(),
		Secret:    apiSecret(),
	}

	client, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
	if err != nil {
		panic(fmt.Sprintf("yifi client init failed: %s", err))
	}
	return client, auth
}

// broadcastClaim signs the unsigned claim tx returned by ClaimFees/ClaimAllFees
// and sends it. The claim tx carries no value; it invokes claim()/claimAll() on
// the distributor (tx.To) from the partner claim address (tx.From), which MUST
// be the signer — it is the msg.sender the contract pays out to.
func broadcastClaim(ctx context.Context, tx models.TxInfo) {
	printJSON("claim tx", tx)

	if tx.Data == "" {
		fmt.Println("nothing to claim: no pending balance on-chain")
		return
	}

	// The signer of the claim tx MUST be the partner's claim address.
	privateKey := privateKey()
	rpcURL := rpcURL("bsc")

	txHash, err := signAndSend(ctx, rpcURL, privateKey, models.Signature{
		From:  tx.From,
		To:    tx.To,
		Data:  tx.Data,
		Value: "0",
	})
	if err != nil {
		panic(fmt.Sprintf("sign and send claim: %s", err))
	}
	fmt.Printf("claim tx sent: https://bscscan.com/tx/%s\n", txHash)
}
