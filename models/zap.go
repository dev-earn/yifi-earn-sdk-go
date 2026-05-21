package models

import "github.com/shopspring/decimal"

type ZapQuoteRequest struct {
	StrategyID                  string          `json:"strategy_id"`
	Sender                      string          `json:"sender"`
	InputToken                  string          `json:"input_token"`
	InputAmount                 string          `json:"input_amount"`
	OutputToken                 string          `json:"output_token"`
	SlippageTolerancePercentage decimal.Decimal `json:"slippage_tolerance_percentage"`
}

type Signature struct {
	OrderID string `json:"orderId"`
	Data    string `json:"data"`
	From    string `json:"from"`
	To      string `json:"to"`
	Value   string `json:"value"`
}
