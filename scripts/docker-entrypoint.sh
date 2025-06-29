#!/bin/sh

# OneGoServer Docker 启动脚本
# 处理容器启动前的初始化工作

set -e

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

# 环境变量默认值
export SERVER_HOST="${SERVER_HOST:-0.0.0.0}"
export SERVER_PORT="${SERVER_PORT:-8080}"
export SERVER_MODE="${SERVER_MODE:-release}"
export CONFIG_FILE="${CONFIG_FILE:-configs/config.yaml}"

# 等待数据库连接
wait_for_db() {
    if [ -n "$DB_HOST" ] && [ -n "$DB_PORT" ]; then
        log_info "等待数据库连接: $DB_HOST:$DB_PORT"
        
        max_attempts=30
        attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; then
                log_info "数据库连接成功"
                return 0
            fi
            
            log_warn "数据库连接失败 (尝试 $attempt/$max_attempts)，5秒后重试..."
            sleep 5
            attempt=$((attempt + 1))
        done
        
        log_error "数据库连接超时，启动失败"
        exit 1
    fi
}

# 等待Redis连接
wait_for_redis() {
    if [ -n "$REDIS_HOST" ] && [ -n "$REDIS_PORT" ]; then
        log_info "等待Redis连接: $REDIS_HOST:$REDIS_PORT"
        
        max_attempts=30
        attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if nc -z "$REDIS_HOST" "$REDIS_PORT" 2>/dev/null; then
                log_info "Redis连接成功"
                return 0
            fi
            
            log_warn "Redis连接失败 (尝试 $attempt/$max_attempts)，5秒后重试..."
            sleep 5
            attempt=$((attempt + 1))
        done
        
        log_warn "Redis连接超时，但继续启动 (Redis为可选依赖)"
    fi
}

# 检查配置文件
check_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_warn "配置文件不存在: $CONFIG_FILE"
        
        # 如果有模板文件，复制一份
        template_file="${CONFIG_FILE}.template"
        if [ -f "$template_file" ]; then
            log_info "使用模板配置文件: $template_file"
            cp "$template_file" "$CONFIG_FILE"
        else
            log_error "配置文件和模板都不存在"
            exit 1
        fi
    fi
    
    log_info "配置文件检查通过: $CONFIG_FILE"
}

# 创建必要目录
create_directories() {
    log_info "创建必要目录..."
    
    mkdir -p logs data
    
    # 确保目录权限正确
    chmod 755 logs data
    
    log_info "目录创建完成"
}

# 数据库迁移 (如果需要)
run_migrations() {
    if [ "$RUN_MIGRATIONS" = "true" ]; then
        log_info "执行数据库迁移..."
        
        # 检查是否有迁移命令
        if [ -x "./migrate" ]; then
            ./migrate up
        elif [ -x "./onego-server" ]; then
            # 如果应用支持迁移命令
            ./onego-server migrate
        else
            log_warn "未找到迁移工具，跳过数据库迁移"
        fi
        
        log_info "数据库迁移完成"
    fi
}

# 显示环境信息
show_env_info() {
    log_info "=== OneGoServer 启动信息 ==="
    log_info "版本: ${APP_VERSION:-unknown}"
    log_info "构建时间: ${BUILD_TIME:-unknown}"
    log_info "Git提交: ${GIT_COMMIT:-unknown}"
    log_info "运行模式: $SERVER_MODE"
    log_info "监听地址: $SERVER_HOST:$SERVER_PORT"
    log_info "配置文件: $CONFIG_FILE"
    log_info "工作目录: $(pwd)"
    log_info "用户: $(whoami)"
    log_info "=============================="
}

# 信号处理
handle_signal() {
    log_info "接收到停止信号，正在优雅关闭..."
    
    # 如果有PID文件，发送信号给主进程
    if [ -f "/tmp/onego-server.pid" ]; then
        pid=$(cat /tmp/onego-server.pid)
        if kill -0 "$pid" 2>/dev/null; then
            kill -TERM "$pid"
            
            # 等待进程结束
            timeout=30
            while [ $timeout -gt 0 ] && kill -0 "$pid" 2>/dev/null; do
                sleep 1
                timeout=$((timeout - 1))
            done
            
            # 如果进程仍然存在，强制杀死
            if kill -0 "$pid" 2>/dev/null; then
                log_warn "进程未在预期时间内结束，强制终止"
                kill -KILL "$pid"
            fi
        fi
    fi
    
    log_info "服务已停止"
    exit 0
}

# 注册信号处理器
trap handle_signal TERM INT

# 主函数
main() {
    log_info "OneGoServer 容器启动中..."
    
    # 显示环境信息
    show_env_info
    
    # 检查配置文件
    check_config
    
    # 创建必要目录
    create_directories
    
    # 等待依赖服务
    wait_for_db
    wait_for_redis
    
    # 执行数据库迁移
    run_migrations
    
    log_info "启动 OneGoServer 应用..."
    
    # 如果没有传入参数，使用默认命令
    if [ $# -eq 0 ]; then
        set -- "/app/onego-server"
    fi
    
    # 启动应用，并将PID保存到文件
    exec "$@" &
    echo $! > /tmp/onego-server.pid
    
    # 等待子进程
    wait
}

# 执行主函数
main "$@" 