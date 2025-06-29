package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 定义错误码常量
const (
	// 通用错误码 1000-1999
	Success         ErrorCode = 0
	InternalError   ErrorCode = 1000
	InvalidParam    ErrorCode = 1001
	NotFound        ErrorCode = 1002
	Unauthorized    ErrorCode = 1003
	Forbidden       ErrorCode = 1004
	TooManyRequests ErrorCode = 1005

	// 用户相关错误码 2000-2999
	UserNotFound          ErrorCode = 2001
	UserAlreadyExists     ErrorCode = 2002
	InvalidCredentials    ErrorCode = 2003
	UserDisabled          ErrorCode = 2004
	WeakPassword          ErrorCode = 2005
	InvalidEmail          ErrorCode = 2006
	EmailAlreadyExists    ErrorCode = 2007
	UsernameAlreadyExists ErrorCode = 2008
	InvalidToken          ErrorCode = 2009
	TokenExpired          ErrorCode = 2010

	// 设备相关错误码 3000-3999
	DeviceNotFound      ErrorCode = 3001
	DeviceAlreadyExists ErrorCode = 3002
	DeviceOffline       ErrorCode = 3003
	DeviceUnauthorized  ErrorCode = 3004
	DeviceBusy          ErrorCode = 3005

	// 任务相关错误码 4000-4999
	TaskNotFound      ErrorCode = 4001
	TaskAlreadyExists ErrorCode = 4002
	TaskInProgress    ErrorCode = 4003
	TaskFailed        ErrorCode = 4004
	TaskCanceled      ErrorCode = 4005
	InvalidTaskStatus ErrorCode = 4006

	// 告警相关错误码 5000-5999
	AlertNotFound        ErrorCode = 5001
	AlertAlreadyExists   ErrorCode = 5002
	InvalidAlertLevel    ErrorCode = 5003
	AlertAlreadyResolved ErrorCode = 5004
)

