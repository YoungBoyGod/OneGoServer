// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Devices is the golang structure of table devices for DAO operations like Where/Data.
type Devices struct {
	g.Meta               `orm:"table:devices, do:true"`
	Id                   interface{} //
	DeviceId             interface{} //
	Name                 interface{} //
	Type                 interface{} //
	Model                interface{} //
	BoardId              interface{} //
	Status               interface{} //
	HealthScore          interface{} //
	LoginUsername        interface{} //
	LoginPort            interface{} //
	LoginPublicKey       interface{} //
	IpAddress            interface{} //
	Port                 interface{} //
	Protocol             interface{} //
	Endpoint             interface{} //
	RegTime              *gtime.Time //
	Metadata             interface{} //
	Tags                 interface{} //
	UptimeHours          interface{} //
	FirstOnlineTime      *gtime.Time //
	LastOnlineTime       *gtime.Time //
	LastOfflineTime      *gtime.Time //
	TotalOnlineDuration  interface{} //
	TotalOfflineDuration interface{} //
	TotalHeartbeats      interface{} //
	TotalAlerts          interface{} //
	TotalTasks           interface{} //
	TotalSuccessTasks    interface{} //
	TotalFailedTasks     interface{} //
	TotalCanceledTasks   interface{} //
	TotalPendingTasks    interface{} //
	TotalRunningTasks    interface{} //
	TotalCompletedTasks  interface{} //
	CreatedAt            *gtime.Time //
	UpdatedAt            *gtime.Time //
	CreatedBy            interface{} //
	UpdatedBy            interface{} //
}
