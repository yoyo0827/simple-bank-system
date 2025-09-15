package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yoyo0827/simple-bank-system/internal/domain"
	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/request"
)

type AccountService struct {
	AccountRepository     *repository.AccountRepository
	TransactionRepository *repository.TransactionRepository
}

// 建立帳號
func (s *AccountService) FindAccount(id string) (*domain.Account, error) {
	acc, err := s.AccountRepository.FindById(id)
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
	err := s.AccountRepository.InsertUser(acc)
	return acc, err
}

// 交易
func (s *AccountService) Transaction(id string, req *request.TransactionRequest) (*domain.Transaction, error) {
	acc, err := s.AccountRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	var txType int // 1=提款, 2=存款
	if req.Amount.IsNegative() {
		txType = 1 // 提款
	} else {
		txType = 2 // 存款
	}

	// 更新帳號餘額
	newBalance := acc.Balance.Add(req.Amount)
	if newBalance.IsNegative() {
		return nil, errors.New("insufficient funds, cannot withdraw more than the current balance")
	}

	err = s.AccountRepository.UpdateBalance(id, newBalance)

	if err != nil {
		return nil, err
	}
	refID := uuid.New().String()
	// 寫入交易紀錄
	tx := &domain.Transaction{
		Name:   acc.Name,
		Type:   txType,
		Amount: req.Amount,
		RefID:  refID,
	}
	err = s.TransactionRepository.InsertTransactions(acc.ID, tx)
	if err != nil {
		return nil, err
	}

	return tx, err
}

// 轉帳
func (s *AccountService) Transfer(req *request.TransferRequest) (string, error) {
	amount := req.Amount
	if amount.IsZero() {
		return "", errors.New("amount cannot be zero")
	}
	if amount.IsNegative() {
		return "", errors.New("amount cannot be negative")
	}
	fromAcc, err := s.AccountRepository.FindById(req.FromID)
	if err != nil {
		return "", err
	}
	toAcc, err := s.AccountRepository.FindById(req.ToID)
	if err != nil {
		return "", err
	}
	refID := uuid.New().String()
	// 更新帳號餘額
	fromAccNewBalance := fromAcc.Balance.Sub(amount)
	toAccNewBalance := toAcc.Balance.Add(amount)

	err = s.AccountRepository.UpdateBalance(fromAcc.ID, fromAccNewBalance)
	if err != nil {
		return "", err
	}
	err = s.AccountRepository.UpdateBalance(toAcc.ID, toAccNewBalance)
	if err != nil {
		return "", err
	}
	// 寫入交易紀錄
	fromAccTx := &domain.Transaction{
		Name:   fromAcc.Name,
		Type:   1, // 提款
		Amount: amount,
		RefID:  refID,
	}
	err = s.TransactionRepository.InsertTransactions(fromAcc.ID, fromAccTx)
	if err != nil {
		return "", err
	}
	toAccTx := &domain.Transaction{
		Name:   toAcc.Name,
		Type:   2, // 存款
		Amount: amount,
		RefID:  refID,
	}
	err = s.TransactionRepository.InsertTransactions(toAcc.ID, toAccTx)
	if err != nil {
		return "", err
	}

	return refID, nil
}
