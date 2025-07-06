package v1

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 通用数据模型定义
// ===============================

// UserInfo 用户基本信息
type UserInfo struct {
	UserId     string      `json:"user_id"`
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	Phone      string      `json:"phone"`
	RealName   string      `json:"real_name"`
	Avatar     string      `json:"avatar"`
	Department string      `json:"department"`
	Position   string      `json:"position"`
	Status     string      `json:"status"`
	IsVerified bool        `json:"is_verified"`
	LastLogin  *gtime.Time `json:"last_login"`
	LoginCount int         `json:"login_count"`
	CreatedAt  *gtime.Time `json:"created_at"`
	UpdatedAt  *gtime.Time `json:"updated_at"`
}

// UserDetailInfo 用户详细信息
type UserDetailInfo struct {
	UserInfo
	Roles               []RoleInfo       `json:"roles"`
	Permissions         []PermissionInfo `json:"permissions"`
	SecuritySettings    SecuritySettings `json:"security_settings"`
	ActivitySummary     ActivitySummary  `json:"activity_summary"`
	SessionCount        int              `json:"session_count"`
	LastPasswordChange  *gtime.Time      `json:"last_password_change"`
	AccountLockedUntil  *gtime.Time      `json:"account_locked_until"`
	FailedLoginAttempts int              `json:"failed_login_attempts"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	RoleId      string      `json:"role_id"`
	RoleName    string      `json:"role_name"`
	RoleDesc    string      `json:"role_desc"`
	RoleType    string      `json:"role_type"`
	Permissions []string    `json:"permissions"`
	AssignedAt  *gtime.Time `json:"assigned_at"`
	AssignedBy  string      `json:"assigned_by"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	PermissionId   string      `json:"permission_id"`
	PermissionName string      `json:"permission_name"`
	Resource       string      `json:"resource"`
	Actions        []string    `json:"actions"`
	Scope          string      `json:"scope"`
	GrantedBy      string      `json:"granted_by"`
	GrantedAt      *gtime.Time `json:"granted_at"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SessionId    string      `json:"session_id"`
	UserId       string      `json:"user_id"`
	DeviceInfo   string      `json:"device_info"`
	IpAddress    string      `json:"ip_address"`
	UserAgent    string      `json:"user_agent"`
	LoginTime    *gtime.Time `json:"login_time"`
	LastActivity *gtime.Time `json:"last_activity"`
	ExpiresAt    *gtime.Time `json:"expires_at"`
	Status       string      `json:"status"`
	Location     string      `json:"location"`
}

// ActivityInfo 活动信息
type ActivityInfo struct {
	ActivityId  string                 `json:"activity_id"`
	UserId      string                 `json:"user_id"`
	ActionType  string                 `json:"action_type"`
	Resource    string                 `json:"resource"`
	ResourceId  string                 `json:"resource_id"`
	Description string                 `json:"description"`
	IpAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   *gtime.Time            `json:"timestamp"`
	Result      string                 `json:"result"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SecurityLogInfo 安全日志信息
type SecurityLogInfo struct {
	LogId     string      `json:"log_id"`
	UserId    string      `json:"user_id"`
	LogType   string      `json:"log_type"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	IpAddress string      `json:"ip_address"`
	UserAgent string      `json:"user_agent"`
	Timestamp *gtime.Time `json:"timestamp"`
	RiskScore float64     `json:"risk_score"`
	Handled   bool        `json:"handled"`
}

// SecuritySettings 安全设置
type SecuritySettings struct {
	EnableTwoFactor         bool `json:"enable_two_factor"`
	EnableEmailNotification bool `json:"enable_email_notification"`
	EnableSmsNotification   bool `json:"enable_sms_notification"`
	SessionTimeout          int  `json:"session_timeout"`
	MaxConcurrentSessions   int  `json:"max_concurrent_sessions"`
	PasswordComplexity      bool `json:"password_complexity"`
	ForcePasswordChange     bool `json:"force_password_change"`
}

// ActivitySummary 活动摘要
type ActivitySummary struct {
	TotalActivities    int     `json:"total_activities"`
	RecentActivities   int     `json:"recent_activities"`
	MostActiveHour     int     `json:"most_active_hour"`
	MostUsedFeature    string  `json:"most_used_feature"`
	AverageSessionTime float64 `json:"average_session_time"`
	LastActivityType   string  `json:"last_activity_type"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	TotalUsers             int                `json:"total_users"`
	ActiveUsers            int                `json:"active_users"`
	InactiveUsers          int                `json:"inactive_users"`
	NewUsers               int                `json:"new_users"`
	StatusDistribution     map[string]int     `json:"status_distribution"`
	RoleDistribution       map[string]int     `json:"role_distribution"`
	DepartmentDistribution map[string]int     `json:"department_distribution"`
	RegistrationTrend      []UserTrendPoint   `json:"registration_trend"`
	LoginTrend             []UserTrendPoint   `json:"login_trend"`
	TopActiveUsers         []UserActivityRank `json:"top_active_users"`
}

// UserTrendPoint 用户趋势数据点
type UserTrendPoint struct {
	Date      string      `json:"date"`
	Count     int         `json:"count"`
	Timestamp *gtime.Time `json:"timestamp"`
}

// UserActivityRank 用户活动排名
type UserActivityRank struct {
	UserId        string      `json:"user_id"`
	Username      string      `json:"username"`
	ActivityCount int         `json:"activity_count"`
	LastActivity  *gtime.Time `json:"last_activity"`
}

// UserBehaviorAnalysis 用户行为分析
type UserBehaviorAnalysis struct {
	UserId              string             `json:"user_id"`
	TotalSessions       int                `json:"total_sessions"`
	AverageSessionTime  float64            `json:"average_session_time"`
	MostActiveTimeSlots []TimeSlotActivity `json:"most_active_time_slots"`
	FeatureUsage        map[string]int     `json:"feature_usage"`
	DeviceUsage         map[string]int     `json:"device_usage"`
	LocationAnalysis    map[string]int     `json:"location_analysis"`
	BehaviorPattern     string             `json:"behavior_pattern"`
	RiskScore           float64            `json:"risk_score"`
	Recommendations     []string           `json:"recommendations"`
}

// TimeSlotActivity 时间段活动
type TimeSlotActivity struct {
	Hour          int `json:"hour"`
	ActivityCount int `json:"activity_count"`
}
