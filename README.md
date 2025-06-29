# OneGoServer
一个使用Go语言开发的分布式任务管理平台，采用DDD(领域驱动设计)架构

## 🚀 技术栈
- **Web框架**: Gin v1.10.1
- **配置管理**: Viper + YAML
- **日志系统**: Zap + Lumberjack
- **数据库**: PostgreSQL + GORM v1.30.0
- **缓存**: Redis v8
- **消息队列**: Kafka
- **认证**: JWT + bcrypt
- **架构模式**: DDD (领域驱动设计)
- **容器化**: Docker + Multi-stage build

## 📊 实现状态总览

| 功能模块 | 状态 | 完成度 | 说明 |
|---------|------|--------|------|
| 🏗️ **基础架构** | ✅ 已完成 | 95% | DDD架构、配置管理、日志系统 |
| 🔧 **中间件系统** | ✅ 已完成 | 85% | CORS、限流、日志中间件 |
| 📡 **API层** | 🔄 部分完成 | 60% | 完整的DTO定义、接口规范 |
| 🏛️ **领域层** | 🔄 部分完成 | 70% | 实体模型、业务用例、仓储接口 |
| 💾 **数据层** | 🔄 规划完成 | 30% | 数据库迁移、DAO接口设计 |
| 🔐 **认证授权** | ❌ 待实现 | 0% | JWT认证、权限管理 |
| 👥 **用户管理** | ❌ 待实现 | 0% | 用户CRUD、角色管理 |
| 📱 **设备管理** | ❌ 待实现 | 0% | 设备注册、心跳、状态管理 |
| 📋 **任务管理** | ❌ 待实现 | 0% | 任务CRUD、调度、执行 |
| 🚨 **告警系统** | 🔄 设计完成 | 40% | 告警模型、业务逻辑、数据库迁移 |
| 🚀 **部署运维** | 🔄 基础完成 | 50% | Docker化、构建脚本 |

## 📁 DDD架构目录结构

### 🟢 已实现的核心架构文件

```
OneGoServer/
├── 📁 cmd/                          # 应用启动层
│   └── server/
│       └── main.go                  # ✅ 启动入口：依赖注入，服务启动
├── 📁 configs/                      # 配置管理层  
│   ├── config.yaml                  # ✅ 默认配置
│   ├── config.development.yaml      # ✅ 开发环境配置
│   └── config.production.yaml       # ✅ 生产环境配置
├── 📁 api/                          # API接口定义层
│   └── v1/
│       ├── user.go                  # ✅ 用户API DTO定义
│       ├── device.go                # ✅ 设备API DTO定义
│       ├── task.go                  # ✅ 任务API DTO定义
│       └── alert.go                 # ✅ 告警API DTO定义
├── 📁 internal/                     # 内部业务层 (DDD核心)
│   ├── 📁 biz/                      # 🎯 领域层 (业务核心)
│   │   ├── models.go                # ✅ 实体定义：User/Device/Task/Alert
│   │   ├── user.go                  # ✅ 用户领域逻辑 + UserRepo接口
│   │   ├── device.go                # ✅ 设备领域逻辑 + DeviceRepo接口  
│   │   ├── task.go                  # ❌ 任务领域逻辑 (待实现)
│   │   └── alert.go                 # ✅ 告警领域逻辑 + AlertRepo接口
│   ├── 📁 data/                     # 🎯 数据层 (持久化)
│   │   ├── dao/                     # DAO数据访问对象
│   │   │   ├── user_dao.go          # ❌ 用户数据访问 (待实现)
│   │   │   ├── device_dao.go        # ❌ 设备数据访问 (待实现)
│   │   │   ├── task_dao.go          # ❌ 任务数据访问 (待实现)
│   │   │   └── alert_dao.go         # ❌ 告警数据访问 (待实现)
│   │   ├── postgres/
│   │   │   └── postgres.go          # ❌ GORM初始化 (待实现)
│   │   ├── redis/
│   │   │   └── redis.go             # ❌ Redis客户端 (待实现)
│   │   └── kafka/
│   │       ├── producer.go          # ❌ Kafka生产者 (待实现)
│   │       └── consumer.go          # ❌ Kafka消费者 (待实现)
│   ├── 📁 service/                  # 🎯 应用层 (业务用例协调)
│   │   ├── auth.go                  # ❌ 认证服务 (待实现)
│   │   ├── user_service.go          # ❌ 用户应用服务 (待实现)
│   │   ├── device_service.go        # ❌ 设备应用服务 (待实现)
│   │   ├── task_service.go          # ❌ 任务应用服务 (待实现)
│   │   └── alert_service.go         # ❌ 告警应用服务 (待实现)
│   └── 📁 server/                   # 🎯 接口层 (HTTP/gRPC)
│       ├── router/
│       │   ├── router.go            # ❌ Gin路由注册 (待实现)
│       │   ├── user_handler.go      # ❌ 用户HTTP处理器 (待实现)
│       │   ├── device_handler.go    # ❌ 设备HTTP处理器 (待实现)
│       │   ├── task_handler.go      # ❌ 任务HTTP处理器 (待实现)
│       │   └── alert_handler.go     # ❌ 告警HTTP处理器 (待实现)
│       └── grpc/
│           └── svc.go               # ❌ gRPC服务实现 (待实现)
├── 📁 pkg/                          # 公共工具包
│   ├── api/dto/                     # DTO数据传输对象
│   ├── errors/                      # 错误码定义
│   ├── util/                        # 通用工具
│   └── validator/                   # 参数验证
├── 📁 scripts/                      # 构建运维脚本
│   ├── migrations/
│   │   └── 20250701_add_alerts_table.sql  # ✅ 告警表迁移文件
│   ├── build.sh                     # ✅ 完整的镜像构建脚本
│   └── docker-entrypoint.sh         # ✅ Docker启动脚本
├── 📁 deployments/                  # 部署配置
│   ├── docker/
│   │   └── Dockerfile               # ✅ 多阶段构建 Dockerfile
│   └── k8s/                         # ❌ K8s配置 (待实现)
├── 📁 test/                         # 测试文件
│   ├── unit/                        # 单元测试
│   └── integration/                 # 集成测试
└── 📁 docs/                         # 项目文档
```

