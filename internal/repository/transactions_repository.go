package repository

import (
	"github.com/yoyo0827/simple-bank-system/internal/domain"
)

type TransactionRepository struct{}

// 寫入交易紀錄
func (r *TransactionRepository) InsertTransactions(db DBTX, accountID string, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (account_id, type, amount, ref_id, description) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.QueryRow(query, accountID, tx.Type, tx.Amount, tx.RefID, tx.Description).Scan(&tx.ID)
}

// 根據 ID 查詢交易紀錄
func (r *TransactionRepository) FindById(db DBTX, id string) ([]*domain.Transaction, error) {
	query := `SELECT t.id, a.name, t.type, t.amount,t.ref_id,COALESCE(t.description, ''),t.created_at FROM transactions t JOIN accounts a ON t.account_id = a.id WHERE t.account_id = $1`
	rows, err := db.Query(query, id)
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
