# OneGoServer Dockerfile
# 多阶段构建，优化镜像大小

# 构建阶段
FROM golang:1.23.3-alpine AS builder

# 设置构建参数
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o /app/bin/onego-server \
    cmd/server/main.go

# 运行阶段
FROM alpine:3.19

# 设置标签
LABEL maintainer="youngboygod@outlook.com"
LABEL version="${VERSION}"
LABEL description="OneGoServer - 分布式任务管理平台"

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata curl

# 创建用户
RUN addgroup -g 1001 -S onego && \
    adduser -u 1001 -S onego -G onego

# 设置时区
ENV TZ=Asia/Shanghai

# 创建必要目录
RUN mkdir -p /app/configs /app/logs /app/data && \
    chown -R onego:onego /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/onego-server /app/

# 复制配置文件
COPY --chown=onego:onego configs/ /app/configs/

# 复制启动脚本
COPY --chown=onego:onego scripts/docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

# 切换用户
USER onego

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 设置入口点
ENTRYPOINT ["/app/docker-entrypoint.sh"]

# 默认命令
CMD ["/app/onego-server"]