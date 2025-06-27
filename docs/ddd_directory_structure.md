# DDD架构目录结构详解

## 完整目录结构

```
OneGoServer/
├── cmd/                          # 应用程序入口点
│   ├── server/                   # 服务器启动程序
│   │   └── main.go              # 主程序入口
│   ├── migrate/                  # 数据库迁移工具
│   │   └── main.go              # 迁移程序
│   └── worker/                   # 后台任务程序
│       └── main.go              # 工作程序
├── config/                       # 配置文件
│   ├── server.yaml              # 服务器配置
│   ├── database.yaml            # 数据库配置
│   └── redis.yaml               # Redis配置
├── internal/                     # 私有应用程序代码
│   ├── domain/                  # 🏛️ 领域层 - 核心业务逻辑
│   │   ├── entity/              # 实体 - 具有唯一标识的业务对象
│   │   │   ├── user.go          # 用户实体
│   │   │   ├── order.go         # 订单实体
│   │   │   └── product.go       # 产品实体
│   │   ├── repository/          # 仓储接口 - 数据访问抽象
│   │   │   ├── user_repository.go    # 用户仓储接口
│   │   │   ├── order_repository.go   # 订单仓储接口
│   │   │   └── product_repository.go # 产品仓储接口
│   │   ├── service/             # 领域服务 - 复杂业务逻辑
│   │   │   ├── user_domain_service.go      # 用户领域服务
│   │   │   ├── order_domain_service.go     # 订单领域服务
│   │   │   └── pricing_service.go          # 定价服务
│   │   ├── valueobject/         # 值对象 - 无唯一标识的业务概念
│   │   │   ├── money.go         # 金额值对象
│   │   │   ├── address.go       # 地址值对象
│   │   │   └── email.go         # 邮箱值对象
│   │   ├── aggregate/           # 聚合根 - 业务一致性边界
│   │   │   ├── user_aggregate.go      # 用户聚合
│   │   │   └── order_aggregate.go     # 订单聚合
│   │   └── event/               # 领域事件
│   │       ├── user_events.go   # 用户相关事件
│   │       └── order_events.go  # 订单相关事件
│   ├── application/             # 🎯 应用层 - 用例编排
│   │   ├── service/             # 应用服务 - 用例实现
│   │   │   ├── user_app_service.go    # 用户应用服务
│   │   │   ├── order_app_service.go   # 订单应用服务
│   │   │   └── auth_service.go        # 认证服务
│   │   ├── dto/                 # 数据传输对象
│   │   │   ├── user_dto.go      # 用户DTO
│   │   │   ├── order_dto.go     # 订单DTO
│   │   │   └── common_dto.go    # 通用DTO
│   │   ├── command/             # 命令对象 - 写操作
│   │   │   ├── user_commands.go # 用户命令
│   │   │   └── order_commands.go # 订单命令
│   │   ├── query/               # 查询对象 - 读操作
│   │   │   ├── user_queries.go  # 用户查询
│   │   │   └── order_queries.go # 订单查询
│   │   └── handler/             # 命令/查询处理器
│   │       ├── command_handlers.go # 命令处理器
│   │       └── query_handlers.go   # 查询处理器
│   ├── infrastructure/          # 🔧 基础设施层 - 技术实现
│   │   ├── database/            # 数据库相关
│   │   │   ├── database.go      # 数据库连接管理
│   │   │   ├── migrations/      # 数据库迁移
│   │   │   └── seeders/         # 数据种子
│   │   ├── repository/          # 仓储实现
│   │   │   ├── user_repository_impl.go    # 用户仓储实现
│   │   │   ├── order_repository_impl.go   # 订单仓储实现
│   │   │   └── product_repository_impl.go # 产品仓储实现
│   │   ├── external/            # 外部服务集成
│   │   │   ├── payment_service.go      # 支付服务
│   │   │   ├── notification_service.go # 通知服务
│   │   │   └── sms_service.go          # 短信服务
│   │   ├── cache/               # 缓存实现
│   │   │   ├── redis_client.go  # Redis客户端
│   │   │   └── memory_cache.go  # 内存缓存
│   │   ├── logger/              # 日志实现
│   │   │   └── zap_logger.go    # Zap日志实现
│   │   └── queue/               # 消息队列
│   │       ├── rabbitmq.go      # RabbitMQ实现
│   │       └── redis_queue.go   # Redis队列实现
│   ├── interfaces/              # 🌐 接口层 - 外部交互
│   │   ├── http/                # HTTP接口
│   │   │   ├── controller/      # 控制器
│   │   │   │   ├── user_controller.go     # 用户控制器
│   │   │   │   ├── order_controller.go    # 订单控制器
│   │   │   │   └── health_controller.go   # 健康检查控制器
│   │   │   ├── middleware/      # 中间件
│   │   │   │   ├── auth_middleware.go     # 认证中间件
│   │   │   │   ├── cors_middleware.go     # CORS中间件
│   │   │   │   ├── rate_limit_middleware.go # 限流中间件
│   │   │   │   └── logging_middleware.go  # 日志中间件
│   │   │   ├── routes/          # 路由配置
│   │   │   │   ├── routes.go    # 路由注册
│   │   │   │   ├── user_routes.go    # 用户路由
│   │   │   │   └── order_routes.go   # 订单路由
│   │   │   └── response/        # 响应格式
│   │   │       ├── response.go  # 统一响应格式
│   │   │       └── error.go     # 错误响应
│   │   ├── grpc/                # gRPC接口
│   │   │   ├── server/          # gRPC服务器
│   │   │   ├── handler/         # gRPC处理器
│   │   │   └── proto/           # Protocol Buffer定义
│   │   └── websocket/           # WebSocket接口
│   │       ├── handler/         # WebSocket处理器
│   │       └── hub.go           # WebSocket连接管理
│   ├── shared/                  # 🔄 共享组件
│   │   ├── constants/           # 常量定义
│   │   │   ├── error_codes.go   # 错误码常量
│   │   │   └── status_codes.go  # 状态码常量
│   │   ├── utils/               # 工具函数
│   │   │   ├── hash.go          # 哈希工具
│   │   │   ├── jwt.go           # JWT工具
│   │   │   ├── validator.go     # 验证工具
│   │   │   └── converter.go     # 转换工具
│   │   ├── errors/              # 自定义错误
│   │   │   ├── domain_errors.go # 领域错误
│   │   │   └── app_errors.go    # 应用错误
│   │   └── types/               # 公共类型
│   │       ├── pagination.go    # 分页类型
│   │       └── sort.go          # 排序类型
│   └── config/                  # 配置管理
│       ├── config.go            # 配置结构定义
│       └── loader.go            # 配置加载器
├── pkg/                         # 公共库（可被外部项目使用）
│   ├── logger/                  # 日志包
│   ├── database/                # 数据库包
│   ├── cache/                   # 缓存包
│   └── validator/               # 验证包
├── api/                         # API文档
│   ├── openapi/                 # OpenAPI规范
│   │   └── spec.yaml            # API规范文件
│   ├── docs/                    # API文档
│   └── examples/                # API示例
├── docs/                        # 项目文档
│   ├── architecture/            # 架构文档
│   ├── deployment/              # 部署文档
│   └── development/             # 开发文档
├── scripts/                     # 脚本文件
│   ├── build.sh                 # 构建脚本
│   ├── deploy.sh                # 部署脚本
│   └── test.sh                  # 测试脚本
├── tests/                       # 测试文件
│   ├── unit/                    # 单元测试
│   ├── integration/             # 集成测试
│   └── e2e/                     # 端到端测试
├── migrations/                  # 数据库迁移文件
├── docker/                      # Docker相关文件
│   ├── Dockerfile               # Docker镜像定义
│   └── docker-compose.yml       # Docker Compose配置
├── .gitignore                   # Git忽略文件
├── go.mod                       # Go模块定义
├── go.sum                       # Go模块校验和
├── Makefile                     # Make构建文件
└── README.md                    # 项目说明
```

