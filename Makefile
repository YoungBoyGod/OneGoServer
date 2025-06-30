# =============================================================================
# OneGoServer Makefile
# =============================================================================

# 项目信息
PROJECT_NAME := OneGoServer
BINARY_NAME := onego
PACKAGE := github.com/YoungBoyGod/OneGoServer

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | cut -d' ' -f3)
BUILD_BY := $(shell whoami)@$(shell hostname)

# 构建标志
LDFLAGS := -ldflags "\
	-X '$(PACKAGE)/cmd.Version=$(VERSION)' \
	-X '$(PACKAGE)/cmd.BuildTime=$(BUILD_TIME)' \
	-X '$(PACKAGE)/cmd.GitCommit=$(GIT_COMMIT)' \
	-X '$(PACKAGE)/cmd.GitBranch=$(GIT_BRANCH)' \
	-X '$(PACKAGE)/cmd.GoVersion=$(GO_VERSION)' \
	-X '$(PACKAGE)/cmd.BuildBy=$(BUILD_BY)' \
	-w -s"

# 目录
BUILD_DIR := bin
DOCKER_DIR := deployments/docker
SCRIPTS_DIR := scripts

# Docker信息  
DOCKER_REGISTRY := youngboygod
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_TAG := $(VERSION)

# 平台信息
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# =============================================================================
# 默认目标
# =============================================================================

.PHONY: help
help: ## 显示帮助信息
	@echo "OneGoServer Build System"
	@echo "======================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make build              # 构建本地版本"
	@echo "  make build-all          # 交叉编译所有平台"
	@echo "  make docker             # 构建Docker镜像"
	@echo "  make test               # 运行测试"
	@echo "  make clean              # 清理构建文件"

.DEFAULT_GOAL := help

# =============================================================================
# 构建目标
# =============================================================================

.PHONY: build
build: clean deps ## 构建本地版本
	@echo "🚀 构建 $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-linux
build-linux: clean deps ## 构建Linux版本  
	@echo "🐧 构建Linux版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./main.go
	@echo "✅ Linux版本构建完成: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

.PHONY: build-windows
build-windows: clean deps ## 构建Windows版本
	@echo "🪟 构建Windows版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./main.go
	@echo "✅ Windows版本构建完成: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

.PHONY: build-darwin
build-darwin: clean deps ## 构建MacOS版本
	@echo "🍎 构建MacOS版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./main.go
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./main.go
	@echo "✅ MacOS版本构建完成: $(BUILD_DIR)/$(BINARY_NAME)-darwin-*"

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## 交叉编译所有平台
	@echo "🌐 所有平台构建完成!"
	@ls -la $(BUILD_DIR)/

# =============================================================================
# 开发目标
# =============================================================================

.PHONY: deps
deps: ## 安装/更新依赖
	@echo "📦 更新依赖..."
	@go mod tidy
	@go mod download
	@echo "✅ 依赖更新完成"

.PHONY: dev
dev: deps ## 开发模式运行
	@echo "🔧 开发模式启动..."
	@go run main.go server --verbose

.PHONY: run
run: build ## 运行构建的二进制文件
	@echo "🚀 运行 $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) server

# =============================================================================
# 测试目标
# =============================================================================

.PHONY: test
test: ## 运行测试
	@echo "🧪 运行测试..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "📊 生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告生成: coverage.html"

.PHONY: bench
bench: ## 运行基准测试
	@echo "⚡ 运行基准测试..."
	@go test -bench=. -benchmem ./...

# =============================================================================
# 代码质量
# =============================================================================

.PHONY: fmt
fmt: ## 格式化代码
	@echo "🎨 格式化代码..."
	@go fmt ./...
	@echo "✅ 代码格式化完成"

.PHONY: lint
lint: ## 代码静态检查
	@echo "🔍 运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint 未安装，跳过检查"; \
		echo "   安装命令: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet: ## Go vet检查
	@echo "🔎 运行 go vet..."
	@go vet ./...
	@echo "✅ Vet检查完成"

.PHONY: check
check: fmt vet lint test ## 完整代码检查流程

# =============================================================================
# Docker目标
# =============================================================================

.PHONY: docker
docker: ## 构建Docker镜像
	@echo "🐳 构建Docker镜像 $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_BRANCH=$(GIT_BRANCH) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg BUILD_BY=$(BUILD_BY) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		-f $(DOCKER_DIR)/Dockerfile .
	@echo "✅ Docker镜像构建完成: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: docker-push
docker-push: docker ## 推送Docker镜像
	@echo "📤 推送Docker镜像..."
	@docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@docker push $(DOCKER_IMAGE):latest
	@echo "✅ Docker镜像推送完成"

.PHONY: docker-run
docker-run: ## 运行Docker容器
	@echo "🚀 运行Docker容器..."
	@docker run --rm -p 8080:8080 \
		--name $(BINARY_NAME)-container \
		-e DB_HOST=host.docker.internal \
		-e REDIS_HOST=host.docker.internal \
		$(DOCKER_IMAGE):$(DOCKER_TAG) server

# =============================================================================
# 发布目标
# =============================================================================

.PHONY: release
release: clean check build-all docker ## 完整发布流程
	@echo "🚀 准备发布 $(VERSION)..."
	@mkdir -p release
	@cp $(BUILD_DIR)/* release/ 2>/dev/null || true
	@echo "✅ 发布文件准备完成: release/"
	@echo ""
	@echo "📋 发布清单:"
	@ls -la release/
	@echo ""
	@echo "🐳 Docker镜像: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: tag
tag: ## 创建Git标签
	@echo "🏷️  创建标签 $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✅ 标签创建完成: $(VERSION)"

# =============================================================================
# 清理目标
# =============================================================================

.PHONY: clean
clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -rf release
	@rm -f coverage.out coverage.html
	@echo "✅ 清理完成"

.PHONY: clean-docker
clean-docker: ## 清理Docker资源
	@echo "🐳 清理Docker资源..."
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true
	@docker system prune -f
	@echo "✅ Docker清理完成"

.PHONY: clean-all
clean-all: clean clean-docker ## 清理所有文件

# =============================================================================
# 信息目标
# =============================================================================

.PHONY: info
info: ## 显示构建信息
	@echo "📋 构建信息"
	@echo "================"
	@echo "项目名称: $(PROJECT_NAME)"
	@echo "二进制名: $(BINARY_NAME)"
	@echo "版本号:   $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git提交:  $(GIT_COMMIT)"
	@echo "Git分支:  $(GIT_BRANCH)"
	@echo "Go版本:   $(GO_VERSION)"
	@echo "构建者:   $(BUILD_BY)"
	@echo "目标平台: $(GOOS)/$(GOARCH)"
	@echo "Docker镜像: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: version
version: build ## 显示版本信息
	@./$(BUILD_DIR)/$(BINARY_NAME) version --verbose

# =============================================================================
# 工具安装
# =============================================================================

.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ 开发工具安装完成" 