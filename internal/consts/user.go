package consts

// ================================
// User 模块常量定义
// ================================

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

// 用户操作类型常量
const (
	UserActionTypeLogin  = "login"
	UserActionTypeLogout = "logout"
	UserActionTypeCreate = "create"
	UserActionTypeUpdate = "update"
	UserActionTypeDelete = "delete"
	UserActionTypeView   = "view"
)

// 用户日志类型常量
const (
	UserLogTypeLogin            = "login"
	UserLogTypePasswordChange   = "password_change"
	UserLogTypePermissionChange = "permission_change"
	UserLogTypeSecurityEvent    = "security_event"
)

// 用户日志级别常量
const (
	UserLogLevelInfo     = "info"
	UserLogLevelWarn     = "warn"
	UserLogLevelError    = "error"
	UserLogLevelCritical = "critical"
)
