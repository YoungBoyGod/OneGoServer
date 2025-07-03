package task

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/consts"
	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
	"gorm.io/gorm"
)

// JSONB 自定义JSON类型，用于PostgreSQL的JSONB字段
type JSONB map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Task 任务主表模型
type Task struct {
	ID          int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Type        string  `gorm:"type:varchar(50);not null" json:"type"`
	Status      string  `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	Priority    int     `gorm:"not null;default:5" json:"priority"`

	// 执行配置
	ExecuteTime *time.Time `gorm:"type:timestamp" json:"execute_time,omitempty"`
	Timeout     int        `gorm:"default:300" json:"timeout"`
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	MaxRetries  int        `gorm:"default:3" json:"max_retries"`
	IsUrgent    bool       `gorm:"default:false" json:"is_urgent"`

	// 任务参数和结果 (JSON格式)
	Parameters   JSONB   `gorm:"type:jsonb" json:"parameters,omitempty"`
	Result       JSONB   `gorm:"type:jsonb" json:"result,omitempty"`
	ErrorMessage *string `gorm:"type:text" json:"error_message,omitempty"`

	// 执行信息
	ExecutorType *string `gorm:"type:varchar(50)" json:"executor_type,omitempty"`
	ExecutorID   *string `gorm:"type:varchar(100)" json:"executor_id,omitempty"`
	DeviceID     *int64  `gorm:"index" json:"device_id,omitempty"`

	// 审计字段
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy *int64    `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *int64    `gorm:"index" json:"updated_by,omitempty"`

	// 关联关系
	Device     *Device         `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
	Executions []TaskExecution `gorm:"foreignKey:TaskID;references:ID" json:"executions,omitempty"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// BeforeUpdate GORM钩子，更新时自动设置更新时间
func (t *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}

// TaskExecution 任务执行记录模型
type TaskExecution struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      int64  `gorm:"not null;index" json:"task_id"`
	ExecutionID string `gorm:"type:varchar(100);uniqueIndex;not null" json:"execution_id"`

	// 分配信息
	DeviceESN *string `gorm:"type:varchar(100)" json:"device_esn,omitempty"`

	// 执行状态
	Status    string     `gorm:"type:varchar(20);not null;default:started" json:"status"`
	StartTime time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"start_time"`
	EndTime   *time.Time `gorm:"type:timestamp" json:"end_time,omitempty"`
	Duration  *int       `gorm:"comment:执行耗时(秒)" json:"duration,omitempty"`

	// 执行详情
	ExecutorInfo *JSONB  `gorm:"type:jsonb" json:"executor_info,omitempty"`
	Logs         *string `gorm:"type:text" json:"logs,omitempty"`
	Metrics      *JSONB  `gorm:"type:jsonb" json:"metrics,omitempty"`
	Output       *JSONB  `gorm:"type:jsonb" json:"output,omitempty"`
	ErrorDetails *JSONB  `gorm:"type:jsonb" json:"error_details,omitempty"`

	// 资源使用统计 (聚合数据)
	CPUUsageAvg       *float64 `gorm:"type:numeric(5,2)" json:"cpu_usage_avg,omitempty"`
	CPUUsagePeak      *float64 `gorm:"type:numeric(5,2)" json:"cpu_usage_peak,omitempty"`
	MemoryUsageAvg    *float64 `gorm:"type:numeric(10,2)" json:"memory_usage_avg,omitempty"`
	MemoryUsagePeak   *float64 `gorm:"type:numeric(10,2)" json:"memory_usage_peak,omitempty"`
	IOOperationsTotal *int64   `json:"io_operations_total,omitempty"`
	IOBytesTotal      *int64   `json:"io_bytes_total,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	Task *Task `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`
}

// TableName 指定表名
func (TaskExecution) TableName() string {
	return "task_executions"
}

// BeforeCreate GORM钩子，创建前自动生成执行ID
func (te *TaskExecution) BeforeCreate(tx *gorm.DB) (err error) {
	if te.ExecutionID == "" {
		te.ExecutionID = utils.GenerateExecutionID()
	}
	return
}

