package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 设备任务管理 API
// ==============================================

// GetDeviceTaskList 获取设备任务列表请求
type GetDeviceTaskListReq struct {
	g.Meta   `path:"/device/{deviceId}/tasks" method:"get" tags:"设备任务" summary:"获取设备任务列表"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status   string `json:"status,omitempty" v:"in:pending,running,completed,failed,canceled#状态只能是pending,running,completed,failed,canceled"`
}

// GetDeviceTaskListRes 获取设备任务列表响应
type GetDeviceTaskListRes struct {
	List  []DeviceTaskInfo `json:"list"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// GetDeviceTaskQueue 获取设备任务队列请求
type GetDeviceTaskQueueReq struct {
	g.Meta   `path:"/device/{deviceId}/task-queue" method:"get" tags:"设备任务" summary:"获取设备任务队列"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size     int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
}

// GetDeviceTaskQueueRes 获取设备任务队列响应
type GetDeviceTaskQueueRes struct {
	List  []DeviceTaskQueueInfo `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// GetDeviceTaskDetail 获取设备任务详情请求
type GetDeviceTaskDetailReq struct {
	g.Meta   `path:"/device/{deviceId}/task/{taskId}" method:"get" tags:"设备任务" summary:"获取设备任务详情"`
	DeviceId string `json:"deviceId" v:"required#设备ID不能为空"`
	TaskId   string `json:"taskId" v:"required#任务ID不能为空"`
}

// GetDeviceTaskDetailRes 获取设备任务详情响应
type GetDeviceTaskDetailRes struct {
	TaskId        string      `json:"task_id"`
	DeviceId      string      `json:"device_id"`
	TaskName      string      `json:"task_name"`
	TaskType      string      `json:"task_type"`
	Description   string      `json:"description"`
	Status        string      `json:"status"`
	Priority      int         `json:"priority"`
	Progress      int         `json:"progress"`
	ResultData    string      `json:"result_data"`
	ErrorMessage  string      `json:"error_message"`
	RetryCount    int         `json:"retry_count"`
	MaxRetries    int         `json:"max_retries"`
	CreatedAt     *gtime.Time `json:"created_at"`
	StartTime     *gtime.Time `json:"start_time"`
	EndTime       *gtime.Time `json:"end_time"`
	EstimatedTime int64       `json:"estimated_time"`
	ActualTime    int64       `json:"actual_time"`
}
