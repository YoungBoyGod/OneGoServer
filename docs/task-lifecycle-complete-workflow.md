# 任务完整生命周期流程指南

## 概述

本文档详细描述OneGo系统中任务从创建到完成的完整生命周期流程，基于v2.0优化架构，涵盖任务创建、队列分配、设备执行、优先级调整、重发机制等全流程。

## 🔄 完整流程概览

### 核心流程阶段

```
任务创建 → 全局队列 → 设备分配 → 设备队列 → 执行处理 → 结果处理 → 历史记录
```

### 涉及的核心表

| 表名 | 作用 | 关键字段 |
|------|------|----------|
| `tasks` | 任务主表 | task_id, name, type, priority |
| `task_assignment_queue` | 全局分配队列 | task_id, assigned_device_id, queue_status |
| `device_task_queue` | 设备执行队列 | device_id, task_id, queue_priority, queue_position |
| `device_load_monitor` | 设备负载监控 | device_id, load_score, current_tasks |
| `task_executions` | 执行记录 | task_id, status, start_time, end_time |

## 📋 详细流程步骤

### 阶段1: 任务创建和初始化

#### 1.1 任务创建
```sql
-- 创建新任务
INSERT INTO tasks (
    task_id, name, description, type, status, priority,
    parameters, executor_type, timeout, max_retries
) VALUES (
    'task_001', '系统检查任务', '检查设备系统状态', 'system_check', 
    'pending', 5, '{"check_type": "full"}', 'shell', 3600, 3
);
```

**关键点**：
- `task_id` 使用业务唯一标识
- `status` 初始为 'pending'
- `priority` 数字越小优先级越高（1-10）

#### 1.2 加入全局分配队列
```sql
-- 任务进入全局分配队列
INSERT INTO task_assignment_queue (
    task_id, priority, queue_status,
    required_device_type, required_capabilities
) VALUES (
    'task_001', 5, 'queued',
    'server', '{"min_cpu": 2, "min_memory": 4096}'
);
```

**触发机制**：
- 自动触发：任务创建时自动加入队列
- 手动触发：管理员手动分配特定任务

### 阶段2: 设备分配算法

#### 2.1 调度器扫描
```sql
-- 调度器查询待分配任务
SELECT taq.*, t.timeout, t.estimated_duration
FROM task_assignment_queue taq
JOIN tasks t ON taq.task_id = t.task_id
WHERE taq.queue_status = 'queued'
ORDER BY taq.priority ASC, taq.queued_at ASC
LIMIT 100;
```

#### 2.2 设备选择策略
```sql
-- 查找可用设备
SELECT d.id, d.device_id, dlm.load_score, dlm.current_tasks
FROM devices d
JOIN device_load_monitor dlm ON d.id = dlm.device_id
WHERE d.status = 'online'
  AND d.type = '${required_device_type}'
  AND dlm.current_tasks < dlm.max_concurrent_tasks
  AND dlm.load_score < 80.0
ORDER BY dlm.load_score ASC, dlm.current_tasks ASC
LIMIT 5;
```

**分配策略**：
1. **负载优先**: 选择load_score最低的设备
2. **任务数平衡**: 优先选择当前任务数少的设备
3. **能力匹配**: 确保设备满足任务的capability要求
4. **排除策略**: 跳过excluded_device_ids中的设备

#### 2.3 确认分配
```sql
-- 更新分配结果
UPDATE task_assignment_queue 
SET assigned_device_id = ${selected_device_id},
    queue_status = 'assigned',
    assigned_at = CURRENT_TIMESTAMP,
    assignment_score = ${calculated_score}
WHERE task_id = 'task_001';

-- 记录分配历史
INSERT INTO task_assignment_history (
    task_id, device_id, action, new_status, 
    details, operation_source
) VALUES (
    'task_001', ${selected_device_id}, 'assigned', 'assigned',
    '{"algorithm_score": 85.5, "selection_reason": "lowest_load"}',
    'system'
);
```

### 阶段3: 设备队列管理