// BusinessError 业务错误结构
type BusinessError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// Error 实现error接口
func (e *BusinessError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("业务错误[%d]: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("业务错误[%d]: %s", e.Code, e.Message)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code ErrorCode, message string, details ...string) *BusinessError {
	err := &BusinessError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// 错误码到HTTP状态码的映射
var errorCodeToHTTPStatus = map[ErrorCode]int{
	Success:         http.StatusOK,
	InternalError:   http.StatusInternalServerError,
	InvalidParam:    http.StatusBadRequest,
	NotFound:        http.StatusNotFound,
	Unauthorized:    http.StatusUnauthorized,
	Forbidden:       http.StatusForbidden,
	TooManyRequests: http.StatusTooManyRequests,

	UserNotFound:          http.StatusNotFound,
	UserAlreadyExists:     http.StatusConflict,
	InvalidCredentials:    http.StatusUnauthorized,
	UserDisabled:          http.StatusForbidden,
	WeakPassword:          http.StatusBadRequest,
	InvalidEmail:          http.StatusBadRequest,
	EmailAlreadyExists:    http.StatusConflict,
	UsernameAlreadyExists: http.StatusConflict,
	InvalidToken:          http.StatusUnauthorized,
	TokenExpired:          http.StatusUnauthorized,

	DeviceNotFound:      http.StatusNotFound,
	DeviceAlreadyExists: http.StatusConflict,
	DeviceOffline:       http.StatusServiceUnavailable,
	DeviceUnauthorized:  http.StatusUnauthorized,
	DeviceBusy:          http.StatusConflict,

	TaskNotFound:      http.StatusNotFound,
	TaskAlreadyExists: http.StatusConflict,
	TaskInProgress:    http.StatusConflict,
	TaskFailed:        http.StatusInternalServerError,
	TaskCanceled:      http.StatusGone,
	InvalidTaskStatus: http.StatusBadRequest,

	AlertNotFound:        http.StatusNotFound,
	AlertAlreadyExists:   http.StatusConflict,
	InvalidAlertLevel:    http.StatusBadRequest,
	AlertAlreadyResolved: http.StatusConflict,
}

// 错误码到错误消息的映射
var errorCodeToMessage = map[ErrorCode]string{
	Success:         "操作成功",
	InternalError:   "内部服务器错误",
	InvalidParam:    "参数错误",
	NotFound:        "资源不存在",
	Unauthorized:    "未授权",
	Forbidden:       "禁止访问",
	TooManyRequests: "请求过于频繁",

	UserNotFound:          "用户不存在",
	UserAlreadyExists:     "用户已存在",
	InvalidCredentials:    "用户名或密码错误",
	UserDisabled:          "用户已被禁用",
	WeakPassword:          "密码强度不够",
	InvalidEmail:          "邮箱格式错误",
	EmailAlreadyExists:    "邮箱已存在",
	UsernameAlreadyExists: "用户名已存在",
	InvalidToken:          "无效的token",
	TokenExpired:          "token已过期",

	DeviceNotFound:      "设备不存在",
	DeviceAlreadyExists: "设备已存在",
	DeviceOffline:       "设备离线",
	DeviceUnauthorized:  "设备未授权",
	DeviceBusy:          "设备忙碌",

	TaskNotFound:      "任务不存在",
	TaskAlreadyExists: "任务已存在",
	TaskInProgress:    "任务正在执行中",
	TaskFailed:        "任务执行失败",
	TaskCanceled:      "任务已取消",
	InvalidTaskStatus: "无效的任务状态",

	AlertNotFound:        "告警不存在",
	AlertAlreadyExists:   "告警已存在",
	InvalidAlertLevel:    "无效的告警级别",
	AlertAlreadyResolved: "告警已解决",
}

// GetHTTPStatus 获取错误码对应的HTTP状态码
func (e *BusinessError) GetHTTPStatus() int {
	if status, exists := errorCodeToHTTPStatus[e.Code]; exists {
		return status
	}
	return http.StatusInternalServerError
}

// GetMessage 获取错误码对应的默认消息
func GetMessage(code ErrorCode) string {
	if message, exists := errorCodeToMessage[code]; exists {
		return message
	}
	return "未知错误"
}

// 快捷创建常用错误的函数

// NewInternalError 创建内部错误
func NewInternalError(details string) *BusinessError {
	return NewBusinessError(InternalError, GetMessage(InternalError), details)
}

// NewInvalidParamError 创建参数错误
func NewInvalidParamError(details string) *BusinessError {
	return NewBusinessError(InvalidParam, GetMessage(InvalidParam), details)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(resource string) *BusinessError {
	message := fmt.Sprintf("%s不存在", resource)
	return NewBusinessError(NotFound, message)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(details string) *BusinessError {
	return NewBusinessError(Unauthorized, GetMessage(Unauthorized), details)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(details string) *BusinessError {
	return NewBusinessError(Forbidden, GetMessage(Forbidden), details)
}

// 用户相关错误快捷函数

// NewUserNotFoundError 创建用户不存在错误
func NewUserNotFoundError() *BusinessError {
	return NewBusinessError(UserNotFound, GetMessage(UserNotFound))
}

// NewUserAlreadyExistsError 创建用户已存在错误
func NewUserAlreadyExistsError(field string) *BusinessError {
	message := fmt.Sprintf("%s已存在", field)
	return NewBusinessError(UserAlreadyExists, message)
}

// NewInvalidCredentialsError 创建凭据错误
func NewInvalidCredentialsError() *BusinessError {
	return NewBusinessError(InvalidCredentials, GetMessage(InvalidCredentials))
}

// NewInvalidTokenError 创建无效token错误
func NewInvalidTokenError() *BusinessError {
	return NewBusinessError(InvalidToken, GetMessage(InvalidToken))
}

// NewTokenExpiredError 创建token过期错误
func NewTokenExpiredError() *BusinessError {
	return NewBusinessError(TokenExpired, GetMessage(TokenExpired))
}

// 设备相关错误快捷函数

// NewDeviceNotFoundError 创建设备不存在错误
func NewDeviceNotFoundError() *BusinessError {
	return NewBusinessError(DeviceNotFound, GetMessage(DeviceNotFound))
}

// NewDeviceOfflineError 创建设备离线错误
func NewDeviceOfflineError() *BusinessError {
	return NewBusinessError(DeviceOffline, GetMessage(DeviceOffline))
}

// 任务相关错误快捷函数

// NewTaskNotFoundError 创建任务不存在错误
func NewTaskNotFoundError() *BusinessError {
	return NewBusinessError(TaskNotFound, GetMessage(TaskNotFound))
}

// NewTaskInProgressError 创建任务进行中错误
func NewTaskInProgressError() *BusinessError {
	return NewBusinessError(TaskInProgress, GetMessage(TaskInProgress))
}

// 告警相关错误快捷函数

// NewAlertNotFoundError 创建告警不存在错误
func NewAlertNotFoundError() *BusinessError {
	return NewBusinessError(AlertNotFound, GetMessage(AlertNotFound))
}

// NewAlertAlreadyResolvedError 创建告警已解决错误
func NewAlertAlreadyResolvedError() *BusinessError {
	return NewBusinessError(AlertAlreadyResolved, GetMessage(AlertAlreadyResolved))
}

// IsBusinessError 检查是否为业务错误
func IsBusinessError(err error) (*BusinessError, bool) {
	if bizErr, ok := err.(*BusinessError); ok {
		return bizErr, true
	}
	return nil, false
}
