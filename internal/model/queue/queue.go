package queue

import (
	"time"

	"OneGfServer/internal/consts"
	"OneGfServer/internal/model/common"

	"github.com/gogf/gf/v2/os/gtime"
)

// 队列状态常量 - 使用统一常量
const (
	QueueStatusActive  = consts.QueueStatusActive
	QueueStatusPaused  = consts.QueueStatusPaused
	QueueStatusStopped = consts.QueueStatusStopped
	QueueStatusError   = consts.QueueStatusError
)

// 队列类型常量 - 使用统一常量
const (
	QueueTypeFIFO     = consts.QueueTypeFIFO
	QueueTypePriority = consts.QueueTypePriority
	QueueTypeDelay    = consts.QueueTypeDelay
	QueueTypeLIFO     = consts.QueueTypeLIFO
)

// 队列操作类型常量 - 使用统一常量
const (
	QueueOperationAdd            = consts.QueueOperationAdd
	QueueOperationRemove         = consts.QueueOperationRemove
	QueueOperationPriorityChange = consts.QueueOperationPriorityChange
	QueueOperationPositionChange = consts.QueueOperationPositionChange
	QueueOperationStart          = consts.QueueOperationStart
	QueueOperationPause          = consts.QueueOperationPause
	QueueOperationResume         = consts.QueueOperationResume
	QueueOperationCancel         = consts.QueueOperationCancel
)

// 调度策略常量 - 使用统一常量
const (
	SchedulingStrategyPriorityFirst = consts.SchedulingStrategyPriorityFirst
	SchedulingStrategyFIFO          = consts.SchedulingStrategyFIFO
	SchedulingStrategyLIFO          = consts.SchedulingStrategyLIFO
	SchedulingStrategyWeighted      = consts.SchedulingStrategyWeighted
)

// 操作来源常量 - 使用统一常量
const (
	OperationSourceManual    = consts.OperationSourceManual
	OperationSourceSystem    = consts.OperationSourceSystem
	OperationSourceAPI       = consts.OperationSourceAPI
	OperationSourceScheduler = consts.OperationSourceScheduler
)

// 队列优先级常量 - 使用统一常量
const (
	QueuePriorityLowest  = consts.QueuePriorityLowest
	QueuePriorityLow     = consts.QueuePriorityLow
	QueuePriorityNormal  = consts.QueuePriorityNormal
	QueuePriorityHigh    = consts.QueuePriorityHigh
	QueuePriorityUrgent  = consts.QueuePriorityUrgent
	QueuePriorityHighest = consts.QueuePriorityHighest
)

// ===============================
// 基础CRUD操作 Input/Output
// ===============================

// CreateQueueInput 创建队列输入
type CreateQueueInput struct {
	Queue *Queue `json:"queue"`
}

// CreateQueueOutput 创建队列输出
type CreateQueueOutput struct {
	QueueID string `json:"queue_id"`
	Message string `json:"message"`
}

// GetQueueByIDInput 根据ID获取队列输入
type GetQueueByIDInput struct {
	ID int64 `json:"id"`
}

// GetQueueByIDOutput 根据ID获取队列输出
type GetQueueByIDOutput struct {
	Queue *Queue `json:"queue"`
}

// GetQueueByQueueIDInput 根据队列ID获取队列输入
type GetQueueByQueueIDInput struct {
	QueueID string `json:"queue_id"`
}

// GetQueueByQueueIDOutput 根据队列ID获取队列输出
type GetQueueByQueueIDOutput struct {
	Queue *Queue `json:"queue"`
}

// UpdateQueueInput 更新队列输入
type UpdateQueueInput struct {
	Queue *Queue `json:"queue"`
}

// UpdateQueueOutput 更新队列输出
type UpdateQueueOutput struct {
	Message string `json:"message"`
}

// DeleteQueueInput 删除队列输入
type DeleteQueueInput struct {
	QueueID string `json:"queue_id"`
	Force   bool   `json:"force"`
}

// DeleteQueueOutput 删除队列输出
type DeleteQueueOutput struct {
	Message string `json:"message"`
}

// ===============================
// 查询操作 Input/Output
// ===============================