#### 3.1 加入设备执行队列
```sql
-- 任务加入设备队列
INSERT INTO device_task_queue (
    device_id, task_id, queue_priority, status,
    estimated_start_time, estimated_duration,
    max_retry_count, timeout_seconds
) VALUES (
    ${device_id}, 'task_001', 5, 'queued',
    CURRENT_TIMESTAMP + INTERVAL '10 minutes', 1800,
    3, 3600
);
```

#### 3.2 触发器自动处理
```sql
-- fn_auto_position_and_manual_flags() 自动执行：
-- 1. 自动分配queue_position到队列末尾
-- 2. 设置original_priority = queue_priority
-- 3. 记录操作历史
```

**自动处理逻辑**：
- **位置分配**: 新任务自动排到该设备队列的末尾
- **优先级记录**: 保存original_priority用于重置
- **历史记录**: 自动记录'add'操作到operation_history

### 阶段4: 优先级调整机制

#### 4.1 全局优先级调整
```sql
-- 调整全局队列优先级
UPDATE task_assignment_queue 
SET priority = 1  -- 提升到最高优先级
WHERE task_id = 'task_001';
```

#### 4.2 设备队列优先级调整
```sql
-- 调整设备队列优先级
UPDATE device_task_queue 
SET queue_priority = 1
WHERE device_id = ${device_id} AND task_id = 'task_001';
```

**触发器处理**：
```sql
-- fn_auto_position_and_manual_flags() 自动标记：
NEW.is_manual_priority := TRUE;
NEW.last_priority_change_at := CURRENT_TIMESTAMP;
NEW.last_action := 'priority_changed';

-- fn_log_device_queue_op() 自动记录：
INSERT INTO device_queue_operation_history(...) VALUES (...);
```

#### 4.3 队列位置调整
```sql
-- 调整队列位置（插队）
UPDATE device_task_queue 
SET queue_position = 2  -- 插到第2位
WHERE device_id = ${device_id} AND task_id = 'task_001';
```

**自动位置调整**：
- 原第2位及以后的任务位置自动+1
- 标记`is_manual_position = true`
- 记录`last_position_change_at`时间戳

### 阶段5: 执行器处理

#### 5.1 任务拉取
```sql
-- 执行器拉取待执行任务
SELECT dtq.*, t.parameters, t.executor_type
FROM device_task_queue dtq
JOIN tasks t ON dtq.task_id = t.task_id
WHERE dtq.device_id = ${executor_device_id}
  AND dtq.status = 'queued'
ORDER BY dtq.queue_priority ASC, dtq.queue_position ASC
LIMIT 1;
```

#### 5.2 开始执行
```sql
-- 标记任务开始执行
UPDATE device_task_queue 
SET status = 'executing',
    actual_start_time = CURRENT_TIMESTAMP
WHERE device_id = ${device_id} AND task_id = 'task_001';

-- 更新设备负载
UPDATE device_load_monitor 
SET current_tasks = current_tasks + 1
WHERE device_id = ${device_id};
```

#### 5.3 执行监控
```sql
-- 检查超时任务
SELECT dtq.task_id, dtq.actual_start_time, dtq.timeout_seconds
FROM device_task_queue dtq
WHERE dtq.status = 'executing'
  AND dtq.actual_start_time + INTERVAL '1 second' * dtq.timeout_seconds < CURRENT_TIMESTAMP;
```

### 阶段6: 结果处理

#### 6.1 成功完成
```sql
-- 标记任务完成
UPDATE device_task_queue 
SET status = 'completed',
    actual_end_time = CURRENT_TIMESTAMP
WHERE device_id = ${device_id} AND task_id = 'task_001';

-- 创建执行记录
INSERT INTO task_executions (
    task_id, execution_id, device_esn, status,
    start_time, end_time, duration,
    output, metrics
) VALUES (
    'task_001', 'exec_001_001', '${device_esn}', 'completed',
    '${start_time}', CURRENT_TIMESTAMP, ${duration},
    '{"result": "success", "output": "系统检查完成"}',
    '{"cpu_usage": 15.5, "memory_usage": 512}'
);
```

