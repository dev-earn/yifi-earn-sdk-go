package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TypeYield string

const (
	FixedYield      TypeYield = "fixed"
	FLoatingYield   TypeYield = "floating"
	CompoundedYield TypeYield = "compounded"
)

type RateFormat string

const (
	FormatAPR RateFormat = "apr"
	FormatAPY RateFormat = "apy"
)

type MarketsRequest struct {
	Search      string   `json:"search,omitempty"`
	IDs         string   `json:"ids,omitempty"`
	CategoryIDs []string `json:"category_ids,omitempty"`
	PlatformIDs []string `json:"platform_ids,omitempty"`
	ChainIDs    []string `json:"chain_ids,omitempty"`
	TokenIDs    []string `json:"token_ids,omitempty"`

	MinTVLUsd *decimal.Decimal `json:"min_tvl_usd,omitempty"`
	MaxTVLUsd *decimal.Decimal `json:"max_tvl_usd,omitempty"`
	MinAPR    *decimal.Decimal `json:"min_apr,omitempty"`
	MaxAPR    *decimal.Decimal `json:"max_apr,omitempty"`

	SortBy        string `json:"sort_by,omitempty"`        // tvlUsd | apr | apy | createdAt | updatedAt
	SortDirection string `json:"sort_direction,omitempty"` // asc/desc

	Limit  int64 `json:"limit,omitempty"`
	Offset *int  `json:"offset"`
}

type MarketsResponse struct {
	MarketID string         `json:"market_id"`
	Name     string         `json:"name"`
	Chain    string         `json:"chain"`
	Platform PlatformType   `json:"platform"`
	Category CategoryID     `json:"category"`
	Tokens   []TokenAddress `json:"participating_tokens"`

	YieldType  TypeYield       `json:"yield_type"`
	RateFormat RateFormat      `json:"rate_format"`
	AnnualRate decimal.Decimal `json:"annual_rate"`
	TVL        decimal.Decimal `json:"tvl_usd"`
}

type Meta interface {
	Category() CategoryID
}

type MarketCurrentResponse struct {
	MarketID string         `json:"market_id"`
	Name     string         `json:"name"`
	Chain    string         `json:"chain"`
	Platform PlatformType   `json:"platform"`
	Category CategoryID     `json:"category"`
	Tokens   []TokenAddress `json:"participating_tokens"`

	YieldType  TypeYield       `json:"yield_type"`
	RateFormat RateFormat      `json:"rate_format"`
	AnnualRate decimal.Decimal `json:"annual_rate"`
	TVL        decimal.Decimal `json:"tvl_usd"`

	Meta Meta `json:"meta"`
}

type MarketHistoryResponse struct {
	MarketID   string     `json:"market_id"`
	YieldType  TypeYield  `json:"yield_type,omitempty"`
	RateFormat RateFormat `json:"rate_format,omitempty"`
	History    []History  `json:"history"`
}

type History struct {
	AnnualRate decimal.Decimal `json:"annual_rate"`
	TVL        decimal.Decimal `json:"tvl_usd"`
	Timestamp  time.Time       `json:"timestamp"`
}
