package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 队列任务管理 API (6个)
// ===============================

// EnqueueTaskReq 任务入队请求
type EnqueueTaskReq struct {
	g.Meta     `path:"/queue/{queueId}/enqueue" method:"post" tags:"队列任务" summary:"任务入队"`
	QueueId    string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskData   map[string]interface{} `json:"task_data" v:"required#任务数据不能为空"`
	Priority   int                    `json:"priority" d:"5" v:"between:1,10#任务优先级为1-10"`
	DelayTime  int64                  `json:"delay_time,omitempty" v:"min:0#延迟时间不能小于0"`
	MaxRetries int                    `json:"max_retries" d:"3" v:"between:0,10#最大重试次数为0-10"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type EnqueueTaskRes struct {
	TaskId   string `json:"task_id"`
	Position int    `json:"position"`
	Message  string `json:"message"`
}

// DequeueTaskReq 任务出队请求
type DequeueTaskReq struct {
	g.Meta  `path:"/queue/{queueId}/dequeue" method:"post" tags:"队列任务" summary:"任务出队"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Count   int    `json:"count" d:"1" v:"between:1,10#出队数量为1-10"`
	Timeout int    `json:"timeout" d:"30" v:"between:1,300#超时时间为1-300秒"`
}

type DequeueTaskRes struct {
	Tasks   []QueueTaskInfo `json:"tasks"`
	Message string          `json:"message"`
}

// GetQueueTasksReq 获取队列任务请求
type GetQueueTasksReq struct {
	g.Meta                   `path:"/queue/{queueId}/tasks" method:"get" tags:"队列任务" summary:"获取队列任务"`
	QueueId                  string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	common.PaginationRequest `json:",inline"`
	Status                   string `json:"status,omitempty" v:"in:pending,processing,completed,failed,retry#状态值无效"`
	Priority                 int    `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
}

type GetQueueTasksRes struct {
	common.PaginationResponse[QueueTaskInfo] `json:",inline"`
}

// ReorderQueueTasksReq 重新排序队列任务请求
type ReorderQueueTasksReq struct {
	g.Meta    `path:"/queue/{queueId}/reorder" method:"post" tags:"队列任务" summary:"重新排序队列任务"`
	QueueId   string   `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskIds   []string `json:"task_ids" v:"required|max:100#任务ID列表不能为空|最多100个任务"`
	Strategy  string   `json:"strategy" v:"required|in:manual,priority,fifo,lifo#排序策略无效"`
	Positions []int    `json:"positions,omitempty"`
}

type ReorderQueueTasksRes struct {
	AffectedCount int    `json:"affected_count"`
	Message       string `json:"message"`
}

// RemoveQueueTaskReq 移除队列任务请求
type RemoveQueueTaskReq struct {
	g.Meta  `path:"/queue/{queueId}/task/{taskId}/remove" method:"post" tags:"队列任务" summary:"移除队列任务"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskId  string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason  string `json:"reason,omitempty" v:"max-length:200#移除原因最大200字符"`
}

type RemoveQueueTaskRes struct {
	Message string `json:"message"`
}

// UpdateTaskPriorityReq 更新任务优先级请求
type UpdateTaskPriorityReq struct {
	g.Meta   `path:"/queue/{queueId}/task/{taskId}/priority" method:"put" tags:"队列任务" summary:"更新任务优先级"`
	QueueId  string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	TaskId   string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Priority int    `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type UpdateTaskPriorityRes struct {
	NewPosition int    `json:"new_position"`
	Message     string `json:"message"`
}
