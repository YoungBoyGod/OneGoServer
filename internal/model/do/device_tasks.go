// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceTasks is the golang structure of table device_tasks for DAO operations like Where/Data.
type DeviceTasks struct {
	g.Meta    `orm:"table:device_tasks, do:true"`
	Id        interface{} //
	DeviceId  interface{} //
	TaskId    interface{} //
	CreatedAt *gtime.Time //
	CreatedBy interface{} //
}
