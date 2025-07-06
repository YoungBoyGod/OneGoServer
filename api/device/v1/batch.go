package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ==============================================
// 批量操作管理 API
// ==============================================

// BatchOperateDevices 批量操作设备请求
type BatchOperateDevicesReq struct {
	g.Meta      `path:"/device/batch" method:"post" tags:"设备管理" summary:"批量操作设备"`
	DeviceIds   []string `json:"device_ids" v:"required|length:1,100#设备ID列表不能为空|最多支持100个设备"`
	Operation   string   `json:"operation" v:"required|in:start,stop,restart,delete,update_status,send_command#操作类型不能为空"`
	Parameters  string   `json:"parameters,omitempty"`
	BatchReason string   `json:"batch_reason,omitempty"`
}

// BatchOperateDevicesRes 批量操作设备响应
type BatchOperateDevicesRes struct {
	BatchId      string                 `json:"batch_id"`
	TotalCount   int                    `json:"total_count"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	SkippedCount int                    `json:"skipped_count"`
	Results      []BatchOperationResult `json:"results"`
	Message      string                 `json:"message"`
}

// GetBatchOperationStatus 获取批量操作状态请求
type GetBatchOperationStatusReq struct {
	g.Meta  `path:"/device/batch/{batchId}" method:"get" tags:"设备管理" summary:"获取批量操作状态"`
	BatchId string `json:"batchId" v:"required#批次ID不能为空"`
}

// GetBatchOperationStatusRes 获取批量操作状态响应
type GetBatchOperationStatusRes struct {
	BatchId        string                 `json:"batch_id"`
	Operation      string                 `json:"operation"`
	Status         string                 `json:"status"` // pending, processing, completed, failed
	TotalCount     int                    `json:"total_count"`
	ProcessedCount int                    `json:"processed_count"`
	SuccessCount   int                    `json:"success_count"`
	FailedCount    int                    `json:"failed_count"`
	Progress       float64                `json:"progress"`
	StartTime      *gtime.Time            `json:"start_time"`
	EndTime        *gtime.Time            `json:"end_time"`
	Results        []BatchOperationResult `json:"results"`
}
