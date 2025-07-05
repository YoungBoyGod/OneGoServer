package v1

import (
	"github.com/gogf/gf/v2/frame/g"
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

// DeviceAssignmentOption 设备分配选项
type DeviceAssignmentOption struct {
	DeviceId           string  `json:"deviceId"`
	DeviceName         string  `json:"deviceName"`
	DeviceType         string  `json:"deviceType"`
	LoadScore          float64 `json:"loadScore"`
	CompatibilityScore float64 `json:"compatibilityScore"`
	AssignmentScore    float64 `json:"assignmentScore"`
	CurrentTasks       int     `json:"currentTasks"`
	MaxTasks           int     `json:"maxTasks"`
	EstimatedWaitTime  int     `json:"estimatedWaitTime"`
	Status             string  `json:"status"`
	Recommendation     string  `json:"recommendation"`
}

// BatchAssignTasksReq 批量分配任务请求
type BatchAssignTasksReq struct {
	g.Meta      `path:"/task/assign/batch" method:"post" tags:"任务管理" summary:"批量分配任务"`
	Assignments []TaskAssignment `json:"assignments" v:"required#分配列表不能为空"`
	Strategy    string           `json:"strategy" d:"auto" v:"in:auto,manual,load_balanced,priority#分配策略只能是auto,manual,load_balanced,priority"`
}

// TaskAssignment 任务分配
type TaskAssignment struct {
	TaskId    string   `json:"taskId" v:"required#任务ID不能为空"`
	DeviceIds []string `json:"deviceIds" v:"required#设备ID列表不能为空"`
	Priority  int      `json:"priority" d:"0"`
}

// BatchAssignTasksRes 批量分配任务响应
type BatchAssignTasksRes struct {
	Total       int                    `json:"total"`
	Success     int                    `json:"success"`
	Failed      int                    `json:"failed"`
	Assignments []TaskAssignmentResult `json:"assignments"`
}

// TaskAssignmentResult 任务分配结果
type TaskAssignmentResult struct {
	TaskId           string `json:"taskId"`
	AssignedDeviceId string `json:"assignedDeviceId"`
	Status           string `json:"status"`
	Error            string `json:"error,omitempty"`
}

// GetTaskAssignmentHistoryReq 获取任务分配历史请求
type GetTaskAssignmentHistoryReq struct {
	g.Meta    `path:"/task/assignment/history" method:"get" tags:"任务管理" summary:"获取任务分配历史"`
	TaskId    string `json:"taskId,omitempty"`
	DeviceId  string `json:"deviceId,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
}

// GetTaskAssignmentHistoryRes 获取任务分配历史响应
type GetTaskAssignmentHistoryRes struct {
	List  []TaskAssignmentHistory `json:"list"`
	Total int64                   `json:"total"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
}

// TaskAssignmentHistory 任务分配历史
type TaskAssignmentHistory struct {
	Id              int64   `json:"id"`
	TaskId          string  `json:"taskId"`
	DeviceId        string  `json:"deviceId"`
	AssignmentScore float64 `json:"assignmentScore"`
	Strategy        string  `json:"strategy"`
	Reason          string  `json:"reason"`
	AssignedAt      string  `json:"assignedAt"`
	CompletedAt     string  `json:"completedAt,omitempty"`
	Status          string  `json:"status"`
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
	g.Meta   `path:"/task/queue/device" method:"get" tags:"任务管理" summary:"获取设备任务队列"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Status   string `json:"status,omitempty"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
}

// GetDeviceTaskQueueRes 获取设备任务队列响应
type GetDeviceTaskQueueRes struct {
	DeviceId string                `json:"deviceId"`
	Queue    []DeviceTaskQueueItem `json:"queue"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	Size     int                   `json:"size"`
}

// DeviceTaskQueueItem 设备任务队列项
type DeviceTaskQueueItem struct {
	Id                 int64  `json:"id"`
	TaskId             string `json:"taskId"`
	TaskName           string `json:"taskName"`
	TaskType           string `json:"taskType"`
	Priority           int    `json:"priority"`
	QueuePosition      int    `json:"queuePosition"`
	EstimatedStartTime string `json:"estimatedStartTime"`
	EstimatedDuration  int    `json:"estimatedDuration"`
	Status             string `json:"status"`
	QueuedAt           string `json:"queuedAt"`
}
