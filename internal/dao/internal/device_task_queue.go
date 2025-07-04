// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceTaskQueueDao is the data access object for the table device_task_queue.
type DeviceTaskQueueDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  DeviceTaskQueueColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// DeviceTaskQueueColumns defines and stores column names for the table device_task_queue.
type DeviceTaskQueueColumns struct {
	Id                   string //
	DeviceId             string //
	TaskId               string //
	QueuePriority        string //
	OriginalPriority     string //
	IsManualPriority     string //
	QueuePosition        string //
	IsManualPosition     string //
	Status               string //
	EstimatedStartTime   string //
	EstimatedDuration    string //
	ActualStartTime      string //
	ActualEndTime        string //
	MaxRetryCount        string //
	CurrentRetry         string //
	TimeoutSeconds       string //
	DependsOnTaskIds     string //
	BlocksTaskIds        string //
	RequeueCount         string //
	LastRequeueAt        string //
	IsRequeued           string //
	CancelReason         string //
	LastPriorityChangeAt string //
	LastPositionChangeAt string //
	QueuedBy             string //
	LastModifiedBy       string //
	LastAction           string //
	CreatedAt            string //
	UpdatedAt            string //
}

// deviceTaskQueueColumns holds the columns for the table device_task_queue.
var deviceTaskQueueColumns = DeviceTaskQueueColumns{
	Id:                   "id",
	DeviceId:             "device_id",
	TaskId:               "task_id",
	QueuePriority:        "queue_priority",
	OriginalPriority:     "original_priority",
	IsManualPriority:     "is_manual_priority",
	QueuePosition:        "queue_position",
	IsManualPosition:     "is_manual_position",
	Status:               "status",
	EstimatedStartTime:   "estimated_start_time",
	EstimatedDuration:    "estimated_duration",
	ActualStartTime:      "actual_start_time",
	ActualEndTime:        "actual_end_time",
	MaxRetryCount:        "max_retry_count",
	CurrentRetry:         "current_retry",
	TimeoutSeconds:       "timeout_seconds",
	DependsOnTaskIds:     "depends_on_task_ids",
	BlocksTaskIds:        "blocks_task_ids",
	RequeueCount:         "requeue_count",
	LastRequeueAt:        "last_requeue_at",
	IsRequeued:           "is_requeued",
	CancelReason:         "cancel_reason",
	LastPriorityChangeAt: "last_priority_change_at",
	LastPositionChangeAt: "last_position_change_at",
	QueuedBy:             "queued_by",
	LastModifiedBy:       "last_modified_by",
	LastAction:           "last_action",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewDeviceTaskQueueDao creates and returns a new DAO object for table data access.
func NewDeviceTaskQueueDao(handlers ...gdb.ModelHandler) *DeviceTaskQueueDao {
	return &DeviceTaskQueueDao{
		group:    "default",
		table:    "device_task_queue",
		columns:  deviceTaskQueueColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceTaskQueueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceTaskQueueDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceTaskQueueDao) Columns() DeviceTaskQueueColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceTaskQueueDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceTaskQueueDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeviceTaskQueueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
