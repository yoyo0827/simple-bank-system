package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/request"
)

// 單元測試 Transaction (存款)
func TestTransaction_Deposit(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()

	accountRepo := &repository.AccountRepository{}
	transactionRepo := &repository.TransactionRepository{}
	svc := &AccountService{DB: db, AccountRepository: accountRepo, TransactionRepository: transactionRepo}

	// 模擬帳號查詢
	rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
		AddRow("acc1", "Alice", "100")
	mock.ExpectQuery(`SELECT (.+) FROM accounts WHERE id = \$1`).
		WithArgs("acc1").
		WillReturnRows(rows)

	// 模擬 update balance
	mock.ExpectExec(`UPDATE accounts SET balance = .* WHERE id = .*`).
		WithArgs(sqlmock.AnyArg(), "acc1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 模擬 insert transaction
	mock.ExpectQuery(`INSERT INTO transactions`).
		WithArgs("acc1", 2, "50", sqlmock.AnyArg(), "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	req := &request.TransactionRequest{Amount: decimal.NewFromInt(50)}
	refID, err := svc.CreateTransaction("acc1", req)

	assert.NoError(t, err)
	assert.NotEmpty(t, refID)
}

// 單元測試 Transaction (提款)
func TestTransaction_Withdraw(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()

	accountRepo := &repository.AccountRepository{}
	transactionRepo := &repository.TransactionRepository{}
	svc := &AccountService{DB: db, AccountRepository: accountRepo, TransactionRepository: transactionRepo}

	// 模擬帳號查詢
	rows := sqlmock.NewRows([]string{"id", "name", "balance"}).
		AddRow("acc1", "Alice", "100")
	mock.ExpectQuery(`SELECT (.+) FROM accounts WHERE id = \$1`).
		WithArgs("acc1").
		WillReturnRows(rows)

	// 模擬 update balance
	mock.ExpectExec(`UPDATE accounts SET balance = .* WHERE id = .*`).
		WithArgs(sqlmock.AnyArg(), "acc1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 模擬 insert transaction
	mock.ExpectQuery(`INSERT INTO transactions`).
		WithArgs("acc1", 1, "50", sqlmock.AnyArg(), "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	req := &request.TransactionRequest{Amount: decimal.NewFromInt(-50)}
	refID, err := svc.CreateTransaction("acc1", req)

	assert.NoError(t, err)
	assert.NotEmpty(t, refID)
}

// 單元測試 Transfer (轉帳)
func TestTransfer(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	accountRepo := &repository.AccountRepository{}
	transactionRepo := &repository.TransactionRepository{}
	svc := &AccountService{DB: db, AccountRepository: accountRepo, TransactionRepository: transactionRepo}

	// 開始 transaction
	mock.ExpectBegin()

	// 查詢 from 帳號
	fromRows := sqlmock.NewRows([]string{"id", "name", "balance"}).
		AddRow("from1", "Alice", "100")
	mock.ExpectQuery(`SELECT (.+) FROM accounts WHERE id = \$1`).
		WithArgs("from1").
		WillReturnRows(fromRows)

	// 查詢 to 帳號
	toRows := sqlmock.NewRows([]string{"id", "name", "balance"}).
		AddRow("to1", "Bob", "50")
	mock.ExpectQuery(`SELECT (.+) FROM accounts WHERE id = \$1`).
		WithArgs("to1").
		WillReturnRows(toRows)

	// 更新餘額
	mock.ExpectExec(`UPDATE accounts SET balance = .* WHERE id = .*`).
		WithArgs(sqlmock.AnyArg(), "from1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`UPDATE accounts SET balance = .* WHERE id = .*`).
		WithArgs(sqlmock.AnyArg(), "to1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	//  寫入交易紀錄
	mock.ExpectQuery(`INSERT INTO transactions`).
		WithArgs("from1", 1, "30", sqlmock.AnyArg(), "Transfer to Bob").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`INSERT INTO transactions`).
		WithArgs("to1", 2, "30", sqlmock.AnyArg(), "Transfer from Alice").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	// mock commit
	mock.ExpectCommit()

	// 執行轉帳
	req := &request.TransferRequest{FromID: "from1", ToID: "to1", Amount: decimal.NewFromInt(30)}
	refID, err := svc.Transfer(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, refID)
	assert.True(t, uuid.Validate(refID) == nil)
}
