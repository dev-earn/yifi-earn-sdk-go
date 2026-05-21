package models

import "github.com/shopspring/decimal"

type WalletTx struct {
	OrderID    string `json:"order_id"`
	PartnerID  string `json:"partner_id"`
	StrategyID string `json:"strategy_id"`
	TokenID    string `json:"token_id"`
	ChainID    string `json:"chain_id"`
	WalletID   string `json:"wallet_id"`
	IsDeposit  bool   `json:"is_deposit"`

	TokenPriceUSD   decimal.Decimal `json:"token_price_usd"`
	InputAmount     decimal.Decimal `json:"input_amount"`
	InputAmountUSD  decimal.Decimal `json:"input_amount_usd"`
	OutputAmount    decimal.Decimal `json:"output_amount"`
	OutputAmountUSD decimal.Decimal `json:"output_amount_usd"`
	TxHash          string          `json:"tx_hash"`

	CreatedAt int64 `json:"created_at"`

	PartnerFee decimal.Decimal `json:"partner_fee"`
	YifiFee    decimal.Decimal `json:"yifi_fee"`
}
