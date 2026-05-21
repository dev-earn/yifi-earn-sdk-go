package main

import (
	"context"
	"encoding/json"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

func Partner() {
	auth := models.PartnerAuth{
		PartnerID: "PARTNER_ID",
		APIKey:    "API_KEY",
		Secret:    "SECRET",
	}

	client, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	ctx := context.Background()

	// summary, err := client.GetPartnerSummary(ctx, auth.PartnerID, models.FromTo{})
	// if err != nil {
	// 	panic(fmt.Sprintf("get partner summary error: %s", err.Error()))
	// }
	// printJSON("partner summary", summary)

	// metrics, err := client.GetPartnerMetrics(ctx, auth.PartnerID, models.FromTo{})
	// if err != nil {
	// 	panic(fmt.Sprintf("get partner metrics error: %s", err.Error()))
	// }
	// printJSON("partner metrics", metrics)

	policy, err := client.GetPartnerPolicy(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get partner policy error: %s", err.Error()))
	}
	printJSON("partner policy", policy)

	partnerFeePref := decimal.NewFromFloat(0.7)
	err = client.UpdatePartnerPolicy(ctx, auth.PartnerID, models.Policy{
		Fees: &models.FeePolicy{
			CommonFee: &partnerFeePref,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("update partner policy error: %s", err.Error()))
	}

	policyUPD, err := client.GetPartnerPolicy(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get partner policy error: %s", err.Error()))
	}
	printJSON("partner policy", policyUPD)

	// apiKeys, err := client.GetAPIKeys(ctx, auth.PartnerID)
	// if err != nil {
	// 	panic(fmt.Sprintf("get api keys error: %s", err.Error()))
	// }
	// printJSON("api keys", apiKeys)

	// apiKey, err := client.CreateAPIKey(ctx, auth.PartnerID, "demo-key")
	// if err != nil {
	// 	panic(fmt.Sprintf("create api key error: %s", err.Error()))
	// }
	// printJSON("create api key", apiKey)

	// err = client.RevokeAPIKey(ctx, auth.PartnerID, "api-key-id-example")
	// if err != nil {
	// 	panic(fmt.Sprintf("revoke api key error: %s", err.Error()))
	// }

	// feeBalance, err := client.GetFeeBalance(ctx, auth.PartnerID)
	// if err != nil {
	// 	panic(fmt.Sprintf("get fee balance error: %s", err.Error()))
	// }
	// printJSON("fee balance", feeBalance)

	// txInfo, err := client.ClaimFees(ctx, models.PartnerFeeClaimRequest{
	// 	PartnerID: auth.PartnerID,
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("claim fees error: %s", err.Error()))
	// }
	// printJSON("claim fees", txInfo)

	// supportedChains, err := client.GetSupportedChains(ctx)
	// if err != nil {
	// 	panic(fmt.Sprintf("get supported chains error: %s", err.Error()))
	// }
	// printJSON("supported chains", supportedChains)

	// claimAddresses, err := client.GetAllClaimAddresses(ctx, auth.PartnerID)
	// if err != nil {
	// 	panic(fmt.Sprintf("get claim addresses error: %s", err.Error()))
	// }
	// printJSON("claim addresses", claimAddresses)

	// err = client.AttachNewClaimAddress(ctx, models.ClaimAddressPublic{
	// 	PartnerID: auth.PartnerID,
	// 	ChainID:   1,
	// 	Address:   "0x0000000000000000000000000000000000000000",
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("attach new claim address error: %s", err.Error()))
	// }

	// err = client.UpdateExistingClaimAddress(ctx, models.ClaimAddress{
	// 	PartnerID: auth.PartnerID,
	// 	ChainID:   1,
	// 	Address:   "0x0000000000000000000000000000000000000000",
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("update existing claim address error: %s", err.Error()))
	// }

}

func printJSON(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal %s error: %s", label, err.Error()))
	}
	fmt.Printf("\n %s \n%s\n", label, string(b))
}
