# 项目结构重构指南

## 概述

本文档记录了 OneGoTask 项目从单文件结构重构为标准 Go 项目布局的过程，实现了命令分离、包模块化和代码复用。

## 重构目标

### 1. 标准化项目结构
- 采用 Go 社区标准的项目布局
- 分离服务器和客户端逻辑
- 提高代码可维护性和可扩展性

### 2. 命令分离
- 服务器命令独立部署
- 客户端工具独立使用  
- 统一入口支持多种使用方式

### 3. 包模块化
- 配置管理包
- 服务器逻辑包
- 客户端功能包

## 项目结构对比

### 重构前
```
OneGoTask/
├── main.go              (包含所有逻辑)
├── config/
│   └── server.yaml
├── go.mod
└── go.sum
```

### 重构后
```
OneGoTask/
├── main.go              (统一入口)
├── cmd/
│   ├── server/
│   │   └── main.go      (服务器命令)
│   └── client/
│       └── main.go      (客户端命令)
├── pkg/
│   ├── config/
│   │   └── config.go    (配置管理)
│   ├── server/
│   │   └── server.go    (服务器逻辑)
│   └── client/
│       └── client.go    (客户端逻辑)
├── config/
│   └── server.yaml      (配置文件)
├── docs/                (文档目录)
├── go.mod
└── go.sum
```

## 包设计说明

### 1. pkg/config 配置管理包

**功能特性:**
- 配置结构体定义
- YAML 文件加载
- 命令行参数合并
- 默认值设置

**核心接口:**
```go
type Config struct {
    Server ServerConfig `yaml:"server"`
    Client ClientConfig `yaml:"client"`
    Log    LogConfig    `yaml:"log"`
}

func Load(configPath string) (*Config, error)
func (c *Config) MergeServerFlags(port string)
func (c *Config) MergeClientFlags(serverURL, outputFormat string, timeout int)
```

### 2. pkg/server 服务器包

**功能特性:**
- HTTP 服务器管理
- 路由配置
- 优雅启动和关闭
- 信号处理

**核心接口:**
```go
type Server struct {
    config *config.Config
    srv    *http.Server
    router *gin.Engine
}

func New(cfg *config.Config) *Server
func (s *Server) Start()
```

### 3. pkg/client 客户端包

**功能特性:**
- HTTP 客户端封装
- 重试机制
- 多种输出格式 (JSON, YAML, Table)
- 错误处理

**核心接口:**
```go
type Client struct {
    config     *config.Config
    httpClient *http.Client
}

func New(cfg *config.Config) *Client
func (c *Client) Ping() error
func (c *Client) Health() error
// ... 其他方法
```

## 命令结构设计

### 1. 统一入口 (main.go)

**支持的命令:**
```bash
# 服务器相关
OneGoTask server start [--port 8080] [--config config.yaml]

# 客户端相关
OneGoTask client ping [--server http://localhost:8080]
OneGoTask client health [--output yaml]
OneGoTask client config [--timeout 30]
OneGoTask client status
OneGoTask client version
OneGoTask client restart
```

### 2. 独立可执行文件

**服务器独立部署:**
```bash
# 编译服务器
go build -o server cmd/server/main.go

# 独立运行
./server --port 8080 --config config.yaml
```

**客户端独立使用:**
```bash
# 编译客户端
go build -o client cmd/client/main.go

# 独立运行
./client ping --server http://remote-server:8080
./client health --output table
```

## 配置文件增强

### 完整配置示例
```yaml
server:
  port: 30000
  mode: debug
  host: "0.0.0.0"
  shutdown_timeout: 30
  read_timeout: 30
  write_timeout: 30

client:
  server_url: "http://localhost:30000"
  timeout: 30
  retry_count: 3
  retry_delay: 1
  output_format: "json"

log:
  level: "info"
  file: "logs/app.log"
```

### 配置优先级
1. **命令行参数** (最高优先级)
2. **配置文件参数**
3. **默认配置** (最低优先级)

## 客户端功能特性

### 1. 多种输出格式

**JSON 格式 (默认):**
```json
{
  "status": "healthy",
  "service": "OneGoTask",
  "version": "1.0.0"
}
```

**YAML 格式:**
```yaml
status: healthy
service: OneGoTask
version: 1.0.0
```

