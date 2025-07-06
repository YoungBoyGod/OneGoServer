package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统指标监控相关API
// ===============================

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
