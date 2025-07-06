package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 用户权限管理相关API
// ===============================

// GetUserRolesReq 获取用户角色请求
type GetUserRolesReq struct {
	g.Meta `path:"/user/{userId}/roles" method:"get" tags:"用户权限" summary:"获取用户角色"`
	UserId string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
}

type GetUserRolesRes struct {
	Roles []RoleInfo `json:"roles"`
}

// AssignUserRoleReq 分配用户角色请求
type AssignUserRoleReq struct {
	g.Meta  `path:"/user/{userId}/roles" method:"post" tags:"用户权限" summary:"分配用户角色"`
	UserId  string   `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	RoleIds []string `json:"role_ids" v:"required|max:10#角色ID列表不能为空|最多10个角色"`
	Reason  string   `json:"reason,omitempty" v:"max-length:200#分配原因最大200字符"`
}

type AssignUserRoleRes struct {
	Message string `json:"message"`
}

// RemoveUserRoleReq 移除用户角色请求
type RemoveUserRoleReq struct {
	g.Meta  `path:"/user/{userId}/roles/remove" method:"post" tags:"用户权限" summary:"移除用户角色"`
	UserId  string   `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	RoleIds []string `json:"role_ids" v:"required|max:10#角色ID列表不能为空|最多10个角色"`
	Reason  string   `json:"reason,omitempty" v:"max-length:200#移除原因最大200字符"`
}

type RemoveUserRoleRes struct {
	Message string `json:"message"`
}

// GetUserPermissionsReq 获取用户权限请求
type GetUserPermissionsReq struct {
	g.Meta   `path:"/user/{userId}/permissions" method:"get" tags:"用户权限" summary:"获取用户权限"`
	UserId   string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Resource string `json:"resource,omitempty" v:"length:1,100#资源名称长度为1-100字符"`
}

type GetUserPermissionsRes struct {
	Permissions []PermissionInfo `json:"permissions"`
}

// CheckUserPermissionReq 检查用户权限请求
type CheckUserPermissionReq struct {
	g.Meta     `path:"/user/{userId}/permission/check" method:"post" tags:"用户权限" summary:"检查用户权限"`
	UserId     string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	Resource   string `json:"resource" v:"required|length:1,100#资源名称不能为空"`
	Action     string `json:"action" v:"required|length:1,50#操作名称不能为空"`
	ResourceId string `json:"resource_id,omitempty" v:"max-length:50#资源ID最大50字符"`
}

type CheckUserPermissionRes struct {
	HasPermission bool   `json:"has_permission"`
	Reason        string `json:"reason"`
}
