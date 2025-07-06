package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 用户基础管理相关API
// ===============================

// RegisterUserReq 用户注册请求
type RegisterUserReq struct {
	g.Meta     `path:"/user/register" method:"post" tags:"用户管理" summary:"用户注册"`
	Username   string `json:"username" v:"required|length:3,20|regex:^[a-zA-Z0-9_]+$#用户名不能为空|用户名长度为3-20字符|用户名只能包含字母数字下划线"`
	Password   string `json:"password" v:"required|length:6,32#密码不能为空|密码长度为6-32字符"`
	Email      string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确"`
	Phone      string `json:"phone,omitempty" v:"phone#手机号格式不正确"`
	RealName   string `json:"real_name,omitempty" v:"length:1,50#真实姓名长度为1-50字符"`
	Avatar     string `json:"avatar,omitempty" v:"url#头像URL格式不正确"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Position   string `json:"position,omitempty" v:"length:1,50#职位名称长度为1-50字符"`
}

type RegisterUserRes struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

// LoginUserReq 用户登录请求
type LoginUserReq struct {
	g.Meta     `path:"/user/login" method:"post" tags:"用户管理" summary:"用户登录"`
	Username   string `json:"username" v:"required|length:3,50#用户名不能为空"`
	Password   string `json:"password" v:"required|length:6,32#密码不能为空"`
	LoginType  string `json:"login_type" d:"password" v:"in:password,email,phone#登录方式无效"`
	DeviceInfo string `json:"device_info,omitempty" v:"max-length:200#设备信息最大200字符"`
	IpAddress  string `json:"ip_address,omitempty" v:"ip#IP地址格式不正确"`
}

type LoginUserRes struct {
	UserId       string   `json:"user_id"`
	Username     string   `json:"username"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Permissions  []string `json:"permissions"`
	Roles        []string `json:"roles"`
	UserInfo     UserInfo `json:"user_info"`
}

// LogoutUserReq 用户登出请求
type LogoutUserReq struct {
	g.Meta     `path:"/user/logout" method:"post" tags:"用户管理" summary:"用户登出"`
	UserId     string `json:"user_id,omitempty" v:"max-length:50#用户ID最大50字符"`
	AllDevices bool   `json:"all_devices,omitempty"`
}

type LogoutUserRes struct {
	Message string `json:"message"`
}

// GetUserListReq 获取用户列表请求
type GetUserListReq struct {
	g.Meta                   `path:"/user/list" method:"get" tags:"用户管理" summary:"获取用户列表"`
	common.PaginationRequest `json:",inline"`
	Status                   string      `json:"status,omitempty" v:"in:active,inactive,locked,pending#状态值无效"`
	Role                     string      `json:"role,omitempty" v:"length:1,50#角色名称长度为1-50字符"`
	Department               string      `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Keyword                  string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
}

type GetUserListRes struct {
	common.PaginationResponse[UserInfo] `json:",inline"`
}

// GetUserDetailReq 获取用户详情请求
type GetUserDetailReq struct {
	g.Meta `path:"/user/{userId}" method:"get" tags:"用户管理" summary:"获取用户详情"`
	UserId string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
}

type GetUserDetailRes struct {
	UserDetail UserDetailInfo `json:"user_detail"`
}

// UpdateUserReq 更新用户信息请求
type UpdateUserReq struct {
	g.Meta     `path:"/user/{userId}" method:"put" tags:"用户管理" summary:"更新用户信息"`
	UserId     string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Email      string `json:"email,omitempty" v:"email#邮箱格式不正确"`
	Phone      string `json:"phone,omitempty" v:"phone#手机号格式不正确"`
	RealName   string `json:"real_name,omitempty" v:"length:1,50#真实姓名长度为1-50字符"`
	Avatar     string `json:"avatar,omitempty" v:"url#头像URL格式不正确"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Position   string `json:"position,omitempty" v:"length:1,50#职位名称长度为1-50字符"`
	Status     string `json:"status,omitempty" v:"in:active,inactive,locked#状态值无效"`
}

type UpdateUserRes struct {
	Message string `json:"message"`
}