// AfterUpdate GORM钩子，结束时自动计算执行时长
func (te *TaskExecution) AfterUpdate(tx *gorm.DB) (err error) {
	if te.EndTime != nil && te.Duration == nil {
		duration := int(te.EndTime.Sub(te.StartTime).Seconds())
		te.Duration = &duration
		// 静默更新，避免无限递归
		tx.Model(te).UpdateColumn("duration", duration)
	}
	return
}

// Device 设备模型 (在task包中引用)
type Device struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID string `gorm:"type:varchar(100);uniqueIndex;not null" json:"device_id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Type     string `gorm:"type:varchar(50);not null" json:"type"`
	Status   string `gorm:"type:varchar(20);not null;default:offline" json:"status"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// 任务状态常量 - 使用统一常量
const (
	TaskStatusPending     = consts.TaskStatusPending
	TaskStatusQueued      = consts.TaskStatusQueued
	TaskStatusAssigning   = consts.TaskStatusAssigning
	TaskStatusAssigned    = consts.TaskStatusAssigned
	TaskStatusDispatching = consts.TaskStatusDispatching
	TaskStatusRunning     = consts.TaskStatusRunning
	TaskStatusCompleted   = consts.TaskStatusCompleted
	TaskStatusFailed      = consts.TaskStatusFailed
	TaskStatusCanceled    = consts.TaskStatusCanceled
)

// 任务类型常量 - 使用统一常量
const (
	TaskTypeBackup  = consts.TaskTypeBackup
	TaskTypeSync    = consts.TaskTypeSync
	TaskTypeMonitor = consts.TaskTypeMonitor
	TaskTypeCustom  = consts.TaskTypeCustom
)

// 执行器类型常量 - 使用统一常量
const (
	ExecutorTypeLocal     = consts.ExecutorTypeLocal
	ExecutorTypeRemote    = consts.ExecutorTypeRemote
	ExecutorTypeContainer = consts.ExecutorTypeContainer
	ExecutorTypeLambda    = consts.ExecutorTypeLambda
)

// 执行状态常量 - 使用统一常量
const (
	ExecutionStatusStarted   = consts.ExecutionStatusStarted
	ExecutionStatusRunning   = consts.ExecutionStatusRunning
	ExecutionStatusCompleted = consts.ExecutionStatusCompleted
	ExecutionStatusFailed    = consts.ExecutionStatusFailed
	ExecutionStatusCanceled  = consts.ExecutionStatusCanceled
)

// TaskFilter 任务查询过滤器
type TaskFilter struct {
	Status    []string   `json:"status,omitempty"`
	Type      []string   `json:"type,omitempty"`
	Priority  *int       `json:"priority,omitempty"`
	DeviceID  *int64     `json:"device_id,omitempty"`
	CreatedBy *int64     `json:"created_by,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Keyword   *string    `json:"keyword,omitempty"` // 搜索名称和描述
}

// TaskSortOption 任务排序选项
type TaskSortOption struct {
	Field string `json:"field"` // id, name, status, priority, created_at, updated_at
	Order string `json:"order"` // asc, desc
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page int `json:"page" default:"1"`
	Size int `json:"size" default:"20"`
}

// TaskStatistics 任务统计信息
type TaskStatistics struct {
	Total     int64            `json:"total"`
	ByStatus  map[string]int64 `json:"by_status"`
	ByType    map[string]int64 `json:"by_type"`
	Recent24h int64            `json:"recent_24h"`
	Success   int64            `json:"success"`
	Failed    int64            `json:"failed"`
}

// TaskCreateRequest 创建任务请求
type TaskCreateRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Description  *string                `json:"description,omitempty"`
	Type         string                 `json:"type" binding:"required"`
	Priority     *int                   `json:"priority,omitempty"`
	ExecuteTime  *time.Time             `json:"execute_time,omitempty"`
	Timeout      *int                   `json:"timeout,omitempty"`
	MaxRetries   *int                   `json:"max_retries,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	ExecutorType *string                `json:"executor_type,omitempty"`
	ExecutorID   *string                `json:"executor_id,omitempty"`
	DeviceID     *int64                 `json:"device_id,omitempty"`
}

