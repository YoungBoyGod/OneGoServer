package consts

// ================================
// Task 模块常量定义
// ================================

// 任务状态常量
const (
	TaskStatusPending     = "pending"
	TaskStatusQueued      = "queued"      // 已入队待分配
	TaskStatusAssigning   = "assigning"   // 分配中
	TaskStatusAssigned    = "assigned"    // 已分配待执行
	TaskStatusDispatching = "dispatching" // 派发中
	TaskStatusRunning     = "running"
	TaskStatusCompleted   = "completed"
	TaskStatusFailed      = "failed"
	TaskStatusCanceled    = "canceled"
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

// 队列状态常量
const (
	QueueStatusQueued    = "queued"
	QueueStatusAssigning = "assigning"
	QueueStatusAssigned  = "assigned"
	QueueStatusFailed    = "failed"
	QueueStatusCanceled  = "canceled"
)

// 分配历史动作常量
const (
	AssignmentActionQueued     = "queued"
	AssignmentActionAssigned   = "assigned"
	AssignmentActionReassigned = "reassigned"
	AssignmentActionFailed     = "failed"
	AssignmentActionCompleted  = "completed"
	AssignmentActionCanceled   = "canceled"
)

// 任务依赖类型常量
const (
	DependencyTypeBefore   = "before"   // 前置依赖
	DependencyTypeAfter    = "after"    // 后置依赖
	DependencyTypeParallel = "parallel" // 并行依赖
)

// 任务优先级常量
const (
	TaskPriorityLowest  = 1
	TaskPriorityLow     = 3
	TaskPriorityNormal  = 5
	TaskPriorityHigh    = 7
	TaskPriorityUrgent  = 9
	TaskPriorityHighest = 10
)

// 任务重试常量
const (
	DefaultMaxRetries = 3
	MaxAllowedRetries = 10
	DefaultTimeout    = 300 // 5分钟
)

// 任务调度常量
const (
	DefaultSchedulingInterval = 30  // 30秒
	MaxConcurrentTasks        = 100 // 最大并发任务数
	TaskBatchSize             = 50  // 批处理大小
)
