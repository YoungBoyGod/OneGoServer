package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 1. 用户基础管理 API (6个)
// ===============================

// RegisterUserReq 用户注册请求
type RegisterUserReq struct {
	g.Meta     `path:"/user/register" method:"post" tags:"用户管理" summary:"用户注册"`
	Username   string `json:"username" v:"required|length:3,20|regex:^[a-zA-Z0-9_]+$#用户名不能为空|用户名长度为3-20字符|用户名只能包含字母数字下划线"`
	Password   string `json:"password" v:"required|length:6,32#密码不能为空|密码长度为6-32字符"`
	Email      string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确"`
	Phone      string `json:"phone,omitempty" v:"phone#手机号格式不正确"`
	RealName   string `json:"real_name,omitempty" v:"length:1,50#真实姓名长度为1-50字符"`
	Avatar     string `json:"avatar,omitempty" v:"url#头像URL格式不正确"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Position   string `json:"position,omitempty" v:"length:1,50#职位名称长度为1-50字符"`
}

type RegisterUserRes struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

// LoginUserReq 用户登录请求
type LoginUserReq struct {
	g.Meta     `path:"/user/login" method:"post" tags:"用户管理" summary:"用户登录"`
	Username   string `json:"username" v:"required|length:3,50#用户名不能为空"`
	Password   string `json:"password" v:"required|length:6,32#密码不能为空"`
	LoginType  string `json:"login_type" d:"password" v:"in:password,email,phone#登录方式无效"`
	DeviceInfo string `json:"device_info,omitempty" v:"max-length:200#设备信息最大200字符"`
	IpAddress  string `json:"ip_address,omitempty" v:"ip#IP地址格式不正确"`
}

type LoginUserRes struct {
	UserId       string   `json:"user_id"`
	Username     string   `json:"username"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Permissions  []string `json:"permissions"`
	Roles        []string `json:"roles"`
	UserInfo     UserInfo `json:"user_info"`
}

// LogoutUserReq 用户登出请求
type LogoutUserReq struct {
	g.Meta     `path:"/user/logout" method:"post" tags:"用户管理" summary:"用户登出"`
	UserId     string `json:"user_id,omitempty" v:"max-length:50#用户ID最大50字符"`
	AllDevices bool   `json:"all_devices,omitempty"`
}

type LogoutUserRes struct {
	Message string `json:"message"`
}

// GetUserListReq 获取用户列表请求
type GetUserListReq struct {
	g.Meta     `path:"/user/list" method:"get" tags:"用户管理" summary:"获取用户列表"`
	Page       int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size       int         `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status     string      `json:"status,omitempty" v:"in:active,inactive,locked,pending#状态值无效"`
	Role       string      `json:"role,omitempty" v:"length:1,50#角色名称长度为1-50字符"`
	Department string      `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Keyword    string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	StartTime  *gtime.Time `json:"start_time,omitempty"`
	EndTime    *gtime.Time `json:"end_time,omitempty"`
	SortBy     string      `json:"sort_by" d:"created_at" v:"in:created_at,updated_at,last_login#排序字段无效"`
	SortOrder  string      `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"`
}

