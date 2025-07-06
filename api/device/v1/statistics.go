package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备统计管理 API
// ==============================================

// GetDeviceStatistics 获取设备统计信息请求
type GetDeviceStatisticsReq struct {
	g.Meta     `path:"/device/statistics" method:"get" tags:"设备统计" summary:"获取设备统计信息"`
	DeviceType string `json:"device_type,omitempty" v:"in:sensor,camera,actuator,gateway#设备类型"`
	Status     string `json:"status,omitempty" v:"in:online,offline,maintenance,error#状态"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	GroupBy    string `json:"group_by,omitempty" v:"in:type,status,date,hour#分组方式"`
}

// GetDeviceStatisticsRes 获取设备统计信息响应
type GetDeviceStatisticsRes struct {
	TotalDevices       int64                   `json:"total_devices"`
	OnlineDevices      int64                   `json:"online_devices"`
	OfflineDevices     int64                   `json:"offline_devices"`
	MaintenanceDevices int64                   `json:"maintenance_devices"`
	ErrorDevices       int64                   `json:"error_devices"`
	ByType             map[string]int64        `json:"by_type"`
	ByStatus           map[string]int64        `json:"by_status"`
	TrendData          []DeviceStatisticsTrend `json:"trend_data"`
	TopActiveDevices   []DeviceActivityInfo    `json:"top_active_devices"`
	AverageUptime      float64                 `json:"average_uptime"`
	AverageHealthScore float64                 `json:"average_health_score"`
}

// GetDevicePerformanceReport 获取设备性能报告请求
type GetDevicePerformanceReportReq struct {
	g.Meta     `path:"/device/{deviceId}/performance" method:"get" tags:"设备统计" summary:"获取设备性能报告"`
	DeviceId   string `json:"deviceId" v:"required#设备ID不能为空"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	ReportType string `json:"report_type" d:"summary" v:"in:summary,detailed,comparison#报告类型"`
}

// GetDevicePerformanceReportRes 获取设备性能报告响应
type GetDevicePerformanceReportRes struct {
	DeviceId            string                    `json:"device_id"`
	DeviceName          string                    `json:"device_name"`
	ReportPeriod        string                    `json:"report_period"`
	GeneratedAt         *gtime.Time               `json:"generated_at"`
	OverallScore        float64                   `json:"overall_score"`
	UptimePercentage    float64                   `json:"uptime_percentage"`
	TaskMetrics         DeviceTaskMetrics         `json:"task_metrics"`
	PerformanceMetrics  DevicePerformanceMetrics  `json:"performance_metrics"`
	AvailabilityMetrics DeviceAvailabilityMetrics `json:"availability_metrics"`
	TrendAnalysis       []PerformanceTrendPoint   `json:"trend_analysis"`
	Recommendations     []string                  `json:"recommendations"`
}
