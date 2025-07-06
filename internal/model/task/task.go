package task

import (
	"time"

	"OneGfServer/internal/consts"

	"github.com/gogf/gf/v2/os/gtime"
)

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

// 队列状态常量 - 使用统一常量
const (
	QueueStatusQueued    = consts.QueueStatusQueued
	QueueStatusAssigning = consts.QueueStatusAssigning
	QueueStatusAssigned  = consts.QueueStatusAssigned
	QueueStatusFailed    = consts.QueueStatusFailed
	QueueStatusCanceled  = consts.QueueStatusCanceled
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

// 任务依赖类型常量 - 使用统一常量
const (
	DependencyTypeBefore   = consts.DependencyTypeBefore
	DependencyTypeAfter    = consts.DependencyTypeAfter
	DependencyTypeParallel = consts.DependencyTypeParallel
)

// 任务优先级常量 - 使用统一常量
const (
	TaskPriorityLowest  = 1
	TaskPriorityLow     = 3
	TaskPriorityNormal  = 5
	TaskPriorityHigh    = 7
	TaskPriorityUrgent  = 9
	TaskPriorityHighest = 10
)

// ===============================
// 基础CRUD操作 Input/Output
// ===============================

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
	Task *Task `json:"task"`
}

// CreateTaskOutput 创建任务输出
type CreateTaskOutput struct {
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
}

// GetTaskByIDInput 根据ID获取任务输入
type GetTaskByIDInput struct {
	ID int64 `json:"id"`
}

// GetTaskByIDOutput 根据ID获取任务输出
type GetTaskByIDOutput struct {
	Task *Task `json:"task"`
}

// GetTaskByTaskIDInput 根据任务ID获取任务输入
type GetTaskByTaskIDInput struct {
	TaskID string `json:"task_id"`
}

// GetTaskByTaskIDOutput 根据任务ID获取任务输出
type GetTaskByTaskIDOutput struct {
	Task *Task `json:"task"`
}

// UpdateTaskInput 更新任务输入
type UpdateTaskInput struct {
	Task *Task `json:"task"`
}

// UpdateTaskOutput 更新任务输出
type UpdateTaskOutput struct {
	Message string `json:"message"`
}

// DeleteTaskInput 删除任务输入
type DeleteTaskInput struct {
	TaskID string `json:"task_id"`
	Force  bool   `json:"force"`
}

// DeleteTaskOutput 删除任务输出
type DeleteTaskOutput struct {
	Message string `json:"message"`
}

// ===============================
// 查询操作 Input/Output
// ===============================

// GetTaskListInput 获取任务列表输入
type GetTaskListInput struct {
	Filter     *TaskFilter       `json:"filter"`
	Sort       *TaskSortOption   `json:"sort"`
	Pagination *PaginationOption `json:"pagination"`
}

// GetTaskListOutput 获取任务列表输出
type GetTaskListOutput struct {
	List  []Task `json:"list"`
	Total int64  `json:"total"`
}

// GetTasksByTypeInput 根据任务类型获取任务输入
type GetTasksByTypeInput struct {
	TaskTypes []string `json:"task_types"`
}

// GetTasksByTypeOutput 根据任务类型获取任务输出
type GetTasksByTypeOutput struct {
	Tasks []Task `json:"tasks"`
}

// GetTasksByStatusInput 根据状态获取任务输入
type GetTasksByStatusInput struct {
	Statuses []string `json:"statuses"`
}

// GetTasksByStatusOutput 根据状态获取任务输出
type GetTasksByStatusOutput struct {
	Tasks []Task `json:"tasks"`
}

// GetTasksByDeviceInput 根据设备获取任务输入
type GetTasksByDeviceInput struct {
	DeviceID string `json:"device_id"`
}

// GetTasksByDeviceOutput 根据设备获取任务输出
type GetTasksByDeviceOutput struct {
	Tasks []Task `json:"tasks"`
}

