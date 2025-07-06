package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 基础设备管理 API
// ==============================================

// RegisterDeviceReq 设备注册请求
type RegisterDeviceReq struct {
	g.Meta       `path:"/device/register" method:"post" tags:"设备管理" summary:"设备注册"`
	DeviceName   string                 `json:"device_name" v:"required|length:1,100#设备名称不能为空|设备名称长度为1-100字符"`
	DeviceType   string                 `json:"device_type" v:"required|in:server,workstation,mobile,iot#设备类型不能为空"`
	DeviceModel  string                 `json:"device_model,omitempty" v:"length:1,100#设备型号长度为1-100字符"`
	SerialNumber string                 `json:"serial_number,omitempty" v:"length:1,100#序列号长度为1-100字符"`
	MacAddress   string                 `json:"mac_address,omitempty" v:"mac#MAC地址格式不正确"`
	IpAddress    string                 `json:"ip_address,omitempty" v:"ip#IP地址格式不正确"`
	Location     string                 `json:"location,omitempty" v:"length:1,200#位置信息长度为1-200字符"`
	Description  string                 `json:"description,omitempty" v:"max-length:500#设备描述最大500字符"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Tags         []string               `json:"tags,omitempty" v:"max:10#标签数量最多10个"`
	OwnerId      string                 `json:"owner_id,omitempty" v:"max-length:50#所有者ID最大50字符"`
	Department   string                 `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
}

// RegisterDeviceRes 注册设备响应
type RegisterDeviceRes struct {
	DeviceId string `json:"device_id"`
	Message  string `json:"message"`
}

// GetDeviceListReq 获取设备列表请求
type GetDeviceListReq struct {
	g.Meta                   `path:"/device/list" method:"get" tags:"设备管理" summary:"获取设备列表"`
	common.PaginationRequest `json:",inline"`
	Status                   string      `json:"status,omitempty" v:"in:online,offline,maintenance,error#状态值无效"`
	DeviceType               string      `json:"device_type,omitempty" v:"in:server,workstation,mobile,iot#设备类型无效"`
	Department               string      `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Keyword                  string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
}

// GetDeviceListRes 获取设备列表响应
type GetDeviceListRes struct {
	common.PaginationResponse[DeviceInfo] `json:",inline"`
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
