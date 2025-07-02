# 数据库迁移功能实现文档

## 功能概述

实现了一个完整的数据库迁移(migration)系统,用于管理数据库结构的版本控制。主要功能包括:

1. 自动执行未应用的migration
2. 记录migration执行历史
3. 支持回滚操作
4. 提供migration状态查询

## 核心组件

### 1. Migration表结构

```sql
CREATE TABLE schema_migrations (
    version         VARCHAR(255) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    applied_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    applied_by      VARCHAR(255) NOT NULL DEFAULT CURRENT_USER,
    checksum        VARCHAR(64) NOT NULL,
    execution_time  INTEGER NOT NULL,
    is_success      BOOLEAN NOT NULL DEFAULT TRUE,
    error_message   TEXT,
    rollback_sql    TEXT
);
```

### 2. 主要功能函数

- `RunMigrations(db *sql.DB) error`: 执行数据库迁移
- `RollbackMigration(db *sql.DB, version string) error`: 回滚指定版本
- `GetMigrationStatus(db *sql.DB) ([]Migration, error)`: 获取迁移状态

### 3. 工作流程

1. 创建migrations表(如果不存在)
2. 获取已执行的migrations
3. 读取migrations目录下的所有SQL文件
4. 按版本号排序
5. 执行未应用的migrations
6. 记录执行结果

## 使用说明

### 1. Migration文件命名规范

Migration文件必须遵循以下命名规范:
```
XXX_description.sql
```
其中:
- XXX: 三位数字版本号
- description: 迁移描述

例如: `001_create_users_table.sql`

### 2. Migration文件结构

每个migration文件应包含:
- 向前迁移的SQL语句
- 回滚SQL语句(可选)

### 3. 执行迁移

```go
db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname")
if err != nil {
    log.Fatal(err)
}

if err := utils.RunMigrations(db); err != nil {
    log.Fatal(err)
}
```

### 4. 回滚迁移

```go
if err := utils.RollbackMigration(db, "001"); err != nil {
    log.Fatal(err)
}
```

### 5. 查看迁移状态

```go
status, err := utils.GetMigrationStatus(db)
if err != nil {
    log.Fatal(err)
}
```

## 安全性考虑

1. 使用事务确保迁移的原子性
2. 记录SQL文件的MD5校验和防止文件被篡改
3. 记录执行时间和执行者信息便于审计
4. 详细的错误信息记录

## 最佳实践

1. 每个迁移文件应该是幂等的
2. 提供回滚SQL以支持版本回退
3. 定期清理旧的迁移记录
4. 在生产环境执行迁移前先在测试环境验证
5. 重要的数据库更改应该提供回滚计划 