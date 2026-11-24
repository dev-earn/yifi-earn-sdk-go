package models

import "github.com/shopspring/decimal"

type WalletTx struct {
	OrderID       string `json:"order_id"`
	PartnerID     string `json:"partner_id"`
	StrategyID    string `json:"strategy_id"`
	TokenID       string `json:"token_id"`
	TokenIDActual string `json:"token_id_actual,omitempty"`
	ChainID       string `json:"chain_id"`
	WalletID      string `json:"wallet_id"`
	IsDeposit     bool   `json:"is_deposit"`

	Symbol    string `json:"symbol"`
	IsNative  bool   `json:"is_native"`
	IsWrapped bool   `json:"is_wrapped"`

	TokenPriceUSD   decimal.Decimal `json:"token_price_usd"`
	Decimals        int32           `json:"decimals"`
	InputAmount     decimal.Decimal `json:"input_amount"`
	InputAmountUSD  decimal.Decimal `json:"input_amount_usd"`
	OutputAmount    decimal.Decimal `json:"output_amount"`
	OutputAmountUSD decimal.Decimal `json:"output_amount_usd"`
	TxHash          string          `json:"tx_hash"`

	CreatedAt int64 `json:"created_at"`

	PartnerFee decimal.Decimal `json:"partner_fee"`
	YifiFee    decimal.Decimal `json:"yifi_fee"`

	Meta *FullTxView `json:"meta,omitempty"`
}

type WalletTxResponse struct {
	OrderID       string `json:"order_id"`
	PartnerID     string `json:"partner_id,omitempty"`
	StrategyID    string `json:"strategy_id,omitempty"`
	TokenID       string `json:"token_id,omitempty"`
	TokenIDActual string `json:"token_id_actual,omitempty"`
	ChainID       string `json:"chain_id,omitempty"`
	WalletID      string `json:"wallet_id,omitempty"`
	IsDeposit     bool   `json:"is_deposit,omitempty"`

	Symbol    string `json:"symbol,omitempty"`
	Decimals  int32  `json:"decimals,omitempty"`
	IsNative  bool   `json:"is_native,omitempty"`
	IsWrapped bool   `json:"is_wrapped,omitempty"`

	TokenPriceUSD   *decimal.Decimal `json:"token_price_usd,omitempty"`
	InputAmount     *decimal.Decimal `json:"input_amount,omitempty"`
	InputAmountUSD  *decimal.Decimal `json:"input_amount_usd,omitempty"`
	OutputAmount    *decimal.Decimal `json:"output_amount,omitempty"`
	OutputAmountUSD *decimal.Decimal `json:"output_amount_usd,omitempty"`
	TxHash          string           `json:"tx_hash,omitempty"`
	CreatedAt       int64            `json:"created_at,omitempty"`
	PartnerFee      *decimal.Decimal `json:"partner_fee,omitempty"`
	YifiFee         *decimal.Decimal `json:"yifi_fee,omitempty"`
	Meta            *FullTxView      `json:"meta,omitempty"`
}

type FullTxView struct {
	TxHash            string                `json:"tx_hash"`
	From              string                `json:"from"`
	To                string                `json:"to"`
	Value             string                `json:"value"`
	Nonce             uint64                `json:"nonce"`
	BlockNumber       uint64                `json:"block_number"`
	BlockHash         string                `json:"block_hash"`
	Status            uint64                `json:"status"`
	GasLimit          uint64                `json:"gas_limit"`
	GasUsed           uint64                `json:"gas_used"`
	GasPrice          string                `json:"gas_price"`
	EffectiveGasPrice string                `json:"effective_gas_price"`
	TransactionFee    string                `json:"transaction_fee"`
	Type              uint8                 `json:"type"`
	PortalInput       *PrintablePortalInput `json:"portal_input,omitempty"`
	TokenTransfers    []TokenTransfer       `json:"token_transfers,omitempty"`
	Zap               string                `json:"zap,omitempty"`
}

type TokenTransfer struct {
	Token  string `json:"token"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type PrintablePortalInput struct {
	InputToken      string                `json:"input_token"`
	InputAmount     string                `json:"input_amount"`
	OutputToken     string                `json:"output_token"`
	MinOutputAmount string                `json:"min_output_amount"`
	Recipient       string                `json:"recipient"`
	Partner         string                `json:"partner"`
	Calls           []PrintablePortalCall `json:"calls"`
}

type PrintablePortalCall struct {
	InputToken  string `json:"input_token"`
	Target      string `json:"target"`
	Data        string `json:"data"`
	AmountIndex string `json:"amount_index"`
}
