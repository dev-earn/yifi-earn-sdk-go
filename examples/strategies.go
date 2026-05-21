package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go"
)

func Strategy() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	// // GetStrategies
	// strategies, err := client.GetStrategies(context.Background(), models.StrategySearch{
	// 	ChainIDs: []string{"ethereum"},
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("get strategies error: %s", err.Error()))
	// }

	// bstrategies, err := json.MarshalIndent(strategies, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal strategy error: %s", err.Error()))
	// }

	// fmt.Println(string(bstrategies))

	// // GetStrategy
	// strategy, err := client.GetStrategy(context.Background(), "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0")
	// if err != nil {
	// 	panic(fmt.Sprintf("get strategy error: %s", err.Error()))
	// }

	// bstrategy, err := json.MarshalIndent(strategy, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal strategy error: %s", err.Error()))
	// }

	// fmt.Println(string(bstrategy))

	// GetStrategyHistory
	// strategyHistory, err := client.GetStrategyHistory(context.Background(), "ethereum:0x08c6F91e2B681FaF5e17227F2a44C307b3C1364C", models.FromTo{})
	// if err != nil {
	// 	panic(fmt.Sprintf("get strategy history error: %s", err.Error()))
	// }

	// bStrategyHistory, err := json.MarshalIndent(strategyHistory, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal strategy history error: %s", err.Error()))
	// }

	// fmt.Println(string(bStrategyHistory))

	// // GetEngines
	// endgines, err := client.GetEngines(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get engines error: %s", err.Error()))
	// }

	// bEngines, err := json.MarshalIndent(endgines, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal engines error: %s", err.Error()))
	// }

	// fmt.Println(string(bEngines))

	// // GetRiskScores
	// riskScores, err := client.GetRiskScores(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get risk scores error: %s", err.Error()))
	// }

	// briskScores, err := json.MarshalIndent(riskScores, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal reis scores error: %s", err.Error()))
	// }

	// fmt.Println(string(briskScores))

	// // GetStatuses
	// statuses, err := client.GetStrategyStatuses(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get statuses error: %s", err.Error()))
	// }

	// bStatuses, err := json.MarshalIndent(statuses, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal statuses error: %s", err.Error()))
	// }

	// fmt.Println(string(bStatuses))

	// // GetYields
	// yields, err := client.GetYields(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get yields error: %s", err.Error()))
	// }

	// byields, err := json.MarshalIndent(yields, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal yields error: %s", err.Error()))
	// }

	// fmt.Println(string(byields))

	// // GetRates
	// rates, err := client.GetRates(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get rates error: %s", err.Error()))
	// }

	// bRates, err := json.MarshalIndent(rates, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal rates error: %s", err.Error()))
	// }

	// fmt.Println(string(bRates))

	// BuildZap
	// zap, err := client.BuildZap(context.Background(), models.ZapQuoteRequest{
	// 	StrategyID:                  "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
	// 	InputToken:                  "ethereum:0x0000000000000000000000000000000000000000",
	// 	InputAmount:                 "10000000000",
	// 	OutputToken:                 "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
	// 	Sender:                      "0x053806E0642F6A6aAFE6176bd264cbbE446C8c92",
	// 	SlippageTolerancePercentage: decimal.NewFromFloat(1),
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("build zap error: %s", err.Error()))
	// }

	// bZap, err := json.MarshalIndent(zap, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal build zap error: %s", err.Error()))
	// }

	// fmt.Println(string(bZap))

	// GetWalletTransactions
	walletID, err := client.GetWalletTransactions(context.Background(), "0x053806E0642F6A6aAFE6176bd264cbbE446C8c92")
	if err != nil {
		panic(fmt.Sprintf("get wallet error: %s", err.Error()))
	}

	bWalletID, err := json.MarshalIndent(walletID, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal wallet error: %s", err.Error()))
	}

	fmt.Println(string(bWalletID))

}
