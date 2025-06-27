package response

import (
	"net/http"

	v1 "github.com/YoungBoyGod/OneGoServer/api/v1"
	"github.com/gin-gonic/gin"
)

// 响应码定义
const (
	CodeSuccess      = 200 // 成功
	CodeInvalidParam = 400 // 参数错误
	CodeUnauthorized = 401 // 未授权
	CodeForbidden    = 403 // 禁止访问
	CodeNotFound     = 404 // 未找到
	CodeServerError  = 500 // 服务器错误
)

// 响应消息定义
const (
	MsgSuccess      = "success"
	MsgInvalidParam = "参数错误"
	MsgUnauthorized = "未授权"
	MsgForbidden    = "禁止访问"
	MsgNotFound     = "未找到"
	MsgServerError  = "服务器内部错误"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, v1.CommonRes{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, v1.CommonRes{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	var httpStatus int
	switch code {
	case CodeInvalidParam:
		httpStatus = http.StatusBadRequest
	case CodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case CodeForbidden:
		httpStatus = http.StatusForbidden
	case CodeNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
		code = CodeServerError
		if message == "" {
			message = MsgServerError
		}
	}

	c.JSON(httpStatus, v1.CommonRes{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// InvalidParam 参数错误响应
func InvalidParam(c *gin.Context, message string) {
	if message == "" {
		message = MsgInvalidParam
	}
	Error(c, CodeInvalidParam, message)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	Error(c, CodeUnauthorized, message)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MsgForbidden
	}
	Error(c, CodeForbidden, message)
}

// NotFound 未找到响应
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MsgNotFound
	}
	Error(c, CodeNotFound, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = MsgServerError
	}
	Error(c, CodeServerError, message)
}

// Page 分页响应
func Page(c *gin.Context, data interface{}, page, pageSize, total int) {
	c.JSON(http.StatusOK, v1.PageRes{
		CommonRes: v1.CommonRes{
			Code:    CodeSuccess,
			Message: MsgSuccess,
			Data:    data,
		},
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}
