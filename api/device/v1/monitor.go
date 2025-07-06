package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备监控管理 API
// ==============================================

// GetDeviceLoadStatus 获取设备负载状态请求
type GetDeviceLoadStatusReq struct {
	g.Meta   `path:"/device/{deviceId}/load" method:"get" tags:"设备监控" summary:"获取设备负载状态"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceLoadStatusRes 获取设备负载状态响应
type GetDeviceLoadStatusRes struct {
	DeviceId           int64       `json:"device_id"`
	CurrentTasks       int         `json:"current_tasks"`
	MaxConcurrentTasks int         `json:"max_concurrent_tasks"`
	CpuLoad            float64     `json:"cpu_load"`
	MemoryUsage        float64     `json:"memory_usage"`
	DiskUsage          float64     `json:"disk_usage"`
	NetworkLatency     int         `json:"network_latency"`
	Status             string      `json:"status"`
	LastHeartbeat      *gtime.Time `json:"last_heartbeat"`
	LoadScore          float64     `json:"load_score"`
	TotalAssigned      int         `json:"total_assigned"`
	TotalCompleted     int         `json:"total_completed"`
	TotalFailed        int         `json:"total_failed"`
	SuccessRate        float64     `json:"success_rate"`
	AvgTaskDuration    float64     `json:"avg_task_duration"`
	LastTaskCompletion *gtime.Time `json:"last_task_completion"`
	UpdatedAt          *gtime.Time `json:"updated_at"`
}

// DeviceLoadHistoryInfo 设备负载历史信息
type DeviceLoadHistoryInfo struct {
	DeviceId        int64       `json:"device_id"`
	Timestamp       *gtime.Time `json:"timestamp"`
	CpuLoad         float64     `json:"cpu_load"`
	MemoryUsage     float64     `json:"memory_usage"`
	DiskUsage       float64     `json:"disk_usage"`
	NetworkLatency  int         `json:"network_latency"`
	CurrentTasks    int         `json:"current_tasks"`
	LoadScore       float64     `json:"load_score"`
	SuccessRate     float64     `json:"success_rate"`
	AvgTaskDuration float64     `json:"avg_task_duration"`
}

// UpdateDeviceLoadConfig 更新设备负载配置请求
type UpdateDeviceLoadConfigReq struct {
	g.Meta             `path:"/device/{deviceId}/load/config" method:"put" tags:"设备监控" summary:"更新设备负载配置"`
	DeviceId           string  `json:"deviceId" v:"required#设备ID不能为空"`
	MaxConcurrentTasks int     `json:"max_concurrent_tasks,omitempty" v:"min:1#最大并发任务数最小为1"`
	CpuThreshold       float64 `json:"cpu_threshold,omitempty" v:"between:0,100#CPU阈值范围为0-100"`
	MemoryThreshold    float64 `json:"memory_threshold,omitempty" v:"between:0,100#内存阈值范围为0-100"`
	DiskThreshold      float64 `json:"disk_threshold,omitempty" v:"between:0,100#磁盘阈值范围为0-100"`
	NetworkThreshold   int     `json:"network_threshold,omitempty" v:"min:0#网络延迟阈值不能为负数"`
}

// UpdateDeviceLoadConfigRes 更新设备负载配置响应
type UpdateDeviceLoadConfigRes struct {
	Message string `json:"message"`
}

// CalculateDeviceLoadScore 计算设备负载评分请求
type CalculateDeviceLoadScoreReq struct {
	g.Meta   `path:"/device/load/calculate" method:"post" tags:"设备管理" summary:"计算设备负载评分"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// CalculateDeviceLoadScoreRes 计算设备负载评分响应
type CalculateDeviceLoadScoreRes struct {
	DeviceId   string  `json:"deviceId"`
	LoadScore  float64 `json:"loadScore"`
	Components struct {
		CpuScore     float64 `json:"cpuScore"`
		MemoryScore  float64 `json:"memoryScore"`
		DiskScore    float64 `json:"diskScore"`
		NetworkScore float64 `json:"networkScore"`
		TaskScore    float64 `json:"taskScore"`
	} `json:"components"`
	Recommendations []string `json:"recommendations"`
}

// GetDeviceLoadMetrics 获取设备负载指标请求
type GetDeviceLoadMetricsReq struct {
	g.Meta   `path:"/device/load/metrics" method:"get" tags:"设备管理" summary:"获取设备负载指标"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Period   string `json:"period" d:"1h" v:"in:1h,6h,24h,7d#时间周期只能是1h,6h,24h,7d"`
}

// GetDeviceLoadMetricsRes 获取设备负载指标响应
type GetDeviceLoadMetricsRes struct {
	DeviceId string `json:"deviceId"`
	Period   string `json:"period"`
	Metrics  struct {
		CpuUsage       []TimePoint `json:"cpuUsage"`
		MemoryUsage    []TimePoint `json:"memoryUsage"`
		DiskUsage      []TimePoint `json:"diskUsage"`
		NetworkLatency []TimePoint `json:"networkLatency"`
		CurrentTasks   []TimePoint `json:"currentTasks"`
		LoadScore      []TimePoint `json:"loadScore"`
	} `json:"metrics"`
}

// OptimizeDeviceLoad 优化设备负载请求
type OptimizeDeviceLoadReq struct {
	g.Meta   `path:"/device/load/optimize" method:"post" tags:"设备管理" summary:"优化设备负载"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Strategy string `json:"strategy" v:"in:auto,manual,conservative,aggressive#优化策略只能是auto,manual,conservative,aggressive"`
}

// OptimizeDeviceLoadRes 优化设备负载响应
type OptimizeDeviceLoadRes struct {
	DeviceId       string   `json:"deviceId"`
	OriginalScore  float64  `json:"originalScore"`
	OptimizedScore float64  `json:"optimizedScore"`
	Improvement    float64  `json:"improvement"`
	Actions        []string `json:"actions"`
	Status         string   `json:"status"`
}

// GetDeviceLoadHistory 获取设备负载历史请求
type GetDeviceLoadHistoryReq struct {
	g.Meta                   `path:"/device/load/history" method:"get" tags:"设备管理" summary:"获取设备负载历史"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	StartTime                string `json:"startTime" v:"required#开始时间不能为空"`
	EndTime                  string `json:"endTime" v:"required#结束时间不能为空"`
	common.PaginationRequest `json:",inline"`
}

// GetDeviceLoadHistoryRes 获取设备负载历史响应
type GetDeviceLoadHistoryRes struct {
	common.PaginationResponse[DeviceLoadRecord] `json:",inline"`
}

// SetDeviceLoadThreshold 设置设备负载阈值请求
type SetDeviceLoadThresholdReq struct {
	g.Meta   `path:"/device/load/threshold" method:"post" tags:"设备管理" summary:"设置设备负载阈值"`
	DeviceId string  `json:"deviceId" v:"required#设备ID不能为空"`
	Warning  float64 `json:"warning" v:"between:0,100#警告阈值必须在0-100之间"`
	Critical float64 `json:"critical" v:"between:0,100#严重阈值必须在0-100之间"`
	MaxTasks int     `json:"maxTasks" v:"min:1#最大任务数最小为1"`
}

// SetDeviceLoadThresholdRes 设置设备负载阈值响应
type SetDeviceLoadThresholdRes struct {
	DeviceId string  `json:"deviceId"`
	Warning  float64 `json:"warning"`
	Critical float64 `json:"critical"`
	MaxTasks int     `json:"maxTasks"`
	Status   string  `json:"status"`
}

// GetDeviceLoadThreshold 获取设备负载阈值请求
type GetDeviceLoadThresholdReq struct {
	g.Meta   `path:"/device/load/threshold" method:"get" tags:"设备管理" summary:"获取设备负载阈值"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceLoadThresholdRes 获取设备负载阈值响应
type GetDeviceLoadThresholdRes struct {
	DeviceId    string  `json:"deviceId"`
	Warning     float64 `json:"warning"`
	Critical    float64 `json:"critical"`
	MaxTasks    int     `json:"maxTasks"`
	CurrentLoad float64 `json:"currentLoad"`
	Status      string  `json:"status"`
}
