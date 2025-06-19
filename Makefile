.PHONY: run dev build clean test install help

# 默认目标
help:
	@echo "可用的命令："
	@echo "  dev     - 使用air启动开发模式（热重载）"
	@echo "  run     - 正常运行程序"
	@echo "  build   - 构建可执行文件"
	@echo "  clean   - 清理临时文件"
	@echo "  test    - 运行测试"
	@echo "  install - 安装air工具"

# 开发模式，使用air热重载
dev:
	@echo "启动开发模式（热重载）..."
	@mkdir -p tmp
	~/go/bin/air

# 正常运行
run:
	@echo "启动服务器..."
	go run .

# 构建
build:
	@echo "构建应用程序..."
	go build -o bin/app .

# 清理
clean:
	@echo "清理临时文件..."
	rm -rf tmp/
	rm -rf bin/
	rm -f build-errors.log

# 测试
test:
	@echo "运行测试..."
	go test ./...

# 安装air工具
install:
	@echo "安装air工具..."
	go install github.com/air-verse/air@latest

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./... 