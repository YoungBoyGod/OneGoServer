// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Devices is the golang structure for table devices.
type Devices struct {
	Id                   int64       `json:"id"                   orm:"id"                     description:""` //
	DeviceId             string      `json:"deviceId"             orm:"device_id"              description:""` //
	Name                 string      `json:"name"                 orm:"name"                   description:""` //
	Type                 string      `json:"type"                 orm:"type"                   description:""` //
	Model                string      `json:"model"                orm:"model"                  description:""` //
	BoardId              string      `json:"boardId"              orm:"board_id"               description:""` //
	Status               string      `json:"status"               orm:"status"                 description:""` //
	HealthScore          int         `json:"healthScore"          orm:"health_score"           description:""` //
	LoginUsername        string      `json:"loginUsername"        orm:"login_username"         description:""` //
	LoginPort            int         `json:"loginPort"            orm:"login_port"             description:""` //
	LoginPublicKey       string      `json:"loginPublicKey"       orm:"login_public_key"       description:""` //
	IpAddress            string      `json:"ipAddress"            orm:"ip_address"             description:""` //
	Port                 int         `json:"port"                 orm:"port"                   description:""` //
	Protocol             string      `json:"protocol"             orm:"protocol"               description:""` //
	Endpoint             string      `json:"endpoint"             orm:"endpoint"               description:""` //
	RegTime              *gtime.Time `json:"regTime"              orm:"reg_time"               description:""` //
	Metadata             string      `json:"metadata"             orm:"metadata"               description:""` //
	Tags                 string      `json:"tags"                 orm:"tags"                   description:""` //
	UptimeHours          float64     `json:"uptimeHours"          orm:"uptime_hours"           description:""` //
	FirstOnlineTime      *gtime.Time `json:"firstOnlineTime"      orm:"first_online_time"      description:""` //
	LastOnlineTime       *gtime.Time `json:"lastOnlineTime"       orm:"last_online_time"       description:""` //
	LastOfflineTime      *gtime.Time `json:"lastOfflineTime"      orm:"last_offline_time"      description:""` //
	TotalOnlineDuration  int64       `json:"totalOnlineDuration"  orm:"total_online_duration"  description:""` //
	TotalOfflineDuration int64       `json:"totalOfflineDuration" orm:"total_offline_duration" description:""` //
	TotalHeartbeats      int64       `json:"totalHeartbeats"      orm:"total_heartbeats"       description:""` //
	TotalAlerts          int64       `json:"totalAlerts"          orm:"total_alerts"           description:""` //
	TotalTasks           int64       `json:"totalTasks"           orm:"total_tasks"            description:""` //
	TotalSuccessTasks    int64       `json:"totalSuccessTasks"    orm:"total_success_tasks"    description:""` //
	TotalFailedTasks     int64       `json:"totalFailedTasks"     orm:"total_failed_tasks"     description:""` //
	TotalCanceledTasks   int64       `json:"totalCanceledTasks"   orm:"total_canceled_tasks"   description:""` //
	TotalPendingTasks    int64       `json:"totalPendingTasks"    orm:"total_pending_tasks"    description:""` //
	TotalRunningTasks    int64       `json:"totalRunningTasks"    orm:"total_running_tasks"    description:""` //
	TotalCompletedTasks  int64       `json:"totalCompletedTasks"  orm:"total_completed_tasks"  description:""` //
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:""` //
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:""` //
	CreatedBy            string      `json:"createdBy"            orm:"created_by"             description:""` //
	UpdatedBy            string      `json:"updatedBy"            orm:"updated_by"             description:""` //
}
