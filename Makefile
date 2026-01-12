# Makefile for User Center Go Application

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME := user-center
DOCKER_IMAGE := user-center:latest
GO_VERSION := 1.21

## help: 显示帮助信息
help:
	@echo "用户中心 - Make 命令"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用命令:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: 编译应用
build:
	@echo "编译应用..."
	go build -ldflags="-w -s" -o $(APP_NAME) cmd/main.go
	@echo "✓ 编译完成: $(APP_NAME)"

## run: 运行应用
run:
	@echo "启动应用..."
	go run cmd/main.go

## test: 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

## clean: 清理编译产物
clean:
	@echo "清理编译产物..."
	rm -f $(APP_NAME)
	go clean
	@echo "✓ 清理完成"

## deps: 安装依赖
deps:
	@echo "安装依赖..."
	go mod download
	go mod tidy
	@echo "✓ 依赖安装完成"

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...
	@echo "✓ 格式化完成"

## vet: 代码静态检查
vet:
	@echo "代码检查..."
	go vet ./...
	@echo "✓ 检查完成"

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "✓ Docker 镜像构建完成"

## docker-up: 启动 Docker Compose 服务
docker-up:
	@echo "启动 Docker Compose 服务..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d --build; \
	else \
		docker compose up -d --build; \
	fi
	@echo "✓ 服务启动完成"
	@echo ""
	@echo "查看状态: make docker-status"
	@echo "查看日志: make docker-logs"

## docker-down: 停止 Docker Compose 服务
docker-down:
	@echo "停止 Docker Compose 服务..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	else \
		docker compose down; \
	fi
	@echo "✓ 服务已停止"

## docker-down-v: 停止服务并删除数据卷
docker-down-v:
	@echo "停止服务并删除数据卷..."
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down -v; \
	else \
		docker compose down -v; \
	fi
	@echo "✓ 服务已停止，数据卷已删除"

## docker-logs: 查看 Docker 日志
docker-logs:
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose logs -f app; \
	else \
		docker compose logs -f app; \
	fi

## docker-status: 查看 Docker 服务状态
docker-status:
	@if command -v docker-compose >/dev/null 2>&1; then \
		docker-compose ps; \
	else \
		docker compose ps; \
	fi

## docker-restart: 重启 Docker 服务
docker-restart: docker-down docker-up

## docker-exec: 进入应用容器
docker-exec:
	@docker exec -it user-center-app sh

## docker-mysql: 进入 MySQL 容器
docker-mysql:
	@docker exec -it user-center-mysql mysql -uroot -p123456 yupi

## dev: 开发模式（热重载需要安装 air）
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "请先安装 air: go install github.com/cosmtrek/air@latest"; \
		echo "或直接运行: make run"; \
	fi

## install-tools: 安装开发工具
install-tools:
	@echo "安装开发工具..."
	go install github.com/cosmtrek/air@latest
	@echo "✓ 工具安装完成"
