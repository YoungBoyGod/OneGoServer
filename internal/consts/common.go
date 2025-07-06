package consts

import "time"

// ================================
// 通用常量定义
// ================================

// 系统常量
const (
	SystemName    = "OneGoServer"
	SystemVersion = "v2.0.0"
	APIVersion    = "v1"
)

// HTTP状态码相关常量
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusPending = "pending"
)

// 通用状态常量
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusRunning   = "running"
	StatusStopped   = "stopped"
	StatusPaused    = "paused"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCanceled  = "canceled"
	StatusLocked    = "locked"
)

// 通用排序常量
const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// 通用操作类型常量
const (
	ActionTypeCreate   = "create"
	ActionTypeUpdate   = "update"
	ActionTypeDelete   = "delete"
	ActionTypeView     = "view"
	ActionTypeStart    = "start"
	ActionTypeStop     = "stop"
	ActionTypePause    = "pause"
	ActionTypeResume   = "resume"
	ActionTypeCancel   = "cancel"
	ActionTypeRestart  = "restart"
	ActionTypeRetry    = "retry"
	ActionTypeAssign   = "assign"
	ActionTypeUnassign = "unassign"
)

// 通用优先级常量
const (
	PriorityLowest  = 1
	PriorityLow     = 3
	PriorityNormal  = 5
	PriorityHigh    = 7
	PriorityUrgent  = 9
	PriorityHighest = 10
)

// 分页常量
const (
	DefaultPageSize = 20
	DefaultPage     = 1
	MaxPageSize     = 100
)

// 时间相关常量
const (
	TimeFormatRFC3339 = time.RFC3339
	TimeFormatDate    = "2006-01-02"
	TimeFormatTime    = "15:04:05"
	TimeFormatFull    = "2006-01-02 15:04:05"
)

// 时间间隔常量（秒）
const (
	SecondMinute = 60
	SecondHour   = 3600
	SecondDay    = 86400
	SecondWeek   = 604800
)

// 文件大小常量（字节）
const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
)

// 数据库相关常量
const (
	DBDefaultTimeout  = 30 * time.Second
	DBMaxConnections  = 100
	DBMaxIdleConns    = 10
	DBConnMaxLifetime = time.Hour
)

// 缓存相关常量
const (
	CacheDefaultTTL    = 300  // 5分钟
	CacheLongTTL       = 3600 // 1小时
	CacheShortTTL      = 60   // 1分钟
	CacheMaxRetries    = 3
	CacheRetryInterval = time.Second
)

// 队列相关常量
const (
	QueueDefaultCapacity = 1000
	QueueMaxCapacity     = 10000
	QueueBatchSize       = 100
	QueueProcessInterval = 5 * time.Second
)

// 网络相关常量
const (
	NetworkTimeoutDefault = 30 * time.Second
	NetworkRetryMax       = 3
	NetworkRetryInterval  = 2 * time.Second
)

// 日志相关常量
const (
	LogMaxSize    = 100 // MB
	LogMaxBackups = 30  // 个
	LogMaxAge     = 30  // 天
	LogCompress   = true
)

// 验证相关常量
const (
	ValidationMaxStringLength = 255
	ValidationMaxTextLength   = 10000
	ValidationMinPassword     = 8
	ValidationMaxPassword     = 128
)

// 错误相关常量
const (
	ErrorCodeSuccess        = 0
	ErrorCodeBadRequest     = 400
	ErrorCodeUnauthorized   = 401
	ErrorCodeForbidden      = 403
	ErrorCodeNotFound       = 404
	ErrorCodeConflict       = 409
	ErrorCodeInternalError  = 500
	ErrorCodeServiceUnavail = 503
)

// ID生成相关常量
const (
	IDPrefixTask      = "task_"
	IDPrefixDevice    = "dev_"
	IDPrefixExecution = "exec_"
	IDPrefixCommand   = "cmd_"
	IDLength          = 20
)

// 正则表达式常量
const (
	RegexEmail    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	RegexPhone    = `^1[3-9]\d{9}$`
	RegexIPv4     = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`
	RegexDeviceID = `^[a-zA-Z0-9_-]+$`
)

// 环境相关常量
const (
	EnvDevelopment = "development"
	EnvTesting     = "testing"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// 日志级别常量
const (
	LogLevelDebug    = "debug"
	LogLevelInfo     = "info"
	LogLevelWarn     = "warn"
	LogLevelError    = "error"
	LogLevelCritical = "critical"
)
