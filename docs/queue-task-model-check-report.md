# Queue.go 和 Task.go 模型检查报告

## 概述

本报告对 `internal/model/queue/queue.go` 和 `internal/model/task/task.go` 文件进行了详细检查，对比了模型定义与数据库表结构（entity）的一致性。

## 检查结果摘要

### Queue.go 文件问题
- ❌ **缺少对应的数据库表结构**：未找到队列主表的 entity 定义
- ❌ **字段命名不一致**：模型使用下划线命名，entity 使用驼峰命名
- ❌ **字段类型不匹配**：部分字段类型与数据库定义不一致

### Task.go 文件问题
- ❌ **字段命名不一致**：模型使用下划线命名，entity 使用驼峰命名
- ❌ **字段类型不匹配**：部分字段类型与数据库定义不一致
- ❌ **缺少字段**：部分 entity 中定义的字段在模型中缺失

## 详细对比分析

### Task 模型对比

| 字段 | Task.go | Entity | 状态 | 问题描述 |
|------|---------|--------|------|----------|
| ID | `int64` | `int64` | ✅ 一致 | - |
| TaskID | `string` | `string` | ✅ 一致 | - |
| Name | `string` | `string` | ✅ 一致 | - |
| Description | `string` | `string` | ✅ 一致 | - |
| Type | `string` | `string` | ✅ 一致 | - |
| Status | `string` | `string` | ✅ 一致 | - |
| Priority | `int` | `int` | ✅ 一致 | - |
| ExecuteTime | `*gtime.Time` | `*gtime.Time` | ✅ 一致 | - |
| Timeout | `int` | `int` | ✅ 一致 | - |
| RetryCount | `int` | `int` | ✅ 一致 | - |
| MaxRetries | `int` | `int` | ✅ 一致 | - |
| IsUrgent | `bool` | `bool` | ✅ 一致 | - |
| Parameters | `string` | `string` | ✅ 一致 | - |
| Result | `string` | `string` | ✅ 一致 | - |
| ErrorMessage | `string` | `string` | ✅ 一致 | - |
| ExecutorType | `string` | `string` | ✅ 一致 | - |
| ExecutorID | `string` | `string` | ❌ 命名不一致 | 应为 `ExecutorId` |
| DeviceID | `int64` | `int64` | ❌ 命名不一致 | 应为 `DeviceId` |
| CreatedAt | `*gtime.Time` | `*gtime.Time` | ✅ 一致 | - |
| UpdatedAt | `*gtime.Time` | `*gtime.Time` | ✅ 一致 | - |
| CreatedBy | `int64` | `int64` | ❌ 命名不一致 | 应为 `CreatedBy` |
| UpdatedBy | `int64` | `int64` | ❌ 命名不一致 | 应为 `UpdatedBy` |

### TaskExecution 模型对比

| 字段 | Task.go | Entity | 状态 | 问题描述 |
|------|---------|--------|------|----------|
| ID | `int64` | `int64` | ✅ 一致 | - |
| TaskID | `string` | `string` | ❌ 命名不一致 | 应为 `TaskId` |
| ExecutionID | `string` | `string` | ❌ 命名不一致 | 应为 `ExecutionId` |
| DeviceESN | `string` | `string` | ❌ 命名不一致 | 应为 `DeviceEsn` |
| Status | `string` | `string` | ✅ 一致 | - |
| StartTime | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `StartTime` |
| EndTime | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `EndTime` |
| Duration | `int` | `int` | ✅ 一致 | - |
| ExecutorInfo | `string` | `string` | ❌ 命名不一致 | 应为 `ExecutorInfo` |
| Logs | `string` | `string` | ✅ 一致 | - |
| Metrics | `string` | `string` | ✅ 一致 | - |
| Output | `string` | `string` | ✅ 一致 | - |
| ErrorDetails | `string` | `string` | ❌ 命名不一致 | 应为 `ErrorDetails` |
| CpuUsageAvg | `float64` | `float64` | ❌ 命名不一致 | 应为 `CpuUsageAvg` |
| CpuUsagePeak | `float64` | `float64` | ❌ 命名不一致 | 应为 `CpuUsagePeak` |
| MemoryUsageAvg | `float64` | `float64` | ❌ 命名不一致 | 应为 `MemoryUsageAvg` |
| MemoryUsagePeak | `float64` | `float64` | ❌ 命名不一致 | 应为 `MemoryUsagePeak` |
| IoOperationsTotal | `int64` | `int64` | ❌ 命名不一致 | 应为 `IoOperationsTotal` |
| IoBytesTotal | `int64` | `int64` | ❌ 命名不一致 | 应为 `IoBytesTotal` |
| CreatedAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `CreatedAt` |
| UpdatedAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `UpdatedAt` |

