package device

import "OneGfServer/internal/consts"

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

// Device 设备主表模型

// DeviceHeartbeat 设备心跳模型

// DeviceLog 设备日志模型

// DeviceCommand 设备命令模型

type GetDeviceListInput struct {
}

type GetDeviceListOutput struct {
}
