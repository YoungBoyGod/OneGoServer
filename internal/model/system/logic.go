package system

// ===============================
// System模块逻辑层Input/Output结构体定义
// ===============================

// HealthCheckInput 健康检查输入
type HealthCheckInput struct {
	// 可以添加健康检查的特定参数
}

// HealthCheckOutput 健康检查输出
type HealthCheckOutput struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
}

// GetSystemMetricsInput 获取系统指标输入
type GetSystemMetricsInput struct {
	Period string `json:"period"`
}

// GetSystemMetricsOutput 获取系统指标输出
type GetSystemMetricsOutput struct {
	Period  string                 `json:"period"`
	Metrics map[string]interface{} `json:"metrics"`
}

// GetSystemStatusInput 获取系统状态输入
type GetSystemStatusInput struct {
	// 可以添加状态查询的特定参数
}

// GetSystemStatusOutput 获取系统状态输出
type GetSystemStatusOutput struct {
	Status      string                 `json:"status"`
	Version     string                 `json:"version"`
	Uptime      string                 `json:"uptime"`
	StartTime   string                 `json:"start_time"`
	Environment string                 `json:"environment"`
	Database    map[string]interface{} `json:"database"`
	Cache       map[string]interface{} `json:"cache"`
	Devices     map[string]interface{} `json:"devices"`
	Tasks       map[string]interface{} `json:"tasks"`
	Users       map[string]interface{} `json:"users"`
}

// GetSystemLogsInput 获取系统日志输入
type GetSystemLogsInput struct {
	Level     string `json:"level"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

// GetSystemLogsOutput 获取系统日志输出
type GetSystemLogsOutput struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

// GetSystemAlertsInput 获取系统告警输入
type GetSystemAlertsInput struct {
	Level     string `json:"level"`
	Status    string `json:"status"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

// GetSystemAlertsOutput 获取系统告警输出
type GetSystemAlertsOutput struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

// AcknowledgeAlertInput 确认告警输入
type AcknowledgeAlertInput struct {
	AlertID string `json:"alert_id"`
	Comment string `json:"comment"`
}

// AcknowledgeAlertOutput 确认告警输出
type AcknowledgeAlertOutput struct {
	AlertID string `json:"alert_id"`
	Status  string `json:"status"`
}

// ResolveAlertInput 解决告警输入
type ResolveAlertInput struct {
	AlertID string `json:"alert_id"`
	Comment string `json:"comment"`
}

// ResolveAlertOutput 解决告警输出
type ResolveAlertOutput struct {
	AlertID string `json:"alert_id"`
	Status  string `json:"status"`
}

// GetSystemPerformanceInput 获取系统性能输入
type GetSystemPerformanceInput struct {
	Period string `json:"period"`
}

// GetSystemPerformanceOutput 获取系统性能输出
type GetSystemPerformanceOutput struct {
	Period      string                 `json:"period"`
	Performance map[string]interface{} `json:"performance"`
}

// 辅助方法相关结构体

// ParsePeriodInput 解析时间周期输入
type ParsePeriodInput struct {
	Period string `json:"period"`
}

// ParsePeriodOutput 解析时间周期输出
type ParsePeriodOutput struct {
	Duration string `json:"duration"`
}

// GenerateTimeSeriesDataInput 生成时间序列数据输入
type GenerateTimeSeriesDataInput struct {
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	IntervalSeconds int     `json:"interval_seconds"`
	MinValue        float64 `json:"min_value"`
	MaxValue        float64 `json:"max_value"`
}

// GenerateTimeSeriesDataOutput 生成时间序列数据输出
type GenerateTimeSeriesDataOutput struct {
	Data []map[string]interface{} `json:"data"`
}

// CalculateSystemPerformanceInput 计算系统性能输入
type CalculateSystemPerformanceInput struct {
	Duration string `json:"duration"`
}

// CalculateSystemPerformanceOutput 计算系统性能输出
type CalculateSystemPerformanceOutput struct {
	Performance map[string]interface{} `json:"performance"`
}

// GenerateSystemLogsInput 生成系统日志输入
type GenerateSystemLogsInput struct {
	Level string `json:"level"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// GenerateSystemLogsOutput 生成系统日志输出
type GenerateSystemLogsOutput struct {
	Logs []map[string]interface{} `json:"logs"`
}

// GenerateSystemAlertsInput 生成系统告警输入
type GenerateSystemAlertsInput struct {
	Level  string `json:"level"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

// GenerateSystemAlertsOutput 生成系统告警输出
type GenerateSystemAlertsOutput struct {
	Alerts []map[string]interface{} `json:"alerts"`
}

// GetSystemInfoInput 获取系统信息输入
type GetSystemInfoInput struct {
	// 可以添加系统信息查询的特定参数
}

// GetSystemInfoOutput 获取系统信息输出
type GetSystemInfoOutput struct {
	Version     string `json:"version"`
	Uptime      string `json:"uptime"`
	StartTime   string `json:"start_time"`
	Environment string `json:"environment"`
}

// GetResourceUsageInput 获取资源使用情况输入
type GetResourceUsageInput struct {
	// 可以添加资源使用查询的特定参数
}

// GetResourceUsageOutput 获取资源使用情况输出
type GetResourceUsageOutput struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// CheckServiceStatusInput 检查服务状态输入
type CheckServiceStatusInput struct {
	// 可以添加服务状态检查的特定参数
}

// CheckServiceStatusOutput 检查服务状态输出
type CheckServiceStatusOutput struct {
	Status   string                 `json:"status"`
	Services map[string]interface{} `json:"services"`
}

// CheckDatabaseStatusInput 检查数据库状态输入
type CheckDatabaseStatusInput struct {
	// 可以添加数据库状态检查的特定参数
}

// CheckDatabaseStatusOutput 检查数据库状态输出
type CheckDatabaseStatusOutput struct {
	Status         string `json:"status"`
	Connections    int    `json:"connections"`
	MaxConnections int    `json:"max_connections"`
}

// CheckCacheStatusInput 检查缓存状态输入
type CheckCacheStatusInput struct {
	// 可以添加缓存状态检查的特定参数
}

// CheckCacheStatusOutput 检查缓存状态输出
type CheckCacheStatusOutput struct {
	Status string `json:"status"`
	Size   int64  `json:"size"`
	Hits   int64  `json:"hits"`
	Misses int64  `json:"misses"`
}

// GetDeviceStatsInput 获取设备统计输入
type GetDeviceStatsInput struct {
	// 可以添加设备统计查询的特定参数
}

// GetDeviceStatsOutput 获取设备统计输出
type GetDeviceStatsOutput struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
	Error   int `json:"error"`
}

// GetTaskStatsInput 获取任务统计输入
type GetTaskStatsInput struct {
	// 可以添加任务统计查询的特定参数
}

// GetTaskStatsOutput 获取任务统计输出
type GetTaskStatsOutput struct {
	Total     int `json:"total"`
	Running   int `json:"running"`
	Pending   int `json:"pending"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

// GetUserStatsInput 获取用户统计输入
type GetUserStatsInput struct {
	// 可以添加用户统计查询的特定参数
}

// GetUserStatsOutput 获取用户统计输出
type GetUserStatsOutput struct {
	Total  int `json:"total"`
	Active int `json:"active"`
	Online int `json:"online"`
}
