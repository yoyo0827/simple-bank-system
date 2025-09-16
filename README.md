
# Simple Bank System

本專案為 **Golang** 實作的簡易 bank system，提供 **RESTful API**，實現基本的帳務功能，並透過 **PostgreSQL** 進行持久化儲存。系統設計支援 原子性交易，確保交易過程中發生異常時會回滾。同時提供 單元測試 (Unit Test) 與 整合測試 (Integration Test)，並可透過 **Docker-Compose** 快速部署。
另有整合 **Swagger UI**，可即時檢視與測試 API 文件。  

- 建立帳戶（Create Account）
- 存款（Deposit）
- 提款（Withdraw）
- 轉帳（Transfer）
- 查詢帳戶資訊（Get Account）
- 查詢交易紀錄（Transaction Logs）
- 支援 **原子性交易**（Atomic Transaction）
- 資料庫使用 PostgreSQL 進行持久化儲存
- 提供 **Swagger API** 文件
- 提供 **Unit Test** 與 **Integration Test**
- 使用 **Docker-Compose** 快速部署

---

##  環境需求
- Go 1.25
- Docker & Docker Compose
- PostgreSQL 16

---

##  專案啟動方式

### 1. Clone 專案
```bash
git clone https://github.com/yoyo0827/simple-bank-system.git
cd simple-bank-system
```

### 2.使用 Docker-Compose 啟動服務

```bash
docker-compose up --build -d
## 伺服器啟動後，預設提供服務於：
## API: http://localhost:8080
## Swagger UI: http://localhost:8080/swagger/index.html 
```
---

##  API 使用方式

### 建立帳戶

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name":"Kevin","balance":1000}'
```

### 查詢帳戶

```bash
curl http://localhost:8080/accounts/<id>
```

### 存款

```bash
curl -X POST http://localhost:8080/accounts/<id>/transaction \
  -H "Content-Type: application/json" \
  -d '{"amount":200}'
```

### 提款

```bash
curl -X POST http://localhost:8080/accounts/<id>/transaction \
  -H "Content-Type: application/json" \
  -d '{"amount":-50}'
```

### 轉帳

```bash
curl -X POST http://localhost:8080/accounts/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_id":"<from_id>","to_id":"<to_id>","amount":100}'
```

### 取得交易紀錄

```bash
curl http://localhost:8080/accounts/<id>/transactions
```

---

##  測試

#### 執行單元測試 (Unit Tests)
```bash
go test ./internal/service -v
```
#### 執行整合測試 (Integration Tests)
```bash
docker-compose --profile test up -d ## 先啟動測試資料庫
go test ./test -tags=integration -v
```

---

##  專案結構

```
simple-bank-system/
 ├── main.go                     # 程式進入點 (server 啟動)
 │
 ├── internal/
 │   ├── api/                    # API handlers (RESTful endpoints)
 │   ├── domain/                 # Domain models (Account, Transaction)
 │   ├── repository/             # 資料存取層 (DB 操作, SQL 實作)
 │   ├── request/                # API 請求參數結構
 │   ├── response/               # API 回傳格式 (共用回應物件)
 │   └── service/                # 商業邏輯 (交易、轉帳、帳號管理)
 │       └── account_service_test.go  # 單元測試 (Unit Tests, 使用 sqlmock)
 │
 ├── test/
 │   └── integration_test.go     # 整合測試 (Integration Tests, 連接真實 DB)
 │
 ├── db/
 │   └── init.sql                # PostgreSQL 初始化 schema
 │
 ├── docs/                       # Swagger 文件
 │   ├── docs.go
 │   ├── swagger.json
 │   └── swagger.yaml
 │
 ├── Dockerfile                  # Docker 映像檔設定
 ├── docker-compose.yml          # Docker-Compose 設定 (DB + API + 測試DB)
 └── go.mod                      # Go Modules 設定
```



