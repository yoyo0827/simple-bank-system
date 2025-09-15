package domain

import "github.com/shopspring/decimal"

type Transaction struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Type        int             `json:"type"` // 1=提款, 2=存款
	Amount      decimal.Decimal `json:"amount"`
	RefID       string          `json:"ref_id"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"created_at"`
}
