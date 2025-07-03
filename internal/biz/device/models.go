package device

import (
	"database/sql/driver"
	"encoding/json"
	"net"
	"time"

	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
	"gorm.io/gorm"
)

// JSONB 自定义JSON类型，用于PostgreSQL的JSONB字段
type JSONB map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// Device 设备主表模型
type Device struct {
	ID           int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID     string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"device_id"`
	Name         string  `gorm:"type:varchar(255);not null" json:"name"`
	Type         string  `gorm:"type:varchar(50);not null" json:"type"`
	Model        *string `gorm:"type:varchar(100)" json:"model,omitempty"`
	Manufacturer *string `gorm:"type:varchar(100)" json:"manufacturer,omitempty"`

	// 连接信息
	IPAddress *net.IP `gorm:"type:inet" json:"ip_address,omitempty"`
	Port      *int    `json:"port,omitempty"`
	Protocol  *string `gorm:"type:varchar(20)" json:"protocol,omitempty"`
	Endpoint  *string `gorm:"type:varchar(500)" json:"endpoint,omitempty"`

	// 认证信息
	AuthType    *string `gorm:"type:varchar(20)" json:"auth_type,omitempty"`
	Credentials *JSONB  `gorm:"type:jsonb;comment:加密存储" json:"credentials,omitempty"`

	// 状态信息
	Status      string     `gorm:"type:varchar(20);not null;default:offline" json:"status"`
	HealthScore int        `gorm:"default:100;comment:健康度(0-100)" json:"health_score"`
	LastSeen    *time.Time `gorm:"type:timestamp" json:"last_seen,omitempty"`

	// 配置信息
	Config       *JSONB `gorm:"type:jsonb" json:"config,omitempty"`
	Capabilities *JSONB `gorm:"type:jsonb" json:"capabilities,omitempty"`
	Metadata     *JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`

	// 统计信息
	TotalCommands   int64   `gorm:"default:0" json:"total_commands"`
	SuccessCommands int64   `gorm:"default:0" json:"success_commands"`
	FailedCommands  int64   `gorm:"default:0" json:"failed_commands"`
	UptimeHours     float64 `gorm:"type:numeric(10,2);default:0" json:"uptime_hours"`

	// 审计字段
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	CreatedBy *int64    `gorm:"index" json:"created_by,omitempty"`
	UpdatedBy *int64    `gorm:"index" json:"updated_by,omitempty"`

	// 关联关系
	Heartbeats []DeviceHeartbeat `gorm:"foreignKey:DeviceID;references:ID" json:"heartbeats,omitempty"`
	Logs       []DeviceLog       `gorm:"foreignKey:DeviceID;references:ID" json:"logs,omitempty"`
	Commands   []DeviceCommand   `gorm:"foreignKey:DeviceID;references:ID" json:"commands,omitempty"`
}

// TableName 指定表名
func (Device) TableName() string {
	return "devices"
}

// BeforeUpdate GORM钩子，更新时自动设置更新时间
func (d *Device) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = time.Now()
	return
}

// DeviceHeartbeat 设备心跳模型
type DeviceHeartbeat struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID      int64     `gorm:"not null;index" json:"device_id"`
	HeartbeatTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"heartbeat_time"`

	// 心跳数据
	Status     string `gorm:"type:varchar(20);not null" json:"status"`
	Metrics    *JSONB `gorm:"type:jsonb" json:"metrics,omitempty"`
	SystemInfo *JSONB `gorm:"type:jsonb" json:"system_info,omitempty"`

	// 网络信息
	IPAddress    *net.IP `gorm:"type:inet" json:"ip_address,omitempty"`
	ResponseTime *int    `gorm:"comment:响应时间(毫秒)" json:"response_time,omitempty"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (DeviceHeartbeat) TableName() string {
	return "device_heartbeats"
}

