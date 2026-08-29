package main

import (
	"context"
	"fmt"
)

// FeeBalance prints the partner's current fee balance across chains.
func FeeBalance() {
	client, auth := partnerTestClient()
	ctx := context.Background()

	bal, err := client.GetFeeBalance(ctx, auth.PartnerID)
	if err != nil {
		panic(fmt.Sprintf("get fee balance: %s", err))
	}
	printJSON("fee balance", bal)
}
