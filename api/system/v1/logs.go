package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统日志管理相关API
// ===============================

// GetSystemLogsReq 获取系统日志请求
type GetSystemLogsReq struct {
	g.Meta    `path:"/logs" method:"get" tags:"系统管理" summary:"获取系统日志"`
	Level     string `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别只能是debug,info,warn,error"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"100" v:"between:1,1000#每页数量为1-1000"`
}

// GetSystemLogsRes 获取系统日志响应
type GetSystemLogsRes struct {
	List  []SystemLog `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
