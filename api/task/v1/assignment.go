package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务分配相关API
// ===============================

// AssignTaskToDeviceReq 分配任务给设备请求
type AssignTaskToDeviceReq struct {
	g.Meta    `path:"/task/assign" method:"post" tags:"任务管理" summary:"分配任务给设备"`
	TaskId    string   `json:"taskId" v:"required#任务ID不能为空"`
	DeviceIds []string `json:"deviceIds" v:"required#设备ID列表不能为空"`
	Strategy  string   `json:"strategy" d:"auto" v:"in:auto,manual,load_balanced,priority#分配策略只能是auto,manual,load_balanced,priority"`
	Force     bool     `json:"force" d:"false"`
}

// AssignTaskToDeviceRes 分配任务给设备响应
type AssignTaskToDeviceRes struct {
	TaskId           string  `json:"taskId"`
	AssignedDeviceId string  `json:"assignedDeviceId"`
	AssignmentScore  float64 `json:"assignmentScore"`
	Reason           string  `json:"reason"`
	Status           string  `json:"status"`
}

// GetTaskAssignmentOptionsReq 获取任务分配选项请求
type GetTaskAssignmentOptionsReq struct {
	g.Meta `path:"/task/assignment/options" method:"get" tags:"任务管理" summary:"获取任务分配选项"`
	TaskId string `json:"taskId" v:"required#任务ID不能为空"`
}

// GetTaskAssignmentOptionsRes 获取任务分配选项响应
type GetTaskAssignmentOptionsRes struct {
	TaskId  string                   `json:"taskId"`
	Options []DeviceAssignmentOption `json:"options"`
	Total   int                      `json:"total"`
}

// BatchAssignTasksReq 批量分配任务请求
type BatchAssignTasksReq struct {
	g.Meta      `path:"/task/assign/batch" method:"post" tags:"任务管理" summary:"批量分配任务"`
	Assignments []TaskAssignment `json:"assignments" v:"required#分配列表不能为空"`
	Strategy    string           `json:"strategy" d:"auto" v:"in:auto,manual,load_balanced,priority#分配策略只能是auto,manual,load_balanced,priority"`
}

// BatchAssignTasksRes 批量分配任务响应
type BatchAssignTasksRes struct {
	Total       int                    `json:"total"`
	Success     int                    `json:"success"`
	Failed      int                    `json:"failed"`
	Assignments []TaskAssignmentResult `json:"assignments"`
}

// GetTaskAssignmentHistoryReq 获取任务分配历史请求
type GetTaskAssignmentHistoryReq struct {
	g.Meta                   `path:"/task/assignment/history" method:"get" tags:"任务管理" summary:"获取任务分配历史"`
	TaskId                   string `json:"taskId,omitempty"`
	DeviceId                 string `json:"deviceId,omitempty"`
	StartTime                string `json:"startTime,omitempty"`
	EndTime                  string `json:"endTime,omitempty"`
	common.PaginationRequest `json:",inline"`
}

// GetTaskAssignmentHistoryRes 获取任务分配历史响应
type GetTaskAssignmentHistoryRes struct {
	common.PaginationResponse[TaskAssignmentHistory] `json:",inline"`
}

// OptimizeTaskAssignmentReq 优化任务分配请求
type OptimizeTaskAssignmentReq struct {
	g.Meta   `path:"/task/assignment/optimize" method:"post" tags:"任务管理" summary:"优化任务分配"`
	TaskIds  []string `json:"taskIds" v:"required#任务ID列表不能为空"`
	Strategy string   `json:"strategy" d:"load_balanced" v:"in:load_balanced,priority,round_robin,least_loaded#优化策略只能是load_balanced,priority,round_robin,least_loaded"`
}

// OptimizeTaskAssignmentRes 优化任务分配响应
type OptimizeTaskAssignmentRes struct {
	OriginalAssignments  []TaskAssignmentResult `json:"originalAssignments"`
	OptimizedAssignments []TaskAssignmentResult `json:"optimizedAssignments"`
	Improvement          float64                `json:"improvement"`
	Actions              []string               `json:"actions"`
	Status               string                 `json:"status"`
}

// GetDeviceTaskQueueReq 获取设备任务队列请求
type GetDeviceTaskQueueReq struct {
	g.Meta                   `path:"/task/queue/device" method:"get" tags:"任务管理" summary:"获取设备任务队列"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	Status                   string `json:"status,omitempty"`
	common.PaginationRequest `json:",inline"`
}

// GetDeviceTaskQueueRes 获取设备任务队列响应
type GetDeviceTaskQueueRes struct {
	DeviceId                                       string                `json:"deviceId"`
	Queue                                          []DeviceTaskQueueItem `json:"queue"`
	common.PaginationResponse[DeviceTaskQueueItem] `json:",inline"`
}

// ===============================
// 任务分配管理相关API
// ===============================

// AssignTaskReq 分配任务请求
type AssignTaskReq struct {
	g.Meta    `path:"/task/{taskId}/assign" method:"post" tags:"任务分配" summary:"分配任务"`
	TaskId    string   `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	DeviceIds []string `json:"device_ids" v:"required|max:10#设备ID列表不能为空|最多10个设备"`
	Priority  int      `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
	Deadline  string   `json:"deadline,omitempty" v:"datetime#截止时间格式不正确"`
}

type AssignTaskRes struct {
	Message string `json:"message"`
}

// GetTaskAssignmentsReq 获取任务分配请求
type GetTaskAssignmentsReq struct {
	g.Meta                   `path:"/task/{taskId}/assignments" method:"get" tags:"任务分配" summary:"获取任务分配"`
	TaskId                   string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	common.PaginationRequest `json:",inline"`
	Status                   string      `json:"status,omitempty" v:"in:pending,running,completed,failed#状态值无效"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
}

type GetTaskAssignmentsRes struct {
	common.PaginationResponse[TaskAssignmentInfo] `json:",inline"`
}

// ReassignTaskReq 重新分配任务请求
type ReassignTaskReq struct {
	g.Meta       `path:"/task/{taskId}/reassign" method:"post" tags:"任务分配" summary:"重新分配任务"`
	TaskId       string   `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	OldDeviceIds []string `json:"old_device_ids" v:"required|max:10#原设备ID列表不能为空"`
	NewDeviceIds []string `json:"new_device_ids" v:"required|max:10#新设备ID列表不能为空"`
	Reason       string   `json:"reason,omitempty" v:"max-length:200#重新分配原因最大200字符"`
}

type ReassignTaskRes struct {
	Message string `json:"message"`
}

// GetDeviceTaskAssignmentsReq 获取设备任务分配请求
type GetDeviceTaskAssignmentsReq struct {
	g.Meta                   `path:"/device/{deviceId}/task-assignments" method:"get" tags:"任务分配" summary:"获取设备任务分配"`
	DeviceId                 string `json:"device_id" v:"required|max-length:50#设备ID不能为空"`
	common.PaginationRequest `json:",inline"`
	Status                   string      `json:"status,omitempty" v:"in:pending,running,completed,failed#状态值无效"`
	TaskType                 string      `json:"task_type,omitempty" v:"in:backup,sync,monitor,custom#任务类型无效"`
	StartTime                *gtime.Time `json:"start_time,omitempty"`
	EndTime                  *gtime.Time `json:"end_time,omitempty"`
}

type GetDeviceTaskAssignmentsRes struct {
	common.PaginationResponse[DeviceTaskAssignmentInfo] `json:",inline"`
}