### 🎯 DDD架构层次说明

| 层次 | 目录 | 职责 | 依赖方向 | 状态 |
|------|------|------|----------|------|
| **接口层** | `internal/server/` | HTTP/gRPC接口，路由处理 | → 应用层 | 🔄 框架完成 |
| **应用层** | `internal/service/` | 业务用例协调，事务管理 | → 领域层 | ❌ 待实现 |
| **领域层** | `internal/biz/` | 核心业务逻辑，实体定义 | 不依赖其他层 | ✅ 大部分完成 |
| **数据层** | `internal/data/` | 数据持久化，外部接口 | ← 领域层接口 | 🔄 接口设计完成 |

## 🎯 详细功能需求表

### 1. 用户管理模块 (37个API接口)
| 功能分类 | API端点 | 方法 | 状态 | 领域模型 | DTO定义 |
|---------|---------|------|------|----------|---------|
| **认证授权** | `/api/v1/auth/register` | POST | ❌ 待实现 | ✅ User实体 | ✅ UserRegisterReq |
| **认证授权** | `/api/v1/auth/login` | POST | ❌ 待实现 | ✅ User实体 | ✅ UserLoginReq |
| **认证授权** | `/api/v1/auth/logout` | POST | ❌ 待实现 | ✅ User实体 | - |
| **用户管理** | `/api/v1/users` | GET | ❌ 待实现 | ✅ User实体 | ✅ UserListReq |
| **用户管理** | `/api/v1/users` | POST | ❌ 待实现 | ✅ User实体 | ✅ UserRegisterReq |
| **用户管理** | `/api/v1/users/{id}` | GET | ❌ 待实现 | ✅ User实体 | ✅ UserResp |
| **用户管理** | `/api/v1/users/{id}` | PUT | ❌ 待实现 | ✅ User实体 | ✅ UserUpdateReq |
| **用户管理** | `/api/v1/users/{id}` | DELETE | ❌ 待实现 | ✅ User实体 | - |
| **密码管理** | `/api/v1/users/{id}/password` | PUT | ❌ 待实现 | ✅ User实体 | ✅ ChangePasswordReq |

