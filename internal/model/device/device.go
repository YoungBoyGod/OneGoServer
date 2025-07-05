package device

import (
	"time"

	"OneGfServer/internal/consts"

	"github.com/gogf/gf/v2/os/gtime"
)

// 设备状态常量 - 使用统一常量
const (
	DeviceStatusOnline      = consts.DeviceStatusOnline
	DeviceStatusOffline     = consts.DeviceStatusOffline
	DeviceStatusMaintenance = consts.DeviceStatusMaintenance
	DeviceStatusError       = consts.DeviceStatusError
)

// 设备类型常量 - 使用统一常量
const (
	DeviceTypeSensor   = consts.DeviceTypeSensor
	DeviceTypeCamera   = consts.DeviceTypeCamera
	DeviceTypeActuator = consts.DeviceTypeActuator
	DeviceTypeGateway  = consts.DeviceTypeGateway
)

// 认证类型常量 - 使用统一常量
const (
	AuthTypeNone        = consts.AuthTypeNone
	AuthTypeBasic       = consts.AuthTypeBasic
	AuthTypeToken       = consts.AuthTypeToken
	AuthTypeCertificate = consts.AuthTypeCertificate
)

// 通信协议常量 - 使用统一常量
const (
	ProtocolHTTP = consts.ProtocolHTTP
	ProtocolMQTT = consts.ProtocolMQTT
	ProtocolTCP  = consts.ProtocolTCP
	ProtocolUDP  = consts.ProtocolUDP
)

// 日志级别常量 - 使用统一常量
const (
	LogLevelDEBUG = "DEBUG" // 保持原有大写格式
	LogLevelINFO  = "INFO"
	LogLevelWARN  = "WARN"
	LogLevelERROR = "ERROR"
)

// 命令状态常量 - 使用统一常量
const (
	CommandStatusPending   = consts.CommandStatusPending
	CommandStatusSent      = consts.CommandStatusSent
	CommandStatusExecuted  = consts.CommandStatusExecuted
	CommandStatusCompleted = consts.CommandStatusCompleted
	CommandStatusFailed    = consts.CommandStatusFailed
	CommandStatusTimeout   = consts.CommandStatusTimeout
)

// ===============================
// 基础CRUD操作 Input/Output
// ===============================

// CreateDeviceInput 创建设备输入
type CreateDeviceInput struct {
	Device *Device `json:"device"`
}

// CreateDeviceOutput 创建设备输出
type CreateDeviceOutput struct {
	DeviceId string `json:"device_id"`
	Message  string `json:"message"`
}

// GetDeviceByIDInput 根据ID获取设备输入
type GetDeviceByIDInput struct {
	ID int64 `json:"id"`
}

// GetDeviceByIDOutput 根据ID获取设备输出
type GetDeviceByIDOutput struct {
	Device *Device `json:"device"`
}

// GetDeviceByDeviceIDInput 根据设备ID获取设备输入
type GetDeviceByDeviceIDInput struct {
	DeviceID string `json:"device_id"`
}

// GetDeviceByDeviceIDOutput 根据设备ID获取设备输出
type GetDeviceByDeviceIDOutput struct {
	Device *Device `json:"device"`
}

// UpdateDeviceInput 更新设备输入
type UpdateDeviceInput struct {
	Device *Device `json:"device"`
}

// UpdateDeviceOutput 更新设备输出
type UpdateDeviceOutput struct {
	Message string `json:"message"`
}

// DeleteDeviceInput 删除设备输入
type DeleteDeviceInput struct {
	ID int64 `json:"id"`
}

// DeleteDeviceOutput 删除设备输出
type DeleteDeviceOutput struct {
	Message string `json:"message"`
}

// ===============================
// 查询操作 Input/Output
// ===============================

// GetDeviceListInput 获取设备列表输入
type GetDeviceListInput struct {
	Filter     *DeviceFilter     `json:"filter"`
	Sort       *DeviceSortOption `json:"sort"`
	Pagination *PaginationOption `json:"pagination"`
}

// GetDeviceListOutput 获取设备列表输出
type GetDeviceListOutput struct {
	List  []Device `json:"list"`
	Total int64    `json:"total"`
}

// GetDevicesByTypeInput 根据设备类型获取设备输入
type GetDevicesByTypeInput struct {
	DeviceTypes []string `json:"device_types"`
}

// GetDevicesByTypeOutput 根据设备类型获取设备输出
type GetDevicesByTypeOutput struct {
	Devices []Device `json:"devices"`
}

// GetDevicesByStatusInput 根据状态获取设备输入
type GetDevicesByStatusInput struct {
	Statuses []string `json:"statuses"`
}

