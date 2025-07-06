package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统日志管理相关API
// ===============================

// GetSystemLogsReq 获取系统日志请求
type GetSystemLogsReq struct {
	g.Meta                   `path:"/logs" method:"get" tags:"系统管理" summary:"获取系统日志"`
	common.PaginationRequest `json:",inline"`
	Level                    string `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别只能是debug,info,warn,error"`
	StartTime                string `json:"startTime,omitempty"`
	EndTime                  string `json:"endTime,omitempty"`
}

// GetSystemLogsRes 获取系统日志响应
type GetSystemLogsRes struct {
	common.PaginationResponse[SystemLog] `json:",inline"`
}
