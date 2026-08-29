// Package yifi is the Go client for the yifi Earn API.
//
// It exposes read-only market data (chains, tokens, markets, strategies) and
// authenticated partner operations (policies, API keys, fee balances and
// claims, zap building). All calls go through the SDK interface, implemented
// by Client.
//
// # Getting started
//
//	client, err := yifi.NewClient(
//	    yifi.WithPartnerAuth(models.PartnerAuth{
//	        PartnerID: "your_partner_id",
//	        APIKey:    "yifi_pk_live_...",
//	        Secret:    "...",
//	    }),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	chains, err := client.GetChains(context.Background())
//
// # Authentication
//
// Public market-data endpoints need no credentials. Partner endpoints require
// WithPartnerAuth: every request is signed with an HMAC-SHA256 of the API key,
// a per-request nonce and the body hash, so a NewClient without partner auth
// will be rejected by those endpoints.
//
// # Testnet
//
// Set UseTestnet = true before calling NewClient to target the testnet API.
//
// # Identifiers
//
// Tokens, markets and strategies are addressed by a "chain:address" string,
// e.g. "bsc:0xb0b84d294e0c75a6abe60171b70edeb2efd14a1b". Native assets use the
// zero address, e.g. "bsc:0x0000000000000000000000000000000000000000".
//
// Runnable examples live in the examples/ directory (a separate module).
package yifi