**表格格式:**
```
┌─────────────────────┬─────────────────────────────────────┐
│        字段         │                 值                  │
├─────────────────────┼─────────────────────────────────────┤
│ status              │ healthy                             │
│ service             │ OneGoTask                           │
│ version             │ 1.0.0                               │
└─────────────────────┴─────────────────────────────────────┘
```

### 2. 重试机制
- 可配置重试次数
- 可配置重试延迟
- 详细的错误信息

### 3. 超时控制
- 可配置请求超时时间
- 防止长时间等待

## 使用示例

### 1. 统一入口使用

```bash
# 启动服务器
./OneGoTask server start --port 8080

# 测试连接
./OneGoTask client ping

# 检查健康状态 (YAML 格式)
./OneGoTask client health --output yaml

# 获取配置信息 (表格格式)
./OneGoTask client config --output table

# 连接远程服务器
./OneGoTask client status --server http://remote:8080
```

### 2. 独立可执行文件使用

```bash
# 编译所有可执行文件
make build

# 服务器独立部署
./server --config production.yaml

# 客户端独立使用
./client ping --server http://production-server:8080
./client health --output json --timeout 10
```

### 3. 开发环境使用

```bash
# 开发模式启动
go run main.go server start --config config/dev.yaml

# 快速测试
go run main.go client ping
go run main.go client health
```

## 编译和部署

### 1. 多种编译方式

**统一可执行文件:**
```bash
go build -o OneGoTask
```

**独立可执行文件:**
```bash
# 服务器
go build -o server cmd/server/main.go

# 客户端  
go build -o client cmd/client/main.go
```

**交叉编译:**
```bash
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o OneGoTask-linux

# Windows 64位
GOOS=windows GOARCH=amd64 go build -o OneGoTask.exe

# macOS 64位
GOOS=darwin GOARCH=amd64 go build -o OneGoTask-macos
```

### 2. 部署方式

**单机部署:**
```bash
# 复制可执行文件
cp OneGoTask /usr/local/bin/

# 启动服务
OneGoTask server start --config /etc/onegotask/config.yaml
```

**容器化部署:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o OneGoTask

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/OneGoTask .
COPY --from=builder /app/config ./config
CMD ["./OneGoTask", "server", "start"]
```

**微服务部署:**
```bash
# 服务器容器
docker build -t onegotask-server -f Dockerfile.server .

# 客户端工具容器
docker build -t onegotask-client -f Dockerfile.client .
```

## 重构优势

### 1. 代码组织优势
- **模块化**: 功能按包分离，职责清晰
- **可复用**: 包可以被其他项目引用
- **可测试**: 每个包可以独立测试
- **可维护**: 修改影响范围小

### 2. 部署灵活性
- **统一部署**: 一个可执行文件包含所有功能
- **分离部署**: 服务器和客户端可独立部署
- **轻量客户端**: 只需要客户端功能时体积更小

### 3. 开发效率
- **并行开发**: 不同模块可以并行开发
- **快速测试**: 可以独立测试各个功能
- **易于扩展**: 新功能可以独立添加

### 4. 运维友好
- **监控**: 服务器和客户端可独立监控
- **日志**: 日志可以按模块分离
- **配置**: 配置管理更加灵活

## 迁移指南

### 1. 从旧版本迁移

**命令兼容性:**
```bash
# 旧版本
./OneGoTask server

# 新版本 (兼容)
./OneGoTask server start
```

**配置文件兼容性:**
- 现有配置文件可以直接使用
- 新增客户端配置项为可选

### 2. 渐进式迁移
1. 先使用统一入口，保持现有使用习惯
2. 逐步引入客户端工具功能
3. 最后考虑分离部署

## 后续扩展

### 1. 功能扩展
- 添加更多客户端命令
- 支持插件机制
- 集成更多第三方服务

### 2. 工具链扩展
- 添加 Makefile
- 集成 CI/CD 脚本
- 添加性能测试工具

### 3. 监控扩展
- 集成 Prometheus 指标
- 添加分布式追踪
- 健康检查增强

## 修改记录

- **2024年**: 完成项目结构重构
- **功能**: cmd 目录分离、包模块化、统一入口
- **新增**: 独立客户端工具、多种输出格式、重试机制 