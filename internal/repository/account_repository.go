package repository

import (
	"github.com/shopspring/decimal"
	"github.com/yoyo0827/simple-bank-system/internal/domain"
)

type AccountRepository struct{}

// 建立帳號
func (r *AccountRepository) FindById(db DBTX, id string) (*domain.Account, error) {
	query := `SELECT id, name, balance FROM accounts WHERE id = $1`
	acc := &domain.Account{}
	err := db.QueryRow(query, id).Scan(&acc.ID, &acc.Name, &acc.Balance)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// 建立帳號
func (r *AccountRepository) CreateUser(db DBTX, account *domain.Account) error {
	query := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, account.Name, account.Balance).Scan(&account.ID)
}

// 更新帳號餘額
func (r *AccountRepository) UpdateBalance(db DBTX, id string, balance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err := db.Exec(query, balance, id)
	return err
}
