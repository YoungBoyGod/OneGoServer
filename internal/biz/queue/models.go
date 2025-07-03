package queue

import (
	"time"
)

// 注意：JSONB、Device、Task等类型在其他包中已定义，这里创建局部引用

// ===== 设备队列相关模型 =====

// DeviceTaskQueue 设备任务队列模型
type DeviceTaskQueue struct {
	// 主键
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  int64  `gorm:"not null;index" json:"device_id"`
	DeviceESN string `gorm:"type:varchar(100);not null" json:"device_esn"`
	TaskID    int64  `gorm:"not null;index" json:"task_id"`

	// 队列管理
	QueuePriority    int  `gorm:"not null;default:5" json:"queue_priority"`
	QueuePosition    int  `gorm:"not null" json:"queue_position"`
	OriginalPriority *int `json:"original_priority,omitempty"`
	IsManualPriority bool `gorm:"default:false" json:"is_manual_priority"`
	IsManualPosition bool `gorm:"default:false" json:"is_manual_position"`

	// 状态信息
	Status             string     `gorm:"type:varchar(20);not null;default:queued" json:"status"`
	EstimatedStartTime *time.Time `gorm:"type:timestamp" json:"estimated_start_time,omitempty"`
	EstimatedDuration  *int       `json:"estimated_duration,omitempty"`
	ActualStartTime    *time.Time `gorm:"type:timestamp" json:"actual_start_time,omitempty"`
	ActualEndTime      *time.Time `gorm:"type:timestamp" json:"actual_end_time,omitempty"`

	// 执行配置
	MaxRetryCount  int `gorm:"default:3" json:"max_retry_count"`
	CurrentRetry   int `gorm:"default:0" json:"current_retry"`
	TimeoutSeconds int `gorm:"default:3600" json:"timeout_seconds"`

	// 依赖关系
	DependsOnTaskIds *string `gorm:"type:text" json:"depends_on_task_ids,omitempty"` // 逗号分隔的ID列表
	BlocksTaskIds    *string `gorm:"type:text" json:"blocks_task_ids,omitempty"`     // 逗号分隔的ID列表

	// 操作记录
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	QueuedBy       *int64    `json:"queued_by,omitempty"`
	LastModifiedBy *int64    `json:"last_modified_by,omitempty"`
	LastAction     *string   `gorm:"type:varchar(50)" json:"last_action,omitempty"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
	Task   *Task   `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`
}

// TableName 指定表名
func (DeviceTaskQueue) TableName() string {
	return "device_task_queue"
}

// DeviceQueueOperationHistory 设备队列操作历史模型
type DeviceQueueOperationHistory struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  int64  `gorm:"not null;index" json:"device_id"`
	DeviceESN string `gorm:"type:varchar(100);not null" json:"device_esn"`
	TaskID    *int64 `gorm:"index" json:"task_id,omitempty"`

	// 操作信息
	OperationType string    `gorm:"type:varchar(30);not null" json:"operation_type"`
	OperationBy   *int64    `json:"operation_by,omitempty"`
	OperationTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"operation_time"`

	// 变更详情
	OldPriority *int    `json:"old_priority,omitempty"`
	NewPriority *int    `json:"new_priority,omitempty"`
	OldPosition *int    `json:"old_position,omitempty"`
	NewPosition *int    `json:"new_position,omitempty"`
	OldStatus   *string `gorm:"type:varchar(20)" json:"old_status,omitempty"`
	NewStatus   *string `gorm:"type:varchar(20)" json:"new_status,omitempty"`

	// 操作原因和备注
	Reason          *string `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Notes           *string `gorm:"type:text" json:"notes,omitempty"`
	OperationSource string  `gorm:"type:varchar(20);default:manual" json:"operation_source"`

	// 批量操作支持
	BatchID          *string `gorm:"type:varchar(50)" json:"batch_id,omitempty"`
	IsBatchOperation bool    `gorm:"default:false" json:"is_batch_operation"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
	Task   *Task   `gorm:"foreignKey:TaskID;references:ID" json:"task,omitempty"`
}

// TableName 指定表名
func (DeviceQueueOperationHistory) TableName() string {
	return "device_queue_operation_history"
}

