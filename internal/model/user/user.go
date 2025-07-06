package user

import (
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
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
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
	UserPermissions    []string `json:"user_permissions"`
	RequiredPermission string   `json:"required_permission"`
	Resource           string   `json:"resource"`
}

// CheckUserPermissionOutput 检查用户权限输出
type CheckUserPermissionOutput struct {
	HasPermission bool `json:"has_permission"`
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
	UserId    string                 `json:"user_id"`
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
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
}

// RefreshUserSessionInput 刷新用户会话输入
type RefreshUserSessionInput struct {
	SessionData map[string]interface{} `json:"session_data"`
}

// RefreshUserSessionOutput 刷新用户会话输出
type RefreshUserSessionOutput struct {
	UpdatedSessionData map[string]interface{} `json:"updated_session_data"`
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
	Message string `json:"message"`
}

// ===============================
// 用户认证授权相关Input/Output结构体
// ===============================

// HashPasswordInput 密码加密输入
type HashPasswordInput struct {
	Password string `json:"password"`
}

// HashPasswordOutput 密码加密输出
type HashPasswordOutput struct {
	HashedPassword string `json:"hashed_password"`
}

// VerifyPasswordInput 验证密码输入
type VerifyPasswordInput struct {
	Password       string `json:"password"`
	HashedPassword string `json:"hashed_password"`
}

// VerifyPasswordOutput 验证密码输出
type VerifyPasswordOutput struct {
	IsValid bool `json:"is_valid"`
}

// GenerateSaltInput 生成盐值输入
type GenerateSaltInput struct{}

// GenerateSaltOutput 生成盐值输出
type GenerateSaltOutput struct {
	Salt string `json:"salt"`
}

// HandleFailedLoginInput 处理登录失败输入
type HandleFailedLoginInput struct {
	UserData map[string]interface{} `json:"user_data"`
}

// HandleFailedLoginOutput 处理登录失败输出
type HandleFailedLoginOutput struct {
	UpdatedUserData map[string]interface{} `json:"updated_user_data"`
	IsLocked        bool                   `json:"is_locked"`
	LockDuration    string                 `json:"lock_duration,omitempty"`
}

// CalculateUserRiskScoreInput 计算用户风险评分输入
type CalculateUserRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateUserRiskScoreOutput 计算用户风险评分输出
type CalculateUserRiskScoreOutput struct {
	RiskScore       float64            `json:"risk_score"`
	RiskFactors     map[string]float64 `json:"risk_factors"`
	RiskLevel       string             `json:"risk_level"`
	Recommendations []string           `json:"recommendations"`
}

// ===============================
// 用户验证相关Input/Output结构体
// ===============================

// ValidateUserRegistrationInput 验证用户注册输入
type ValidateUserRegistrationInput struct {
	UserData map[string]interface{} `json:"user_data"`
}

// ValidateUserRegistrationOutput 验证用户注册输出
type ValidateUserRegistrationOutput struct {
	IsValid bool              `json:"is_valid"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// ValidateUsernameInput 验证用户名输入
type ValidateUsernameInput struct {
	Username string `json:"username"`
}

// ValidateUsernameOutput 验证用户名输出
type ValidateUsernameOutput struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
}

// ValidatePasswordStrengthInput 验证密码强度输入
type ValidatePasswordStrengthInput struct {
	Password string `json:"password"`
}

// ValidatePasswordStrengthOutput 验证密码强度输出
type ValidatePasswordStrengthOutput struct {
	IsValid     bool     `json:"is_valid"`
	Strength    string   `json:"strength"`
	Score       float64  `json:"score"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// ValidateEmailInput 验证邮箱输入
type ValidateEmailInput struct {
	Email string `json:"email"`
}

// ValidateEmailOutput 验证邮箱输出
type ValidateEmailOutput struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
}

// ValidatePhoneInput 验证手机号输入
type ValidatePhoneInput struct {
	Phone string `json:"phone"`
}

// ValidatePhoneOutput 验证手机号输出
type ValidatePhoneOutput struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
}

// ===============================
// 用户会话管理相关Input/Output结构体
// ===============================

// GenerateSessionIdInput 生成会话ID输入
type GenerateSessionIdInput struct{}

// GenerateSessionIdOutput 生成会话ID输出
type GenerateSessionIdOutput struct {
	SessionId string `json:"session_id"`
}

