package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go"
)

func Token() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	// GetToken
	// token, err := client.GetToken(context.Background(), "ethereum:0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	// if err != nil {
	// 	panic(fmt.Sprintf("get token error: %s", err.Error()))
	// }

	// b, err := json.MarshalIndent(token, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal token error: %s", err.Error()))
	// }

	// fmt.Println(string(b))

	// GetTokens
	// tokens, err := client.GetTokens(context.Background(), models.TokensRequest{
	// 	Chains: []string{"ethereum"},
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("get token error: %s", err.Error()))
	// }

	// bTokens, err := json.MarshalIndent(tokens, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal token error: %s", err.Error()))
	// }

	// fmt.Println(string(bTokens))

	// GetTokenHistoryMetrics
	// tokenHistory, err := client.GetTokenHistoryMetrics(context.Background(), "ethereum:0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", models.TokenHistoryMetricsRequest{})
	// if err != nil {
	// 	panic(fmt.Sprintf("get token history error: %s", err.Error()))
	// }
	// bTokenHistory, err := json.MarshalIndent(tokenHistory, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal token error: %s", err.Error()))
	// }

	// fmt.Println(string(bTokenHistory))

	// GetTokenGroup
	// tokenGroup, err := client.GetTokenGroup(context.Background(), "ethereum:0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	// if err != nil {
	// 	panic(fmt.Sprintf("get token group error: %s", err.Error()))
	// }

	// bTokenGroup, err := json.MarshalIndent(tokenGroup, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal token group error: %s", err.Error()))
	// }

	// fmt.Println(string(bTokenGroup))

	// GetTokenTags
	// tokenTags, err := client.GetTokenTags(context.Background())
	// if err != nil {
	// 	panic(fmt.Sprintf("get token tags error: %s", err.Error()))
	// }

	// bTokenTags, err := json.MarshalIndent(tokenTags, "", "  ")
	// if err != nil {
	// 	panic(fmt.Sprintf("marshal token group error: %s", err.Error()))
	// }

	// fmt.Println(string(bTokenTags))

	// GetTokenGroups
	tokenGroups, err := client.GetTokenGroups(context.Background())
	if err != nil {
		panic(fmt.Sprintf("get token tags error: %s", err.Error()))
	}

	bTokenGroups, err := json.MarshalIndent(tokenGroups, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal token group error: %s", err.Error()))
	}

	fmt.Println(string(bTokenGroups))

}
