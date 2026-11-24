package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func signAndSend(ctx context.Context, rpcURL, privateKeyHex string, sig models.Signature) (string, error) {
	ethClient, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return "", fmt.Errorf("dial rpc: %w", err)
	}
	defer ethClient.Close()

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	fromAddr := crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey))
	toAddr := common.HexToAddress(sig.To)

	value := new(big.Int)
	if sig.Value != "" && sig.Value != "0" {
		if hexVal := strings.TrimPrefix(sig.Value, "0x"); len(hexVal) < len(sig.Value) {
			value.SetString(hexVal, 16) // "0x"-prefixed => hex
		} else {
			value.SetString(sig.Value, 10) // bare decimal wei string
		}
	}

	data := common.FromHex(sig.Data)

	nonce, err := ethClient.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("get nonce: %w", err)
	}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("get chain id: %w", err)
	}

	tip, err := ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("suggest tip cap: %w", err)
	}

	head, err := ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("get latest header: %w", err)
	}

	// BSC can have a zero base fee with a nonzero priority tip.
	maxFee := new(big.Int).Set(head.BaseFee)
	if maxFee.Cmp(tip) < 0 {
		maxFee = new(big.Int).Set(tip)
	}

	var gasLimit uint64
	if sig.GasLimit != "" {
		limit := new(big.Int)
		base := 10
		raw := sig.GasLimit
		if strings.HasPrefix(raw, "0x") {
			base = 16
			raw = strings.TrimPrefix(raw, "0x")
		}
		if _, ok := limit.SetString(raw, base); !ok {
			return "", fmt.Errorf("invalid engine gas limit: %q", sig.GasLimit)
		}
		gasLimit = limit.Uint64()
	} else {
		gasLimit, err = ethClient.EstimateGas(ctx, ethereum.CallMsg{
			From:     fromAddr,
			To:       &toAddr,
			Value:    value,
			Data:     data,
			GasPrice: maxFee,
		})
		if err != nil {
			var dataErr rpc.DataError
			if errors.As(err, &dataErr) {
				return "", fmt.Errorf("estimate gas: %w (data: %v)", err, dataErr.ErrorData())
			}
			return "", fmt.Errorf("estimate gas: %w", err)
		}
	}
	if sig.GasLimit == "" {
		// +20% buffer when estimating locally.
		gasLimit = gasLimit * 12 / 10
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: maxFee,
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     value,
		Data:      data,
	})

	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("sign tx: %w", err)
	}

	if err := ethClient.SendTransaction(ctx, signed); err != nil {
		return "", fmt.Errorf("send tx: %w", err)
	}

	return signed.Hash().Hex(), nil
}
