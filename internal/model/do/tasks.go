// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tasks is the golang structure of table tasks for DAO operations like Where/Data.
type Tasks struct {
	g.Meta       `orm:"table:tasks, do:true"`
	Id           interface{} //
	DeviceId     interface{} //
	TaskId       interface{} //
	Name         interface{} //
	Description  interface{} //
	Type         interface{} //
	Status       interface{} //
	Priority     interface{} //
	ExecuteTime  *gtime.Time //
	Timeout      interface{} //
	RetryCount   interface{} //
	MaxRetries   interface{} //
	IsUrgent     interface{} //
	Parameters   interface{} //
	Result       interface{} //
	ErrorMessage interface{} //
	ExecutorType interface{} //
	ExecutorId   interface{} //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	CreatedBy    interface{} //
	UpdatedBy    interface{} //
}
