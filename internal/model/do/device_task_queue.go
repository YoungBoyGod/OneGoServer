// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceTaskQueue is the golang structure of table device_task_queue for DAO operations like Where/Data.
type DeviceTaskQueue struct {
	g.Meta               `orm:"table:device_task_queue, do:true"`
	Id                   interface{} //
	DeviceId             interface{} //
	TaskId               interface{} //
	QueuePriority        interface{} //
	OriginalPriority     interface{} //
	IsManualPriority     interface{} //
	QueuePosition        interface{} //
	IsManualPosition     interface{} //
	Status               interface{} //
	EstimatedStartTime   *gtime.Time //
	EstimatedDuration    interface{} //
	ActualStartTime      *gtime.Time //
	ActualEndTime        *gtime.Time //
	MaxRetryCount        interface{} //
	CurrentRetry         interface{} //
	TimeoutSeconds       interface{} //
	DependsOnTaskIds     []string    //
	BlocksTaskIds        []string    //
	RequeueCount         interface{} //
	LastRequeueAt        *gtime.Time //
	IsRequeued           interface{} //
	CancelReason         interface{} //
	LastPriorityChangeAt *gtime.Time //
	LastPositionChangeAt *gtime.Time //
	QueuedBy             interface{} //
	LastModifiedBy       interface{} //
	LastAction           interface{} //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
}
