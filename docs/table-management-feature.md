# PostgreSQL 表管理功能详解

## 📋 概述

本文档详细介绍了新增到 `pkg/sql/postgresql.go` 模块中的**表管理功能**，包括动态表创建、删除、查询和测试表管理等企业级功能。

## 🚀 新增功能列表

### 1. 核心数据结构

#### 📊 TableColumn - 表列定义
```go
type TableColumn struct {
    Name         string `json:"name"`           // 列名
    Type         string `json:"type"`           // 数据类型 (VARCHAR, INT, TIMESTAMP等)
    NotNull      bool   `json:"not_null"`       // 是否非空
    PrimaryKey   bool   `json:"primary_key"`    // 是否主键
    Unique       bool   `json:"unique"`         // 是否唯一
    DefaultValue string `json:"default_value"`  // 默认值
}
```

#### 📋 TableDefinition - 表定义结构
```go
type TableDefinition struct {
    Name    string        `json:"name"`     // 表名
    Columns []TableColumn `json:"columns"`  // 列定义数组
}
```

### 2. 表管理API函数

#### 🔨 CreateTable - 创建表
```go
func CreateTable(tableDef *TableDefinition) error
```
**功能**: 根据表定义创建数据库表
**参数**: 
- `tableDef`: 表定义结构指针
**返回**: 错误信息(如果有)
**特性**:
- 自动构建 CREATE TABLE SQL 语句
- 支持主键、唯一索引、非空约束、默认值
- 包含详细的操作日志记录
- 执行时间监控

#### 🗑️ DropTable - 删除表
```go
func DropTable(tableName string) error
```
**功能**: 删除指定名称的数据库表
**参数**:
- `tableName`: 要删除的表名
**返回**: 错误信息(如果有)
**特性**:
- 使用 IF EXISTS 安全删除
- 详细的操作日志记录
- 执行时间监控

#### ✅ TableExists - 检查表是否存在
```go
func TableExists(tableName string) (bool, error)
```
**功能**: 检查指定表是否存在于数据库中
**参数**:
- `tableName`: 要检查的表名
**返回**: 
- `bool`: 表是否存在
- `error`: 错误信息(如果有)
**特性**:
- 查询 information_schema.tables
- 精确的存在性检查

#### 📋 GetTableInfo - 获取表信息
```go
func GetTableInfo(tableName string) (*TableDefinition, error)
```
**功能**: 获取指定表的完整结构信息
**参数**:
- `tableName`: 表名
**返回**:
- `*TableDefinition`: 表定义结构
- `error`: 错误信息(如果有)
**特性**:
- 从 information_schema.columns 获取列信息
- 包含列类型、约束、默认值等完整信息
- 返回结构化的表定义

#### 🧪 CreateTestTable - 创建测试表
```go
func CreateTestTable() (string, error)
```
**功能**: 创建带时间戳的测试表
**返回**:
- `string`: 生成的测试表名
- `error`: 错误信息(如果有)
**特性**:
- 自动生成唯一表名 `test_<timestamp>`
- 预定义的测试表结构
- 包含常用字段 (id, name, email, created_at, updated_at)

### 3. 内部辅助函数

#### 🔧 buildCreateTableSQL - SQL构建器
```go
func buildCreateTableSQL(tableDef *TableDefinition) string
```
**功能**: 根据表定义构建 CREATE TABLE SQL 语句
**特性**:
- 智能约束处理 (PRIMARY KEY, NOT NULL, UNIQUE, DEFAULT)
- 规范的SQL格式输出
- 支持复杂的列定义组合

## 🧪 测试功能详解

### TestAddAndDeleteTable 测试函数

#### 📝 测试覆盖范围
1. **创建和删除表操作**
   - 表名生成验证
   - SQL语句格式验证
   - 操作逻辑测试

2. **表名唯一性测试** 
   - 生成10个不同的时间戳表名
   - 验证表名的唯一性
   - 确保无重复表名

3. **表名格式验证**
   - PostgreSQL表名长度限制 (≤63字符)
   - 表名格式正则验证 `^test_\d+$`
   - 时间戳有效性验证

