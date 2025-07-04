// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tasks is the golang structure for table tasks.
type Tasks struct {
	Id           int64       `json:"id"           orm:"id"            description:""` //
	TaskId       string      `json:"taskId"       orm:"task_id"       description:""` //
	Name         string      `json:"name"         orm:"name"          description:""` //
	Description  string      `json:"description"  orm:"description"   description:""` //
	Type         string      `json:"type"         orm:"type"          description:""` //
	Status       string      `json:"status"       orm:"status"        description:""` //
	Priority     int         `json:"priority"     orm:"priority"      description:""` //
	ExecuteTime  *gtime.Time `json:"executeTime"  orm:"execute_time"  description:""` //
	Timeout      int         `json:"timeout"      orm:"timeout"       description:""` //
	RetryCount   int         `json:"retryCount"   orm:"retry_count"   description:""` //
	MaxRetries   int         `json:"maxRetries"   orm:"max_retries"   description:""` //
	IsUrgent     bool        `json:"isUrgent"     orm:"is_urgent"     description:""` //
	Parameters   string      `json:"parameters"   orm:"parameters"    description:""` //
	Result       string      `json:"result"       orm:"result"        description:""` //
	ErrorMessage string      `json:"errorMessage" orm:"error_message" description:""` //
	ExecutorType string      `json:"executorType" orm:"executor_type" description:""` //
	ExecutorId   string      `json:"executorId"   orm:"executor_id"   description:""` //
	DeviceId     int64       `json:"deviceId"     orm:"device_id"     description:""` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""` //
	CreatedBy    int64       `json:"createdBy"    orm:"created_by"    description:""` //
	UpdatedBy    int64       `json:"updatedBy"    orm:"updated_by"    description:""` //
}
