package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务状态管理 API (2个)
// ===============================

// GetTaskStatusReq 获取任务状态请求
type GetTaskStatusReq struct {
	g.Meta `path:"/task/{taskId}/status" method:"get" tags:"任务状态" summary:"获取任务状态"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
}

type GetTaskStatusRes struct {
	TaskId         string      `json:"task_id"`
	Status         string      `json:"status"`
	Progress       float64     `json:"progress"`
	ExecutionId    string      `json:"execution_id"`
	StartTime      *gtime.Time `json:"start_time"`
	EndTime        *gtime.Time `json:"end_time"`
	Duration       int64       `json:"duration"`
	ErrorMessage   string      `json:"error_message"`
	ExecutionCount int         `json:"execution_count"`
	LastCheckTime  *gtime.Time `json:"last_check_time"`
}

// UpdateTaskStatusReq 更新任务状态请求
type UpdateTaskStatusReq struct {
	g.Meta       `path:"/task/{taskId}/status" method:"put" tags:"任务状态" summary:"更新任务状态"`
	TaskId       string  `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Status       string  `json:"status" v:"required|in:pending,running,completed,failed,cancelled#状态值无效"`
	Progress     float64 `json:"progress,omitempty" v:"between:0,100#进度为0-100"`
	ErrorMessage string  `json:"error_message,omitempty" v:"max-length:1000#错误信息最大1000字符"`
}

type UpdateTaskStatusRes struct {
	Message string `json:"message"`
}
