package response

import (
	"net/http"

	"github.com/YoungBoyGod/OneGoServer/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 错误码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
	Meta    *Meta       `json:"meta,omitempty"` // 元数据（分页信息等）
}

// Meta 元数据结构
type Meta struct {
	Page     int   `json:"page,omitempty"`      // 当前页码
	PageSize int   `json:"page_size,omitempty"` // 每页大小
	Total    int64 `json:"total,omitempty"`     // 总记录数
	Pages    int   `json:"pages,omitempty"`     // 总页数
}

// PageData 分页数据结构
type PageData struct {
	List interface{} `json:"list"` // 数据列表
	Meta *Meta       `json:"meta"` // 分页信息
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.Success),
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    int(errors.Success),
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, page, pageSize int, total int64) {
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages == 0 {
		pages = 1
	}

	meta := &Meta{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
	}

	c.JSON(http.StatusOK, Response{
		Code:    int(errors.Success),
		Message: "获取成功",
		Data: PageData{
			List: list,
			Meta: meta,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	// 检查是否为业务错误
	if bizErr, ok := errors.IsBusinessError(err); ok {
		c.JSON(bizErr.GetHTTPStatus(), Response{
			Code:    int(bizErr.Code),
			Message: bizErr.Message,
		})
		return
	}

	// 默认内部服务器错误
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errors.InternalError),
		Message: "内部服务器错误",
	})
}

// ErrorWithCode 错误响应（指定错误码）
func ErrorWithCode(c *gin.Context, code errors.ErrorCode, message string) {
	httpStatus := http.StatusInternalServerError
	if bizErr := errors.NewBusinessError(code, message); bizErr != nil {
		httpStatus = bizErr.GetHTTPStatus()
	}

	c.JSON(httpStatus, Response{
		Code:    int(code),
		Message: message,
	})
}

// BadRequest 400错误响应
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    int(errors.InvalidParam),
		Message: message,
	})
}

// Unauthorized 401错误响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权访问"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    int(errors.Unauthorized),
		Message: message,
	})
}

// Forbidden 403错误响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    int(errors.Forbidden),
		Message: message,
	})
}

// NotFound 404错误响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    int(errors.NotFound),
		Message: message,
	})
}

// InternalError 500错误响应
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "内部服务器错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    int(errors.InternalError),
		Message: message,
	})
}

// TooManyRequests 429错误响应
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = "请求过于频繁"
	}
	c.JSON(http.StatusTooManyRequests, Response{
		Code:    int(errors.TooManyRequests),
		Message: message,
	})
}

// 用户相关错误响应

// UserNotFound 用户不存在
func UserNotFound(c *gin.Context) {
	Error(c, errors.NewUserNotFoundError())
}

// UserAlreadyExists 用户已存在
func UserAlreadyExists(c *gin.Context, field string) {
	Error(c, errors.NewUserAlreadyExistsError(field))
}

// InvalidCredentials 无效凭据
func InvalidCredentials(c *gin.Context) {
	Error(c, errors.NewInvalidCredentialsError())
}

// InvalidToken 无效token
func InvalidToken(c *gin.Context) {
	Error(c, errors.NewInvalidTokenError())
}

// TokenExpired token过期
func TokenExpired(c *gin.Context) {
	Error(c, errors.NewTokenExpiredError())
}

// 设备相关错误响应

// DeviceNotFound 设备不存在
func DeviceNotFound(c *gin.Context) {
	Error(c, errors.NewDeviceNotFoundError())
}

// DeviceOffline 设备离线
func DeviceOffline(c *gin.Context) {
	Error(c, errors.NewDeviceOfflineError())
}

// 任务相关错误响应

// TaskNotFound 任务不存在
func TaskNotFound(c *gin.Context) {
	Error(c, errors.NewTaskNotFoundError())
}

// TaskInProgress 任务进行中
func TaskInProgress(c *gin.Context) {
	Error(c, errors.NewTaskInProgressError())
}

// 告警相关错误响应

// AlertNotFound 告警不存在
func AlertNotFound(c *gin.Context) {
	Error(c, errors.NewAlertNotFoundError())
}

// AlertAlreadyResolved 告警已解决
func AlertAlreadyResolved(c *gin.Context) {
	Error(c, errors.NewAlertAlreadyResolvedError())
}

// AbortWithError 中断请求并返回错误
func AbortWithError(c *gin.Context, err error) {
	Error(c, err)
	c.Abort()
}

// AbortWithBusinessError 中断请求并返回业务错误
func AbortWithBusinessError(c *gin.Context, code errors.ErrorCode, message string) {
	ErrorWithCode(c, code, message)
	c.Abort()
}

// BindingError 参数绑定错误
func BindingError(c *gin.Context, err error) {
	BadRequest(c, "参数格式错误: "+err.Error())
}

// ValidationError 参数验证错误
func ValidationError(c *gin.Context, message string) {
	BadRequest(c, "参数验证失败: "+message)
}

// GetPage 从查询参数获取分页信息
func GetPage(c *gin.Context) (page, pageSize int) {
	page = 1
	pageSize = 20

	if p := c.Query("page"); p != "" {
		if pageVal := parseInt(p); pageVal > 0 {
			page = pageVal
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if pageSizeVal := parseInt(ps); pageSizeVal > 0 && pageSizeVal <= 100 {
			pageSize = pageSizeVal
		}
	}

	return page, pageSize
}

// parseInt 字符串转整数
func parseInt(s string) int {
	result := 0
	for _, digit := range s {
		if digit >= '0' && digit <= '9' {
			result = result*10 + int(digit-'0')
		} else {
			return 0
		}
	}
	return result
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUserRoleFromContext 从上下文获取用户角色
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	if role, exists := c.Get("user_role"); exists {
		if roleStr, ok := role.(string); ok {
			return roleStr, true
		}
	}
	return "", false
}

// SetUserToContext 设置用户信息到上下文
func SetUserToContext(c *gin.Context, userID uint, username, role string) {
	c.Set("user_id", userID)
	c.Set("username", username)
	c.Set("user_role", role)
}

// HealthCheck 健康检查响应
func HealthCheck(c *gin.Context, status string, details map[string]interface{}) {
	response := gin.H{
		"status": status,
		"time":   gin.H{"startup": "2024-12-29T10:00:00Z"},
	}

	if details != nil {
		response["details"] = details
	}

	if status == "OK" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
