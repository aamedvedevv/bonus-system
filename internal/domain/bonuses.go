package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrNoWithdraws = errors.New("the user has no withdraws")
	ErrNoBonuses   = errors.New("not enough bonuses")
)

type Withdraw struct {
	OrderID    string          `json:"order"`
	Bonuses    decimal.Decimal `json:"sum"`
	UploadedAt string          `json:"processed_at"`
	UserID     int64           `json:"-"`
}

type BalanceOutput struct {
	Bonuses  decimal.Decimal `json:"current"`
	Withdraw decimal.Decimal `json:"withdrawn"`
}
