package yifi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

var _ SDK = (*Client)(nil)

// Option configures a Client in NewClient.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client (timeouts, transport, proxy).
// The default has a 30s timeout. A nil client is ignored.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithBaseURL overrides the API host
func WithBaseURL(raw string) Option {
	return func(c *Client) {
		u, err := url.Parse(strings.TrimSpace(raw))
		if err == nil && u.Scheme != "" && u.Host != "" {
			c.baseURL = u
		}
	}
}

// WithBasePath overrides the API base path (default "/api/v1"). Rarely needed.
func WithBasePath(basePath string) Option {
	return func(c *Client) {
		c.basePath = normalizeBasePath(basePath)
	}
}

// WithPartnerAuth attaches partner credentials so the client can call the
// authenticated partner endpoints. Without it, those calls are rejected.
func WithPartnerAuth(auth models.PartnerAuth) Option {
	return func(c *Client) {
		c.partnerAuth = &models.PartnerAuth{
			PartnerID: auth.PartnerID,
			APIKey:    auth.APIKey,
			Secret:    auth.Secret,
		}
	}
}

// Client is the concrete SDK implementation. Create it with NewClient; it is
// safe for concurrent use.
type Client struct {
	baseURL     *url.URL
	basePath    string
	httpClient  *http.Client
	partnerAuth *models.PartnerAuth
}

