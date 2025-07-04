// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceLoadMonitorDao is the data access object for the table device_load_monitor.
type DeviceLoadMonitorDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  DeviceLoadMonitorColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// DeviceLoadMonitorColumns defines and stores column names for the table device_load_monitor.
type DeviceLoadMonitorColumns struct {
	Id                 string //
	DeviceId           string //
	CurrentTasks       string //
	MaxConcurrentTasks string //
	CpuLoad            string //
	MemoryUsage        string //
	DiskUsage          string //
	NetworkLatency     string //
	Status             string //
	LastHeartbeat      string //
	LoadScore          string //
	TotalAssigned      string //
	TotalCompleted     string //
	TotalFailed        string //
	SuccessRate        string //
	AvgTaskDuration    string //
	LastTaskCompletion string //
	UpdatedAt          string //
}

// deviceLoadMonitorColumns holds the columns for the table device_load_monitor.
var deviceLoadMonitorColumns = DeviceLoadMonitorColumns{
	Id:                 "id",
	DeviceId:           "device_id",
	CurrentTasks:       "current_tasks",
	MaxConcurrentTasks: "max_concurrent_tasks",
	CpuLoad:            "cpu_load",
	MemoryUsage:        "memory_usage",
	DiskUsage:          "disk_usage",
	NetworkLatency:     "network_latency",
	Status:             "status",
	LastHeartbeat:      "last_heartbeat",
	LoadScore:          "load_score",
	TotalAssigned:      "total_assigned",
	TotalCompleted:     "total_completed",
	TotalFailed:        "total_failed",
	SuccessRate:        "success_rate",
	AvgTaskDuration:    "avg_task_duration",
	LastTaskCompletion: "last_task_completion",
	UpdatedAt:          "updated_at",
}

// NewDeviceLoadMonitorDao creates and returns a new DAO object for table data access.
func NewDeviceLoadMonitorDao(handlers ...gdb.ModelHandler) *DeviceLoadMonitorDao {
	return &DeviceLoadMonitorDao{
		group:    "default",
		table:    "device_load_monitor",
		columns:  deviceLoadMonitorColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceLoadMonitorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceLoadMonitorDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceLoadMonitorDao) Columns() DeviceLoadMonitorColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceLoadMonitorDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceLoadMonitorDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DeviceLoadMonitorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
