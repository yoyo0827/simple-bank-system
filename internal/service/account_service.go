package service

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yoyo0827/simple-bank-system/internal/domain"
	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/request"
)

type AccountService struct {
	DB                    *sql.DB
	AccountRepository     *repository.AccountRepository
	TransactionRepository *repository.TransactionRepository
}

// 查詢帳號
func (s *AccountService) FindAccount(id string) (*domain.Account, error) {
	acc, err := s.AccountRepository.FindById(s.DB, id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// 建立帳號
func (s *AccountService) CreateAccount(name string, balance float64) (*domain.Account, error) {
	if balance < 0 {
		return nil, errors.New("balance cannot be negative")
	}
	acc := &domain.Account{
		Name:    name,
		Balance: decimal.NewFromFloat(balance), // float64 -> decimal
	}
	err := s.AccountRepository.CreateUser(s.DB, acc)
	return acc, err
}

// 交易
func (s *AccountService) CreateTransaction(id string, req *request.TransactionRequest) (string, error) {
	// 交易安全，使用 transaction
	transaction, err := s.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		}
	}()
	// 查詢帳號
	acc, err := s.AccountRepository.FindById(s.DB, id)
	if err != nil {
		return "", err
	}
	// 定義交易類型 1=提款, 2=存款
	var txType int
	if req.Amount.IsNegative() {
		txType = 1
	} else {
		txType = 2
	}

	// 更新帳號餘額
	newBalance := acc.Balance.Add(req.Amount)
	if newBalance.IsNegative() {
		return "", errors.New("insufficient funds, cannot withdraw more than the current balance")
	}
	if err := s.AccountRepository.UpdateBalance(transaction, id, newBalance); err != nil {
		return "", err
	}
	// 寫入交易紀錄
	refID := uuid.New().String()
	tx := newTransaction(acc.Name, txType, req.Amount, refID, "")
	if err := s.TransactionRepository.InsertTransactions(transaction, acc.ID, tx); err != nil {
		return "", errors.New("failed to insert transaction record: " + err.Error())
	}
	// 印出交易紀錄 log
	log.Printf(
		"[Transaction] ref_id=%s | acc=%s | amount=%s | at=%s",
		refID, acc.ID, req.Amount.String(), time.Now().Format(time.RFC3339),
	)
	// 提交交易
	if err := transaction.Commit(); err != nil {
		return "", err
	}
	return refID, nil
}

// 轉帳
func (s *AccountService) Transfer(req *request.TransferRequest) (string, error) {
	// 驗證轉帳金額
	amount := req.Amount
	if err := validateAmount(amount); err != nil {
		return "", err
	}

	// 交易安全，使用 transaction
	transaction, err := s.DB.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if p := recover(); p != nil {
			transaction.Rollback()
			panic(p)
		}
	}()
	defer transaction.Rollback()

	// 查詢雙方帳號
	fromAcc, err := s.AccountRepository.FindById(s.DB, req.FromID)
	if err != nil {
		return "", err
	}
	toAcc, err := s.AccountRepository.FindById(s.DB, req.ToID)
	if err != nil {
		return "", err
	}
	// 檢查餘額是否足夠
	if fromAcc.Balance.Cmp(req.Amount) < 0 {
		return "", errors.New("insufficient funds, cannot transfer more than the current balance")
	}
	// 更新雙方帳號餘額
	if err := s.AccountRepository.UpdateBalance(transaction, fromAcc.ID, fromAcc.Balance.Sub(amount)); err != nil {
		return "", err
	}
	if err := s.AccountRepository.UpdateBalance(transaction, toAcc.ID, toAcc.Balance.Add(amount)); err != nil {
		return "", err
	}
	// 寫入交易紀錄
	refID := uuid.New().String()
	if err := s.TransactionRepository.InsertTransactions(transaction, fromAcc.ID,
		newTransaction(fromAcc.Name, 1, req.Amount, refID, "Transfer to "+toAcc.Name)); err != nil {
		return "", err
	}
	if err := s.TransactionRepository.InsertTransactions(transaction, toAcc.ID,
		newTransaction(toAcc.Name, 2, req.Amount, refID, "Transfer from "+fromAcc.Name)); err != nil {
		return "", err
	}
	// 印出轉帳紀錄 log
	log.Printf(
		"[Transfer] ref_id=%s | from_acc=%s | to_acc=%s | amount=%s | at=%s",
		refID, fromAcc.ID, toAcc.ID, amount.String(), time.Now().Format(time.RFC3339),
	)
	// 提交交易
	if err := transaction.Commit(); err != nil {
		return "", err
	}
	return refID, nil
}

// 查詢帳號交易紀錄
func (s *AccountService) FindAccountTransactions(id string) ([]*domain.Transaction, error) {
	transactions, err := s.TransactionRepository.FindById(s.DB, id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// 驗證轉帳金額
func validateAmount(amount decimal.Decimal) error {
	if amount.IsZero() {
		return errors.New("amount cannot be zero")
	}
	if amount.IsNegative() {
		return errors.New("amount cannot be negative")
	}
	return nil
}

// 建立交易紀錄
func newTransaction(name string, txType int, amount decimal.Decimal, refID, desc string) *domain.Transaction {
	return &domain.Transaction{
		Name:        name,
		Type:        txType,
		Amount:      amount.Abs(),
		RefID:       refID,
		Description: desc,
	}
}