// ===============================
// 任务执行控制 Input/Output
// ===============================

// StartTaskInput 启动任务输入
type StartTaskInput struct {
	TaskID string `json:"task_id"`
	Force  bool   `json:"force"`
}

// StartTaskOutput 启动任务输出
type StartTaskOutput struct {
	Message string `json:"message"`
}

// StopTaskInput 停止任务输入
type StopTaskInput struct {
	TaskID string `json:"task_id"`
	Reason string `json:"reason"`
}

// StopTaskOutput 停止任务输出
type StopTaskOutput struct {
	Message string `json:"message"`
}

// RestartTaskInput 重启任务输入
type RestartTaskInput struct {
	TaskID string `json:"task_id"`
	Force  bool   `json:"force"`
}

// RestartTaskOutput 重启任务输出
type RestartTaskOutput struct {
	Message string `json:"message"`
}

// CancelTaskInput 取消任务输入
type CancelTaskInput struct {
	TaskID string `json:"task_id"`
	Reason string `json:"reason"`
}

// CancelTaskOutput 取消任务输出
type CancelTaskOutput struct {
	Message string `json:"message"`
}

// RetryTaskInput 重试任务输入
type RetryTaskInput struct {
	TaskID string `json:"task_id"`
	Force  bool   `json:"force"`
}

// RetryTaskOutput 重试任务输出
type RetryTaskOutput struct {
	Message string `json:"message"`
}

// ===============================
// 任务优先级管理 Input/Output
// ===============================

// GetTaskPrioritiesInput 获取任务优先级列表输入
type GetTaskPrioritiesInput struct{}

// GetTaskPrioritiesOutput 获取任务优先级列表输出
type GetTaskPrioritiesOutput struct {
	Priorities []TaskPriorityInfo `json:"priorities"`
}

// UpdateTaskPriorityInput 更新任务优先级输入
type UpdateTaskPriorityInput struct {
	TaskID   string `json:"task_id"`
	Priority int    `json:"priority"`
	Reason   string `json:"reason"`
}

// UpdateTaskPriorityOutput 更新任务优先级输出
type UpdateTaskPriorityOutput struct {
	Message string `json:"message"`
}

// BatchUpdateTaskPriorityInput 批量更新任务优先级输入
type BatchUpdateTaskPriorityInput struct {
	TaskIDs  []string `json:"task_ids"`
	Priority int      `json:"priority"`
	Reason   string   `json:"reason"`
}

// BatchUpdateTaskPriorityOutput 批量更新任务优先级输出
type BatchUpdateTaskPriorityOutput struct {
	AffectedCount int    `json:"affected_count"`
	Message       string `json:"message"`
}

// ===============================
// 任务状态管理 Input/Output
// ===============================

// GetTaskStatusInput 获取任务状态输入
type GetTaskStatusInput struct {
	TaskID string `json:"task_id"`
}

// GetTaskStatusOutput 获取任务状态输出
type GetTaskStatusOutput struct {
	Status *TaskStatusInfo `json:"status"`
}

// UpdateTaskStatusInput 更新任务状态输入
type UpdateTaskStatusInput struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// UpdateTaskStatusOutput 更新任务状态输出
type UpdateTaskStatusOutput struct {
	Message string `json:"message"`
}

// ===============================
// 任务日志管理 Input/Output
// ===============================

// GetTaskLogsInput 获取任务日志列表输入
type GetTaskLogsInput struct {
	TaskID string     `json:"task_id"`
	Filter *LogFilter `json:"filter"`
}

// GetTaskLogsOutput 获取任务日志列表输出
type GetTaskLogsOutput struct {
	Logs []TaskLog `json:"logs"`
}

// GetTaskLogDetailInput 获取任务日志详情输入
type GetTaskLogDetailInput struct {
	LogID string `json:"log_id"`
}

