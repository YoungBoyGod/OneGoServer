package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务调度管理 API (4个)
// ===============================

// GetTaskScheduleListReq 获取任务调度列表请求
type GetTaskScheduleListReq struct {
	g.Meta                   `path:"/task/schedule/list" method:"get" tags:"任务调度" summary:"获取任务调度列表"`
	common.PaginationRequest `json:",inline"`
	ScheduleType             string `json:"schedule_type,omitempty" v:"in:once,recurring,cron#调度类型无效"`
	Status                   string `json:"status,omitempty" v:"in:active,inactive,paused#状态无效"`
}

type GetTaskScheduleListRes struct {
	common.PaginationResponse[TaskScheduleInfo] `json:",inline"`
}

// CreateTaskScheduleReq 创建任务调度请求
type CreateTaskScheduleReq struct {
	g.Meta       `path:"/task/schedule/create" method:"post" tags:"任务调度" summary:"创建任务调度"`
	ScheduleName string                 `json:"schedule_name" v:"required|length:1,100#调度名称不能为空"`
	TaskTemplate map[string]interface{} `json:"task_template" v:"required#任务模板不能为空"`
	ScheduleType string                 `json:"schedule_type" v:"required|in:once,recurring,cron#调度类型无效"`
	CronExpr     string                 `json:"cron_expr,omitempty" v:"max-length:100#Cron表达式最大100字符"`
	StartTime    *gtime.Time            `json:"start_time,omitempty"`
	EndTime      *gtime.Time            `json:"end_time,omitempty"`
	Interval     int                    `json:"interval,omitempty" v:"min:1#间隔时间最小为1"`
	Enabled      bool                   `json:"enabled" d:"true"`
}

type CreateTaskScheduleRes struct {
	ScheduleId string `json:"schedule_id"`
	Message    string `json:"message"`
}

// UpdateTaskScheduleReq 更新任务调度请求
type UpdateTaskScheduleReq struct {
	g.Meta       `path:"/task/schedule/{scheduleId}" method:"put" tags:"任务调度" summary:"更新任务调度"`
	ScheduleId   string                 `json:"schedule_id" v:"required|max-length:50#调度ID不能为空"`
	ScheduleName string                 `json:"schedule_name,omitempty" v:"length:1,100#调度名称长度为1-100字符"`
	TaskTemplate map[string]interface{} `json:"task_template,omitempty"`
	CronExpr     string                 `json:"cron_expr,omitempty" v:"max-length:100#Cron表达式最大100字符"`
	StartTime    *gtime.Time            `json:"start_time,omitempty"`
	EndTime      *gtime.Time            `json:"end_time,omitempty"`
	Interval     int                    `json:"interval,omitempty" v:"min:1#间隔时间最小为1"`
	Enabled      bool                   `json:"enabled,omitempty"`
}

type UpdateTaskScheduleRes struct {
	Message string `json:"message"`
}

// DeleteTaskScheduleReq 删除任务调度请求
type DeleteTaskScheduleReq struct {
	g.Meta     `path:"/task/schedule/{scheduleId}" method:"delete" tags:"任务调度" summary:"删除任务调度"`
	ScheduleId string `json:"schedule_id" v:"required|max-length:50#调度ID不能为空"`
	Force      bool   `json:"force,omitempty"`
}

type DeleteTaskScheduleRes struct {
	Message string `json:"message"`
}
