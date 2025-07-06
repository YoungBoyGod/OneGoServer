package consts

// ================================
// Queue 模块常量定义
// ================================

// 队列状态常量
const (
	QueueStatusActive  = "active"
	QueueStatusPaused  = "paused"
	QueueStatusStopped = "stopped"
	QueueStatusError   = "error"
)

// 队列类型常量
const (
	QueueTypeFIFO     = "fifo"
	QueueTypePriority = "priority"
	QueueTypeDelay    = "delay"
	QueueTypeLIFO     = "lifo"
)

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

// 队列优先级常量
const (
	QueuePriorityLowest  = 1
	QueuePriorityLow     = 3
	QueuePriorityNormal  = 5
	QueuePriorityHigh    = 7
	QueuePriorityUrgent  = 9
	QueuePriorityHighest = 10
)

// 队列配置常量
const (
	QueueMaxSizeDefault    = 100
	QueueMaxSizeLimit      = 1000
	QueueConcurrentDefault = 1
	QueueConcurrentLimit   = 10
	QueueBatchSizeDefault  = 10
	QueueBatchSizeLimit    = 100
)

// 队列超时常量
const (
	QueueTimeoutDefault = 3600  // 1小时
	QueueTimeoutMin     = 60    // 1分钟
	QueueTimeoutMax     = 86400 // 24小时
)

// 队列重试常量
const (
	QueueRetryDefault = 3
	QueueRetryMax     = 10
	QueueRetryMin     = 0
)

// 队列监控常量
const (
	QueueMonitorInterval     = 30  // 监控间隔30秒
	QueueHealthCheckInterval = 60  // 健康检查间隔1分钟
	QueueMetricsInterval     = 300 // 指标收集间隔5分钟
)
