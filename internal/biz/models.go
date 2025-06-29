package biz

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"unique;not null;size:50" json:"username"`
	Email       string         `gorm:"unique;not null;size:100" json:"email"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Password    string         `gorm:"not null;size:255" json:"-"` // 不在JSON中返回密码
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Status      UserStatus     `gorm:"default:1" json:"status"`
	Role        UserRole       `gorm:"default:2" json:"role"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Devices []Device `gorm:"foreignKey:UserID" json:"devices,omitempty"`
	Tasks   []Task   `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
}

// UserStatus 用户状态枚举
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 活跃
	UserStatusLocked   UserStatus = 2 // 锁定
	UserStatusBanned   UserStatus = 3 // 禁用
)

// UserRole 用户角色枚举
type UserRole int

const (
	UserRoleAdmin UserRole = 1 // 管理员
	UserRoleUser  UserRole = 2 // 普通用户
	UserRoleGuest UserRole = 3 // 访客
)

// Device 设备实体
type Device struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	DeviceID    string         `gorm:"unique;not null;size:100" json:"device_id"` // 设备唯一标识
	Name        string         `gorm:"not null;size:100" json:"name"`
	Type        string         `gorm:"size:50" json:"type"`
	Description string         `gorm:"size:500" json:"description"`
	Status      DeviceStatus   `gorm:"default:1" json:"status"`
	IP          string         `gorm:"size:45" json:"ip"`
	MAC         string         `gorm:"size:17" json:"mac"`
	OS          string         `gorm:"size:100" json:"os"`
	Version     string         `gorm:"size:50" json:"version"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	LastSeenAt  *time.Time     `json:"last_seen_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	User   User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Tasks  []Task  `gorm:"foreignKey:DeviceID" json:"tasks,omitempty"`
	Alerts []Alert `gorm:"foreignKey:DeviceID" json:"alerts,omitempty"`
}

// DeviceStatus 设备状态枚举
type DeviceStatus int

const (
	DeviceStatusOffline     DeviceStatus = 0 // 离线
	DeviceStatusOnline      DeviceStatus = 1 // 在线
	DeviceStatusMaintenance DeviceStatus = 2 // 维护中
	DeviceStatusError       DeviceStatus = 3 // 错误
)

// Task 任务实体
type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"not null;size:200" json:"name"`
	Description string         `gorm:"size:1000" json:"description"`
	Type        TaskType       `gorm:"not null" json:"type"`
	Status      TaskStatus     `gorm:"default:1" json:"status"`
	Priority    TaskPriority   `gorm:"default:2" json:"priority"`
	Command     string         `gorm:"size:2000" json:"command"`
	Parameters  string         `gorm:"type:text" json:"parameters"` // JSON格式参数
	Result      string         `gorm:"type:text" json:"result"`     // 执行结果
	ErrorMsg    string         `gorm:"size:1000" json:"error_msg"`
	RetryCount  int            `gorm:"default:0" json:"retry_count"`
	MaxRetries  int            `gorm:"default:3" json:"max_retries"`
	Timeout     int            `gorm:"default:3600" json:"timeout"` // 超时时间(秒)
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	DeviceID    uint           `gorm:"not null;index" json:"device_id"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Device Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

// TaskType 任务类型枚举
type TaskType int

const (
	TaskTypeCommand TaskType = 1 // 命令执行
	TaskTypeScript  TaskType = 2 // 脚本执行
	TaskTypeFile    TaskType = 3 // 文件操作
	TaskTypeSystem  TaskType = 4 // 系统操作
)

// TaskStatus 任务状态枚举
type TaskStatus int

const (
	TaskStatusPending   TaskStatus = 1 // 待执行
	TaskStatusRunning   TaskStatus = 2 // 执行中
	TaskStatusCompleted TaskStatus = 3 // 已完成
	TaskStatusFailed    TaskStatus = 4 // 失败
	TaskStatusCancelled TaskStatus = 5 // 已取消
	TaskStatusTimeout   TaskStatus = 6 // 超时
)

// TaskPriority 任务优先级枚举
type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 1 // 低优先级
	TaskPriorityNormal TaskPriority = 2 // 普通优先级
	TaskPriorityHigh   TaskPriority = 3 // 高优先级
	TaskPriorityUrgent TaskPriority = 4 // 紧急优先级
)

// Alert 告警实体
type Alert struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Title      string         `gorm:"not null;size:200" json:"title"`
	Message    string         `gorm:"not null;size:1000" json:"message"`
	Type       AlertType      `gorm:"not null" json:"type"`
	Level      AlertLevel     `gorm:"not null" json:"level"`
	Status     AlertStatus    `gorm:"default:1" json:"status"`
	Source     string         `gorm:"size:100" json:"source"`    // 告警来源
	Metadata   string         `gorm:"type:text" json:"metadata"` // JSON格式元数据
	DeviceID   uint           `gorm:"not null;index" json:"device_id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	AckedBy    *uint          `gorm:"index" json:"acked_by"` // 确认人
	AckedAt    *time.Time     `json:"acked_at"`              // 确认时间
	ResolvedAt *time.Time     `json:"resolved_at"`           // 解决时间
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联关系
	Device  Device `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AckUser *User  `gorm:"foreignKey:AckedBy" json:"ack_user,omitempty"`
}

// AlertType 告警类型枚举
type AlertType int

const (
	AlertTypeSystem      AlertType = 1 // 系统告警
	AlertTypeDevice      AlertType = 2 // 设备告警
	AlertTypeTask        AlertType = 3 // 任务告警
	AlertTypeSecurity    AlertType = 4 // 安全告警
	AlertTypePerformance AlertType = 5 // 性能告警
)

// AlertLevel 告警级别枚举
type AlertLevel int

const (
	AlertLevelInfo     AlertLevel = 1 // 信息
	AlertLevelWarning  AlertLevel = 2 // 警告
	AlertLevelError    AlertLevel = 3 // 错误
	AlertLevelCritical AlertLevel = 4 // 严重
)

// AlertStatus 告警状态枚举
type AlertStatus int

const (
	AlertStatusOpen         AlertStatus = 1 // 开放
	AlertStatusAcknowledged AlertStatus = 2 // 已确认
	AlertStatusResolved     AlertStatus = 3 // 已解决
	AlertStatusClosed       AlertStatus = 4 // 已关闭
)

// TableName 指定表名
func (User) TableName() string   { return "users" }
func (Device) TableName() string { return "devices" }
func (Task) TableName() string   { return "tasks" }
func (Alert) TableName() string  { return "alerts" }
