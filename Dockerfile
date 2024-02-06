# 第一階段
FROM golang:1.18.1-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

#ENV GOPROXY=https://goproxy.io // go mod download 失敗請把註解拿掉

# 先將所有原始碼複製到 Docker 內
COPY . .

# 然後再執行 go mod download
RUN go mod download
RUN go build -o scannba

# 第二階段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/scannba .
# COPY --from=builder /app/static ./static
# COPY --from=builder /app/blackbox_exporter ./blackbox_exporter

EXPOSE 8080

CMD ["./scannba"]
