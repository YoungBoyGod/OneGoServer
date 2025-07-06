# 任务模块常量迁移总结

## 迁移概述

本次迁移将 `internal/model/task/task.go` 中定义的常量迁移到 `internal/consts/task.go` 中，实现任务相关常量的统一管理。

## 迁移内容

### 1. 检查 consts/task.go

#### 已存在的常量类型
- **任务状态常量**: 9个状态常量（pending, queued, assigning, assigned, dispatching, running, completed, failed, canceled）
- **任务类型常量**: 4个类型常量（backup, sync, monitor, custom）
- **执行器类型常量**: 4个类型常量（local, remote, container, lambda）
- **执行状态常量**: 5个状态常量（started, running, completed, failed, canceled）
- **队列状态常量**: 5个状态常量（queued, assigning, assigned, failed, canceled）
- **分配历史动作常量**: 6个动作常量（queued, assigned, reassigned, failed, completed, canceled）
- **任务依赖类型常量**: 3个类型常量（before, after, parallel）
- **任务优先级常量**: 6个优先级常量（1, 3, 5, 7, 9, 10）
- **任务重试常量**: 3个重试常量（默认重试次数、最大重试次数、默认超时）
- **任务调度常量**: 3个调度常量（默认调度间隔、最大并发任务数、批处理大小）

#### 常量列表
```go
// 任务状态常量
TaskStatusPending, TaskStatusQueued, TaskStatusAssigning, TaskStatusAssigned,
TaskStatusDispatching, TaskStatusRunning, TaskStatusCompleted, TaskStatusFailed, TaskStatusCanceled

// 任务类型常量
TaskTypeBackup, TaskTypeSync, TaskTypeMonitor, TaskTypeCustom

// 执行器类型常量
ExecutorTypeLocal, ExecutorTypeRemote, ExecutorTypeContainer, ExecutorTypeLambda

// 执行状态常量
ExecutionStatusStarted, ExecutionStatusRunning, ExecutionStatusCompleted,
ExecutionStatusFailed, ExecutionStatusCanceled

// 队列状态常量
QueueStatusQueued, QueueStatusAssigning, QueueStatusAssigned, QueueStatusFailed, QueueStatusCanceled

// 分配历史动作常量
AssignmentActionQueued, AssignmentActionAssigned, AssignmentActionReassigned,
AssignmentActionFailed, AssignmentActionCompleted, AssignmentActionCanceled

// 任务依赖类型常量
DependencyTypeBefore, DependencyTypeAfter, DependencyTypeParallel

// 任务优先级常量
TaskPriorityLowest(1), TaskPriorityLow(3), TaskPriorityNormal(5), 
TaskPriorityHigh(7), TaskPriorityUrgent(9), TaskPriorityHighest(10)

// 任务重试常量
DefaultMaxRetries(3), MaxAllowedRetries(10), DefaultTimeout(300)

// 任务调度常量
DefaultSchedulingInterval(30), MaxConcurrentTasks(100), TaskBatchSize(50)
```

### 2. 更新 model/task/task.go

#### 修改内容
- **移除重复定义**: 删除原有的任务优先级常量定义
- **使用consts引用**: 改为使用 `consts.XXX` 的方式引用常量
- **保持向后兼容**: 保留原有的常量名称，确保其他模块不受影响

#### 修改示例
```go
// 修改前
const (
    TaskPriorityLowest  = 1
    TaskPriorityLow     = 3
    TaskPriorityNormal  = 5
    TaskPriorityHigh    = 7
    TaskPriorityUrgent  = 9
    TaskPriorityHighest = 10
)

// 修改后
const (
    TaskPriorityLowest  = consts.TaskPriorityLowest
    TaskPriorityLow     = consts.TaskPriorityLow
    TaskPriorityNormal  = consts.TaskPriorityNormal
    TaskPriorityHigh    = consts.TaskPriorityHigh
    TaskPriorityUrgent  = consts.TaskPriorityUrgent
    TaskPriorityHighest = consts.TaskPriorityHighest
)
```

#### 已使用consts的常量保持不变
- 任务状态常量
- 任务类型常量
- 执行器类型常量
- 执行状态常量
- 队列状态常量
- 分配历史动作常量
- 任务依赖类型常量

## 迁移优点

### 1. 统一管理
- 所有任务相关常量集中在 `consts/task.go` 管理
- 避免重复定义，减少代码冗余
- 便于维护和修改

### 2. 提高一致性
- 确保任务常量定义的一致性
- 避免不同模块使用不同的常量值
- 减少因常量不一致导致的错误

### 3. 便于扩展
- 新增任务常量只需要在 `consts/task.go` 添加
- 其他模块可以直接引用，无需重复定义
- 支持跨模块的常量共享

### 4. 保持兼容性
- 保留原有的常量名称
- 其他模块无需修改代码
- 平滑迁移，无破坏性变更

## 迁移影响

### 1. 影响范围
- `internal/model/task/task.go` - 主要修改文件
- `internal/consts/task.go` - 已存在，无需修改

### 2. 兼容性
- ✅ 向后兼容：其他模块无需修改
- ✅ 编译通过：修复了所有编译错误
- ✅ 功能正常：常量功能保持不变

### 3. 性能影响
- 无性能影响：常量在编译时确定
- 无运行时开销：只是引用方式改变

## 常量分类总结

### 1. 任务状态相关（14个）
- **任务状态**: 9个（pending, queued, assigning, assigned, dispatching, running, completed, failed, canceled）
- **执行状态**: 5个（started, running, completed, failed, canceled）

### 2. 任务类型相关（7个）
- **任务类型**: 4个（backup, sync, monitor, custom）
- **执行器类型**: 4个（local, remote, container, lambda）

### 3. 任务队列相关（11个）
- **队列状态**: 5个（queued, assigning, assigned, failed, canceled）
- **分配历史动作**: 6个（queued, assigned, reassigned, failed, completed, canceled）

### 4. 任务配置相关（12个）
- **优先级**: 6个（1, 3, 5, 7, 9, 10）
- **重试设置**: 3个（默认重试次数、最大重试次数、默认超时）
- **调度设置**: 3个（默认调度间隔、最大并发任务数、批处理大小）

### 5. 任务依赖相关（3个）
- **依赖类型**: 3个（before, after, parallel）

## 后续建议

### 1. 继续迁移
- 检查其他model文件中的任务相关常量
- 将业务相关常量迁移到对应的consts文件
- 创建缺失的consts文件

### 2. 规范管理
- 制定任务常量命名规范
- 建立常量审查机制
- 避免重复定义

### 3. 文档更新
- 更新API文档中的任务常量说明
- 维护任务常量使用指南
- 记录常量变更历史

## 总结

本次迁移成功实现了任务常量的统一管理，解决了重复定义问题，提高了代码的一致性和可维护性。迁移过程平滑，无破坏性变更，为后续的任务模块开发奠定了良好基础。

### 迁移统计
- **修改文件**: 1个
- **迁移常量**: 6个（任务优先级常量）
- **保持兼容**: 100%
- **编译状态**: ✅ 成功
- **迁移状态**: ✅ 完成（所有常量已使用consts层） 