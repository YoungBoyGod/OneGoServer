// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceHeartbeats is the golang structure of table device_heartbeats for DAO operations like Where/Data.
type DeviceHeartbeats struct {
	g.Meta        `orm:"table:device_heartbeats, do:true"`
	Id            interface{} //
	DeviceId      interface{} //
	HeartbeatTime *gtime.Time //
	Status        interface{} //
	Metadata      interface{} //
	IpAddress     interface{} //
	ResponseTime  interface{} //
}
