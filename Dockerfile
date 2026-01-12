FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置Go代理（加速依赖下载）
ENV GOPROXY=https://goproxy.cn,direct

# 复制依赖文件
COPY go.mod ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o user-center cmd/main.go

# 使用轻量级镜像
FROM alpine:latest

# 安装必要的工具和时区数据
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/user-center .

# 创建日志目录
RUN mkdir -p logs

EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/api/health || exit 1

CMD ["./user-center"]
