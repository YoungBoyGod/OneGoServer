package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ==============================================
// 设备控制操作 API
// ==============================================

// SendDeviceCommand 发送设备命令请求
type SendDeviceCommandReq struct {
	g.Meta      `path:"/device/{deviceId}/command" method:"post" tags:"设备控制" summary:"发送设备命令"`
	DeviceId    string `json:"deviceId" v:"required#设备ID不能为空"`
	CommandType string `json:"command_type" v:"required#命令类型不能为空"`
	CommandData string `json:"command_data" v:"required#命令数据不能为空"`
}

// SendDeviceCommandRes 发送设备命令响应
type SendDeviceCommandRes struct {
	CommandId string `json:"command_id"`
	Message   string `json:"message"`
}

// GetDeviceHeartbeat 获取设备心跳请求
type GetDeviceHeartbeatReq struct {
	g.Meta                   `path:"/device/{deviceId}/heartbeat" method:"get" tags:"设备控制" summary:"获取设备心跳"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	common.PaginationRequest `json:",inline"`
}

// GetDeviceHeartbeatRes 获取设备心跳响应
type GetDeviceHeartbeatRes struct {
	common.PaginationResponse[DeviceHeartbeatInfo] `json:",inline"`
}

// UpdateDeviceHeartbeat 更新设备心跳请求
type UpdateDeviceHeartbeatReq struct {
	g.Meta       `path:"/device/{deviceId}/heartbeat" method:"post" tags:"设备控制" summary:"更新设备心跳"`
	DeviceId     string `json:"deviceId" v:"required#设备ID不能为空"`
	Status       string `json:"status" v:"required|in:online,offline,error#状态不能为空|状态只能是online,offline,error"`
	ResponseTime int    `json:"response_time,omitempty" v:"min:0#响应时间不能为负数"`
	Metadata     string `json:"metadata,omitempty"`
}

// UpdateDeviceHeartbeatRes 更新设备心跳响应
type UpdateDeviceHeartbeatRes struct {
	Message string `json:"message"`
}

// GetDeviceCommandHistory 获取设备命令执行历史请求
type GetDeviceCommandHistoryReq struct {
	g.Meta                   `path:"/device/{deviceId}/commands" method:"get" tags:"设备控制" summary:"获取设备命令执行历史"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	common.PaginationRequest `json:",inline"`
	CommandType              string `json:"command_type,omitempty"`
	Status                   string `json:"status,omitempty" v:"in:pending,sent,executed,completed,failed,timeout#状态限制"`
	StartTime                string `json:"start_time,omitempty"`
	EndTime                  string `json:"end_time,omitempty"`
}

// GetDeviceCommandHistoryRes 获取设备命令执行历史响应
type GetDeviceCommandHistoryRes struct {
	common.PaginationResponse[DeviceCommandInfo] `json:",inline"`
}

// GetDeviceCommandDetail 获取设备命令详情请求
type GetDeviceCommandDetailReq struct {
	g.Meta    `path:"/device/{deviceId}/command/{commandId}" method:"get" tags:"设备控制" summary:"获取设备命令详情"`
	DeviceId  string `json:"deviceId" v:"required#设备ID不能为空"`
	CommandId string `json:"commandId" v:"required#命令ID不能为空"`
}

// GetDeviceCommandDetailRes 获取设备命令详情响应
type GetDeviceCommandDetailRes struct {
	DeviceCommandInfo
}