// GetDevicesByStatusOutput 根据状态获取设备输出
type GetDevicesByStatusOutput struct {
	Devices []Device `json:"devices"`
}

// GetOnlineDevicesInput 获取在线设备输入
type GetOnlineDevicesInput struct{}

// GetOnlineDevicesOutput 获取在线设备输出
type GetOnlineDevicesOutput struct {
	Devices []Device `json:"devices"`
}

// GetOfflineDevicesInput 获取离线设备输入
type GetOfflineDevicesInput struct {
	Duration time.Duration `json:"duration"`
}

// GetOfflineDevicesOutput 获取离线设备输出
type GetOfflineDevicesOutput struct {
	Devices []Device `json:"devices"`
}

// ===============================
// 状态管理 Input/Output
// ===============================

// UpdateDeviceStatusInput 更新设备状态输入
type UpdateDeviceStatusInput struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
}

// UpdateDeviceStatusOutput 更新设备状态输出
type UpdateDeviceStatusOutput struct {
	Message string `json:"message"`
}

// UpdateDeviceHealthScoreInput 更新设备健康度输入
type UpdateDeviceHealthScoreInput struct {
	DeviceID string `json:"device_id"`
	Score    int    `json:"score"`
}

// UpdateDeviceHealthScoreOutput 更新设备健康度输出
type UpdateDeviceHealthScoreOutput struct {
	Message string `json:"message"`
}

// UpdateDeviceLastSeenInput 更新设备最后活跃时间输入
type UpdateDeviceLastSeenInput struct {
	DeviceID string    `json:"device_id"`
	LastSeen time.Time `json:"last_seen"`
}

// UpdateDeviceLastSeenOutput 更新设备最后活跃时间输出
type UpdateDeviceLastSeenOutput struct {
	Message string `json:"message"`
}

// BatchUpdateDeviceStatusInput 批量更新设备状态输入
type BatchUpdateDeviceStatusInput struct {
	DeviceIDs []string `json:"device_ids"`
	Status    string   `json:"status"`
}

// BatchUpdateDeviceStatusOutput 批量更新设备状态输出
type BatchUpdateDeviceStatusOutput struct {
	Message string `json:"message"`
}

// ===============================
// 统计查询 Input/Output
// ===============================

// GetDeviceStatisticsInput 获取设备统计信息输入
type GetDeviceStatisticsInput struct {
	Filter *DeviceFilter `json:"filter"`
}

// GetDeviceStatisticsOutput 获取设备统计信息输出
type GetDeviceStatisticsOutput struct {
	Statistics *DeviceStatistics `json:"statistics"`
}

// CountDevicesByStatusInput 根据状态统计设备数量输入
type CountDevicesByStatusInput struct{}

// CountDevicesByStatusOutput 根据状态统计设备数量输出
type CountDevicesByStatusOutput struct {
	Counts map[string]int64 `json:"counts"`
}

// CountDevicesByTypeInput 根据类型统计设备数量输入
type CountDevicesByTypeInput struct{}

// CountDevicesByTypeOutput 根据类型统计设备数量输出
type CountDevicesByTypeOutput struct {
	Counts map[string]int64 `json:"counts"`
}

// GetDeviceHealthReportInput 获取设备健康报告输入
type GetDeviceHealthReportInput struct{}

// GetDeviceHealthReportOutput 获取设备健康报告输出
type GetDeviceHealthReportOutput struct {
	Report map[string]interface{} `json:"report"`
}

// ===============================
// 心跳管理 Input/Output
// ===============================

// CreateDeviceHeartbeatInput 创建设备心跳输入
type CreateDeviceHeartbeatInput struct {
	Heartbeat *DeviceHeartbeat `json:"heartbeat"`
}

// CreateDeviceHeartbeatOutput 创建设备心跳输出
type CreateDeviceHeartbeatOutput struct {
	Message string `json:"message"`
}

// GetLatestDeviceHeartbeatInput 获取最新设备心跳输入
type GetLatestDeviceHeartbeatInput struct {
	DeviceID int64 `json:"device_id"`
}

// GetLatestDeviceHeartbeatOutput 获取最新设备心跳输出
type GetLatestDeviceHeartbeatOutput struct {
	Heartbeat *DeviceHeartbeat `json:"heartbeat"`
}

// GetDeviceHeartbeatHistoryInput 获取设备心跳历史输入
type GetDeviceHeartbeatHistoryInput struct {
	DeviceID int64 `json:"device_id"`
	Hours    int   `json:"hours"`
}

// GetDeviceHeartbeatHistoryOutput 获取设备心跳历史输出
type GetDeviceHeartbeatHistoryOutput struct {
	Heartbeats []DeviceHeartbeat `json:"heartbeats"`
}

