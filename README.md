
---

````markdown
# Simple Bank System

本專案為 **Golang** 實作的簡易銀行系統，提供 RESTful API，功能包含：

- 建立帳戶（Create Account）
- 存款（Deposit）
- 提款（Withdraw）
- 轉帳（Transfer）
- 查詢帳戶資訊（Get Account）
- 查詢交易紀錄（Transaction Logs）
- 支援 **原子性交易**（Atomic Transaction，轉帳要嘛全部成功、要嘛全部失敗）
- 資料庫使用 PostgreSQL 進行持久化儲存
- 提供 **單元測試** 與 **整合測試**
- 使用 **Docker** 快速部署

---

## 🛠 環境需求
- Go 1.22+
- Docker（可選，用於容器化）

---

##  專案啟動方式

### 1. Clone 專案
```bash
git clone https://github.com/<your-username>/simple-bank-system.git
cd simple-bank-system
````

### 2. 本地執行

```bash
go run .
```

預設啟動在 `http://localhost:8080`

---

## 🔗 API 使用方式

### 建立帳戶

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name":"Kevin","balance":1000}'
```

### 查詢帳戶

```bash
curl http://localhost:8080/accounts/<account_id>
```

### 存款

```bash
curl -X POST http://localhost:8080/accounts/<account_id>/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount":200}'
```

### 提款

```bash
curl -X POST http://localhost:8080/accounts/<account_id>/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount":50}'
```

### 轉帳

```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"from":"<from_id>","to":"<to_id>","amount":100}'
```

### 查詢交易紀錄

```bash
curl http://localhost:8080/accounts/<account_id>/transfers
```

---

##  測試

### 執行單元測試 & 整合測試

```bash
go test ./...
```

---

##  使用 Docker

### 建立映像檔

```bash
docker build -t simple-bank-system .
```

### 執行容器

```bash
docker run -p 8080:8080 simple-bank-system
```

系統會在 `http://localhost:8080` 提供服務。

---

##  專案結構

```
simple-bank-system/
 ├── main.go          # 程式進入點
 ├── model.go         # 資料結構 (Account, TransferLog)
 ├── store.go         # 資料存取邏輯 (in-memory)
 ├── handlers.go      # RESTful API handler
 ├── handlers_test.go # 單元測試與整合測試
 ├── Dockerfile       # Docker 容器化設定
 └── go.mod           # Go modules 設定
```

---

## 📌 備註

* 本專案僅作為練習與展示用途

```