type GetUserListRes struct {
	List  []UserInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// GetUserDetailReq 获取用户详情请求
type GetUserDetailReq struct {
	g.Meta `path:"/user/{userId}" method:"get" tags:"用户管理" summary:"获取用户详情"`
	UserId string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
}

type GetUserDetailRes struct {
	UserDetail UserDetailInfo `json:"user_detail"`
}

// UpdateUserReq 更新用户信息请求
type UpdateUserReq struct {
	g.Meta     `path:"/user/{userId}" method:"put" tags:"用户管理" summary:"更新用户信息"`
	UserId     string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Email      string `json:"email,omitempty" v:"email#邮箱格式不正确"`
	Phone      string `json:"phone,omitempty" v:"phone#手机号格式不正确"`
	RealName   string `json:"real_name,omitempty" v:"length:1,50#真实姓名长度为1-50字符"`
	Avatar     string `json:"avatar,omitempty" v:"url#头像URL格式不正确"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Position   string `json:"position,omitempty" v:"length:1,50#职位名称长度为1-50字符"`
	Status     string `json:"status,omitempty" v:"in:active,inactive,locked#状态值无效"`
}

type UpdateUserRes struct {
	Message string `json:"message"`
}

// ===============================
// 2. 用户权限管理 API (5个)
// ===============================

// GetUserRolesReq 获取用户角色请求
type GetUserRolesReq struct {
	g.Meta `path:"/user/{userId}/roles" method:"get" tags:"用户权限" summary:"获取用户角色"`
	UserId string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
}

type GetUserRolesRes struct {
	Roles []RoleInfo `json:"roles"`
}

// AssignUserRoleReq 分配用户角色请求
type AssignUserRoleReq struct {
	g.Meta  `path:"/user/{userId}/roles" method:"post" tags:"用户权限" summary:"分配用户角色"`
	UserId  string   `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	RoleIds []string `json:"role_ids" v:"required|max:10#角色ID列表不能为空|最多10个角色"`
	Reason  string   `json:"reason,omitempty" v:"max-length:200#分配原因最大200字符"`
}

type AssignUserRoleRes struct {
	Message string `json:"message"`
}

// RemoveUserRoleReq 移除用户角色请求
type RemoveUserRoleReq struct {
	g.Meta  `path:"/user/{userId}/roles/remove" method:"post" tags:"用户权限" summary:"移除用户角色"`
	UserId  string   `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	RoleIds []string `json:"role_ids" v:"required|max:10#角色ID列表不能为空|最多10个角色"`
	Reason  string   `json:"reason,omitempty" v:"max-length:200#移除原因最大200字符"`
}

type RemoveUserRoleRes struct {
	Message string `json:"message"`
}

// GetUserPermissionsReq 获取用户权限请求
type GetUserPermissionsReq struct {
	g.Meta   `path:"/user/{userId}/permissions" method:"get" tags:"用户权限" summary:"获取用户权限"`
	UserId   string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Resource string `json:"resource,omitempty" v:"length:1,100#资源名称长度为1-100字符"`
}

type GetUserPermissionsRes struct {
	Permissions []PermissionInfo `json:"permissions"`
}

// CheckUserPermissionReq 检查用户权限请求
type CheckUserPermissionReq struct {
	g.Meta     `path:"/user/{userId}/permission/check" method:"post" tags:"用户权限" summary:"检查用户权限"`
	UserId     string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Resource   string `json:"resource" v:"required|length:1,100#资源名称不能为空"`
	Action     string `json:"action" v:"required|length:1,50#操作名称不能为空"`
	ResourceId string `json:"resource_id,omitempty" v:"max-length:50#资源ID最大50字符"`
}

type CheckUserPermissionRes struct {
	HasPermission bool   `json:"has_permission"`
	Reason        string `json:"reason"`
}

// ===============================
// 3. 用户会话管理 API (4个)
// ===============================

// GetUserSessionsReq 获取用户会话请求
type GetUserSessionsReq struct {
	g.Meta `path:"/user/{userId}/sessions" method:"get" tags:"用户会话" summary:"获取用户会话"`
	UserId string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Page   int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size   int    `json:"size" d:"10" v:"between:1,50#每页数量为1-50"`
}

type GetUserSessionsRes struct {
	List  []SessionInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// RefreshTokenReq 刷新令牌请求
type RefreshTokenReq struct {
	g.Meta       `path:"/user/token/refresh" method:"post" tags:"用户会话" summary:"刷新令牌"`
	RefreshToken string `json:"refresh_token" v:"required|length:1,500#刷新令牌不能为空"`
	DeviceInfo   string `json:"device_info,omitempty" v:"max-length:200#设备信息最大200字符"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RevokeUserSessionReq 注销用户会话请求
type RevokeUserSessionReq struct {
	g.Meta    `path:"/user/{userId}/session/{sessionId}/revoke" method:"post" tags:"用户会话" summary:"注销用户会话"`
	UserId    string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	SessionId string `json:"session_id" v:"required|max-length:50#会话ID不能为空"`
	Reason    string `json:"reason,omitempty" v:"max-length:200#注销原因最大200字符"`
}

type RevokeUserSessionRes struct {
	Message string `json:"message"`
}

// GetUserActivityReq 获取用户活动请求
type GetUserActivityReq struct {
	g.Meta     `path:"/user/{userId}/activity" method:"get" tags:"用户会话" summary:"获取用户活动"`
	UserId     string      `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Page       int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size       int         `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	ActionType string      `json:"action_type,omitempty" v:"in:login,logout,create,update,delete,view#操作类型无效"`
	StartTime  *gtime.Time `json:"start_time,omitempty"`
	EndTime    *gtime.Time `json:"end_time,omitempty"`
}

type GetUserActivityRes struct {
	List  []ActivityInfo `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// ===============================
// 4. 用户安全管理 API (4个)
// ===============================

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	g.Meta          `path:"/user/{userId}/password" method:"put" tags:"用户安全" summary:"修改密码"`
	UserId          string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	OldPassword     string `json:"old_password" v:"required|length:6,32#原密码不能为空"`
	NewPassword     string `json:"new_password" v:"required|length:6,32#新密码不能为空"`
	ConfirmPassword string `json:"confirm_password" v:"required|length:6,32#确认密码不能为空"`
}

type ChangePasswordRes struct {
	Message string `json:"message"`
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	g.Meta      `path:"/user/password/reset" method:"post" tags:"用户安全" summary:"重置密码"`
	Username    string `json:"username" v:"required|length:3,50#用户名不能为空"`
	Email       string `json:"email" v:"required|email#邮箱不能为空"`
	Code        string `json:"code,omitempty" v:"length:4,10#验证码长度为4-10字符"`
	NewPassword string `json:"new_password,omitempty" v:"length:6,32#新密码长度为6-32字符"`
}

type ResetPasswordRes struct {
	Message  string `json:"message"`
	CodeSent bool   `json:"code_sent"`
}

// GetUserSecurityLogReq 获取用户安全日志请求
type GetUserSecurityLogReq struct {
	g.Meta    `path:"/user/{userId}/security/logs" method:"get" tags:"用户安全" summary:"获取用户安全日志"`
	UserId    string      `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	LogType   string      `json:"log_type,omitempty" v:"in:login,password_change,permission_change,security_event#日志类型无效"`
	Level     string      `json:"level,omitempty" v:"in:info,warn,error,critical#日志级别无效"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
}

type GetUserSecurityLogRes struct {
	List  []SecurityLogInfo `json:"list"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

// UpdateUserSecuritySettingsReq 更新用户安全设置请求
type UpdateUserSecuritySettingsReq struct {
	g.Meta                  `path:"/user/{userId}/security/settings" method:"put" tags:"用户安全" summary:"更新用户安全设置"`
	UserId                  string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	EnableTwoFactor         bool   `json:"enable_two_factor,omitempty"`
	EnableEmailNotification bool   `json:"enable_email_notification,omitempty"`
	EnableSmsNotification   bool   `json:"enable_sms_notification,omitempty"`
	SessionTimeout          int    `json:"session_timeout,omitempty" v:"between:300,86400#会话超时时间为300-86400秒"`
	MaxConcurrentSessions   int    `json:"max_concurrent_sessions,omitempty" v:"between:1,10#最大并发会话为1-10个"`
}

type UpdateUserSecuritySettingsRes struct {
	Message string `json:"message"`
}

// ===============================
// 5. 用户统计分析 API (2个)
// ===============================

// GetUserStatisticsReq 获取用户统计信息请求
type GetUserStatisticsReq struct {
	g.Meta     `path:"/user/statistics" method:"get" tags:"用户统计" summary:"获取用户统计信息"`
	TimeRange  string `json:"time_range" d:"30d" v:"in:24h,7d,30d,90d#时间范围无效"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Role       string `json:"role,omitempty" v:"length:1,50#角色名称长度为1-50字符"`
	GroupBy    string `json:"group_by" d:"status" v:"in:status,role,department,registration_date#分组字段无效"`
}

type GetUserStatisticsRes struct {
	Statistics UserStatistics `json:"statistics"`
}

// GetUserBehaviorAnalysisReq 获取用户行为分析请求
type GetUserBehaviorAnalysisReq struct {
	g.Meta    `path:"/user/{userId}/behavior" method:"get" tags:"用户统计" summary:"获取用户行为分析"`
	UserId    string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	TimeRange string `json:"time_range" d:"30d" v:"in:7d,30d,90d#时间范围无效"`
}

type GetUserBehaviorAnalysisRes struct {
	BehaviorAnalysis UserBehaviorAnalysis `json:"behavior_analysis"`
}

// ===============================
// 6. 用户数据结构定义
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
