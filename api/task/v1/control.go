package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 任务执行控制 API (5个)
// ===============================

// StartTaskReq 启动任务请求
type StartTaskReq struct {
	g.Meta `path:"/task/{taskId}/start" method:"post" tags:"任务执行" summary:"启动任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Force  bool   `json:"force,omitempty"`
}

type StartTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}

// StopTaskReq 停止任务请求
type StopTaskReq struct {
	g.Meta `path:"/task/{taskId}/stop" method:"post" tags:"任务执行" summary:"停止任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#停止原因最大200字符"`
}

type StopTaskRes struct {
	Message string `json:"message"`
}

// RestartTaskReq 重启任务请求
type RestartTaskReq struct {
	g.Meta `path:"/task/{taskId}/restart" method:"post" tags:"任务执行" summary:"重启任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#重启原因最大200字符"`
}

type RestartTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}

// CancelTaskReq 取消任务请求
type CancelTaskReq struct {
	g.Meta `path:"/task/{taskId}/cancel" method:"post" tags:"任务执行" summary:"取消任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#取消原因最大200字符"`
}

type CancelTaskRes struct {
	Message string `json:"message"`
}

// RetryTaskReq 重试任务请求
type RetryTaskReq struct {
	g.Meta `path:"/task/{taskId}/retry" method:"post" tags:"任务执行" summary:"重试任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#重试原因最大200字符"`
}

type RetryTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}
