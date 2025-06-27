# OneGoServer 业务分层架构设计

## 设计理念

本项目采用了**领域驱动设计(DDD)**和**洋葱架构(Clean Architecture)**的设计思想，将业务逻辑按照职责进行清晰的分层。

## 核心设计原则

### 1. 依赖倒置原则
- 高层模块不依赖低层模块，两者都依赖于抽象
- 抽象不依赖于具体实现，具体实现依赖于抽象

### 2. 单一职责原则
- 每个层次都有明确的职责边界
- 每个包、类、函数都只做一件事

### 3. 开闭原则
- 对扩展开放，对修改封闭
- 通过接口和依赖注入实现可扩展性

## 分层架构详解

### 📁 Domain 层（领域层）- 核心业务逻辑
```
internal/domain/
├── entity/          # 领域实体
├── repository/      # 仓储接口
└── service/         # 领域服务
```

**职责：**
- 定义核心业务实体和业务规则
- 不依赖任何外部技术细节
- 包含纯粹的业务逻辑

**优势：**
- 业务逻辑独立，可复用
- 测试简单，不依赖外部系统
- 符合业务专家的思维模型

**示例：**
```go
// entity/user.go - 用户实体，包含业务属性和约束
type User struct {
    ID       uint   
    Username string  // 业务约束：唯一
    Email    string  // 业务约束：唯一，邮箱格式
    Password string  // 业务约束：加密存储
}

// service/user_service.go - 领域服务，处理业务规则
func (s *UserService) CreateUser(user *entity.User) error {
    // 业务规则：检查用户名唯一性
    // 业务规则：密码加密
}
```

### 📁 Application 层（应用层）- 应用逻辑协调
```
internal/application/
├── service/         # 应用服务
└── dto/            # 数据传输对象
```

**职责：**
- 协调领域对象完成应用功能
- 定义应用程序的用例
- 处理事务边界

**优势：**
- 将复杂的业务流程拆分为简单的步骤
- 提供稳定的API给接口层
- 便于测试和维护

**示例：**
```go
// service/user_app_service.go - 应用服务，协调多个领域服务
func (s *UserAppService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
    // 1. 转换DTO为领域对象
    // 2. 调用领域服务处理业务逻辑
    // 3. 转换领域对象为响应DTO
}
```

### 📁 Infrastructure 层（基础设施层）- 技术实现
```
internal/infrastructure/
├── database/        # 数据库连接和配置
└── repository/      # 仓储具体实现
```

**职责：**
- 实现领域层定义的接口
- 处理外部系统集成（数据库、消息队列等）
- 技术细节的具体实现

**优势：**
- 技术实现可替换
- 领域层不依赖具体技术
- 便于技术栈升级

**示例：**
```go
// repository/user_repository_impl.go - 仓储实现
type userRepositoryImpl struct {
    db *gorm.DB  // 具体的技术实现
}

func (r *userRepositoryImpl) Create(user *entity.User) error {
    return r.db.Create(user).Error  // GORM具体实现
}
```

### 📁 Interfaces 层（接口层）- 外部交互
```
internal/interfaces/
└── http/
    ├── controller/  # HTTP控制器
    ├── response/    # 响应格式
    └── routes/      # 路由配置
```

**职责：**
- 处理外部请求（HTTP、gRPC等）
- 参数验证和格式转换
- 错误处理和响应格式化

**优势：**
- 接口协议可替换（HTTP、gRPC、GraphQL）
- 统一的错误处理和响应格式
- 清晰的API边界

**示例：**
```go
// controller/user_controller.go - HTTP控制器
func (c *UserController) Register(ctx *gin.Context) {
    // 1. 参数验证和绑定
    // 2. 调用应用服务
    // 3. 格式化响应
}
```

## 为什么这样设计？

### 1. **可测试性**
```
Domain层   -> 纯业务逻辑，单元测试简单
Application层 -> Mock依赖，集成测试清晰  
Infrastructure层 -> 可单独测试技术实现
Interfaces层 -> API测试和端到端测试
```

### 2. **可维护性**
- 职责清晰，修改影响范围可控
- 依赖方向单一，避免循环依赖
- 代码结构预期一致，降低理解成本

### 3. **可扩展性**
- 新增功能只需要在对应层添加代码
- 技术栈变更只影响Infrastructure层
- 业务逻辑变更只影响Domain层

### 4. **团队协作**
```
前端开发  -> 关注Interfaces层API设计
后端开发  -> 关注Application和Domain层
运维/DBA  -> 关注Infrastructure层
产品经理  -> 关注Domain层业务规则
```

## 数据流向

```
HTTP请求 
  ↓
Controller (参数验证)
  ↓  
Application Service (协调业务流程)
  ↓
Domain Service (执行业务规则)
  ↓
Repository Interface (抽象数据访问)
  ↓
Repository Implementation (具体数据库操作)
  ↓
Database
```

## 依赖关系

```
Interfaces → Application → Domain ← Infrastructure
```

- **Interfaces层** 依赖 Application层
- **Application层** 依赖 Domain层  
- **Infrastructure层** 实现 Domain层接口
- **Domain层** 不依赖任何其他层

## 与传统MVC的对比

| 方面 | 传统MVC | 分层架构 |
|------|---------|----------|
| **业务逻辑位置** | Controller中 | Domain层中 |
| **数据库依赖** | Model直接依赖 | 通过接口抽象 |
| **测试难度** | 需要完整环境 | 可分层测试 |
| **代码复用** | 较难复用 | 高复用性 |
| **技术绑定** | 强绑定框架 | 技术无关 |

## 实践建议

### 1. 从Domain层开始设计
- 先定义业务实体和规则
- 再考虑如何实现和暴露

### 2. 接口优先
- 先定义接口，再实现
- 便于测试和替换实现

### 3. 保持层次边界
- 不要跨层调用
- 通过依赖注入管理依赖

### 4. 渐进式重构
- 可以从简单的三层架构开始
- 随着复杂度增加再细化分层

## 总结

这种架构设计的核心思想是**"业务逻辑与技术实现分离"**，让代码更加：
- 🎯 **聚焦** - 每层职责明确
- 🔧 **灵活** - 技术栈可替换  
- 🧪 **可测** - 分层测试策略
- 👥 **协作** - 团队并行开发
- 📈 **扩展** - 适应业务增长

通过这样的分层设计，我们能够构建出既满足当前需求，又能够适应未来变化的高质量软件系统。 