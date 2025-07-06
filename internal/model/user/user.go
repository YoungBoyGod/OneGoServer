package user

import (
	"time"

	"OneGfServer/internal/consts"
	"OneGfServer/internal/model/common"

	"github.com/gogf/gf/v2/os/gtime"
)

// 用户状态常量 - 使用统一常量
const (
	UserStatusActive   = consts.UserStatusActive
	UserStatusInactive = consts.UserStatusInactive
	UserStatusLocked   = consts.UserStatusLocked
	UserStatusPending  = consts.UserStatusPending
)

// 用户角色常量 - 使用统一常量
const (
	UserRoleAdmin    = consts.UserRoleAdmin
	UserRoleManager  = consts.UserRoleManager
	UserRoleOperator = consts.UserRoleOperator
	UserRoleViewer   = consts.UserRoleViewer
)

// 登录类型常量 - 使用统一常量
const (
	LoginTypePassword = consts.LoginTypePassword
	LoginTypeEmail    = consts.LoginTypeEmail
	LoginTypePhone    = consts.LoginTypePhone
)

// 用户操作类型常量 - 使用统一常量
const (
	ActionTypeLogin  = consts.UserActionTypeLogin
	ActionTypeLogout = consts.UserActionTypeLogout
	ActionTypeCreate = consts.UserActionTypeCreate
	ActionTypeUpdate = consts.UserActionTypeUpdate
	ActionTypeDelete = consts.UserActionTypeDelete
	ActionTypeView   = consts.UserActionTypeView
)

// 用户日志类型常量 - 使用统一常量
const (
	LogTypeLogin            = consts.UserLogTypeLogin
	LogTypePasswordChange   = consts.UserLogTypePasswordChange
	LogTypePermissionChange = consts.UserLogTypePermissionChange
	LogTypeSecurityEvent    = consts.UserLogTypeSecurityEvent
)

// 用户日志级别常量 - 使用统一常量
const (
	LogLevelInfo     = consts.UserLogLevelInfo
	LogLevelWarn     = consts.UserLogLevelWarn
	LogLevelError    = consts.UserLogLevelError
	LogLevelCritical = consts.UserLogLevelCritical
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
	Filter     *UserFilter               `json:"filter"`
	Sort       *UserSortOption           `json:"sort"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetUserListOutput 获取用户列表输出
type GetUserListOutput struct {
	common.PaginationResponse[User] `json:",inline"`
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
	UserID     string                    `json:"user_id"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetUserSessionsOutput 获取用户会话输出
type GetUserSessionsOutput struct {
	common.PaginationResponse[SessionInfo] `json:",inline"`
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
	UserID     string                    `json:"user_id"`
	ActionType string                    `json:"action_type"`
	StartTime  *gtime.Time               `json:"start_time"`
	EndTime    *gtime.Time               `json:"end_time"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetUserActivityOutput 获取用户活动输出
type GetUserActivityOutput struct {
	common.PaginationResponse[ActivityInfo] `json:",inline"`
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
	UserID     string                    `json:"user_id"`
	LogType    string                    `json:"log_type"`
	Level      string                    `json:"level"`
	StartTime  *gtime.Time               `json:"start_time"`
	EndTime    *gtime.Time               `json:"end_time"`
	Pagination *common.PaginationRequest `json:"pagination"`
}

// GetUserSecurityLogOutput 获取用户安全日志输出
type GetUserSecurityLogOutput struct {
	common.PaginationResponse[SecurityLogInfo] `json:",inline"`
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
	Message string `