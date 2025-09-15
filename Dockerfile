
FROM golang:1.25 AS builder
WORKDIR /app
# 複製 go.mod / go.sum 並安裝依賴
COPY go.mod go.sum ./
RUN go mod download
# 複製專案程式碼
COPY . .
# 編譯成可執行檔
RUN go build -o bank-server ./cmd/server
FROM debian:bookworm-slim
WORKDIR /app
# 複製編譯好的 binary
COPY --from=builder /app/bank-server .
# 複製 Swagger docs
COPY --from=builder /app/docs ./docs
# server.port
ENV SERVER_PORT=8080

# 對外開放 8080 port
EXPOSE 8080

# 啟動服務
CMD ["./bank-server"]