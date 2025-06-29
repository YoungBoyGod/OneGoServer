# OneGoServer
一个使用Go写的分布式任务管理平台

## 🚀 技术栈
- **Web框架**: Gin v1.10.1
- **配置管理**: Viper + YAML
- **日志系统**: Zap + Lumberjack
- **数据库**: 规划中 (MySQL/PostgreSQL + GORM)
- **缓存**: 规划中 (Redis)
- **消息队列**: 规划中 (Kafka)
- **认证**: 规划中 (JWT)

## 📊 实现状态

| 模块 | 状态 | 完成度 | 说明 |
|------|------|--------|------|
| 🏗️ **基础架构** | ✅ 已完成 | 90% | 服务器启动、配置管理、日志系统 |
| 🔧 **中间件** | 🔄 部分完成 | 60% | CORS、限流、日志中间件已实现 |
| 🔐 **认证授权** | ❌ 未开始 | 0% | JWT、用户认证等待实现 |
| 👥 **用户管理** | ❌ 未开始 | 0% | 注册、登录、用户CRUD |
| 📱 **设备管理** | ❌ 未开始 | 0% | 设备注册、心跳、状态管理 |
| 📋 **任务管理** | ❌ 未开始 | 0% | 任务CRUD、调度、执行 |
| 💾 **数据层** | ❌ 未开始 | 0% | 数据库模型、DAO层 |

## 📁 项目结构

| 文件路径 | 功能描述 | 实现状态 | 说明 |
|---------|----------|----------|------|
| **cmd/server.go** | 应用入口，负责全局初始化与启动 | ✅ 已实现 | Gin引擎、中间件、路由注册 |
| **config/server.yaml** | 服务器配置文件 | ✅ 已实现 | 数据库、日志、安全配置 |
| **config/server.yaml.template** | 配置模板文件 | ✅ 已实现 | 开发环境配置示例 |
| **internal/config/config.go** | 配置管理与验证 | ✅ 已实现 | Viper配置、环境变量绑定 |
| **internal/controller/health.go** | 健康检查控制器 | ✅ 已实现 | 基础健康检查接口 |
| **internal/middleware/cors.go** | 跨域配置中间件 | ✅ 已实现 | CORS头设置 |
| **internal/middleware/logger.go** | 日志中间件 | ✅ 已实现 | HTTP请求日志记录 |
| **internal/middleware/ratelimit.go** | 限流中间件 | ✅ 已实现 | 基于令牌桶的限流 |
| **api/v1/*.go** | API接口定义 | 🔄 基础实现 | 健康检查、占位接口 |
| **pkg/response/response.go** | 统一响应格式 | ✅ 已实现 | API响应标准化 |
| | | | |
| **internal/model/user.go** | 用户表ORM定义 | ❌ 待实现 | User结构体、数据库映射 |
| **internal/model/device.go** | 设备表ORM定义 | ❌ 待实现 | Device结构体、状态枚举 |
| **internal/model/task.go** | 任务表ORM定义 | ❌ 待实现 | Task结构体、状态管理 |
| **internal/service/auth.go** | 认证服务 | ❌ 待实现 | 注册/登录/JWT管理 |
| **internal/service/user.go** | 用户管理服务 | ❌ 待实现 | 用户CRUD、缓存 |
| **internal/service/device.go** | 设备管理服务 | ❌ 待实现 | 设备注册、心跳检测 |
| **internal/service/task.go** | 任务管理服务 | ❌ 待实现 | 任务调度、执行管理 |
| **internal/dao/** | 数据访问层 | ❌ 待实现 | 数据库操作封装 |

## 🚀 快速开始

### 环境要求
- Go 1.23.3+
- Git

### 1. 克隆项目
```bash
git clone https://github.com/YoungBoyGod/OneGoServer.git
cd OneGoServer
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置文件
```bash
# 复制配置模板
cp config/server.yaml.template config/server.yaml

# 编辑配置文件 (可选，使用默认配置即可启动)
# vim config/server.yaml
```

### 4. 启动服务
```bash
go run main.go
```

### 5. 验证服务
```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 预期响应
{"code":0,"message":"success","data":"OK"}
```

## 📡 API接口

### 当前可用接口
- `GET /api/v1/health` - 健康检查 ✅
- `GET /api/v1/device` - 设备接口占位 🔄
- `GET /api/v1/client` - 客户端接口占位 🔄  
- `GET /api/v1/task` - 任务接口占位 🔄
- `GET /api/v1/queue` - 队列接口占位 🔄

### 规划中的接口
- `POST /api/v1/auth/register` - 用户注册 ❌
- `POST /api/v1/auth/login` - 用户登录 ❌
- `POST /api/v1/auth/logout` - 用户登出 ❌
- `GET /api/v1/users` - 用户列表 ❌
- `POST /api/v1/devices` - 设备注册 ❌
- `PUT /api/v1/devices/:id/heartbeat` - 设备心跳 ❌
- `POST /api/v1/tasks` - 创建任务 ❌
- `GET /api/v1/tasks/:id` - 任务详情 ❌

## 🔧 配置说明

### 服务器配置
```yaml
server:
  host: "0.0.0.0"      # 监听地址
  port: 8080           # 监听端口
  mode: "debug"        # 运行模式: debug/release
  name: "OneGoServer"  # 服务名称
```

### 日志配置
```yaml
logging:
  level: "info"                    # 日志级别
  output: "both"                   # 输出方式: console/file/both
  file_path: "logs/app_{timestamp}.log"  # 日志文件路径
  max_size: 100                    # 单文件最大大小(MB)
  max_backups: 10                  # 备份文件数量
  max_age: 30                      # 保留天数
```

### 安全配置
```yaml
security:
  cors:
    allow_origins: ["*"]
    allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  rate_limit:
    enabled: true
    requests_per_minute: 60
    burst: 10
```

## 🗺️ 开发路线图

### Phase 1: 数据层 (规划中)
- [ ] 数据库连接与ORM集成
- [ ] 用户、设备、任务模型定义
- [ ] 数据迁移脚本
- [ ] DAO层实现

### Phase 2: 业务层 (规划中)  
- [ ] JWT认证系统
- [ ] 用户管理服务
- [ ] 设备管理服务
- [ ] 任务管理服务

### Phase 3: 增强功能 (规划中)
- [ ] Redis缓存集成
- [ ] Kafka消息队列
- [ ] 监控与指标
- [ ] 单元测试

### Phase 4: 部署与运维 (规划中)
- [ ] Docker镜像
- [ ] Kubernetes部署
- [ ] CI/CD流水线
- [ ] 文档完善

## 🛠️ 开发指南

### 添加新接口
1. 在 `api/v1/` 目录下定义接口函数
2. 在 `cmd/server.go` 中注册路由
3. 实现对应的service层逻辑（待service层完成后）

### 添加中间件
1. 在 `internal/middleware/` 目录下实现中间件
2. 在 `cmd/server.go` 的 `setupMiddleware()` 中注册

### 修改配置
1. 更新 `internal/config/config.go` 中的结构体
2. 更新 `config/server.yaml.template` 配置模板
3. 在 `bindEnvironmentVariables()` 中添加环境变量绑定

## 📞 联系方式

- 项目地址: https://github.com/YoungBoyGod/OneGoServer
- 问题反馈: Issues页面

---

**当前版本**: v0.1.0 (开发中)  
**最后更新**: 2024-12-29
