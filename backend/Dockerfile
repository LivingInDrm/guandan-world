FROM golang:1.22-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o server main.go

# 使用轻量级镜像运行
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./server"] 