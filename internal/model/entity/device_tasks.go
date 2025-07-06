// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceTasks is the golang structure for table device_tasks.
type DeviceTasks struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	DeviceId  string      `json:"deviceId"  orm:"device_id"  description:""` //
	TaskId    string      `json:"taskId"    orm:"task_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	CreatedBy string      `json:"createdBy" orm:"created_by" description:""` //
}
