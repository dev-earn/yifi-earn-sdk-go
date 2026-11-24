package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/shopspring/decimal"
)

// Partner reads the partner's current configuration: policy and claim
// addresses. Both are read-only calls, safe to run as-is.
func Partner() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	policy, err := client.GetPartnerPolicy(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get partner policy: %s", err))
	}
	printJSON("partner policy", policy)

	claimAddresses, err := client.GetAllClaimAddresses(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get claim addresses: %s", err))
	}
	printJSON("claim addresses", claimAddresses)
}

// PartnerStats fetches the partner's aggregate summary and per-strategy
// metrics. An empty FromTo requests the default time range. Read-only.
func PartnerStats() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	summary, err := client.GetPartnerSummary(ctx, auth.PartnerID, models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get partner summary: %s", err))
	}
	printJSON("partner summary", summary)

	metrics, err := client.GetPartnerMetrics(ctx, auth.PartnerID, models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get partner metrics: %s", err))
	}
	printJSON("partner metrics", metrics)
}

// UpdatePolicy sets the partner's common fee to 0.3% and reads it back.
// MUTATING: this overwrites the live partner policy.
func UpdatePolicy() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	commonFee := decimal.NewFromFloat(0.3) // allowed range: 0 to 0.7
	if err := client.UpdatePartnerPolicy(ctx, auth.PartnerID, models.Policy{
		Fees: &models.FeePolicy{CommonFee: &commonFee},
	}); err != nil {
		panic(fmt.Sprintf("update partner policy: %s", err))
	}

	policy, err := client.GetPartnerPolicy(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get partner policy: %s", err))
	}
	printJSON("updated policy", policy)
}

// APIKeys lists the partner's keys, then creates and revokes a demo key.
// MUTATING: creating/revoking changes live partner credentials. The secret is
// only returned once, at creation.
func APIKeys() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	keys, err := client.GetAPIKeys(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get api keys: %s", err))
	}
	printJSON("api keys", keys)

	created, err := client.CreateAPIKey(ctx, auth.PartnerID, "demo-key")
	if err != nil {
		panic(fmt.Sprintf("create api key: %s", err))
	}
	printJSON("created api key", created)

	if err := client.RevokeAPIKey(ctx, auth.PartnerID, created.APIKey.APIKey); err != nil {
		panic(fmt.Sprintf("revoke api key: %s", err))
	}
	fmt.Println("demo api key revoked")
}

// SupportedChains lists the chains on which partner fees can be claimed.
// Read-only.
func SupportedChains() {
	client, _ := partnerTestClient()
	ctx := context.Background()

	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		panic(fmt.Sprintf("get supported chains: %s", err))
	}
	printJSON("supported chains", chains)
}

// AttachClaimAddress registers a new fee claim address on ethereum (chain 1).
// MUTATING.
func AttachClaimAddress() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	if err := client.AttachNewClaimAddress(ctx, models.ClaimAddressPublic{
		PartnerID: auth.PartnerID,
		ChainID:   1,
		Address:   walletAddr(),
		TOTPCode:  getenv("YIFI_TOTP_CODE", ""),
	}); err != nil {
		panic(fmt.Sprintf("attach claim address: %s", err))
	}
	fmt.Println("claim address attached")
}

// UpdateClaimAddress renames an existing claim address. MUTATING.
func UpdateClaimAddress() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	name := "main"
	if err := client.UpdateExistingClaimAddress(ctx, models.ClaimAddress{
		PartnerID: auth.PartnerID,
		ChainID:   1,
		Address:   walletAddr(),
		Name:      &name,
		TOTPCode:  getenv("YIFI_TOTP_CODE", ""),
	}); err != nil {
		panic(fmt.Sprintf("update claim address: %s", err))
	}
	fmt.Println("claim address updated")
}

// printJSON pretty-prints any value under a label; used across the examples.
func printJSON(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal %s: %s", label, err))
	}
	fmt.Printf("\n %s \n%s\n", label, string(b))
}
