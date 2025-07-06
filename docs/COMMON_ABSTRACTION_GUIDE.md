# 通用抽象层设计指南

## 概述

本文档描述了OneGoServer项目中通用抽象层的设计理念、使用方法和最佳实践。通过抽象通用部分，我们实现了代码复用、一致性维护和开发效率提升。

## 设计目标

### 1. 代码复用
- 消除重复的结构体定义
- 统一常量和枚举值
- 标准化接口和工具函数

### 2. 一致性维护
- 统一的命名规范
- 标准化的字段定义
- 一致的验证规则

### 3. 开发效率
- 减少样板代码
- 提供通用工具函数
- 简化新模块开发

## 通用抽象层结构

### 1. 通用常量 (`internal/model/common/common.go`)

#### 状态常量
```go
const (
    StatusActive   = "active"
    StatusInactive = "inactive"
    StatusPending  = "pending"
    StatusRunning  = "running"
    StatusStopped  = "stopped"
    StatusPaused   = "paused"
    StatusError    = "error"
    StatusCompleted = "completed"
    StatusFailed   = "failed"
    StatusCanceled = "canceled"
    StatusLocked   = "locked"
)
```

#### 排序常量
```go
const (
    SortOrderAsc  = "asc"
    SortOrderDesc = "desc"
)
```

#### 操作类型常量
```go
const (
    ActionTypeCreate   = "create"
    ActionTypeUpdate   = "update"
    ActionTypeDelete   = "delete"
    ActionTypeView     = "view"
    ActionTypeStart    = "start"
    ActionTypeStop     = "stop"
    ActionTypePause    = "pause"
    ActionTypeResume   = "resume"
    ActionTypeCancel   = "cancel"
    ActionTypeRestart  = "restart"
    ActionTypeRetry    = "retry"
    ActionTypeAssign   = "assign"
    ActionTypeUnassign = "unassign"
)
```

#### 优先级常量
```go
const (
    PriorityLowest  = 1
    PriorityLow     = 3
    PriorityNormal  = 5
    PriorityHigh    = 7
    PriorityUrgent  = 9
    PriorityHighest = 10
)
```

### 2. 通用基础结构体

#### BaseEntity - 基础实体
```go
type BaseEntity struct {
    ID        int64       `json:"id"`
    CreatedAt *gtime.Time `json:"created_at"`
    UpdatedAt *gtime.Time `json:"updated_at"`
    CreatedBy int64       `json:"created_by"`
    UpdatedBy int64       `json:"updated_by"`
}
```

#### BaseEntityWithStatus - 带状态的基础实体
```go
type BaseEntityWithStatus struct {
    BaseEntity
    Status string `json:"status"`
}
```

#### BaseEntityWithName - 带名称的基础实体
```go
type BaseEntityWithName struct {
    BaseEntity
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

#### BaseEntityWithPriority - 带优先级的基础实体
```go
type BaseEntityWithPriority struct {
    BaseEntity
    Priority int `json:"priority"`
}
```

### 3. 通用Input/Output结构体

#### BaseCreateInput/Output - 基础创建操作
```go
type BaseCreateInput struct {
    Entity interface{} `json:"entity"`
}

type BaseCreateOutput struct {
    ID      string `json:"id"`
    Message string `json:"message"`
}
```

#### BaseGetByIDInput/Output - 基础根据ID获取操作
```go
type BaseGetByIDInput struct {
    ID int64 `json:"id"`
}

type BaseGetByIDOutput struct {
    Entity interface{} `json:"entity"`
}
```

#### BaseListInput/Output - 基础列表操作
```go
type BaseListInput struct {
    Filter     interface{} `json:"filter"`
    Sort       interface{} `json:"sort"`
    Pagination *PaginationOption `json:"pagination"`
}

