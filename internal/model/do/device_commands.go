// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceCommands is the golang structure of table device_commands for DAO operations like Where/Data.
type DeviceCommands struct {
	g.Meta        `orm:"table:device_commands, do:true"`
	Id            interface{} //
	DeviceId      interface{} //
	CommandType   interface{} //
	CommandData   interface{} //
	Status        interface{} //
	SentTime      *gtime.Time //
	ExecutedTime  *gtime.Time //
	CompletedTime *gtime.Time //
	ResponseData  interface{} //
	ErrorMessage  interface{} //
	CreatedAt     *gtime.Time //
	CreatedBy     interface{} //
}
