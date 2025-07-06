package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统性能监控相关API
// ===============================

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
