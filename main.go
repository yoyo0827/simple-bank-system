package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/yoyo0827/simple-bank-system/internal/api"
	"github.com/yoyo0827/simple-bank-system/internal/config"
	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/service"

	_ "github.com/yoyo0827/simple-bank-system/docs" // swagger docs
)

// @title Simple Bank System API
// @version 1.0
// @description A simple banking system implemented in Go with RESTful APIs.
// @host localhost:8080
// @BasePath /
func main() {
	// 本機開發使用 .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	// 初始化 DB
	config.InitDatabase()
	defer config.DB.Close()

	// 初始化三層
	accountRepo := &repository.AccountRepository{DB: config.DB}
	transactionRepo := &repository.TransactionRepository{DB: config.DB}
	svc := &service.AccountService{AccountRepository: accountRepo, TransactionRepository: transactionRepo}
	handler := &api.ApiHandler{AccountService: svc}
	// 啟動簡單的 server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", handler.CreateAccount)
	mux.HandleFunc("GET /accounts/{id}", handler.FindAccount)
	mux.HandleFunc("POST /accounts/{id}/transactions", handler.CreateTransaction)
	mux.HandleFunc("PUT /accounts/transfer", handler.CreateTransfer)
	mux.HandleFunc("GET /accounts/{id}/transactions", handler.FindTransactionDetail)

	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
