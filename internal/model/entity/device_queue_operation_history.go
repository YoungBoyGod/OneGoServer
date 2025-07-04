// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceQueueOperationHistory is the golang structure for table device_queue_operation_history.
type DeviceQueueOperationHistory struct {
	Id               int64       `json:"id"               orm:"id"                 description:""` //
	DeviceId         int64       `json:"deviceId"         orm:"device_id"          description:""` //
	TaskId           string      `json:"taskId"           orm:"task_id"            description:""` //
	OperationType    string      `json:"operationType"    orm:"operation_type"     description:""` //
	OperationBy      int64       `json:"operationBy"      orm:"operation_by"       description:""` //
	OperationTime    *gtime.Time `json:"operationTime"    orm:"operation_time"     description:""` //
	OldPriority      int         `json:"oldPriority"      orm:"old_priority"       description:""` //
	NewPriority      int         `json:"newPriority"      orm:"new_priority"       description:""` //
	OldPosition      int         `json:"oldPosition"      orm:"old_position"       description:""` //
	NewPosition      int         `json:"newPosition"      orm:"new_position"       description:""` //
	OldStatus        string      `json:"oldStatus"        orm:"old_status"         description:""` //
	NewStatus        string      `json:"newStatus"        orm:"new_status"         description:""` //
	Reason           string      `json:"reason"           orm:"reason"             description:""` //
	Notes            string      `json:"notes"            orm:"notes"              description:""` //
	OperationSource  string      `json:"operationSource"  orm:"operation_source"   description:""` //
	BatchId          string      `json:"batchId"          orm:"batch_id"           description:""` //
	IsBatchOperation bool        `json:"isBatchOperation" orm:"is_batch_operation" description:""` //
}