// NewClient builds a Client. It reads UseTestnet to pick the API host, applies
// the given options, and returns an error if the resolved endpoint is invalid.
func NewClient(opts ...Option) (SDK, error) {
	u, err := url.Parse(strings.TrimSpace(getAPIEndpoint()))
	if err != nil {
		return nil, fmt.Errorf("parse api endpoint: %w", err)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid api endpoint: %q", getAPIEndpoint())
	}

	c := &Client{
		baseURL:    u,
		basePath:   "/api/v1",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.basePath = normalizeBasePath(c.basePath)

	return c, nil
}
func (c *Client) GetChains(ctx context.Context) ([]models.Chain, error) {
	req, err := c.do(ctx, http.MethodGet, "/chains", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Chain](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetPlatforms(ctx context.Context) ([]models.Platform, error) {
	req, err := c.do(ctx, http.MethodGet, "/platforms", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Platform](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetCategories(ctx context.Context) ([]models.Category, error) {
	req, err := c.do(ctx, http.MethodGet, "/categories", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Category](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTokens(ctx context.Context, req models.TokensRequest) (*models.Pagination[models.GetTokensResponse], error) {
	q := encodeTokensRequest(req)

	r, err := c.do(ctx, http.MethodGet, "/tokens", q, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[models.Pagination[models.GetTokensResponse]](c.httpClient, r)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (c *Client) GetToken(ctx context.Context, tokenID string) (models.TokenFullCurrent, error) {
	req, err := c.do(ctx, http.MethodGet, "/tokens/"+url.PathEscape(tokenID), nil, nil)
	if err != nil {
		return models.TokenFullCurrent{}, err
	}

	res, err := DoReq[models.TokenFullCurrent](c.httpClient, req)
	if err != nil {
		return models.TokenFullCurrent{}, err
	}

	if res.Data == nil {
		return models.TokenFullCurrent{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTokenHistoryMetrics(ctx context.Context, tokenID string, req models.FromTo) (models.TokenStatsHistory, error) {
	q := encodeFromTo(req)

	httpReq, err := c.do(
		ctx,
		http.MethodGet,
		"/tokens/"+url.PathEscape(tokenID)+"/metrics",
		q,
		nil,
	)
	if err != nil {
		return models.TokenStatsHistory{}, err
	}

	res, err := DoReq[models.TokenStatsHistory](c.httpClient, httpReq)
	if err != nil {
		return models.TokenStatsHistory{}, err
	}

	if res.Data == nil {
		return models.TokenStatsHistory{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTokenGroup(ctx context.Context, tokenID string) (models.TokenGroup, error) {
	req, err := c.do(ctx, http.MethodGet, "/tokens/"+url.PathEscape(tokenID)+"/group", nil, nil)
	if err != nil {
		return models.TokenGroup{}, err
	}

	res, err := DoReq[models.TokenGroup](c.httpClient, req)
	if err != nil {
		return models.TokenGroup{}, err
	}

	if res.Data == nil {
		return models.TokenGroup{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTokenTags(ctx context.Context) ([]models.TokenTag, error) {
	req, err := c.do(ctx, http.MethodGet, "/tags", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.TokenTag](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTokenGroups(ctx context.Context) ([]models.TokenGroup, error) {
	req, err := c.do(ctx, http.MethodGet, "/groups", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.TokenGroup](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetMarkets(ctx context.Context, request models.MarketsRequest) (*models.Pagination[models.MarketsResponse], error) {
	q := encodeMarketsRequest(request)

	req, err := c.do(ctx, http.MethodGet, "/markets", q, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[models.Pagination[models.MarketsResponse]](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (c *Client) GetMarket(ctx context.Context, marketID string) (models.MarketCurrentResponse, error) {
	req, err := c.do(ctx, http.MethodGet, "/markets/"+url.PathEscape(marketID), nil, nil)
	if err != nil {
		return models.MarketCurrentResponse{}, err
	}

	res, err := DoReq[models.MarketCurrentResponse](c.httpClient, req)
	if err != nil {
		return models.MarketCurrentResponse{}, err
	}

	if res.Data == nil {
		return models.MarketCurrentResponse{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetMarketHistoryMetrics(ctx context.Context, marketID string, req models.FromTo) (models.MarketHistoryResponse, error) {
	q := encodeFromTo(req)

	httpReq, err := c.do(ctx, http.MethodGet, "/markets/"+url.PathEscape(marketID)+"/metrics", q, nil)
	if err != nil {
		return models.MarketHistoryResponse{}, err
	}

	res, err := DoReq[models.MarketHistoryResponse](c.httpClient, httpReq)
	if err != nil {
		return models.MarketHistoryResponse{}, err
	}

	if res.Data == nil {
		return models.MarketHistoryResponse{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetMarketTokenRoles(ctx context.Context) ([]models.MarketTokenRoles, error) {
	req, err := c.do(ctx, http.MethodGet, "/roles", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.MarketTokenRoles](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetStrategies(ctx context.Context, req models.StrategySearch) (*models.Pagination[models.StrategiesResponse], error) {
	q := encodeStrategiesRequest(req)

	httpReq, err := c.do(ctx, http.MethodGet, "/strategies", q, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[models.Pagination[models.StrategiesResponse]](c.httpClient, httpReq)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (c *Client) GetStrategy(ctx context.Context, strategyID string) (models.StrategyResponse, error) {
	req, err := c.do(ctx, http.MethodGet, "/strategies/"+url.PathEscape(strategyID), nil, nil)
	if err != nil {
		return models.StrategyResponse{}, err
	}

	res, err := DoReq[models.StrategyResponse](c.httpClient, req)
	if err != nil {
		return models.StrategyResponse{}, err
	}

	if res.Data == nil {
		return models.StrategyResponse{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetStrategyHistory(ctx context.Context, strategyID string, req models.FromTo) (models.StrategyHistoryResponse, error) {
	q := encodeFromTo(req)

	httpReq, err := c.do(ctx, http.MethodGet, "/strategies/"+url.PathEscape(strategyID)+"/history", q, nil)
	if err != nil {
		return models.StrategyHistoryResponse{}, err
	}

	res, err := DoReq[models.StrategyHistoryResponse](c.httpClient, httpReq)
	if err != nil {
		return models.StrategyHistoryResponse{}, err
	}

	if res.Data == nil {
		return models.StrategyHistoryResponse{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetEngines(ctx context.Context) ([]models.Engine, error) {
	req, err := c.do(ctx, http.MethodGet, "/engines", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Engine](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetRiskScores(ctx context.Context) ([]models.Risk, error) {
	req, err := c.do(ctx, http.MethodGet, "/risks", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Risk](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetStrategyStatuses(ctx context.Context) ([]models.Status, error) {
	req, err := c.do(ctx, http.MethodGet, "/statuses", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.Status](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetYields(ctx context.Context) ([]models.TypeYield, error) {
	req, err := c.do(ctx, http.MethodGet, "/yields", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.TypeYield](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetRates(ctx context.Context) ([]models.RateFormat, error) {
	req, err := c.do(ctx, http.MethodGet, "/rates", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.RateFormat](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) QuoteZap(ctx context.Context, quote models.ZapQuoteRequest) (models.ZapQuoteResponse, error) {
	req, err := c.do(ctx, http.MethodPost, "/zap/quote", nil, quote)
	if err != nil {
		return models.ZapQuoteResponse{}, err
	}
	res, err := DoReq[models.ZapQuoteResponse](c.httpClient, req)
	if err != nil {
		return models.ZapQuoteResponse{}, err
	}
	if res.Data == nil {
		return models.ZapQuoteResponse{}, nil
	}
	return *res.Data, nil
}

func (c *Client) BuildZap(ctx context.Context, quote models.ZapQuoteRequest) (models.Signature, error) {
	req, err := c.do(ctx, http.MethodPost, "/zap/build", nil, quote)
	if err != nil {
		return models.Signature{}, err
	}

	res, err := DoReq[models.Signature](c.httpClient, req)
	if err != nil {
		return models.Signature{}, err
	}

	if res.Data == nil {
		return models.Signature{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetTxInfo(ctx context.Context, orderID string) (models.WalletTxResponse, error) {
	req, err := c.do(ctx, http.MethodGet, "/tx/"+url.PathEscape(orderID), nil, nil)
	if err != nil {
		return models.WalletTxResponse{}, err
	}
	res, err := DoReq[models.WalletTxResponse](c.httpClient, req)
	if err != nil {
		return models.WalletTxResponse{}, err
	}
	if res.Data == nil {
		return models.WalletTxResponse{}, nil
	}
	return *res.Data, nil
}

func (c *Client) GetWalletTransaction(ctx context.Context, walletID, orderID string) (models.WalletTxResponse, error) {
	path := "/transactions/" + url.PathEscape(walletID) + "/stats/" + url.PathEscape(orderID)
	req, err := c.do(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return models.WalletTxResponse{}, err
	}
	res, err := DoReq[models.WalletTxResponse](c.httpClient, req)
	if err != nil {
		return models.WalletTxResponse{}, err
	}
	if res.Data == nil {
		return models.WalletTxResponse{}, nil
	}
	return *res.Data, nil
}

func (c *Client) GetWalletTransactions(ctx context.Context, walletID string) ([]models.WalletTx, error) {
	req, err := c.do(ctx, http.MethodGet, "/transactions/"+url.PathEscape(walletID)+"/stats", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.WalletTx](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

// Partner

func (c *Client) GetPartnerSummary(ctx context.Context, partnerID string, request models.FromTo) (models.Summary, error) {
	q := encodeFromTo(request)

	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/summary", q, nil)
	if err != nil {
		return models.Summary{}, err
	}

	res, err := DoReq[models.Summary](c.httpClient, req)
	if err != nil {
		return models.Summary{}, err
	}

	if res.Data == nil {
		return models.Summary{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetPartnerTransactions(ctx context.Context, partnerID string, request models.PartnerTransactionsRequest) (models.PartnerTransactionsResponse, error) {
	q := encodeFromTo(request.FromTo)
	if len(request.WalletIDs) > 0 {
		q.Set("wallet_ids", strings.Join(request.WalletIDs, ","))
	}
	if len(request.StrategyIDs) > 0 {
		q.Set("strategy_ids", strings.Join(request.StrategyIDs, ","))
	}
	if len(request.ChainIDs) > 0 {
		q.Set("chain_ids", strings.Join(request.ChainIDs, ","))
	}
	if request.Type != "" {
		q.Set("type", request.Type)
	}
	if request.Hash != "" {
		q.Set("hash", request.Hash)
	}
	if request.Limit > 0 {
		q.Set("limit", strconv.Itoa(request.Limit))
	}
	if request.Offset > 0 {
		q.Set("offset", strconv.Itoa(request.Offset))
	}
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/transactions", q, nil)
	if err != nil {
		return models.PartnerTransactionsResponse{}, err
	}
	res, err := DoReq[models.PartnerTransactionsResponse](c.httpClient, req)
	if err != nil {
		return models.PartnerTransactionsResponse{}, err
	}
	if res.Data == nil {
		return models.PartnerTransactionsResponse{}, nil
	}
	return *res.Data, nil
}

func (c *Client) GetStrategyPerformance(ctx context.Context, partnerID string, request models.FromTo) (models.StrategyPerformanceResponse, error) {
	q := encodeFromTo(request)
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/performance", q, nil)
	if err != nil {
		return models.StrategyPerformanceResponse{}, err
	}
	res, err := DoReq[models.StrategyPerformanceResponse](c.httpClient, req)
	if err != nil {
		return models.StrategyPerformanceResponse{}, err
	}
	if res.Data == nil {
		return models.StrategyPerformanceResponse{}, nil
	}
	return *res.Data, nil
}

func (c *Client) GetPartnerMetrics(ctx context.Context, partnerID string, request models.FromTo) (models.MetricsPartner, error) {
	q := encodeFromTo(request)

	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/metrics", q, nil)
	if err != nil {
		return models.MetricsPartner{}, err
	}

	res, err := DoReq[models.MetricsPartner](c.httpClient, req)
	if err != nil {
		return models.MetricsPartner{}, err
	}

	if res.Data == nil {
		return models.MetricsPartner{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetPartnerPolicy(ctx context.Context, partnerID string) (models.Policy, error) {
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/policy", nil, nil)
	if err != nil {
		return models.Policy{}, err
	}

	res, err := DoReq[models.Policy](c.httpClient, req)
	if err != nil {
		return models.Policy{}, err
	}

	if res.Data == nil {
		return models.Policy{}, nil
	}

	return *res.Data, nil
}

func (c *Client) UpdatePartnerPolicy(ctx context.Context, partnerID string, policy models.Policy) error {
	req, err := c.do(ctx, http.MethodPut, "/partners/"+url.PathEscape(partnerID)+"/policy", nil, policy)
	if err != nil {
		return err
	}

	_, err = DoReq[any](c.httpClient, req)
	return err
}

func (c *Client) GetAPIKeys(ctx context.Context, partnerID string) ([]models.APIKey, error) {
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/api-keys", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.APIKey](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) CreateAPIKey(ctx context.Context, partnerID, name string) (models.APIKeyFull, error) {
	req, err := c.do(ctx, http.MethodPost, "/partners/"+url.PathEscape(partnerID)+"/api-keys", nil, models.CreateAPIKeyRequest{Name: name})
	if err != nil {
		return models.APIKeyFull{}, err
	}

	res, err := DoReq[models.APIKeyFull](c.httpClient, req)
	if err != nil {
		return models.APIKeyFull{}, err
	}

	if res.Data == nil {
		return models.APIKeyFull{}, nil
	}

	return *res.Data, nil
}

func (c *Client) RevokeAPIKey(ctx context.Context, partnerID, apiKey string) error {
	req, err := c.do(ctx, http.MethodDelete, "/partners/"+url.PathEscape(partnerID)+"/api-keys/"+url.PathEscape(apiKey), nil, nil)
	if err != nil {
		return err
	}

	_, err = DoReq[any](c.httpClient, req)
	return err
}

func (c *Client) GetFeeBalance(ctx context.Context, partnerID string) (models.Balance, error) {
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/fee-balance", nil, nil)
	if err != nil {
		return models.Balance{}, err
	}

	res, err := DoReq[models.Balance](c.httpClient, req)
	if err != nil {
		return models.Balance{}, err
	}

	if res.Data == nil {
		return models.Balance{}, nil
	}

	return *res.Data, nil
}

func (c *Client) ClaimFees(ctx context.Context, reqBody models.PartnerFeeClaimRequest) (models.TxInfo, error) {
	req, err := c.do(ctx, http.MethodPost, "/partners/"+url.PathEscape(reqBody.PartnerID)+"/fees/claim", nil, reqBody)
	if err != nil {
		return models.TxInfo{}, err
	}

	res, err := DoReq[models.TxInfo](c.httpClient, req)
	if err != nil {
		return models.TxInfo{}, err
	}

	if res.Data == nil {
		return models.TxInfo{}, nil
	}

	return *res.Data, nil
}

func (c *Client) ClaimAllFees(ctx context.Context, reqBody models.PartnerFeeClaimAllRequest) (models.TxInfo, error) {
	req, err := c.do(ctx, http.MethodPost, "/partners/"+url.PathEscape(reqBody.PartnerID)+"/fees/claim-all", nil, reqBody)
	if err != nil {
		return models.TxInfo{}, err
	}

	res, err := DoReq[models.TxInfo](c.httpClient, req)
	if err != nil {
		return models.TxInfo{}, err
	}

	if res.Data == nil {
		return models.TxInfo{}, nil
	}

	return *res.Data, nil
}

func (c *Client) GetSupportedChains(ctx context.Context) ([]models.PartnerChain, error) {
	req, err := c.do(ctx, http.MethodGet, "/partners/supported-chains", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.PartnerChain](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) GetAllClaimAddresses(ctx context.Context, partnerID string) ([]models.ClaimAddress, error) {
	req, err := c.do(ctx, http.MethodGet, "/partners/"+url.PathEscape(partnerID)+"/claim-addresses", nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := DoReq[[]models.ClaimAddress](c.httpClient, req)
	if err != nil {
		return nil, err
	}

	if res.Data == nil {
		return nil, nil
	}

	return *res.Data, nil
}

func (c *Client) AttachNewClaimAddress(ctx context.Context, reqBody models.ClaimAddressPublic) error {
	req, err := c.do(ctx, http.MethodPost, "/partners/"+url.PathEscape(reqBody.PartnerID)+"/claim-addresses", nil, reqBody)
	if err != nil {
		return err
	}

	_, err = DoReq[any](c.httpClient, req)
	return err
}

func (c *Client) UpdateExistingClaimAddress(ctx context.Context, reqBody models.ClaimAddress) error {
	chainID := strconv.Itoa(reqBody.ChainID)
	req, err := c.do(ctx, http.MethodPut, "/partners/"+url.PathEscape(reqBody.PartnerID)+"/claim-addresses/"+chainID, nil, reqBody)
	if err != nil {
		return err
	}

	_, err = DoReq[any](c.httpClient, req)
	return err
}
