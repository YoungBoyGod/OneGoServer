package v1

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	Id                  int64       `json:"id"`
	DeviceId            string      `json:"device_id"`
	Name                string      `json:"name"`
	Type                string      `json:"type"`
	Model               string      `json:"model"`
	BoardId             string      `json:"board_id"`
	Status              string      `json:"status"`
	HealthScore         int         `json:"health_score"`
	IpAddress           string      `json:"ip_address"`
	Port                int         `json:"port"`
	Protocol            string      `json:"protocol"`
	Endpoint            string      `json:"endpoint"`
	RegTime             *gtime.Time `json:"reg_time"`
	LastOnlineTime      *gtime.Time `json:"last_online_time"`
	LastOfflineTime     *gtime.Time `json:"last_offline_time"`
	TotalOnlineDuration int64       `json:"total_online_duration"`
	TotalHeartbeats     int64       `json:"total_heartbeats"`
	TotalAlerts         int64       `json:"total_alerts"`
	TotalTasks          int64       `json:"total_tasks"`
	CreatedAt           *gtime.Time `json:"created_at"`
	UpdatedAt           *gtime.Time `json:"updated_at"`
}

// DeviceHeartbeatInfo 设备心跳信息
type DeviceHeartbeatInfo struct {
	DeviceId      string      `json:"device_id"`
	HeartbeatTime *gtime.Time `json:"heartbeat_time"`
	Status        string      `json:"status"`
	IpAddress     string      `json:"ip_address"`
	ResponseTime  int         `json:"response_time"`
	Metadata      string      `json:"metadata"`
}

// DeviceTaskInfo 设备任务信息
type DeviceTaskInfo struct {
	Id        int64       `json:"id"`
	DeviceId  string      `json:"device_id"`
	TaskId    string      `json:"task_id"`
	TaskName  string      `json:"task_name"`
	TaskType  string      `json:"task_type"`
	Status    string      `json:"status"`
	Priority  int         `json:"priority"`
	CreatedAt *gtime.Time `json:"created_at"`
	StartTime *gtime.Time `json:"start_time"`
	EndTime   *gtime.Time `json:"end_time"`
}

// DeviceTaskQueueInfo 设备任务队列信息
type DeviceTaskQueueInfo struct {
	Id         int64       `json:"id"`
	DeviceId   string      `json:"device_id"`
	TaskId     string      `json:"task_id"`
	Priority   int         `json:"priority"`
	Status     string      `json:"status"`
	RetryCount int         `json:"retry_count"`
	CreatedAt  *gtime.Time `json:"created_at"`
	UpdatedAt  *gtime.Time `json:"updated_at"`
}