// GetTaskLogDetailOutput 获取任务日志详情输出
type GetTaskLogDetailOutput struct {
	Log *TaskLogDetail `json:"log"`
}

// ClearTaskLogsInput 清空任务日志输入
type ClearTaskLogsInput struct {
	TaskID string     `json:"task_id"`
	Before *time.Time `json:"before"`
}

// ClearTaskLogsOutput 清空任务日志输出
type ClearTaskLogsOutput struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ExportTaskLogsInput 导出任务日志输入
type ExportTaskLogsInput struct {
	TaskID    string     `json:"task_id"`
	Format    string     `json:"format"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// ExportTaskLogsOutput 导出任务日志输出
type ExportTaskLogsOutput struct {
	DownloadURL string `json:"download_url"`
	Message     string `json:"message"`
}

// ===============================
// 任务队列管理 Input/Output
// ===============================

// GetTaskQueuesInput 获取任务队列列表输入
type GetTaskQueuesInput struct {
	Filter *QueueFilter `json:"filter"`
}

// GetTaskQueuesOutput 获取任务队列列表输出
type GetTaskQueuesOutput struct {
	Queues []TaskQueue `json:"queues"`
}

// GetTaskQueueDetailInput 获取任务队列详情输入
type GetTaskQueueDetailInput struct {
	QueueID string `json:"queue_id"`
}

// GetTaskQueueDetailOutput 获取任务队列详情输出
type GetTaskQueueDetailOutput struct {
	Queue *TaskQueueDetail `json:"queue"`
}

// ManageTaskQueueInput 管理任务队列输入
type ManageTaskQueueInput struct {
	QueueID string                 `json:"queue_id"`
	Action  string                 `json:"action"`
	Config  map[string]interface{} `json:"config"`
}

// ManageTaskQueueOutput 管理任务队列输出
type ManageTaskQueueOutput struct {
	Message string `json:"message"`
}

// GetTaskQueueStatsInput 获取任务队列统计输入
type GetTaskQueueStatsInput struct{}

// GetTaskQueueStatsOutput 获取任务队列统计输出
type GetTaskQueueStatsOutput struct {
	Stats *TaskQueueStats `json:"stats"`
}

// ===============================
// 任务调度管理 Input/Output
// ===============================

// GetTaskSchedulesInput 获取任务调度列表输入
type GetTaskSchedulesInput struct {
	Filter *ScheduleFilter `json:"filter"`
}

// GetTaskSchedulesOutput 获取任务调度列表输出
type GetTaskSchedulesOutput struct {
	Schedules []TaskSchedule `json:"schedules"`
}

// CreateTaskScheduleInput 创建任务调度输入
type CreateTaskScheduleInput struct {
	Schedule *TaskSchedule `json:"schedule"`
}

// CreateTaskScheduleOutput 创建任务调度输出
type CreateTaskScheduleOutput struct {
	ScheduleID string `json:"schedule_id"`
	Message    string `json:"message"`
}

// UpdateTaskScheduleInput 更新任务调度输入
type UpdateTaskScheduleInput struct {
	Schedule *TaskSchedule `json:"schedule"`
}

// UpdateTaskScheduleOutput 更新任务调度输出
type UpdateTaskScheduleOutput struct {
	Message string `json:"message"`
}

// DeleteTaskScheduleInput 删除任务调度输入
type DeleteTaskScheduleInput struct {
	ScheduleID string `json:"schedule_id"`
}

// DeleteTaskScheduleOutput 删除任务调度输出
type DeleteTaskScheduleOutput struct {
	Message string `json:"message"`
}

// ===============================
// 任务统计分析 Input/Output
// ===============================

// GetTaskStatisticsInput 获取任务统计信息输入
type GetTaskStatisticsInput struct {
	Filter *TaskFilter `json:"filter"`
}

// GetTaskStatisticsOutput 获取任务统计信息输出
type GetTaskStatisticsOutput struct {
	Statistics *TaskStatistics `json:"statistics"`
}

// GetTaskPerformanceReportInput 获取任务性能报告输入
type GetTaskPerformanceReportInput struct {
	TaskID string `json:"task_id"`
}

// GetTaskPerformanceReportOutput 获取任务性能报告输出
type GetTaskPerformanceReportOutput struct {
	Report *TaskPerformanceReport `json:"report"`
}

// ===============================
// 任务执行管理 Input/Output
// ===============================

// CreateTaskExecutionInput 创建任务执行输入
type CreateTaskExecutionInput struct {
	Execution *TaskExecution `json:"execution"`
}

// CreateTaskExecutionOutput 创建任务执行输出
type CreateTaskExecutionOutput struct {
	ExecutionID string `json:"execution_id"`
	Message     string `json:"message"`
}

// GetTaskExecutionInput 获取任务执行输入
type GetTaskExecutionInput struct {
	ExecutionID string `json:"execution_id"`
}

// GetTaskExecutionOutput 获取任务执行输出
type GetTaskExecutionOutput struct {
	Execution *TaskExecution `json:"execution"`
}

// UpdateTaskExecutionInput 更新任务执行输入
type UpdateTaskExecutionInput struct {
	Execution *TaskExecution `json:"execution"`
}

// UpdateTaskExecutionOutput 更新任务执行输出
type UpdateTaskExecutionOutput struct {
	Message string `json:"message"`
}

// GetTaskExecutionsInput 获取任务执行历史输入
type GetTaskExecutionsInput struct {
	TaskID string           `json:"task_id"`
	Filter *ExecutionFilter `json:"filter"`
}

// GetTaskExecutionsOutput 获取任务执行历史输出
type GetTaskExecutionsOutput struct {
	Executions []TaskExecution `json:"executions"`
}

// ===============================
// 任务分配管理 Input/Output
// ===============================

// CreateTaskAssignmentInput 创建任务分配输入
type CreateTaskAssignmentInput struct {
	Assignment *TaskAssignment `json:"assignment"`
}

// CreateTaskAssignmentOutput 创建任务分配输出
type CreateTaskAssignmentOutput struct {
	AssignmentID string `json:"assignment_id"`
	Message      string `json:"message"`
}

// GetTaskAssignmentInput 获取任务分配输入
type GetTaskAssignmentInput struct {
	AssignmentID string `json:"assignment_id"`
}

// GetTaskAssignmentOutput 获取任务分配输出
type GetTaskAssignmentOutput struct {
	Assignment *TaskAssignment `json:"assignment"`
}

// UpdateTaskAssignmentInput 更新任务分配输入
type UpdateTaskAssignmentInput struct {
	Assignment *TaskAssignment `json:"assignment"`
}

// UpdateTaskAssignmentOutput 更新任务分配输出
type UpdateTaskAssignmentOutput struct {
	Message string `json:"message"`
}

// GetTaskAssignmentsInput 获取任务分配历史输入
type GetTaskAssignmentsInput struct {
	TaskID string            `json:"task_id"`
	Filter *AssignmentFilter `json:"filter"`
}

// GetTaskAssignmentsOutput 获取任务分配历史输出
type GetTaskAssignmentsOutput struct {
	Assignments []TaskAssignment `json:"assignments"`
}

// ===============================
// 数据模型定义
// ===============================

// Task 任务主表模型
type Task struct {
	ID           int64       `json:"id"`
	TaskID       string      `json:"task_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Status       string      `json:"status"`
	Priority     int         `json:"priority"`
	ExecuteTime  *gtime.Time `json:"execute_time"`
	Timeout      int         `json:"timeout"`
	RetryCount   int         `json:"retry_count"`
	MaxRetries   int         `json:"max_retries"`
	IsUrgent     bool        `json:"is_urgent"`
	Parameters   string      `json:"parameters"`
	Result       string      `json:"result"`
	ErrorMessage string      `json:"error_message"`
	ExecutorType string      `json:"executor_type"`
	ExecutorID   string      `json:"executor_id"`
	DeviceID     string      `json:"device_id"`
	CreatedAt    *gtime.Time `json:"created_at"`
	UpdatedAt    *gtime.Time `json:"updated_at"`
	CreatedBy    int64       `json:"created_by"`
	UpdatedBy    int64       `json:"updated_by"`
}

