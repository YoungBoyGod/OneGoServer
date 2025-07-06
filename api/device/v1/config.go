package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备配置管理 API
// ==============================================

// GetDeviceConfig 获取设备配置请求
type GetDeviceConfigReq struct {
	g.Meta   `path:"/device/{deviceId}/config" method:"get" tags:"设备配置" summary:"获取设备配置"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceConfigRes 获取设备配置响应
type GetDeviceConfigRes struct {
	DeviceId       string                 `json:"device_id"`
	Configurations map[string]interface{} `json:"configurations"`
	Version        string                 `json:"version"`
	LastUpdated    *gtime.Time            `json:"last_updated"`
	UpdatedBy      string                 `json:"updated_by"`
}

// UpdateDeviceConfig 更新设备配置请求
type UpdateDeviceConfigReq struct {
	g.Meta         `path:"/device/{deviceId}/config" method:"put" tags:"设备配置" summary:"更新设备配置"`
	DeviceId       string                 `json:"deviceId" v:"required#设备ID不能为空"`
	Configurations map[string]interface{} `json:"configurations" v:"required#配置信息不能为空"`
	Reason         string                 `json:"reason,omitempty"`
	ApplyNow       bool                   `json:"apply_now" d:"true"`
}

// UpdateDeviceConfigRes 更新设备配置响应
type UpdateDeviceConfigRes struct {
	Message     string `json:"message"`
	Version     string `json:"version"`
	ConfigId    string `json:"config_id"`
	ApplyStatus string `json:"apply_status"`
}

// GetDeviceConfigHistory 获取设备配置历史请求
type GetDeviceConfigHistoryReq struct {
	g.Meta    `path:"/device/{deviceId}/config/history" method:"get" tags:"设备配置" summary:"获取设备配置历史"`
	DeviceId  string `json:"deviceId" v:"required#设备ID不能为空"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"10" v:"between:1,50#每页数量为1-50"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// GetDeviceConfigHistoryRes 获取设备配置历史响应
type GetDeviceConfigHistoryRes struct {
	List  []DeviceConfigHistoryInfo `json:"list"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"size"`
}
