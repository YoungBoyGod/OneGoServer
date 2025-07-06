package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ==============================================
// 设备队列管理 API
// ==============================================

// GetDeviceQueueHistory 获取设备队列操作历史请求
type GetDeviceQueueHistoryReq struct {
	g.Meta        `path:"/device/{deviceId}/queue/history" method:"get" tags:"设备队列" summary:"获取设备队列操作历史"`
	DeviceId      string `json:"deviceId" v:"required#设备ID不能为空"`
	Page          int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size          int    `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	OperationType string `json:"operation_type,omitempty" v:"in:enqueue,dequeue,priority_change,cancel,restart#操作类型"`
	TaskId        string `json:"task_id,omitempty"`
	StartTime     string `json:"start_time,omitempty"`
	EndTime       string `json:"end_time,omitempty"`
	BatchId       string `json:"batch_id,omitempty"`
}

// GetDeviceQueueHistoryRes 获取设备队列操作历史响应
type GetDeviceQueueHistoryRes struct {
	List  []DeviceQueueHistoryInfo `json:"list"`
	Total int64                    `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}
