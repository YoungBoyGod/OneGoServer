# 任务执行资源使用字段设计解决方案

## 问题分析

在 `task_executions` 表中，资源使用字段面临的挑战：

```sql
-- 当前设计
cpu_usage       NUMERIC(5,2), -- cpu使用率
memory_usage    NUMERIC(10,2), -- 内存使用率  
io_operations   BIGINT, -- 磁盘IO操作次数
```

**核心问题**: 资源使用是动态变化的，任务执行过程中会实时变化，单个静态字段无法完整反映资源消耗情况。

## 解决方案

### 方案一：聚合统计值 (推荐)

将资源使用字段改为存储关键统计指标：

```sql
-- 修改后的设计
cpu_usage_avg    NUMERIC(5,2), -- CPU平均使用率(%)
cpu_usage_peak   NUMERIC(5,2), -- CPU峰值使用率(%)
memory_usage_avg NUMERIC(10,2), -- 内存平均使用量(MB)
memory_usage_peak NUMERIC(10,2), -- 内存峰值使用量(MB)
io_operations_total BIGINT, -- IO操作总次数
io_bytes_total   BIGINT, -- IO字节总数
```

**优点**: 
- 提供关键性能指标
- 数据库查询高效
- 存储空间合理

**数据来源**:
- 任务执行器定期采样(每秒或每5秒)
- 计算平均值和峰值
- 任务结束时写入数据库

### 方案二：详细轨迹存储

使用 `metrics` JSONB字段存储完整资源使用轨迹：

```sql
-- 在metrics字段中存储
metrics JSONB -- 示例数据:
{
  "resource_timeline": [
    {"timestamp": "2024-01-15T10:00:01Z", "cpu": 25.5, "memory": 1024, "io_ops": 10},
    {"timestamp": "2024-01-15T10:00:06Z", "cpu": 67.2, "memory": 1280, "io_ops": 25},
    {"timestamp": "2024-01-15T10:00:11Z", "cpu": 45.8, "memory": 1150, "io_ops": 18}
  ],
  "resource_summary": {
    "cpu_avg": 46.17, "cpu_peak": 67.2,
    "memory_avg": 1151.33, "memory_peak": 1280,
    "io_total": 53
  }
}
```

**优点**:
- 保留完整数据轨迹
- 支持详细分析和图表展示
- 灵活的查询能力

**缺点**:
- 存储空间较大
- 查询性能可能受影响

### 方案三：独立资源监控表

创建专门的资源监控表：

```sql
CREATE TABLE task_resource_metrics (
    id              BIGSERIAL PRIMARY KEY,
    execution_id    VARCHAR(100) NOT NULL,
    timestamp       TIMESTAMP NOT NULL,
    cpu_usage       NUMERIC(5,2),
    memory_usage    NUMERIC(10,2),
    io_read_ops     BIGINT,
    io_write_ops    BIGINT,
    io_read_bytes   BIGINT,
    io_write_bytes  BIGINT,
    network_in      BIGINT,
    network_out     BIGINT,
    
    FOREIGN KEY (execution_id) REFERENCES task_executions(execution_id)
);
```

**优点**:
- 时序数据完整存储
- 支持复杂查询和分析
- 可接入监控系统

**缺点**:
- 数据量大，需要数据清理策略
- 查询复杂度增加

## 推荐实现

### 最佳实践组合方案

1. **task_executions表**: 存储聚合统计值(方案一)
2. **metrics字段**: 存储摘要轨迹数据(方案二的简化版)
3. **可选**: 对于重要任务，启用详细监控表(方案三)

### 数据采集流程

```
任务执行器 → 定期采样(5秒间隔) → 本地缓存 → 任务结束时 → 计算统计值 → 写入数据库
     ↓
  实时监控(可选) → 写入时序数据库 → 图表展示
```

### 示例数据结构

```json
// metrics字段示例
{
  "sampling_interval": 5,
  "total_samples": 12,
  "resource_summary": {
    "cpu": {"avg": 45.2, "peak": 78.5, "min": 12.3},
    "memory": {"avg": 1024.5, "peak": 1536.2, "min": 512.1},
    "io": {"total_ops": 1250, "total_bytes": 52428800}
  },
  "critical_events": [
    {"timestamp": "2024-01-15T10:02:15Z", "event": "cpu_spike", "value": 78.5},
    {"timestamp": "2024-01-15T10:03:45Z", "event": "memory_peak", "value": 1536.2}
  ]
}
```

## 实施建议

1. **阶段一**: 实施聚合统计值方案(立即可用)
2. **阶段二**: 完善metrics字段数据结构
3. **阶段三**: 根据需要增加详细监控表

这样既保证了基础的资源使用监控，又为将来的扩展留出了空间。 