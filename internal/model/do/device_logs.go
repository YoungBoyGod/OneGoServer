// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceLogs is the golang structure of table device_logs for DAO operations like Where/Data.
type DeviceLogs struct {
	g.Meta        `orm:"table:device_logs, do:true"`
	Id            interface{} //
	DeviceId      interface{} //
	LogTime       *gtime.Time //
	Level         interface{} //
	Category      interface{} //
	Message       interface{} //
	Details       interface{} //
	Source        interface{} //
	CorrelationId interface{} //
}
