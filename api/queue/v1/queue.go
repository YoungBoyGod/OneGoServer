package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 1. 队列基础管理 API (5个)
// ===============================

// CreateQueueReq 创建队列请求
type CreateQueueReq struct {
	g.Meta         `path:"/queue/create" method:"post" tags:"队列管理" summary:"创建队列"`
	QueueName      string                 `json:"queue_name" v:"required|length:1,100#队列名称不能为空|队列名称长度为1-100字符"`
	QueueType      string                 `json:"queue_type" v:"required|in:fifo,priority,delay,lifo#队列类型不能为空"`
	Description    string                 `json:"description,omitempty" v:"max-length:500#队列描述最大500字符"`
	MaxSize        int                    `json:"max_size" d:"1000" v:"between:1,10000#队列最大大小为1-10000"`
	MaxConcurrency int                    `json:"max_concurrency" d:"10" v:"between:1,100#最大并发数为1-100"`
	Priority       int                    `json:"priority" d:"5" v:"between:1,10#队列优先级为1-10"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Enabled        bool                   `json:"enabled" d:"true"`
}

type CreateQueueRes struct {
	QueueId string `json:"queue_id"`
	Message string `json:"message"`
}

// GetQueueListReq 获取队列列表请求
type GetQueueListReq struct {
	g.Meta    `path:"/queue/list" method:"get" tags:"队列管理" summary:"获取队列列表"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status    string      `json:"status,omitempty" v:"in:active,paused,stopped,error#状态值无效"`
	QueueType string      `json:"queue_type,omitempty" v:"in:fifo,priority,delay,lifo#队列类型无效"`
	Priority  int         `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
	Keyword   string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	SortBy    string      `json:"sort_by" d:"created_at" v:"in:created_at,updated_at,priority,task_count#排序字段无效"`
	SortOrder string      `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"`
}

type GetQueueListRes struct {
	List  []QueueInfo `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// GetQueueDetailReq 获取队列详情请求
type GetQueueDetailReq struct {
	g.Meta  `path:"/queue/{queueId}" method:"get" tags:"队列管理" summary:"获取队列详情"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueDetailRes struct {
	QueueDetail QueueDetailInfo `json:"queue_detail"`
}

// UpdateQueueReq 更新队列请求
type UpdateQueueReq struct {
	g.Meta         `path:"/queue/{queueId}" method:"put" tags:"队列管理" summary:"更新队列"`
	QueueId        string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	QueueName      string                 `json:"queue_name,omitempty" v:"length:1,100#队列名称长度为1-100字符"`
	Description    string                 `json:"description,omitempty" v:"max-length:500#队列描述最大500字符"`
	MaxSize        int                    `json:"max_size,omitempty" v:"between:1,10000#队列最大大小为1-10000"`
	MaxConcurrency int                    `json:"max_concurrency,omitempty" v:"between:1,100#最大并发数为1-100"`
	Priority       int                    `json:"priority,omitempty" v:"between:1,10#队列优先级为1-10"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Enabled        bool                   `json:"enabled,omitempty"`
}

type UpdateQueueRes struct {
	Message string `json:"message"`
}

// DeleteQueueReq 删除队列请求
type DeleteQueueReq struct {
	g.Meta  `path:"/queue/{queueId}" method:"delete" tags:"队列管理" summary:"删除队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Force   bool   `json:"force,omitempty"`
}

type DeleteQueueRes struct {
	Message string `json:"message"`
}

// ===============================
// 2. 队列操作控制 API (6个)
// ===============================

// StartQueueReq 启动队列请求
type StartQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/start" method:"post" tags:"队列控制" summary:"启动队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Force   bool   `json:"force,omitempty"`
}

type StartQueueRes struct {
	Message string `json:"message"`
}

// PauseQueueReq 暂停队列请求
type PauseQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/pause" method:"post" tags:"队列控制" summary:"暂停队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#暂停原因最大200字符"`
}

type PauseQueueRes struct {
	Message string `json:"message"`
}

// ResumeQueueReq 恢复队列请求
type ResumeQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/resume" method:"post" tags:"队列控制" summary:"恢复队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#恢复原因最大200字符"`
}

type ResumeQueueRes struct {
	Message string `json:"message"`
}

