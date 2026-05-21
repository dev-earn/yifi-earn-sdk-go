package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TokenStatsHistory struct {
	TokenID string         `json:"token_id"`
	History []TokenHistory `json:"history"`
}

type TokenHistory struct {
	Timestamp time.Time  `json:"timestamp"`
	Metrics   TokenStats `json:"metrics"`
}

type TokenFullCurrent struct {
	Token
	TokenStats
}

type TokenStats struct {
	PriceUSD          decimal.Decimal `json:"price_usd"`
	Volume24hUSD      decimal.Decimal `json:"volume_24h_usd"`
	CirculatingSupply decimal.Decimal `json:"circulating_supply"`
	MarketCap         decimal.Decimal `json:"market_cap"`
	FDV               decimal.Decimal `json:"fdv"`
	Holders           decimal.Decimal `json:"holders"`
}

type Token struct {
	TokenCore
	TagIDs []string `json:"tag_ids"`
	TokenGroup
	Addresses []Address `json:"addresses"`
}

type TokenMarket struct {
	Token
	Chain    string     `json:"chain"`
	Platform string     `json:"platform"`
	Category CategoryID `json:"category"`
	Metadata struct {
		MarketAddress string `json:"market_address"`
	} `json:"metadata"`
}

type TokenCore struct {
	TokenID string `json:"token_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Image   string `json:"image,omitempty"`
}

type TokenGroup struct {
	GroupID          string `json:"token_group"`
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
	Logo             string `json:"logo_url"`
}

type TokenCoreGroup struct {
	TokenID string `json:"token_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Image   string `json:"image,omitempty"`
	GroupID string `json:"token_group"`
}

type Address struct {
	Chain    string `json:"chain_id"`
	Addr     string `json:"address"`
	Decimals int    `json:"decimals"`
}

type TokenTag struct {
	TagID       string `json:"tag_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TokenAddress struct {
	TokenID  string          `json:"token_id"`
	Chain    string          `json:"chain"`
	Address  string          `json:"address"`
	Decimals int             `json:"decimals"`
	LogoUrl  string          `json:"logo_url"`
	Role     string          `json:"role"`
	Symbol   string          `json:"symbol"`
	PriceUSD decimal.Decimal `json:"price_usd"`
}

type SearchTokensMarkets struct {
	SearchSymbol string          `json:"search_symbol"`
	MinLiquidity decimal.Decimal `json:"min_liquidity"`
	Tag          string          `json:"tag"`
	MinAPY       decimal.Decimal `json:"min_apy"`
	Platform     string          `json:"platform,omitempty"`
	Network      string          `json:"network,omitempty"`
}
