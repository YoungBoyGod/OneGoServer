package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统告警管理相关API
// ===============================

// GetSystemAlertsReq 获取系统告警请求
type GetSystemAlertsReq struct {
	g.Meta    `path:"/alerts" method:"get" tags:"系统管理" summary:"获取系统告警"`
	Level     string `json:"level,omitempty" v:"in:info,warn,error,critical#告警级别只能是info,warn,error,critical"`
	Status    string `json:"status,omitempty" v:"in:active,resolved,acknowledged#告警状态只能是active,resolved,acknowledged"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
}

// GetSystemAlertsRes 获取系统告警响应
type GetSystemAlertsRes struct {
	List  []SystemAlert `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// AcknowledgeAlertReq 确认告警请求
type AcknowledgeAlertReq struct {
	g.Meta  `path:"/alerts/acknowledge" method:"post" tags:"系统管理" summary:"确认告警"`
	AlertId string `json:"alertId" v:"required#告警ID不能为空"`
	Comment string `json:"comment,omitempty"`
}

// AcknowledgeAlertRes 确认告警响应
type AcknowledgeAlertRes struct {
	AlertId string `json:"alertId"`
	Status  string `json:"status"`
}

// ResolveAlertReq 解决告警请求
type ResolveAlertReq struct {
	g.Meta  `path:"/alerts/resolve" method:"post" tags:"系统管理" summary:"解决告警"`
	AlertId string `json:"alertId" v:"required#告警ID不能为空"`
	Comment string `json:"comment,omitempty"`
}

// ResolveAlertRes 解决告警响应
type ResolveAlertRes struct {
	AlertId string `json:"alertId"`
	Status  string `json:"status"`
}