// StopQueueReq 停止队列请求
type StopQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/stop" method:"post" tags:"队列控制" summary:"停止队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#停止原因最大200字符"`
	Force   bool   `json:"force,omitempty"`
}

type StopQueueRes struct {
	Message string `json:"message"`
}

// ClearQueueReq 清空队列请求
type ClearQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/clear" method:"post" tags:"队列控制" summary:"清空队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#清空原因最大200字符"`
	Force   bool   `json:"force,omitempty"`
}

type ClearQueueRes struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ResetQueueReq 重置队列请求
type ResetQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/reset" method:"post" tags:"队列控制" summary:"重置队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#重置原因最大200字符"`
}

type ResetQueueRes struct {
	Message string `json:"message"`
}

// ===============================
// 3. 队列任务管理 API (6个)
// ===============================

// EnqueueTaskReq 任务入队请求
type EnqueueTaskReq struct {
	g.Meta     `path:"/queue/{queueId}/enqueue" method:"post" tags:"队列任务" summary:"任务入队"`
	QueueId    string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskData   map[string]interface{} `json:"task_data" v:"required#任务数据不能为空"`
	Priority   int                    `json:"priority" d:"5" v:"between:1,10#任务优先级为1-10"`
	DelayTime  int64                  `json:"delay_time,omitempty" v:"min:0#延迟时间不能小于0"`
	MaxRetries int                    `json:"max_retries" d:"3" v:"between:0,10#最大重试次数为0-10"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type EnqueueTaskRes struct {
	TaskId   string `json:"task_id"`
	Position int    `json:"position"`
	Message  string `json:"message"`
}

// DequeueTaskReq 任务出队请求
type DequeueTaskReq struct {
	g.Meta  `path:"/queue/{queueId}/dequeue" method:"post" tags:"队列任务" summary:"任务出队"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Count   int    `json:"count" d:"1" v:"between:1,10#出队数量为1-10"`
	Timeout int    `json:"timeout" d:"30" v:"between:1,300#超时时间为1-300秒"`
}

type DequeueTaskRes struct {
	Tasks   []QueueTaskInfo `json:"tasks"`
	Message string          `json:"message"`
}

// GetQueueTasksReq 获取队列任务请求
type GetQueueTasksReq struct {
	g.Meta   `path:"/queue/{queueId}/tasks" method:"get" tags:"队列任务" summary:"获取队列任务"`
	QueueId  string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	Status   string `json:"status,omitempty" v:"in:pending,processing,completed,failed,retry#状态值无效"`
	Priority int    `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
}

type GetQueueTasksRes struct {
	List  []QueueTaskInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// ReorderQueueTasksReq 重新排序队列任务请求
type ReorderQueueTasksReq struct {
	g.Meta    `path:"/queue/{queueId}/reorder" method:"post" tags:"队列任务" summary:"重新排序队列任务"`
	QueueId   string   `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskIds   []string `json:"task_ids" v:"required|max:100#任务ID列表不能为空|最多100个任务"`
	Strategy  string   `json:"strategy" v:"required|in:manual,priority,fifo,lifo#排序策略无效"`
	Positions []int    `json:"positions,omitempty"`
}

type ReorderQueueTasksRes struct {
	AffectedCount int    `json:"affected_count"`
	Message       string `json:"message"`
}

// RemoveQueueTaskReq 移除队列任务请求
type RemoveQueueTaskReq struct {
	g.Meta  `path:"/queue/{queueId}/task/{taskId}/remove" method:"post" tags:"队列任务" summary:"移除队列任务"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskId  string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#移除原因最大200字符"`
}

type RemoveQueueTaskRes struct {
	Message string `json:"message"`
}

