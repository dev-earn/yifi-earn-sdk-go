package models

import (
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
	Timestamp  time.Time       `json:"timestamp"`
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
