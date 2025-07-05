package user

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusLocked   = "locked"
	UserStatusPending  = "pending"
)

// 用户角色常量
const (
	UserRoleAdmin    = "admin"
	UserRoleManager  = "manager"
	UserRoleOperator = "operator"
	UserRoleViewer   = "viewer"
)

// 登录类型常量
const (
	LoginTypePassword = "password"
	LoginTypeEmail    = "email"
	LoginTypePhone    = "phone"
)

// 操作类型常量
const (
	ActionTypeLogin  = "login"
	ActionTypeLogout = "logout"
	ActionTypeCreate = "create"
	ActionTypeUpdate = "update"
	ActionTypeDelete = "delete"
	ActionTypeView   = "view"
)

// 日志类型常量
const (
	LogTypeLogin            = "login"
	LogTypePasswordChange   = "password_change"
	LogTypePermissionChange = "permission_change"
	LogTypeSecurityEvent    = "security_event"
)

// 日志级别常量
const (
	LogLevelInfo     = "info"
	LogLevelWarn     = "warn"
	LogLevelError    = "error"
	LogLevelCritical = "critical"
)

// ===============================
// 基础CRUD操作 Input/Output
// ===============================

// CreateUserInput 创建用户输入
type CreateUserInput struct {
	User *User `json:"user"`
}

// CreateUserOutput 创建用户输出
type CreateUserOutput struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// GetUserByIDInput 根据ID获取用户输入
type GetUserByIDInput struct {
	ID int64 `json:"id"`
}

// GetUserByIDOutput 根据ID获取用户输出
type GetUserByIDOutput struct {
	User *User `json:"user"`
}

// GetUserByUserIDInput 根据用户ID获取用户输入
type GetUserByUserIDInput struct {
	UserID string `json:"user_id"`
}

// GetUserByUserIDOutput 根据用户ID获取用户输出
type GetUserByUserIDOutput struct {
	User *User `json:"user"`
}

// GetUserByUsernameInput 根据用户名获取用户输入
type GetUserByUsernameInput struct {
	Username string `json:"username"`
}

// GetUserByUsernameOutput 根据用户名获取用户输出
type GetUserByUsernameOutput struct {
	User *User `json:"user"`
}

// GetUserByEmailInput 根据邮箱获取用户输入
type GetUserByEmailInput struct {
	Email string `json:"email"`
}

// GetUserByEmailOutput 根据邮箱获取用户输出
type GetUserByEmailOutput struct {
	User *User `json:"user"`
}

// UpdateUserInput 更新用户输入
type UpdateUserInput struct {
	User *User `json:"user"`
}

// UpdateUserOutput 更新用户输出
type UpdateUserOutput struct {
	Message string `json:"message"`
}

// DeleteUserInput 删除用户输入
type DeleteUserInput struct {
	UserID string `json:"user_id"`
	Force  bool   `json:"force"`
}

// DeleteUserOutput 删除用户输出
type DeleteUserOutput struct {
	Message string `json:"message"`
}

// ===============================
// 查询操作 Input/Output
// ===============================

// GetUserListInput 获取用户列表输入
type GetUserListInput struct {
	Filter     *UserFilter       `json:"filter"`
	Sort       *UserSortOption   `json:"sort"`
	Pagination *PaginationOption `json:"pagination"`
}

// GetUserListOutput 获取用户列表输出
type GetUserListOutput struct {
	List  []User `json:"list"`
	Total int64  `json:"total"`
}

// GetUsersByStatusInput 根据状态获取用户输入
type GetUsersByStatusInput struct {
	Statuses []string `json:"statuses"`
}

// GetUsersByStatusOutput 根据状态获取用户输出
type GetUsersByStatusOutput struct {
	Users []User `json:"users"`
}

// GetUsersByRoleInput 根据角色获取用户输入
type GetUsersByRoleInput struct {
	Roles []string `json:"roles"`
}

// GetUsersByRoleOutput 根据角色获取用户输出
type GetUsersByRoleOutput struct {
	Users []User `json:"users"`
}

// GetUsersByDepartmentInput 根据部门获取用户输入
type GetUsersByDepartmentInput struct {
	Departments []string `json:"departments"`
}

// GetUsersByDepartmentOutput 根据部门获取用户输出
type GetUsersByDepartmentOutput struct {
	Users []User `json:"users"`
}

// ===============================
// 用户认证授权 Input/Output
// ===============================

// LoginUserInput 用户登录输入
type LoginUserInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	LoginType  string `json:"login_type"`
	DeviceInfo string `json:"device_info"`
	IPAddress  string `json:"ip_address"`
}

