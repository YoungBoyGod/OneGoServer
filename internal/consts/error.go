package consts

// 业务错误码
const (
	// 用户相关错误
	ErrUserNotFound       = "USER_NOT_FOUND"
	ErrUserAlreadyExists  = "USER_ALREADY_EXISTS"
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrUsernameTaken      = "USERNAME_TAKEN"
	ErrEmailTaken         = "EMAIL_TAKEN"
	ErrInvalidPassword    = "INVALID_PASSWORD"
	ErrUserDisabled       = "USER_DISABLED"
	ErrAccountLocked      = "ACCOUNT_LOCKED"

	// 认证相关错误
	ErrInvalidToken     = "INVALID_TOKEN"
	ErrTokenExpired     = "TOKEN_EXPIRED"
	ErrMissingToken     = "MISSING_TOKEN"
	ErrPermissionDenied = "PERMISSION_DENIED"

	// 参数验证错误
	ErrInvalidParams = "INVALID_PARAMS"
	ErrMissingParams = "MISSING_PARAMS"
	ErrInvalidFormat = "INVALID_FORMAT"
	ErrOutOfRange    = "OUT_OF_RANGE"

	// 系统错误
	ErrInternalServer = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError  = "DATABASE_ERROR"
	ErrNetworkError   = "NETWORK_ERROR"
	ErrServiceTimeout = "SERVICE_TIMEOUT"
)

// 错误消息映射
var ErrorMessages = map[string]string{
	ErrUserNotFound:       "用户不存在",
	ErrUserAlreadyExists:  "用户已存在",
	ErrInvalidCredentials: "用户名或密码错误",
	ErrUsernameTaken:      "用户名已被占用",
	ErrEmailTaken:         "邮箱已被占用",
	ErrInvalidPassword:    "密码格式不正确",
	ErrUserDisabled:       "用户已被禁用",
	ErrAccountLocked:      "账户已被锁定",

	ErrInvalidToken:     "无效的令牌",
	ErrTokenExpired:     "令牌已过期",
	ErrMissingToken:     "缺少令牌",
	ErrPermissionDenied: "权限不足",

	ErrInvalidParams: "参数不正确",
	ErrMissingParams: "缺少必要参数",
	ErrInvalidFormat: "格式不正确",
	ErrOutOfRange:    "参数超出范围",

	ErrInternalServer: "服务器内部错误",
	ErrDatabaseError:  "数据库错误",
	ErrNetworkError:   "网络错误",
	ErrServiceTimeout: "服务超时",
}
