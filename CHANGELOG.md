# Changelog

## v0.0.3

### Added
- Configurable base URL support.
- Partner transaction reporting, strategy performance, and full-cycle examples.
- Quote routing metadata, gas estimates, approval requirements, complete wallet transaction metadata, and quote `expires_in` timestamps.
- Partner transaction filters and token/strategy search options.
- Public read-only SDK verification runner (`YIFI_EXAMPLE=public`).

### Changed
- Synchronized policy, supported-chain, strategy, claim, and transaction models with the gateway contracts.
- `GetTxInfo` and `GetWalletTransaction` return `WalletTxResponse` with optional execution fields.
- Claim requests now use the current partner API fields and return transaction value.
- LI.FI status uses the upstream transaction ID directly; `RegisterTxHash` was removed from the SDK.

### Fixed
- `UpdateExistingClaimAddress` now puts the numeric chain ID in the route instead of the wallet address.
- Mutation endpoints accept empty success responses without JSON decoding errors.
- The partner example waits for webhook settlement and verifies the final zero fee balance.

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
