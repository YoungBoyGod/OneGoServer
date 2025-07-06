// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskAssignmentHistory is the golang structure for table task_assignment_history.
type TaskAssignmentHistory struct {
	Id              int64       `json:"id"              orm:"id"               description:""` //
	TaskId          string      `json:"taskId"          orm:"task_id"          description:""` //
	DeviceId        string      `json:"deviceId"        orm:"device_id"        description:""` //
	Action          string      `json:"action"          orm:"action"           description:""` //
	PreviousStatus  string      `json:"previousStatus"  orm:"previous_status"  description:""` //
	NewStatus       string      `json:"newStatus"       orm:"new_status"       description:""` //
	Reason          string      `json:"reason"          orm:"reason"           description:""` //
	Details         string      `json:"details"         orm:"details"          description:""` //
	OperationSource string      `json:"operationSource" orm:"operation_source" description:""` //
	OperatorId      string      `json:"operatorId"      orm:"operator_id"      description:""` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""` //
}