// DeviceQueueConfig 设备队列配置模型
type DeviceQueueConfig struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  int64  `gorm:"not null;uniqueIndex" json:"device_id"`
	DeviceESN string `gorm:"type:varchar(100);not null" json:"device_esn"`

	// 队列配置
	MaxQueueSize       int  `gorm:"default:100" json:"max_queue_size"`
	MaxConcurrentTasks int  `gorm:"default:1" json:"max_concurrent_tasks"`
	AutoStartTasks     bool `gorm:"default:true" json:"auto_start_tasks"`
	PriorityScheduling bool `gorm:"default:true" json:"priority_scheduling"`

	// 调度策略
	SchedulingStrategy string `gorm:"type:varchar(20);default:priority_first" json:"scheduling_strategy"`
	LoadBalancing      bool   `gorm:"default:true" json:"load_balancing"`

	// 时间窗口配置
	WorkStartTime *time.Time `gorm:"type:time" json:"work_start_time,omitempty"`
	WorkEndTime   *time.Time `gorm:"type:time" json:"work_end_time,omitempty"`
	Timezone      string     `gorm:"type:varchar(50);default:UTC" json:"timezone"`

	// 资源限制
	MaxCPUUsage    float64 `gorm:"type:numeric(5,2);default:80.0" json:"max_cpu_usage"`
	MaxMemoryUsage float64 `gorm:"type:numeric(5,2);default:80.0" json:"max_memory_usage"`
	MinFreeDisk    int64   `gorm:"default:1073741824" json:"min_free_disk"`

	// 通知配置
	NotifyOnCompletion  bool    `gorm:"default:false" json:"notify_on_completion"`
	NotifyOnFailure     bool    `gorm:"default:true" json:"notify_on_failure"`
	NotificationWebhook *string `gorm:"type:varchar(255)" json:"notification_webhook,omitempty"`

	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (DeviceQueueConfig) TableName() string {
	return "device_queue_config"
}

// ===== 状态常量 =====

// 设备队列任务状态常量
const (
	DeviceQueueStatusQueued    = "queued"
	DeviceQueueStatusExecuting = "executing"
	DeviceQueueStatusPaused    = "paused"
	DeviceQueueStatusCompleted = "completed"
	DeviceQueueStatusFailed    = "failed"
	DeviceQueueStatusCanceled  = "canceled"
)

// 队列操作类型常量
const (
	QueueOperationAdd            = "add"
	QueueOperationRemove         = "remove"
	QueueOperationPriorityChange = "priority_change"
	QueueOperationPositionChange = "position_change"
	QueueOperationStart          = "start"
	QueueOperationPause          = "pause"
	QueueOperationResume         = "resume"
	QueueOperationCancel         = "cancel"
)

// 调度策略常量
const (
	SchedulingStrategyPriorityFirst = "priority_first"
	SchedulingStrategyFIFO          = "fifo"
	SchedulingStrategyLIFO          = "lifo"
	SchedulingStrategyWeighted      = "weighted"
)

// 操作来源常量
const (
	OperationSourceManual    = "manual"
	OperationSourceSystem    = "system"
	OperationSourceAPI       = "api"
	OperationSourceScheduler = "scheduler"
)

// ===== 数据传输对象 =====

// DeviceQueueFilter 设备队列查询过滤器
type DeviceQueueFilter struct {
	DeviceID         *int64     `json:"device_id,omitempty"`
	DeviceESN        *string    `json:"device_esn,omitempty"`
	TaskID           *int64     `json:"task_id,omitempty"`
	Status           []string   `json:"status,omitempty"`
	Priority         *int       `json:"priority,omitempty"`
	CreatedAfter     *time.Time `json:"created_after,omitempty"`
	CreatedBefore    *time.Time `json:"created_before,omitempty"`
	IsManualPriority *bool      `json:"is_manual_priority,omitempty"`
	IsManualPosition *bool      `json:"is_manual_position,omitempty"`
}

// DeviceQueueSortOption 设备队列排序选项
type DeviceQueueSortOption struct {
	Field string `json:"field"` // queue_position, queue_priority, created_at, estimated_start_time
	Order string `json:"order"` // asc, desc
}

