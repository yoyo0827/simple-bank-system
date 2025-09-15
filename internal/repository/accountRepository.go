package repository

import (
	"database/sql"

	"github.com/shopspring/decimal"
	"github.com/yoyo0827/simple-bank-system/internal/domain"
)

type AccountRepository struct {
	DB *sql.DB
}

// 建立帳號
func (r *AccountRepository) FindById(id string) (*domain.Account, error) {
	query := `SELECT id, name, balance FROM accounts WHERE id = $1`
	acc := &domain.Account{}
	err := r.DB.QueryRow(query, id).Scan(&acc.ID, &acc.Name, &acc.Balance)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// 建立帳號
func (r *AccountRepository) InsertUser(account *domain.Account) error {
	query := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id`
	return r.DB.QueryRow(query, account.Name, account.Balance).Scan(&account.ID)
}

// 更新帳號餘額
func (r *AccountRepository) UpdateBalance(id string, balance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, balance, id)
	return err
}
