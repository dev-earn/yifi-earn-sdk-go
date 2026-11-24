package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

func PublicReadOnlyCycle() {
	yifi.UseTestnet = true
	client, err := yifi.NewClient()
	must(err)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	readCheck("GetChains", func() (any, error) { return client.GetChains(ctx) })
	readCheck("GetPlatforms", func() (any, error) { return client.GetPlatforms(ctx) })
	readCheck("GetCategories", func() (any, error) { return client.GetCategories(ctx) })

	const tokenID = "ethereum:0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	readCheck("GetTokens", func() (any, error) {
		return client.GetTokens(ctx, models.TokensRequest{Chains: []string{"ethereum"}, Limit: 1})
	})
	readCheck("GetToken", func() (any, error) { return client.GetToken(ctx, tokenID) })
	readCheck("GetTokenHistoryMetrics", func() (any, error) {
		return client.GetTokenHistoryMetrics(ctx, tokenID, models.FromTo{})
	})
	readCheck("GetTokenGroup", func() (any, error) { return client.GetTokenGroup(ctx, tokenID) })
	readCheck("GetTokenTags", func() (any, error) { return client.GetTokenTags(ctx) })
	readCheck("GetTokenGroups", func() (any, error) { return client.GetTokenGroups(ctx) })

	markets := readCheck("GetMarkets", func() (any, error) {
		return client.GetMarkets(ctx, models.MarketsRequest{ChainIDs: []string{"ethereum"}, Limit: 1})
	}).(*models.Pagination[models.MarketsResponse])
	if markets == nil || len(markets.Data) == 0 {
		must(errors.New("no market available for read-only checks"))
	}
	marketID := markets.Data[0].MarketID
	readCheck("GetMarket", func() (any, error) { return client.GetMarket(ctx, marketID) })
	readCheck("GetMarketHistoryMetrics", func() (any, error) {
		return client.GetMarketHistoryMetrics(ctx, marketID, models.FromTo{})
	})
	readCheck("GetMarketTokenRoles", func() (any, error) { return client.GetMarketTokenRoles(ctx) })

	readCheck("GetStrategies", func() (any, error) {
		return client.GetStrategies(ctx, models.StrategySearch{ChainIDs: []string{"bsc"}, Limit: 1})
	})
	readCheck("GetStrategy", func() (any, error) { return client.GetStrategy(ctx, cycleStrategyID) })
	readCheck("GetStrategyHistory", func() (any, error) {
		return client.GetStrategyHistory(ctx, cycleStrategyID, models.FromTo{})
	})
	readCheck("GetEngines", func() (any, error) { return client.GetEngines(ctx) })
	readCheck("GetRiskScores", func() (any, error) { return client.GetRiskScores(ctx) })
	readCheck("GetStrategyStatuses", func() (any, error) { return client.GetStrategyStatuses(ctx) })
	readCheck("GetYields", func() (any, error) { return client.GetYields(ctx) })
	readCheck("GetRates", func() (any, error) { return client.GetRates(ctx) })

	req := models.ZapQuoteRequest{
		StrategyID:                  cycleStrategyID,
		Sender:                      walletAddr(),
		InputToken:                  cycleInputToken,
		InputAmount:                 cycleInputAmount,
		OutputToken:                 cycleStrategyID,
		SlippageTolerancePercentage: decimal.NewFromInt(1),
		Engine:                      "lifi",
	}
	quote := readCheck("QuoteZap", func() (any, error) { return client.QuoteZap(ctx, req) }).(models.ZapQuoteResponse)
	if quote.OrderID == "" || quote.MinOutputAmount == "" {
		must(errors.New("quote is missing order_id or min_output_amount"))
	}
	req.OrderID = quote.OrderID
	build := readCheck("BuildZap", func() (any, error) { return client.BuildZap(ctx, req) }).(models.Signature)
	if build.OrderID != quote.OrderID || build.MinOutputAmount != quote.MinOutputAmount || build.OutputAmount != quote.OutputAmount {
		must(errors.New("build differs from quote"))
	}

	const orderID = "1d76c85c02cf2a69d30a9faf2480f287"
	readCheck("GetTxInfo", func() (any, error) { return client.GetTxInfo(ctx, orderID) })
	readCheck("GetWalletTransaction", func() (any, error) {
		return client.GetWalletTransaction(ctx, walletAddr(), orderID)
	})
	readCheck("GetWalletTransactions", func() (any, error) {
		return client.GetWalletTransactions(ctx, walletAddr())
	})
	fmt.Println("public SDK cycle passed; no transaction was broadcast")
}

func readCheck(name string, call func() (any, error)) any {
	result, err := call()
	if err != nil {
		must(fmt.Errorf("%s: %w", name, err))
	}
	fmt.Printf("ok %s\n", name)
	if os.Getenv("YIFI_VERBOSE") == "1" {
		printJSON(name, result)
	}
	return result
}
