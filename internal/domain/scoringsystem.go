package domain

import "github.com/shopspring/decimal"

type ScoringSystem struct {
	OrderID string          `json:"order"`
	Status  OrderStatus     `json:"status"`
	Bonuses decimal.Decimal `json:"accrual"`
}
