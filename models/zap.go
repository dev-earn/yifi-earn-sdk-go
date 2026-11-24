package models

import "github.com/shopspring/decimal"

// ZapQuoteRequest describes a zap: swap InputToken into the strategy and
// deposit, in one transaction. Token fields are "chain:address" IDs; use the
// zero address for the native coin. Pass it to SDK.BuildZap.
type ZapQuoteRequest struct {
	// StrategyID is the target strategy, e.g. "bsc:0xb0b8...a1b".
	StrategyID string `json:"strategy_id"`
	// Sender is the wallet that will sign and send the returned transaction.
	Sender string `json:"sender"`
	// InputToken is the asset being spent; native coin is the zero address.
	InputToken string `json:"input_token"`
	// InputAmount is the amount of InputToken in its smallest unit (wei), as a
	// decimal string — e.g. "1000000" for 1 USDT (6 decimals).
	InputAmount string `json:"input_amount"`
	// OutputToken is the strategy's deposit token (often equal to StrategyID).
	OutputToken string `json:"output_token"`
	// SlippageTolerancePercentage is the max acceptable slippage, in percent
	// (e.g. 1 = 1%). The zap reverts if exceeded.
	SlippageTolerancePercentage decimal.Decimal `json:"slippage_tolerance_percentage"`
	// Engine optionally selects a strategy's declared fallback routing engine.
	Engine string `json:"engine,omitempty"`
	// OrderID is required by the current API and links build to a prior QuoteZap.
	OrderID string `json:"order_id,omitempty"`
}

// ZapQuoteResponse is the preview returned by QuoteZap. Amounts use the output
// token's smallest unit.
type ZapQuoteResponse struct {
	OrderID                  string          `json:"order_id"`
	ExpiresIn                int64           `json:"expires_in"`
	MinOutputAmount          string          `json:"min_output_amount"`
	MinOutputAmountUsd       float64         `json:"min_output_amount_usd"`
	SlippageTolerancePercent decimal.Decimal `json:"slippage_tolerance_percentage"`
	OutputToken              string          `json:"output_token"`
	OutputAmount             string          `json:"output_amount"`
	OutputAmountUsd          float64         `json:"output_amount_usd"`
	FeeToken                 string          `json:"fee_token"`
	FeeAmount                string          `json:"fee_amount"`
	FeeAmountUsd             float64         `json:"fee_amount_usd"`
	Steps                    []string        `json:"steps"`
	Engine                   string          `json:"engine"`
	RouteTool                string          `json:"engine_route_tool,omitempty"`
	ApprovalAddress          string          `json:"approval_address"`
	ApprovalRequired         bool            `json:"approval_required,omitempty"`
	GasEstimate              *GasEstimate    `json:"gas_estimate,omitempty"`
	FeeRecipient             string          `json:"fee_recipient"`
}

type GasEstimate struct {
	Amount    string  `json:"amount"`
	Token     string  `json:"token"`
	AmountUSD float64 `json:"amount_usd,omitempty"`
}

// Signature is an unsigned EVM transaction returned by BuildZap (and the fee
// claim calls). Sign it with the Sender's key and broadcast it; see
// examples/sign.go. Note Value is a decimal wei string, not hex.
type Signature struct {
	// OrderID identifies the quote server-side (for support/tracing).
	OrderID string `json:"orderId"`
	// Data is the transaction calldata, 0x-prefixed hex.
	Data string `json:"data"`
	// From is the expected sender; must match the signing key.
	From string `json:"from"`
	// To is the contract to call (e.g. the Portals router).
	To string `json:"to"`
	// Value is the native coin to attach, as a decimal wei string ("0" if none).
	Value string `json:"value"`
	// Engine is the routing engine that produced the frozen quote.
	Engine string `json:"engine,omitempty"`
	// ApprovalAddress is the spender required for ERC-20 inputs.
	ApprovalAddress string `json:"approval_address,omitempty"`
	// OutputAmount and MinOutputAmount are taken from the same frozen quote.
	OutputAmount    string `json:"output_amount,omitempty"`
	MinOutputAmount string `json:"min_output_amount,omitempty"`
	// GasLimit and ChainID are supplied by engines that simulate the route.
	GasLimit string `json:"gasLimit,omitempty"`
	ChainID  string `json:"chainId,omitempty"`
}
