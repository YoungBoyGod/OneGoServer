package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

/*
基础设备管理：
	1. 设备主动注册
	2. 获取设备列表
	3. 设备白名单管理
	4. 删除设备

设备状态管理：
	1. 获取设备状态
	2. 更新设备状态

设备控制操作：
	1. 发送设备命令
	2. 获取设备心跳
	3. 更新设备心跳

设备信息查询：
	1. 获取设备详情
	2. 更新设备信息

设备任务管理：
	1. 获取设备任务列表
	2. 获取设备任务队列
	3. 获取设备任务详情

设备告警管理：
	1. 获取设备告警列表
	2. 获取设备告警详情
	3. 更新设备告警
*/

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

// DeviceInfo 设备信息
type DeviceInfo struct {
	Id                  int64       `json:"id"`
	DeviceId            string      `json:"device_id"`
	Name                string      `json:"name"`
	Type                string      `json:"type"`
	Model               string      `json:"model"`
	BoardId             string      `json:"board_id"`
	Status              string      `json:"status"`
	HealthScore         int         `json:"health_score"`
	IpAddress           string      `json:"ip_address"`
	Port                int         `json:"port"`
	Protocol            string      `json:"protocol"`
	Endpoint            string      `json:"endpoint"`
	RegTime             *gtime.Time `json:"reg_time"`
	LastOnlineTime      *gtime.Time `json:"last_online_time"`
	LastOfflineTime     *gtime.Time `json:"last_offline_time"`
	TotalOnlineDuration int64       `json:"total_online_duration"`
	TotalHeartbeats     int64       `json:"total_heartbeats"`
	TotalAlerts         int64       `json:"total_alerts"`
	TotalTasks          int64       `json:"total_tasks"`
	CreatedAt           *gtime.Time `json:"created_at"`
	UpdatedAt           *gtime.Time `json:"updated_at"`
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
	g.Meta   `path:"/device/{deviceId}/heartbeat" method:"get" tags:"设备控制" summary:"获取设备心跳"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
}

// GetDeviceHeartbeatRes 获取设备心跳响应
type GetDeviceHeartbeatRes struct {
	List  []DeviceHeartbeatInfo `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// DeviceHeartbeatInfo 设备心跳信息
type DeviceHeartbeatInfo struct {
	Id            int64       `json:"id"`
	DeviceId      string      `json:"device_id"`
	HeartbeatTime *gtime.Time `json:"heartbeat_time"`
	Status        string      `json:"status"`
	IpAddress     string      `json:"ip_address"`
	ResponseTime  int         `json:"response_time"`
	Metadata      string      `json:"metadata"`
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

// ==============================================
// 设备任务管理 API
// ==============================================

// GetDeviceTaskList 获取设备任务列表请求
type GetDeviceTaskListReq struct {
	g.Meta   `path:"/device/{deviceId}/tasks" method:"get" tags:"设备任务" summary:"获取设备任务列表"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status   string `json:"status,omitempty" v:"in:pending,running,completed,failed,canceled#状态只能是pending,running,completed,failed,canceled"`
}

// GetDeviceTaskListRes 获取设备任务列表响应
type GetDeviceTaskListRes struct {
	List  []DeviceTaskInfo `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// DeviceTaskInfo 设备任务信息
type DeviceTaskInfo struct {
	Id        int64       `json:"id"`
	DeviceId  string      `json:"device_id"`
	TaskId    string      `json:"task_id"`
	TaskName  string      `json:"task_name"`
	TaskType  string      `json:"task_type"`
	Status    string      `json:"status"`
	Priority  int         `json:"priority"`
	CreatedAt *gtime.Time `json:"created_at"`
	StartTime *gtime.Time `json:"start_time"`
	EndTime   *gtime.Time `json:"end_time"`
}

// GetDeviceTaskQueue 获取设备任务队列请求
type GetDeviceTaskQueueReq struct {
	g.Meta   `path:"/device/{deviceId}/task-queue" method:"get" tags:"设备任务" summary:"获取设备任务队列"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
}

// GetDeviceTaskQueueRes 获取设备任务队列响应
type GetDeviceTaskQueueRes struct {
	List  []DeviceTaskQueueInfo `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// DeviceTaskQueueInfo 设备任务队列信息
type DeviceTaskQueueInfo struct {
	Id         int64       `json:"id"`
	DeviceId   string      `json:"device_id"`
	TaskId     string      `json:"task_id"`
	Priority   int         `json:"priority"`
	Status     string      `json:"status"`
	RetryCount int         `json:"retry_count"`
	CreatedAt  *gtime.Time `json:"created_at"`
	UpdatedAt  *gtime.Time `json:"updated_at"`
}

// GetDeviceTaskDetail 获取设备任务详情请求
type GetDeviceTaskDetailReq struct {
	g.Meta   `path:"/device/{deviceId}/task/{taskId}" method:"get" tags:"设备任务" summary:"获取设备任务详情"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	TaskId   string `json:"taskId" v:"required#任务ID不能为空"`
}

// GetDeviceTaskDetailRes 获取设备任务详情响应
type GetDeviceTaskDetailRes struct {
	TaskId        string      `json:"task_id"`
	DeviceId      string      `json:"device_id"`
	TaskName      string      `json:"task_name"`
	TaskType      string      `json:"task_type"`
	Description   string      `json:"description"`
	Status        string      `json:"status"`
	Priority      int         `json:"priority"`
	Progress      int         `json:"progress"`
	ResultData    string      `json:"result_data"`
	ErrorMessage  string      `json:"error_message"`
	RetryCount    int         `json:"retry_count"`
	MaxRetries    int         `json:"max_retries"`
	CreatedAt     *gtime.Time `json:"created_at"`
	StartTime     *gtime.Time `json:"start_time"`
	EndTime       *gtime.Time `json:"end_time"`
	EstimatedTime int64       `json:"estimated_time"`
	ActualTime    int64       `json:"actual_time"`
}

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

// DeviceAlertInfo 设备告警信息
type DeviceAlertInfo struct {
	Id          int64       `json:"id"`
	DeviceId    string      `json:"device_id"`
	AlertType   string      `json:"alert_type"`
	AlertLevel  string      `json:"alert_level"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	CreatedAt   *gtime.Time `json:"created_at"`
	ResolvedAt  *gtime.Time `json:"resolved_at"`
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