## 各层详细说明

### 🏛️ Domain层（领域层）

**职责**：包含核心业务逻辑，不依赖任何外部技术

#### Entity（实体）
```go
// internal/domain/entity/user.go
type User struct {
    ID       UserID    // 唯一标识
    Username string    // 业务属性
    Email    Email     // 值对象
    // ... 业务方法
    func (u *User) ChangeEmail(newEmail Email) error
}
```

#### Repository（仓储接口）
```go
// internal/domain/repository/user_repository.go
type UserRepository interface {
    Save(user *User) error
    FindByID(id UserID) (*User, error)
    FindByEmail(email Email) (*User, error)
}
```

#### Domain Service（领域服务）
```go
// internal/domain/service/user_domain_service.go
type UserDomainService struct {
    userRepo UserRepository
}

func (s *UserDomainService) RegisterUser(username, email, password string) (*User, error) {
    // 复杂业务逻辑：检查唯一性、验证规则等
}
```

#### Value Object（值对象）
```go
// internal/domain/valueobject/email.go
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    // 验证邮箱格式
}
```

### 🎯 Application层（应用层）

**职责**：编排用例流程，不包含业务规则

#### Application Service（应用服务）
```go
// internal/application/service/user_app_service.go
type UserAppService struct {
    userDomainService *domain.UserDomainService
    emailService      EmailService
}

func (s *UserAppService) RegisterUser(req RegisterRequest) (*UserResponse, error) {
    // 1. 调用领域服务注册用户
    // 2. 发送欢迎邮件
    // 3. 发布领域事件
    // 4. 返回DTO
}
```