#### 6.2 执行失败
```sql
-- 标记任务失败
UPDATE device_task_queue 
SET status = 'failed'
WHERE device_id = ${device_id} AND task_id = 'task_001';
```

#### 6.3 重发判断
```sql
-- 检查是否需要重发
SELECT dtq.current_retry, dtq.max_retry_count
FROM device_task_queue dtq
WHERE dtq.device_id = ${device_id} 
  AND dtq.task_id = 'task_001'
  AND dtq.status = 'failed';
```

### 阶段7: 重发机制

#### 7.1 自动重发
```sql
-- 重发失败任务
UPDATE device_task_queue 
SET status = 'queued'  -- failed -> queued 触发重发
WHERE device_id = ${device_id} AND task_id = 'task_001';
```

#### 7.2 触发器自动处理
```sql
-- fn_handle_requeue() 自动执行：
IF OLD.status IN ('canceled','failed') AND NEW.status = 'queued' THEN
    NEW.requeue_count   := OLD.requeue_count + 1;
    NEW.is_requeued     := TRUE;
    NEW.last_requeue_at := CURRENT_TIMESTAMP;
    NEW.current_retry   := 0;  -- 重置重试次数
    NEW.actual_start_time := NULL;  -- 清除执行时间
    NEW.actual_end_time := NULL;
END IF;
```

#### 7.3 手动重发
```sql
-- 使用内置函数重发
SELECT requeue_failed_task(${device_id}, 'task_001', '手动重试');
```

## 🔍 监控和查询

### 全局队列状态
```sql
-- 查看全局分配队列状态
SELECT 
    COUNT(*) as total_tasks,
    COUNT(*) FILTER (WHERE queue_status = 'queued') as pending_assignments,
    COUNT(*) FILTER (WHERE queue_status = 'assigned') as assigned_tasks,
    AVG(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - queued_at))/60) as avg_wait_minutes
FROM task_assignment_queue;
```

### 设备队列状态
```sql
-- 使用优化视图查看设备队列状态
SELECT * FROM device_queue_status 
WHERE device_id = ${device_id};

-- 详细队列信息
SELECT * FROM device_queue_details 
WHERE device_id = ${device_id}
ORDER BY queue_position ASC;
```

### 操作历史追踪
```sql
-- 查看任务操作历史
SELECT 
    operation_type, operation_time, 
    old_status, new_status, operation_source
FROM device_queue_operation_history 
WHERE task_id = 'task_001'
ORDER BY operation_time DESC;
```

## 🛠️ 批量管理操作

### 批量优先级调整
```sql
-- 批量调整多个任务的优先级
SELECT batch_update_priority(
    ${device_id}, 
    ARRAY['task_001', 'task_002', 'task_003'], 
    1,  -- 新优先级
    ${operator_id}
);
```

### 批量重发操作
```sql
-- 批量重发失败任务
WITH failed_tasks AS (
    SELECT device_id, task_id 
    FROM device_task_queue 
    WHERE status = 'failed' 
      AND current_retry < max_retry_count
)
UPDATE device_task_queue 
SET status = 'queued' 
WHERE (device_id, task_id) IN (SELECT device_id, task_id FROM failed_tasks);
```

### 批量取消操作
```sql
-- 批量取消指定设备的待执行任务
UPDATE device_task_queue 
SET status = 'canceled',
    cancel_reason = '系统维护，批量取消任务'
WHERE device_id = ${device_id} 
  AND status = 'queued';
```

## 📊 性能监控指标

### 队列性能指标
```sql
-- 队列长度监控
SELECT 
    device_id,
    COUNT(*) as total_queue_length,
    COUNT(*) FILTER (WHERE status = 'queued') as pending_count,
    COUNT(*) FILTER (WHERE is_requeued = true) as requeue_count,
    AVG(queue_priority) as avg_priority
FROM device_task_queue 
GROUP BY device_id;
```

