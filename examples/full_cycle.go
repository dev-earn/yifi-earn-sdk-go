package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

const (
	cycleStrategyID  = "bsc:0xfd5840cd36d94d7229439859c0112a4185bc0255"
	cycleInputToken  = "bsc:0x55d398326f99059ff775485246999027b3197955"
	cycleInputAmount = "2000000000000000000"
)

func FullCycle() {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	// The configuration selects the engine, strategy, chain, and input amount.
	// Environment overrides allow the same example to run LI.FI and Portals.
	cfg := cycleConfigForRun()

	yifi.UseTestnet = true
	client, err := yifi.NewClient()
	must(err)

	// Load the strategy first so the user can verify its engine and active status
	// before requesting a quote.
	strategy, err := client.GetStrategy(ctx, cfg.strategyID)
	must(err)
	fmt.Printf("strategy: id=%s name=%s engine=%s status=%s\n",
		strategy.StrategyID, strategy.Name, strategy.RoutingEngine, strategy.Status)

	// The private key must resolve to the configured wallet. This prevents
	// accidentally signing from a different account than the quote sender.
	key, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey(), "0x"))
	must(err)
	signer := crypto.PubkeyToAddress(key.PublicKey)
	if !strings.EqualFold(signer.Hex(), walletAddr()) {
		must(fmt.Errorf("signer mismatch: key resolves to %s, configured wallet is %s", signer.Hex(), walletAddr()))
	}

	rpc, err := ethclient.DialContext(ctx, rpcURL(cfg.chain))
	must(err)
	defer rpc.Close()
	chainID, err := rpc.ChainID(ctx)
	must(err)
	if chainID.Int64() != cfg.chainID {
		must(fmt.Errorf("unexpected chain id: %s", chainID.String()))
	}

	tokenAddress := strings.TrimPrefix(cfg.inputToken, cfg.chain+":")
	requiredAmount, ok := new(big.Int).SetString(cfg.inputAmount, 10)
	if !ok {
		must(fmt.Errorf("invalid test amount: %s", cfg.inputAmount))
	}
	balance, err := erc20Uint256(ctx, rpc, tokenAddress, "0x70a08231", signer)
	must(err)
	nativeBalance, err := rpc.BalanceAt(ctx, signer, nil)
	must(err)
	fmt.Printf("wallet: %s native(wei)=%s input(base units)=%s\n", signer.Hex(), nativeBalance.String(), balance.String())
	// Native inputs use the zero address and do not have an ERC-20 balance.
	// They still require enough native coin to pay transaction gas.
	if tokenAddress != "0x0000000000000000000000000000000000000000" && balance.Cmp(requiredAmount) < 0 {
		must(fmt.Errorf("insufficient input balance for test amount %s", cfg.inputAmount))
	}
	if nativeBalance.Sign() == 0 {
		must(errors.New("insufficient native coin for gas"))
	}

	// Quote freezes the route and its minimum output for one minute. Build must
	// use the same request plus this order ID so it cannot silently re-route.
	req := models.ZapQuoteRequest{
		StrategyID:                  cfg.strategyID,
		Sender:                      signer.Hex(),
		InputToken:                  cfg.inputToken,
		InputAmount:                 cfg.inputAmount,
		OutputToken:                 cfg.strategyID,
		SlippageTolerancePercentage: decimal.NewFromInt(1),
		Engine:                      cfg.engine,
	}
	quote, err := client.QuoteZap(ctx, req)
	must(err)
	printJSON("quote", quote)
	if quote.OrderID == "" || quote.MinOutputAmount == "" {
		must(errors.New("quote is missing order_id or min_output_amount"))
	}

	req.OrderID = quote.OrderID
	build, err := client.BuildZap(ctx, req)
	must(err)
	printJSON("build", build)
	// Verify that the returned transaction belongs to the quote the user saw,
	// including its expected and minimum output amounts and sender.
	if build.OrderID != quote.OrderID {
		must(fmt.Errorf("order id changed: quote=%s build=%s", quote.OrderID, build.OrderID))
	}
	if build.MinOutputAmount != quote.MinOutputAmount || build.OutputAmount != quote.OutputAmount {
		must(errors.New("build output amounts differ from quote"))
	}
	if !strings.EqualFold(build.From, signer.Hex()) {
		must(fmt.Errorf("build sender mismatch: %s", build.From))
	}
	// ERC-20 inputs need an allowance for the exact spender returned by build.
	// Native inputs skip this check and attach value to the transaction instead.
	if tokenAddress != "0x0000000000000000000000000000000000000000" {
		if build.ApprovalAddress == "" {
			must(errors.New("build has no approval address"))
		}
		allowance, err := erc20Allowance(ctx, rpc, tokenAddress, signer, common.HexToAddress(build.ApprovalAddress))
		must(err)
		fmt.Printf("allowance for %s: %s input base units\n", build.ApprovalAddress, allowance.String())
		if allowance.Cmp(requiredAmount) < 0 {
			must(fmt.Errorf("insufficient allowance; approve %s before running this test", build.ApprovalAddress))
		}
	}

	// Before broadcast, status may legitimately be unavailable because the
	// engine has not seen an execution yet. This is informational, not a failure.
	if status, err := client.GetTxInfo(ctx, quote.OrderID); err != nil {
		fmt.Printf("pre-broadcast YiFi status: unavailable: %v\n", err)
	} else {
		printJSON("pre-broadcast YiFi status", status)
	}

	// Sign and broadcast the frozen transaction. The order ID is enough for
	// subsequent status checks; no transaction-hash registration is required.
	txHash, err := signAndSend(ctx, rpcURL(cfg.chain), privateKey(), build)
	must(err)
	fmt.Printf("transaction sent: %s\n", txHash)

	receipt := waitForReceipt(ctx, rpc, common.HexToHash(txHash))
	printJSON("on-chain receipt", receipt)
	if receipt.Status != types.ReceiptStatusSuccessful {
		must(fmt.Errorf("transaction reverted on-chain: %s", txHash))
	}

	// Engine status and stored wallet history arrive asynchronously through the
	// engine and distributor webhooks, so poll both independently.
	for attempt := 1; attempt <= 20; attempt++ {
		status, err := client.GetTxInfo(ctx, quote.OrderID)
		if err == nil {
			printJSON("YiFi status", status)
			break
		}
		if attempt == 20 {
			fmt.Printf("YiFi status unavailable after polling: %v\n", err)
		}
		time.Sleep(3 * time.Second)
	}

	for attempt := 1; attempt <= 20; attempt++ {
		stored, err := client.GetWalletTransaction(ctx, signer.Hex(), quote.OrderID)
		if err == nil {
			printJSON("stored YiFi transaction", stored)
			break
		}
		if attempt == 20 {
			fmt.Printf("stored YiFi transaction unavailable after polling: %v\n", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func erc20Uint256(ctx context.Context, rpc *ethclient.Client, token, selector string, owner common.Address) (*big.Int, error) {
	data := append(common.FromHex(selector), common.LeftPadBytes(owner.Bytes(), 32)...)
	to := common.HexToAddress(token)
	result, err := rpc.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(result), nil
}

func erc20Allowance(ctx context.Context, rpc *ethclient.Client, token string, owner, spender common.Address) (*big.Int, error) {
	data := append(common.FromHex("0xdd62ed3e"), common.LeftPadBytes(owner.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(spender.Bytes(), 32)...)
	to := common.HexToAddress(token)
	result, err := rpc.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(result), nil
}

func waitForReceipt(ctx context.Context, rpc *ethclient.Client, hash common.Hash) *types.Receipt {
	for attempt := 1; attempt <= 60; attempt++ {
		receipt, err := rpc.TransactionReceipt(ctx, hash)
		if err == nil {
			return receipt
		}
		if !errors.Is(err, ethereum.NotFound) {
			fmt.Printf("receipt lookup error: %v\n", err)
		}
		if attempt == 60 {
			must(fmt.Errorf("transaction not mined in time: %s", hash.Hex()))
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}
