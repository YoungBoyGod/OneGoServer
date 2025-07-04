// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceLogs is the golang structure for table device_logs.
type DeviceLogs struct {
	Id            int64       `json:"id"            orm:"id"             description:""` //
	DeviceId      string      `json:"deviceId"      orm:"device_id"      description:""` //
	LogTime       *gtime.Time `json:"logTime"       orm:"log_time"       description:""` //
	Level         string      `json:"level"         orm:"level"          description:""` //
	Category      string      `json:"category"      orm:"category"       description:""` //
	Message       string      `json:"message"       orm:"message"        description:""` //
	Details       string      `json:"details"       orm:"details"        description:""` //
	Source        string      `json:"source"        orm:"source"         description:""` //
	CorrelationId string      `json:"correlationId" orm:"correlation_id" description:""` //
}
