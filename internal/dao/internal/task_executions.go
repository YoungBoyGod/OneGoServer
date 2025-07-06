// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaskExecutionsDao is the data access object for the table task_executions.
type TaskExecutionsDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  TaskExecutionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// TaskExecutionsColumns defines and stores column names for the table task_executions.
type TaskExecutionsColumns struct {
	Id                string //
	TaskId            string //
	ExecutionId       string //
	DeviceId          string //
	Status            string //
	StartTime         string //
	EndTime           string //
	Duration          string //
	ExecutorInfo      string //
	Logs              string //
	Metrics           string //
	Output            string //
	ErrorDetails      string //
	CpuUsageAvg       string //
	CpuUsagePeak      string //
	MemoryUsageAvg    string //
	MemoryUsagePeak   string //
	IoOperationsTotal string //
	IoBytesTotal      string //
	CreatedAt         string //
	UpdatedAt         string //
}

// taskExecutionsColumns holds the columns for the table task_executions.
var taskExecutionsColumns = TaskExecutionsColumns{
	Id:                "id",
	TaskId:            "task_id",
	ExecutionId:       "execution_id",
	DeviceId:          "device_id",
	Status:            "status",
	StartTime:         "start_time",
	EndTime:           "end_time",
	Duration:          "duration",
	ExecutorInfo:      "executor_info",
	Logs:              "logs",
	Metrics:           "metrics",
	Output:            "output",
	ErrorDetails:      "error_details",
	CpuUsageAvg:       "cpu_usage_avg",
	CpuUsagePeak:      "cpu_usage_peak",
	MemoryUsageAvg:    "memory_usage_avg",
	MemoryUsagePeak:   "memory_usage_peak",
	IoOperationsTotal: "io_operations_total",
	IoBytesTotal:      "io_bytes_total",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewTaskExecutionsDao creates and returns a new DAO object for table data access.
func NewTaskExecutionsDao(handlers ...gdb.ModelHandler) *TaskExecutionsDao {
	return &TaskExecutionsDao{
		group:    "default",
		table:    "task_executions",
		columns:  taskExecutionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TaskExecutionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TaskExecutionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TaskExecutionsDao) Columns() TaskExecutionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TaskExecutionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TaskExecutionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TaskExecutionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
