package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Summary struct {
	PartnerID      string          `json:"partner_id"`
	UsersCount     int64           `json:"users_count"`
	DepositsUSD    decimal.Decimal `json:"deposits_usd"`
	WithdrawalsUSD decimal.Decimal `json:"withdrawals_usd"`
	FeesEarnedUSD  decimal.Decimal `json:"fees_earned_usd"`
	From           time.Time       `json:"from"`
	To             time.Time       `json:"to"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