// DeviceQueueStatistics 设备队列统计信息
type DeviceQueueStatistics struct {
	DeviceID               int64            `json:"device_id"`
	DeviceESN              string           `json:"device_esn"`
	TotalTasks             int64            `json:"total_tasks"`
	QueuedTasks            int64            `json:"queued_tasks"`
	ExecutingTasks         int64            `json:"executing_tasks"`
	CompletedTasks         int64            `json:"completed_tasks"`
	FailedTasks            int64            `json:"failed_tasks"`
	ManualPriorityTasks    int64            `json:"manual_priority_tasks"`
	ManualPositionTasks    int64            `json:"manual_position_tasks"`
	ByStatus               map[string]int64 `json:"by_status"`
	ByPriority             map[string]int64 `json:"by_priority"`
	AvgWaitTime            float64          `json:"avg_wait_time"`
	TotalEstimatedDuration int64            `json:"total_estimated_duration"`
	NextTaskStartTime      *time.Time       `json:"next_task_start_time,omitempty"`
}

// BatchQueueOperation 批量队列操作请求
type BatchQueueOperation struct {
	DeviceID     int64   `json:"device_id" binding:"required"`
	TaskIDs      []int64 `json:"task_ids" binding:"required"`
	Operation    string  `json:"operation" binding:"required"` // priority_change, position_change, cancel, pause, resume
	NewPriority  *int    `json:"new_priority,omitempty"`
	NewPositions []int   `json:"new_positions,omitempty"`
	Reason       *string `json:"reason,omitempty"`
	BatchID      *string `json:"batch_id,omitempty"`
}

// QueuePositionChangeRequest 队列位置调整请求
type QueuePositionChangeRequest struct {
	DeviceID    int64  `json:"device_id" binding:"required"`
	TaskID      int64  `json:"task_id" binding:"required"`
	NewPosition int    `json:"new_position" binding:"required"`
	Reason      string `json:"reason"`
}

// QueuePriorityChangeRequest 队列优先级调整请求
type QueuePriorityChangeRequest struct {
	DeviceID    int64  `json:"device_id" binding:"required"`
	TaskID      int64  `json:"task_id" binding:"required"`
	NewPriority int    `json:"new_priority" binding:"required,min=1,max=10"`
	Reason      string `json:"reason"`
}

// DeviceQueueResponse 设备队列响应
type DeviceQueueResponse struct {
	DeviceInfo  DeviceInfo            `json:"device_info"`
	QueueConfig *DeviceQueueConfig    `json:"queue_config,omitempty"`
	QueueItems  []DeviceTaskQueue     `json:"queue_items"`
	Statistics  DeviceQueueStatistics `json:"statistics"`
	Pagination  PaginationInfo        `json:"pagination"`
}

// DeviceInfo 设备基本信息
type DeviceInfo struct {
	ID        int64  `json:"id"`
	DeviceESN string `json:"device_esn"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

// PaginationInfo 分页信息 (引用公共类型)

// DeviceQueueListResponse 设备队列列表响应
type DeviceQueueListResponse struct {
	Devices []DeviceQueueSummary `json:"devices"`
	Summary GlobalQueueSummary   `json:"summary"`
}

// DeviceQueueSummary 设备队列摘要
type DeviceQueueSummary struct {
	DeviceInfo        DeviceInfo `json:"device_info"`
	TotalQueuedTasks  int64      `json:"total_queued_tasks"`
	ExecutingTasks    int64      `json:"executing_tasks"`
	ManualAdjustments int64      `json:"manual_adjustments"`
	NextTaskStartTime *time.Time `json:"next_task_start_time,omitempty"`
	EstimatedDuration int64      `json:"estimated_duration"`
	LastQueueUpdate   time.Time  `json:"last_queue_update"`
}

// GlobalQueueSummary 全局队列摘要
type GlobalQueueSummary struct {
	TotalDevices           int64   `json:"total_devices"`
	DevicesWithTasks       int64   `json:"devices_with_tasks"`
	TotalQueuedTasks       int64   `json:"total_queued_tasks"`
	TotalExecutingTasks    int64   `json:"total_executing_tasks"`
	TotalManualAdjustments int64   `json:"total_manual_adjustments"`
	AvgQueueLength         float64 `json:"avg_queue_length"`
}

// ===== 引用的外部模型 =====
// 注意：Device和Task类型在其他包中已定义，在实际使用时需要导入相应的包
