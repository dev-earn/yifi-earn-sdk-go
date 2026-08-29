package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

// Token demonstrates every public token endpoint. All are read-only and need
// no partner auth, so the whole function is safe to run as-is.
func Token() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err))
	}
	ctx := context.Background()

	// WETH on ethereum, as a "chain:address" token ID.
	const wethID = "ethereum:0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"

	// GetToken: current state of a single token.
	token, err := client.GetToken(ctx, wethID)
	if err != nil {
		panic(fmt.Sprintf("get token: %s", err))
	}
	printJSON("token", token)

	// GetTokens: paginated, filterable list.
	tokens, err := client.GetTokens(ctx, models.TokensRequest{Chains: []string{"ethereum"}})
	if err != nil {
		panic(fmt.Sprintf("get tokens: %s", err))
	}
	printJSON("tokens", tokens)

	// GetTokenHistoryMetrics: historical stats over a time range (empty = default).
	history, err := client.GetTokenHistoryMetrics(ctx, wethID, models.FromTo{})
	if err != nil {
		panic(fmt.Sprintf("get token history: %s", err))
	}
	printJSON("token history", history)

	// GetTokenGroup: the group a token belongs to (e.g. bridged variants).
	group, err := client.GetTokenGroup(ctx, wethID)
	if err != nil {
		panic(fmt.Sprintf("get token group: %s", err))
	}
	printJSON("token group", group)

	// GetTokenTags: all available tags.
	tags, err := client.GetTokenTags(ctx)
	if err != nil {
		panic(fmt.Sprintf("get token tags: %s", err))
	}
	printJSON("token tags", tags)

	// GetTokenGroups: all groups.
	groups, err := client.GetTokenGroups(ctx)
	if err != nil {
		panic(fmt.Sprintf("get token groups: %s", err))
	}
	printJSON("token groups", groups)
}
