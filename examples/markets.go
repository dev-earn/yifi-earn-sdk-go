package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

// Market demonstrates every public market endpoint. All read-only.
func Market() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err))
	}
	ctx := context.Background()

	const marketID = "ethereum:0x08c6F91e2B681FaF5e17227F2a44C307b3C1364C"

	// GetMarkets: paginated, filterable list.
	markets, err := client.GetMarkets(ctx, models.MarketsRequest{})
	if err != nil {
		panic(fmt.Sprintf("get markets: %s", err))
	}
	printJSON("markets", markets)

	// GetMarket: current state of a single market.
	market, err := client.GetMarket(ctx, marketID)
	if err != nil {
		panic(fmt.Sprintf("get market: %s", err))
	}
	printJSON("market", market)

	// GetMarketTokenRoles: roles a token can play in a market (supply, borrow, …).
	roles, err := client.GetMarketTokenRoles(ctx)
	if err != nil {
		panic(fmt.Sprintf("get market token roles: %s", err))
	}
	printJSON("market token roles", roles)

	// GetMarketHistoryMetrics: historical stats over a time range.
	history, err := client.GetMarketHistoryMetrics(ctx, marketID, models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get market history: %s", err))
	}
	printJSON("market history", history)
}