### 2. 设备管理模块 (18个API接口)
| 功能分类 | API端点 | 方法 | 状态 | 领域模型 | DTO定义 |
|---------|---------|------|------|----------|---------|
| **设备注册** | `/api/v1/devices` | POST | ❌ 待实现 | ✅ Device实体 | ✅ DeviceRegisterReq |
| **设备管理** | `/api/v1/devices` | GET | ❌ 待实现 | ✅ Device实体 | ✅ DeviceListReq |
| **设备管理** | `/api/v1/devices/{id}` | GET | ❌ 待实现 | ✅ Device实体 | ✅ DeviceResp |
| **设备管理** | `/api/v1/devices/{id}` | PUT | ❌ 待实现 | ✅ Device实体 | ✅ DeviceUpdateReq |
| **设备管理** | `/api/v1/devices/{id}` | DELETE | ❌ 待实现 | ✅ Device实体 | - |
| **设备心跳** | `/api/v1/devices/{id}/heartbeat` | PUT | ❌ 待实现 | ✅ Device实体 | ✅ DeviceHeartbeatReq |
| **设备状态** | `/api/v1/devices/{id}/status` | PUT | ❌ 待实现 | ✅ Device实体 | ✅ DeviceStatusUpdateReq |
| **设备统计** | `/api/v1/devices/stats` | GET | ❌ 待实现 | ✅ Device实体 | ✅ DeviceStatsResp |

### 3. 任务管理模块 (24个API接口)
| 功能分类 | API端点 | 方法 | 状态 | 领域模型 | DTO定义 |
|---------|---------|------|------|----------|---------|
| **任务创建** | `/api/v1/tasks` | POST | ❌ 待实现 | ✅ Task实体 | ✅ TaskCreateReq |
| **任务管理** | `/api/v1/tasks` | GET | ❌ 待实现 | ✅ Task实体 | ✅ TaskListReq |
| **任务管理** | `/api/v1/tasks/{id}` | GET | ❌ 待实现 | ✅ Task实体 | ✅ TaskResp |
| **任务管理** | `/api/v1/tasks/{id}` | PUT | ❌ 待实现 | ✅ Task实体 | ✅ TaskUpdateReq |
| **任务管理** | `/api/v1/tasks/{id}` | DELETE | ❌ 待实现 | ✅ Task实体 | - |
| **任务控制** | `/api/v1/tasks/{id}/cancel` | POST | ❌ 待实现 | ✅ Task实体 | - |
| **任务控制** | `/api/v1/tasks/{id}/retry` | POST | ❌ 待实现 | ✅ Task实体 | - |
| **任务日志** | `/api/v1/tasks/{id}/logs` | GET | ❌ 待实现 | ✅ Task实体 | ✅ TaskLogResp |

### 4. 告警系统模块 🆕 (15个API接口)
| 功能分类 | API端点 | 方法 | 状态 | 领域模型 | DTO定义 |
|---------|---------|------|------|----------|---------|
| **告警创建** | `/api/v1/alerts` | POST | ❌ 待实现 | ✅ Alert实体 | ✅ AlertCreateReq |
| **告警管理** | `/api/v1/alerts` | GET | ❌ 待实现 | ✅ Alert实体 | ✅ AlertListReq |
| **告警管理** | `/api/v1/alerts/{id}` | GET | ❌ 待实现 | ✅ Alert实体 | ✅ AlertResp |
| **告警处理** | `/api/v1/alerts/{id}/acknowledge` | POST | ❌ 待实现 | ✅ Alert实体 | ✅ AlertAckReq |
| **告警处理** | `/api/v1/alerts/{id}/resolve` | POST | ❌ 待实现 | ✅ Alert实体 | - |
| **告警处理** | `/api/v1/alerts/{id}/close` | POST | ❌ 待实现 | ✅ Alert实体 | - |
| **告警统计** | `/api/v1/alerts/stats` | GET | ❌ 待实现 | ✅ Alert实体 | ✅ AlertStatsResp |

## 🗺️ 开发路线图 (已更新)

