package domain

import "github.com/shopspring/decimal"

type Account struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}
