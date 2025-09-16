package api

import (
	"encoding/json"
	"net/http"

	"github.com/yoyo0827/simple-bank-system/internal/request"
	"github.com/yoyo0827/simple-bank-system/internal/response"
	"github.com/yoyo0827/simple-bank-system/internal/service"
)

type ApiHandler struct {
	AccountService *service.AccountService
}

// CreateAccount godoc
// @Summary 建立帳號
// @Description 建立一個新的帳號，初始餘額必須 >= 0
// @Tags 帳號相關
// @Accept json
// @Produce json
// @Param account body request.CreateAccountRequest true "Account Info"
// @Success 200 {object} response.ApiResponse
// @Router /accounts [post]
func (h *ApiHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req request.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	acc, err := h.AccountService.CreateAccount(req.Name, req.Balance)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, acc)
}

// findAccount godoc
// @Summary 查詢帳號
// @Description 根據帳號 ID 查詢帳號資訊
// @Tags 帳號相關
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} response.ApiResponse
// @Router /accounts/{id} [get]
func (h *ApiHandler) FindAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	acc, err := h.AccountService.FindAccount(id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, acc)
}

// Transaction godoc
// @Summary 交易
// @Description 對指定帳號進行存款或提款操作，金額為正數表示存款，負數表示提款
// @Tags 交易相關
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param transaction body request.TransactionRequest true "Transaction Info"
// @Success 200 {object} response.ApiResponse
// @Router /accounts/{id}/transactions [post]
func (h *ApiHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req request.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	refID, err := h.AccountService.Transaction(id, &req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, map[string]string{"ref_id": refID})
}

// Transfer godoc
// @Summary 轉帳
// @Description 由指定帳號進行轉帳操作
// @Tags 交易相關
// @Accept json
// @Produce json
// @Param transaction body request.TransferRequest true "Transfer Info"
// @Success 200 {object} response.ApiResponse
// @Router /accounts/transfer [post]
func (h *ApiHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var req request.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	refID, err := h.AccountService.Transfer(&req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, map[string]string{"ref_id": refID})
}

// TransactionDetail godoc
// @Summary 取得交易紀錄
// @Description 取得指定帳號的所有交易資訊
// @Tags 交易相關
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} response.ApiResponse
// @Router /accounts/{id}/transactions [get]
func (h *ApiHandler) FindTransactionDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	transactions, err := h.AccountService.FindAccountTransactions(id)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	response.WriteSuccess(w, http.StatusOK, transactions)
}
