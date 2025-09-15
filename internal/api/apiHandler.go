package api

import (
	"encoding/json"
	"net/http"

	"github.com/yoyo0827/simple-bank-system/internal/request"
	"github.com/yoyo0827/simple-bank-system/internal/service"
)

type ApiHandler struct {
	AccountService *service.AccountService
}

// CreateAccount godoc
// @Summary 建立帳號
// @Description 建立一個新的帳號，初始餘額必須 >= 0
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body request.CreateAccountRequest true "Account Info"
// @Success 200 {object} domain.Account
// @Router /accounts [post]
func (h *ApiHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req request.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.AccountService.CreateAccount(req.Name, req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(acc)
}

// findAccount godoc
// @Summary 查詢帳號
// @Description 根據帳號 ID 查詢帳號資訊
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} domain.Account
// @Router /accounts/{id} [get]
func (h *ApiHandler) FindAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	acc, err := h.AccountService.FindAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(acc)
}

// Transaction godoc
// @Summary 進行交易
// @Description 對指定帳號進行存款或提款操作，金額為正數表示存款，負數表示提款
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param transaction body request.TransactionRequest true "Transaction Info"
// @Success 200 {object} domain.Account
// @Router /accounts/{id}/transactions [post]
func (h *ApiHandler) Transaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	acc, err := h.AccountService.Transaction(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(acc)
}

// Transfer godoc
// @Summary 進行轉帳
// @Description 由指定帳號進行轉帳操作
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body request.TransferRequest true "Transfer Info"
// @Success 200 {object} domain.Account
// @Router /accounts/transfer [put]
func (h *ApiHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req request.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	refID, err := h.AccountService.Transfer(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 或 200 都行
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"ref_id": refID,
	})
}