// TaskUpdateRequest 更新任务请求
type TaskUpdateRequest struct {
	Name         *string                `json:"name,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Priority     *int                   `json:"priority,omitempty"`
	ExecuteTime  *time.Time             `json:"execute_time,omitempty"`
	Timeout      *int                   `json:"timeout,omitempty"`
	MaxRetries   *int                   `json:"max_retries,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	ExecutorType *string                `json:"executor_type,omitempty"`
	ExecutorID   *string                `json:"executor_id,omitempty"`
	DeviceID     *int64                 `json:"device_id,omitempty"`
}

// TaskExecuteRequest 执行任务请求
type TaskExecuteRequest struct {
	ExecutorType *string                `json:"executor_type,omitempty"`
	ExecutorID   *string                `json:"executor_id,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	DeviceID     *int64                 `json:"device_id,omitempty"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks      []Task          `json:"tasks"`
	Pagination PaginationInfo  `json:"pagination"`
	Statistics *TaskStatistics `json:"statistics,omitempty"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ===== 任务分配队列相关模型 =====

// TaskAssignmentQueue 任务分配队列模型
type TaskAssignmentQueue struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      int64  `gorm:"not null;uniqueIndex" json:"task_id"`
	Priority    int    `gorm:"not null;default:5" json:"priority"`
	QueueStatus string `gorm:"type:varchar(20);not null;default:queued" json:"queue_status"`

	// 分配条件
	RequiredDeviceType   *string `gorm:"type:varchar(50)" json:"required_device_type,omitempty"`
	RequiredCapabilities *JSONB  `gorm:"type:jsonb" json:"required_capabilities,omitempty"`
	PreferredDeviceIds   *string `gorm:"type:text" json:"preferred_device_ids,omitempty"` // 以逗号分隔的ID列表
	ExcludedDeviceIds    *string `gorm:"type:text" json:"excluded_device_ids,omitempty"`  // 以逗号分隔的ID列表

	// 分配结果
	AssignedDeviceID  *int64     `gorm:"index" json:"assigned_device_id,omitempty"`
	AssignedDeviceESN *string    `gorm:"type:varchar(100)" json:"assigned_device_esn,omitempty"`
	AssignedAt        *time.Time `gorm:"type:timestamp" json:"assigned_at,omitempty"`
	AssignmentScore   *float64   `gorm:"type:numeric(5,2)" json:"assignment_score,omitempty"`

	// 队列信息
	QueuePosition     *int `json:"queue_position,omitempty"`
	EstimatedWaitTime *int `json:"estimated_wait_time,omitempty"`
	RetryCount        int  `gorm:"default:0" json:"retry_count"`
	MaxRetries        int  `gorm:"default:3" json:"max_retries"`

	// 时间戳
	QueuedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"queued_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Task           *Task   `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`
	AssignedDevice *Device `gorm:"foreignKey:AssignedDeviceID;references:ID" json:"assigned_device,omitempty"`
}

// TableName 指定表名
func (TaskAssignmentQueue) TableName() string {
	return "task_assignment_queue"
}

// DeviceLoadMonitor 设备负载监控模型
type DeviceLoadMonitor struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  int64  `gorm:"not null;index" json:"device_id"`
	DeviceESN string `gorm:"type:varchar(100);not null" json:"device_esn"`

	// 负载指标
	CurrentTasks       int      `gorm:"default:0" json:"current_tasks"`
	MaxConcurrentTasks int      `gorm:"default:1" json:"max_concurrent_tasks"`
	CPULoad            *float64 `gorm:"type:numeric(5,2)" json:"cpu_load,omitempty"`
	MemoryUsage        *float64 `gorm:"type:numeric(5,2)" json:"memory_usage,omitempty"`
	DiskUsage          *float64 `gorm:"type:numeric(5,2)" json:"disk_usage,omitempty"`
	NetworkLatency     *int     `json:"network_latency,omitempty"`

	// 设备状态
	Status        string     `gorm:"type:varchar(20);not null;default:online" json:"status"`
	LastHeartbeat *time.Time `gorm:"type:timestamp" json:"last_heartbeat,omitempty"`
	LoadScore     *float64   `gorm:"type:numeric(5,2)" json:"load_score,omitempty"`

	// 统计信息
	TotalAssigned  int      `gorm:"default:0" json:"total_assigned"`
	TotalCompleted int      `gorm:"default:0" json:"total_completed"`
	TotalFailed    int      `gorm:"default:0" json:"total_failed"`
	SuccessRate    *float64 `gorm:"type:numeric(5,2)" json:"success_rate,omitempty"`

	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (DeviceLoadMonitor) TableName() string {
	return "device_load_monitor"
}

// TaskAssignmentHistory 任务分配历史记录模型
type TaskAssignmentHistory struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    int64   `gorm:"not null;index" json:"task_id"`
	DeviceID  *int64  `gorm:"index" json:"device_id,omitempty"`
	DeviceESN *string `gorm:"type:varchar(100)" json:"device_esn,omitempty"`

	Action         string  `gorm:"type:varchar(20);not null" json:"action"`
	PreviousStatus *string `gorm:"type:varchar(20)" json:"previous_status,omitempty"`
	NewStatus      *string `gorm:"type:varchar(20)" json:"new_status,omitempty"`
	Reason         *string `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Details        *JSONB  `gorm:"type:jsonb" json:"details,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`

	// 关联关系
	Task   *Task   `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (TaskAssignmentHistory) TableName() string {
	return "task_assignment_history"
}

// 队列状态常量 - 使用统一常量
const (
	QueueStatusQueued    = consts.QueueStatusQueued
	QueueStatusAssigning = consts.QueueStatusAssigning
	QueueStatusAssigned  = consts.QueueStatusAssigned
	QueueStatusFailed    = consts.QueueStatusFailed
	QueueStatusCanceled  = consts.QueueStatusCanceled
)

// 设备负载状态常量 - 使用统一常量
const (
	DeviceStatusOnline      = consts.DeviceStatusOnline
	DeviceStatusOffline     = consts.DeviceStatusOffline
	DeviceStatusBusy        = "busy" // 保持原有定义
	DeviceStatusMaintenance = consts.DeviceStatusMaintenance
)

// 分配历史动作常量 - 使用统一常量
const (
	AssignmentActionQueued     = consts.AssignmentActionQueued
	AssignmentActionAssigned   = consts.AssignmentActionAssigned
	AssignmentActionReassigned = consts.AssignmentActionReassigned
	AssignmentActionFailed     = consts.AssignmentActionFailed
	AssignmentActionCompleted  = consts.AssignmentActionCompleted
	AssignmentActionCanceled   = consts.AssignmentActionCanceled
)

// TaskQueueFilter 任务队列查询过滤器
type TaskQueueFilter struct {
	QueueStatus    []string   `json:"queue_status,omitempty"`
	Priority       *int       `json:"priority,omitempty"`
	DeviceType     *string    `json:"device_type,omitempty"`
	AssignedDevice *int64     `json:"assigned_device,omitempty"`
	QueuedAfter    *time.Time `json:"queued_after,omitempty"`
	QueuedBefore   *time.Time `json:"queued_before,omitempty"`
}

// DeviceLoadFilter 设备负载查询过滤器
type DeviceLoadFilter struct {
	Status         []string `json:"status,omitempty"`
	MinLoadScore   *float64 `json:"min_load_score,omitempty"`
	MaxLoadScore   *float64 `json:"max_load_score,omitempty"`
	MinSuccessRate *float64 `json:"min_success_rate,omitempty"`
	DeviceType     *string  `json:"device_type,omitempty"`
}

// QueueStatistics 队列统计信息
type QueueStatistics struct {
	Total         int64            `json:"total"`
	ByStatus      map[string]int64 `json:"by_status"`
	ByPriority    map[string]int64 `json:"by_priority"`
	AvgWaitTime   float64          `json:"avg_wait_time"`
	TotalAssigned int64            `json:"total_assigned"`
	TotalFailed   int64            `json:"total_failed"`
}

// DeviceLoadSummary 设备负载汇总
type DeviceLoadSummary struct {
	TotalDevices   int64   `json:"total_devices"`
	OnlineDevices  int64   `json:"online_devices"`
	BusyDevices    int64   `json:"busy_devices"`
	OfflineDevices int64   `json:"offline_devices"`
	AvgLoadScore   float64 `json:"avg_load_score"`
	AvgSuccessRate float64 `json:"avg_success_rate"`
}
