package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

func EthZap() {
	privateKey := privateKey()
	rpcURL := rpcURL("ethereum")

	auth := models.PartnerAuth{
		PartnerID: partnerID(),
		APIKey:    apiKey(),
		Secret:    apiSecret(),
	}

	yifi.UseTestnet = true

	testClient, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
	if err != nil {
		panic(fmt.Sprintf("yifi test client init failed: %s", err))
	}

	ctx := context.Background()

	strategyID := "ethereum:0xae7ab96520de3a18e5e111b5eaab095312d7fe84" // Lido stETH

	req := models.ZapQuoteRequest{
		StrategyID:                  strategyID,
		Sender:                      walletAddr(),
		InputToken:                  "ethereum:0x0000000000000000000000000000000000000000",
		InputAmount:                 "4500000000000000", // 0.0045 ETH
		OutputToken:                 strategyID,
		SlippageTolerancePercentage: decimal.NewFromInt(3),
	}
	quote, err := testClient.QuoteZap(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("quote zap: %s", err))
	}
	req.OrderID = quote.OrderID
	sig, err := testClient.BuildZap(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("build zap: %s", err))
	}
	printJSON("zap build", sig)

	txHash, err := signAndSend(ctx, rpcURL, privateKey, sig)
	if err != nil {
		panic(fmt.Sprintf("sign and send: %s", err))
	}
	fmt.Printf("tx sent: https://etherscan.io/tx/%s\n", txHash)
}
