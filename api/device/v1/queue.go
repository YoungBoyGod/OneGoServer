package v1

import (
	"OneGfServer/api/common"

	"github.com/gogf/gf/v2/frame/g"
)

// ==============================================
// 设备队列管理 API
// ==============================================

// GetDeviceQueueHistory 获取设备队列操作历史请求
type GetDeviceQueueHistoryReq struct {
	g.Meta                   `path:"/device/{deviceId}/queue/history" method:"get" tags:"设备队列" summary:"获取设备队列操作历史"`
	DeviceId                 string `json:"deviceId" v:"required#设备ID不能为空"`
	common.PaginationRequest `json:",inline"`
	OperationType            string `json:"operation_type,omitempty" v:"in:enqueue,dequeue,priority_change,cancel,restart#操作类型"`
	TaskId                   string `json:"task_id,omitempty"`
	StartTime                string `json:"start_time,omitempty"`
	EndTime                  string `json:"end_time,omitempty"`
	BatchId                  string `json:"batch_id,omitempty"`
}

// GetDeviceQueueHistoryRes 获取设备队列操作历史响应
type GetDeviceQueueHistoryRes struct {
	common.PaginationResponse[DeviceQueueHistoryInfo] `json:",inline"`
}
