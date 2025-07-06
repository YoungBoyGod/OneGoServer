package v1

// ===============================
// 通用数据模型定义
// ===============================

// TimePoint 时间点数据
type TimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
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
