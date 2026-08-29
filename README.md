# yifi-earn-sdk-go

Go SDK for the yifi Earn API: public market data (chains, tokens, markets,
strategies), zap building, and authenticated partner operations (policies, API
keys, fee balances and claims).

## Install

```bash
go get github.com/dev-earn/yifi-earn-sdk-go
```

Requires Go 1.26+.

## Quick start (public data)

Public endpoints need no credentials.

```go
package main

import (
	"context"
	"fmt"

	yifi "github.com/dev-earn/yifi-earn-sdk-go"
)

func main() {
	client, err := yifi.NewClient()
	if err != nil {
		panic(err)
	}

	chains, err := client.GetChains(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", chains)
}
```

### Identifiers

Tokens, markets and strategies are addressed by a `chain:address` string, e.g.
`bsc:0xb0b84d294e0c75a6abe60171b70edeb2efd14a1b`. Native coins use the zero
address, e.g. `bsc:0x0000000000000000000000000000000000000000`.

### Testnet

```go
client, _ := yifi.NewClient()
```

## Authentication

Partner endpoints require credentials via `WithPartnerAuth`. Each request is
signed automatically (HMAC-SHA256 over the API key, a per-request nonce and the
body hash). Public endpoints keep working without auth.

```go
client, err := yifi.NewClient(yifi.WithPartnerAuth(models.PartnerAuth{
	PartnerID: "your_partner_id",
	APIKey:    "yifi_pk_live_...",
	Secret:    "...",
}))
```

Create and manage API keys via the frontend or support.

## Signing transactions

`BuildZap`, `ClaimFees` and `ClaimAllFees` return an **unsigned** transaction
(`models.Signature` / `models.TxInfo`), not a broadcast one. You sign it with
the sender's key and submit it to an RPC yourself. A few gotchas:

- `Signature.Value` is a **decimal wei string** (e.g. `"1000000000000000"`), not
  hex — parse it with base 10.
- `InputAmount` is also decimal wei in the token's smallest unit
  (e.g. `"1000000"` = 1 USDT at 6 decimals).
- For a native-coin input (zero address) the router wraps the coin using the
  transaction's `value`; you do not need the wrapped token. For an ERC-20 input
  you must `approve` the router first, or the zap reverts with
  `TRANSFER_FROM_FAILED`.

See [`examples/sign.go`](examples/sign.go) for a complete sign-and-send helper.

## Partner flow

1. Get an API key (frontend or support).
2. Initialize the client with `WithPartnerAuth`.
3. Set your fee via `UpdatePartnerPolicy` (default `0.3`, max `0.7`).
4. Build zaps with `BuildZap`; your fee is credited when the user signs.
5. Check accrued fees with `GetFeeBalance` — split into `free` (accrued) and
   `pending` (pushed on-chain, claimable) buckets.
6. Register a claim address with `AttachNewClaimAddress`.
7. Claim with `ClaimFees` (one token) or `ClaimAllFees` (all pending on a
   chain), then sign and broadcast the returned transaction from your claim
   address — it must be the signer.

```go
auth := models.PartnerAuth{PartnerID: "PARTNER_ID", APIKey: "API_KEY", Secret: "SECRET"}
client, err := yifi.NewClient(yifi.WithPartnerAuth(auth))
if err != nil {
	panic(err)
}
ctx := context.Background()

// Set the partner fee.
fee := decimal.NewFromFloat(0.5)
if err := client.UpdatePartnerPolicy(ctx, auth.PartnerID, models.Policy{
	Fees: &models.FeePolicy{CommonFee: &fee},
}); err != nil {
	panic(err)
}

// Build a zap for your user to sign.
zap, err := client.BuildZap(ctx, models.ZapQuoteRequest{
	StrategyID:                  "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
	InputToken:                  "ethereum:0x0000000000000000000000000000000000000000", // native ETH
	InputAmount:                 "10000000000000000",                                    // 0.01 ETH (wei)
	OutputToken:                 "ethereum:0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0",
	Sender:                      "0xYourUserWallet", // the wallet that signs
	SlippageTolerancePercentage: decimal.NewFromInt(1),
})
if err != nil {
	panic(err)
}

// Check accrued fees, then claim.
balance, err := client.GetFeeBalance(ctx, auth.PartnerID)
if err != nil {
	panic(err)
}
_ = balance
_ = zap
```

## Examples

Runnable examples live in [`examples/`](examples/), which is a **separate Go
module** so its on-chain dependencies (go-ethereum, godotenv) don't leak into
apps that import the SDK.

```bash
cd examples
cp .env.example .env   # fill in your credentials / RPC / signer key
```

`examples/main.go` is a dispatcher — uncomment one scenario and run it:

```bash
go run .
```

Read-only scenarios (`Token`, `Market`, `Strategy`, `Partner`, `FeeBalance`,
`SupportedChains`) are safe to run. Scenarios marked **MUTATING** in their doc
comment (`UpdatePolicy`, `APIKeys`, `Zap`, `Claim`, `ClaimAll`,
`AttachClaimAddress`) change live state or broadcast transactions — read them
before running.

### RPC configuration (multi-chain)

The on-chain examples are EVM-only and work across any number of chains. Add one
`RPC_URL_<CHAIN>` per chain you operate on to `examples/.env`; the `<CHAIN>`
suffix is the SDK's `chain:address` prefix, upper-cased:

```bash
RPC_URL_ETHEREUM=https://eth-mainnet.your-rpc-provider/v2/your-key
RPC_URL_ARBITRUM=https://arb-mainnet.your-rpc-provider/v2/your-key
RPC_URL_BSC=https://bnb-mainnet.your-rpc-provider/v2/your-key
# ...add optimism, base, polygon, etc. as needed

RPC_URL=https://your-rpc-provider/v2/your-key   # optional catch-all fallback
```

Each example resolves its endpoint by chain (e.g. an Arbitrum zap reads
`RPC_URL_ARBITRUM`), falling back to `RPC_URL`. If a chain an example needs has
no endpoint configured, it fails fast with a message telling you which
`RPC_URL_<CHAIN>` to set.

## API reference

Full docs on [pkg.go.dev](https://pkg.go.dev/github.com/dev-earn/yifi-earn-sdk-go).
The complete surface is the `SDK` interface in [`sdk.go`](sdk.go).
```
