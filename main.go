package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yoyo0827/simple-bank-system/internal/api"
	"github.com/yoyo0827/simple-bank-system/internal/config"
	"github.com/yoyo0827/simple-bank-system/internal/repository"
	"github.com/yoyo0827/simple-bank-system/internal/router"
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
	if err := godotenv.Load(); err == nil {
		log.Println("loaded .env file")
	}
	// 初始化 DB
	config.InitDatabase()
	defer config.DB.Close()

	// 初始化 Handler & Service & Repository
	accountRepo := &repository.AccountRepository{}
	transactionRepo := &repository.TransactionRepository{}
	service := &service.AccountService{
		DB:                    config.DB,
		AccountRepository:     accountRepo,
		TransactionRepository: transactionRepo,
	}
	handler := &api.ApiHandler{AccountService: service}
	// 啟動 server
	mux := router.NewRouter(handler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
