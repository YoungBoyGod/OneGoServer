package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 用户会话管理相关API
// ===============================

// GetUserSessionsReq 获取用户会话请求
type GetUserSessionsReq struct {
	g.Meta                   `path:"/user/{userId}/sessions" method:"get" tags:"用户会话" summary:"获取用户会话"`
	UserId                   string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	common.PaginationRequest `json:",inline"`
}

type GetUserSessionsRes struct {
	common.PaginationResponse[SessionInfo] `json:",inline"`
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
	g.Meta                   `path:"/user/{userId}/activity" method:"get" tags:"用户会话" summary:"获取用户活动"`
	UserId                   string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	common.PaginationRequest `json:",inline"`
	ActionType               string      `json:"action_type,omitempty" v:"in:login,logout,create,update,delete,view#操作类型无效"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
}

type GetUserActivityRes struct {
	common.PaginationResponse[ActivityInfo] `json:",inline"`
}
