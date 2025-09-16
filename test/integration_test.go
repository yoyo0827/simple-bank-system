package test

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/request"
	"github.com/yoyo0827/simple-bank-system/internal/service"
)

func setupIntegrationDB(t *testing.T) *service.AccountService {
	dsn := "postgres://test:test@localhost:6001/testDB?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("cannot ping test db: %v", err)
	}

	return &service.AccountService{
		DB:                    db,
		AccountRepository:     &repository.AccountRepository{},
		TransactionRepository: &repository.TransactionRepository{},
	}
}

func TestIntegration(t *testing.T) {
	svc := setupIntegrationDB(t)

	// === 建立帳號 ===
	acc1, err := svc.CreateAccount("test1", 100)
	assert.NoError(t, err)
	acc2, err := svc.CreateAccount("test2", 50)
	assert.NoError(t, err)

	// 驗證帳號正確建立
	got1, _ := svc.FindAccount(acc1.ID)
	assert.Equal(t, "test1", got1.Name)
	assert.Equal(t, "100", got1.Balance.String())

	got2, _ := svc.FindAccount(acc2.ID)
	assert.Equal(t, "test2", got2.Name)
	assert.Equal(t, "50", got2.Balance.String())

	// === 存款 ===
	depReq := &request.TransactionRequest{Amount: decimal.NewFromInt(40)} // 存 40
	depRefID, err := svc.CreateTransaction(acc1.ID, depReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, depRefID)

	afterDeposit, _ := svc.FindAccount(acc1.ID)
	assert.Equal(t, "140", afterDeposit.Balance.String()) // 100 + 40

	// === 提款 ===
	wdReq := &request.TransactionRequest{Amount: decimal.NewFromInt(-20)} // 提 20
	wdRefID, err := svc.CreateTransaction(acc1.ID, wdReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, wdRefID)

	afterWithdraw, _ := svc.FindAccount(acc1.ID)
	assert.Equal(t, "120", afterWithdraw.Balance.String()) // 140 - 20

	// === 轉帳 ===
	tfReq := &request.TransferRequest{FromID: acc1.ID, ToID: acc2.ID, Amount: decimal.NewFromInt(30)}
	tfRefID, err := svc.Transfer(tfReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, tfRefID)

	afterTF1, _ := svc.FindAccount(acc1.ID)
	afterTF2, _ := svc.FindAccount(acc2.ID)
	assert.Equal(t, "90", afterTF1.Balance.String()) // 120 - 30
	assert.Equal(t, "80", afterTF2.Balance.String()) // 50 + 30

	// === 查交易紀錄 ===
	txs1, err := svc.FindAccountTransactions(acc1.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(txs1), 3) // 存款 / 提款 / 轉帳

	lastTx := txs1[len(txs1)-1]
	assert.Equal(t, tfRefID, lastTx.RefID) // 最新一筆應該是轉帳
	assert.Equal(t, 1, lastTx.Type)        // acc1 這邊的轉帳是提款

	txs2, err := svc.FindAccountTransactions(acc2.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(txs2), 1) // 至少一筆轉帳紀錄
	assert.Equal(t, tfRefID, txs2[len(txs2)-1].RefID)
	assert.Equal(t, 2, txs2[len(txs2)-1].Type) // acc2 收到的是存款
}