### Phase 1: 数据层建设 (1-2周) 🔥 高优先级
- [x] **DDD架构设计**: 完整的目录结构和分层设计
- [x] **领域模型定义**: User/Device/Task/Alert实体 + 枚举
- [x] **API接口规范**: 37个用户 + 18个设备 + 24个任务 + 15个告警接口DTO
- [x] **数据库迁移**: alerts表结构设计和SQL迁移文件
- [ ] **GORM集成**: PostgreSQL连接池 + 自动迁移
- [ ] **DAO层实现**: 4个核心实体的CRUD操作
- [ ] **Redis集成**: 缓存客户端 + 分布式锁

### Phase 2: 应用层实现 (2-3周) 🔥 高优先级  
- [ ] **认证服务**: JWT生成验证 + bcrypt密码加密
- [ ] **用户应用服务**: 注册登录 + 用户管理 + 密码管理
- [ ] **设备应用服务**: 设备注册 + 心跳处理 + 状态管理
- [ ] **任务应用服务**: 任务CRUD + 状态流转 + 重试机制
- [ ] **告警应用服务**: 告警生成 + 确认处理 + 统计分析

### Phase 3: 接口层实现 (1-2周) 🟡 中优先级
- [ ] **HTTP路由**: Gin路由注册 + 中间件集成  
- [ ] **HTTP处理器**: 4个核心模块的处理器实现
- [ ] **参数验证**: 请求参数验证 + 错误处理
- [ ] **响应封装**: 统一响应格式 + 错误码管理
- [ ] **gRPC接口**: 高性能内部通信接口 (可选)

### Phase 4: 高级功能 (2-3周) 🟡 中优先级
- [ ] **Kafka集成**: 异步任务处理 + 事件驱动
- [ ] **任务调度**: 定时任务 + 延时任务 + 任务队列
- [ ] **实时推送**: WebSocket + 设备状态推送
- [ ] **监控告警**: 系统监控 + 业务告警 + 通知推送
- [ ] **性能优化**: 数据库优化 + 缓存策略 + 并发处理

### Phase 5: 运维部署 (1-2周) 🟠 低优先级
- [x] **Docker化**: 多阶段构建 + 最小镜像 + 健康检查
- [x] **构建脚本**: 完整的CI/CD构建脚本
- [ ] **K8s部署**: Deployment + Service + ConfigMap
- [ ] **Helm Chart**: 参数化部署 + 环境管理
- [ ] **监控体系**: Prometheus + Grafana + 链路追踪

## 🚀 快速开始

### 环境要求
- Go 1.23.3+
- PostgreSQL 13+
- Redis 6+
- Docker & Docker Compose
- Git

### 1. 克隆项目
```bash
git clone https://github.com/YoungBoyGod/OneGoServer.git
cd OneGoServer
```

### 2. 环境配置
```bash
# 复制配置文件
cp configs/config.yaml.template configs/config.yaml

# 启动依赖服务 (PostgreSQL + Redis + Kafka)
docker-compose up -d postgres redis kafka

# 等待服务启动完成
sleep 10
```

### 3. 数据库初始化
```bash
# 执行数据库迁移
psql -h localhost -U onegouser -d onegodb -f scripts/migrations/20250701_add_alerts_table.sql
```

### 4. 启动应用
```bash
# 安装依赖
go mod download

# 启动开发服务器
go run cmd/server/main.go
```

### 5. Docker部署 (推荐)
```bash
# 构建镜像
./scripts/build.sh -v v0.1.0

# 启动服务
docker run -p 8080:8080 --name onego-server \
  -e DB_HOST=host.docker.internal \
  -e REDIS_HOST=host.docker.internal \
  youngboygod/onego-server:v0.1.0
```

### 6. 验证部署
```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 预期响应
{"code":200,"message":"success","data":"OK"}
```

## 📡 当前可用API

| 端点 | 方法 | 描述 | 状态 | 完成度 |
|------|------|------|------|--------|
| `/api/v1/health` | GET | 健康检查 | ✅ 可用 | 100% |
| `/api/v1/device` | GET | 设备接口占位 | 🔄 开发中 | 20% |
| `/api/v1/client` | GET | 客户端接口占位 | 🔄 开发中 | 20% |
| `/api/v1/task` | GET | 任务接口占位 | 🔄 开发中 | 20% |
| `/api/v1/queue` | GET | 队列接口占位 | 🔄 开发中 | 20% |

