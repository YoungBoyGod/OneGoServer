package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 设备负载评分相关API
// ===============================

// CalculateDeviceLoadScoreReq 计算设备负载评分请求
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

// GetDeviceLoadMetricsReq 获取设备负载指标请求
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

// TimePoint 时间点数据
type TimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// OptimizeDeviceLoadReq 优化设备负载请求
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

// GetDeviceLoadHistoryReq 获取设备负载历史请求
type GetDeviceLoadHistoryReq struct {
	g.Meta    `path:"/device/load/history" method:"get" tags:"设备管理" summary:"获取设备负载历史"`
	DeviceId  string `json:"deviceId" v:"required#设备ID不能为空"`
	StartTime string `json:"startTime" v:"required#开始时间不能为空"`
	EndTime   string `json:"endTime" v:"required#结束时间不能为空"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"100" v:"between:1,1000#每页数量为1-1000"`
}

// GetDeviceLoadHistoryRes 获取设备负载历史响应
type GetDeviceLoadHistoryRes struct {
	List  []DeviceLoadRecord `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
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

// SetDeviceLoadThresholdReq 设置设备负载阈值请求
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

// GetDeviceLoadThresholdReq 获取设备负载阈值请求
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