4. **表生命周期模拟**
   - 完整的表操作流程模拟
   - CREATE → INSERT → SELECT → UPDATE → DELETE → DROP
   - 操作成功性验证

5. **表管理功能测试**
   - TableDefinition 结构验证
   - TableColumn 字段验证
   - 主键和唯一约束验证

6. **SQL构建测试**
   - 验证SQL构建逻辑的组件
   - 列定义组合测试
   - 预期SQL格式验证

7. **表管理功能单元测试**
   - 时间戳生成逻辑验证
   - TableColumn 结构完整性测试
   - 数据类型和约束验证

#### 🏆 测试结果
```
✅ TestAddAndDeleteTable PASSED
   ├── ✅ 创建和删除表操作 (0.00s)
   ├── ✅ 表名唯一性测试 (0.00s)  
   ├── ✅ 表名格式验证 (0.00s)
   ├── ✅ 模拟表生命周期管理 (0.00s)
   ├── ✅ 表管理功能测试 (0.00s)
   ├── ✅ SQL构建测试 (0.00s)
   └── ✅ 表管理功能单元测试 (0.00s)
```

## 💻 使用示例

### 基础表创建示例

```go
// 定义表结构
tableDef := &sql.TableDefinition{
    Name: "users",
    Columns: []sql.TableColumn{
        {
            Name:       "id",
            Type:       "SERIAL",
            PrimaryKey: true,
            NotNull:    true,
        },
        {
            Name:    "username",
            Type:    "VARCHAR(50)",
            NotNull: true,
            Unique:  true,
        },
        {
            Name: "email",
            Type: "VARCHAR(255)",
            Unique: true,
        },
        {
            Name:         "created_at",
            Type:         "TIMESTAMP",
            DefaultValue: "CURRENT_TIMESTAMP",
        },
    },
}

// 创建表
if err := sql.CreateTable(tableDef); err != nil {
    log.Printf("Failed to create table: %v", err)
}
```

### 测试表管理示例

```go
// 创建测试表
tableName, err := sql.CreateTestTable()
if err != nil {
    log.Printf("Failed to create test table: %v", err)
    return
}
fmt.Printf("Created test table: %s\n", tableName)

// 检查表是否存在
exists, err := sql.TableExists(tableName)
if err != nil {
    log.Printf("Failed to check table existence: %v", err)
    return
}
fmt.Printf("Table exists: %t\n", exists)

// 获取表信息
tableInfo, err := sql.GetTableInfo(tableName)
if err != nil {
    log.Printf("Failed to get table info: %v", err)
    return
}
fmt.Printf("Table info: %+v\n", tableInfo)

// 删除表
if err := sql.DropTable(tableName); err != nil {
    log.Printf("Failed to drop table: %v", err)
    return
}
fmt.Printf("Table %s dropped successfully\n", tableName)
```

### 高级表定义示例

```go
// 复杂表结构定义
ordersDef := &sql.TableDefinition{
    Name: "orders",
    Columns: []sql.TableColumn{
        {
            Name:       "order_id",
            Type:       "BIGSERIAL",
            PrimaryKey: true,
        },
        {
            Name:    "user_id",
            Type:    "BIGINT",
            NotNull: true,
        },
        {
            Name:    "product_id",
            Type:    "BIGINT", 
            NotNull: true,
        },
        {
            Name:         "order_status",
            Type:         "VARCHAR(20)",
            NotNull:      true,
            DefaultValue: "'pending'",
        },
        {
            Name: "order_amount",
            Type: "DECIMAL(10,2)",
            NotNull: true,
        },
        {
            Name:         "created_at",
            Type:         "TIMESTAMP",
            NotNull:      true,
            DefaultValue: "CURRENT_TIMESTAMP",
        },
        {
            Name:         "updated_at",
            Type:         "TIMESTAMP", 
            DefaultValue: "CURRENT_TIMESTAMP",
        },
    },
}

// 创建订单表
if err := sql.CreateTable(ordersDef); err != nil {
    log.Fatal("Failed to create orders table:", err)
}
```

