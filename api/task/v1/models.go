package v1

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务数据结构定义
// ===============================

// TaskInfo 任务基本信息
type TaskInfo struct {
	TaskId         string      `json:"task_id"`
	TaskName       string      `json:"task_name"`
	TaskType       string      `json:"task_type"`
	Description    string      `json:"description"`
	Status         string      `json:"status"`
	Priority       int         `json:"priority"`
	DeviceId       string      `json:"device_id"`
	DeviceName     string      `json:"device_name"`
	Progress       float64     `json:"progress"`
	ExecutionId    string      `json:"execution_id"`
	CreatedAt      *gtime.Time `json:"created_at"`
	UpdatedAt      *gtime.Time `json:"updated_at"`
	ScheduleAt     *gtime.Time `json:"schedule_at"`
	DeadlineAt     *gtime.Time `json:"deadline_at"`
	StartTime      *gtime.Time `json:"start_time"`
	EndTime        *gtime.Time `json:"end_time"`
	Duration       int64       `json:"duration"`
	RetryCount     int         `json:"retry_count"`
	ExecutionCount int         `json:"execution_count"`
}

// TaskDetailInfo 任务详细信息
type TaskDetailInfo struct {
	TaskInfo
	Config           map[string]interface{} `json:"config"`
	ErrorMessage     string                 `json:"error_message"`
	ExecutionHistory []TaskExecutionInfo    `json:"execution_history"`
	DependsOn        []string               `json:"depends_on"`
	DependedBy       []string               `json:"depended_by"`
	ResourceUsage    TaskResourceUsage      `json:"resource_usage"`
	Tags             []string               `json:"tags"`
	CreatedBy        string                 `json:"created_by"`
	UpdatedBy        string                 `json:"updated_by"`
	LastCheckTime    *gtime.Time            `json:"last_check_time"`
	NextScheduleTime *gtime.Time            `json:"next_schedule_time"`
}

// TaskExecutionInfo 任务执行信息
type TaskExecutionInfo struct {
	ExecutionId   string            `json:"execution_id"`
	Status        string            `json:"status"`
	StartTime     *gtime.Time       `json:"start_time"`
	EndTime       *gtime.Time       `json:"end_time"`
	Duration      int64             `json:"duration"`
	ErrorMessage  string            `json:"error_message"`
	Progress      float64           `json:"progress"`
	ExecutedBy    string            `json:"executed_by"`
	ResourceUsage TaskResourceUsage `json:"resource_usage"`
}

// TaskResourceUsage 任务资源使用情况
type TaskResourceUsage struct {
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskIO      float64 `json:"disk_io"`
	NetworkIO   float64 `json:"network_io"`
}

// TaskPriorityInfo 任务优先级信息
type TaskPriorityInfo struct {
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TaskCount   int    `json:"task_count"`
	Color       string `json:"color"`
}