### 重发频率监控
```sql
-- 重发频率统计
SELECT 
    DATE(last_requeue_at) as requeue_date,
    COUNT(*) as total_requeues,
    COUNT(DISTINCT task_id) as unique_tasks,
    AVG(requeue_count) as avg_requeue_count
FROM device_task_queue 
WHERE is_requeued = true
GROUP BY DATE(last_requeue_at)
ORDER BY requeue_date DESC;
```

### 执行效率统计
```sql
-- 任务执行效率分析
SELECT 
    device_id,
    COUNT(*) as total_executed,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_count,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_count,
    ROUND(
        COUNT(*) FILTER (WHERE status = 'completed') * 100.0 / COUNT(*), 2
    ) as success_rate,
    AVG(
        EXTRACT(EPOCH FROM (actual_end_time - actual_start_time))
    ) as avg_execution_seconds
FROM device_task_queue 
WHERE actual_start_time IS NOT NULL
GROUP BY device_id;
```

## 🚨 异常处理和告警

### 队列堆积告警
```sql
-- 检查队列堆积情况
SELECT 
    device_id,
    COUNT(*) as queue_length,
    MIN(created_at) as oldest_task_time
FROM device_task_queue 
WHERE status = 'queued'
GROUP BY device_id
HAVING COUNT(*) > 50  -- 队列长度超过50告警
   OR MIN(created_at) < CURRENT_TIMESTAMP - INTERVAL '1 hour';  -- 最老任务超过1小时
```

### 重发频率异常
```sql
-- 检查重发频率异常
SELECT 
    task_id, device_id, requeue_count,
    last_requeue_at, cancel_reason
FROM device_task_queue 
WHERE requeue_count > 5  -- 重发次数超过5次
   OR (is_requeued = true AND last_requeue_at > CURRENT_TIMESTAMP - INTERVAL '10 minutes');
```

### 设备负载异常
```sql
-- 检查设备负载异常
SELECT 
    device_id, load_score, current_tasks,
    max_concurrent_tasks, last_heartbeat
FROM device_load_monitor 
WHERE load_score > 90  -- 负载超过90%
   OR current_tasks > max_concurrent_tasks  -- 超过最大并发数
   OR last_heartbeat < CURRENT_TIMESTAMP - INTERVAL '5 minutes';  -- 5分钟无心跳
```

## 📋 最佳实践建议

### 1. 任务设计原则
- **幂等性**: 确保任务可以安全重复执行
- **超时设置**: 合理设置timeout_seconds避免任务卡死
- **参数验证**: 在tasks.parameters中包含完整的执行参数

### 2. 优先级管理
- **谨慎调整**: 避免频繁手动调整优先级
- **批量操作**: 使用批量函数提高操作效率
- **监控影响**: 观察优先级调整对系统整体性能的影响

### 3. 重发策略
- **合理重试**: 根据任务类型设置合适的max_retry_count
- **失败分析**: 定期分析高频重发任务的失败原因
- **自动清理**: 定期清理长期失败的任务记录

### 4. 性能优化
- **索引利用**: 充分利用已优化的索引进行查询
- **视图查询**: 使用预计算视图减少复杂JOIN
- **批量处理**: 避免单条记录操作，优先使用批量函数

## 🎯 总结

这个完整的任务生命周期流程基于OneGo系统v2.0优化架构，提供了：

### ✅ 完整流程覆盖
- 从任务创建到最终完成的全流程管理
- 智能设备分配和负载均衡
- 灵活的优先级调整机制
- 可靠的重发和异常处理

### ✅ 自动化处理
- 触发器自动管理队列位置和优先级
- 自动记录操作历史和审计信息
- 智能重发机制和故障恢复

### ✅ 监控和维护
- 丰富的监控视图和统计查询
- 完善的异常检测和告警机制
- 高效的批量管理操作

这个流程设计确保了任务处理的高效性、可靠性和可维护性，为OneGo系统提供了坚实的任务管理基础。 