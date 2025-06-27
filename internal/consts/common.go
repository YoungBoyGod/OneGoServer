package consts

// HTTP状态码
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500
)

// 分页默认值
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// JWT相关常量
const (
	JWTSigningKey     = "OneGo_JWT_Secret_Key_2024"
	JWTExpireDuration = 24 * 60 * 60 // 24小时（秒）
	JWTIssuer         = "OneGo"
)

// 用户相关常量
const (
	MinUsernameLength = 3
	MaxUsernameLength = 20
	MinPasswordLength = 6
	MaxPasswordLength = 50
	MaxLoginAttempts  = 5
	LockoutDuration   = 30 * 60 // 30分钟（秒）
)

// 数据库相关常量
const (
	DBMaxOpenConns    = 100
	DBMaxIdleConns    = 10
	DBConnMaxLifetime = 300 // 5分钟（秒）
)

// 缓存相关常量
const (
	CacheDefaultExpiration = 300 // 5分钟（秒）
	CacheCleanupInterval   = 600 // 10分钟（秒）
	UserCacheKeyPrefix     = "user:"
	TokenCacheKeyPrefix    = "token:"
)

// 服务器相关常量
const (
	DefaultServerPort  = ":8080"
	MaxRequestBodySize = 32 << 20 // 32MB
	RequestTimeout     = 30       // 30秒
	ReadTimeout        = 30       // 30秒
	WriteTimeout       = 30       // 30秒
)

// 响应状态
const (
	ResponseStatusSuccess = "success"
	ResponseStatusError   = "error"
	ResponseStatusFail    = "fail"
)
