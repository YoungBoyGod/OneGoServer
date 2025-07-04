// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceQueueOperationHistory is the golang structure of table device_queue_operation_history for DAO operations like Where/Data.
type DeviceQueueOperationHistory struct {
	g.Meta           `orm:"table:device_queue_operation_history, do:true"`
	Id               interface{} //
	DeviceId         interface{} //
	TaskId           interface{} //
	OperationType    interface{} //
	OperationBy      interface{} //
	OperationTime    *gtime.Time //
	OldPriority      interface{} //
	NewPriority      interface{} //
	OldPosition      interface{} //
	NewPosition      interface{} //
	OldStatus        interface{} //
	NewStatus        interface{} //
	Reason           interface{} //
	Notes            interface{} //
	OperationSource  interface{} //
	BatchId          interface{} //
	IsBatchOperation interface{} //
}
