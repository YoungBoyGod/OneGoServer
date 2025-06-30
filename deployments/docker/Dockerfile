# =============================================================================
# OneGoServer Docker 多阶段构建
# =============================================================================

# 构建阶段
FROM golang:1.23.3-alpine AS builder

# 构建参数
ARG VERSION=v0.1.0
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
ARG GO_VERSION=unknown
ARG BUILD_BY=docker

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X 'github.com/YoungBoyGod/OneGoServer/cmd.Version=${VERSION}' \
              -X 'github.com/YoungBoyGod/OneGoServer/cmd.BuildTime=${BUILD_TIME}' \
              -X 'github.com/YoungBoyGod/OneGoServer/cmd.GitCommit=${GIT_COMMIT}' \
              -X 'github.com/YoungBoyGod/OneGoServer/cmd.GitBranch=${GIT_BRANCH}' \
              -X 'github.com/YoungBoyGod/OneGoServer/cmd.GoVersion=${GO_VERSION}' \
              -X 'github.com/YoungBoyGod/OneGoServer/cmd.BuildBy=${BUILD_BY}' \
              -w -s" \
    -o onego ./main.go

# 运行阶段
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S onego && \
    adduser -u 1001 -S onego -G onego

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/onego .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 复制启动脚本
COPY --from=builder /app/scripts/docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# 设置文件权限
RUN chown -R onego:onego /app

# 切换到非root用户
USER onego

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# 启动应用
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["server"]