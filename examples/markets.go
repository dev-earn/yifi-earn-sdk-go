package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

func Market() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	// GetMarkets
	// markets, err := client.GetMarkets(context.Background(), models.MarketsRequest{})
	// if err != nil {
	// 	panic(fmt.Sprintf("get markets error: %s", err.Error()))
	// }

	// bMarkets, err := json.MarshalIndent(markets, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal markets error: %s", err.Error()))
	// }

	// fmt.Println(string(bMarkets))

	// GetMarket
	// market, err := client.GetMarket(context.Background(), "ethereum:0x08c6F91e2B681FaF5e17227F2a44C307b3C1364C")
	// if err != nil {
	// 	panic(fmt.Sprintf("get market error: %s", err.Error()))
	// }

	// bMarket, err := json.MarshalIndent(market, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal market error: %s", err.Error()))
	// }

	// fmt.Println(string(bMarket))

	// GetMarketTokenRoles
	// roles, err := client.GetMarketTokenRoles(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get roles error: %s", err.Error()))
	// }

	// broles, err := json.MarshalIndent(roles, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal roles error: %s", err.Error()))
	// }

	// fmt.Println(string(broles))

	// GetMarketHistoryMetrics
	marketHistory, err := client.GetMarketHistoryMetrics(context.Background(), "ethereum:0x08c6F91e2B681FaF5e17227F2a44C307b3C1364C", models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get market history error: %s", err.Error()))
	}

	bMarketHistory, err := json.MarshalIndent(marketHistory, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal market history error: %s", err.Error()))
	}

	fmt.Println(string(bMarketHistory))
}
