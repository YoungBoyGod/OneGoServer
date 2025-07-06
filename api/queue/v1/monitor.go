package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 队列监控统计 API (4个)
// ===============================

// GetQueueStatusReq 获取队列状态请求
type GetQueueStatusReq struct {
	g.Meta  `path:"/queue/{queueId}/status" method:"get" tags:"队列监控" summary:"获取队列状态"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueStatusRes struct {
	QueueStatus QueueStatusInfo `json:"queue_status"`
}

// GetQueueMetricsReq 获取队列指标请求
type GetQueueMetricsReq struct {
	g.Meta    `path:"/queue/{queueId}/metrics" method:"get" tags:"队列监控" summary:"获取队列指标"`
	QueueId   string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TimeRange string `json:"time_range" d:"1h" v:"in:5m,15m,1h,6h,24h,7d#时间范围无效"`
}

type GetQueueMetricsRes struct {
	Metrics QueueMetrics `json:"metrics"`
}

// GetQueueStatisticsReq 获取队列统计请求
type GetQueueStatisticsReq struct {
	g.Meta    `path:"/queue/statistics" method:"get" tags:"队列监控" summary:"获取队列统计"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	QueueType string `json:"queue_type,omitempty" v:"in:fifo,priority,delay,lifo#队列类型无效"`
	Status    string `json:"status,omitempty" v:"in:active,paused,stopped,error#状态值无效"`
	GroupBy   string `json:"group_by" d:"status" v:"in:status,type,priority#分组字段无效"`
}

type GetQueueStatisticsRes struct {
	Statistics QueueStatistics `json:"statistics"`
}

// GetQueuePerformanceReportReq 获取队列性能报告请求
type GetQueuePerformanceReportReq struct {
	g.Meta    `path:"/queue/{queueId}/performance" method:"get" tags:"队列监控" summary:"获取队列性能报告"`
	QueueId   string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TimeRange string `json:"time_range" d:"7d" v:"in:24h,7d,30d#时间范围无效"`
}

type GetQueuePerformanceReportRes struct {
	PerformanceReport QueuePerformanceReport `json:"performance_report"`
}
