// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskAssignmentHistory is the golang structure of table task_assignment_history for DAO operations like Where/Data.
type TaskAssignmentHistory struct {
	g.Meta          `orm:"table:task_assignment_history, do:true"`
	Id              interface{} //
	TaskId          interface{} //
	DeviceId        interface{} //
	Action          interface{} //
	PreviousStatus  interface{} //
	NewStatus       interface{} //
	Reason          interface{} //
	Details         interface{} //
	OperationSource interface{} //
	OperatorId      interface{} //
	CreatedAt       *gtime.Time //
}
