package main

import (
	"log"
	"net/http"

	// godotenv
	_ "github.com/yoyo0827/simple-bank-system/docs" // swagger docs
	"github.com/yoyo0827/simple-bank-system/internal/api"
	"github.com/yoyo0827/simple-bank-system/internal/config"
)

// @title Simple Bank System API
// @version 1.0
// @description A simple banking system implemented in Go (with standard library ServeMux)
// @host localhost:8080
// @BasePath /
func main() {
	// 初始化 DB
	config.InitDatabase()
	defer config.DB.Close()
	// 啟動簡單的 server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", api.HelloHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
