package repository

import "database/sql"

// DBTX 是一個介面，抽象化 sql.DB 和 sql.Tx 的共同行為
type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