// ===============================
// 用户行为分析相关Input/Output结构体
// ===============================

// AnalyzeUserBehaviorInput 分析用户行为输入
type AnalyzeUserBehaviorInput struct {
	ActivityData []map[string]interface{} `json:"activity_data"`
}

// AnalyzeUserBehaviorOutput 分析用户行为输出
type AnalyzeUserBehaviorOutput struct {
	Analysis map[string]interface{} `json:"analysis"`
}

// FindMostActiveHourInput 找到最活跃时间段输入
type FindMostActiveHourInput struct {
	HourCounts map[int]int `json:"hour_counts"`
}

// FindMostActiveHourOutput 找到最活跃时间段输出
type FindMostActiveHourOutput struct {
	MostActiveHour int `json:"most_active_hour"`
}

// IdentifyBehaviorPatternInput 识别行为模式输入
type IdentifyBehaviorPatternInput struct {
	ActionCounts map[string]int `json:"action_counts"`
	HourCounts   map[int]int    `json:"hour_counts"`
}

// IdentifyBehaviorPatternOutput 识别行为模式输出
type IdentifyBehaviorPatternOutput struct {
	Pattern string `json:"pattern"`
}

// CalculateActivityScoreInput 计算活跃度评分输入
type CalculateActivityScoreInput struct {
	ActivityData []map[string]interface{} `json:"activity_data"`
}

// CalculateActivityScoreOutput 计算活跃度评分输出
type CalculateActivityScoreOutput struct {
	Score float64 `json:"score"`
}

// ===============================
// 风险评分相关Input/Output结构体
// ===============================

// CalculateLoginRiskScoreInput 计算登录风险评分输入
type CalculateLoginRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateLoginRiskScoreOutput 计算登录风险评分输出
type CalculateLoginRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateLocationRiskScoreInput 计算地理位置风险评分输入
type CalculateLocationRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateLocationRiskScoreOutput 计算地理位置风险评分输出
type CalculateLocationRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateDeviceRiskScoreInput 计算设备风险评分输入
type CalculateDeviceRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateDeviceRiskScoreOutput 计算设备风险评分输出
type CalculateDeviceRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateTimeRiskScoreInput 计算时间异常风险评分输入
type CalculateTimeRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateTimeRiskScoreOutput 计算时间异常风险评分输出
type CalculateTimeRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateBehaviorRiskScoreInput 计算行为异常风险评分输入
type CalculateBehaviorRiskScoreInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// CalculateBehaviorRiskScoreOutput 计算行为异常风险评分输出
type CalculateBehaviorRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateHistoryRiskScoreInput 计算历史安全事件风险评分输入
type CalculateHistoryRiskScoreInput struct {
	UserData map[string]interface{} `json:"user_data"`
}

// CalculateHistoryRiskScoreOutput 计算历史安全事件风险评分输出
type CalculateHistoryRiskScoreOutput struct {
	Score float64 `json:"score"`
}

// ===============================
// 辅助功能相关Input/Output结构体
// ===============================

// IsSameNetworkInput 检查网络输入
type IsSameNetworkInput struct {
	IP1 string `json:"ip1"`
	IP2 string `json:"ip2"`
}

// IsSameNetworkOutput 检查网络输出
type IsSameNetworkOutput struct {
	IsSame bool `json:"is_same"`
}

// IsKnownMaliciousIPInput 检查恶意IP输入
type IsKnownMaliciousIPInput struct {
	IP string `json:"ip"`
}

// IsKnownMaliciousIPOutput 检查恶意IP输出
type IsKnownMaliciousIPOutput struct {
	IsMalicious bool `json:"is_malicious"`
}

// IsMobileDeviceInput 检查移动设备输入
type IsMobileDeviceInput struct {
	UserAgent string `json:"user_agent"`
}

// IsMobileDeviceOutput 检查移动设备输出
type IsMobileDeviceOutput struct {
	IsMobile bool `json:"is_mobile"`
}

// IsAbnormalLoginPatternInput 检查异常登录模式输入
type IsAbnormalLoginPatternInput struct {
	UserData  map[string]interface{} `json:"user_data"`
	LoginData map[string]interface{} `json:"login_data"`
}

// IsAbnormalLoginPatternOutput 检查异常登录模式输出
type IsAbnormalLoginPatternOutput struct {
	IsAbnormal bool `json:"is_abnormal"`
}