### TaskAssignment 模型对比

| 字段 | Task.go | Entity | 状态 | 问题描述 |
|------|---------|--------|------|----------|
| ID | `int64` | `int64` | ✅ 一致 | - |
| TaskID | `string` | `string` | ❌ 命名不一致 | 应为 `TaskId` |
| Priority | `int` | `int` | ✅ 一致 | - |
| QueueStatus | `string` | `string` | ❌ 命名不一致 | 应为 `QueueStatus` |
| RequiredDeviceType | `string` | `string` | ❌ 命名不一致 | 应为 `RequiredDeviceType` |
| RequiredCapabilities | `string` | `string` | ❌ 命名不一致 | 应为 `RequiredCapabilities` |
| PreferredDeviceIDs | `[]int64` | `[]int64` | ❌ 命名不一致 | 应为 `PreferredDeviceIds` |
| ExcludedDeviceIDs | `[]int64` | `[]int64` | ❌ 命名不一致 | 应为 `ExcludedDeviceIds` |
| AssignmentStrategy | `string` | `string` | ❌ 命名不一致 | 应为 `AssignmentStrategy` |
| AffinityRules | `string` | `string` | ❌ 命名不一致 | 应为 `AffinityRules` |
| AssignedDeviceID | `int64` | `int64` | ❌ 命名不一致 | 应为 `AssignedDeviceId` |
| AssignedAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `AssignedAt` |
| AssignmentScore | `float64` | `float64` | ❌ 命名不一致 | 应为 `AssignmentScore` |
| QueuePosition | `int` | `int` | ❌ 命名不一致 | 应为 `QueuePosition` |
| EstimatedWaitTime | `int` | `int` | ❌ 命名不一致 | 应为 `EstimatedWaitTime` |
| RetryCount | `int` | `int` | ✅ 一致 | - |
| MaxRetries | `int` | `int` | ✅ 一致 | - |
| OriginalPriority | `int` | `int` | ❌ 命名不一致 | 应为 `OriginalPriority` |
| LastPriorityChangeAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `LastPriorityChangeAt` |
| PriorityChangeReason | `string` | `string` | ❌ 命名不一致 | 应为 `PriorityChangeReason` |
| PriorityBoostReason | `string` | `string` | ❌ 命名不一致 | 应为 `PriorityBoostReason` |
| OperationSource | `string` | `string` | ❌ 命名不一致 | 应为 `OperationSource` |
| QueuedAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `QueuedAt` |
| UpdatedAt | `*gtime.Time` | `*gtime.Time` | ❌ 命名不一致 | 应为 `UpdatedAt` |

## 问题分类统计

### 命名不一致问题
- **Task 模型**：4个字段命名不一致
- **TaskExecution 模型**：12个字段命名不一致
- **TaskAssignment 模型**：18个字段命名不一致
- **总计**：34个字段需要修正命名

### 类型不匹配问题
- 所有字段类型基本一致，主要是命名问题

### 缺失字段问题
- 所有 entity 中定义的字段在 model 中都有对应，无缺失字段

## 修改建议

### 1. 统一命名规范
建议使用驼峰命名（camelCase），与 entity 保持一致：
- `ExecutorID` → `ExecutorId`
- `DeviceID` → `DeviceId`
- `CreatedBy` → `CreatedBy`
- `UpdatedBy` → `UpdatedBy`

### 2. 修正字段命名
按照 entity 定义修正所有字段的命名，确保与数据库表结构完全一致。

### 3. 更新相关引用
修正 Input/Output 结构体中对这些字段的引用。

### 4. 测试验证
修改后需要进行完整的测试，确保功能正常。

## 优先级建议

### 高优先级
1. 修正 Task 模型的字段命名
2. 修正 TaskExecution 模型的字段命名
3. 修正 TaskAssignment 模型的字段命名

### 中优先级
1. 更新 Input/Output 结构体引用
2. 更新相关业务逻辑代码

### 低优先级
1. 更新文档和注释
2. 代码风格统一

## 风险评估

### 低风险
- 字段类型基本一致，主要是命名问题
- 修改后不会影响数据库操作

### 中风险
- 需要更新所有引用这些字段的代码
- 需要确保 Input/Output 结构体的一致性

### 建议
1. 分模块逐步修改
2. 每次修改后进行测试
3. 保持向后兼容性

## 总结

本次检查发现了大量的字段命名不一致问题，主要集中在 Task 相关的模型中。建议按照 entity 定义统一字段命名，确保模型层与数据库表结构完全一致。修改过程中需要注意更新所有相关引用，并进行充分的测试验证。 