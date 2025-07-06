package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 队列操作控制 API (6个)
// ===============================

// StartQueueReq 启动队列请求
type StartQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/start" method:"post" tags:"队列控制" summary:"启动队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Force   bool   `json:"force,omitempty"`
}

type StartQueueRes struct {
	Message string `json:"message"`
}

// PauseQueueReq 暂停队列请求
type PauseQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/pause" method:"post" tags:"队列控制" summary:"暂停队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#暂停原因最大200字符"`
}

type PauseQueueRes struct {
	Message string `json:"message"`
}

// ResumeQueueReq 恢复队列请求
type ResumeQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/resume" method:"post" tags:"队列控制" summary:"恢复队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#恢复原因最大200字符"`
}

type ResumeQueueRes struct {
	Message string `json:"message"`
}

// StopQueueReq 停止队列请求
type StopQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/stop" method:"post" tags:"队列控制" summary:"停止队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#停止原因最大200字符"`
	Force   bool   `json:"force,omitempty"`
}

type StopQueueRes struct {
	Message string `json:"message"`
}

// ClearQueueReq 清空队列请求
type ClearQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/clear" method:"post" tags:"队列控制" summary:"清空队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#清空原因最大200字符"`
	Force   bool   `json:"force,omitempty"`
}

type ClearQueueRes struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ResetQueueReq 重置队列请求
type ResetQueueReq struct {
	g.Meta  `path:"/queue/{queueId}/reset" method:"post" tags:"队列控制" summary:"重置队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#重置原因最大200字符"`
}

type ResetQueueRes struct {
	Message string `json:"message"`
}
