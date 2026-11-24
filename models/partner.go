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
	StrategyID        string          `json:"strategy_id"`
	UsersCount        int64           `json:"users_count"`
	TransactionsCount int64           `json:"transactions_count"`
	DepositsUSD       decimal.Decimal `json:"deposits_usd"`
	WithdrawalsUSD    decimal.Decimal `json:"withdrawals_usd"`
	FeesEarnedUSD     decimal.Decimal `json:"fees_earned_usd"`
}

type Policy struct {
	ChainIDs         []string         `json:"chain_ids"`
	CategoryIDs      []string         `json:"category_ids"`
	MinTVLUsd        *decimal.Decimal `json:"minTvlUsd"`
	MinAPR           *decimal.Decimal `json:"minApr"`
	RateFormat       string           `json:"rate_format"`
	MaxRiskScore     *int             `json:"maxRiskScore"`
	AllowStrategyIDs []string         `json:"allow_strategy_ids"`
	DenyStrategyIDs  []string         `json:"deny_strategy_ids"`
	UpdatedAt        *time.Time       `json:"updated_at"`

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

// TokenAmount splits a token's fee balance into its buckets.
type TokenAmount struct {
	// Total = Free + Pending.
	Total decimal.Decimal `json:"total"`
	// Free is accrued but not yet pushed to the TokenDistributor.
	Free decimal.Decimal `json:"free"`
	// Pending is pushed to the TokenDistributor and claimable on-chain.
	Pending decimal.Decimal `json:"pending"`
}

type TokenBalance struct {
	TokenID      string          `json:"token_id"`
	ChainID      int             `json:"chain_id"`
	Symbol       string          `json:"symbol"`
	Decimals     int32           `json:"decimals"`
	IsNative     bool            `json:"is_native"`
	IsWrapped    bool            `json:"is_wrapped"`
	Amount       TokenAmount     `json:"amount"`
	AmountUSD    decimal.Decimal `json:"amount_usd"`
	ClaimableUSD decimal.Decimal `json:"claimable_usd"`
	PendingUSD   decimal.Decimal `json:"pending_usd"`
}

type Balance struct {
	PartnerID    string          `json:"partner_id"`
	ClaimableUsd decimal.Decimal `json:"claimable_usd"`
	PendingUsd   decimal.Decimal `json:"pending_usd"`
	TotalUsd     decimal.Decimal `json:"total_usd"`
	Tokens       []TokenBalance  `json:"tokens"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type PartnerFeeClaimRequest struct {
	PartnerID    string `json:"partner_id"`
	TokenID      string `json:"token_id"`
	ChainID      string `json:"chain_id"`
	ClaimAddress string `json:"claim_address"`
}

// PartnerFeeClaimAllRequest claims every token's whole fee balance on a chain
// (claimAll). PartnerID goes in the path; ChainID is the chain name, e.g. "bsc".
type PartnerFeeClaimAllRequest struct {
	PartnerID    string `json:"partner_id"`
	ChainID      string `json:"chain_id"`
	ClaimAddress string `json:"claim_address,omitempty"`
}

type TxInfo struct {
	ClaimID    string          `json:"claim_id"`
	TokenID    string          `json:"token_id"`
	From       string          `json:"from"`
	To         string          `json:"to"`
	TxHash     string          `json:"tx_hash"`
	AmountSent decimal.Decimal `json:"amount_sent"`

	// Data is the unsigned claim()/claimAll() calldata the partner signs from
	// From to To (value 0) to pull the fees out of the distributor.
	Data    string          `json:"data,omitempty"`
	ChainID int64           `json:"chain_id,omitempty"`
	Value   decimal.Decimal `json:"value"`
}

type PartnerChain struct {
	ChainID     int    `json:"chain_id"`
	ChainName   string `json:"chain_name"`
	Distributor string `json:"distributor,omitempty"`
	IsEVM       bool   `json:"is_evm,omitempty"`
}

type ClaimAddress struct {
	PartnerID string    `json:"partner_id"`
	ChainID   int       `json:"chain_id"`
	Address   string    `json:"address"`
	Name      *string   `json:"name,omitempty"`
	TOTPCode  string    `json:"totp_code,omitempty"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClaimAddressPublic struct {
	PartnerID string  `json:"partner_id"`
	ChainID   int     `json:"chain_id"`
	Address   string  `json:"address"`
	Name      *string `json:"name"`
	TOTPCode  string  `json:"totp_code,omitempty"`
}
