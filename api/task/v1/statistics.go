package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 任务统计分析 API (2个)
// ===============================

// GetTaskStatisticsReq 获取任务统计信息请求
type GetTaskStatisticsReq struct {
	g.Meta    `path:"/task/statistics" method:"get" tags:"任务统计" summary:"获取任务统计信息"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	TaskType  string `json:"task_type,omitempty" v:"in:backup,sync,monitor,custom#任务类型无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	GroupBy   string `json:"group_by" d:"status" v:"in:status,type,device,priority#分组字段无效"`
}

type GetTaskStatisticsRes struct {
	Statistics TaskStatistics `json:"statistics"`
}

// GetTaskPerformanceReportReq 获取任务性能报告请求
type GetTaskPerformanceReportReq struct {
	g.Meta    `path:"/task/{taskId}/performance" method:"get" tags:"任务统计" summary:"获取任务性能报告"`
	TaskId    string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	TimeRange string `json:"time_range" d:"7d" v:"in:24h,7d,30d#时间范围无效"`
}

type GetTaskPerformanceReportRes struct {
	PerformanceReport TaskPerformanceReport `json:"performance_report"`
}
