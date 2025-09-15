package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase 會在服務啟動時初始化資料庫連線
func InitDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	for i := 0; i < 10; i++ { // 1s一次，最多重試 10 次
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			if pingErr := DB.Ping(); pingErr == nil {
				log.Println(" Database connected")
				return
			} else {
				err = pingErr
			}
		}
		log.Printf(" Waiting for database... (%d/10) %v", i+1, err)
		time.Sleep(1 * time.Second)
	}

	log.Fatalf("Could not connect to database: %v", err)
}
