// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceLoadMonitor is the golang structure of table device_load_monitor for DAO operations like Where/Data.
type DeviceLoadMonitor struct {
	g.Meta             `orm:"table:device_load_monitor, do:true"`
	Id                 interface{} //
	DeviceId           interface{} //
	CurrentTasks       interface{} //
	MaxConcurrentTasks interface{} //
	CpuLoad            interface{} //
	MemoryUsage        interface{} //
	DiskUsage          interface{} //
	NetworkLatency     interface{} //
	Status             interface{} //
	LastHeartbeat      *gtime.Time //
	LoadScore          interface{} //
	TotalAssigned      interface{} //
	TotalCompleted     interface{} //
	TotalFailed        interface{} //
	SuccessRate        interface{} //
	AvgTaskDuration    interface{} //
	LastTaskCompletion *gtime.Time //
	UpdatedAt          *gtime.Time //
}
