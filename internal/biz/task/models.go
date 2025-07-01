package task

import (
	"database/sql/driver"
	"encoding/json"
	"time"

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

	// 资源使用
	CPUUsage     *float64 `gorm:"type:numeric(5,2)" json:"cpu_usage,omitempty"`
	MemoryUsage  *float64 `gorm:"type:numeric(10,2)" json:"memory_usage,omitempty"`
	IOOperations *int64   `json:"io_operations,omitempty"`

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
		te.ExecutionID = generateExecutionID()
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

// 任务状态常量
const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
	TaskStatusCanceled  = "canceled"
)

// 任务类型常量
const (
	TaskTypeBackup  = "backup"
	TaskTypeSync    = "sync"
	TaskTypeMonitor = "monitor"
	TaskTypeCustom  = "custom"
)

// 执行器类型常量
const (
	ExecutorTypeLocal     = "local"
	ExecutorTypeRemote    = "remote"
	ExecutorTypeContainer = "container"
	ExecutorTypeLambda    = "lambda"
)

// 执行状态常量
const (
	ExecutionStatusStarted   = "started"
	ExecutionStatusRunning   = "running"
	ExecutionStatusCompleted = "completed"
	ExecutionStatusFailed    = "failed"
	ExecutionStatusCanceled  = "canceled"
)

// generateExecutionID 生成唯一的执行ID
func generateExecutionID() string {
	// 使用时间戳和随机数生成唯一ID
	return "exec_" + time.Now().Format("20060102_150405") + "_" + generateRandomString(6)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

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
