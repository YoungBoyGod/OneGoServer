// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeviceHeartbeatsDao is the data access object for the table device_heartbeats.
type DeviceHeartbeatsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  DeviceHeartbeatsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// DeviceHeartbeatsColumns defines and stores column names for the table device_heartbeats.
type DeviceHeartbeatsColumns struct {
	Id            string //
	DeviceId      string //
	HeartbeatTime string //
	Status        string //
	Metadata      string //
	IpAddress     string //
	ResponseTime  string //
}

// deviceHeartbeatsColumns holds the columns for the table device_heartbeats.
var deviceHeartbeatsColumns = DeviceHeartbeatsColumns{
	Id:            "id",
	DeviceId:      "device_id",
	HeartbeatTime: "heartbeat_time",
	Status:        "status",
	Metadata:      "metadata",
	IpAddress:     "ip_address",
	ResponseTime:  "response_time",
}

// NewDeviceHeartbeatsDao creates and returns a new DAO object for table data access.
func NewDeviceHeartbeatsDao(handlers ...gdb.ModelHandler) *DeviceHeartbeatsDao {
	return &DeviceHeartbeatsDao{
		group:    "default",
		table:    "device_heartbeats",
		columns:  deviceHeartbeatsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeviceHeartbeatsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeviceHeartbeatsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeviceHeartbeatsDao) Columns() DeviceHeartbeatsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeviceHeartbeatsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeviceHeartbeatsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeviceHeartbeatsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
