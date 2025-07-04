// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceQueueOperationHistoryDao is the data access object for the table device_queue_operation_history.
type DeviceQueueOperationHistoryDao struct {
	table    string                             // table is the underlying table name of the DAO.
	group    string                             // group is the database configuration group name of the current DAO.
	columns  DeviceQueueOperationHistoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler                 // handlers for customized model modification.
}

// DeviceQueueOperationHistoryColumns defines and stores column names for the table device_queue_operation_history.
type DeviceQueueOperationHistoryColumns struct {
	Id               string //
	DeviceId         string //
	TaskId           string //
	OperationType    string //
	OperationBy      string //
	OperationTime    string //
	OldPriority      string //
	NewPriority      string //
	OldPosition      string //
	NewPosition      string //
	OldStatus        string //
	NewStatus        string //
	Reason           string //
	Notes            string //
	OperationSource  string //
	BatchId          string //
	IsBatchOperation string //
}

// deviceQueueOperationHistoryColumns holds the columns for the table device_queue_operation_history.
var deviceQueueOperationHistoryColumns = DeviceQueueOperationHistoryColumns{
	Id:               "id",
	DeviceId:         "device_id",
	TaskId:           "task_id",
	OperationType:    "operation_type",
	OperationBy:      "operation_by",
	OperationTime:    "operation_time",
	OldPriority:      "old_priority",
	NewPriority:      "new_priority",
	OldPosition:      "old_position",
	NewPosition:      "new_position",
	OldStatus:        "old_status",
	NewStatus:        "new_status",
	Reason:           "reason",
	Notes:            "notes",
	OperationSource:  "operation_source",
	BatchId:          "batch_id",
	IsBatchOperation: "is_batch_operation",
}

// NewDeviceQueueOperationHistoryDao creates and returns a new DAO object for table data access.
func NewDeviceQueueOperationHistoryDao(handlers ...gdb.ModelHandler) *DeviceQueueOperationHistoryDao {
	return &DeviceQueueOperationHistoryDao{
		group:    "default",
		table:    "device_queue_operation_history",
		columns:  deviceQueueOperationHistoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceQueueOperationHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceQueueOperationHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceQueueOperationHistoryDao) Columns() DeviceQueueOperationHistoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceQueueOperationHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceQueueOperationHistoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeviceQueueOperationHistoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