// UpdateTaskPriorityReq 更新任务优先级请求
type UpdateTaskPriorityReq struct {
	g.Meta   `path:"/queue/{queueId}/task/{taskId}/priority" method:"put" tags:"队列任务" summary:"更新任务优先级"`
	QueueId  string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskId   string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Priority int    `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type UpdateTaskPriorityRes struct {
	NewPosition int    `json:"new_position"`
	Message     string `json:"message"`
}

// ===============================
// 4. 队列监控统计 API (4个)
// ===============================

// GetQueueStatusReq 获取队列状态请求
type GetQueueStatusReq struct {
	g.Meta  `path:"/queue/{queueId}/status" method:"get" tags:"队列监控" summary:"获取队列状态"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueStatusRes struct {
	QueueStatus QueueStatusInfo `json:"queue_status"`
}

// GetQueueMetricsReq 获取队列指标请求
type GetQueueMetricsReq struct {
	g.Meta    `path:"/queue/{queueId}/metrics" method:"get" tags:"队列监控" summary:"获取队列指标"`
	QueueId   string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TimeRange string `json:"time_range" d:"1h" v:"in:5m,15m,1h,6h,24h,7d#时间范围无效"`
}

type GetQueueMetricsRes struct {
	Metrics QueueMetrics `json:"metrics"`
}

// GetQueueStatisticsReq 获取队列统计请求
type GetQueueStatisticsReq struct {
	g.Meta    `path:"/queue/statistics" method:"get" tags:"队列监控" summary:"获取队列统计"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	QueueType string `json:"queue_type,omitempty" v:"in:fifo,priority,delay,lifo#队列类型无效"`
	Status    string `json:"status,omitempty" v:"in:active,paused,stopped,error#状态值无效"`
	GroupBy   string `json:"group_by" d:"status" v:"in:status,type,priority#分组字段无效"`
}

type GetQueueStatisticsRes struct {
	Statistics QueueStatistics `json:"statistics"`
}

// GetQueuePerformanceReportReq 获取队列性能报告请求
type GetQueuePerformanceReportReq struct {
	g.Meta    `path:"/queue/{queueId}/performance" method:"get" tags:"队列监控" summary:"获取队列性能报告"`
	QueueId   string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TimeRange string `json:"time_range" d:"7d" v:"in:24h,7d,30d#时间范围无效"`
}

type GetQueuePerformanceReportRes struct {
	PerformanceReport QueuePerformanceReport `json:"performance_report"`
}

// ===============================
// 5. 队列配置管理 API (3个)
// ===============================

// GetQueueConfigReq 获取队列配置请求
type GetQueueConfigReq struct {
	g.Meta  `path:"/queue/{queueId}/config" method:"get" tags:"队列配置" summary:"获取队列配置"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueConfigRes struct {
	Config QueueConfig `json:"config"`
}

// UpdateQueueConfigReq 更新队列配置请求
type UpdateQueueConfigReq struct {
	g.Meta               `path:"/queue/{queueId}/config" method:"put" tags:"队列配置" summary:"更新队列配置"`
	QueueId              string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	ProcessingStrategy   string                 `json:"processing_strategy,omitempty" v:"in:parallel,sequential,batch#处理策略无效"`
	RetryStrategy        string                 `json:"retry_strategy,omitempty" v:"in:immediate,exponential,linear,fixed#重试策略无效"`
	RetryDelayMs         int                    `json:"retry_delay_ms,omitempty" v:"min:100#重试延迟最小100毫秒"`
	MaxRetryInterval     int                    `json:"max_retry_interval,omitempty" v:"min:1000#最大重试间隔最小1000毫秒"`
	DeadLetterEnabled    bool                   `json:"dead_letter_enabled,omitempty"`
	DeadLetterMaxRetries int                    `json:"dead_letter_max_retries,omitempty" v:"between:0,10#死信最大重试次数为0-10"`
	BatchSize            int                    `json:"batch_size,omitempty" v:"between:1,100#批量大小为1-100"`
	TimeoutMs            int                    `json:"timeout_ms,omitempty" v:"min:1000#超时时间最小1000毫秒"`
	CustomSettings       map[string]interface{} `json:"custom_settings,omitempty"`
}

type UpdateQueueConfigRes struct {
	Message string `json:"message"`
}

// GetQueueConfigHistoryReq 获取队列配置历史请求
type GetQueueConfigHistoryReq struct {
	g.Meta  `path:"/queue/{queueId}/config/history" method:"get" tags:"队列配置" summary:"获取队列配置历史"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Page    int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size    int    `json:"size" d:"10" v:"between:1,50#每页数量为1-50"`
}

type GetQueueConfigHistoryRes struct {
	List  []QueueConfigHistory `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// ===============================
// 6. 队列数据结构定义
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