// LoginUserOutput 用户登录输出
type LoginUserOutput struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Permissions  []string  `json:"permissions"`
	Roles        []string  `json:"roles"`
	UserInfo     *UserInfo `json:"user_info"`
}

// LogoutUserInput 用户登出输入
type LogoutUserInput struct {
	UserID     string `json:"user_id"`
	AllDevices bool   `json:"all_devices"`
}

// LogoutUserOutput 用户登出输出
type LogoutUserOutput struct {
	Message string `json:"message"`
}

// ValidateUserLoginInput 验证用户登录输入
type ValidateUserLoginInput struct {
	Username string                 `json:"username"`
	Password string                 `json:"password"`
	UserData map[string]interface{} `json:"user_data"`
}

// ValidateUserLoginOutput 验证用户登录输出
type ValidateUserLoginOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ===============================
// 用户权限管理 Input/Output
// ===============================

// GetUserRolesInput 获取用户角色输入
type GetUserRolesInput struct {
	UserID string `json:"user_id"`
}

// GetUserRolesOutput 获取用户角色输出
type GetUserRolesOutput struct {
	Roles []RoleInfo `json:"roles"`
}

// AssignUserRoleInput 分配用户角色输入
type AssignUserRoleInput struct {
	UserID  string   `json:"user_id"`
	RoleIDs []string `json:"role_ids"`
	Reason  string   `json:"reason"`
}

// AssignUserRoleOutput 分配用户角色输出
type AssignUserRoleOutput struct {
	Message string `json:"message"`
}

// RemoveUserRoleInput 移除用户角色输入
type RemoveUserRoleInput struct {
	UserID  string   `json:"user_id"`
	RoleIDs []string `json:"role_ids"`
	Reason  string   `json:"reason"`
}

// RemoveUserRoleOutput 移除用户角色输出
type RemoveUserRoleOutput struct {
	Message string `json:"message"`
}

// GetUserPermissionsInput 获取用户权限输入
type GetUserPermissionsInput struct {
	UserID   string `json:"user_id"`
	Resource string `json:"resource"`
}

// GetUserPermissionsOutput 获取用户权限输出
type GetUserPermissionsOutput struct {
	Permissions []PermissionInfo `json:"permissions"`
}

// CheckUserPermissionInput 检查用户权限输入
type CheckUserPermissionInput struct {
	UserID     string `json:"user_id"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	ResourceID string `json:"resource_id"`
}

// CheckUserPermissionOutput 检查用户权限输出
type CheckUserPermissionOutput struct {
	HasPermission bool   `json:"has_permission"`
	Reason        string `json:"reason"`
}

// CalculateUserPermissionsInput 计算用户权限输入
type CalculateUserPermissionsInput struct {
	UserRoles []map[string]interface{} `json:"user_roles"`
}

// CalculateUserPermissionsOutput 计算用户权限输出
type CalculateUserPermissionsOutput struct {
	Permissions []string `json:"permissions"`
}

// ===============================
// 用户会话管理 Input/Output
// ===============================

// GetUserSessionsInput 获取用户会话输入
type GetUserSessionsInput struct {
	UserID string `json:"user_id"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

// GetUserSessionsOutput 获取用户会话输出
type GetUserSessionsOutput struct {
	List  []SessionInfo `json:"list"`
	Total int64         `json:"total"`
}

// RefreshTokenInput 刷新令牌输入
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
	DeviceInfo   string `json:"device_info"`
}

// RefreshTokenOutput 刷新令牌输出
type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RevokeUserSessionInput 注销用户会话输入
type RevokeUserSessionInput struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// RevokeUserSessionOutput 注销用户会话输出
type RevokeUserSessionOutput struct {
	Message string `json:"message"`
}

