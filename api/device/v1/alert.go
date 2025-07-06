package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备告警管理 API
// ==============================================

// GetDeviceAlertList 获取设备告警列表请求
type GetDeviceAlertListReq struct {
	g.Meta     `path:"/device/{deviceId}/alerts" method:"get" tags:"设备告警" summary:"获取设备告警列表"`
	DeviceId   string `json:"deviceId" v:"required#设备ID不能为空"`
	Page       int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size       int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	AlertLevel string `json:"alert_level,omitempty" v:"in:low,medium,high,critical#告警级别只能是low,medium,high,critical"`
	Status     string `json:"status,omitempty" v:"in:active,resolved,ignored#状态只能是active,resolved,ignored"`
}

// GetDeviceAlertListRes 获取设备告警列表响应
type GetDeviceAlertListRes struct {
	List  []DeviceAlertInfo `json:"list"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

// GetDeviceAlertDetail 获取设备告警详情请求
type GetDeviceAlertDetailReq struct {
	g.Meta   `path:"/device/{deviceId}/alert/{alertId}" method:"get" tags:"设备告警" summary:"获取设备告警详情"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	AlertId  string `json:"alertId" v:"required#告警ID不能为空"`
}

// GetDeviceAlertDetailRes 获取设备告警详情响应
type GetDeviceAlertDetailRes struct {
	Id             int64       `json:"id"`
	DeviceId       string      `json:"device_id"`
	AlertType      string      `json:"alert_type"`
	AlertLevel     string      `json:"alert_level"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Status         string      `json:"status"`
	AlertData      string      `json:"alert_data"`
	TriggerValue   string      `json:"trigger_value"`
	ThresholdValue string      `json:"threshold_value"`
	CreatedAt      *gtime.Time `json:"created_at"`
	ResolvedAt     *gtime.Time `json:"resolved_at"`
	ResolvedBy     string      `json:"resolved_by"`
	Resolution     string      `json:"resolution"`
}

// UpdateDeviceAlert 更新设备告警请求
type UpdateDeviceAlertReq struct {
	g.Meta     `path:"/device/{deviceId}/alert/{alertId}" method:"put" tags:"设备告警" summary:"更新设备告警"`
	DeviceId   string `json:"deviceId" v:"required#设备ID不能为空"`
	AlertId    string `json:"alertId" v:"required#告警ID不能为空"`
	Status     string `json:"status" v:"required|in:active,resolved,ignored#状态不能为空|状态只能是active,resolved,ignored"`
	Resolution string `json:"resolution,omitempty"`
}

// UpdateDeviceAlertRes 更新设备告警响应
type UpdateDeviceAlertRes struct {
	Message string `json:"message"`
}
