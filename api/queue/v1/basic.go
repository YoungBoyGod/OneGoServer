package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 队列基础管理 API (5个)
// ===============================

// CreateQueueReq 创建队列请求
type CreateQueueReq struct {
	g.Meta         `path:"/queue/create" method:"post" tags:"队列管理" summary:"创建队列"`
	QueueName      string                 `json:"queue_name" v:"required|length:1,100#队列名称不能为空|队列名称长度为1-100字符"`
	QueueType      string                 `json:"queue_type" v:"required|in:fifo,priority,delay,lifo#队列类型不能为空"`
	Description    string                 `json:"description,omitempty" v:"max-length:500#队列描述最大500字符"`
	MaxSize        int                    `json:"max_size" d:"1000" v:"between:1,10000#队列最大大小为1-10000"`
	MaxConcurrency int                    `json:"max_concurrency" d:"10" v:"between:1,100#最大并发数为1-100"`
	Priority       int                    `json:"priority" d:"5" v:"between:1,10#队列优先级为1-10"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Enabled        bool                   `json:"enabled" d:"true"`
}

type CreateQueueRes struct {
	QueueId string `json:"queue_id"`
	Message string `json:"message"`
}

// GetQueueListReq 获取队列列表请求
type GetQueueListReq struct {
	g.Meta    `path:"/queue/list" method:"get" tags:"队列管理" summary:"获取队列列表"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status    string      `json:"status,omitempty" v:"in:active,paused,stopped,error#状态值无效"`
	QueueType string      `json:"queue_type,omitempty" v:"in:fifo,priority,delay,lifo#队列类型无效"`
	Priority  int         `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
	Keyword   string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	SortBy    string      `json:"sort_by" d:"created_at" v:"in:created_at,updated_at,priority,task_count#排序字段无效"`
	SortOrder string      `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"`
}

type GetQueueListRes struct {
	List  []QueueInfo `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// GetQueueDetailReq 获取队列详情请求
type GetQueueDetailReq struct {
	g.Meta  `path:"/queue/{queueId}" method:"get" tags:"队列管理" summary:"获取队列详情"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueDetailRes struct {
	QueueDetail QueueDetailInfo `json:"queue_detail"`
}

// UpdateQueueReq 更新队列请求
type UpdateQueueReq struct {
	g.Meta         `path:"/queue/{queueId}" method:"put" tags:"队列管理" summary:"更新队列"`
	QueueId        string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	QueueName      string                 `json:"queue_name,omitempty" v:"length:1,100#队列名称长度为1-100字符"`
	Description    string                 `json:"description,omitempty" v:"max-length:500#队列描述最大500字符"`
	MaxSize        int                    `json:"max_size,omitempty" v:"between:1,10000#队列最大大小为1-10000"`
	MaxConcurrency int                    `json:"max_concurrency,omitempty" v:"between:1,100#最大并发数为1-100"`
	Priority       int                    `json:"priority,omitempty" v:"between:1,10#队列优先级为1-10"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Enabled        bool                   `json:"enabled,omitempty"`
}

type UpdateQueueRes struct {
	Message string `json:"message"`
}

// DeleteQueueReq 删除队列请求
type DeleteQueueReq struct {
	g.Meta  `path:"/queue/{queueId}" method:"delete" tags:"队列管理" summary:"删除队列"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Force   bool   `json:"force,omitempty"`
}

type DeleteQueueRes struct {
	Message string `json:"message"`
}
