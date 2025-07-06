package user

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// User 用户实体
type User struct {
	ID                   string      `json:"id" db:"id"`
	Username             string      `json:"username" db:"username"`
	Email                string      `json:"email" db:"email"`
	Phone                string      `json:"phone" db:"phone"`
	Password             string      `json:"-" db:"password"`
	Status               string      `json:"status" db:"status"`
	Role                 string      `json:"role" db:"role"`
	Avatar               string      `json:"avatar" db:"avatar"`
	Nickname             string      `json:"nickname" db:"nickname"`
	RealName             string      `json:"real_name" db:"real_name"`
	Gender               string      `json:"gender" db:"gender"`
	Birthday             *gtime.Time `json:"birthday" db:"birthday"`
	Department           string      `json:"department" db:"department"`
	Position             string      `json:"position" db:"position"`
	EmployeeID           string      `json:"employee_id" db:"employee_id"`
	LastLoginAt          *gtime.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIP          string      `json:"last_login_ip" db:"last_login_ip"`
	LastUserAgent        string      `json:"last_user_agent" db:"last_user_agent"`
	FailedLoginAttempts  int         `json:"failed_login_attempts" db:"failed_login_attempts"`
	AccountLockedUntil   *gtime.Time `json:"account_locked_until" db:"account_locked_until"`
	ForcePasswordChange  bool        `json:"force_password_change" db:"force_password_change"`
	PasswordChangedAt    *gtime.Time `json:"password_changed_at" db:"password_changed_at"`
	RecentLoginCount     int         `json:"recent_login_count" db:"recent_login_count"`
	RecentSecurityEvents int         `json:"recent_security_events" db:"recent_security_events"`
	CreatedAt            *gtime.Time `json:"created_at" db:"created_at"`
	UpdatedAt            *gtime.Time `json:"updated_at" db:"updated_at"`
	DeletedAt            *gtime.Time `json:"deleted_at" db:"deleted_at"`
}

// UserSession 用户会话实体
type UserSession struct {
	ID           string      `json:"id" db:"id"`
	UserID       string      `json:"user_id" db:"user_id"`
	SessionID    string      `json:"session_id" db:"session_id"`
	IPAddress    string      `json:"ip_address" db:"ip_address"`
	UserAgent    string      `json:"user_agent" db:"user_agent"`
	DeviceInfo   string      `json:"device_info" db:"device_info"`
	Status       string      `json:"status" db:"status"`
	CreatedAt    *gtime.Time `json:"created_at" db:"created_at"`
	LastActivity *gtime.Time `json:"last_activity" db:"last_activity"`
	ExpiresAt    *gtime.Time `json:"expires_at" db:"expires_at"`
}

// UserRole 用户角色实体
type UserRole struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	RoleID      string      `json:"role_id" db:"role_id"`
	RoleName    string      `json:"role_name" db:"role_name"`
	Permissions []string    `json:"permissions" db:"permissions"`
	CreatedAt   *gtime.Time `json:"created_at" db:"created_at"`
	UpdatedAt   *gtime.Time `json:"updated_at" db:"updated_at"`
}

// UserActivity 用户活动实体
type UserActivity struct {
	ID              string      `json:"id" db:"id"`
	UserID          string      `json:"user_id" db:"user_id"`
	ActionType      string      `json:"action_type" db:"action_type"`
	ActionDetail    string      `json:"action_detail" db:"action_detail"`
	IPAddress       string      `json:"ip_address" db:"ip_address"`
	UserAgent       string      `json:"user_agent" db:"user_agent"`
	SessionID       string      `json:"session_id" db:"session_id"`
	ResourceType    string      `json:"resource_type" db:"resource_type"`
	ResourceID      string      `json:"resource_id" db:"resource_id"`
	SessionDuration float64     `json:"session_duration" db:"session_duration"`
	CreatedAt       *gtime.Time `json:"created_at" db:"created_at"`
}

// UserSecurityLog 用户安全日志实体
type UserSecurityLog struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	EventType   string      `json:"event_type" db:"event_type"`
	EventDetail string      `json:"event_detail" db:"event_detail"`
	IPAddress   string      `json:"ip_address" db:"ip_address"`
	UserAgent   string      `json:"user_agent" db:"user_agent"`
	RiskScore   float64     `json:"risk_score" db:"risk_score"`
	RiskLevel   string      `json:"risk_level" db:"risk_level"`
	CreatedAt   *gtime.Time `json:"created_at" db:"created_at"`
}

// UserFilter 用户过滤条件
type UserFilter struct {
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Status         string     `json:"status"`
	Role           string     `json:"role"`
	Department     string     `json:"department"`
	Position       string     `json:"position"`
	EmployeeID     string     `json:"employee_id"`
	CreatedAtStart *time.Time `json:"created_at_start"`
	CreatedAtEnd   *time.Time `json:"created_at_end"`
}

// UserSortOption 用户排序选项
type UserSortOption struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // asc, desc
}
