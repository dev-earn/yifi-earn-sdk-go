package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// init loads a local .env file if present so the examples can be configured
// without exporting variables in the shell. A missing file is not an error —
// values may come from the real environment instead.
func init() {
	_ = godotenv.Load()
}

// getenv returns the value of the environment variable named by key, falling
// back to placeholder when it is unset. The placeholders make the examples
// compile out of the box; supply real values via a .env file (see
// .env.example) or your shell environment before running anything that hits
// the network or signs a transaction.
func getenv(key, placeholder string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return placeholder
}

// Common configuration pulled from the environment. Never commit real keys —
// see .env.example for the expected variables.
func partnerID() string  { return getenv("YIFI_PARTNER_ID", "your_partner_id") }
func apiKey() string     { return getenv("YIFI_API_KEY", "yifi_pk_live_YOUR_API_KEY") }
func apiSecret() string  { return getenv("YIFI_API_SECRET", "YOUR_API_SECRET") }
func privateKey() string { return getenv("PRIVATE_KEY", "YOUR_WALLET_PRIVATE_KEY") }
func walletAddr() string { return getenv("WALLET_ADDRESS", "0xYOUR_WALLET_ADDRESS") }

// rpcURLs collects one JSON-RPC endpoint per chain from every RPC_URL_<CHAIN>
// environment variable (e.g. RPC_URL_ARBITRUM, RPC_URL_BSC). All supported
// chains are EVM, so operate on as many as you like — just add one line per
// chain to your .env. The map key is the lower-cased chain name, matching the
// "chain:address" prefix used throughout the SDK.
func rpcURLs() map[string]string {
	m := make(map[string]string)
	for _, kv := range os.Environ() {
		key, val, ok := strings.Cut(kv, "=")
		if !ok || val == "" {
			continue
		}
		if chain, found := strings.CutPrefix(key, "RPC_URL_"); found {
			m[strings.ToLower(chain)] = val
		}
	}
	return m
}

// rpcURL returns the JSON-RPC endpoint for the given chain (e.g. "arbitrum"),
// from its RPC_URL_<CHAIN> entry, falling back to a generic RPC_URL. It panics
// with an actionable message if neither is set, so a misconfigured chain fails
// loudly instead of dialing a placeholder host.
func rpcURL(chain string) string {
	if u := rpcURLs()[strings.ToLower(chain)]; u != "" {
		return u
	}
	if u := os.Getenv("RPC_URL"); u != "" {
		return u
	}
	panic(fmt.Sprintf("no RPC endpoint for chain %q: set RPC_URL_%s (or a generic RPC_URL) in examples/.env",
		chain, strings.ToUpper(chain)))
}
