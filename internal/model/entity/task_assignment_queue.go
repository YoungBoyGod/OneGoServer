// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskAssignmentQueue is the golang structure for table task_assignment_queue.
type TaskAssignmentQueue struct {
	Id                   int64       `json:"id"                   orm:"id"                      description:""` //
	TaskId               string      `json:"taskId"               orm:"task_id"                 description:""` //
	Priority             int         `json:"priority"             orm:"priority"                description:""` //
	QueueStatus          string      `json:"queueStatus"          orm:"queue_status"            description:""` //
	RequiredDeviceType   string      `json:"requiredDeviceType"   orm:"required_device_type"    description:""` //
	RequiredCapabilities string      `json:"requiredCapabilities" orm:"required_capabilities"   description:""` //
	PreferredDeviceIds   []int64     `json:"preferredDeviceIds"   orm:"preferred_device_ids"    description:""` //
	ExcludedDeviceIds    []int64     `json:"excludedDeviceIds"    orm:"excluded_device_ids"     description:""` //
	AssignmentStrategy   string      `json:"assignmentStrategy"   orm:"assignment_strategy"     description:""` //
	AffinityRules        string      `json:"affinityRules"        orm:"affinity_rules"          description:""` //
	AssignedDeviceId     string      `json:"assignedDeviceId"     orm:"assigned_device_id"      description:""` //
	AssignedAt           *gtime.Time `json:"assignedAt"           orm:"assigned_at"             description:""` //
	AssignmentScore      float64     `json:"assignmentScore"      orm:"assignment_score"        description:""` //
	QueuePosition        int         `json:"queuePosition"        orm:"queue_position"          description:""` //
	EstimatedWaitTime    int         `json:"estimatedWaitTime"    orm:"estimated_wait_time"     description:""` //
	RetryCount           int         `json:"retryCount"           orm:"retry_count"             description:""` //
	MaxRetries           int         `json:"maxRetries"           orm:"max_retries"             description:""` //
	OriginalPriority     int         `json:"originalPriority"     orm:"original_priority"       description:""` //
	LastPriorityChangeAt *gtime.Time `json:"lastPriorityChangeAt" orm:"last_priority_change_at" description:""` //
	PriorityChangeReason string      `json:"priorityChangeReason" orm:"priority_change_reason"  description:""` //
	PriorityBoostReason  string      `json:"priorityBoostReason"  orm:"priority_boost_reason"   description:""` //
	OperationSource      string      `json:"operationSource"      orm:"operation_source"        description:""` //
	QueuedAt             *gtime.Time `json:"queuedAt"             orm:"queued_at"               description:""` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              description:""` //
}
