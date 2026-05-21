# yifi-earn-sdk-go

Go SDK for the YIFI gateway API, including public core, zap and partner-related endpoints.

## Install

```bash
go get github.com/dev-earn/yifi-earn-sdk-go
```

## Quick start

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

func main() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	chains, err := client.GetChains(context.Background())
	if err != nil {
		panic(fmt.Sprintf("get chains error: %s", err.Error()))
	}

	b, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal chains error: %s", err.Error()))
	}

	fmt.Println(string(b))
}
```

## Partner flow

If you are a partner, the recommended flow is:

1. Generate an API key via the frontend or by contacting support.
2. Initialize the SDK client with your `partnerID`, `apiKey`, and `secret`.
3. Update your partner policy, especially the fee percentage.
4. Build zaps using the partner-authenticated client.
5. Check your fee balance either via the API or the frontend.
6. Attach your claim address to the allowed chains.
7. Claim your accumulated fees to that address.

### Important note about fees

By default, the partner fee is set to `0.3`. It can be increased up to `0.7`.

Once a zap is built and the generated instruction is signed by the user, your fee share is credited to your fee balance.

## Partner example

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

func Partners() {
	auth := models.PartnerAuth{
		PartnerID: "PARTNER_ID",
		APIKey:    "API_KEY",
		Secret:    "SECRET",
	}

	client, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	// 1. Update partner policy
    partnerFeePref := decimal.NewFromFloat(0.7)
	err = client.UpdatePartnerPolicy(ctx, auth.PartnerID, models.Policy{
		Fees: &models.FeePolicy{
			CommonFee: &partnerFeePref,
		},
	})

	policy, err := client.GetPartnerPolicy(context.Background(), auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get partner policy error: %s", err.Error()))
	}
	printJSON("current policy", policy)

	// 2. Build zaps as a partner
	//    The partner headers are attached automatically when partner auth is configured.
	buildResp, err := client.BuildZap(context.Background(), models.ZapQuoteRequest{
		StrategyID:                  "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
		InputToken:                  "ethereum:0x0000000000000000000000000000000000000000",
		InputAmount:                 "10000000000",
		OutputToken:                 "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
		Sender:                      "0x053806E0642F6A6aAFE6176bd264cbbE446C8c92", // your user, who signs
		SlippageTolerancePercentage: decimal.NewFromFloat(1),
	})
	printJSON("build zap", buildResp)

	// 3. Check fee balance
	balance, err := client.GetFeeBalance(context.Background(), auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get fee balance error: %s", err.Error()))
	}
	printJSON("fee balance", balance)

	// 4. Attach claim address
	err = client.AttachNewClaimAddress(context.Background(), models.ClaimAddressPublic{
		PartnerID: auth.PartnerID,
		ChainID:   1,
		Address:   "0x0000000000000000000000000000000000000000",
	})
	if err != nil {
		panic(fmt.Sprintf("attach new claim address error: %s", err.Error()))
	}

	// 5. Claim fees
	txInfo, err := client.ClaimFees(context.Background(), models.PartnerFeeClaimRequest{
		PartnerID: auth.PartnerID,
	})
	if err != nil {
		panic(fmt.Sprintf("claim fees error: %s", err.Error()))
	}
	printJSON("claim fees", txInfo)
}

func printJSON(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal %s error: %s", label, err.Error()))
	}
	fmt.Printf("\n %s \n%s\n", label, string(b))
}
```

## Notes

* Partner authentication is optional.
* If partner auth is configured, the SDK signs partner requests automatically.
* If partner auth is not configured, public endpoints continue to work normally.
* Use the frontend or support to create and manage your first partner API key.
