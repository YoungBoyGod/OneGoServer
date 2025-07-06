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

// ===============================
// 队列Logic层相关 Input/Output
// ===============================

// ValidateQueueCreationInput 验证队列创建输入
type ValidateQueueCreationInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// ValidateQueueCreationOutput 验证队列创建输出
type ValidateQueueCreationOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ValidateQueueOperationInput 验证队列操作输入
type ValidateQueueOperationInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
	Operation string                 `json:"operation"`  // 操作类型
}

// ValidateQueueOperationOutput 验证队列操作输出
type ValidateQueueOperationOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// CalculateQueueHealthScoreInput 计算队列健康度评分输入
type CalculateQueueHealthScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateQueueHealthScoreOutput 计算队列健康度评分输出
type CalculateQueueHealthScoreOutput struct {
	HealthScore float64                `json:"health_score"` // 健康度评分
	Components  map[string]interface{} `json:"components"`   // 评分组件
}

// ValidateTaskEnqueueInput 验证任务入队输入
type ValidateTaskEnqueueInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
	TaskData  map[string]interface{} `json:"task_data"`  // 任务数据
}

// ValidateTaskEnqueueOutput 验证任务入队输出
type ValidateTaskEnqueueOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// SortQueueTasksInput 排序队列任务输入
type SortQueueTasksInput struct {
	QueueType string                   `json:"queue_type"` // 队列类型
	Tasks     []map[string]interface{} `json:"tasks"`      // 任务列表
}

// SortQueueTasksOutput 排序队列任务输出
type SortQueueTasksOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// CalculateLoadBalanceInput 计算负载均衡输入
type CalculateLoadBalanceInput struct {
	Queues []map[string]interface{} `json:"queues"` // 队列列表
}

// CalculateLoadBalanceOutput 计算负载均衡输出
type CalculateLoadBalanceOutput struct {
	LoadBalance map[string]interface{} `json:"load_balance"` // 负载均衡结果
}

// SelectOptimalQueueInput 选择最优队列输入
type SelectOptimalQueueInput struct {
	TaskData        map[string]interface{}   `json:"task_data"`        // 任务数据
	AvailableQueues []map[string]interface{} `json:"available_queues"` // 可用队列列表
}

// SelectOptimalQueueOutput 选择最优队列输出
type SelectOptimalQueueOutput struct {
	OptimalQueue map[string]interface{} `json:"optimal_queue"` // 最优队列
	Score        float64                `json:"score"`         // 选择评分
	Reason       string                 `json:"reason"`        // 选择原因
}

