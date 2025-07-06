package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ==============================================
// 设备日志管理 API
// ==============================================

// GetDeviceLogList 获取设备日志列表请求
type GetDeviceLogListReq struct {
	g.Meta                   `path:"/device/{deviceId}/logs" method:"get" tags:"设备日志" summary:"获取设备日志列表"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	common.PaginationRequest `json:",inline"`
	Level                    string `json:"level,omitempty" v:"in:DEBUG,INFO,WARN,ERROR#日志级别只能是DEBUG,INFO,WARN,ERROR"`
	Category                 string `json:"category,omitempty"`
	StartTime                string `json:"start_time,omitempty"`
	EndTime                  string `json:"end_time,omitempty"`
	Keyword                  string `json:"keyword,omitempty"`
	CorrelationId            string `json:"correlation_id,omitempty"`
}

// GetDeviceLogListRes 获取设备日志列表响应
type GetDeviceLogListRes struct {
	common.PaginationResponse[DeviceLogInfo] `json:",inline"`
}

// GetDeviceLogDetail 获取设备日志详情请求
type GetDeviceLogDetailReq struct {
	g.Meta   `path:"/device/{deviceId}/log/{logId}" method:"get" tags:"设备日志" summary:"获取设备日志详情"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	LogId    string `json:"logId" v:"required#日志ID不能为空"`
}

// GetDeviceLogDetailRes 获取设备日志详情响应
type GetDeviceLogDetailRes struct {
	DeviceLogInfo
}

// ClearDeviceLogs 清空设备日志请求
type ClearDeviceLogsReq struct {
	g.Meta     `path:"/device/{deviceId}/logs" method:"delete" tags:"设备日志" summary:"清空设备日志"`
	DeviceId   string `json:"deviceId" v:"required#设备ID不能为空"`
	BeforeTime string `json:"before_time,omitempty"`
	Level      string `json:"level,omitempty" v:"in:DEBUG,INFO,WARN,ERROR#日志级别只能是DEBUG,INFO,WARN,ERROR"`
	Category   string `json:"category,omitempty"`
	KeepRecent int    `json:"keep_recent,omitempty" v:"min:0#保留最近记录数不能为负数"`
}

// ClearDeviceLogsRes 清空设备日志响应
type ClearDeviceLogsRes struct {
	DeletedCount int64  `json:"deleted_count"`
	Message      string `json:"message"`
}
