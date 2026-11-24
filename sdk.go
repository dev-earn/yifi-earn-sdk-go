package yifi

import (
	"context"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

// SDK is the full surface of the yifi Earn API. Client is the concrete
// implementation returned by NewClient. Methods split into public market-data
// calls (no auth) and partner calls (require WithPartnerAuth).
//
// Every method takes a context.Context for cancellation and deadlines. IDs are
// "chain:address" strings (see the package overview).
type SDK interface {
	// --- Reference data (public) ---

	// GetChains lists every chain supported by the API.
	GetChains(ctx context.Context) ([]models.Chain, error)
	// GetPlatforms lists the integrated DeFi platforms (protocols).
	GetPlatforms(ctx context.Context) ([]models.Platform, error)
	// GetCategories lists the strategy/market categories.
	GetCategories(ctx context.Context) ([]models.Category, error)

	// --- Tokens (public) ---

	// GetTokens returns a paginated, filterable list of tokens.
	GetTokens(ctx context.Context, req models.TokensRequest) (*models.Pagination[models.GetTokensResponse], error)
	// GetToken returns the current state of a single token by "chain:address" ID.
	GetToken(ctx context.Context, tokenID string) (models.TokenFullCurrent, error)
	// GetTokenHistoryMetrics returns historical metrics for a token over a time range.
	GetTokenHistoryMetrics(ctx context.Context, tokenID string, req models.FromTo) (models.TokenStatsHistory, error)
	// GetTokenGroup returns the group a token belongs to (e.g. bridged variants of one asset).
	GetTokenGroup(ctx context.Context, tokenID string) (models.TokenGroup, error)
	// GetTokenTags lists the available token tags.
	GetTokenTags(ctx context.Context) ([]models.TokenTag, error)
	// GetTokenGroups lists all token groups.
	GetTokenGroups(ctx context.Context) ([]models.TokenGroup, error)

	// --- Markets (public) ---

	// GetMarkets returns a paginated, filterable list of markets.
	GetMarkets(ctx context.Context, request models.MarketsRequest) (*models.Pagination[models.MarketsResponse], error)
	// GetMarket returns the current state of a single market by ID.
	GetMarket(ctx context.Context, marketID string) (models.MarketCurrentResponse, error)
	// GetMarketHistoryMetrics returns historical metrics for a market over a time range.
	GetMarketHistoryMetrics(ctx context.Context, marketID string, req models.FromTo) (models.MarketHistoryResponse, error)
	// GetMarketTokenRoles lists the roles a token can play within a market (e.g. supply, borrow).
	GetMarketTokenRoles(ctx context.Context) ([]models.MarketTokenRoles, error)

	// --- Strategies (public) ---

	// GetStrategies returns a paginated, filterable list of strategies.
	GetStrategies(ctx context.Context, req models.StrategySearch) (*models.Pagination[models.StrategiesResponse], error)
	// GetStrategy returns the current state of a single strategy by ID.
	GetStrategy(ctx context.Context, strategyID string) (models.StrategyResponse, error)
	// GetStrategyHistory returns historical metrics for a strategy over a time range.
	GetStrategyHistory(ctx context.Context, strategyID string, req models.FromTo) (models.StrategyHistoryResponse, error)
	// GetEngines lists the strategy engines (yield mechanisms).
	GetEngines(ctx context.Context) ([]models.Engine, error)
	// GetRiskScores lists the available risk scores.
	GetRiskScores(ctx context.Context) ([]models.Risk, error)
	// GetStrategyStatuses lists the possible strategy statuses.
	GetStrategyStatuses(ctx context.Context) ([]models.Status, error)
	// GetYields lists the yield types used to classify strategies.
	GetYields(ctx context.Context) ([]models.TypeYield, error)
	// GetRates lists the supported rate formats (e.g. APR, APY).
	GetRates(ctx context.Context) ([]models.RateFormat, error)

	// --- Zap & wallet ---

	// QuoteZap freezes a quote and its transaction for one minute.
	QuoteZap(ctx context.Context, quote models.ZapQuoteRequest) (models.ZapQuoteResponse, error)
	// BuildZap returns the frozen unsigned transaction for quote.OrderID.
	BuildZap(ctx context.Context, quote models.ZapQuoteRequest) (models.Signature, error)
	// GetTxInfo returns execution information by the quote/build order id.
	GetTxInfo(ctx context.Context, orderID string) (models.WalletTxResponse, error)
	// GetWalletTransactions returns the yifi-tracked transactions for a wallet address.
	GetWalletTransactions(ctx context.Context, walletID string) ([]models.WalletTx, error)
	// GetWalletTransaction returns one yifi-tracked transaction by wallet and order id.
	GetWalletTransaction(ctx context.Context, walletID, orderID string) (models.WalletTxResponse, error)

	// --- Partner (require WithPartnerAuth) ---

	// GetPartnerSummary returns aggregate partner stats over a time range.
	GetPartnerSummary(ctx context.Context, partnerID string, req models.FromTo) (models.Summary, error)
	// GetPartnerMetrics returns per-strategy partner metrics over a time range.
	GetPartnerMetrics(ctx context.Context, partnerID string, req models.FromTo) (models.MetricsPartner, error)
	// GetPartnerTransactions returns attributed user transactions, optionally filtered by hash.
	GetPartnerTransactions(ctx context.Context, partnerID string, req models.PartnerTransactionsRequest) (models.PartnerTransactionsResponse, error)
	// GetStrategyPerformance returns attributed performance grouped by strategy.
	GetStrategyPerformance(ctx context.Context, partnerID string, req models.FromTo) (models.StrategyPerformanceResponse, error)
	// GetPartnerPolicy returns the partner's current policy (fees, allow/deny lists, limits).
	GetPartnerPolicy(ctx context.Context, partnerID string) (models.Policy, error)
	// UpdatePartnerPolicy overwrites the partner's policy. Only non-nil fields are applied.
	UpdatePartnerPolicy(ctx context.Context, partnerID string, policy models.Policy) error
	// GetAPIKeys lists the partner's API keys (secrets are not returned).
	GetAPIKeys(ctx context.Context, partnerID string) ([]models.APIKey, error)
	// CreateAPIKey creates a named API key and returns it, including the secret (shown once).
	CreateAPIKey(ctx context.Context, partnerID, name string) (models.APIKeyFull, error)
	// RevokeAPIKey permanently disables an API key.
	RevokeAPIKey(ctx context.Context, partnerID, apiKey string) error
	// GetFeeBalance returns the partner's accrued fees, split into free and pending buckets.
	GetFeeBalance(ctx context.Context, partnerID string) (models.Balance, error)
	// ClaimFees pushes one token's free balance to pending and returns the unsigned
	// claim() calldata to sign from the partner's claim address.
	ClaimFees(ctx context.Context, reqBody models.PartnerFeeClaimRequest) (models.TxInfo, error)
	// ClaimAllFees does the same as ClaimFees for every pending token on a chain (claimAll).
	ClaimAllFees(ctx context.Context, reqBody models.PartnerFeeClaimAllRequest) (models.TxInfo, error)
	// GetSupportedChains lists the chains on which partner fees can be claimed.
	GetSupportedChains(ctx context.Context) ([]models.PartnerChain, error)
	// GetAllClaimAddresses lists the partner's configured fee claim addresses.
	GetAllClaimAddresses(ctx context.Context, partnerID string) ([]models.ClaimAddress, error)
	// AttachNewClaimAddress registers a new claim address for a chain.
	AttachNewClaimAddress(ctx context.Context, reqBody models.ClaimAddressPublic) error
	// UpdateExistingClaimAddress updates an existing claim address (e.g. its name).
	UpdateExistingClaimAddress(ctx context.Context, reqBody models.ClaimAddress) error
}