## 🔧 配置说明

### 多环境配置支持
```yaml
# 配置文件优先级
configs/config.yaml              # 基础配置
configs/config.development.yaml  # 开发环境覆盖  
configs/config.production.yaml   # 生产环境覆盖
```

### 核心配置项
```yaml
server:
  host: "0.0.0.0"
  port: "8080"
  mode: "debug"  # debug/release/test

database:
  type: "postgres"
  host: "localhost" 
  port: 5432
  username: "onegouser"
  password: "onegopass"
  dbname: "onegodb"
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: "localhost"
  port: 6379
  db: 0
  pool_size: 10

kafka:
  brokers: ["localhost:9092"]
  topics:
    task_events: "task-events"
    device_events: "device-events"
    alert_events: "alert-events"

jwt:
  secret: "your-super-secret-jwt-key"
  expire_hours: 24

logging:
  level: "info"
  output: "both"  # console/file/both
  file_path: "logs/app_{timestamp}.log"
```

## 🛠️ 开发指南

### DDD分层架构开发规范

#### 1. 领域层开发 (`internal/biz/`)
```go
// 实体定义 - 核心业务对象
type User struct {
    ID       uint   `gorm:"primarykey"`
    Username string `gorm:"unique;not null"`
    // ...
}

// 仓储接口 - 数据访问抽象
type UserRepo interface {
    CreateUser(user *User) error
    GetUserByID(id uint) (*User, error)
    // ...
}

// 用例实现 - 业务逻辑封装
type UserUsecase struct {
    repo UserRepo
}
```

#### 2. 数据层开发 (`internal/data/`)
```go
// DAO实现 - 仓储接口具体实现
type userDAO struct {
    db *gorm.DB
}

func (d *userDAO) CreateUser(user *biz.User) error {
    return d.db.Create(user).Error
}
```

#### 3. 应用层开发 (`internal/service/`)
```go
// 应用服务 - 协调多个用例
type UserService struct {
    userUC   *biz.UserUsecase
    deviceUC *biz.DeviceUsecase
}
```

#### 4. 接口层开发 (`internal/server/`)
```go
// HTTP处理器 - 请求响应处理
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req v1.UserRegisterReq
    // 参数绑定 -> 调用应用服务 -> 返回响应
}
```

### 添加新功能模块的步骤

1. **领域层**: 在 `internal/biz/` 定义实体、仓储接口、用例
2. **API层**: 在 `api/v1/` 定义请求响应DTO
3. **数据层**: 在 `internal/data/dao/` 实现仓储接口  
4. **应用层**: 在 `internal/service/` 实现应用服务
5. **接口层**: 在 `internal/server/router/` 实现HTTP处理器
6. **测试**: 在 `test/` 编写单元测试和集成测试

### 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 Go 命名约定
- 编写完整的单元测试 (目标覆盖率 80%+)
- 添加必要的注释和文档
- 使用 `golangci-lint` 进行静态检查

## 📈 项目统计

| 指标 | 数量 | 说明 |
|------|------|------|
| **Go文件** | 12个 | 核心业务逻辑文件 |
| **API接口** | 94个 | 完整的REST API设计 |
| **数据模型** | 4个 | User/Device/Task/Alert |
| **配置文件** | 3个 | 多环境配置支持 |
| **Docker文件** | 3个 | 镜像构建+部署脚本 |
| **迁移文件** | 1个 | 数据库结构变更 |
| **代码行数** | 2000+ | 包含注释和文档 |

## 📞 联系方式

- **项目地址**: https://github.com/YoungBoyGod/OneGoServer
- **问题反馈**: [Issues页面](https://github.com/YoungBoyGod/OneGoServer/issues)
- **开发文档**: [docs/](./docs/)
- **架构设计**: 采用DDD分层架构，确保代码的可维护性和可扩展性

---

**当前版本**: v0.1.0 (开发中)  
**架构模式**: DDD (领域驱动设计) + Clean Architecture  
**开发进度**: 基础架构完成，业务逻辑开发中  
**下个里程碑**: Phase 1 数据层建设完成  
**最后更新**: 2024-12-29