// TaskExecution 任务执行模型
type TaskExecution struct {
	ID                int64       `json:"id"`
	TaskID            string      `json:"task_id"`
	ExecutionID       string      `json:"execution_id"`
	DeviceESN         string      `json:"device_esn"`
	Status            string      `json:"status"`
	StartTime         *gtime.Time `json:"start_time"`
	EndTime           *gtime.Time `json:"end_time"`
	Duration          int         `json:"duration"`
	ExecutorInfo      string      `json:"executor_info"`
	Logs              string      `json:"logs"`
	Metrics           string      `json:"metrics"`
	Output            string      `json:"output"`
	ErrorDetails      string      `json:"error_details"`
	CpuUsageAvg       float64     `json:"cpu_usage_avg"`
	CpuUsagePeak      float64     `json:"cpu_usage_peak"`
	MemoryUsageAvg    float64     `json:"memory_usage_avg"`
	MemoryUsagePeak   float64     `json:"memory_usage_peak"`
	IoOperationsTotal int64       `json:"io_operations_total"`
	IoBytesTotal      int64       `json:"io_bytes_total"`
	CreatedAt         *gtime.Time `json:"created_at"`
	UpdatedAt         *gtime.Time `json:"updated_at"`
}

// TaskAssignment 任务分配模型
type TaskAssignment struct {
	ID                   int64       `json:"id"`
	TaskID               string      `json:"task_id"`
	Priority             int         `json:"priority"`
	QueueStatus          string      `json:"queue_status"`
	RequiredDeviceType   string      `json:"required_device_type"`
	RequiredCapabilities string      `json:"required_capabilities"`
	PreferredDeviceIDs   []int64     `json:"preferred_device_ids"`
	ExcludedDeviceIDs    []int64     `json:"excluded_device_ids"`
	AssignmentStrategy   string      `json:"assignment_strategy"`
	AffinityRules        string      `json:"affinity_rules"`
	AssignedDeviceID     int64       `json:"assigned_device_id"`
	AssignedAt           *gtime.Time `json:"assigned_at"`
	AssignmentScore      float64     `json:"assignment_score"`
	QueuePosition        int         `json:"queue_position"`
	EstimatedWaitTime    int         `json:"estimated_wait_time"`
	RetryCount           int         `json:"retry_count"`
	MaxRetries           int         `json:"max_retries"`
	OriginalPriority     int         `json:"original_priority"`
	LastPriorityChangeAt *gtime.Time `json:"last_priority_change_at"`
	PriorityChangeReason string      `json:"priority_change_reason"`
	PriorityBoostReason  string      `json:"priority_boost_reason"`
	OperationSource      string      `json:"operation_source"`
	QueuedAt             *gtime.Time `json:"queued_at"`
	UpdatedAt            *gtime.Time `json:"updated_at"`
}

