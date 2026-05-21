package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dev-earn/yifi-earn-sdk-go"
)

func Chain() {
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
