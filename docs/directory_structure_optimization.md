# OneGoServer目录结构优化方案

## 概述

基于分布式任务系统的目录结构设计理念，结合OneGoServer的实际需求，制定渐进式的目录结构优化方案。

## 当前结构分析

### 现有结构
```
OneGoServer/
├── api/v1/                 # API定义
├── internal/               # 内部业务逻辑
│   ├── cmd/               # 服务器启动逻辑
│   ├── consts/            # 常量定义
│   ├── controller/        # 控制器
│   ├── config/            # 配置管理
│   └── service/           # 业务服务
├── pkg/                   # 公共组件
│   ├── middleware/        # 中间件
│   └── response/          # 响应处理
├── config/                # 配置文件
└── main.go               # 入口文件
```

### 存在问题
1. 缺乏数据访问层（DAO/Repository）
2. 模型定义不够清晰
3. 缺乏完整的工程化支持
4. 日志文件管理需要优化

## 优化方案

### 阶段一：基础架构重构（当前实施）

#### 目标结构
```
OneGoServer/
├── cmd/                          # 应用程序入口
│   └── server/
│       └── main.go              # 服务端主程序
├── api/                         # 对外接口定义
│   └── v1/                      # API版本管理
│       ├── auth.go              # 认证相关API
│       ├── user.go              # 用户相关API
│       ├── health.go            # 健康检查API
│       └── common.go            # 通用API定义
├── internal/                    # 内部业务逻辑
│   ├── app/                     # 应用核心
│   │   ├── server.go           # 服务器启动逻辑
│   │   └── routes.go           # 路由配置
│   ├── controller/              # 控制器层
│   │   ├── auth.go             # 认证控制器
│   │   ├── user.go             # 用户控制器
│   │   └── health.go           # 健康检查控制器
│   ├── service/                 # 业务服务层
│   │   ├── interfaces/         # 服务接口定义
│   │   ├── auth.go             # 认证服务
│   │   └── user.go             # 用户服务
│   ├── repository/              # 数据访问层
│   │   ├── interfaces/         # 数据接口定义
│   │   ├── mysql/              # MySQL实现
│   │   ├── sqlite/             # SQLite实现
│   │   └── redis/              # Redis实现
│   ├── model/                   # 数据模型
│   │   ├── entity/             # 数据实体
│   │   ├── dto/                # 数据传输对象
│   │   └── vo/                 # 值对象
│   ├── config/                  # 配置管理
│   │   ├── config.go           # 配置结构
│   │   └── loader.go           # 配置加载器
│   └── consts/                  # 常量定义
│       ├── error.go            # 错误常量
│       └── common.go           # 通用常量
├── pkg/                         # 可复用公共包
│   ├── middleware/              # 中间件
│   │   ├── auth.go             # 认证中间件
│   │   ├── cors.go             # CORS中间件
│   │   ├── logger.go           # 日志中间件
│   │   └── ratelimit.go        # 限流中间件
│   ├── response/                # 响应处理
│   │   └── response.go         # 统一响应格式
│   ├── database/                # 数据库工具
│   │   ├── mysql.go            # MySQL连接
│   │   ├── sqlite.go           # SQLite连接
│   │   └── redis.go            # Redis连接
│   ├── logger/                  # 日志工具
│   │   ├── zap.go              # Zap日志实现
│   │   └── interface.go        # 日志接口
│   └── utils/                   # 工具函数
│       ├── jwt.go              # JWT工具
│       ├── hash.go             # 哈希工具
│       └── validator.go        # 验证工具
├── configs/                     # 配置文件
│   ├── development.yaml         # 开发环境
│   ├── production.yaml          # 生产环境
│   └── test.yaml               # 测试环境
├── scripts/                     # 脚本文件
│   ├── build.sh                # 构建脚本
│   ├── migrate.sh              # 数据库迁移
│   └── deploy.sh               # 部署脚本
├── migrations/                  # 数据库迁移文件
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── logs/                        # 日志文件目录
│   ├── app_20250627155344.log   # 应用日志
│   ├── http_20250627155344.log  # HTTP访问日志
│   └── error_20250627155344.log # 错误日志
├── docs/                        # 文档
│   ├── api.md                  # API文档
│   ├── deployment.md           # 部署文档
│   └── architecture.md         # 架构文档
├── tests/                       # 测试文件
│   ├── integration/            # 集成测试
│   └── unit/                   # 单元测试
├── docker/                      # Docker文件
│   ├── Dockerfile              # Docker构建文件
│   └── docker-compose.yml      # 容器编排
├── go.mod                       # Go模块文件
├── go.sum                       # Go依赖锁定
├── Makefile                     # 构建脚本
└── README.md                    # 项目说明
```

