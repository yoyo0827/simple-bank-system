package request

import "github.com/shopspring/decimal"

type TransferRequest struct {
	FromID string          `json:"from_id"`
	ToID   string          `json:"to_id"`
	Amount decimal.Decimal `json:"amount"`
}