// DeviceLog 设备日志模型
type DeviceLog struct {
	ID       int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID int64     `gorm:"not null;index" json:"device_id"`
	LogTime  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index" json:"log_time"`

	// 日志分类
	Level    string  `gorm:"type:varchar(10);not null;index" json:"level"`
	Category *string `gorm:"type:varchar(50);index" json:"category,omitempty"`

	// 日志内容
	Message string `gorm:"type:text;not null" json:"message"`
	Details *JSONB `gorm:"type:jsonb" json:"details,omitempty"`

	// 上下文信息
	Source        *string `gorm:"type:varchar(100)" json:"source,omitempty"`
	CorrelationID *string `gorm:"type:varchar(100)" json:"correlation_id,omitempty"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (DeviceLog) TableName() string {
	return "device_logs"
}

// DeviceCommand 设备命令模型
type DeviceCommand struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  int64  `gorm:"not null;index" json:"device_id"`
	CommandID string `gorm:"type:varchar(100);uniqueIndex;not null" json:"command_id"`

	// 命令信息
	CommandType string `gorm:"type:varchar(50);not null" json:"command_type"`
	CommandData JSONB  `gorm:"type:jsonb;not null" json:"command_data"`

	// 执行状态
	Status        string     `gorm:"type:varchar(20);not null;default:pending" json:"status"`
	SentTime      time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"sent_time"`
	ExecutedTime  *time.Time `gorm:"type:timestamp" json:"executed_time,omitempty"`
	CompletedTime *time.Time `gorm:"type:timestamp" json:"completed_time,omitempty"`

	// 结果信息
	ResponseData *JSONB  `gorm:"type:jsonb" json:"response_data,omitempty"`
	ErrorMessage *string `gorm:"type:text" json:"error_message,omitempty"`

	// 审计信息
	CreatedBy *int64 `gorm:"index" json:"created_by,omitempty"`

	// 关联关系
	Device *Device `gorm:"foreignKey:DeviceID;references:ID" json:"device,omitempty"`
}

// TableName 指定表名
func (DeviceCommand) TableName() string {
	return "device_commands"
}

// BeforeCreate GORM钩子，创建前自动生成命令ID
func (dc *DeviceCommand) BeforeCreate(tx *gorm.DB) (err error) {
	if dc.CommandID == "" {
		dc.CommandID = generateCommandID()
	}
	return
}

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
	AuthTypeCertificate = "certificate"
)

// 通信协议常量
const (
	ProtocolHTTP = "HTTP"
	ProtocolMQTT = "MQTT"
	ProtocolTCP  = "TCP"
	ProtocolUDP  = "UDP"
)

// 日志级别常量
const (
	LogLevelDEBUG = "DEBUG"
	LogLevelINFO  = "INFO"
	LogLevelWARN  = "WARN"
	LogLevelERROR = "ERROR"
)

// 命令状态常量
const (
	CommandStatusPending   = "pending"
	CommandStatusSent      = "sent"
	CommandStatusExecuted  = "executed"
	CommandStatusCompleted = "completed"
	CommandStatusFailed    = "failed"
	CommandStatusTimeout   = "timeout"
)

// generateCommandID 生成唯一的命令ID
func generateCommandID() string {
	return "cmd_" + utils.GenerateExecutionID()
}

// DeviceFilter 设备查询过滤器
type DeviceFilter struct {
	Status       []string   `json:"status,omitempty"`
	Type         []string   `json:"type,omitempty"`
	Manufacturer *string    `json:"manufacturer,omitempty"`
	Model        *string    `json:"model,omitempty"`
	Protocol     *string    `json:"protocol,omitempty"`
	MinHealth    *int       `json:"min_health,omitempty"`
	MaxHealth    *int       `json:"max_health,omitempty"`
	LastSeenFrom *time.Time `json:"last_seen_from,omitempty"`
	LastSeenTo   *time.Time `json:"last_seen_to,omitempty"`
	Keyword      *string    `json:"keyword,omitempty"` // 搜索名称和设备ID
}

