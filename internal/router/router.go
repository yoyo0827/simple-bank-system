package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/yoyo0827/simple-bank-system/internal/api"
)

func NewRouter(handler *api.ApiHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// 路由定義
	mux.HandleFunc("POST /accounts", handler.CreateAccount)
	mux.HandleFunc("GET /accounts/{id}", handler.FindAccount)
	mux.HandleFunc("POST /accounts/{id}/transactions", handler.CreateTransaction)
	mux.HandleFunc("POST /accounts/transfer", handler.CreateTransfer)
	mux.HandleFunc("GET /accounts/{id}/transactions", handler.FindTransactionDetail)

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
