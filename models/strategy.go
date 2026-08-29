package models

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type StrategySearch struct {
	Search      *string  `json:"search"`
	IDs         *string  `json:"ids"`
	CategoryIDs []string `json:"category_ids"`
	ChainIDs    []string `json:"chain_ids"`
	PlatformIDs []string `json:"platform_ids"`
	TokenIDs    []string `json:"token_ids"`
	Group       string   `json:"group"`
	Status      *Status  `json:"status"`

	MinTVLUsd *decimal.Decimal `json:"min_tvl_usd,omitempty"`
	MaxTVLUsd *decimal.Decimal `json:"max_tvl_usd,omitempty"`
	MinAPR    *decimal.Decimal `json:"min_apr,omitempty"`
	MaxAPR    *decimal.Decimal `json:"max_apr,omitempty"`

	SortBy        string `json:"sort_by,omitempty"`        // updatedAt|createdAt|apr|apy|tvlUsd
	SortDirection string `json:"sort_direction,omitempty"` // asc/desc

	Limit  int64 `json:"limit"`
	Offset *int  `json:"offset"`
}

type StrategyResponse struct {
	StrategyID      string           `json:"strategy_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	YieldType       TypeYield        `json:"yield_type"`
	RateFormat      RateFormat       `json:"rate_format"`
	CategoryIDs     []string         `json:"category_ids"`
	Platforms       []Platform       `json:"platforms"`
	ChainIDs        []Chain          `json:"chain_ids"`
	InputTokens     []TokenAddress   `json:"input_tokens"`
	OutputTokens    []TokenAddress   `json:"output_tokens"`
	Status          string           `json:"status"`
	Limits          Limits           `json:"limits"`
	RoutingEngine   string           `json:"routing_engine"`
	Fees            Fees             `json:"fee_breakdown"`
	ExecutionConfig *ExecutionConfig `json:"execution_config,omitempty"`
	Metrics         Metrics          `json:"metrics"`
	Risk            Risk             `json:"risk"`
	CreatedAt       int64            `json:"created_at"`
	UpdatedAt       int64            `json:"updated_at"`
}

type StrategiesResponse struct {
	StrategyID   string         `json:"strategy_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Status       string         `json:"status"`
	UpdatedAt    int64          `json:"updated_at"`
	InputTokens  []TokenAddress `json:"input_tokens"`
	OutputTokens []TokenAddress `json:"output_tokens"`
	ChainIDs     []Chain        `json:"chain_ids"` //
	CategoryIDs  []string       `json:"category_ids"`
	Metrics      Metrics        `json:"metrics"`
	Risk         Risk           `json:"risk"`
}

type Limits struct {
	DepositsPaused    bool `json:"deposits_paused"`
	WithdrawalsPaused bool `json:"withdrawals_paused"`
}

type Fees struct {
	MaxFee     decimal.Decimal `json:"max_fee"`
	TotalFee   decimal.Decimal `json:"total_fee"`
	YIFIFee    decimal.Decimal `json:"yifi_fee"`
	PartnerFee decimal.Decimal `json:"partner_fee"`
}

type ExecutionConfig struct {
	InputTokenIDs  []string `json:"input_token_ids"`
	OutputTokenIDs []string `json:"output_token_ids"`
}

type Metrics struct {
	AnnualRate decimal.Decimal `json:"annual_rate"`
	TvlUsd     decimal.Decimal `json:"tvl_usd"`
	// Meta carries the point-in-time reserve snapshot for lending strategies
	// (e.g. Aave). It is populated on the history endpoint for lending
	// strategies and omitted otherwise, so it is nil for non-lending metrics.
	Meta      *MetaLending `json:"meta,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// MetaLending is the lending-strategy meta: a single Aave v3 reserve snapshot
// (rates / config / interest-rate model) keyed by network, chain, underlying
// and aToken.
type MetaLending struct {
	Network    string          `json:"network"`
	ChainID    int             `json:"chain_id"`
	Underlying string          `json:"underlying"`
	AToken     string          `json:"atoken"`
	APYTVL     APYTVLMeta      `json:"apy_tvl"`
	Metadata   ReserveMetadata `json:"metadata"`
	IRM        IRMMeta         `json:"irm"`
}

// APYTVLMeta holds the rate and liquidity snapshot for a reserve.
type APYTVLMeta struct {
	SupplyAPRPct      float64  `json:"supply_apr_pct"`
	SupplyAPYPct      float64  `json:"supply_apy_pct"`
	VarBorrowAPRPct   float64  `json:"var_borrow_apr_pct"`
	VarBorrowAPYPct   float64  `json:"var_borrow_apy_pct"`
	UtilizationPct    float64  `json:"utilization_pct"`
	TVLNative         float64  `json:"tvl_native"`
	PriceUSD          float64  `json:"price_usd"`
	LiquidityIndexRay *big.Int `json:"liquidity_index_ray"`
	UpdatedAtTS       int64    `json:"updated_at_ts"`
	UpdatedAtBlock    int64    `json:"updated_at_block"`
}

// ReserveMetadata is the reserve configuration in human-readable units.
type ReserveMetadata struct {
	Decimals             int     `json:"decimals"`
	LTVPct               float64 `json:"ltv_pct"`
	LiquidationThreshold float64 `json:"liquidation_threshold"`
	LiquidationBonus     float64 `json:"liquidation_bonus"`
	ReserveFactorPct     float64 `json:"reserve_factor_pct"`
	UseAsCollateral      bool    `json:"use_as_collateral"`
	BorrowingEnabled     bool    `json:"borrowing_enabled"`
	StableBorrowEnabled  bool    `json:"stable_borrow_enabled"`
	IsActive             bool    `json:"is_active"`
	IsFrozen             bool    `json:"is_frozen"`
}

// IRMMeta is the interest-rate-model snapshot.
type IRMMeta struct {
	BaseRatePct   float64 `json:"base_rate_pct"`
	Slope1Pct     float64 `json:"slope1_pct"`
	Slope2Pct     float64 `json:"slope2_pct"`
	UOptimalPct   float64 `json:"u_optimal_pct"`
	UMaxExcessPct float64 `json:"u_max_excess_pct"`
}

type Risk struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	RiskScore   RiskLevelType `json:"risk_score"`
}

type RiskLevelType int

const (
	RiskLevelOne   RiskLevelType = 1
	RiskLevelTwo   RiskLevelType = 2
	RiskLevelThree RiskLevelType = 3
	RiskLevelFour  RiskLevelType = 4
	RiskLevelFive  RiskLevelType = 5
)

type Status string

const (
	StatusDraft    Status = "draft"
	StatusActive   Status = "active"
	StatusPaused   Status = "paused"
	StatusArchived Status = "archived"
)

type StrategyHistoryResponse struct {
	StrategyID string     `json:"strategy_id"`
	YieldType  TypeYield  `json:"yield_type,omitempty"`
	RateFormat RateFormat `json:"rate_format,omitempty"`
	History    []Metrics  `json:"history"`
}
