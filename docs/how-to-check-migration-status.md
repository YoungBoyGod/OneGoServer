# 如何查看数据库迁移状态

## 方法一：使用迁移状态查看工具

我们提供了一个专门的工具来查看已执行的迁移状态：

### 运行工具
```bash
cd tools/migration-status
go run main.go
```

### 工具输出示例
```
=== 数据库迁移状态查看工具 ===

📊 已执行的迁移列表 (共 3 个):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
版本号     迁移名称                        执行时间              状态      执行耗时(ms)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
001        create_tasks_table              2025-07-03 09:48:51   ✅ 成功   1523
002        create_devices_table            2025-07-03 09:48:52   ✅ 成功   856
003        create_task_assignment_queue    2025-07-03 09:48:52   ✅ 成功   1205
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📋 详细信息:

🔹 版本 001 (create_tasks_table):
   执行时间: 2025-07-03 09:48:51
   状态: ✅ 成功执行
   执行耗时: 1523 ms
   校验和: a1b2c3d4...

🔹 版本 002 (create_devices_table):
   执行时间: 2025-07-03 09:48:52
   状态: ✅ 成功执行
   执行耗时: 856 ms
   校验和: e5f6g7h8...

✨ 查看完成!
```

## 方法二：直接查询数据库

如果您有数据库访问权限，可以直接查询 `schema_migrations` 表：

### 查看所有迁移记录
```sql
SELECT 
    version,
    name,
    applied_at,
    is_success,
    execution_time,
    checksum
FROM schema_migrations 
ORDER BY applied_at;
```

### 查看特定迁移的详细信息
```sql
SELECT * FROM schema_migrations WHERE version = '001';
```

### 查看失败的迁移
```sql
SELECT 
    version,
    name,
    error_message,
    applied_at
FROM schema_migrations 
WHERE is_success = false
ORDER BY applied_at DESC;
```

## 方法三：在代码中查看

您也可以在Go代码中使用我们提供的API：

```go
package main

import (
    "github.com/YoungBoyGod/OneGoServer/pkg/utils"
    "github.com/YoungBoyGod/OneGoServer/pkg/sql"
)

func checkMigrations() {
    gormDB := sql.GetDB()
    db, err := gormDB.DB()
    if err != nil {
        log.Fatal(err)
    }
    
    migrations, err := utils.GetMigrationStatus(db)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, migration := range migrations {
        fmt.Printf("版本 %s: %s (执行时间: %s)\n", 
            migration.Version, 
            migration.Name, 
            migration.AppliedAt.Format("2006-01-02 15:04:05"))
    }
}
```

## 数据库表结构

`schema_migrations` 表包含以下字段：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| version | VARCHAR(255) | 迁移版本号 |
| name | VARCHAR(255) | 迁移名称 |
| applied_at | TIMESTAMP | 执行时间 |
| applied_by | VARCHAR(255) | 执行者 |
| checksum | VARCHAR(64) | 文件校验和 |
| execution_time | INTEGER | 执行耗时(毫秒) |
| is_success | BOOLEAN | 是否成功 |
| error_message | TEXT | 错误信息 |
| rollback_sql | TEXT | 回滚SQL |

## 常见问题

### Q: 如果没有显示任何迁移记录怎么办？
A: 这可能意味着：
- 还没有执行过任何迁移
- `schema_migrations` 表还未创建
- 数据库连接配置有问题

### Q: 如何查看迁移的SQL内容？
A: SQL文件保存在 `internal/data/migrations/` 目录下，可以直接查看文件内容。

### Q: 如何重新执行某个迁移？
A: 可以使用回滚功能先回滚，然后重新执行：
```go
utils.RollbackMigration(db, "001")
utils.RunMigrations(db)
``` 