type BaseListOutput struct {
    List  interface{} `json:"list"`
    Total int64       `json:"total"`
}
```

### 4. 通用过滤和排序选项

#### PaginationOption - 分页选项
```go
type PaginationOption struct {
    Page int `json:"page"`
    Size int `json:"size"`
}
```

#### BaseFilter - 基础过滤条件
```go
type BaseFilter struct {
    Status    []string    `json:"status"`
    Type      []string    `json:"type"`
    Keyword   *string     `json:"keyword"`
    StartTime *time.Time  `json:"start_time"`
    EndTime   *time.Time  `json:"end_time"`
    Enabled   *bool       `json:"enabled"`
}
```

#### BaseSortOption - 基础排序选项
```go
type BaseSortOption struct {
    Field string `json:"field"` // id, name, status, priority, created_at, updated_at
    Order string `json:"order"` // asc, desc
}
```

### 5. 通用接口定义

#### Entity - 实体接口
```go
type Entity interface {
    GetID() int64
    GetCreatedAt() *gtime.Time
    GetUpdatedAt() *gtime.Time
}
```

#### Filter - 过滤接口
```go
type Filter interface {
    GetStatus() []string
    GetType() []string
    GetKeyword() *string
    GetStartTime() *time.Time
    GetEndTime() *time.Time
}
```

#### Sort - 排序接口
```go
type Sort interface {
    GetField() string
    GetOrder() string
}
```

### 6. 通用工具函数

#### 验证函数
```go
// 检查状态是否有效
func IsValidStatus(status string) bool

// 检查排序方向是否有效
func IsValidSortOrder(order string) bool

// 检查优先级是否有效
func IsValidPriority(priority int) bool

// 检查日志级别是否有效
func IsValidLogLevel(level string) bool
```

#### 默认值函数
```go
// 获取默认分页选项
func GetDefaultPagination() *PaginationOption

// 获取默认排序选项
func GetDefaultSort(field string) *BaseSortOption
```

## 使用方法

### 1. 继承通用结构体

#### 设备模块示例
```go
// Device 设备主表模型 - 继承BaseEntityWithName
type Device struct {
    common.BaseEntityWithName
    DeviceID         string      `json:"device_id"`
    DeviceType       string      `json:"device_type"`
    Model            string      `json:"model"`
    // ... 其他设备特定字段
}

// DeviceFilter 设备过滤条件 - 继承BaseFilter
type DeviceFilter struct {
    common.BaseFilter
    Manufacturer *string    `json:"manufacturer"`
    Model        *string    `json:"model"`
    Protocol     *string    `json:"protocol"`
    // ... 其他设备特定过滤条件
}
```

#### 任务模块示例
```go
// Task 任务模型 - 继承BaseEntityWithPriority
type Task struct {
    common.BaseEntityWithPriority
    TaskID       string      `json:"task_id"`
    TaskType     string      `json:"task_type"`
    Status       string      `json:"status"`
    // ... 其他任务特定字段
}

// TaskFilter 任务过滤条件 - 继承BaseFilter
type TaskFilter struct {
    common.BaseFilter
    Priority     *int       `json:"priority"`
    DeviceID     *string    `json:"device_id"`
    ExecutorType *string    `json:"executor_type"`
    // ... 其他任务特定过滤条件
}
```

### 2. 使用通用Input/Output结构体

#### 创建操作
```go
// CreateDeviceInput 创建设备输入 - 继承BaseCreateInput
type CreateDeviceInput struct {
    common.BaseCreateInput
    Device *Device `json:"device"`
}

// CreateDeviceOutput 创建设备输出 - 继承BaseCreateOutput
type CreateDeviceOutput struct {
    common.BaseCreateOutput
    DeviceID string `json:"device_id"`
}
```

#### 列表操作
```go
// GetDeviceListInput 获取设备列表输入 - 继承BaseListInput
type GetDeviceListInput struct {
    common.BaseListInput
    Filter     *DeviceFilter     `json:"filter"`
    Sort       *DeviceSortOption `json:"sort"`
    Pagination *common.PaginationOption `json:"pagination"`
}

// GetDeviceListOutput 获取设备列表输出 - 继承BaseListOutput
type GetDeviceListOutput struct {
    common.BaseListOutput
    List  []Device `json:"list"`
    Total int64    `json:"total"`
}
```

### 3. 实现通用接口

#### Entity接口实现
```go
// GetID 实现Entity接口
func (d *Device) GetID() int64 {
    return d.ID
}

// GetCreatedAt 实现Entity接口
func (d *Device) GetCreatedAt() *gtime.Time {
    return d.CreatedAt
}

