package consts

// ================================
// Device 模块常量定义
// ================================

// 设备状态常量
const (
	DeviceStatusOnline      = "online"
	DeviceStatusOffline     = "offline"
	DeviceStatusMaintenance = "maintenance"
	DeviceStatusError       = "error"
)

// 设备类型常量
const (
	DeviceTypeSensor   = "sensor"
	DeviceTypeCamera   = "camera"
	DeviceTypeActuator = "actuator"
	DeviceTypeGateway  = "gateway"
)

// 认证类型常量
const (
	AuthTypeNone        = "none"
	AuthTypeBasic       = "basic"
	AuthTypeToken       = "token"
	AuthTypeBearer      = "bearer"
	AuthTypeOAuth2      = "oauth2"
	AuthTypeAPIKey      = "apikey"
	AuthTypeCertificate = "certificate"
)

// 通信协议常量
const (
	ProtocolHTTP   = "http"
	ProtocolHTTPS  = "https"
	ProtocolTCP    = "tcp"
	ProtocolUDP    = "udp"
	ProtocolMQTT   = "mqtt"
	ProtocolCoAP   = "coap"
	ProtocolModbus = "modbus"
	ProtocolBACnet = "bacnet"
)

// 设备日志级别常量 (使用通用日志级别常量)
const (
	DeviceLogLevelDebug = LogLevelDebug
	DeviceLogLevelInfo  = LogLevelInfo
	DeviceLogLevelWarn  = LogLevelWarn
	DeviceLogLevelError = LogLevelError
	DeviceLogLevelFatal = "fatal" // 设备特有的致命级别
)

// 设备命令状态常量
const (
	CommandStatusPending   = "pending"
	CommandStatusSent      = "sent"
	CommandStatusExecuted  = "executed"
	CommandStatusCompleted = "completed"
	CommandStatusFailed    = "failed"
	CommandStatusTimeout   = "timeout"
)

// 设备命令类型常量
const (
	CommandTypeStart     = "start"
	CommandTypeStop      = "stop"
	CommandTypeRestart   = "restart"
	CommandTypeStatus    = "status"
	CommandTypeConfigure = "configure"
	CommandTypeReset     = "reset"
	CommandTypeUpdate    = "update"
	CommandTypeSync      = "sync"
)

// 设备健康度常量
const (
	HealthScoreMin       = 0
	HealthScoreMax       = 100
	HealthScoreDefault   = 100
	HealthScoreThreshold = 50 // 健康度阈值
)

// 设备心跳常量
const (
	HeartbeatIntervalSensor   = 120 // 传感器2分钟
	HeartbeatIntervalCamera   = 180 // 摄像头3分钟
	HeartbeatIntervalActuator = 60  // 执行器1分钟
	HeartbeatIntervalGateway  = 300 // 网关5分钟
	HeartbeatTimeoutMultiple  = 3   // 超时倍数（连续3次心跳丢失后认为离线）
)

// 设备连接常量
const (
	ConnectionTimeoutSeconds = 30 // 连接超时30秒
	MaxConnectionRetries     = 3  // 最大连接重试次数
	ConnectionRetryInterval  = 5  // 重试间隔5秒
)

// 设备负载常量
const (
	LoadThresholdNormal   = 0.6  // 正常负载阈值60%
	LoadThresholdWarning  = 0.8  // 警告负载阈值80%
	LoadThresholdCritical = 0.95 // 严重负载阈值95%
)

// 设备端口范围常量
const (
	PortMin = 1
	PortMax = 65535
)

// 设备配置常量
const (
	DefaultQueueSize       = 100  // 默认队列大小
	DefaultConcurrentTasks = 1    // 默认并发任务数
	DefaultTimeoutSeconds  = 3600 // 默认超时1小时
)
