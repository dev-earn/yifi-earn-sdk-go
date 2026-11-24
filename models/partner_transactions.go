package models

type PartnerTransactionsRequest struct {
	FromTo
	WalletIDs   []string
	StrategyIDs []string
	ChainIDs    []string
	Type        string
	Hash        string
	Limit       int
	Offset      int
}

type PartnerTransaction struct {
	Date            string `json:"date"`
	Type            string `json:"type"`
	Wallet          string `json:"wallet"`
	Strategy        string `json:"strategy"`
	StrategyName    string `json:"strategy_name"`
	Chain           string `json:"chain"`
	InputAmount     string `json:"input_amount"`
	InputAmountUSD  string `json:"input_amount_usd"`
	InputToken      string `json:"input_token"`
	OutputAmount    string `json:"output_amount"`
	OutputAmountUSD string `json:"output_amount_usd"`
	OutputToken     string `json:"output_token"`
	FeeAmount       string `json:"fee_amount"`
	FeeAmountUSD    string `json:"fee_amount_usd"`
	FeeToken        string `json:"fee_token"`
	TxHash          string `json:"tx_hash"`
}

type PartnerTransactionsResponse struct {
	Items  []PartnerTransaction `json:"items"`
	Total  int                  `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type StrategyFeeBreakdown struct {
	TokenID      string `json:"token_id"`
	Symbol       string `json:"symbol"`
	FeeAmount    string `json:"fee_amount"`
	FeeAmountUSD string `json:"fee_amount_usd"`
}

type StrategyPerformance struct {
	StrategyID          string                 `json:"strategy_id"`
	StrategyName        string                 `json:"strategy_name"`
	WalletsParticipated int                    `json:"wallets_participated"`
	DepositsAmount      string                 `json:"deposits_amount"`
	DepositsToken       string                 `json:"deposits_token"`
	DepositsAmountUSD   string                 `json:"deposits_amount_usd"`
	DepositsCount       int64                  `json:"deposits_count"`
	WithdrawalsAmount   string                 `json:"withdrawals_amount"`
	WithdrawalsToken    string                 `json:"withdrawals_token"`
	WithdrawalsUSD      string                 `json:"withdrawals_amount_usd"`
	WithdrawalsCount    int64                  `json:"withdrawals_count"`
	FeesEarnedUSD       string                 `json:"fees_earned_usd"`
	FeesByToken         []StrategyFeeBreakdown `json:"fees_by_token"`
}

type StrategyPerformanceResponse struct {
	Items  []StrategyPerformance `json:"items"`
	Total  int                   `json:"total"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}