// DeviceSortOption 设备排序选项
type DeviceSortOption struct {
	Field string `json:"field"` // id, name, status, health_score, last_seen, created_at
	Order string `json:"order"` // asc, desc
}

// DeviceStatistics 设备统计信息
type DeviceStatistics struct {
	Total        int64            `json:"total"`
	ByStatus     map[string]int64 `json:"by_status"`
	ByType       map[string]int64 `json:"by_type"`
	Online       int64            `json:"online"`
	Offline      int64            `json:"offline"`
	AvgHealth    float64          `json:"avg_health"`
	RecentActive int64            `json:"recent_active"` // 最近24小时活跃设备
}

// DeviceCreateRequest 创建设备请求
type DeviceCreateRequest struct {
	DeviceID     string                 `json:"device_id" binding:"required"`
	Name         string                 `json:"name" binding:"required"`
	Type         string                 `json:"type" binding:"required"`
	Model        *string                `json:"model,omitempty"`
	Manufacturer *string                `json:"manufacturer,omitempty"`
	IPAddress    *string                `json:"ip_address,omitempty"`
	Port         *int                   `json:"port,omitempty"`
	Protocol     *string                `json:"protocol,omitempty"`
	Endpoint     *string                `json:"endpoint,omitempty"`
	AuthType     *string                `json:"auth_type,omitempty"`
	Credentials  map[string]interface{} `json:"credentials,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceUpdateRequest 更新设备请求
type DeviceUpdateRequest struct {
	Name         *string                `json:"name,omitempty"`
	Model        *string                `json:"model,omitempty"`
	Manufacturer *string                `json:"manufacturer,omitempty"`
	IPAddress    *string                `json:"ip_address,omitempty"`
	Port         *int                   `json:"port,omitempty"`
	Protocol     *string                `json:"protocol,omitempty"`
	Endpoint     *string                `json:"endpoint,omitempty"`
	AuthType     *string                `json:"auth_type,omitempty"`
	Credentials  map[string]interface{} `json:"credentials,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// DeviceCommandRequest 发送设备命令请求
type DeviceCommandRequest struct {
	CommandType string                 `json:"command_type" binding:"required"`
	CommandData map[string]interface{} `json:"command_data" binding:"required"`
}

// DeviceHeartbeatRequest 设备心跳请求
type DeviceHeartbeatRequest struct {
	Status       string                 `json:"status" binding:"required"`
	Metrics      map[string]interface{} `json:"metrics,omitempty"`
	SystemInfo   map[string]interface{} `json:"system_info,omitempty"`
	ResponseTime *int                   `json:"response_time,omitempty"`
}

// DeviceLogRequest 设备日志请求
type DeviceLogRequest struct {
	Level         string                 `json:"level" binding:"required"`
	Category      *string                `json:"category,omitempty"`
	Message       string                 `json:"message" binding:"required"`
	Details       map[string]interface{} `json:"details,omitempty"`
	Source        *string                `json:"source,omitempty"`
	CorrelationID *string                `json:"correlation_id,omitempty"`
}

// DeviceListResponse 设备列表响应
type DeviceListResponse struct {
	Devices    []Device          `json:"devices"`
	Pagination PaginationInfo    `json:"pagination"`
	Statistics *DeviceStatistics `json:"statistics,omitempty"`
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page int `json:"page"` // 页码，从1开始
	Size int `json:"size"` // 每页数量
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// DeviceStatusResponse 设备状态响应
type DeviceStatusResponse struct {
	DeviceID    string     `json:"device_id"`
	Status      string     `json:"status"`
	HealthScore int        `json:"health_score"`
	LastSeen    *time.Time `json:"last_seen,omitempty"`
	Uptime      float64    `json:"uptime_hours"`
	Online      bool       `json:"online"`
}

// LogFilter 日志过滤器
type LogFilter struct {
	Level     *string    `json:"level,omitempty"`
	Category  *string    `json:"category,omitempty"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Keyword   *string    `json:"keyword,omitempty"`
}
