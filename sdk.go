package yifi

import (
	"context"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

type SDK interface {
	GetChains(ctx context.Context) ([]models.Chain, error)
	GetPlatforms(ctx context.Context) ([]models.Platform, error)
	GetCategories(ctx context.Context) ([]models.Category, error)

	GetTokens(ctx context.Context, req models.TokensRequest) (*models.Pagination[models.GetTokensResponse], error)
	GetToken(ctx context.Context, tokenID string) (models.TokenFullCurrent, error)
	GetTokenHistoryMetrics(ctx context.Context, tokenID string, req models.FromTo) (models.TokenStatsHistory, error)
	GetTokenGroup(ctx context.Context, tokenID string) (models.TokenGroup, error)
	GetTokenTags(ctx context.Context) ([]models.TokenTag, error)
	GetTokenGroups(ctx context.Context) ([]models.TokenGroup, error)

	GetMarkets(ctx context.Context, request models.MarketsRequest) (*models.Pagination[models.MarketsResponse], error)
	GetMarket(ctx context.Context, marketID string) (models.MarketCurrentResponse, error)
	GetMarketHistoryMetrics(ctx context.Context, marketID string, req models.FromTo) (models.MarketHistoryResponse, error)
	GetMarketTokenRoles(ctx context.Context) ([]models.MarketTokenRoles, error)

	GetStrategies(ctx context.Context, req models.StrategySearch) (*models.Pagination[models.StrategiesResponse], error)
	GetStrategy(ctx context.Context, strategyID string) (models.StrategyResponse, error)
	GetStrategyHistory(ctx context.Context, strategyID string, req models.FromTo) (models.StrategyHistoryResponse, error)
	GetEngines(ctx context.Context) ([]models.Engine, error)
	GetRiskScores(ctx context.Context) ([]models.Risk, error)
	GetStrategyStatuses(ctx context.Context) ([]models.Status, error)
	GetYields(ctx context.Context) ([]models.TypeYield, error)
	GetRates(ctx context.Context) ([]models.RateFormat, error)

	BuildZap(ctx context.Context, quote models.ZapQuoteRequest) (models.Signature, error)
	GetWalletTransactions(ctx context.Context, walletID string) ([]models.WalletTx, error)

	GetPartnerSummary(ctx context.Context, partnerID string, req models.FromTo) (models.Summary, error)
	GetPartnerMetrics(ctx context.Context, partnerID string, req models.FromTo) (models.MetricsPartner, error)
	GetPartnerPolicy(ctx context.Context, partnerID string) (models.Policy, error)
	UpdatePartnerPolicy(ctx context.Context, partnerID string, policy models.Policy) error
	GetAPIKeys(ctx context.Context, partnerID string) ([]models.APIKey, error)
	CreateAPIKey(ctx context.Context, partnerID, name string) (models.APIKeyFull, error)
	RevokeAPIKey(ctx context.Context, partnerID, apiKey string) error
	GetFeeBalance(ctx context.Context, partnerID string) (models.Balance, error)
	ClaimFees(ctx context.Context, reqBody models.PartnerFeeClaimRequest) (models.TxInfo, error)
	GetSupportedChains(ctx context.Context) ([]models.Chain, error)
	GetAllClaimAddresses(ctx context.Context, partnerID string) ([]models.ClaimAddress, error)
	AttachNewClaimAddress(ctx context.Context, reqBody models.ClaimAddressPublic) error
	UpdateExistingClaimAddress(ctx context.Context, reqBody models.ClaimAddress) error
}