// GetUpdatedAt 实现Entity接口
func (d *Device) GetUpdatedAt() *gtime.Time {
    return d.UpdatedAt
}
```

#### Filter接口实现
```go
// GetStatus 实现Filter接口
func (f *DeviceFilter) GetStatus() []string {
    return f.Status
}

// GetType 实现Filter接口
func (f *DeviceFilter) GetType() []string {
    return f.Type
}

// GetKeyword 实现Filter接口
func (f *DeviceFilter) GetKeyword() *string {
    return f.Keyword
}

// GetStartTime 实现Filter接口
func (f *DeviceFilter) GetStartTime() *time.Time {
    return f.StartTime
}

// GetEndTime 实现Filter接口
func (f *DeviceFilter) GetEndTime() *time.Time {
    return f.EndTime
}
```

### 4. 使用通用工具函数

#### 验证数据
```go
// 验证设备状态
if !common.IsValidStatus(device.Status) {
    return errors.New("invalid device status")
}

// 验证排序方向
if !common.IsValidSortOrder(sort.Order) {
    return errors.New("invalid sort order")
}

// 验证优先级
if !common.IsValidPriority(task.Priority) {
    return errors.New("invalid priority")
}
```

#### 设置默认值
```go
// 设置默认分页
if input.Pagination == nil {
    input.Pagination = common.GetDefaultPagination()
}

// 设置默认排序
if input.Sort == nil {
    input.Sort = common.GetDefaultSort("created_at")
}
```

## 最佳实践

### 1. 命名规范

#### 常量命名
- 使用模块前缀：`DeviceStatusOnline`、`TaskStatusRunning`
- 继承通用常量：`DeviceStatusOnline = common.StatusActive`

#### 结构体命名
- 继承通用结构体：`Device struct { common.BaseEntityWithName }`
- 特定字段：`DeviceID`、`DeviceType`、`Model`

#### 接口命名
- 使用通用接口：`Entity`、`Filter`、`Sort`
- 实现接口方法：`GetID()`、`GetStatus()`、`GetField()`

### 2. 字段组织

#### 基础字段在前
```go
type Device struct {
    common.BaseEntityWithName  // 基础字段
    // 业务字段
    DeviceID   string `json:"device_id"`
    DeviceType string `json:"device_type"`
    // 配置字段
    Config     *DeviceConfig `json:"config"`
}
```

#### 过滤条件扩展
```go
type DeviceFilter struct {
    common.BaseFilter  // 通用过滤条件
    // 设备特定过滤条件
    Manufacturer *string    `json:"manufacturer"`
    Model        *string    `json:"model"`
    Protocol     *string    `json:"protocol"`
}
```

### 3. 验证规则

#### 使用通用验证
```go
// 验证状态
if !common.IsValidStatus(status) {
    return common.NewCommonError("INVALID_STATUS", "invalid status", "")
}

// 验证优先级
if !common.IsValidPriority(priority) {
    return common.NewCommonError("INVALID_PRIORITY", "invalid priority", "")
}
```

#### 模块特定验证
```go
// 验证设备类型
func IsValidDeviceType(deviceType string) bool {
    validTypes := []string{
        DeviceTypeSensor, DeviceTypeCamera, DeviceTypeActuator, DeviceTypeGateway,
    }
    for _, validType := range validTypes {
        if deviceType == validType {
            return true
        }
    }
    return false
}
```

### 4. 错误处理

#### 使用通用错误
```go
// 创建通用错误
func NewDeviceNotFoundError(deviceID string) *common.CommonError {
    return common.NewCommonError(
        "DEVICE_NOT_FOUND",
        "device not found",
        fmt.Sprintf("device_id: %s", deviceID),
    )
}

// 使用通用错误
if device == nil {
    return NewDeviceNotFoundError(deviceID)
}
```

### 5. 响应处理

#### 使用通用响应
```go
// 成功响应
func (s *sDevice) GetDevice(ctx context.Context, input *GetDeviceInput) (*GetDeviceOutput, error) {
    device, err := s.getDeviceByID(ctx, input.ID)
    if err != nil {
        return nil, err
    }
    
    return &GetDeviceOutput{
        Device: device,
    }, nil
}

