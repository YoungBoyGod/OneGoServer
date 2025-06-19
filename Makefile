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

# 多环境服务器启动命令
server-dev: build ## 启动开发环境服务器
	@echo "🔧 启动开发环境服务器..."
	@./$(BINARY_NAME) server --env dev --verbose

server-prod: build ## 启动生产环境服务器
	@echo "🚀 启动生产环境服务器..."
	@./$(BINARY_NAME) server --env prod --verbose

server-test: build ## 启动测试环境服务器
	@echo "🧪 启动测试环境服务器..."
	@./$(BINARY_NAME) server --env test --verbose

# 环境变量方式启动
dev-env: build ## 使用环境变量启动开发环境
	@echo "🔧 使用环境变量启动开发环境..."
	@APP_ENV=dev ./$(BINARY_NAME) server --verbose

prod-env: build ## 使用环境变量启动生产环境
	@echo "🚀 使用环境变量启动生产环境..."
	@APP_ENV=prod ./$(BINARY_NAME) server --verbose

test-env: build ## 使用环境变量启动测试环境
	@echo "🧪 使用环境变量启动测试环境..."
	@APP_ENV=test ./$(BINARY_NAME) server --verbose

# 多环境编译
build-dev: ## 编译开发版本
	@echo "编译开发版本..."
	@APP_ENV=dev go build -ldflags="-X learngo0619/cmd.buildEnv=dev" -o $(BINARY_NAME)-dev .
	@echo "✅ 开发版本编译完成: $(BINARY_NAME)-dev"

build-prod: ## 编译生产版本
	@echo "编译生产版本..."
	@APP_ENV=prod go build -ldflags="-X learngo0619/cmd.buildEnv=prod -w -s" -o $(BINARY_NAME)-prod .
	@echo "✅ 生产版本编译完成: $(BINARY_NAME)-prod"

build-test: ## 编译测试版本
	@echo "编译测试版本..."
	@APP_ENV=test go build -ldflags="-X learngo0619/cmd.buildEnv=test" -o $(BINARY_NAME)-test .
	@echo "✅ 测试版本编译完成: $(BINARY_NAME)-test"

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

# 多环境配置验证
config-validate-dev: build ## 验证开发环境配置
	@echo "验证开发环境配置..."
	@./$(BINARY_NAME) config validate --env dev

config-validate-prod: build ## 验证生产环境配置
	@echo "验证生产环境配置..."
	@./$(BINARY_NAME) config validate --env prod

config-validate-test: build ## 验证测试环境配置
	@echo "验证测试环境配置..."
	@./$(BINARY_NAME) config validate --env test

config-validate-all: config-validate-dev config-validate-prod config-validate-test ## 验证所有环境配置
	@echo "✅ 所有环境配置验证完成"

test: ## 运行测试
	@echo "运行测试..."
	@go test -v ./...

clean: ## 清理编译文件
	@echo "清理编译文件..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-dev $(BINARY_NAME)-prod $(BINARY_NAME)-test
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

# 多环境发布编译
build-release-all: build-release ## 编译所有环境的发布版本
	@echo "编译所有环境的发布版本..."
	@mkdir -p $(BUILD_DIR)
	# 开发环境版本
	@APP_ENV=dev CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X learngo0619/cmd.buildEnv=dev" -o $(BUILD_DIR)/$(BINARY_NAME)-dev-linux-amd64 .
	@APP_ENV=dev CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X learngo0619/cmd.buildEnv=dev" -o $(BUILD_DIR)/$(BINARY_NAME)-dev-darwin-amd64 .
	# 生产环境版本
	@APP_ENV=prod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X learngo0619/cmd.buildEnv=prod -w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-prod-linux-amd64 .
	@APP_ENV=prod CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X learngo0619/cmd.buildEnv=prod -w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-prod-darwin-amd64 .
	@echo "✅ 所有环境发布版本编译完成"

# 快速开发命令
quick-start: install build server ## 快速开始（安装依赖+编译+启动服务器）

quick-dev: install build server-dev ## 快速开发环境启动

quick-prod: install build server-prod ## 快速生产环境启动

all: clean install lint test build ## 完整构建流程 