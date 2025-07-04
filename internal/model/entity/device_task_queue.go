// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceTaskQueue is the golang structure for table device_task_queue.
type DeviceTaskQueue struct {
	Id                   int64       `json:"id"                   orm:"id"                      description:""` //
	DeviceId             int64       `json:"deviceId"             orm:"device_id"               description:""` //
	TaskId               string      `json:"taskId"               orm:"task_id"                 description:""` //
	QueuePriority        int         `json:"queuePriority"        orm:"queue_priority"          description:""` //
	OriginalPriority     int         `json:"originalPriority"     orm:"original_priority"       description:""` //
	IsManualPriority     bool        `json:"isManualPriority"     orm:"is_manual_priority"      description:""` //
	QueuePosition        int         `json:"queuePosition"        orm:"queue_position"          description:""` //
	IsManualPosition     bool        `json:"isManualPosition"     orm:"is_manual_position"      description:""` //
	Status               string      `json:"status"               orm:"status"                  description:""` //
	EstimatedStartTime   *gtime.Time `json:"estimatedStartTime"   orm:"estimated_start_time"    description:""` //
	EstimatedDuration    int         `json:"estimatedDuration"    orm:"estimated_duration"      description:""` //
	ActualStartTime      *gtime.Time `json:"actualStartTime"      orm:"actual_start_time"       description:""` //
	ActualEndTime        *gtime.Time `json:"actualEndTime"        orm:"actual_end_time"         description:""` //
	MaxRetryCount        int         `json:"maxRetryCount"        orm:"max_retry_count"         description:""` //
	CurrentRetry         int         `json:"currentRetry"         orm:"current_retry"           description:""` //
	TimeoutSeconds       int         `json:"timeoutSeconds"       orm:"timeout_seconds"         description:""` //
	DependsOnTaskIds     []string    `json:"dependsOnTaskIds"     orm:"depends_on_task_ids"     description:""` //
	BlocksTaskIds        []string    `json:"blocksTaskIds"        orm:"blocks_task_ids"         description:""` //
	RequeueCount         int         `json:"requeueCount"         orm:"requeue_count"           description:""` //
	LastRequeueAt        *gtime.Time `json:"lastRequeueAt"        orm:"last_requeue_at"         description:""` //
	IsRequeued           bool        `json:"isRequeued"           orm:"is_requeued"             description:""` //
	CancelReason         string      `json:"cancelReason"         orm:"cancel_reason"           description:""` //
	LastPriorityChangeAt *gtime.Time `json:"lastPriorityChangeAt" orm:"last_priority_change_at" description:""` //
	LastPositionChangeAt *gtime.Time `json:"lastPositionChangeAt" orm:"last_position_change_at" description:""` //
	QueuedBy             int64       `json:"queuedBy"             orm:"queued_by"               description:""` //
	LastModifiedBy       int64       `json:"lastModifiedBy"       orm:"last_modified_by"        description:""` //
	LastAction           string      `json:"lastAction"           orm:"last_action"             description:""` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"              description:""` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              description:""` //
}
