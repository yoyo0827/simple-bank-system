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
	query := `INSERT INTO transactions (account_id, type, amount, ref_id, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.DB.QueryRow(query, accountID, tx.Type, tx.Amount, tx.RefID, tx.Description).Scan(&tx.ID)
}

// 根據 ID 查詢交易紀錄
func (r *TransactionRepository) FindById(id string) ([]*domain.Transaction, error) {
	query := `SELECT t.id, a.name, t.type, t.amount,t.ref_id,COALESCE(t.description, ''),t.created_at FROM transactions t JOIN accounts a ON t.account_id = a.id WHERE t.account_id = $1`
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		tx := &domain.Transaction{}
		if err := rows.Scan(&tx.ID, &tx.Name, &tx.Type, &tx.Amount, &tx.RefID, &tx.Description, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