// CalculateQueueStatisticsInput 计算队列统计信息输入
type CalculateQueueStatisticsInput struct {
	QueueData      map[string]interface{}   `json:"queue_data"`      // 队列数据
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// CalculateQueueStatisticsOutput 计算队列统计信息输出
type CalculateQueueStatisticsOutput struct {
	Statistics map[string]interface{} `json:"statistics"` // 统计信息
}

// ValidateQueueConfigurationInput 验证队列配置输入
type ValidateQueueConfigurationInput struct {
	Config map[string]interface{} `json:"config"` // 配置数据
}

// ValidateQueueConfigurationOutput 验证队列配置输出
type ValidateQueueConfigurationOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// CalculateStatusScoreInput 计算状态评分输入
type CalculateStatusScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateStatusScoreOutput 计算状态评分输出
type CalculateStatusScoreOutput struct {
	Score float64 `json:"score"` // 状态评分
}

// CalculatePerformanceScoreInput 计算性能评分输入
type CalculatePerformanceScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculatePerformanceScoreOutput 计算性能评分输出
type CalculatePerformanceScoreOutput struct {
	Score float64 `json:"score"` // 性能评分
}

// CalculateErrorScoreInput 计算错误评分输入
type CalculateErrorScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateErrorScoreOutput 计算错误评分输出
type CalculateErrorScoreOutput struct {
	Score float64 `json:"score"` // 错误评分
}

// CalculateResourceScoreInput 计算资源评分输入
type CalculateResourceScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateResourceScoreOutput 计算资源评分输出
type CalculateResourceScoreOutput struct {
	Score float64 `json:"score"` // 资源评分
}

// CalculateResponseScoreInput 计算响应评分输入
type CalculateResponseScoreInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateResponseScoreOutput 计算响应评分输出
type CalculateResponseScoreOutput struct {
	Score float64 `json:"score"` // 响应评分
}

// ValidateTaskDataInput 验证任务数据输入
type ValidateTaskDataInput struct {
	TaskData map[string]interface{} `json:"task_data"` // 任务数据
}

// ValidateTaskDataOutput 验证任务数据输出
type ValidateTaskDataOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// SortFIFOInput FIFO排序输入
type SortFIFOInput struct {
	Tasks []map[string]interface{} `json:"tasks"` // 任务列表
}

// SortFIFOOutput FIFO排序输出
type SortFIFOOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// SortLIFOInput LIFO排序输入
type SortLIFOInput struct {
	Tasks []map[string]interface{} `json:"tasks"` // 任务列表
}

// SortLIFOOutput LIFO排序输出
type SortLIFOOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// SortByPriorityInput 按优先级排序输入
type SortByPriorityInput struct {
	Tasks []map[string]interface{} `json:"tasks"` // 任务列表
}

// SortByPriorityOutput 按优先级排序输出
type SortByPriorityOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// SortRoundRobinInput 轮询排序输入
type SortRoundRobinInput struct {
	Tasks []map[string]interface{} `json:"tasks"` // 任务列表
}

// SortRoundRobinOutput 轮询排序输出
type SortRoundRobinOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// SortByWeightInput 按权重排序输入
type SortByWeightInput struct {
	Tasks []map[string]interface{} `json:"tasks"` // 任务列表
}

// SortByWeightOutput 按权重排序输出
type SortByWeightOutput struct {
	SortedTasks []map[string]interface{} `json:"sorted_tasks"` // 排序后的任务列表
}

// CalculateTaskWeightInput 计算任务权重输入
type CalculateTaskWeightInput struct {
	Task map[string]interface{} `json:"task"` // 任务数据
}

// CalculateTaskWeightOutput 计算任务权重输出
type CalculateTaskWeightOutput struct {
	Weight float64 `json:"weight"` // 任务权重
}

// CalculateQueueScoreInput 计算队列评分输入
type CalculateQueueScoreInput struct {
	Queue    map[string]interface{} `json:"queue"`     // 队列数据
	TaskData map[string]interface{} `json:"task_data"` // 任务数据
}

// CalculateQueueScoreOutput 计算队列评分输出
type CalculateQueueScoreOutput struct {
	Score float64 `json:"score"` // 队列评分
}

// CalculateUtilizationRateInput 计算利用率输入
type CalculateUtilizationRateInput struct {
	QueueData map[string]interface{} `json:"queue_data"` // 队列数据
}

// CalculateUtilizationRateOutput 计算利用率输出
type CalculateUtilizationRateOutput struct {
	UtilizationRate float64 `json:"utilization_rate"` // 利用率
}

// CalculateAverageProcessingTimeInput 计算平均处理时间输入
type CalculateAverageProcessingTimeInput struct {
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// CalculateAverageProcessingTimeOutput 计算平均处理时间输出
type CalculateAverageProcessingTimeOutput struct {
	AverageProcessingTime float64 `json:"average_processing_time"` // 平均处理时间
}

// CalculateThroughputInput 计算吞吐量输入
type CalculateThroughputInput struct {
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// CalculateThroughputOutput 计算吞吐量输出
type CalculateThroughputOutput struct {
	Throughput float64 `json:"throughput"` // 吞吐量
}

// CalculateErrorRateInput 计算错误率输入
type CalculateErrorRateInput struct {
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// CalculateErrorRateOutput 计算错误率输出
type CalculateErrorRateOutput struct {
	ErrorRate float64 `json:"error_rate"` // 错误率
}

// AnalyzeQueueTrendInput 分析队列趋势输入
type AnalyzeQueueTrendInput struct {
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// AnalyzeQueueTrendOutput 分析队列趋势输出
type AnalyzeQueueTrendOutput struct {
	Trend map[string]interface{} `json:"trend"` // 趋势分析结果
}

// PredictQueueBehaviorInput 预测队列行为输入
type PredictQueueBehaviorInput struct {
	HistoricalData []map[string]interface{} `json:"historical_data"` // 历史数据
}

// PredictQueueBehaviorOutput 预测队列行为输出
type PredictQueueBehaviorOutput struct {
	Prediction map[string]interface{} `json:"prediction"` // 预测结果
}

// ValidateBasicConfigInput 验证基础配置输入
type ValidateBasicConfigInput struct {
	Config map[string]interface{} `json:"config"` // 配置数据
}

// ValidateBasicConfigOutput 验证基础配置输出
type ValidateBasicConfigOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ValidatePerformanceConfigInput 验证性能配置输入
type ValidatePerformanceConfigInput struct {
	Config map[string]interface{} `json:"config"` // 配置数据
}

// ValidatePerformanceConfigOutput 验证性能配置输出
type ValidatePerformanceConfigOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ValidateSecurityConfigInput 验证安全配置输入
type ValidateSecurityConfigInput struct {
	Config map[string]interface{} `json:"config"` // 配置数据
}

// ValidateSecurityConfigOutput 验证安全配置输出
type ValidateSecurityConfigOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ValidateAccessControlInput 验证访问控制输入
type ValidateAccessControlInput struct {
	AccessControl map[string]interface{} `json:"access_control"` // 访问控制数据
}

// ValidateAccessControlOutput 验证访问控制输出
type ValidateAccessControlOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}
