// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TasksDao is the data access object for the table tasks.
type TasksDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TasksColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TasksColumns defines and stores column names for the table tasks.
type TasksColumns struct {
	Id           string //
	TaskId       string //
	Name         string //
	Description  string //
	Type         string //
	Status       string //
	Priority     string //
	ExecuteTime  string //
	Timeout      string //
	RetryCount   string //
	MaxRetries   string //
	IsUrgent     string //
	Parameters   string //
	Result       string //
	ErrorMessage string //
	ExecutorType string //
	ExecutorId   string //
	DeviceId     string //
	CreatedAt    string //
	UpdatedAt    string //
	CreatedBy    string //
	UpdatedBy    string //
}

// tasksColumns holds the columns for the table tasks.
var tasksColumns = TasksColumns{
	Id:           "id",
	TaskId:       "task_id",
	Name:         "name",
	Description:  "description",
	Type:         "type",
	Status:       "status",
	Priority:     "priority",
	ExecuteTime:  "execute_time",
	Timeout:      "timeout",
	RetryCount:   "retry_count",
	MaxRetries:   "max_retries",
	IsUrgent:     "is_urgent",
	Parameters:   "parameters",
	Result:       "result",
	ErrorMessage: "error_message",
	ExecutorType: "executor_type",
	ExecutorId:   "executor_id",
	DeviceId:     "device_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	CreatedBy:    "created_by",
	UpdatedBy:    "updated_by",
}

// NewTasksDao creates and returns a new DAO object for table data access.
func NewTasksDao(handlers ...gdb.ModelHandler) *TasksDao {
	return &TasksDao{
		group:    "default",
		table:    "tasks",
		columns:  tasksColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TasksDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TasksDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TasksDao) Columns() TasksColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TasksDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TasksDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TasksDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
