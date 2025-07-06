package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ==============================================
// 基础设备管理 API
// ==============================================

// RegisterDevice 注册设备请求
type RegisterDeviceReq struct {
	g.Meta        `path:"/device/register" method:"post" tags:"设备管理" summary:"注册设备"`
	DeviceName    string `json:"device_name" v:"required|length:1,100#设备名称不能为空|设备名称长度为1-100字符"`
	DeviceType    string `json:"device_type" v:"required|in:sensor,camera,actuator,gateway#设备类型不能为空|设备类型只能是sensor,camera,actuator,gateway"`
	Model         string `json:"model" v:"required|length:1,50#设备型号不能为空|设备型号长度为1-50字符"`
	BoardId       string `json:"board_id" v:"required|length:1,50#主板ID不能为空|主板ID长度为1-50字符"`
	IpAddress     string `json:"ip_address" v:"required|ip#IP地址不能为空|IP地址格式不正确"`
	Port          int    `json:"port" v:"required|between:1,65535#端口不能为空|端口范围为1-65535"`
	Protocol      string `json:"protocol" v:"required|in:http,mqtt,tcp,udp#协议不能为空|协议只能是http,mqtt,tcp,udp"`
	LoginUsername string `json:"login_username,omitempty"`
	LoginPassword string `json:"login_password,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	Tags          string `json:"tags,omitempty"`
}

// RegisterDeviceRes 注册设备响应
type RegisterDeviceRes struct {
	DeviceId string `json:"device_id"`
	Message  string `json:"message"`
}

// GetDeviceList 获取设备列表请求
type GetDeviceListReq struct {
	g.Meta     `path:"/device/list" method:"get" tags:"设备管理" summary:"获取设备列表"`
	Page       int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size       int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	DeviceType string `json:"device_type,omitempty" v:"in:sensor,camera,actuator,gateway#设备类型只能是sensor,camera,actuator,gateway"`
	Status     string `json:"status,omitempty" v:"in:online,offline,maintenance,error#状态只能是online,offline,maintenance,error"`
	Keyword    string `json:"keyword,omitempty"`
}

// GetDeviceListRes 获取设备列表响应
type GetDeviceListRes struct {
	List  []DeviceInfo `json:"list"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// ManageDeviceWhitelist 管理设备白名单请求
type ManageDeviceWhitelistReq struct {
	g.Meta   `path:"/device/whitelist" method:"post" tags:"设备管理" summary:"管理设备白名单"`
	DeviceId string `json:"device_id" v:"required#设备ID不能为空"`
	Action   string `json:"action" v:"required|in:add,remove#操作类型不能为空|操作类型只能是add,remove"`
}

// ManageDeviceWhitelistRes 管理设备白名单响应
type ManageDeviceWhitelistRes struct {
	Message string `json:"message"`
}

// DeleteDevice 删除设备请求
type DeleteDeviceReq struct {
	g.Meta   `path:"/device/{deviceId}" method:"delete" tags:"设备管理" summary:"删除设备"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// DeleteDeviceRes 删除设备响应
type DeleteDeviceRes struct {
	Message string `json:"message"`
}