// CleanupOldDeviceHeartbeatsInput 清理过期心跳输入
type CleanupOldDeviceHeartbeatsInput struct {
	Days int `json:"days"`
}

// CleanupOldDeviceHeartbeatsOutput 清理过期心跳输出
type CleanupOldDeviceHeartbeatsOutput struct {
	Message string `json:"message"`
}

// ===============================
// 日志管理 Input/Output
// ===============================

// CreateDeviceLogInput 创建设备日志输入
type CreateDeviceLogInput struct {
	Log *DeviceLog `json:"log"`
}

// CreateDeviceLogOutput 创建设备日志输出
type CreateDeviceLogOutput struct {
	Message string `json:"message"`
}

// GetDeviceLogsInput 根据设备ID获取设备日志输入
type GetDeviceLogsInput struct {
	DeviceID int64      `json:"device_id"`
	Filter   *LogFilter `json:"filter"`
}

// GetDeviceLogsOutput 根据设备ID获取设备日志输出
type GetDeviceLogsOutput struct {
	Logs []DeviceLog `json:"logs"`
}

// GetDeviceLogsByLevelInput 根据日志级别获取设备日志输入
type GetDeviceLogsByLevelInput struct {
	Level string `json:"level"`
	Limit int    `json:"limit"`
}

// GetDeviceLogsByLevelOutput 根据日志级别获取设备日志输出
type GetDeviceLogsByLevelOutput struct {
	Logs []DeviceLog `json:"logs"`
}

// CleanupOldDeviceLogsInput 清理过期日志输入
type CleanupOldDeviceLogsInput struct {
	Days int `json:"days"`
}

// CleanupOldDeviceLogsOutput 清理过期日志输出
type CleanupOldDeviceLogsOutput struct {
	Message string `json:"message"`
}

// ===============================
// 命令管理 Input/Output
// ===============================

// CreateDeviceCommandInput 创建设备命令输入
type CreateDeviceCommandInput struct {
	Command *DeviceCommand `json:"command"`
}

// CreateDeviceCommandOutput 创建设备命令输出
type CreateDeviceCommandOutput struct {
	Message string `json:"message"`
}

// GetDeviceCommandInput 根据命令ID获取设备命令输入
type GetDeviceCommandInput struct {
	CommandID string `json:"command_id"`
}

// GetDeviceCommandOutput 根据命令ID获取设备命令输出
type GetDeviceCommandOutput struct {
	Command *DeviceCommand `json:"command"`
}

// GetDeviceCommandsInput 根据设备ID获取设备命令输入
type GetDeviceCommandsInput struct {
	DeviceID int64    `json:"device_id"`
	Status   []string `json:"status"`
}

// GetDeviceCommandsOutput 根据设备ID获取设备命令输出
type GetDeviceCommandsOutput struct {
	Commands []DeviceCommand `json:"commands"`
}

// UpdateDeviceCommandStatusInput 更新设备命令状态输入
type UpdateDeviceCommandStatusInput struct {
	CommandID string `json:"command_id"`
	Status    string `json:"status"`
}

// UpdateDeviceCommandStatusOutput 更新设备命令状态输出
type UpdateDeviceCommandStatusOutput struct {
	Message string `json:"message"`
}

// UpdateDeviceCommandResponseInput 更新设备命令响应输入
type UpdateDeviceCommandResponseInput struct {
	CommandID string  `json:"command_id"`
	Response  JSONB   `json:"response"`
	ErrorMsg  *string `json:"error_msg"`
}

// UpdateDeviceCommandResponseOutput 更新设备命令响应输出
type UpdateDeviceCommandResponseOutput struct {
	Message string `json:"message"`
}

// ===============================
// 数据模型定义
// ===============================