// TaskStatusInfo 任务状态信息
type TaskStatusInfo struct {
	TaskID       string      `json:"task_id"`
	Status       string      `json:"status"`
	Progress     float64     `json:"progress"`
	StartTime    *gtime.Time `json:"start_time"`
	EndTime      *gtime.Time `json:"end_time"`
	Duration     int64       `json:"duration"`
	RetryCount   int         `json:"retry_count"`
	ErrorMessage string      `json:"error_message"`
}

// TaskPriorityInfo 任务优先级信息
type TaskPriorityInfo struct {
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TaskCount   int    `json:"task_count"`
	Color       string `json:"color"`
}

// TaskLog 任务日志模型
type TaskLog struct {
	LogID       string      `json:"log_id"`
	TaskID      string      `json:"task_id"`
	ExecutionID string      `json:"execution_id"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	Timestamp   *gtime.Time `json:"timestamp"`
	Source      string      `json:"source"`
}

// TaskLogDetail 任务日志详情模型
type TaskLogDetail struct {
	TaskLog
	StackTrace string                 `json:"stack_trace"`
	Context    map[string]interface{} `json:"context"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// TaskQueue 任务队列模型
type TaskQueue struct {
	QueueID   string      `json:"queue_id"`
	QueueName string      `json:"queue_name"`
	QueueType string      `json:"queue_type"`
	DeviceID  string      `json:"device_id"`
	TaskCount int         `json:"task_count"`
	Status    string      `json:"status"`
	Priority  int         `json:"priority"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

// TaskQueueDetail 任务队列详情模型
type TaskQueueDetail struct {
	TaskQueue
	Tasks           []Task      `json:"tasks"`
	MaxConcurrency  int         `json:"max_concurrency"`
	CurrentRunning  int         `json:"current_running"`
	ProcessingSpeed float64     `json:"processing_speed"`
	AverageWaitTime float64     `json:"average_wait_time"`
	LastProcessTime *gtime.Time `json:"last_process_time"`
}

// TaskQueueStats 任务队列统计信息
type TaskQueueStats struct {
	TotalQueues     int     `json:"total_queues"`
	TotalTasks      int     `json:"total_tasks"`
	PendingTasks    int     `json:"pending_tasks"`
	RunningTasks    int     `json:"running_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	FailedTasks     int     `json:"failed_tasks"`
	AverageWaitTime float64 `json:"average_wait_time"`
	Throughput      float64 `json:"throughput"`
}

// TaskSchedule 任务调度模型
type TaskSchedule struct {
	ScheduleID   string      `json:"schedule_id"`
	ScheduleName string      `json:"schedule_name"`
	ScheduleType string      `json:"schedule_type"`
	CronExpr     string      `json:"cron_expr"`
	Status       string      `json:"status"`
	Enabled      bool        `json:"enabled"`
	TaskCount    int         `json:"task_count"`
	LastRun      *gtime.Time `json:"last_run"`
	NextRun      *gtime.Time `json:"next_run"`
	CreatedAt    *gtime.Time `json:"created_at"`
	UpdatedAt    *gtime.Time `json:"updated_at"`
}

// TaskStatistics 任务统计信息
type TaskStatistics struct {
	TotalTasks           int               `json:"total_tasks"`
	StatusDistribution   map[string]int    `json:"status_distribution"`
	TypeDistribution     map[string]int    `json:"type_distribution"`
	PriorityDistribution map[string]int    `json:"priority_distribution"`
	DeviceDistribution   map[string]int    `json:"device_distribution"`
	SuccessRate          float64           `json:"success_rate"`
	AverageExecutionTime float64           `json:"average_execution_time"`
	TrendData            []TaskTrendPoint  `json:"trend_data"`
	TopFailedTasks       []TaskFailureInfo `json:"top_failed_tasks"`
}

// TaskTrendPoint 任务趋势数据点
type TaskTrendPoint struct {
	Timestamp    *gtime.Time `json:"timestamp"`
	TaskCount    int         `json:"task_count"`
	SuccessCount int         `json:"success_count"`
	FailureCount int         `json:"failure_count"`
	SuccessRate  float64     `json:"success_rate"`
	AvgDuration  float64     `json:"avg_duration"`
}

// TaskFailureInfo 任务失败信息
type TaskFailureInfo struct {
	TaskID       string  `json:"task_id"`
	TaskName     string  `json:"task_name"`
	FailureCount int     `json:"failure_count"`
	FailureRate  float64 `json:"failure_rate"`
	LastError    string  `json:"last_error"`
}

// TaskPerformanceReport 任务性能报告
type TaskPerformanceReport struct {
	TaskID             string                  `json:"task_id"`
	TaskName           string                  `json:"task_name"`
	OverallScore       float64                 `json:"overall_score"`
	ExecutionMetrics   *TaskExecutionMetrics   `json:"execution_metrics"`
	ResourceMetrics    *TaskResourceMetrics    `json:"resource_metrics"`
	ReliabilityMetrics *TaskReliabilityMetrics `json:"reliability_metrics"`
	PerformanceTrend   []TaskTrendPoint        `json:"performance_trend"`
	Recommendations    []string                `json:"recommendations"`
}

// TaskExecutionMetrics 任务执行指标
type TaskExecutionMetrics struct {
	TotalExecutions      int     `json:"total_executions"`
	SuccessfulExecutions int     `json:"successful_executions"`
	FailedExecutions     int     `json:"failed_executions"`
	SuccessRate          float64 `json:"success_rate"`
	AverageExecutionTime float64 `json:"average_execution_time"`
	MinExecutionTime     float64 `json:"min_execution_time"`
	MaxExecutionTime     float64 `json:"max_execution_time"`
}

// TaskResourceMetrics 任务资源指标
type TaskResourceMetrics struct {
	AverageCpuUsage    float64 `json:"average_cpu_usage"`
	PeakCpuUsage       float64 `json:"peak_cpu_usage"`
	AverageMemoryUsage float64 `json:"average_memory_usage"`
	PeakMemoryUsage    float64 `json:"peak_memory_usage"`
	TotalDiskIO        float64 `json:"total_disk_io"`
	TotalNetworkIO     float64 `json:"total_network_io"`
}

// TaskReliabilityMetrics 任务可靠性指标
type TaskReliabilityMetrics struct {
	MTBF             float64 `json:"mtbf"` // Mean Time Between Failures
	MTTR             float64 `json:"mttr"` // Mean Time To Recovery
	AvailabilityRate float64 `json:"availability_rate"`
	ErrorRate        float64 `json:"error_rate"`
	RetrySuccessRate float64 `json:"retry_success_rate"`
}

// ===============================
// 过滤和排序选项
// ===============================

// TaskFilter 任务过滤条件
type TaskFilter struct {
	Status       []string   `json:"status"`
	Type         []string   `json:"type"`
	Priority     *int       `json:"priority"`
	DeviceID     *string    `json:"device_id"`
	ExecutorType *string    `json:"executor_type"`
	IsUrgent     *bool      `json:"is_urgent"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Keyword      *string    `json:"keyword"`
}

// TaskSortOption 任务排序选项
type TaskSortOption struct {
	Field string `json:"field"` // id, name, status, priority, created_at, updated_at
	Order string `json:"order"` // asc, desc
}

// LogFilter 日志过滤条件
type LogFilter struct {
	Level     *string    `json:"level"`
	Source    *string    `json:"source"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Keyword   *string    `json:"keyword"`
}

// QueueFilter 队列过滤条件
type QueueFilter struct {
	Status    []string   `json:"status"`
	Type      []string   `json:"type"`
	DeviceID  *string    `json:"device_id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// ScheduleFilter 调度过滤条件
type ScheduleFilter struct {
	Status    []string   `json:"status"`
	Type      []string   `json:"type"`
	Enabled   *bool      `json:"enabled"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// ExecutionFilter 执行过滤条件
type ExecutionFilter struct {
	Status    []string   `json:"status"`
	DeviceESN *string    `json:"device_esn"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// AssignmentFilter 分配过滤条件
type AssignmentFilter struct {
	Status             []string   `json:"status"`
	AssignmentStrategy *string    `json:"assignment_strategy"`
	AssignedDeviceID   *int64     `json:"assigned_device_id"`
	StartTime          *time.Time `json:"start_time"`
	EndTime            *time.Time `json:"end_time"`
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