#### DTO（数据传输对象）
```go
// internal/application/dto/user_dto.go
type RegisterRequest struct {
    Username string `json:"username" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

### 🔧 Infrastructure层（基础设施层）

**职责**：实现技术细节，对接外部系统

#### Repository Implementation（仓储实现）
```go
// internal/infrastructure/repository/user_repository_impl.go
type userRepositoryImpl struct {
    db *gorm.DB
}

func (r *userRepositoryImpl) Save(user *entity.User) error {
    return r.db.Save(user).Error
}
```

#### External Services（外部服务）
```go
// internal/infrastructure/external/email_service.go
type EmailService struct {
    smtpClient SMTPClient
}

func (s *EmailService) SendWelcomeEmail(to, username string) error {
    // 实际发送邮件逻辑
}
```

### 🌐 Interfaces层（接口层）

**职责**：处理外部请求，协议转换

#### HTTP Controller（HTTP控制器）
```go
// internal/interfaces/http/controller/user_controller.go
type UserController struct {
    userAppService *application.UserAppService
}

func (c *UserController) Register(ctx *gin.Context) {
    var req dto.RegisterRequest
    // 1. 绑定请求参数
    // 2. 调用应用服务
    // 3. 返回响应
}
```

#### Middleware（中间件）
```go
// internal/interfaces/http/middleware/auth_middleware.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT验证逻辑
    }
}
```

## DDD vs MVC 对比

### 目录结构对比

| 方面 | MVC | DDD |
|------|-----|-----|
| **组织原则** | 按技术职责分层 | 按业务领域分层 |
| **核心目录** | `/models`, `/views`, `/controllers` | `/domain`, `/application`, `/infrastructure`, `/interfaces` |
| **文件数量** | 相对较少 | 相对较多 |
| **复杂度** | 简单直观 | 复杂但清晰 |

### MVC典型结构
```
app/
├── controllers/
│   ├── UserController.go
│   └── OrderController.go
├── models/
│   ├── User.go
│   └── Order.go
├── views/
│   ├── user/
│   └── order/
└── routes/
    └── web.go
```

### DDD典型结构
```
internal/
├── domain/
│   ├── entity/
│   ├── repository/
│   └── service/
├── application/
│   ├── service/
│   └── dto/
├── infrastructure/
│   └── repository/
└── interfaces/
    └── http/
```

## 文件职责对比

| 文件类型 | MVC职责 | DDD职责 |
|----------|---------|---------|
| **Controller** | 处理HTTP请求 + 业务逻辑 | 只处理HTTP请求转换 |
| **Model** | 数据结构 + 数据库操作 | 纯业务实体，无技术依赖 |
| **Service** | 可选的业务逻辑层 | 明确分为领域服务和应用服务 |
| **Repository** | 通常不存在 | 明确的数据访问抽象 |

## 优势对比

### MVC优势
- ✅ 学习成本低
- ✅ 开发速度快
- ✅ 适合简单业务
- ✅ 团队容易理解

### DDD优势
- ✅ 业务逻辑清晰
- ✅ 易于测试
- ✅ 技术无关性
- ✅ 适合复杂业务
- ✅ 易于维护和扩展

## 选择建议

### 使用MVC的场景
- 📦 **简单CRUD应用**
- 👥 **小团队（<5人）**
- ⏰ **快速原型开发**
- 📚 **学习项目**

### 使用DDD的场景
- 🏢 **复杂业务逻辑**
- 👨‍👩‍👧‍👦 **大团队（>10人）**
- 🔄 **长期维护项目**
- 🧪 **高质量要求**
- 🔧 **技术栈可能变更**

## 实际开发建议

### 1. 渐进式采用
- 从简单的三层架构开始
- 随着复杂度增加逐步引入DDD概念
- 不要一开始就构建完整的DDD架构

### 2. 团队培训
- 确保团队理解DDD概念
- 建立清晰的编码规范
- 定期代码审查

### 3. 工具支持
- 使用代码生成器减少重复工作
- 建立项目模板
- 自动化测试覆盖

这种DDD架构设计虽然复杂，但能够很好地应对业务复杂度的增长，保持代码的可维护性和可扩展性。 