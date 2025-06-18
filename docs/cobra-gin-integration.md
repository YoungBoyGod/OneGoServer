# Cobra 与 Gin 集成指南

## 概述

本文档记录了如何将 Cobra CLI 框架与 Gin Web 框架进行集成，实现命令行方式启动 HTTP 服务器。

## 修改内容

### 1. 主要架构调整

**修改前:**
- `main()` 函数直接启动 Gin 服务器
- 缺少命令行接口
- 端口硬编码为 30000

**修改后:**
- 使用 Cobra 作为主要的命令行框架
- Gin 服务器封装在 `server` 子命令中
- 支持命令行参数配置

### 2. 代码结构变化

#### 新增变量
```go
var (
    port string  // 服务器端口配置
)
```

#### 命令结构
```go
rootCmd         // 根命令: OneGoTask 服务管理工具
└── serverCmd   // 子命令: 启动 HTTP 服务器
```

#### 新增函数
- `startServer()`: 封装 Gin 服务器启动逻辑
- 改进的错误处理和日志输出

### 3. 新增功能

1. **命令行接口**
   - `./OneGoTask server`: 启动服务器
   - `./OneGoTask server -p 8080`: 指定端口启动
   - `./OneGoTask --help`: 查看帮助信息

2. **增强的路由**
   - `/ping`: 基本测试接口
   - `/health`: 健康检查接口

3. **改进的启动信息**
   - 显示服务器启动状态
   - 显示监听地址和端口
   - 提供测试链接

## 使用方法

### 构建项目
```bash
go build -o OneGoTask
```

### 启动服务器
```bash
# 使用默认端口 30000
./OneGoTask server

# 指定端口
./OneGoTask server --port 8080
./OneGoTask server -p 8080
```

### 测试接口
```bash
# 测试基本接口
curl http://localhost:30000/ping

# 健康检查
curl http://localhost:30000/health
```

## 配置文件

在 `config/server.yaml` 中添加了示例配置：

```yaml
server:
  port: 30000
  mode: release
  host: "0.0.0.0"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: ""
  database: "onegotask"

log:
  level: "info"
  file: "logs/app.log"
```

## 优势

1. **模块化**: 命令和服务器逻辑分离
2. **可扩展**: 易于添加新的子命令
3. **配置灵活**: 支持命令行参数和配置文件
4. **用户友好**: 清晰的命令行接口和帮助信息
5. **错误处理**: 完善的错误处理和日志输出

## 后续扩展建议

1. 添加配置文件读取功能
2. 添加数据库连接命令
3. 添加服务停止和重启命令
4. 添加版本信息命令
5. 添加日志配置和轮转

## 修改时间

- 创建时间: 2024年
- 最后修改: 2024年
- 修改人: AI Assistant 