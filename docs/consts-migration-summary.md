# 常量迁移总结

## 迁移概述

本次迁移将 `internal/model/common/common.go` 中定义的常量迁移到 `internal/consts/common.go` 中，实现常量的统一管理和避免重复定义。

## 迁移内容

### 1. 扩展 consts/common.go

#### 新增常量类型
- **通用状态常量**: 11个状态常量（active, inactive, running, stopped等）
- **通用操作类型常量**: 13个操作类型常量（create, update, delete等）
- **通用日志级别常量**: 5个日志级别常量（debug, info, warn, error, critical）
- **通用优先级常量**: 6个优先级常量（1-10的数值常量）

#### 常量列表
```go
// 通用状态常量
StatusActive, StatusInactive, StatusRunning, StatusStopped, 
StatusPaused, StatusCompleted, StatusFailed, StatusCanceled, StatusLocked

// 通用排序常量
SortOrderAsc, SortOrderDesc

// 通用操作类型常量
ActionTypeCreate, ActionTypeUpdate, ActionTypeDelete, ActionTypeView,
ActionTypeStart, ActionTypeStop, ActionTypePause, ActionTypeResume,
ActionTypeCancel, ActionTypeRestart, ActionTypeRetry, ActionTypeAssign, ActionTypeUnassign

// 通用日志级别常量
LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelCritical

// 通用优先级常量
PriorityLowest(1), PriorityLow(3), PriorityNormal(5), PriorityHigh(7), PriorityUrgent(9), PriorityHighest(10)
```

### 2. 更新 model/common/common.go

#### 修改内容
- **添加导入**: 导入 `OneGfServer/internal/consts` 包
- **移除重复定义**: 删除原有的常量定义
- **使用consts引用**: 改为使用 `consts.XXX` 的方式引用常量
- **保持向后兼容**: 保留原有的常量名称，确保其他模块不受影响

#### 修改示例
```go
// 修改前
const (
    StatusActive = "active"
    StatusInactive = "inactive"
    // ...
)

// 修改后
const (
    StatusActive = consts.StatusActive
    StatusInactive = consts.StatusInactive
    // ...
)
```

### 3. 修复冲突问题

#### 问题描述
- `consts/device.go` 中定义了日志级别常量，与 `consts/common.go` 中的通用日志级别常量重复
- 导致编译错误：`LogLevelDebug redeclared in this block`

#### 解决方案
- 移除 `consts/device.go` 中重复的日志级别常量定义
- 改为使用通用日志级别常量，并添加设备特有的常量
- 保留设备特有的 `DeviceLogLevelFatal` 常量

```go
// 修改前
const (
    LogLevelDebug = "debug"
    LogLevelInfo = "info"
    // ...
)

// 修改后
const (
    DeviceLogLevelDebug = LogLevelDebug
    DeviceLogLevelInfo = LogLevelInfo
    DeviceLogLevelFatal = "fatal" // 设备特有的致命级别
)
```

## 迁移优点

### 1. 统一管理
- 所有通用常量集中在 `consts` 层管理
- 避免重复定义，减少代码冗余
- 便于维护和修改

### 2. 提高一致性
- 确保常量定义的一致性
- 避免不同模块使用不同的常量值
- 减少因常量不一致导致的错误

### 3. 便于扩展
- 新增常量只需要在 `consts` 层添加
- 其他模块可以直接引用，无需重复定义
- 支持跨模块的常量共享

### 4. 保持兼容性
- 保留原有的常量名称
- 其他模块无需修改代码
- 平滑迁移，无破坏性变更

## 迁移影响

### 1. 影响范围
- `internal/model/common/common.go` - 主要修改文件
- `internal/consts/common.go` - 扩展文件
- `internal/consts/device.go` - 修复冲突文件

### 2. 兼容性
- ✅ 向后兼容：其他模块无需修改
- ✅ 编译通过：修复了所有编译错误
- ✅ 功能正常：常量功能保持不变

### 3. 性能影响
- 无性能影响：常量在编译时确定
- 无运行时开销：只是引用方式改变

## 后续建议

### 1. 继续迁移
- 检查其他model文件中的常量定义
- 将业务相关常量迁移到对应的consts文件
- 创建缺失的consts文件（如user.go）

### 2. 规范管理
- 制定常量命名规范
- 建立常量审查机制
- 避免重复定义

### 3. 文档更新
- 更新API文档中的常量说明
- 维护常量使用指南
- 记录常量变更历史

## 总结

本次迁移成功实现了常量的统一管理，解决了重复定义问题，提高了代码的一致性和可维护性。迁移过程平滑，无破坏性变更，为后续的常量管理奠定了良好基础。 