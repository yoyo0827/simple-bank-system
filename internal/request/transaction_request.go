package request

import "github.com/shopspring/decimal"

type TransactionRequest struct {
	Amount decimal.Decimal `json:"amount"`
}
