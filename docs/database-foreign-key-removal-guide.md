# 数据库外键约束移除说明文档

## 概述

根据公司要求，禁止在数据库层面使用外键约束（FOREIGN KEY），本次变更移除了所有migration文件中的外键约束，改为在应用层维护数据一致性。

## 变更原因

### 公司政策要求
- 公司明确禁止使用数据库外键约束
- 需要提高数据库性能和可维护性
- 避免外键约束导致的锁竞争问题
- 提高数据库的灵活性和扩展性

### 技术考虑
- **性能优化**: 外键约束会在INSERT/UPDATE/DELETE操作时增加额外检查开销
- **锁竞争**: 外键约束可能导致表级锁，影响并发性能
- **维护复杂性**: 外键约束增加了数据迁移和表结构变更的复杂度
- **分库分表**: 外键约束不利于后续的数据库水平拆分

## 变更范围

### 涉及的Migration文件

#### 1. 001_create_tasks_table.sql
- **移除约束**: `task_executions`表的`task_id`外键约束
- **影响表**: task_executions → tasks

#### 2. 002_create_devices_table.sql
- **移除约束**: 
  - `device_heartbeats`表的`device_id`外键约束
  - `device_logs`表的`device_id`外键约束
  - `device_commands`表的`device_id`外键约束
  - `tasks`表的`device_id`外键约束
- **影响表**: 
  - device_heartbeats → devices
  - device_logs → devices
  - device_commands → devices
  - tasks → devices

#### 3. 003_create_task_assignment_queue.sql
- **移除约束**:
  - `task_assignment_queue`表的`task_id`和`assigned_device_id`外键约束
  - `device_load_monitor`表的`device_id`外键约束
  - `task_assignment_history`表的`task_id`和`device_id`外键约束
- **影响表**:
  - task_assignment_queue → tasks, devices
  - device_load_monitor → devices
  - task_assignment_history → tasks, devices

#### 4. 004_create_device_task_queue.sql
- **移除约束**:
  - `device_task_queue`表的多个外键约束
  - `device_queue_operation_history`表的多个外键约束
  - `device_queue_config`表的`device_id`外键约束
- **影响表**:
  - device_task_queue → devices, tasks, users
  - device_queue_operation_history → devices, tasks, users
  - device_queue_config → devices

### 其他修正
- 修正了`devices`表索引中引用不存在字段`last_seen`的问题
- 将`last_seen`相关索引改为`last_online_time`

## 应对措施

### 1. 应用层数据一致性保证

#### 事务管理
```go
// 示例：删除设备时的数据一致性处理
func DeleteDevice(ctx context.Context, deviceID int64) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 1. 删除相关心跳记录
    if err := deleteDeviceHeartbeats(tx, deviceID); err != nil {
        return err
    }
    
    // 2. 删除相关日志记录
    if err := deleteDeviceLogs(tx, deviceID); err != nil {
        return err
    }
    
    // 3. 删除相关命令记录
    if err := deleteDeviceCommands(tx, deviceID); err != nil {
        return err
    }
    
    // 4. 更新任务表设备ID为NULL
    if err := updateTasksDeviceID(tx, deviceID, nil); err != nil {
        return err
    }
    
    // 5. 删除设备记录
    if err := deleteDevice(tx, deviceID); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

#### 数据验证
```go
// 示例：创建任务时验证设备存在性
func CreateTask(ctx context.Context, task *Task) error {
    if task.DeviceID != nil {
        // 验证设备是否存在
        exists, err := deviceExists(ctx, *task.DeviceID)
        if err != nil {
            return err
        }
        if !exists {
            return errors.New("设备不存在")
        }
    }
    
    return createTask(ctx, task)
}
```

### 2. 数据库约束替代方案

#### 保留的约束
- **唯一约束**: 继续使用UNIQUE约束保证数据唯一性
- **检查约束**: 使用CHECK约束验证数据格式和范围
- **非空约束**: 继续使用NOT NULL约束

#### 索引策略
- 保留所有性能相关的索引
- 在外键字段上保留索引以提高查询性能
- 定期监控查询性能，优化索引策略

### 3. 数据一致性监控

#### 数据校验任务
```sql
-- 定期检查数据一致性的SQL示例

