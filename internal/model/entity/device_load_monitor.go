// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DeviceLoadMonitor is the golang structure for table device_load_monitor.
type DeviceLoadMonitor struct {
	Id                 int64       `json:"id"                 orm:"id"                   description:""` //
	DeviceId           string      `json:"deviceId"           orm:"device_id"            description:""` //
	CurrentTasks       int         `json:"currentTasks"       orm:"current_tasks"        description:""` //
	MaxConcurrentTasks int         `json:"maxConcurrentTasks" orm:"max_concurrent_tasks" description:""` //
	CpuLoad            float64     `json:"cpuLoad"            orm:"cpu_load"             description:""` //
	MemoryUsage        float64     `json:"memoryUsage"        orm:"memory_usage"         description:""` //
	DiskUsage          float64     `json:"diskUsage"          orm:"disk_usage"           description:""` //
	NetworkLatency     int         `json:"networkLatency"     orm:"network_latency"      description:""` //
	Status             string      `json:"status"             orm:"status"               description:""` //
	LastHeartbeat      *gtime.Time `json:"lastHeartbeat"      orm:"last_heartbeat"       description:""` //
	LoadScore          float64     `json:"loadScore"          orm:"load_score"           description:""` //
	TotalAssigned      int         `json:"totalAssigned"      orm:"total_assigned"       description:""` //
	TotalCompleted     int         `json:"totalCompleted"     orm:"total_completed"      description:""` //
	TotalFailed        int         `json:"totalFailed"        orm:"total_failed"         description:""` //
	SuccessRate        float64     `json:"successRate"        orm:"success_rate"         description:""` //
	AvgTaskDuration    float64     `json:"avgTaskDuration"    orm:"avg_task_duration"    description:""` //
	LastTaskCompletion *gtime.Time `json:"lastTaskCompletion" orm:"last_task_completion" description:""` //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:""` //
}
