package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备信息查询 API
// ==============================================

// GetDeviceDetail 获取设备详情请求
type GetDeviceDetailReq struct {
	g.Meta   `path:"/device/{deviceId}" method:"get" tags:"设备信息" summary:"获取设备详情"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceDetailRes 获取设备详情响应
type GetDeviceDetailRes struct {
	DeviceInfo
	LoginUsername        string      `json:"login_username"`
	LoginPort            int         `json:"login_port"`
	LoginPublicKey       string      `json:"login_public_key"`
	Metadata             string      `json:"metadata"`
	Tags                 string      `json:"tags"`
	FirstOnlineTime      *gtime.Time `json:"first_online_time"`
	TotalOfflineDuration int64       `json:"total_offline_duration"`
	TotalSuccessTasks    int64       `json:"total_success_tasks"`
	TotalFailedTasks     int64       `json:"total_failed_tasks"`
	TotalCanceledTasks   int64       `json:"total_canceled_tasks"`
	TotalPendingTasks    int64       `json:"total_pending_tasks"`
	TotalRunningTasks    int64       `json:"total_running_tasks"`
	TotalCompletedTasks  int64       `json:"total_completed_tasks"`
	CreatedBy            string      `json:"created_by"`
	UpdatedBy            string      `json:"updated_by"`
}

// UpdateDeviceInfo 更新设备信息请求
type UpdateDeviceInfoReq struct {
	g.Meta        `path:"/device/{deviceId}" method:"put" tags:"设备信息" summary:"更新设备信息"`
	DeviceId      string `json:"deviceId" v:"required#设备ID不能为空"`
	Name          string `json:"name,omitempty" v:"length:1,100#设备名称长度为1-100字符"`
	Model         string `json:"model,omitempty" v:"length:1,50#设备型号长度为1-50字符"`
	IpAddress     string `json:"ip_address,omitempty" v:"ip#IP地址格式不正确"`
	Port          int    `json:"port,omitempty" v:"between:1,65535#端口范围为1-65535"`
	Protocol      string `json:"protocol,omitempty" v:"in:http,mqtt,tcp,udp#协议只能是http,mqtt,tcp,udp"`
	LoginUsername string `json:"login_username,omitempty"`
	LoginPassword string `json:"login_password,omitempty"`
	Metadata      string `json:"metadata,omitempty"`
	Tags          string `json:"tags,omitempty"`
}

// UpdateDeviceInfoRes 更新设备信息响应
type UpdateDeviceInfoRes struct {
	Message string `json:"message"`
}
