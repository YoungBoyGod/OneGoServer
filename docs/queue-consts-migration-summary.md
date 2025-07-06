# 队列模块常量迁移总结

## 迁移概述

本次迁移将 `internal/model/queue/queue.go` 中定义的常量迁移到 `internal/consts/queue.go` 中，实现队列相关常量的统一管理。

## 迁移内容

### 1. 扩展 consts/queue.go

#### 新增常量类型
- **队列状态常量**: 4个状态常量（active, paused, stopped, error）
- **队列类型常量**: 4个类型常量（fifo, priority, delay, lifo）

#### 常量列表
```go
// 队列状态常量
QueueStatusActive  = "active"
QueueStatusPaused  = "paused"
QueueStatusStopped = "stopped"
QueueStatusError   = "error"

// 队列类型常量
QueueTypeFIFO     = "fifo"
QueueTypePriority = "priority"
QueueTypeDelay    = "delay"
QueueTypeLIFO     = "lifo"
```

#### 现有常量保持不变
- 设备队列任务状态常量（6个）
- 队列操作类型常量（8个）
- 调度策略常量（4个）
- 操作来源常量（4个）
- 队列优先级常量（6个）
- 队列配置常量（6个）
- 队列超时常量（3个）
- 队列重试常量（3个）
- 队列监控常量（3个）

### 2. 更新 model/queue/queue.go

#### 修改内容
- **移除重复定义**: 删除原有的队列状态和类型常量定义
- **使用consts引用**: 改为使用 `consts.XXX` 的方式引用常量
- **保持向后兼容**: 保留原有的常量名称，确保其他模块不受影响

#### 修改示例
```go
// 修改前
const (
    QueueStatusActive  = "active"
    QueueStatusPaused  = "paused"
    QueueStatusStopped = "stopped"
    QueueStatusError   = "error"
)

// 修改后
const (
    QueueStatusActive  = consts.QueueStatusActive
    QueueStatusPaused  = consts.QueueStatusPaused
    QueueStatusStopped = consts.QueueStatusStopped
    QueueStatusError   = consts.QueueStatusError
)
```

#### 已使用consts的常量保持不变
- 队列操作类型常量
- 调度策略常量
- 操作来源常量
- 队列优先级常量

## 迁移优点

### 1. 统一管理
- 所有队列相关常量集中在 `consts/queue.go` 管理
- 避免重复定义，减少代码冗余
- 便于维护和修改

### 2. 提高一致性
- 确保队列常量定义的一致性
- 避免不同模块使用不同的常量值
- 减少因常量不一致导致的错误

### 3. 便于扩展
- 新增队列常量只需要在 `consts/queue.go` 添加
- 其他模块可以直接引用，无需重复定义
- 支持跨模块的常量共享

### 4. 保持兼容性
- 保留原有的常量名称
- 其他模块无需修改代码
- 平滑迁移，无破坏性变更

## 迁移影响

### 1. 影响范围
- `internal/consts/queue.go` - 扩展文件
- `internal/model/queue/queue.go` - 主要修改文件

### 2. 兼容性
- ✅ 向后兼容：其他模块无需修改
- ✅ 编译通过：修复了所有编译错误
- ✅ 功能正常：常量功能保持不变

### 3. 性能影响
- 无性能影响：常量在编译时确定
- 无运行时开销：只是引用方式改变

## 常量分类总结

### 1. 队列状态相关
- **队列状态**: 4个（active, paused, stopped, error）
- **设备队列任务状态**: 6个（queued, executing, paused, completed, failed, canceled）

### 2. 队列类型相关
- **队列类型**: 4个（fifo, priority, delay, lifo）
- **调度策略**: 4个（priority_first, fifo, lifo, weighted）

### 3. 队列操作相关
- **操作类型**: 8个（add, remove, priority_change, position_change, start, pause, resume, cancel）
- **操作来源**: 4个（manual, system, api, scheduler）

### 4. 队列配置相关
- **优先级**: 6个（1, 3, 5, 7, 9, 10）
- **配置参数**: 6个（大小、并发、批处理等）
- **超时设置**: 3个（默认、最小、最大）
- **重试设置**: 3个（默认、最大、最小）
- **监控设置**: 3个（监控间隔、健康检查、指标收集）

## 后续建议

### 1. 继续迁移
- 检查其他model文件中的队列相关常量
- 将业务相关常量迁移到对应的consts文件
- 创建缺失的consts文件

### 2. 规范管理
- 制定队列常量命名规范
- 建立常量审查机制
- 避免重复定义

### 3. 文档更新
- 更新API文档中的队列常量说明
- 维护队列常量使用指南
- 记录常量变更历史

## 总结

本次迁移成功实现了队列常量的统一管理，解决了重复定义问题，提高了代码的一致性和可维护性。迁移过程平滑，无破坏性变更，为后续的队列模块开发奠定了良好基础。

### 迁移统计
- **修改文件**: 2个
- **新增常量**: 8个
- **保持兼容**: 100%
- **编译状态**: ✅ 成功 