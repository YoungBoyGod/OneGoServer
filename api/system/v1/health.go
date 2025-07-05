package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统健康检查和监控相关API
// ===============================

// HealthCheckReq 健康检查请求
type HealthCheckReq struct {
	g.Meta `path:"/health" method:"get" tags:"系统管理" summary:"系统健康检查"`
}

// HealthCheckRes 健康检查响应
type HealthCheckRes struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
}

// GetSystemMetricsReq 获取系统指标请求
type GetSystemMetricsReq struct {
	g.Meta `path:"/metrics" method:"get" tags:"系统管理" summary:"获取系统指标"`
	Period string `json:"period" d:"1h" v:"in:1h,6h,24h,7d#时间周期只能是1h,6h,24h,7d"`
}

// GetSystemMetricsRes 获取系统指标响应
type GetSystemMetricsRes struct {
	Period  string `json:"period"`
	Metrics struct {
		CpuUsage          []TimePoint `json:"cpuUsage"`
		MemoryUsage       []TimePoint `json:"memoryUsage"`
		DiskUsage         []TimePoint `json:"diskUsage"`
		NetworkIO         []TimePoint `json:"networkIO"`
		ActiveConnections []TimePoint `json:"activeConnections"`
		RequestRate       []TimePoint `json:"requestRate"`
		ErrorRate         []TimePoint `json:"errorRate"`
		ResponseTime      []TimePoint `json:"responseTime"`
	} `json:"metrics"`
}

// TimePoint 时间点数据
type TimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// GetSystemStatusReq 获取系统状态请求
type GetSystemStatusReq struct {
	g.Meta `path:"/status" method:"get" tags:"系统管理" summary:"获取系统状态"`
}

// GetSystemStatusRes 获取系统状态响应
type GetSystemStatusRes struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Uptime      string `json:"uptime"`
	StartTime   string `json:"startTime"`
	Environment string `json:"environment"`
	Database    struct {
		Status         string `json:"status"`
		Connections    int    `json:"connections"`
		MaxConnections int    `json:"maxConnections"`
	} `json:"database"`
	Cache struct {
		Status string `json:"status"`
		Size   int64  `json:"size"`
		Hits   int64  `json:"hits"`
		Misses int64  `json:"misses"`
	} `json:"cache"`
	Devices struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
		Error   int `json:"error"`
	} `json:"devices"`
	Tasks struct {
		Total     int `json:"total"`
		Running   int `json:"running"`
		Pending   int `json:"pending"`
		Completed int `json:"completed"`
		Failed    int `json:"failed"`
	} `json:"tasks"`
	Users struct {
		Total  int `json:"total"`
		Active int `json:"active"`
		Online int `json:"online"`
	} `json:"users"`
}

// GetSystemLogsReq 获取系统日志请求
type GetSystemLogsReq struct {
	g.Meta    `path:"/logs" method:"get" tags:"系统管理" summary:"获取系统日志"`
	Level     string `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别只能是debug,info,warn,error"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"100" v:"between:1,1000#每页数量为1-1000"`
}

// GetSystemLogsRes 获取系统日志响应
type GetSystemLogsRes struct {
	List  []SystemLog `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SystemLog 系统日志
type SystemLog struct {
	Id        int64  `json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Module    string `json:"module"`
	TraceId   string `json:"traceId"`
	UserId    string `json:"userId,omitempty"`
	IpAddress string `json:"ipAddress,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	CreatedAt string `json:"createdAt"`
}

// GetSystemAlertsReq 获取系统告警请求
type GetSystemAlertsReq struct {
	g.Meta    `path:"/alerts" method:"get" tags:"系统管理" summary:"获取系统告警"`
	Level     string `json:"level,omitempty" v:"in:info,warn,error,critical#告警级别只能是info,warn,error,critical"`
	Status    string `json:"status,omitempty" v:"in:active,resolved,acknowledged#告警状态只能是active,resolved,acknowledged"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
}

// GetSystemAlertsRes 获取系统告警响应
type GetSystemAlertsRes struct {
	List  []SystemAlert `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// SystemAlert 系统告警
type SystemAlert struct {
	Id             int64  `json:"id"`
	Level          string `json:"level"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	Module         string `json:"module"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	SourceId       string `json:"sourceId"`
	AcknowledgedBy string `json:"acknowledgedBy,omitempty"`
	AcknowledgedAt string `json:"acknowledgedAt,omitempty"`
	ResolvedBy     string `json:"resolvedBy,omitempty"`
	ResolvedAt     string `json:"resolvedAt,omitempty"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

// AcknowledgeAlertReq 确认告警请求
type AcknowledgeAlertReq struct {
	g.Meta  `path:"/alerts/acknowledge" method:"post" tags:"系统管理" summary:"确认告警"`
	AlertId string `json:"alertId" v:"required#告警ID不能为空"`
	Comment string `json:"comment,omitempty"`
}

// AcknowledgeAlertRes 确认告警响应
type AcknowledgeAlertRes struct {
	AlertId string `json:"alertId"`
	Status  string `json:"status"`
}

// ResolveAlertReq 解决告警请求
type ResolveAlertReq struct {
	g.Meta  `path:"/alerts/resolve" method:"post" tags:"系统管理" summary:"解决告警"`
	AlertId string `json:"alertId" v:"required#告警ID不能为空"`
	Comment string `json:"comment,omitempty"`
}

// ResolveAlertRes 解决告警响应
type ResolveAlertRes struct {
	AlertId string `json:"alertId"`
	Status  string `json:"status"`
}

// GetSystemPerformanceReq 获取系统性能请求
type GetSystemPerformanceReq struct {
	g.Meta `path:"/performance" method:"get" tags:"系统管理" summary:"获取系统性能"`
	Period string `json:"period" d:"1h" v:"in:1h,6h,24h,7d#时间周期只能是1h,6h,24h,7d"`
}

// GetSystemPerformanceRes 获取系统性能响应
type GetSystemPerformanceRes struct {
	Period      string `json:"period"`
	Performance struct {
		ApiLatency struct {
			Avg float64 `json:"avg"`
			P95 float64 `json:"p95"`
			P99 float64 `json:"p99"`
			Max float64 `json:"max"`
		} `json:"apiLatency"`
		Throughput struct {
			RequestsPerSecond float64 `json:"requestsPerSecond"`
			TotalRequests     int64   `json:"totalRequests"`
		} `json:"throughput"`
		ErrorRate struct {
			Rate        float64 `json:"rate"`
			TotalErrors int64   `json:"totalErrors"`
		} `json:"errorRate"`
		ResourceUsage struct {
			CpuUsage    float64 `json:"cpuUsage"`
			MemoryUsage float64 `json:"memoryUsage"`
			DiskUsage   float64 `json:"diskUsage"`
		} `json:"resourceUsage"`
	} `json:"performance"`
}