// GetQueueListInput 获取队列列表输入
type GetQueueListInput struct {
	Filter     *QueueFilter              `json:"filter"`
	Sort       *QueueSortOption          `json:"sort"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetQueueListOutput 获取队列列表输出
type GetQueueListOutput struct {
	common.PaginationResponse[Queue] `json:",inline"`
}

// GetQueuesByTypeInput 根据队列类型获取队列输入
type GetQueuesByTypeInput struct {
	QueueTypes []string `json:"queue_types"`
}

// GetQueuesByTypeOutput 根据队列类型获取队列输出
type GetQueuesByTypeOutput struct {
	Queues []Queue `json:"queues"`
}

// GetQueuesByStatusInput 根据状态获取队列输入
type GetQueuesByStatusInput struct {
	Statuses []string `json:"statuses"`
}

// GetQueuesByStatusOutput 根据状态获取队列输出
type GetQueuesByStatusOutput struct {
	Queues []Queue `json:"queues"`
}

// GetActiveQueuesInput 获取活跃队列输入
type GetActiveQueuesInput struct{}

// GetActiveQueuesOutput 获取活跃队列输出
type GetActiveQueuesOutput struct {
	Queues []Queue `json:"queues"`
}

// ===============================
// 队列控制操作 Input/Output
// ===============================

// StartQueueInput 启动队列输入
type StartQueueInput struct {
	QueueID string `json:"queue_id"`
	Force   bool   `json:"force"`
}

// StartQueueOutput 启动队列输出
type StartQueueOutput struct {
	Message string `json:"message"`
}

// PauseQueueInput 暂停队列输入
type PauseQueueInput struct {
	QueueID string `json:"queue_id"`
	Reason  string `json:"reason"`
}

// PauseQueueOutput 暂停队列输出
type PauseQueueOutput struct {
	Message string `json:"message"`
}

// ResumeQueueInput 恢复队列输入
type ResumeQueueInput struct {
	QueueID string `json:"queue_id"`
	Reason  string `json:"reason"`
}

// ResumeQueueOutput 恢复队列输出
type ResumeQueueOutput struct {
	Message string `json:"message"`
}

// StopQueueInput 停止队列输入
type StopQueueInput struct {
	QueueID string `json:"queue_id"`
	Reason  string `json:"reason"`
	Force   bool   `json:"force"`
}

// StopQueueOutput 停止队列输出
type StopQueueOutput struct {
	Message string `json:"message"`
}

// ClearQueueInput 清空队列输入
type ClearQueueInput struct {
	QueueID string `json:"queue_id"`
	Reason  string `json:"reason"`
	Force   bool   `json:"force"`
}

// ClearQueueOutput 清空队列输出
type ClearQueueOutput struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ResetQueueInput 重置队列输入
type ResetQueueInput struct {
	QueueID string `json:"queue_id"`
	Reason  string `json:"reason"`
}

// ResetQueueOutput 重置队列输出
type ResetQueueOutput struct {
	Message string `json:"message"`
}

// ===============================
// 队列任务管理 Input/Output
// ===============================

// GetQueueTasksInput 获取队列任务输入
type GetQueueTasksInput struct {
	QueueID    string                    `json:"queue_id"`
	Filter     *TaskFilter               `json:"filter"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetQueueTasksOutput 获取队列任务输出
type GetQueueTasksOutput struct {
	common.PaginationResponse[QueueTask] `json:",inline"`
}

// ReorderQueueTasksInput 重新排序队列任务输入
type ReorderQueueTasksInput struct {
	QueueID   string   `json:"queue_id"`
	TaskIDs   []string `json:"task_ids"`
	Strategy  string   `json:"strategy"`
	Positions []int    `json:"positions"`
}

// ReorderQueueTasksOutput 重新排序队列任务输出
type ReorderQueueTasksOutput struct {
	AffectedCount int    `json:"affected_count"`
	Message       string `json:"message"`
}

// RemoveQueueTaskInput 移除队列任务输入
type RemoveQueueTaskInput struct {
	QueueID string `json:"queue_id"`
	TaskID  string `json:"task_id"`
	Reason  string `json:"reason"`
}

// RemoveQueueTaskOutput 移除队列任务输出
type RemoveQueueTaskOutput struct {
	Message string `json:"message"`
}

// ===============================
// 队列监控统计 Input/Output
// ===============================

// GetQueueStatisticsInput 获取队列统计信息输入
type GetQueueStatisticsInput struct {
	Filter *QueueFilter `json:"filter"`
}

// GetQueueStatisticsOutput 获取队列统计信息输出
type GetQueueStatisticsOutput struct {
	Statistics *QueueStatistics `json:"statistics"`
}

// GetQueuePerformanceReportInput 获取队列性能报告输入
type GetQueuePerformanceReportInput struct {
	QueueID   string `json:"queue_id"`
	TimeRange string `json:"time_range"`
}

// GetQueuePerformanceReportOutput 获取队列性能报告输出
type GetQueuePerformanceReportOutput struct {
	PerformanceReport *QueuePerformanceReport `json:"performance_report"`
}

// ===============================
// 队列配置管理 Input/Output
// ===============================

// GetQueueConfigInput 获取队列配置输入
type GetQueueConfigInput struct {
	QueueID string `json:"queue_id"`
}

// GetQueueConfigOutput 获取队列配置输出
type GetQueueConfigOutput struct {
	Config *QueueConfig `json:"config"`
}

// UpdateQueueConfigInput 更新队列配置输入
type UpdateQueueConfigInput struct {
	QueueID string       `json:"queue_id"`
	Config  *QueueConfig `json:"config"`
}

// UpdateQueueConfigOutput 更新队列配置输出
type UpdateQueueConfigOutput struct {
	Message string `json:"message"`
}

// GetQueueConfigHistoryInput 获取队列配置历史输入
type GetQueueConfigHistoryInput struct {
	QueueID    string                    `json:"queue_id"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetQueueConfigHistoryOutput 获取队列配置历史输出
type GetQueueConfigHistoryOutput struct {
	common.PaginationResponse[QueueConfigHistory] `json:",inline"`
}

// ===============================
// 数据模型定义
// ===============================

// Queue 队列主表模型
type Queue struct {
	ID             int64        `json:"id"`
	QueueID        string       `json:"queue_id"`
	QueueName      string       `json:"queue_name"`
	QueueType      string       `json:"queue_type"`
	Description    string       `json:"description"`
	Status         string       `json:"status"`
	Priority       int          `json:"priority"`
	MaxSize        int          `json:"max_size"`
	CurrentSize    int          `json:"current_size"`
	MaxConcurrency int          `json:"max_concurrency"`
	ActiveWorkers  int          `json:"active_workers"`
	Enabled        bool         `json:"enabled"`
	Config         *QueueConfig `json:"config"`
	CreatedAt      *gtime.Time  `json:"created_at"`
	UpdatedAt      *gtime.Time  `json:"updated_at"`
	LastProcessed  *gtime.Time  `json:"last_processed"`
	CreatedBy      int64        `json:"created_by"`
	UpdatedBy      int64        `json:"updated_by"`
}

// QueueTask 队列任务模型
type QueueTask struct {
	TaskID       string                 `json:"task_id"`
	QueueID      string                 `json:"queue_id"`
	Status       string                 `json:"status"`
	Priority     int                    `json:"priority"`
	Position     int                    `json:"position"`
	TaskData     map[string]interface{} `json:"task_data"`
	Metadata     map[string]interface{} `json:"metadata"`
	RetryCount   int                    `json:"retry_count"`
	MaxRetries   int                    `json:"max_retries"`
	CreatedAt    *gtime.Time            `json:"created_at"`
	StartedAt    *gtime.Time            `json:"started_at"`
	CompletedAt  *gtime.Time            `json:"completed_at"`
	ScheduledAt  *gtime.Time            `json:"scheduled_at"`
	ProcessTime  float64                `json:"process_time"`
	ErrorMessage string                 `json:"error_message"`
	WorkerID     string                 `json:"worker_id"`
}

// QueueConfig 队列配置模型
type QueueConfig struct {
	ProcessingStrategy   string                 `json:"processing_strategy"`
	RetryStrategy        string                 `json:"retry_strategy"`
	RetryDelayMs         int                    `json:"retry_delay_ms"`
	MaxRetryInterval     int                    `json:"max_retry_interval"`
	DeadLetterEnabled    bool                   `json:"dead_letter_enabled"`
	DeadLetterMaxRetries int                    `json:"dead_letter_max_retries"`
	BatchSize            int                    `json:"batch_size"`
	TimeoutMs            int                    `json:"timeout_ms"`
	AutoScale            bool                   `json:"auto_scale"`
	MinWorkers           int                    `json:"min_workers"`
	MaxWorkers           int                    `json:"max_workers"`
	CustomSettings       map[string]interface{} `json:"custom_settings"`
}

// QueueConfigHistory 队列配置历史模型
type QueueConfigHistory struct {
	HistoryID    string       `json:"history_id"`
	QueueID      string       `json:"queue_id"`
	ConfigData   *QueueConfig `json:"config_data"`
	ChangedBy    string       `json:"changed_by"`
	ChangeReason string       `json:"change_reason"`
	CreatedAt    *gtime.Time  `json:"created_at"`
}

// QueueStatistics 队列统计信息模型
type QueueStatistics struct {
	TotalQueues        int                 `json:"total_queues"`
	ActiveQueues       int                 `json:"active_queues"`
	PausedQueues       int                 `json:"paused_queues"`
	StoppedQueues      int                 `json:"stopped_queues"`
	TotalTasks         int                 `json:"total_tasks"`
	PendingTasks       int                 `json:"pending_tasks"`
	ProcessingTasks    int                 `json:"processing_tasks"`
	CompletedTasks     int                 `json:"completed_tasks"`
	FailedTasks        int                 `json:"failed_tasks"`
	StatusDistribution map[string]int      `json:"status_distribution"`
	TypeDistribution   map[string]int      `json:"type_distribution"`
	OverallThroughput  float64             `json:"overall_throughput"`
	OverallSuccessRate float64             `json:"overall_success_rate"`
	TrendData          []QueueTrendPoint   `json:"trend_data"`
	TopActiveQueues    []QueueActivityRank `json:"top_active_queues"`
}

// QueueTrendPoint 队列趋势数据点
type QueueTrendPoint struct {
	Timestamp      *gtime.Time `json:"timestamp"`
	TaskCount      int         `json:"task_count"`
	SuccessCount   int         `json:"success_count"`
	FailureCount   int         `json:"failure_count"`
	SuccessRate    float64     `json:"success_rate"`
	AvgWaitTime    float64     `json:"avg_wait_time"`
	AvgProcessTime float64     `json:"avg_process_time"`
	Throughput     float64     `json:"throughput"`
}

// QueueActivityRank 队列活跃度排名
type QueueActivityRank struct {
	QueueID      string      `json:"queue_id"`
	QueueName    string      `json:"queue_name"`
	TaskCount    int         `json:"task_count"`
	Throughput   float64     `json:"throughput"`
	SuccessRate  float64     `json:"success_rate"`
	LastActivity *gtime.Time `json:"last_activity"`
}

// QueuePerformanceReport 队列性能报告
type QueuePerformanceReport struct {
	QueueId            string                   `json:"queue_id"`
	QueueName          string                   `json:"queue_name"`
	ReportPeriod       string                   `json:"report_period"`
	OverallScore       float64                  `json:"overall_score"`
	ThroughputMetrics  *ThroughputMetrics       `json:"throughput_metrics"`
	LatencyMetrics     *LatencyMetrics          `json:"latency_metrics"`
	ReliabilityMetrics *QueueReliabilityMetrics `json:"reliability_metrics"`
	ResourceMetrics    *QueueResourceMetrics    `json:"resource_metrics"`
	PerformanceTrend   []QueueTrendPoint        `json:"performance_trend"`
	Issues             []PerformanceIssue       `json:"issues"`
	Recommendations    []string                 `json:"recommendations"`
}

// ThroughputMetrics 吞吐量指标
type ThroughputMetrics struct {
	AverageThroughput  float64 `json:"average_throughput"`
	PeakThroughput     float64 `json:"peak_throughput"`
	MinThroughput      float64 `json:"min_throughput"`
	ThroughputVariance float64 `json:"throughput_variance"`
}

// LatencyMetrics 延迟指标
type LatencyMetrics struct {
	AverageWaitTime    float64 `json:"average_wait_time"`
	AverageProcessTime float64 `json:"average_process_time"`
	P50WaitTime        float64 `json:"p50_wait_time"`
	P95WaitTime        float64 `json:"p95_wait_time"`
	P99WaitTime        float64 `json:"p99_wait_time"`
}

// QueueReliabilityMetrics 队列可靠性指标
type QueueReliabilityMetrics struct {
	SuccessRate      float64 `json:"success_rate"`
	ErrorRate        float64 `json:"error_rate"`
	RetryRate        float64 `json:"retry_rate"`
	DeadLetterRate   float64 `json:"dead_letter_rate"`
	AvailabilityRate float64 `json:"availability_rate"`
}

// QueueResourceMetrics 队列资源指标
type QueueResourceMetrics struct {
	AverageCpuUsage    float64 `json:"average_cpu_usage"`
	PeakCpuUsage       float64 `json:"peak_cpu_usage"`
	AverageMemoryUsage float64 `json:"average_memory_usage"`
	PeakMemoryUsage    float64 `json:"peak_memory_usage"`
	TotalDiskIO        float64 `json:"total_disk_io"`
	TotalNetworkIO     float64 `json:"total_network_io"`
}

// PerformanceIssue 性能问题
type PerformanceIssue struct {
	IssueType      string  `json:"issue_type"`
	Severity       string  `json:"severity"`
	Description    string  `json:"description"`
	Impact         float64 `json:"impact"`
	Recommendation string  `json:"recommendation"`
}

// ===============================
// 过滤和排序选项
// ===============================

// QueueFilter 队列过滤条件
type QueueFilter struct {
	Status    []string   `json:"status"`
	Type      []string   `json:"type"`
	Priority  *int       `json:"priority"`
	Keyword   *string    `json:"keyword"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Enabled   *bool      `json:"enabled"`
}

// QueueSortOption 队列排序选项
type QueueSortOption struct {
	Field string `json:"field"` // id, name, status, priority, created_at, updated_at
	Order string `json:"order"` // asc, desc
}

// TaskFilter 任务过滤条件
type TaskFilter struct {
	Status    *string    `json:"status"`
	Priority  *int       `json:"priority"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}
