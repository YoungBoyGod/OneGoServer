#!/bin/bash

# OneGoServer 镜像构建脚本
# 用途: 构建Docker镜像并推送到仓库

set -e

# 配置变量
APP_NAME="onego-server"
REGISTRY="${DOCKER_REGISTRY:-docker.io}"
NAMESPACE="${DOCKER_NAMESPACE:-youngboygod}"
VERSION="${BUILD_VERSION:-latest}"
DOCKERFILE="${DOCKERFILE:-deployments/docker/Dockerfile}"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查构建依赖..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装或未在PATH中"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或未在PATH中"
        exit 1
    fi
    
    log_info "依赖检查通过"
}

# 代码质量检查
quality_check() {
    log_info "执行代码质量检查..."
    
    # Go格式化检查
    if ! gofmt -l . | grep -q '^$'; then
        log_warn "代码格式不符合标准，正在格式化..."
        gofmt -w .
    fi
    
    # Go模块整理
    go mod tidy
    
    # 静态检查 (可选)
    if command -v golangci-lint &> /dev/null; then
        log_info "运行 golangci-lint..."
        golangci-lint run ./...
    fi
    
    log_info "代码质量检查完成"
}

# 运行测试
run_tests() {
    log_info "运行单元测试..."
    
    # 运行测试并生成覆盖率报告
    go test -v -race -coverprofile=coverage.out ./...
    
    # 显示覆盖率
    if [ -f coverage.out ]; then
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        log_info "测试覆盖率: $COVERAGE"
        
        # 清理覆盖率文件
        rm -f coverage.out
    fi
    
    log_info "测试完成"
}

# 构建二进制文件
build_binary() {
    log_info "构建Go二进制文件..."
    
    # 设置构建信息
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
    
    # 构建标志
    LDFLAGS="-s -w"
    LDFLAGS="$LDFLAGS -X main.version=${VERSION}"
    LDFLAGS="$LDFLAGS -X main.buildTime=${BUILD_TIME}"
    LDFLAGS="$LDFLAGS -X main.gitCommit=${GIT_COMMIT}"
    LDFLAGS="$LDFLAGS -X main.gitBranch=${GIT_BRANCH}"
    
    # 构建
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags "$LDFLAGS" \
        -o bin/$APP_NAME \
        cmd/server/main.go
    
    log_info "二进制文件构建完成: bin/$APP_NAME"
}

# 构建Docker镜像
build_docker() {
    log_info "构建Docker镜像..."
    
    # 镜像标签
    IMAGE_TAG="$REGISTRY/$NAMESPACE/$APP_NAME:$VERSION"
    LATEST_TAG="$REGISTRY/$NAMESPACE/$APP_NAME:latest"
    
    log_info "镜像标签: $IMAGE_TAG"
    
    # 构建镜像
    docker build \
        --file $DOCKERFILE \
        --tag $IMAGE_TAG \
        --tag $LATEST_TAG \
        --build-arg VERSION=$VERSION \
        --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
        --build-arg GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
        .
    
    log_info "Docker镜像构建完成"
}

# 推送镜像
push_docker() {
    if [ "$PUSH_IMAGE" = "true" ]; then
        log_info "推送Docker镜像..."
        
        IMAGE_TAG="$REGISTRY/$NAMESPACE/$APP_NAME:$VERSION"
        LATEST_TAG="$REGISTRY/$NAMESPACE/$APP_NAME:latest"
        
        docker push $IMAGE_TAG
        
        if [ "$VERSION" != "latest" ]; then
            docker push $LATEST_TAG
        fi
        
        log_info "镜像推送完成"
    else
        log_info "跳过镜像推送 (PUSH_IMAGE != true)"
    fi
}

# 清理构建产物
cleanup() {
    log_info "清理构建产物..."
    
    rm -rf bin/
    
    # 清理Docker构建缓存 (可选)
    if [ "$CLEAN_CACHE" = "true" ]; then
        docker builder prune -f
    fi
    
    log_info "清理完成"
}

# 显示帮助信息
show_help() {
    cat << EOF
OneGoServer 构建脚本

用法: $0 [选项]

选项:
    -h, --help          显示帮助信息
    -v, --version       设置版本号 (默认: latest)
    -p, --push          构建后推送镜像
    -t, --test          运行测试
    -c, --clean         构建前清理
    --skip-quality      跳过代码质量检查
    --skip-test         跳过测试
    --clean-cache       清理Docker构建缓存

环境变量:
    BUILD_VERSION       构建版本号
    DOCKER_REGISTRY     Docker仓库地址
    DOCKER_NAMESPACE    Docker命名空间
    DOCKERFILE          Dockerfile路径
    PUSH_IMAGE          是否推送镜像 (true/false)
    CLEAN_CACHE         是否清理缓存 (true/false)

示例:
    $0                                  # 基础构建
    $0 -v v1.0.0 -p                    # 构建v1.0.0并推送
    $0 -t -c                           # 清理后运行测试并构建
    BUILD_VERSION=v1.0.0 PUSH_IMAGE=true $0  # 使用环境变量

EOF
}

# 主函数
main() {
    # 解析参数
    SKIP_QUALITY=false
    SKIP_TEST=false
    RUN_TEST=false
    CLEAN_BEFORE=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -p|--push)
                PUSH_IMAGE=true
                shift
                ;;
            -t|--test)
                RUN_TEST=true
                shift
                ;;
            -c|--clean)
                CLEAN_BEFORE=true
                shift
                ;;
            --skip-quality)
                SKIP_QUALITY=true
                shift
                ;;
            --skip-test)
                SKIP_TEST=true
                shift
                ;;
            --clean-cache)
                CLEAN_CACHE=true
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    log_info "开始构建 $APP_NAME:$VERSION"
    
    # 检查依赖
    check_dependencies
    
    # 清理 (如果需要)
    if [ "$CLEAN_BEFORE" = "true" ]; then
        cleanup
    fi
    
    # 代码质量检查
    if [ "$SKIP_QUALITY" != "true" ]; then
        quality_check
    fi
    
    # 运行测试
    if [ "$RUN_TEST" = "true" ] && [ "$SKIP_TEST" != "true" ]; then
        run_tests
    fi
    
    # 构建
    build_binary
    build_docker
    
    # 推送
    push_docker
    
    log_info "构建完成! 镜像: $REGISTRY/$NAMESPACE/$APP_NAME:$VERSION"
}

# 执行主函数
main "$@" 