#### 关键改进点

1. **清晰的分层架构**
   - Controller → Service → Repository → Model
   - 每层职责明确，依赖关系清晰

2. **接口导向设计**
   - service/interfaces/ 定义业务接口
   - repository/interfaces/ 定义数据接口
   - 便于测试和扩展

3. **数据库适配器模式**
   - 支持多种数据库实现
   - 便于切换和扩展

4. **完善的日志管理**
   - 应用日志、HTTP日志、错误日志分离
   - 时间戳命名，便于管理

5. **工程化支持**
   - 多环境配置
   - 数据库迁移
   - 测试支持
   - 部署脚本

### 阶段二：业务功能扩展

#### 用户管理模块
```
internal/
├── controller/
│   └── user.go                 # 用户CRUD操作
├── service/
│   ├── interfaces/
│   │   └── user.go            # 用户服务接口
│   └── user.go                # 用户业务逻辑
├── repository/
│   ├── interfaces/
│   │   └── user.go            # 用户数据接口
│   └── mysql/
│       └── user.go            # 用户数据操作
└── model/
    ├── entity/
    │   └── user.go            # 用户实体
    └── dto/
        ├── user_create.go     # 创建用户DTO
        └── user_update.go     # 更新用户DTO
```

#### 认证授权模块
```
internal/
├── controller/
│   └── auth.go                # 登录、注册、权限验证
├── service/
│   ├── interfaces/
│   │   └── auth.go           # 认证服务接口
│   └── auth.go               # 认证业务逻辑
└── model/
    └── dto/
        ├── login.go          # 登录请求DTO
        └── token.go          # Token响应DTO
```

### 阶段三：高级特性

#### 缓存层
```
pkg/
└── cache/
    ├── interface.go           # 缓存接口
    ├── redis.go              # Redis实现
    └── memory.go             # 内存实现
```

#### 消息队列
```
pkg/
└── queue/
    ├── interface.go          # 队列接口
    ├── redis.go              # Redis队列
    └── rabbitmq.go           # RabbitMQ实现
```

#### 监控指标
```
pkg/
└── metrics/
    ├── prometheus.go         # Prometheus指标
    └── collector.go          # 指标收集器
```

## 迁移策略

### 1. 渐进式迁移
- 不影响现有功能的前提下逐步重构
- 保持向后兼容性
- 分模块进行迁移

### 2. 迁移步骤
1. **移动配置文件**：config/server.yaml → configs/development.yaml
2. **重构启动逻辑**：main.go → cmd/server/main.go
3. **分离业务层**：创建service层和repository层
4. **完善模型定义**：统一model结构
5. **优化中间件**：重构pkg/middleware
6. **添加测试**：创建tests目录结构

### 3. 风险控制
- 保留原有结构作为备份
- 分阶段验证功能完整性
- 监控性能影响

## 与参考结构的对比

### 相似点
- 清晰的分层架构（cmd、internal、pkg）
- 接口导向设计
- 完整的工程化支持

### 差异点
- **简化程度**：去除了不必要的细分目录
- **实用导向**：专注于Web服务而非分布式任务系统
- **渐进式**：支持逐步演进而非一次性重构

### 适应性调整
- 保留了分布式系统的优秀设计模式
- 简化了复杂的目录层级
- 适合中小型Web服务项目

## 实施建议

### 1. 立即实施
- ✅ 多类型日志文件（已完成）
- 📝 配置文件重组
- 📝 启动逻辑重构

### 2. 短期实施（1-2周）
- 📝 创建repository层
- 📝 完善service层接口
- 📝 统一model定义

### 3. 中期实施（1个月）
- 📝 添加用户管理模块
- 📝 完善认证授权
- 📝 添加测试框架

### 4. 长期规划（3个月）
- 📝 缓存层集成
- 📝 监控指标系统
- 📝 部署自动化

## 总结

这个目录结构优化方案：
- ✅ **实用性强**：适合当前项目规模
- ✅ **可扩展性好**：支持未来功能扩展
- ✅ **风险可控**：渐进式重构策略
- ✅ **工程化完整**：覆盖开发、测试、部署全流程

建议从日志系统优化开始，逐步实施这个重构方案。 