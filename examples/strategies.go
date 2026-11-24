package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

// Strategy demonstrates the public strategy endpoints and wallet lookup. All
// read-only.
func Strategy() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err))
	}
	ctx := context.Background()

	const stETHStrategy = "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0"

	// GetStrategies: paginated, filterable list.
	strategies, err := client.GetStrategies(ctx, models.StrategySearch{ChainIDs: []string{"ethereum"}})
	if err != nil {
		panic(fmt.Sprintf("get strategies: %s", err))
	}
	printJSON("strategies", strategies)

	// GetStrategy: current state of a single strategy.
	strategy, err := client.GetStrategy(ctx, stETHStrategy)
	if err != nil {
		panic(fmt.Sprintf("get strategy: %s", err))
	}
	printJSON("strategy", strategy)

	// GetStrategyHistory: historical stats over a time range.
	history, err := client.GetStrategyHistory(ctx, stETHStrategy, models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get strategy history: %s", err))
	}
	printJSON("strategy history", history)

	// The remaining calls return reference data used to classify strategies.
	engines, err := client.GetEngines(ctx)
	if err != nil {
		panic(fmt.Sprintf("get engines: %s", err))
	}
	printJSON("engines", engines)

	riskScores, err := client.GetRiskScores(ctx)
	if err != nil {
		panic(fmt.Sprintf("get risk scores: %s", err))
	}
	printJSON("risk scores", riskScores)

	statuses, err := client.GetStrategyStatuses(ctx)
	if err != nil {
		panic(fmt.Sprintf("get strategy statuses: %s", err))
	}
	printJSON("strategy statuses", statuses)

	yields, err := client.GetYields(ctx)
	if err != nil {
		panic(fmt.Sprintf("get yields: %s", err))
	}
	printJSON("yields", yields)

	rates, err := client.GetRates(ctx)
	if err != nil {
		panic(fmt.Sprintf("get rates: %s", err))
	}
	printJSON("rates", rates)

	// GetWalletTransactions: yifi-tracked transactions for a wallet.
	txs, err := client.GetWalletTransactions(ctx, walletAddr())
	if err != nil {
		panic(fmt.Sprintf("get wallet transactions: %s", err))
	}
	printJSON("wallet transactions", txs)
}

// StrategyZapQuote builds an (unsigned) zap into a strategy and prints the
// quote. This only fetches the transaction; it does not sign or broadcast — see
// Zap in zap.go for the full signing flow.
func StrategyZapQuote() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err))
	}
	ctx := context.Background()

	const stETHStrategy = "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0"

	req := models.ZapQuoteRequest{
		StrategyID:                  stETHStrategy,
		InputToken:                  "ethereum:0x0000000000000000000000000000000000000000", // native ETH
		InputAmount:                 "10000000000000000",                                   // 0.01 ETH (wei)
		OutputToken:                 stETHStrategy,
		Sender:                      walletAddr(),
		SlippageTolerancePercentage: decimal.NewFromInt(1),
	}
	quote, err := client.QuoteZap(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("quote zap: %s", err))
	}
	req.OrderID = quote.OrderID
	zap, err := client.BuildZap(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("build zap: %s", err))
	}
	printJSON("zap build", zap)
}
