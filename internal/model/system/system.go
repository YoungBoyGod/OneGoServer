package system

// ===============================
// System模块实体定义
// ===============================

// System 系统实体
type System struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	StartTime   string `json:"start_time"`
	Uptime      string `json:"uptime"`
}

// SystemHealth 系统健康状态
type SystemHealth struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	Period  string                 `json:"period"`
	Metrics map[string]interface{} `json:"metrics"`
}

// SystemStatus 系统状态
type SystemStatus struct {
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

// SystemLog 系统日志
type SystemLog struct {
	ID        int64  `json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Module    string `json:"module"`
	TraceID   string `json:"trace_id"`
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
}

// SystemAlert 系统告警
type SystemAlert struct {
	ID             int64  `json:"id"`
	Level          string `json:"level"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	Module         string `json:"module"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	SourceID       string `json:"source_id"`
	AcknowledgedBy string `json:"acknowledged_by"`
	AcknowledgedAt string `json:"acknowledged_at"`
	ResolvedBy     string `json:"resolved_by"`
	ResolvedAt     string `json:"resolved_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// SystemPerformance 系统性能
type SystemPerformance struct {
	Period      string                 `json:"period"`
	Performance map[string]interface{} `json:"performance"`
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Status   string                 `json:"status"`
	Services map[string]interface{} `json:"services"`
}

// DatabaseStatus 数据库状态
type DatabaseStatus struct {
	Status         string `json:"status"`
	Connections    int    `json:"connections"`
	MaxConnections int    `json:"max_connections"`
}

// CacheStatus 缓存状态
type CacheStatus struct {
	Status string `json:"status"`
	Size   int64  `json:"size"`
	Hits   int64  `json:"hits"`
	Misses int64  `json:"misses"`
}

// DeviceStats 设备统计
type DeviceStats struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
	Error   int `json:"error"`
}

// TaskStats 任务统计
type TaskStats struct {
	Total     int `json:"total"`
	Running   int `json:"running"`
	Pending   int `json:"pending"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

// UserStats 用户统计
type UserStats struct {
	Total  int `json:"total"`
	Active int `json:"active"`
	Online int `json:"online"`
}
