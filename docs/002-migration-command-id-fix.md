# 002迁移文件command_id字段错误修复

## 错误描述

执行数据库迁移时出现错误：
```
Failed to run migrations: 执行migration 002 失败: ERROR: column "command_id" does not exist (SQLSTATE 42703)
```

## 问题分析

在 `002_create_devices_table.sql` 文件中，第164行定义了一个索引：
```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_device_commands_command_id ON device_commands(command_id);
```

但是在 `device_commands` 表的定义中，并没有 `command_id` 这个字段。表定义中只有：
- `id` (主键)
- `device_id`
- `command_type`
- `command_data`
- 其他字段...

## 修复方案

删除引用不存在字段的索引定义：

### 修复前
```sql
-- 创建设备命令索引
CREATE INDEX IF NOT EXISTS idx_device_commands_device_id ON device_commands(device_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_status ON device_commands(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_device_commands_command_id ON device_commands(command_id);
```

### 修复后
```sql
-- 创建设备命令索引
CREATE INDEX IF NOT EXISTS idx_device_commands_device_id ON device_commands(device_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_status ON device_commands(status);
```

## 修复内容

- 文件：`internal/data/migrations/002_create_devices_table.sql`
- 修改：删除第164行的错误索引定义
- 影响：`device_commands` 表能够正确创建

## 验证

修复后重新运行数据库迁移应该能够成功执行。 