-- 检查孤立的心跳记录
SELECT COUNT(*) as orphaned_heartbeats
FROM device_heartbeats dh
LEFT JOIN devices d ON dh.device_id = d.device_id
WHERE d.device_id IS NULL;

-- 检查孤立的任务记录
SELECT COUNT(*) as orphaned_tasks
FROM tasks t
LEFT JOIN devices d ON t.device_id = d.id
WHERE t.device_id IS NOT NULL AND d.id IS NULL;
```

#### 监控告警
- 设置定时任务检查数据一致性
- 发现数据不一致时及时告警
- 建立数据修复机制

### 4. 开发规范

#### 代码审查要点
- 确保删除操作包含级联删除逻辑
- 验证关联数据的存在性检查
- 使用事务保证操作的原子性

#### 测试要求
- 单元测试覆盖数据一致性逻辑
- 集成测试验证级联操作
- 性能测试确保无性能退化

## 迁移指南

### 现有数据库迁移
```sql
-- 如果现有数据库已存在外键约束，需要执行以下SQL移除

-- 移除tasks表的外键约束
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS fk_tasks_device_id;

-- 移除device_heartbeats表的外键约束
ALTER TABLE device_heartbeats DROP CONSTRAINT IF EXISTS device_heartbeats_device_id_fkey;

-- 移除device_logs表的外键约束
ALTER TABLE device_logs DROP CONSTRAINT IF EXISTS device_logs_device_id_fkey;

-- 移除device_commands表的外键约束
ALTER TABLE device_commands DROP CONSTRAINT IF EXISTS device_commands_device_id_fkey;

-- 移除其他表的外键约束...
```

### 应用代码修改
1. **Repository层**: 添加数据一致性检查逻辑
2. **Service层**: 实现事务级联操作
3. **Controller层**: 增强参数验证
4. **测试代码**: 补充数据一致性测试用例

## 风险评估

### 潜在风险
1. **数据不一致**: 应用层逻辑错误可能导致数据不一致
2. **性能影响**: 应用层验证可能增加响应时间
3. **开发复杂度**: 需要开发者手动维护数据关联关系

### 风险缓解
1. **完善测试**: 全面的单元测试和集成测试
2. **代码审查**: 严格的代码审查流程
3. **监控告警**: 实时数据一致性监控
4. **文档完善**: 详细的开发规范文档

## 性能影响分析

### 预期性能提升
- **INSERT性能**: 减少外键约束检查开销，预计提升10-20%
- **DELETE性能**: 减少级联检查开销，预计提升15-25%
- **并发性能**: 减少锁竞争，提升并发处理能力

### 性能监控指标
- 数据库操作响应时间
- 事务处理吞吐量
- 锁等待时间统计
- 死锁发生频率

## 后续工作计划

### 短期任务（1-2周）
- [ ] 完成所有相关代码的数据一致性逻辑实现
- [ ] 编写和执行数据一致性测试用例
- [ ] 部署到测试环境验证

### 中期任务（1个月）
- [ ] 建立数据一致性监控系统
- [ ] 完善开发规范文档
- [ ] 培训开发团队新的开发模式

### 长期任务（3个月）
- [ ] 性能优化和调优
- [ ] 建立自动化数据修复机制
- [ ] 总结最佳实践经验

## 总结

移除数据库外键约束是一个重要的架构调整，需要团队协作确保数据一致性。通过合理的应用层设计、完善的测试和监控，可以在提升性能的同时保证数据质量。

## 相关文档
- [设备生命周期字段交互分析](./device-lifecycle-field-analysis.md)
- [设备生命周期API交互详细指南](./device-lifecycle-api-interaction-guide.md)
- [数据库设计分析](./database-design-analysis.md) 