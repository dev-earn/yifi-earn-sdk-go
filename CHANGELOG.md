# Changelog

## v0.0.2

### Added
- `ClaimAllFees` — claim every pending token's fees on a chain in one call (`claimAll`).
- Fee balances now split into `free` / `pending` buckets (`TokenAmount`).
- Lending/Aave strategy history now exposes the full reserve snapshot via
  `Metrics.Meta` (`MetaLending`: apy/tvl, reserve metadata, IRM).
- Multi-chain RPC support in the examples via per-chain `RPC_URL_<CHAIN>` vars.

### Fixed
- `StrategyHistoryResponse` dropped the per-point `meta` object for Aave
  strategies; the model now round-trips the API response.

### Docs / internal
- Godoc comments across the `SDK` interface, client, options and key models.
- Rewrote README (auth, signing flow, IDs, testnet, examples).
- Examples moved to their own Go module so on-chain deps (go-ethereum) no longer
  leak into apps that import the SDK.
- Runnable, compiled example scenarios;