// TaskLogInfo 任务日志信息
type TaskLogInfo struct {
	LogId       string      `json:"log_id"`
	TaskId      string      `json:"task_id"`
	ExecutionId string      `json:"execution_id"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	Timestamp   *gtime.Time `json:"timestamp"`
	Source      string      `json:"source"`
}

// TaskLogDetailInfo 任务日志详细信息
type TaskLogDetailInfo struct {
	TaskLogInfo
	StackTrace string                 `json:"stack_trace"`
	Context    map[string]interface{} `json:"context"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// TaskQueueInfo 任务队列信息
type TaskQueueInfo struct {
	QueueId   string      `json:"queue_id"`
	QueueName string      `json:"queue_name"`
	QueueType string      `json:"queue_type"`
	DeviceId  string      `json:"device_id"`
	TaskCount int         `json:"task_count"`
	Status    string      `json:"status"`
	Priority  int         `json:"priority"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

// TaskQueueDetailInfo 任务队列详细信息
type TaskQueueDetailInfo struct {
	TaskQueueInfo
	Tasks           []TaskInfo  `json:"tasks"`
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

// TaskScheduleInfo 任务调度信息
type TaskScheduleInfo struct {
	ScheduleId   string      `json:"schedule_id"`
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
	TaskId       string  `json:"task_id"`
	TaskName     string  `json:"task_name"`
	FailureCount int     `json:"failure_count"`
	FailureRate  float64 `json:"failure_rate"`
	LastError    string  `json:"last_error"`
}

// TaskPerformanceReport 任务性能报告
type TaskPerformanceReport struct {
	TaskId             string                 `json:"task_id"`
	TaskName           string                 `json:"task_name"`
	OverallScore       float64                `json:"overall_score"`
	ExecutionMetrics   TaskExecutionMetrics   `json:"execution_metrics"`
	ResourceMetrics    TaskResourceMetrics    `json:"resource_metrics"`
	ReliabilityMetrics TaskReliabilityMetrics `json:"reliability_metrics"`
	PerformanceTrend   []TaskTrendPoint       `json:"performance_trend"`
	Recommendations    []string               `json:"recommendations"`
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

// DeviceAssignmentOption 设备分配选项
type DeviceAssignmentOption struct {
	DeviceId           string  `json:"deviceId"`
	DeviceName         string  `json:"deviceName"`
	DeviceType         string  `json:"deviceType"`
	LoadScore          float64 `json:"loadScore"`
	CompatibilityScore float64 `json:"compatibilityScore"`
	AssignmentScore    float64 `json:"assignmentScore"`
	CurrentTasks       int     `json:"currentTasks"`
	MaxTasks           int     `json:"maxTasks"`
	EstimatedWaitTime  int     `json:"estimatedWaitTime"`
	Status             string  `json:"status"`
	Recommendation     string  `json:"recommendation"`
}

// TaskAssignment 任务分配
type TaskAssignment struct {
	TaskId    string   `json:"taskId" v:"required#任务ID不能为空"`
	DeviceIds []string `json:"deviceIds" v:"required#设备ID列表不能为空"`
	Priority  int      `json:"priority" d:"0"`
}

// TaskAssignmentResult 任务分配结果
type TaskAssignmentResult struct {
	TaskId           string `json:"taskId"`
	AssignedDeviceId string `json:"assignedDeviceId"`
	Status           string `json:"status"`
	Error            string `json:"error,omitempty"`
}

// TaskAssignmentHistory 任务分配历史
type TaskAssignmentHistory struct {
	Id              int64   `json:"id"`
	TaskId          string  `json:"taskId"`
	DeviceId        string  `json:"deviceId"`
	AssignmentScore float64 `json:"assignmentScore"`
	Strategy        string  `json:"strategy"`
	Reason          string  `json:"reason"`
	AssignedAt      string  `json:"assignedAt"`
	CompletedAt     string  `json:"completedAt,omitempty"`
	Status          string  `json:"status"`
}

// DeviceTaskQueueItem 设备任务队列项
type DeviceTaskQueueItem struct {
	Id                 int64  `json:"id"`
	TaskId             string `json:"taskId"`
	TaskName           string `json:"taskName"`
	TaskType           string `json:"taskType"`
	Priority           int    `json:"priority"`
	QueuePosition      int    `json:"queuePosition"`
	EstimatedStartTime string `json:"estimatedStartTime"`
	EstimatedDuration  int    `json:"estimatedDuration"`
	Status             string `json:"status"`
	QueuedAt           string `json:"queuedAt"`
}

// TaskAssignmentInfo 任务分配信息
type TaskAssignmentInfo struct {
	TaskId          string      `json:"task_id"`
	DeviceId        string      `json:"device_id"`
	AssignmentScore float64     `json:"assignment_score"`
	Strategy        string      `json:"strategy"`
	Reason          string      `json:"reason"`
	AssignedAt      *gtime.Time `json:"assigned_at"`
	CompletedAt     *gtime.Time `json:"completed_at,omitempty"`
	Status          string      `json:"status"`
}

// DeviceTaskAssignmentInfo 设备任务分配信息
type DeviceTaskAssignmentInfo struct {
	TaskId          string      `json:"task_id"`
	TaskName        string      `json:"task_name"`
	TaskType        string      `json:"task_type"`
	AssignmentScore float64     `json:"assignment_score"`
	Strategy        string      `json:"strategy"`
	Reason          string      `json:"reason"`
	AssignedAt      *gtime.Time `json:"assigned_at"`
	CompletedAt     *gtime.Time `json:"completed_at,omitempty"`
	Status          string      `json:"status"`
	Priority        int         `json:"priority"`
	Progress        float64     `json:"progress"`
}
