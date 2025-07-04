// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskExecutions is the golang structure of table task_executions for DAO operations like Where/Data.
type TaskExecutions struct {
	g.Meta            `orm:"table:task_executions, do:true"`
	Id                interface{} //
	TaskId            interface{} //
	ExecutionId       interface{} //
	DeviceEsn         interface{} //
	Status            interface{} //
	StartTime         *gtime.Time //
	EndTime           *gtime.Time //
	Duration          interface{} //
	ExecutorInfo      interface{} //
	Logs              interface{} //
	Metrics           interface{} //
	Output            interface{} //
	ErrorDetails      interface{} //
	CpuUsageAvg       interface{} //
	CpuUsagePeak      interface{} //
	MemoryUsageAvg    interface{} //
	MemoryUsagePeak   interface{} //
	IoOperationsTotal interface{} //
	IoBytesTotal      interface{} //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}
