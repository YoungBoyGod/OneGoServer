// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TaskAssignmentQueueDao is the data access object for the table task_assignment_queue.
type TaskAssignmentQueueDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  TaskAssignmentQueueColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// TaskAssignmentQueueColumns defines and stores column names for the table task_assignment_queue.
type TaskAssignmentQueueColumns struct {
	Id                   string //
	TaskId               string //
	Priority             string //
	QueueStatus          string //
	RequiredDeviceType   string //
	RequiredCapabilities string //
	PreferredDeviceIds   string //
	ExcludedDeviceIds    string //
	AssignmentStrategy   string //
	AffinityRules        string //
	AssignedDeviceId     string //
	AssignedAt           string //
	AssignmentScore      string //
	QueuePosition        string //
	EstimatedWaitTime    string //
	RetryCount           string //
	MaxRetries           string //
	OriginalPriority     string //
	LastPriorityChangeAt string //
	PriorityChangeReason string //
	PriorityBoostReason  string //
	OperationSource      string //
	QueuedAt             string //
	UpdatedAt            string //
}

// taskAssignmentQueueColumns holds the columns for the table task_assignment_queue.
var taskAssignmentQueueColumns = TaskAssignmentQueueColumns{
	Id:                   "id",
	TaskId:               "task_id",
	Priority:             "priority",
	QueueStatus:          "queue_status",
	RequiredDeviceType:   "required_device_type",
	RequiredCapabilities: "required_capabilities",
	PreferredDeviceIds:   "preferred_device_ids",
	ExcludedDeviceIds:    "excluded_device_ids",
	AssignmentStrategy:   "assignment_strategy",
	AffinityRules:        "affinity_rules",
	AssignedDeviceId:     "assigned_device_id",
	AssignedAt:           "assigned_at",
	AssignmentScore:      "assignment_score",
	QueuePosition:        "queue_position",
	EstimatedWaitTime:    "estimated_wait_time",
	RetryCount:           "retry_count",
	MaxRetries:           "max_retries",
	OriginalPriority:     "original_priority",
	LastPriorityChangeAt: "last_priority_change_at",
	PriorityChangeReason: "priority_change_reason",
	PriorityBoostReason:  "priority_boost_reason",
	OperationSource:      "operation_source",
	QueuedAt:             "queued_at",
	UpdatedAt:            "updated_at",
}

// NewTaskAssignmentQueueDao creates and returns a new DAO object for table data access.
func NewTaskAssignmentQueueDao(handlers ...gdb.ModelHandler) *TaskAssignmentQueueDao {
	return &TaskAssignmentQueueDao{
		group:    "default",
		table:    "task_assignment_queue",
		columns:  taskAssignmentQueueColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TaskAssignmentQueueDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TaskAssignmentQueueDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TaskAssignmentQueueDao) Columns() TaskAssignmentQueueColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TaskAssignmentQueueDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TaskAssignmentQueueDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TaskAssignmentQueueDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