// GetUserActivityInput 获取用户活动输入
type GetUserActivityInput struct {
	UserID     string      `json:"user_id"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	ActionType string      `json:"action_type"`
	StartTime  *gtime.Time `json:"start_time"`
	EndTime    *gtime.Time `json:"end_time"`
}

// GetUserActivityOutput 获取用户活动输出
type GetUserActivityOutput struct {
	List  []ActivityInfo `json:"list"`
	Total int64          `json:"total"`
}

// CreateUserSessionInput 创建用户会话输入
type CreateUserSessionInput struct {
	UserID    string                 `json:"user_id"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CreateUserSessionOutput 创建用户会话输出
type CreateUserSessionOutput struct {
	SessionData map[string]interface{} `json:"session_data"`
}

// ValidateUserSessionInput 验证用户会话输入
type ValidateUserSessionInput struct {
	SessionData map[string]interface{} `json:"session_data"`
}

// ValidateUserSessionOutput 验证用户会话输出
type ValidateUserSessionOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// RefreshUserSessionInput 刷新用户会话输入
type RefreshUserSessionInput struct {
	SessionData map[string]interface{} `json:"session_data"`
}

// RefreshUserSessionOutput 刷新用户会话输出
type RefreshUserSessionOutput struct {
	SessionData map[string]interface{} `json:"session_data"`
}

// ===============================
// 用户安全管理 Input/Output
// ===============================

// ChangePasswordInput 修改密码输入
type ChangePasswordInput struct {
	UserID          string `json:"user_id"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// ChangePasswordOutput 修改密码输出
type ChangePasswordOutput struct {
	Message string `json:"message"`
}

// ResetPasswordInput 重置密码输入
type ResetPasswordInput struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordOutput 重置密码输出
type ResetPasswordOutput struct {
	Message  string `json:"message"`
	CodeSent bool   `json:"code_sent"`
}

// GetUserSecurityLogInput 获取用户安全日志输入
type GetUserSecurityLogInput struct {
	UserID    string      `json:"user_id"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	LogType   string      `json:"log_type"`
	Level     string      `json:"level"`
	StartTime *gtime.Time `json:"start_time"`
	EndTime   *gtime.Time `json:"end_time"`
}

// GetUserSecurityLogOutput 获取用户安全日志输出
type GetUserSecurityLogOutput struct {
	List  []SecurityLogInfo `json:"list"`
	Total int64             `json:"total"`
}

// UpdateUserSecuritySettingsInput 更新用户安全设置输入
type UpdateUserSecuritySettingsInput struct {
	UserID                  string `json:"user_id"`
	EnableTwoFactor         bool   `json:"enable_two_factor"`
	EnableEmailNotification bool   `json:"enable_email_notification"`
	EnableSmsNotification   bool   `json:"enable_sms_notification"`
	SessionTimeout          int    `json:"session_timeout"`
	MaxConcurrentSessions   int    `json:"max_concurrent_sessions"`
}

// UpdateUserSecuritySettingsOutput 更新用户安全设置输出
type UpdateUserSecuritySettingsOutput struct {
	Message string `json:"message"`
}

// CalculateUserRiskScoreInput 计算用户风险评分输入
type CalculateUserRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateUserRiskScoreOutput 计算用户风险评分输出
type CalculateUserRiskScoreOutput struct {
	RiskScore float64 `json:"risk_score"`
}

// ===============================
// 用户统计分析 Input/Output
// ===============================

// GetUserStatisticsInput 获取用户统计信息输入
type GetUserStatisticsInput struct {
	TimeRange  string `json:"time_range"`
	Department string `json:"department"`
	Role       string `json:"role"`
	GroupBy    string `json:"group_by"`
}

// GetUserStatisticsOutput 获取用户统计信息输出
type GetUserStatisticsOutput struct {
	Statistics *UserStatistics `json:"statistics"`
}

// GetUserBehaviorAnalysisInput 获取用户行为分析输入
type GetUserBehaviorAnalysisInput struct {
	UserID    string `json:"user_id"`
	TimeRange string `json:"time_range"`
}

// GetUserBehaviorAnalysisOutput 获取用户行为分析输出
type GetUserBehaviorAnalysisOutput struct {
	BehaviorAnalysis *UserBehaviorAnalysis `json:"behavior_analysis"`
}

// ===============================
// 用户验证 Input/Output
// ===============================

// ValidateUserRegistrationInput 验证用户注册输入
type ValidateUserRegistrationInput struct {
	UserData map[string]interface{} `json:"user_data"`
}

// ValidateUserRegistrationOutput 验证用户注册输出
type ValidateUserRegistrationOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidateUsernameInput 验证用户名输入
type ValidateUsernameInput struct {
	Username string `json:"username"`
}

// ValidateUsernameOutput 验证用户名输出
type ValidateUsernameOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidatePasswordStrengthInput 验证密码强度输入
type ValidatePasswordStrengthInput struct {
	Password string `json:"password"`
}

// ValidatePasswordStrengthOutput 验证密码强度输出
type ValidatePasswordStrengthOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidateEmailInput 验证邮箱输入
type ValidateEmailInput struct {
	Email string `json:"email"`
}

// ValidateEmailOutput 验证邮箱输出
type ValidateEmailOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidatePhoneInput 验证手机号输入
type ValidatePhoneInput struct {
	Phone string `json:"phone"`
}

// ValidatePhoneOutput 验证手机号输出
type ValidatePhoneOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ===============================
// 数据模型定义
// ===============================

// User 用户主表模型
type User struct {
	ID                  int64       `json:"id"`
	UserID              string      `json:"user_id"`
	Username            string      `json:"username"`
	Password            string      `json:"password"`
	Email               string      `json:"email"`
	Phone               string      `json:"phone"`
	RealName            string      `json:"real_name"`
	Avatar              string      `json:"avatar"`
	Department          string      `json:"department"`
	Position            string      `json:"position"`
	Status              string      `json:"status"`
	IsVerified          bool        `json:"is_verified"`
	LastLogin           *gtime.Time `json:"last_login"`
	LoginCount          int         `json:"login_count"`
	FailedLoginAttempts int         `json:"failed_login_attempts"`
	AccountLockedUntil  *gtime.Time `json:"account_locked_until"`
	LastPasswordChange  *gtime.Time `json:"last_password_change"`
	ForcePasswordChange bool        `json:"force_password_change"`
	CreatedAt           *gtime.Time `json:"created_at"`
	UpdatedAt           *gtime.Time `json:"updated_at"`
	CreatedBy           int64       `json:"created_by"`
	UpdatedBy           int64       `json:"updated_by"`
}

// UserInfo 用户基本信息
type UserInfo struct {
	UserID     string      `json:"user_id"`
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
	Roles               []RoleInfo        `json:"roles"`
	Permissions         []PermissionInfo  `json:"permissions"`
	SecuritySettings    *SecuritySettings `json:"security_settings"`
	ActivitySummary     *ActivitySummary  `json:"activity_summary"`
	SessionCount        int               `json:"session_count"`
	LastPasswordChange  *gtime.Time       `json:"last_password_change"`
	AccountLockedUntil  *gtime.Time       `json:"account_locked_until"`
	FailedLoginAttempts int               `json:"failed_login_attempts"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	RoleID      string      `json:"role_id"`
	RoleName    string      `json:"role_name"`
	RoleDesc    string      `json:"role_desc"`
	RoleType    string      `json:"role_type"`
	Permissions []string    `json:"permissions"`
	AssignedAt  *gtime.Time `json:"assigned_at"`
	AssignedBy  string      `json:"assigned_by"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	PermissionID   string      `json:"permission_id"`
	PermissionName string      `json:"permission_name"`
	Resource       string      `json:"resource"`
	Actions        []string    `json:"actions"`
	Scope          string      `json:"scope"`
	GrantedBy      string      `json:"granted_by"`
	GrantedAt      *gtime.Time `json:"granted_at"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SessionID    string      `json:"session_id"`
	UserID       string      `json:"user_id"`
	DeviceInfo   string      `json:"device_info"`
	IPAddress    string      `json:"ip_address"`
	UserAgent    string      `json:"user_agent"`
	LoginTime    *gtime.Time `json:"login_time"`
	LastActivity *gtime.Time `json:"last_activity"`
	ExpiresAt    *gtime.Time `json:"expires_at"`
	Status       string      `json:"status"`
	Location     string      `json:"location"`
}

// ActivityInfo 活动信息
type ActivityInfo struct {
	ActivityID  string                 `json:"activity_id"`
	UserID      string                 `json:"user_id"`
	ActionType  string                 `json:"action_type"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Description string                 `json:"description"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   *gtime.Time            `json:"timestamp"`
	Result      string                 `json:"result"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SecurityLogInfo 安全日志信息
type SecurityLogInfo struct {
	LogID     string      `json:"log_id"`
	UserID    string      `json:"user_id"`
	LogType   string      `json:"log_type"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	IPAddress string      `json:"ip_address"`
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
	UserID        string      `json:"user_id"`
	Username      string      `json:"username"`
	ActivityCount int         `json:"activity_count"`
	LastActivity  *gtime.Time `json:"last_activity"`
}

// UserBehaviorAnalysis 用户行为分析
type UserBehaviorAnalysis struct {
	UserID              string             `json:"user_id"`
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

// ===============================
// 过滤和排序选项
// ===============================

// UserFilter 用户过滤条件
type UserFilter struct {
	Status     []string   `json:"status"`
	Role       []string   `json:"role"`
	Department []string   `json:"department"`
	IsVerified *bool      `json:"is_verified"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	Keyword    *string    `json:"keyword"`
}

// UserSortOption 用户排序选项
type UserSortOption struct {
	Field string `json:"field"` // id, username, email, status, created_at, updated_at, last_login
	Order string `json:"order"` // asc, desc
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
