# OneGoServer GoFrame架构重构方案

## 重构目标

参考[GoFrame工程目录设计](https://goframe.org/docs/design/project-structure)，将当前的DDD分层架构重构为GoFrame推荐的改进三层架构，在保持代码清晰度的同时提高开发效率。

## GoFrame架构核心思想

GoFrame采用的是**改进的三层架构**，主要特点：
- 📁 **简化分层**：相比DDD四层架构，减少为三层核心架构
- 🔄 **清晰流转**：请求按照 `api -> controller -> service -> dao` 流转
- 🎯 **职责明确**：每层都有明确的职责边界
- 🛠️ **工程实用**：更贴近实际项目开发需求

## 当前架构 vs GoFrame架构

### 当前DDD架构
```
internal/
├── domain/          # 领域层
├── application/     # 应用层  
├── infrastructure/  # 基础设施层
└── interfaces/      # 接口层
```

### GoFrame目标架构
```
OneGoServer/
├── api/                    # 对外接口定义
│   └── v1/
│       └── user.go        # 用户接口定义
├── internal/              # 内部逻辑
│   ├── cmd/               # 启动命令
│   ├── consts/            # 常量定义
│   ├── controller/        # 控制器层
│   │   └── user.go       # 用户控制器
│   ├── service/           # 业务逻辑层
│   │   └── user.go       # 用户业务逻辑
│   ├── dao/               # 数据访问层
│   │   └── user.go       # 用户数据访问
│   ├── model/             # 数据模型
│   │   ├── entity/        # 实体模型
│   │   │   └── user.go   # 用户实体
│   │   └── do/            # 数据操作对象
│   │       └── user.go   # 用户DO
│   └── config/            # 配置管理
├── manifest/              # 交付清单
│   ├── config/           # 配置文件
│   └── docker/           # Docker文件
├── resource/              # 静态资源
├── go.mod
└── main.go
```

## 核心层次说明

### 1. api层 - 对外接口定义
**职责**：定义HTTP API的请求和响应结构
```go
// api/v1/user.go
type RegisterReq struct {
    Username string `json:"username" v:"required|length:3,20"`
    Email    string `json:"email" v:"required|email"`
    Password string `json:"password" v:"required|length:6,20"`
}

type RegisterRes struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

### 2. controller层 - 接口处理
**职责**：接收请求，参数验证，调用service，返回响应
```go
// internal/controller/user.go
func (c *UserController) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterRes, error) {
    // 调用service层处理业务逻辑
    user, err := service.User().Register(ctx, req)
    if err != nil {
        return nil, err
    }
    
    return &v1.RegisterRes{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }, nil
}
```

### 3. service层 - 业务逻辑
**职责**：核心业务逻辑处理，数据校验，业务规则
```go
// internal/service/user.go
func (s *sUser) Register(ctx context.Context, req *v1.RegisterReq) (*entity.User, error) {
    // 业务逻辑：检查用户名是否存在
    if exists, err := dao.User.CheckUsernameExists(ctx, req.Username); err != nil {
        return nil, err
    } else if exists {
        return nil, errors.New("用户名已存在")
    }
    
    // 创建用户
    user := &entity.User{
        Username: req.Username,
        Email:    req.Email,
        Password: s.hashPassword(req.Password),
    }
    
    // 调用dao层保存
    if err := dao.User.Insert(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 4. dao层 - 数据访问
**职责**：数据库操作，只包含基础CRUD方法
```go
// internal/dao/user.go
func (d *UserDao) Insert(ctx context.Context, user *entity.User) error {
    _, err := d.DB.Model(&user).Insert()
    return err
}

func (d *UserDao) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
    count, err := d.DB.Model(&entity.User{}).Where("username", username).Count()
    return count > 0, err
}
```

### 5. model层 - 数据模型
**职责**：管理数据结构定义

#### entity - 实体模型
```go
// internal/model/entity/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### do - 数据操作对象
```go
// internal/model/do/user.go
type User struct {
    ID       interface{} // 支持各种查询条件
    Username interface{}
    Email    interface{}
    Password interface{}
}
```

## 请求流转过程

```
HTTP请求 → api验证 → controller处理 → service业务逻辑 → dao数据操作 → 返回响应
```

### 详细流程说明
1. **HTTP请求到达**：Gin路由接收请求
2. **api层验证**：绑定到对应的Req结构体，执行基础验证
3. **controller层处理**：接收验证后的请求，调用service
4. **service层业务逻辑**：执行具体的业务规则和逻辑
5. **dao层数据操作**：执行数据库CRUD操作
6. **响应返回**：将结果封装为Res结构体返回

## 与DDD架构的对比

| 方面 | DDD分层架构 | GoFrame架构 |
|------|-------------|-------------|
| **层次数量** | 4层 | 3层+api |
| **复杂度** | 较高 | 适中 |
| **学习成本** | 高 | 中等 |
| **开发效率** | 初期慢 | 较快 |
| **业务逻辑位置** | Domain层 | Service层 |
| **数据访问** | Repository模式 | DAO模式 |
| **适用场景** | 复杂业务 | 大部分业务场景 |

## 重构优势

### ✅ 保留的优点
- 分层清晰，职责明确
- 易于测试和维护
- 支持大团队协作
- 代码结构规范

### ✅ 新增的优势
- 降低学习成本
- 提高开发效率
- 减少代码量
- 更贴近实际项目需求

### ✅ 解决的问题
- DDD概念过于复杂
- 初期开发速度慢
- 团队理解成本高
- 过度设计问题

## 重构计划

### 阶段1：重构目录结构
1. 创建GoFrame标准目录
2. 迁移现有代码到新结构
3. 调整import路径

### 阶段2：简化分层逻辑
1. 将Domain Service逻辑合并到Service层
2. 将Repository接口简化为DAO
3. 保留Entity和DTO的概念

### 阶段3：优化代码结构
1. 统一错误处理
2. 完善配置管理
3. 添加中间件支持

### 阶段4：测试和文档
1. 更新测试用例
2. 完善API文档
3. 添加使用示例

## 实施建议

1. **渐进式迁移**：不需要一次性重构所有代码
2. **保持功能完整**：确保重构过程中功能不受影响
3. **团队培训**：让团队理解新的架构思想
4. **文档同步**：及时更新相关文档

通过这次重构，我们将获得一个既保持代码质量又提高开发效率的架构方案。 