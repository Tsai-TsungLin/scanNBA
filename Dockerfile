# 第一階段：建置
FROM golang:1.21-alpine AS builder

# 安裝必要工具
RUN apk add --no-cache git

# 設定工作目錄
WORKDIR /app

# 複製 go mod 檔案
COPY go.mod go.sum ./
RUN go mod download

# 複製所有原始碼
COPY . .

# 編譯應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nba-scanner .

# 第二階段：運行
FROM alpine:latest

# 安裝 CA 證書（用於 HTTPS 請求）和時區資料
RUN apk --no-cache add ca-certificates tzdata

# 設定時區為台北
ENV TZ=Asia/Taipei

WORKDIR /root/

# 從 builder 階段複製編譯好的二進制檔案
COPY --from=builder /app/nba-scanner .

# 暴露端口
EXPOSE 8080

# 設定健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# 執行應用程式
CMD ["./nba-scanner", "--server", "--port", "8080"]
