package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备状态管理 API
// ==============================================

// GetDeviceStatus 获取设备状态请求
type GetDeviceStatusReq struct {
	g.Meta   `path:"/device/{deviceId}/status" method:"get" tags:"设备状态" summary:"获取设备状态"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceStatusRes 获取设备状态响应
type GetDeviceStatusRes struct {
	DeviceId            string      `json:"device_id"`
	Status              string      `json:"status"`
	HealthScore         int         `json:"health_score"`
	LastOnlineTime      *gtime.Time `json:"last_online_time"`
	LastOfflineTime     *gtime.Time `json:"last_offline_time"`
	TotalOnlineDuration int64       `json:"total_online_duration"`
	UptimeHours         float64     `json:"uptime_hours"`
}

// UpdateDeviceStatus 更新设备状态请求
type UpdateDeviceStatusReq struct {
	g.Meta      `path:"/device/{deviceId}/status" method:"put" tags:"设备状态" summary:"更新设备状态"`
	DeviceId    string `json:"deviceId" v:"required#设备ID不能为空"`
	Status      string `json:"status" v:"required|in:online,offline,maintenance,error#状态不能为空|状态只能是online,offline,maintenance,error"`
	HealthScore int    `json:"health_score,omitempty" v:"between:0,100#健康分数范围为0-100"`
}

// UpdateDeviceStatusRes 更新设备状态响应
type UpdateDeviceStatusRes struct {
	Message string `json:"message"`
}
