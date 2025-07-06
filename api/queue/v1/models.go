package v1

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 队列数据结构定义
// ===============================

// QueueInfo 队列基本信息
type QueueInfo struct {
	QueueId        string      `json:"queue_id"`
	QueueName      string      `json:"queue_name"`
	QueueType      string      `json:"queue_type"`
	Description    string      `json:"description"`
	Status         string      `json:"status"`
	Priority       int         `json:"priority"`
	MaxSize        int         `json:"max_size"`
	CurrentSize    int         `json:"current_size"`
	MaxConcurrency int         `json:"max_concurrency"`
	ActiveWorkers  int         `json:"active_workers"`
	Enabled        bool        `json:"enabled"`
	CreatedAt      *gtime.Time `json:"created_at"`
	UpdatedAt      *gtime.Time `json:"updated_at"`
	LastProcessed  *gtime.Time `json:"last_processed"`
}

// QueueDetailInfo 队列详细信息
type QueueDetailInfo struct {
	QueueInfo
	Config             QueueConfig         `json:"config"`
	StatusInfo         QueueStatusInfo     `json:"status_info"`
	RecentTasks        []QueueTaskInfo     `json:"recent_tasks"`
	PerformanceSummary PerformanceSummary  `json:"performance_summary"`
	ProcessingHistory  []ProcessingHistory `json:"processing_history"`
	ErrorSummary       ErrorSummary        `json:"error_summary"`
	CreatedBy          string              `json:"created_by"`
	UpdatedBy          string              `json:"updated_by"`
}

// QueueStatusInfo 队列状态信息
type QueueStatusInfo struct {
	QueueId            string      `json:"queue_id"`
	Status             string      `json:"status"`
	IsHealthy          bool        `json:"is_healthy"`
	PendingTasks       int         `json:"pending_tasks"`
	ProcessingTasks    int         `json:"processing_tasks"`
	CompletedTasks     int         `json:"completed_tasks"`
	FailedTasks        int         `json:"failed_tasks"`
	RetryingTasks      int         `json:"retrying_tasks"`
	ActiveWorkers      int         `json:"active_workers"`
	AverageWaitTime    float64     `json:"average_wait_time"`
	AverageProcessTime float64     `json:"average_process_time"`
	Throughput         float64     `json:"throughput"`
	ErrorRate          float64     `json:"error_rate"`
	LastHealthCheck    *gtime.Time `json:"last_health_check"`
	StatusChangedAt    *gtime.Time `json:"status_changed_at"`
}

// QueueTaskInfo 队列任务信息
type QueueTaskInfo struct {
	TaskId       string                 `json:"task_id"`
	QueueId      string                 `json:"queue_id"`
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
	WorkerId     string                 `json:"worker_id"`
}

// QueueConfig 队列配置
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

// QueueConfigHistory 队列配置历史
type QueueConfigHistory struct {
	HistoryId    string      `json:"history_id"`
	QueueId      string      `json:"queue_id"`
	ConfigData   QueueConfig `json:"config_data"`
	ChangedBy    string      `json:"changed_by"`
	ChangeReason string      `json:"change_reason"`
	CreatedAt    *gtime.Time `json:"created_at"`
}

// QueueMetrics 队列指标
type QueueMetrics struct {
	QueueId              string            `json:"queue_id"`
	TimeRange            string            `json:"time_range"`
	TotalProcessed       int               `json:"total_processed"`
	SuccessfulTasks      int               `json:"successful_tasks"`
	FailedTasks          int               `json:"failed_tasks"`
	AverageThroughput    float64           `json:"average_throughput"`
	PeakThroughput       float64           `json:"peak_throughput"`
	AverageWaitTime      float64           `json:"average_wait_time"`
	AverageProcessTime   float64           `json:"average_process_time"`
	SuccessRate          float64           `json:"success_rate"`
	ErrorRate            float64           `json:"error_rate"`
	TrendData            []QueueTrendPoint `json:"trend_data"`
	ErrorDistribution    map[string]int    `json:"error_distribution"`
	PriorityDistribution map[string]int    `json:"priority_distribution"`
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

// QueueStatistics 队列统计信息
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

// QueueActivityRank 队列活跃度排名
type QueueActivityRank struct {
	QueueId      string      `json:"queue_id"`
	QueueName    string      `json:"queue_name"`
	TaskCount    int         `json:"task_count"`
	Throughput   float64     `json:"throughput"`
	SuccessRate  float64     `json:"success_rate"`
	LastActivity *gtime.Time `json:"last_activity"`
}

// QueuePerformanceReport 队列性能报告
type QueuePerformanceReport struct {
	QueueId            string                  `json:"queue_id"`
	QueueName          string                  `json:"queue_name"`
	ReportPeriod       string                  `json:"report_period"`
	OverallScore       float64                 `json:"overall_score"`
	ThroughputMetrics  ThroughputMetrics       `json:"throughput_metrics"`
	LatencyMetrics     LatencyMetrics          `json:"latency_metrics"`
	ReliabilityMetrics QueueReliabilityMetrics `json:"reliability_metrics"`
	ResourceMetrics    QueueResourceMetrics    `json:"resource_metrics"`
	PerformanceTrend   []QueueTrendPoint       `json:"performance_trend"`
	Issues             []PerformanceIssue      `json:"issues"`
	Recommendations    []string                `json:"recommendations"`
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
	AverageWorkerCount int     `json:"average_worker_count"`
	PeakWorkerCount    int     `json:"peak_worker_count"`
	WorkerUtilization  float64 `json:"worker_utilization"`
	QueueUtilization   float64 `json:"queue_utilization"`
	MemoryUsage        float64 `json:"memory_usage"`
}

// PerformanceSummary 性能摘要
type PerformanceSummary struct {
	TasksProcessedToday int     `json:"tasks_processed_today"`
	AverageThroughput   float64 `json:"average_throughput"`
	CurrentSuccessRate  float64 `json:"current_success_rate"`
	AverageWaitTime     float64 `json:"average_wait_time"`
	HealthScore         float64 `json:"health_score"`
}

// ProcessingHistory 处理历史
type ProcessingHistory struct {
	Date           string  `json:"date"`
	TasksProcessed int     `json:"tasks_processed"`
	SuccessCount   int     `json:"success_count"`
	FailureCount   int     `json:"failure_count"`
	SuccessRate    float64 `json:"success_rate"`
	AvgProcessTime float64 `json:"avg_process_time"`
}

// ErrorSummary 错误摘要
type ErrorSummary struct {
	TotalErrors  int               `json:"total_errors"`
	RecentErrors int               `json:"recent_errors"`
	ErrorRate    float64           `json:"error_rate"`
	TopErrors    []ErrorInfo       `json:"top_errors"`
	ErrorTrend   []ErrorTrendPoint `json:"error_trend"`
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	ErrorType    string      `json:"error_type"`
	ErrorMessage string      `json:"error_message"`
	ErrorCount   int         `json:"error_count"`
	LastOccurred *gtime.Time `json:"last_occurred"`
}

// ErrorTrendPoint 错误趋势点
type ErrorTrendPoint struct {
	Date       string      `json:"date"`
	ErrorCount int         `json:"error_count"`
	Timestamp  *gtime.Time `json:"timestamp"`
}

// PerformanceIssue 性能问题
type PerformanceIssue struct {
	IssueType   string `json:"issue_type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Suggestion  string `json:"suggestion"`
}
