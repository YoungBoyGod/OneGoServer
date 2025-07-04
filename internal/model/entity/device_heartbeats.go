// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceHeartbeats is the golang structure for table device_heartbeats.
type DeviceHeartbeats struct {
	Id            int64       `json:"id"            orm:"id"             description:""` //
	DeviceId      string      `json:"deviceId"      orm:"device_id"      description:""` //
	HeartbeatTime *gtime.Time `json:"heartbeatTime" orm:"heartbeat_time" description:""` //
	Status        string      `json:"status"        orm:"status"         description:""` //
	Metadata      string      `json:"metadata"      orm:"metadata"       description:""` //
	IpAddress     string      `json:"ipAddress"     orm:"ip_address"     description:""` //
	ResponseTime  int         `json:"responseTime"  orm:"response_time"  description:""` //
}
