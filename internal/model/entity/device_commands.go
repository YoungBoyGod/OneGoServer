// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceCommands is the golang structure for table device_commands.
type DeviceCommands struct {
	Id            int64       `json:"id"            orm:"id"             description:""` //
	DeviceId      string      `json:"deviceId"      orm:"device_id"      description:""` //
	CommandType   string      `json:"commandType"   orm:"command_type"   description:""` //
	CommandData   string      `json:"commandData"   orm:"command_data"   description:""` //
	Status        string      `json:"status"        orm:"status"         description:""` //
	SentTime      *gtime.Time `json:"sentTime"      orm:"sent_time"      description:""` //
	ExecutedTime  *gtime.Time `json:"executedTime"  orm:"executed_time"  description:""` //
	CompletedTime *gtime.Time `json:"completedTime" orm:"completed_time" description:""` //
	ResponseData  string      `json:"responseData"  orm:"response_data"  description:""` //
	ErrorMessage  string      `json:"errorMessage"  orm:"error_message"  description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:""` //
	CreatedBy     string      `json:"createdBy"     orm:"created_by"     description:""` //
}