// DeviceAlertInfo 设备告警信息
type DeviceAlertInfo struct {
	Id          int64       `json:"id"`
	DeviceId    string      `json:"device_id"`
	AlertType   string      `json:"alert_type"`
	AlertLevel  string      `json:"alert_level"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	CreatedAt   *gtime.Time `json:"created_at"`
	ResolvedAt  *gtime.Time `json:"resolved_at"`
}

// DeviceLogInfo 设备日志信息
type DeviceLogInfo struct {
	Id            int64       `json:"id"`
	DeviceId      string      `json:"device_id"`
	LogTime       *gtime.Time `json:"log_time"`
	Level         string      `json:"level"`
	Category      string      `json:"category"`
	Message       string      `json:"message"`
	Details       string      `json:"details"`
	Source        string      `json:"source"`
	CorrelationId string      `json:"correlation_id"`
}

// DeviceQueueHistoryInfo 设备队列历史信息
type DeviceQueueHistoryInfo struct {
	Id               int64       `json:"id"`
	DeviceId         int64       `json:"device_id"`
	TaskId           string      `json:"task_id"`
	OperationType    string      `json:"operation_type"`
	OperationBy      int64       `json:"operation_by"`
	OperationTime    *gtime.Time `json:"operation_time"`
	OldPriority      int         `json:"old_priority"`
	NewPriority      int         `json:"new_priority"`
	OldPosition      int         `json:"old_position"`
	NewPosition      int         `json:"new_position"`
	OldStatus        string      `json:"old_status"`
	NewStatus        string      `json:"new_status"`
	Reason           string      `json:"reason"`
	Notes            string      `json:"notes"`
	OperationSource  string      `json:"operation_source"`
	BatchId          string      `json:"batch_id"`
	IsBatchOperation bool        `json:"is_batch_operation"`
}

// DeviceCommandInfo 设备命令信息
type DeviceCommandInfo struct {
	Id            int64       `json:"id"`
	DeviceId      string      `json:"device_id"`
	CommandType   string      `json:"command_type"`
	CommandData   string      `json:"command_data"`
	Status        string      `json:"status"`
	SentTime      *gtime.Time `json:"sent_time"`
	ExecutedTime  *gtime.Time `json:"executed_time"`
	CompletedTime *gtime.Time `json:"completed_time"`
	ResponseData  string      `json:"response_data"`
	ErrorMessage  string      `json:"error_message"`
	Duration      int64       `json:"duration"`
	CreatedAt     *gtime.Time `json:"created_at"`
	CreatedBy     int64       `json:"created_by"`
}

// DeviceConfigHistoryInfo 设备配置历史信息
type DeviceConfigHistoryInfo struct {
	Id             int64                  `json:"id"`
	DeviceId       string                 `json:"device_id"`
	Version        string                 `json:"version"`
	Configurations map[string]interface{} `json:"configurations"`
	ChangeType     string                 `json:"change_type"` // create, update, delete
	Reason         string                 `json:"reason"`
	ApplyStatus    string                 `json:"apply_status"` // pending, applied, failed
	CreatedAt      *gtime.Time            `json:"created_at"`
	UpdatedBy      string                 `json:"updated_by"`
}

// DeviceStatisticsTrend 设备统计趋势
type DeviceStatisticsTrend struct {
	Timestamp      *gtime.Time      `json:"timestamp"`
	OnlineCount    int64            `json:"online_count"`
	OfflineCount   int64            `json:"offline_count"`
	NewDevices     int64            `json:"new_devices"`
	TotalTasks     int64            `json:"total_tasks"`
	CompletedTasks int64            `json:"completed_tasks"`
	FailedTasks    int64            `json:"failed_tasks"`
	ByType         map[string]int64 `json:"by_type"`
}

// DeviceActivityInfo 设备活动信息
type DeviceActivityInfo struct {
	DeviceId        string      `json:"device_id"`
	DeviceName      string      `json:"device_name"`
	DeviceType      string      `json:"device_type"`
	TaskCount       int64       `json:"task_count"`
	CompletedTasks  int64       `json:"completed_tasks"`
	FailedTasks     int64       `json:"failed_tasks"`
	SuccessRate     float64     `json:"success_rate"`
	AvgTaskDuration float64     `json:"avg_task_duration"`
	UptimeHours     float64     `json:"uptime_hours"`
	HealthScore     int         `json:"health_score"`
	LastActiveTime  *gtime.Time `json:"last_active_time"`
}

// DeviceTaskMetrics 设备任务指标
type DeviceTaskMetrics struct {
	TotalTasks       int64   `json:"total_tasks"`
	CompletedTasks   int64   `json:"completed_tasks"`
	FailedTasks      int64   `json:"failed_tasks"`
	CanceledTasks    int64   `json:"canceled_tasks"`
	SuccessRate      float64 `json:"success_rate"`
	AvgExecutionTime float64 `json:"avg_execution_time"`
	MaxExecutionTime float64 `json:"max_execution_time"`
	MinExecutionTime float64 `json:"min_execution_time"`
	TasksPerHour     float64 `json:"tasks_per_hour"`
}

// DevicePerformanceMetrics 设备性能指标
type DevicePerformanceMetrics struct {
	AvgCpuUsage       float64 `json:"avg_cpu_usage"`
	PeakCpuUsage      float64 `json:"peak_cpu_usage"`
	AvgMemoryUsage    float64 `json:"avg_memory_usage"`
	PeakMemoryUsage   float64 `json:"peak_memory_usage"`
	AvgDiskUsage      float64 `json:"avg_disk_usage"`
	AvgNetworkLatency float64 `json:"avg_network_latency"`
	MaxNetworkLatency float64 `json:"max_network_latency"`
	AvgResponseTime   float64 `json:"avg_response_time"`
}

// DeviceAvailabilityMetrics 设备可用性指标
type DeviceAvailabilityMetrics struct {
	TotalHours        float64 `json:"total_hours"`
	OnlineHours       float64 `json:"online_hours"`
	OfflineHours      float64 `json:"offline_hours"`
	MaintenanceHours  float64 `json:"maintenance_hours"`
	UptimePercentage  float64 `json:"uptime_percentage"`
	MTBF              float64 `json:"mtbf"` // Mean Time Between Failures
	MTTR              float64 `json:"mttr"` // Mean Time To Recovery
	AvailabilityScore float64 `json:"availability_score"`
}

// PerformanceTrendPoint 性能趋势点
type PerformanceTrendPoint struct {
	Timestamp      *gtime.Time `json:"timestamp"`
	CpuUsage       float64     `json:"cpu_usage"`
	MemoryUsage    float64     `json:"memory_usage"`
	DiskUsage      float64     `json:"disk_usage"`
	NetworkLatency float64     `json:"network_latency"`
	TaskCount      int64       `json:"task_count"`
	SuccessRate    float64     `json:"success_rate"`
	ResponseTime   float64     `json:"response_time"`
	HealthScore    float64     `json:"health_score"`
}

// TimePoint 时间点
type TimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// DeviceLoadRecord 设备负载记录
type DeviceLoadRecord struct {
	Id             int64   `json:"id"`
	DeviceId       string  `json:"deviceId"`
	LoadScore      float64 `json:"loadScore"`
	CpuUsage       float64 `json:"cpuUsage"`
	MemoryUsage    float64 `json:"memoryUsage"`
	DiskUsage      float64 `json:"diskUsage"`
	NetworkLatency float64 `json:"networkLatency"`
	CurrentTasks   int     `json:"currentTasks"`
	MaxTasks       int     `json:"maxTasks"`
	RecordedAt     string  `json:"recordedAt"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	DeviceId string `json:"device_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}
