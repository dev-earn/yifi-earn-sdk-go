package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go"
)

func Platform() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(fmt.Sprintf("yifi new client init failed: %s", err.Error()))
	}

	chains, err := client.GetPlatforms(context.Background())
	if err != nil {
		panic(fmt.Sprintf("get platforms error: %s", err.Error()))

	}

	b, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("marshal platforms error: %s", err.Error()))
	}

	fmt.Println(string(b))
}
