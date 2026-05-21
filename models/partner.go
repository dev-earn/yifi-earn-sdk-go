package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type MetricsPartner struct {
	PartnerID string    `json:"partner_id"`
	Items     []Metric  `json:"items"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Metric struct {
	StrategyID     string          `json:"strategy_id"`
	UsersCount     int64           `json:"users_count"`
	DepositsUSD    decimal.Decimal `json:"deposits_usd"`
	WithdrawalsUSD decimal.Decimal `json:"withdrawals_usd"`
	FeesEarnedUSD  decimal.Decimal `json:"fees_earned_usd"`
}

type Policy struct {
	AllowStrategyIds   []string         `json:"allow_strategy_ids,omitempty"`
	DenyStrategyIds    []string         `json:"deny_strategy_ids,omitempty"`
	AllowedChainIds    []string         `json:"allowed_chain_ids,omitempty"`
	AllowedCategoryIds []string         `json:"allowed_category_ids,omitempty"`
	MinTvlUsd          *decimal.Decimal `json:"min_tvl_usd,omitempty"`
	MinRate            *decimal.Decimal `json:"min_rate,omitempty"`
	RateType           string           `json:"rate_type,omitempty"`
	MaxRiskLevel       string           `json:"max_risk_level,omitempty"`
	UpdatedAt          *time.Time       `json:"updated_at,omitempty"`

	Fees *FeePolicy `json:"fees_policies,omitempty"`
}

type FeePolicy struct {
	CommonFee     *decimal.Decimal `json:"common_fee,omitempty"` //  allowed range of values is from 0 to 0.7
	StrategiesFee []StrategyFee    `json:"strategies_fee,omitempty"`
}

type StrategyFee struct {
	StrategyID           string           `json:"strategy_id,omitempty"`
	PartnerFeePreference *decimal.Decimal `json:"partner_fee_preference"` //  allowed range of values is from 0 to 0.7
}

type APIKey struct {
	APIKey     string     `json:"api_key"`
	SecretByte []byte     `json:"secret_byte,omitempty"`
	Hash       []byte     `json:"hash,omitempty"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Status     bool       `json:"status"`
	Meta       any        `json:"meta,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

type APIKeyFull struct {
	APIKey
	Secret string `json:"secret"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type TokenBalance struct {
	TokenID string          `json:"token_id"`
	ChainID int             `json:"chain_id"`
	Amount  decimal.Decimal `json:"amount"`
}

type Balance struct {
	PartnerID    string          `json:"partner_id"`
	AvailableUsd decimal.Decimal `json:"available_usd"`
	Tokens       []TokenBalance  `json:"tokens"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type PartnerFeeClaimRequest struct {
	PartnerID       string          `json:"partner_id"`
	TokenID         string          `json:"token_id"`
	IsTokenNative   bool            `json:"is_token_native"`
	ChainID         string          `json:"chain_id"`
	ClaimAddress    string          `json:"claim_address"`
	AmountRequested decimal.Decimal `json:"amount_requested"`
	GasLimit        decimal.Decimal `json:"gas_limit"`
	Slippage        decimal.Decimal `json:"slippage"`
}

type TxInfo struct {
	ClaimID    string          `json:"claim_id"`
	TokenID    string          `json:"token_id"`
	From       string          `json:"from"`
	To         string          `json:"to"`
	TxHash     string          `json:"tx_hash"`
	AmountSent decimal.Decimal `json:"amount_sent"`
}

type ClaimAddress struct {
	PartnerID string    `json:"partner_id"`
	ChainID   int       `json:"chain_id"`
	Address   string    `json:"address"`
	Name      *string   `json:"name,omitempty"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClaimAddressPublic struct {
	PartnerID string  `json:"partner_id"`
	ChainID   int     `json:"chain_id"`
	Address   string  `json:"address"`
	Name      *string `json:"name"`
}