## 🔄 完整的表管理流程

### 1. 表创建流程
```
定义表结构 → 验证数据库连接 → 构建SQL → 执行创建 → 记录日志 → 返回结果
```

### 2. 表删除流程  
```
指定表名 → 验证数据库连接 → 构建DROP SQL → 执行删除 → 记录日志 → 返回结果
```

### 3. 表查询流程
```
表名输入 → 查询information_schema → 解析结果 → 构建结构 → 返回信息
```

### 4. 测试表流程
```
生成时间戳 → 构建表名 → 定义测试结构 → 创建表 → 返回表名
```

## 📊 SQL生成示例

### 基础CREATE TABLE
```sql
CREATE TABLE IF NOT EXISTS test_1751269264542411000 (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

### 复杂CREATE TABLE
```sql
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    order_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    order_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

### DROP TABLE
```sql
DROP TABLE IF EXISTS test_1751269264542411000
```

### 表存在性检查
```sql
SELECT COUNT(*) FROM information_schema.tables 
WHERE table_schema = 'public' AND table_name = ?
```

### 表信息查询
```sql
SELECT 
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_schema = 'public' AND table_name = ?
ORDER BY ordinal_position
```

## 🛡️ 安全特性

### 1. 输入验证
- 表名格式验证
- 列名和类型检查
- SQL注入防护

### 2. 错误处理
- 详细的错误信息包装
- 操作失败时的回滚机制
- 异常情况的优雅处理

### 3. 数据库保护
- 使用 IF EXISTS 安全删除
- 事务性操作支持
- 连接状态验证

## 📈 监控和日志

### 操作日志记录
```go
// 表创建日志
pkglog.LogDBOperation("CREATE_TABLE", tableName, duration, err)

// 表删除日志  
pkglog.LogDBOperation("DROP_TABLE", tableName, duration, err)

// 表检查日志
pkglog.LogDBOperation("CHECK_TABLE", tableName, duration, err)

// 成功日志
pkglog.LogInfo("Table created successfully", zap.String("table", tableName))
```

### 性能监控
- 每个操作的执行时间记录
- 成功/失败状态跟踪
- 详细的操作分类统计

## 🎯 最佳实践

### 1. 表命名规范
- 使用小写字母和下划线
- 避免保留关键字
- 表名长度控制在63字符以内

### 2. 列定义建议
- 合理设置 NOT NULL 约束
- 适当使用 UNIQUE 索引
- 为时间字段设置默认值

### 3. 测试表管理
- 使用时间戳确保表名唯一性
- 测试完成后及时清理测试表
- 避免在生产环境创建测试表

### 4. 错误处理
- 始终检查函数返回的错误
- 记录详细的操作日志
- 实现适当的重试机制

## 🔄 与现有系统集成

### 日志系统集成
- 使用项目统一的日志格式
- 结构化日志记录
- 支持不同日志级别

### 配置系统集成
- 继承数据库连接配置
- 支持环境变量配置
- 配置验证和默认值

### 事务系统集成
- 支持在事务中执行表操作
- 与 WithTransaction 函数兼容
- 操作失败时的自动回滚

## 📋 总结

表管理功能为 PostgreSQL 模块增加了强大的动态表管理能力：

### ✨ 核心价值
- **动态性**: 运行时创建和管理表结构
- **安全性**: 完善的输入验证和错误处理
- **可测试性**: 基于时间戳的测试表支持
- **可监控性**: 详细的操作日志和性能监控
- **易用性**: 简洁的API和丰富的使用示例

### 🎯 应用场景
- **动态业务表**: 根据业务需求动态创建数据表
- **测试环境**: 创建临时测试表进行功能验证
- **数据迁移**: 批量创建和管理表结构
- **多租户系统**: 为不同租户创建独立的数据表

### 🚀 未来扩展
- 支持索引管理
- 添加表约束管理
- 实现表结构版本控制
- 支持表数据备份和恢复

这个表管理功能显著增强了 PostgreSQL 模块的企业级能力，为复杂的数据库管理需求提供了完整的解决方案！ 