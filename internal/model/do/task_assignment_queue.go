// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskAssignmentQueue is the golang structure of table task_assignment_queue for DAO operations like Where/Data.
type TaskAssignmentQueue struct {
	g.Meta               `orm:"table:task_assignment_queue, do:true"`
	Id                   interface{} //
	TaskId               interface{} //
	Priority             interface{} //
	QueueStatus          interface{} //
	RequiredDeviceType   interface{} //
	RequiredCapabilities interface{} //
	PreferredDeviceIds   []int64     //
	ExcludedDeviceIds    []int64     //
	AssignmentStrategy   interface{} //
	AffinityRules        interface{} //
	AssignedDeviceId     interface{} //
	AssignedAt           *gtime.Time //
	AssignmentScore      interface{} //
	QueuePosition        interface{} //
	EstimatedWaitTime    interface{} //
	RetryCount           interface{} //
	MaxRetries           interface{} //
	OriginalPriority     interface{} //
	LastPriorityChangeAt *gtime.Time //
	PriorityChangeReason interface{} //
	PriorityBoostReason  interface{} //
	OperationSource      interface{} //
	QueuedAt             *gtime.Time //
	UpdatedAt            *gtime.Time //
}
