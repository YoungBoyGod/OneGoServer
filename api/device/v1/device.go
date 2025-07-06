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

// ==============================================
// 设备日志管理 API
// ==============================================

// GetDeviceLogList 获取设备日志列表请求
type GetDeviceLogListReq struct {
	g.Meta        `path:"/device/{deviceId}/logs" method:"get" tags:"设备日志" summary:"获取设备日志列表"`
	DeviceId      string `json:"deviceId" v:"required#设备ID不能为空"`
	Page          int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size          int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	Level         string `json:"level,omitempty" v:"in:DEBUG,INFO,WARN,ERROR#日志级别只能是DEBUG,INFO,WARN,ERROR"`
	Category      string `json:"category,omitempty"`
	StartTime     string `json:"start_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
	Keyword       string `json:"keyword,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
}

// GetDeviceLogListRes 获取设备日志列表响应
type GetDeviceLogListRes struct {
	List  []DeviceLogInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// DeviceLogInfo 设备日志信息
type DeviceLogInfo struct {
	Id            int64       `json:"id"`
	DeviceId      string      `json:"device_id"`
	LogTime       *gtime.Time `json:"log_time"`
	Level         string      `json:"level"`
	Category      string      `json:"category"`
	Message       string      `json:"message"`
	Details       string      `json:"details"`
	Source        string      `json:"source"`
	CorrelationId string      `json:"correlation_id"`
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

// ==============================================
// 设备负载监控 API
// ==============================================

// GetDeviceLoadStatus 获取设备负载状态请求
type GetDeviceLoadStatusReq struct {
	g.Meta   `path:"/device/{deviceId}/load" method:"get" tags:"设备监控" summary:"获取设备负载状态"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceLoadStatusRes 获取设备负载状态响应
type GetDeviceLoadStatusRes struct {
	DeviceId           int64       `json:"device_id"`
	CurrentTasks       int         `json:"current_tasks"`
	MaxConcurrentTasks int         `json:"max_concurrent_tasks"`
	CpuLoad            float64     `json:"cpu_load"`
	MemoryUsage        float64     `json:"memory_usage"`
	DiskUsage          float64     `json:"disk_usage"`
	NetworkLatency     int         `json:"network_latency"`
	Status             string      `json:"status"`
	LastHeartbeat      *gtime.Time `json:"last_heartbeat"`
	LoadScore          float64     `json:"load_score"`
	TotalAssigned      int         `json:"total_assigned"`
	TotalCompleted     int         `json:"total_completed"`
	TotalFailed        int         `json:"total_failed"`
	SuccessRate        float64     `json:"success_rate"`
	AvgTaskDuration    float64     `json:"avg_task_duration"`
	LastTaskCompletion *gtime.Time `json:"last_task_completion"`
	UpdatedAt          *gtime.Time `json:"updated_at"`
}

// DeviceLoadHistoryInfo 设备负载历史信息
type DeviceLoadHistoryInfo struct {
	DeviceId        int64       `json:"device_id"`
	Timestamp       *gtime.Time `json:"timestamp"`
	CpuLoad         float64     `json:"cpu_load"`
	MemoryUsage     float64     `json:"memory_usage"`
	DiskUsage       float64     `json:"disk_usage"`
	NetworkLatency  int         `json:"network_latency"`
	CurrentTasks    int         `json:"current_tasks"`
	LoadScore       float64     `json:"load_score"`
	SuccessRate     float64     `json:"success_rate"`
	AvgTaskDuration float64     `json:"avg_task_duration"`
}

// UpdateDeviceLoadConfig 更新设备负载配置请求
type UpdateDeviceLoadConfigReq struct {
	g.Meta             `path:"/device/{deviceId}/load/config" method:"put" tags:"设备监控" summary:"更新设备负载配置"`
	DeviceId           string  `json:"deviceId" v:"required#设备ID不能为空"`
	MaxConcurrentTasks int     `json:"max_concurrent_tasks,omitempty" v:"min:1#最大并发任务数最小为1"`
	CpuThreshold       float64 `json:"cpu_threshold,omitempty" v:"between:0,100#CPU阈值范围为0-100"`
	MemoryThreshold    float64 `json:"memory_threshold,omitempty" v:"between:0,100#内存阈值范围为0-100"`
	DiskThreshold      float64 `json:"disk_threshold,omitempty" v:"between:0,100#磁盘阈值范围为0-100"`
	NetworkThreshold   int     `json:"network_threshold,omitempty" v:"min:0#网络延迟阈值不能为负数"`
}

// UpdateDeviceLoadConfigRes 更新设备负载配置响应
type UpdateDeviceLoadConfigRes struct {
	Message string `json:"message"`
}

// ==============================================
// 设备队列操作历史 API
// ==============================================

// GetDeviceQueueHistory 获取设备队列操作历史请求
type GetDeviceQueueHistoryReq struct {
	g.Meta        `path:"/device/{deviceId}/queue/history" method:"get" tags:"设备队列" summary:"获取设备队列操作历史"`
	DeviceId      string `json:"deviceId" v:"required#设备ID不能为空"`
	Page          int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size          int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	OperationType string `json:"operation_type,omitempty" v:"in:enqueue,dequeue,priority_change,cancel,restart#操作类型"`
	TaskId        string `json:"task_id,omitempty"`
	StartTime     string `json:"start_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
	BatchId       string `json:"batch_id,omitempty"`
}

// GetDeviceQueueHistoryRes 获取设备队列操作历史响应
type GetDeviceQueueHistoryRes struct {
	List  []DeviceQueueHistoryInfo `json:"list"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

// DeviceQueueHistoryInfo 设备队列操作历史信息
type DeviceQueueHistoryInfo struct {
	Id               int64       `json:"id"`
	DeviceId         int64       `json:"device_id"`
	TaskId           string      `json:"task_id"`
	OperationType    string      `json:"operation_type"`
	OperationBy      int64       `json:"operation_by"`
	OperationTime    *gtime.Time `json:"operation_time"`
	OldPriority      int         `json:"old_priority"`
	NewPriority      int         `json:"new_priority"`
	OldPosition      int         `json:"old_position"`
	NewPosition      int         `json:"new_position"`
	OldStatus        string      `json:"old_status"`
	NewStatus        string      `json:"new_status"`
	Reason           string      `json:"reason"`
	Notes            string      `json:"notes"`
	OperationSource  string      `json:"operation_source"`
	BatchId          string      `json:"batch_id"`
	IsBatchOperation bool        `json:"is_batch_operation"`
}

// ==============================================
// 设备命令执行历史 API
// ==============================================

// GetDeviceCommandHistory 获取设备命令执行历史请求
type GetDeviceCommandHistoryReq struct {
	g.Meta      `path:"/device/{deviceId}/commands" method:"get" tags:"设备控制" summary:"获取设备命令执行历史"`
	DeviceId    string `json:"deviceId" v:"required#设备ID不能为空"`
	Page        int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size        int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	CommandType string `json:"command_type,omitempty"`
	Status      string `json:"status,omitempty" v:"in:pending,sent,executed,completed,failed,timeout#状态限制"`
	StartTime   string `json:"start_time,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

// GetDeviceCommandHistoryRes 获取设备命令执行历史响应
type GetDeviceCommandHistoryRes struct {
	List  []DeviceCommandInfo `json:"list"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

// DeviceCommandInfo 设备命令信息
type DeviceCommandInfo struct {
	Id            int64       `json:"id"`
	DeviceId      string      `json:"device_id"`
	CommandType   string      `json:"command_type"`
	CommandData   string      `json:"command_data"`
	Status        string      `json:"status"`
	SentTime      *gtime.Time `json:"sent_time"`
	ExecutedTime  *gtime.Time `json:"executed_time"`
	CompletedTime *gtime.Time `json:"completed_time"`
	ResponseData  string      `json:"response_data"`
	ErrorMessage  string      `json:"error_message"`
	Duration      int64       `json:"duration"`
	CreatedAt     *gtime.Time `json:"created_at"`
	CreatedBy     int64       `json:"created_by"`
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

// ==============================================
// 设备批量操作 API
// ==============================================

// BatchOperateDevices 批量操作设备请求
type BatchOperateDevicesReq struct {
	g.Meta      `path:"/device/batch" method:"post" tags:"设备管理" summary:"批量操作设备"`
	DeviceIds   []string `json:"device_ids" v:"required|length:1,100#设备ID列表不能为空|最多支持100个设备"`
	Operation   string   `json:"operation" v:"required|in:start,stop,restart,delete,update_status,send_command#操作类型不能为空"`
	Parameters  string   `json:"parameters,omitempty"`
	BatchReason string   `json:"batch_reason,omitempty"`
}

// BatchOperateDevicesRes 批量操作设备响应
type BatchOperateDevicesRes struct {
	BatchId      string                 `json:"batch_id"`
	TotalCount   int                    `json:"total_count"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	SkippedCount int                    `json:"skipped_count"`
	Results      []BatchOperationResult `json:"results"`
	Message      string                 `json:"message"`
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	DeviceId string `json:"device_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
}

// GetBatchOperationStatus 获取批量操作状态请求
type GetBatchOperationStatusReq struct {
	g.Meta  `path:"/device/batch/{batchId}" method:"get" tags:"设备管理" summary:"获取批量操作状态"`
	BatchId string `json:"batchId" v:"required#批次ID不能为空"`
}

// GetBatchOperationStatusRes 获取批量操作状态响应
type GetBatchOperationStatusRes struct {
	BatchId        string                 `json:"batch_id"`
	Operation      string                 `json:"operation"`
	Status         string                 `json:"status"` // pending, processing, completed, failed
	TotalCount     int                    `json:"total_count"`
	ProcessedCount int                    `json:"processed_count"`
	SuccessCount   int                    `json:"success_count"`
	FailedCount    int                    `json:"failed_count"`
	Progress       float64                `json:"progress"`
	StartTime      *gtime.Time            `json:"start_time"`
	EndTime        *gtime.Time            `json:"end_time"`
	Results        []BatchOperationResult `json:"results"`
}

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

// DeviceConfigHistoryInfo 设备配置历史信息
type DeviceConfigHistoryInfo struct {
	Id             int64                  `json:"id"`
	DeviceId       string                 `json:"device_id"`
	Version        string                 `json:"version"`
	Configurations map[string]interface{} `json:"configurations"`
	ChangeType     string                 `json:"change_type"` // create, update, delete
	Reason         string                 `json:"reason"`
	ApplyStatus    string                 `json:"apply_status"` // pending, applied, failed
	CreatedAt      *gtime.Time            `json:"created_at"`
	UpdatedBy      string                 `json:"updated_by"`
}

// ==============================================
// 设备统计分析 API
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

// DeviceStatisticsTrend 设备统计趋势数据
type DeviceStatisticsTrend struct {
	Timestamp      *gtime.Time      `json:"timestamp"`
	OnlineCount    int64            `json:"online_count"`
	OfflineCount   int64            `json:"offline_count"`
	NewDevices     int64            `json:"new_devices"`
	TotalTasks     int64            `json:"total_tasks"`
	CompletedTasks int64            `json:"completed_tasks"`
	FailedTasks    int64            `json:"failed_tasks"`
	ByType         map[string]int64 `json:"by_type"`
}

// DeviceActivityInfo 设备活跃度信息
type DeviceActivityInfo struct {
	DeviceId        string      `json:"device_id"`
	DeviceName      string      `json:"device_name"`
	DeviceType      string      `json:"device_type"`
	TaskCount       int64       `json:"task_count"`
	CompletedTasks  int64       `json:"completed_tasks"`
	FailedTasks     int64       `json:"failed_tasks"`
	SuccessRate     float64     `json:"success_rate"`
	AvgTaskDuration float64     `json:"avg_task_duration"`
	UptimeHours     float64     `json:"uptime_hours"`
	HealthScore     int         `json:"health_score"`
	LastActiveTime  *gtime.Time `json:"last_active_time"`
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

// DeviceTaskMetrics 设备任务指标
type DeviceTaskMetrics struct {
	TotalTasks       int64   `json:"total_tasks"`
	CompletedTasks   int64   `json:"completed_tasks"`
	FailedTasks      int64   `json:"failed_tasks"`
	CanceledTasks    int64   `json:"canceled_tasks"`
	SuccessRate      float64 `json:"success_rate"`
	AvgExecutionTime float64 `json:"avg_execution_time"`
	MaxExecutionTime float64 `json:"max_execution_time"`
	MinExecutionTime float64 `json:"min_execution_time"`
	TasksPerHour     float64 `json:"tasks_per_hour"`
}

// DevicePerformanceMetrics 设备性能指标
type DevicePerformanceMetrics struct {
	AvgCpuUsage       float64 `json:"avg_cpu_usage"`
	PeakCpuUsage      float64 `json:"peak_cpu_usage"`
	AvgMemoryUsage    float64 `json:"avg_memory_usage"`
	PeakMemoryUsage   float64 `json:"peak_memory_usage"`
	AvgDiskUsage      float64 `json:"avg_disk_usage"`
	AvgNetworkLatency float64 `json:"avg_network_latency"`
	MaxNetworkLatency float64 `json:"max_network_latency"`
	AvgResponseTime   float64 `json:"avg_response_time"`
}

// DeviceAvailabilityMetrics 设备可用性指标
type DeviceAvailabilityMetrics struct {
	TotalHours        float64 `json:"total_hours"`
	OnlineHours       float64 `json:"online_hours"`
	OfflineHours      float64 `json:"offline_hours"`
	MaintenanceHours  float64 `json:"maintenance_hours"`
	UptimePercentage  float64 `json:"uptime_percentage"`
	MTBF              float64 `json:"mtbf"` // Mean Time Between Failures
	MTTR              float64 `json:"mttr"` // Mean Time To Recovery
	AvailabilityScore float64 `json:"availability_score"`
}

// PerformanceTrendPoint 性能趋势点
type PerformanceTrendPoint struct {
	Timestamp      *gtime.Time `json:"timestamp"`
	CpuUsage       float64     `json:"cpu_usage"`
	MemoryUsage    float64     `json:"memory_usage"`
	DiskUsage      float64     `json:"disk_usage"`
	NetworkLatency float64     `json:"network_latency"`
	TaskCount      int64       `json:"task_count"`
	SuccessRate    float64     `json:"success_rate"`
	ResponseTime   float64     `json:"response_time"`
	HealthScore    float64     `json:"health_score"`
}

// ===============================
// 设备负载评分相关API
// ===============================

// CalculateDeviceLoadScoreReq 计算设备负载评分请求
type CalculateDeviceLoadScoreReq struct {
	g.Meta   `path:"/device/load/calculate" method:"post" tags:"设备管理" summary:"计算设备负载评分"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// CalculateDeviceLoadScoreRes 计算设备负载评分响应
type CalculateDeviceLoadScoreRes struct {
	DeviceId   string  `json:"deviceId"`
	LoadScore  float64 `json:"loadScore"`
	Components struct {
		CpuScore     float64 `json:"cpuScore"`
		MemoryScore  float64 `json:"memoryScore"`
		DiskScore    float64 `json:"diskScore"`
		NetworkScore float64 `json:"networkScore"`
		TaskScore    float64 `json:"taskScore"`
	} `json:"components"`
	Recommendations []string `json:"recommendations"`
}

// GetDeviceLoadMetricsReq 获取设备负载指标请求
type GetDeviceLoadMetricsReq struct {
	g.Meta   `path:"/device/load/metrics" method:"get" tags:"设备管理" summary:"获取设备负载指标"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Period   string `json:"period" d:"1h" v:"in:1h,6h,24h,7d#时间周期只能是1h,6h,24h,7d"`
}

// GetDeviceLoadMetricsRes 获取设备负载指标响应
type GetDeviceLoadMetricsRes struct {
	DeviceId string `json:"deviceId"`
	Period   string `json:"period"`
	Metrics  struct {
		CpuUsage       []TimePoint `json:"cpuUsage"`
		MemoryUsage    []TimePoint `json:"memoryUsage"`
		DiskUsage      []TimePoint `json:"diskUsage"`
		NetworkLatency []TimePoint `json:"networkLatency"`
		CurrentTasks   []TimePoint `json:"currentTasks"`
		LoadScore      []TimePoint `json:"loadScore"`
	} `json:"metrics"`
}

// TimePoint 时间点数据
type TimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// OptimizeDeviceLoadReq 优化设备负载请求
type OptimizeDeviceLoadReq struct {
	g.Meta   `path:"/device/load/optimize" method:"post" tags:"设备管理" summary:"优化设备负载"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Strategy string `json:"strategy" v:"in:auto,manual,conservative,aggressive#优化策略只能是auto,manual,conservative,aggressive"`
}

// OptimizeDeviceLoadRes 优化设备负载响应
type OptimizeDeviceLoadRes struct {
	DeviceId       string   `json:"deviceId"`
	OriginalScore  float64  `json:"originalScore"`
	OptimizedScore float64  `json:"optimizedScore"`
	Improvement    float64  `json:"improvement"`
	Actions        []string `json:"actions"`
	Status         string   `json:"status"`
}

// GetDeviceLoadHistoryReq 获取设备负载历史请求
type GetDeviceLoadHistoryReq struct {
	g.Meta    `path:"/device/load/history" method:"get" tags:"设备管理" summary:"获取设备负载历史"`
	DeviceId  string `json:"deviceId" v:"required#设备ID不能为空"`
	StartTime string `json:"startTime" v:"required#开始时间不能为空"`
	EndTime   string `json:"endTime" v:"required#结束时间不能为空"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"100" v:"between:1,1000#每页数量为1-1000"`
}

// GetDeviceLoadHistoryRes 获取设备负载历史响应
type GetDeviceLoadHistoryRes struct {
	List  []DeviceLoadRecord `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// DeviceLoadRecord 设备负载记录
type DeviceLoadRecord struct {
	Id             int64   `json:"id"`
	DeviceId       string  `json:"deviceId"`
	LoadScore      float64 `json:"loadScore"`
	CpuUsage       float64 `json:"cpuUsage"`
	MemoryUsage    float64 `json:"memoryUsage"`
	DiskUsage      float64 `json:"diskUsage"`
	NetworkLatency float64 `json:"networkLatency"`
	CurrentTasks   int     `json:"currentTasks"`
	MaxTasks       int     `json:"maxTasks"`
	RecordedAt     string  `json:"recordedAt"`
}

// SetDeviceLoadThresholdReq 设置设备负载阈值请求
type SetDeviceLoadThresholdReq struct {
	g.Meta   `path:"/device/load/threshold" method:"post" tags:"设备管理" summary:"设置设备负载阈值"`
	DeviceId string  `json:"deviceId" v:"required#设备ID不能为空"`
	Warning  float64 `json:"warning" v:"between:0,100#警告阈值必须在0-100之间"`
	Critical float64 `json:"critical" v:"between:0,100#严重阈值必须在0-100之间"`
	MaxTasks int     `json:"maxTasks" v:"min:1#最大任务数最小为1"`
}

// SetDeviceLoadThresholdRes 设置设备负载阈值响应
type SetDeviceLoadThresholdRes struct {
	DeviceId string  `json:"deviceId"`
	Warning  float64 `json:"warning"`
	Critical float64 `json:"critical"`
	MaxTasks int     `json:"maxTasks"`
	Status   string  `json:"status"`
}

// GetDeviceLoadThresholdReq 获取设备负载阈值请求
type GetDeviceLoadThresholdReq struct {
	g.Meta   `path:"/device/load/threshold" method:"get" tags:"设备管理" summary:"获取设备负载阈值"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
}

// GetDeviceLoadThresholdRes 获取设备负载阈值响应
type GetDeviceLoadThresholdRes struct {
	DeviceId    string  `json:"deviceId"`
	Warning     float64 `json:"warning"`
	Critical    float64 `json:"critical"`
	MaxTasks    int     `json:"maxTasks"`
	CurrentLoad float64 `json:"currentLoad"`
	Status      string  `json:"status"`
}