// 错误响应
if device == nil {
    return nil, common.NewCommonError("DEVICE_NOT_FOUND", "device not found", "")
}
```

## 重构指南

### 1. 识别重复代码

#### 查找重复结构体
- `PaginationOption` 在多个模块中重复
- `BaseFilter` 类似的过滤条件
- `BaseSortOption` 类似的排序选项

#### 查找重复常量
- 状态常量：`active`、`inactive`、`pending` 等
- 排序常量：`asc`、`desc`
- 优先级常量：`1`、`3`、`5`、`7`、`9`、`10`

### 2. 提取通用部分

#### 创建通用结构体
```go
// 提取通用字段
type BaseEntity struct {
    ID        int64       `json:"id"`
    CreatedAt *gtime.Time `json:"created_at"`
    UpdatedAt *gtime.Time `json:"updated_at"`
    CreatedBy int64       `json:"created_by"`
    UpdatedBy int64       `json:"updated_by"`
}
```

#### 创建通用接口
```go
// 提取通用方法
type Entity interface {
    GetID() int64
    GetCreatedAt() *gtime.Time
    GetUpdatedAt() *gtime.Time
}
```

### 3. 重构现有代码

#### 继承通用结构体
```go
// 重构前
type Device struct {
    ID        int64       `json:"id"`
    CreatedAt *gtime.Time `json:"created_at"`
    UpdatedAt *gtime.Time `json:"updated_at"`
    CreatedBy int64       `json:"created_by"`
    UpdatedBy int64       `json:"updated_by"`
    Name      string      `json:"name"`
    // ... 其他字段
}

// 重构后
type Device struct {
    common.BaseEntityWithName
    DeviceID   string `json:"device_id"`
    DeviceType string `json:"device_type"`
    // ... 其他字段
}
```

#### 实现通用接口
```go
// 实现Entity接口
func (d *Device) GetID() int64 {
    return d.ID
}

func (d *Device) GetCreatedAt() *gtime.Time {
    return d.CreatedAt
}

func (d *Device) GetUpdatedAt() *gtime.Time {
    return d.UpdatedAt
}
```

### 4. 更新依赖

#### 更新导入
```go
import (
    "OneGfServer/internal/model/common"
)
```

#### 更新引用
```go
// 使用通用常量
DeviceStatusOnline = common.StatusActive

// 使用通用结构体
PaginationOption = common.PaginationOption

// 使用通用函数
common.IsValidStatus(status)
```

## 优势总结

### 1. 代码复用
- **减少重复代码**：消除重复的结构体定义
- **统一常量管理**：集中管理状态、排序等常量
- **标准化接口**：提供统一的Entity、Filter、Sort接口

### 2. 一致性维护
- **统一命名规范**：所有模块使用相同的命名规则
- **标准化字段**：基础字段在所有模块中保持一致
- **统一验证规则**：使用相同的验证逻辑

### 3. 开发效率
- **快速开发**：新模块可以快速继承通用结构体
- **减少错误**：使用经过验证的通用组件
- **简化维护**：集中管理通用逻辑，降低维护成本
- **增强扩展性**：新模块可以轻松继承通用功能

### 4. 扩展性
- **易于扩展**：新增模块只需继承通用结构体
- **向后兼容**：通用接口保证向后兼容性
- **灵活配置**：支持模块特定的扩展字段

## 注意事项

### 1. 避免过度抽象
- 只抽象真正通用的部分
- 保持模块的独立性
- 避免创建过于复杂的抽象层

### 2. 保持向后兼容
- 通用接口的修改要谨慎
- 提供迁移指南
- 保持API的稳定性

### 3. 文档维护
- 及时更新使用指南
- 提供示例代码
- 记录最佳实践

### 4. 测试覆盖
- 为通用组件编写测试
- 确保接口实现的正确性
- 验证向后兼容性

## 总结

通用抽象层的设计为OneGoServer项目提供了强大的代码复用和一致性维护能力。通过合理使用通用结构体、接口和工具函数，我们可以：

1. **提高开发效率**：减少重复代码，快速开发新功能
2. **保证代码质量**：统一的验证规则和错误处理
3. **简化维护工作**：集中管理通用逻辑，降低维护成本
4. **增强扩展性**：新模块可以轻松继承通用功能

遵循本指南的最佳实践，可以充分发挥通用抽象层的优势，构建高质量、可维护的代码库。 