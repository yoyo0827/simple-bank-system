package repository

import (
	"database/sql"

	"github.com/yoyo0827/simple-bank-system/internal/domain"
)

type TransactionRepository struct {
	DB *sql.DB
}

// 寫入交易紀錄
func (r *TransactionRepository) InsertTransactions(accountID string, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (account_id, type, amount) VALUES ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, accountID, tx.Type, tx.Amount).Scan(&tx.ID)
}

// 根據 ID 查詢交易紀錄
func (r *TransactionRepository) FindById(id string) (*domain.Transaction, error) {
	query := `SELECT id, account_id, type, amount, created_at FROM transactions WHERE id = $1`
	tx := &domain.Transaction{}
	err := r.DB.QueryRow(query, id).Scan(&tx.ID, &tx.Name, &tx.Type, &tx.Amount, &tx.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
