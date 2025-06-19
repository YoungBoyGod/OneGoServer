# Makefile for learngo0619 project
.PHONY: help build run dev clean install test server version config

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BINARY_NAME = learngo0619
BUILD_DIR = bin
CONFIG_FILE = config/server.yaml

help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 编译项目
	@echo "编译项目..."
	@go build -o $(BINARY_NAME) .
	@echo "✅ 编译完成: $(BINARY_NAME)"

run: build ## 编译并运行服务器
	@echo "启动服务器..."
	@./$(BINARY_NAME) server

dev: ## 启动开发模式（热重载）
	@echo "启动开发模式（热重载）..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ air未安装，请运行: make install-air"; \
		exit 1; \
	fi

server: build ## 启动HTTP服务器
	@./$(BINARY_NAME) server

server-verbose: build ## 启动HTTP服务器（详细日志）
	@./$(BINARY_NAME) server --verbose

server-port: build ## 启动HTTP服务器（自定义端口）
	@read -p "请输入端口号: " port; \
	./$(BINARY_NAME) server --port $$port

version: build ## 显示版本信息
	@./$(BINARY_NAME) version

version-short: build ## 显示简短版本号
	@./$(BINARY_NAME) version --short

config-validate: build ## 验证配置文件
	@./$(BINARY_NAME) config validate

config-show: build ## 显示配置内容
	@./$(BINARY_NAME) config show

config-show-json: build ## 以JSON格式显示配置
	@./$(BINARY_NAME) config show --format json

config-show-yaml: build ## 以YAML格式显示配置
	@./$(BINARY_NAME) config show --format yaml

config-generate: build ## 生成示例配置文件
	@./$(BINARY_NAME) config generate

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

clean: ## 清理编译文件
	@echo "清理编译文件..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/
	@echo "✅ 清理完成"

install: ## 安装依赖
	@echo "安装依赖..."
	@go mod tidy
	@go mod download
	@echo "✅ 依赖安装完成"

install-air: ## 安装Air热重载工具
	@echo "安装Air热重载工具..."
	@go install github.com/air-verse/air@latest
	@echo "✅ Air安装完成"

install-tools: install-air ## 安装所有开发工具
	@echo "✅ 所有开发工具安装完成"

# 开发相关命令
fmt: ## 格式化代码
	@echo "格式化代码..."
	@go fmt ./...
	@echo "✅ 代码格式化完成"

vet: ## 静态分析
	@echo "运行静态分析..."
	@go vet ./...
	@echo "✅ 静态分析完成"

lint: fmt vet ## 代码检查（格式化+静态分析）
	@echo "✅ 代码检查完成"

# 部署相关命令
build-release: ## 编译发布版本
	@echo "编译发布版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✅ 发布版本编译完成"

# 快速开发命令
quick-start: install build server ## 快速开始（安装依赖+编译+启动服务器）

all: clean install lint test build ## 完整构建流程 