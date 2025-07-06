package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 用户安全管理相关API
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
