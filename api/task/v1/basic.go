package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务基础管理 API (5个)
// ===============================

// CreateTaskReq 创建任务请求
type CreateTaskReq struct {
	g.Meta      `path:"/task/create" method:"post" tags:"任务管理" summary:"创建任务"`
	TaskName    string                 `json:"task_name" v:"required|length:1,100#任务名称不能为空|任务名称长度为1-100字符"`
	TaskType    string                 `json:"task_type" v:"required|in:backup,sync,monitor,custom#任务类型不能为空"`
	Description string                 `json:"description,omitempty" v:"max-length:500#任务描述最大500字符"`
	Priority    int                    `json:"priority" d:"5" v:"between:1,10#任务优先级为1-10"`
	DeviceId    string                 `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	Config      map[string]interface{} `json:"config,omitempty"`
	ScheduleAt  *gtime.Time            `json:"schedule_at,omitempty"`
	DeadlineAt  *gtime.Time            `json:"deadline_at,omitempty"`
	RetryCount  int                    `json:"retry_count" d:"3" v:"between:0,10#重试次数为0-10"`
}

type CreateTaskRes struct {
	TaskId  string `json:"task_id"`
	Message string `json:"message"`
}

// GetTaskListReq 获取任务列表请求
type GetTaskListReq struct {
	g.Meta                   `path:"/task/list" method:"get" tags:"任务管理" summary:"获取任务列表"`
	common.PaginationRequest `json:",inline"`
	Status                   string      `json:"status,omitempty" v:"in:pending,running,completed,failed,cancelled#状态值无效"`
	TaskType                 string      `json:"task_type,omitempty" v:"in:backup,sync,monitor,custom#任务类型无效"`
	DeviceId                 string      `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	Priority                 int         `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
	Keyword                  string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
}

type GetTaskListRes struct {
	common.PaginationResponse[TaskInfo] `json:",inline"`
}

// GetTaskDetailReq 获取任务详情请求
type GetTaskDetailReq struct {
	g.Meta `path:"/task/{taskId}" method:"get" tags:"任务管理" summary:"获取任务详情"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
}

type GetTaskDetailRes struct {
	TaskDetail TaskDetailInfo `json:"task_detail"`
}

// UpdateTaskReq 更新任务请求
type UpdateTaskReq struct {
	g.Meta      `path:"/task/{taskId}" method:"put" tags:"任务管理" summary:"更新任务"`
	TaskId      string                 `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	TaskName    string                 `json:"task_name,omitempty" v:"length:1,100#任务名称长度为1-100字符"`
	Description string                 `json:"description,omitempty" v:"max-length:500#任务描述最大500字符"`
	Priority    int                    `json:"priority,omitempty" v:"between:1,10#任务优先级为1-10"`
	Config      map[string]interface{} `json:"config,omitempty"`
	ScheduleAt  *gtime.Time            `json:"schedule_at,omitempty"`
	DeadlineAt  *gtime.Time            `json:"deadline_at,omitempty"`
	RetryCount  int                    `json:"retry_count,omitempty" v:"between:0,10#重试次数为0-10"`
}

type UpdateTaskRes struct {
	Message string `json:"message"`
}

// DeleteTaskReq 删除任务请求
type DeleteTaskReq struct {
	g.Meta `path:"/task/{taskId}" method:"delete" tags:"任务管理" summary:"删除任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Force  bool   `json:"force,omitempty"`
}

type DeleteTaskRes struct {
	Message string `json:"message"`
}
