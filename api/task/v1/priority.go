package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 任务优先级管理 API (3个)
// ===============================

// GetTaskPriorityListReq 获取任务优先级列表请求
type GetTaskPriorityListReq struct {
	g.Meta `path:"/task/priority/list" method:"get" tags:"任务优先级" summary:"获取任务优先级列表"`
}

type GetTaskPriorityListRes struct {
	PriorityList []TaskPriorityInfo `json:"priority_list"`
}

// UpdateTaskPriorityReq 更新任务优先级请求
type UpdateTaskPriorityReq struct {
	g.Meta   `path:"/task/{taskId}/priority" method:"put" tags:"任务优先级" summary:"更新任务优先级"`
	TaskId   string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Priority int    `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type UpdateTaskPriorityRes struct {
	Message string `json:"message"`
}

// BatchUpdateTaskPriorityReq 批量更新任务优先级请求
type BatchUpdateTaskPriorityReq struct {
	g.Meta   `path:"/task/priority/batch" method:"put" tags:"任务优先级" summary:"批量更新任务优先级"`
	TaskIds  []string `json:"task_ids" v:"required|max:100#任务ID列表不能为空|最多100个任务"`
	Priority int      `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string   `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type BatchUpdateTaskPriorityRes struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedTasks  []string `json:"failed_tasks"`
	Message      string   `json:"message"`
}
