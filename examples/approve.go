package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Approve approves the Portals router to spend the wallet's USDT on polygon.
// Required before an ERC20-input zap: the router pulls the token via
// transferFrom, which reverts with TRANSFER_FROM_FAILED without an allowance.
func Approve() {
	ctx := context.Background()
	rpcURL := rpcURL("arbitrum")
	privateKey := privateKey()
	token := "0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9"   // USDT on arbitrum
	spender := "0x34b6a821d2f26c6b7cdb01Cd91895170C6574A0D" // Portals router (from the zap quote)

	max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	// approve(spender, maxUint256) selector 0x095ea7b3
	data := append(common.FromHex("0x095ea7b3"),
		append(common.LeftPadBytes(common.HexToAddress(spender).Bytes(), 32),
			common.LeftPadBytes(max.Bytes(), 32)...)...)

	hash, err := signAndSend(ctx, rpcURL, privateKey, models.Signature{
		To:    token,
		Data:  "0x" + common.Bytes2Hex(data),
		Value: "0",
	})
	if err != nil {
		panic(fmt.Sprintf("approve: %s", err))
	}
	fmt.Printf("approve sent: https://polygonscan.com/tx/%s (waiting to mine)\n", hash)

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	th := common.HexToHash(hash)
	for i := 0; i < 60; i++ {
		if r, err := client.TransactionReceipt(ctx, th); err == nil && r != nil {
			fmt.Println("approve mined, status:", r.Status)
			return
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("approve not mined in time")
}