// Device 设备主表模型
type Device struct {
	ID                   int64       `json:"id"`
	DeviceID             string      `json:"device_id"`
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	Model                string      `json:"model"`
	BoardID              string      `json:"board_id"`
	Status               string      `json:"status"`
	HealthScore          int         `json:"health_score"`
	LoginUsername        string      `json:"login_username"`
	LoginPort            int         `json:"login_port"`
	LoginPublicKey       string      `json:"login_public_key"`
	IPAddress            string      `json:"ip_address"`
	Port                 int         `json:"port"`
	Protocol             string      `json:"protocol"`
	Endpoint             string      `json:"endpoint"`
	RegTime              *gtime.Time `json:"reg_time"`
	Metadata             string      `json:"metadata"`
	Tags                 string      `json:"tags"`
	UptimeHours          float64     `json:"uptime_hours"`
	FirstOnlineTime      *gtime.Time `json:"first_online_time"`
	LastOnlineTime       *gtime.Time `json:"last_online_time"`
	LastOfflineTime      *gtime.Time `json:"last_offline_time"`
	TotalOnlineDuration  int64       `json:"total_online_duration"`
	TotalOfflineDuration int64       `json:"total_offline_duration"`
	TotalHeartbeats      int64       `json:"total_heartbeats"`
	TotalAlerts          int64       `json:"total_alerts"`
	TotalTasks           int64       `json:"total_tasks"`
	TotalSuccessTasks    int64       `json:"total_success_tasks"`
	TotalFailedTasks     int64       `json:"total_failed_tasks"`
	TotalCanceledTasks   int64       `json:"total_canceled_tasks"`
	TotalPendingTasks    int64       `json:"total_pending_tasks"`
	TotalRunningTasks    int64       `json:"total_running_tasks"`
	TotalCompletedTasks  int64       `json:"total_completed_tasks"`
	CreatedAt            *gtime.Time `json:"created_at"`
	UpdatedAt            *gtime.Time `json:"updated_at"`
	CreatedBy            int64       `json:"created_by"`
	UpdatedBy            int64       `json:"updated_by"`
}

// DeviceHeartbeat 设备心跳模型
type DeviceHeartbeat struct {
	ID             int64       `json:"id"`
	DeviceID       int64       `json:"device_id"`
	HeartbeatTime  *gtime.Time `json:"heartbeat_time"`
	CPUUsage       float64     `json:"cpu_usage"`
	MemoryUsage    float64     `json:"memory_usage"`
	DiskUsage      float64     `json:"disk_usage"`
	NetworkStatus  string      `json:"network_status"`
	NetworkLatency int         `json:"network_latency"`
	RunningTasks   int         `json:"running_tasks"`
	ErrorCount     int         `json:"error_count"`
	Metadata       string      `json:"metadata"`
	CreatedAt      *gtime.Time `json:"created_at"`
}

// DeviceLog 设备日志模型
type DeviceLog struct {
	ID        int64       `json:"id"`
	DeviceID  int64       `json:"device_id"`
	Level     string      `json:"level"`
	Category  string      `json:"category"`
	Message   string      `json:"message"`
	LogTime   *gtime.Time `json:"log_time"`
	Metadata  string      `json:"metadata"`
	CreatedAt *gtime.Time `json:"created_at"`
}

// DeviceCommand 设备命令模型
type DeviceCommand struct {
	ID            int64       `json:"id"`
	DeviceID      int64       `json:"device_id"`
	CommandID     string      `json:"command_id"`
	CommandType   string      `json:"command_type"`
	CommandData   string      `json:"command_data"`
	Status        string      `json:"status"`
	SentTime      *gtime.Time `json:"sent_time"`
	ExecutedTime  *gtime.Time `json:"executed_time"`
	CompletedTime *gtime.Time `json:"completed_time"`
	ResponseData  JSONB       `json:"response_data"`
	ErrorMessage  string      `json:"error_message"`
	CreatedAt     *gtime.Time `json:"created_at"`
	CreatedBy     int64       `json:"created_by"`
}

// ===============================
// 过滤和排序选项
// ===============================

// DeviceFilter 设备过滤条件
type DeviceFilter struct {
	Status       []string   `json:"status"`
	Type         []string   `json:"type"`
	Manufacturer *string    `json:"manufacturer"`
	Model        *string    `json:"model"`
	Protocol     *string    `json:"protocol"`
	MinHealth    *int       `json:"min_health"`
	MaxHealth    *int       `json:"max_health"`
	LastSeenFrom *time.Time `json:"last_seen_from"`
	LastSeenTo   *time.Time `json:"last_seen_to"`
	Keyword      *string    `json:"keyword"`
}

// DeviceSortOption 设备排序选项
type DeviceSortOption struct {
	Field string `json:"field"` // id, name, status, health_score, last_seen, created_at
	Order string `json:"order"` // asc, desc
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// LogFilter 日志过滤条件
type LogFilter struct {
	Level     *string    `json:"level"`
	Category  *string    `json:"category"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Keyword   *string    `json:"keyword"`
}

// ===============================
// 统计和报告模型
// ===============================

// DeviceStatistics 设备统计信息
type DeviceStatistics struct {
	Total        int64            `json:"total"`
	Online       int64            `json:"online"`
	Offline      int64            `json:"offline"`
	ByStatus     map[string]int64 `json:"by_status"`
	ByType       map[string]int64 `json:"by_type"`
	AvgHealth    float64          `json:"avg_health"`
	RecentActive int64            `json:"recent_active"`
}

// JSONB JSON二进制类型
type JSONB map[string]interface{}
