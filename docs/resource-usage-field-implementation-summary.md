# 任务执行资源使用字段实现总结

## 问题背景

用户提出关键问题：**"资源使用一般是动态的，这里如何填写呢？"**

这个问题指出了在 `task_executions` 表中，资源使用字段（cpu_usage、memory_usage、io_operations）面临的核心挑战：
- 任务执行过程中资源使用是实时变化的
- 单个静态字段无法完整反映动态资源消耗情况
- 需要设计合理的数据结构来存储和展示资源使用信息

## 解决方案

### 采用聚合统计值方案

将原有的单值字段改为多个聚合统计字段：

**原设计 (问题)**:
```sql
cpu_usage       NUMERIC(5,2), -- cpu使用率
memory_usage    NUMERIC(10,2), -- 内存使用率
io_operations   BIGINT, -- 磁盘IO操作次数
```

**新设计 (解决)**:
```sql
-- 资源使用统计 (聚合数据)
cpu_usage_avg    NUMERIC(5,2), -- CPU平均使用率(%)
cpu_usage_peak   NUMERIC(5,2), -- CPU峰值使用率(%)
memory_usage_avg NUMERIC(10,2), -- 内存平均使用量(MB)
memory_usage_peak NUMERIC(10,2), -- 内存峰值使用量(MB)
io_operations_total BIGINT, -- IO操作总次数
io_bytes_total   BIGINT, -- IO字节总数
```

### 数据采集流程

```
任务执行器启动 
    ↓
定期采样 (5秒间隔)
    ↓
本地缓存统计数据
    ↓
任务结束时计算聚合值
    ↓
写入数据库
```

### 实际数据填写示例

**场景**: 一个备份任务执行60秒，每5秒采样一次

| 时间点 | CPU使用率(%) | 内存使用(MB) | IO操作次数 |
|--------|-------------|-------------|-----------|
| 0s     | 15.2        | 512         | 0         |
| 5s     | 45.8        | 1024        | 125       |
| 10s    | 67.3        | 1280        | 250       |
| 15s    | 78.9        | 1536        | 380       |
| 20s    | 65.4        | 1280        | 445       |
| ...    | ...         | ...         | ...       |
| 60s    | 23.1        | 768         | 1250      |

**计算聚合值**:
- `cpu_usage_avg`: 45.67% (所有采样点平均值)
- `cpu_usage_peak`: 78.9% (最高值)
- `memory_usage_avg`: 1024.5MB (平均值)
- `memory_usage_peak`: 1536MB (峰值)
- `io_operations_total`: 1250 (累计总数)
- `io_bytes_total`: 52428800 (累计字节数)

## 技术实现

### 1. 数据库迁移脚本修改

**文件**: `internal/data/migrations/001_create_tasks_table.sql`

修改了task_executions表的资源使用字段定义，从3个简单字段扩展为6个聚合统计字段。

### 2. GORM模型同步更新

**文件**: `internal/biz/task/models.go`

```go
// TaskExecution结构体字段更新
type TaskExecution struct {
    // ... 其他字段
    
    // 资源使用统计 (聚合数据)
    CPUUsageAvg       *float64 `gorm:"type:numeric(5,2)" json:"cpu_usage_avg,omitempty"`
    CPUUsagePeak      *float64 `gorm:"type:numeric(5,2)" json:"cpu_usage_peak,omitempty"`
    MemoryUsageAvg    *float64 `gorm:"type:numeric(10,2)" json:"memory_usage_avg,omitempty"`
    MemoryUsagePeak   *float64 `gorm:"type:numeric(10,2)" json:"memory_usage_peak,omitempty"`
    IOOperationsTotal *int64   `json:"io_operations_total,omitempty"`
    IOBytesTotal      *int64   `json:"io_bytes_total,omitempty"`
}
```

### 3. Metrics字段配合使用

除了基础的聚合统计字段，还可以利用 `metrics` JSONB字段存储更详细的资源使用轨迹：

```json
{
  "sampling_interval": 5,
  "total_samples": 12,
  "resource_summary": {
    "cpu": {"avg": 45.67, "peak": 78.9, "min": 15.2},
    "memory": {"avg": 1024.5, "peak": 1536, "min": 512},
    "io": {"total_ops": 1250, "total_bytes": 52428800}
  },
  "critical_events": [
    {"timestamp": "2024-01-15T10:02:15Z", "event": "cpu_spike", "value": 78.9},
    {"timestamp": "2024-01-15T10:03:45Z", "event": "memory_peak", "value": 1536}
  ]
}
```

## 优势分析

### ✅ 解决的问题
1. **动态资源监控**: 通过聚合统计有效反映资源使用变化
2. **性能分析**: 提供平均值和峰值，便于性能分析
3. **存储效率**: 平衡了数据完整性和存储空间
4. **查询效率**: 避免了复杂的时序数据查询

### ✅ 业务价值
1. **任务优化**: 通过资源使用统计优化任务执行策略
2. **容量规划**: 基于历史数据进行资源容量规划
3. **异常监控**: 通过峰值监控发现异常任务
4. **成本分析**: 准确计算任务执行的资源成本

## 后续扩展

### 阶段规划
1. **当前阶段**: 基础聚合统计 ✅
2. **下一阶段**: 完善metrics字段数据结构
3. **未来阶段**: 集成Prometheus等监控系统

### 可选增强
- 创建独立的资源监控表存储详细时序数据
- 集成图表组件展示资源使用趋势
- 添加资源使用告警机制

## Git提交记录

1. `197c218` - 优化任务执行资源使用字段设计：从单值改为聚合统计
2. `73737f9` - 同步GORM模型：更新TaskExecution资源使用字段为聚合统计

## 结论

通过将动态的资源使用数据转换为聚合统计值，我们有效解决了用户提出的"动态资源如何填写"的问题。这种设计既保留了关键的性能指标，又保证了数据库的查询效率和存储优化。 