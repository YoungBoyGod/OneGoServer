package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 任务队列管理 API (4个)
// ===============================

// GetTaskQueueListReq 获取任务队列列表请求
type GetTaskQueueListReq struct {
	g.Meta    `path:"/task/queue/list" method:"get" tags:"任务队列" summary:"获取任务队列列表"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	QueueType string `json:"queue_type,omitempty" v:"in:pending,running,priority#队列类型无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
}

type GetTaskQueueListRes struct {
	List  []TaskQueueInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// GetTaskQueueDetailReq 获取任务队列详情请求
type GetTaskQueueDetailReq struct {
	g.Meta  `path:"/task/queue/{queueId}" method:"get" tags:"任务队列" summary:"获取任务队列详情"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetTaskQueueDetailRes struct {
	QueueDetail TaskQueueDetailInfo `json:"queue_detail"`
}

// ManageTaskQueueReq 管理任务队列请求
type ManageTaskQueueReq struct {
	g.Meta   `path:"/task/queue/{queueId}/manage" method:"post" tags:"任务队列" summary:"管理任务队列"`
	QueueId  string   `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Action   string   `json:"action" v:"required|in:pause,resume,clear,reorder#操作类型无效"`
	TaskIds  []string `json:"task_ids,omitempty"`
	Priority int      `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
}

type ManageTaskQueueRes struct {
	Message string `json:"message"`
}

// GetTaskQueueStatsReq 获取任务队列统计请求
type GetTaskQueueStatsReq struct {
	g.Meta    `path:"/task/queue/stats" method:"get" tags:"任务队列" summary:"获取任务队列统计"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
}

type GetTaskQueueStatsRes struct {
	QueueStats TaskQueueStats `json:"queue_stats"`
}
