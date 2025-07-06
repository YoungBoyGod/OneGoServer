// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TaskExecutions is the golang structure for table task_executions.
type TaskExecutions struct {
	Id                int64       `json:"id"                orm:"id"                  description:""` //
	TaskId            string      `json:"taskId"            orm:"task_id"             description:""` //
	ExecutionId       string      `json:"executionId"       orm:"execution_id"        description:""` //
	DeviceId          string      `json:"deviceId"          orm:"device_id"           description:""` //
	Status            string      `json:"status"            orm:"status"              description:""` //
	StartTime         *gtime.Time `json:"startTime"         orm:"start_time"          description:""` //
	EndTime           *gtime.Time `json:"endTime"           orm:"end_time"            description:""` //
	Duration          int         `json:"duration"          orm:"duration"            description:""` //
	ExecutorInfo      string      `json:"executorInfo"      orm:"executor_info"       description:""` //
	Logs              string      `json:"logs"              orm:"logs"                description:""` //
	Metrics           string      `json:"metrics"           orm:"metrics"             description:""` //
	Output            string      `json:"output"            orm:"output"              description:""` //
	ErrorDetails      string      `json:"errorDetails"      orm:"error_details"       description:""` //
	CpuUsageAvg       float64     `json:"cpuUsageAvg"       orm:"cpu_usage_avg"       description:""` //
	CpuUsagePeak      float64     `json:"cpuUsagePeak"      orm:"cpu_usage_peak"      description:""` //
	MemoryUsageAvg    float64     `json:"memoryUsageAvg"    orm:"memory_usage_avg"    description:""` //
	MemoryUsagePeak   float64     `json:"memoryUsagePeak"   orm:"memory_usage_peak"   description:""` //
	IoOperationsTotal int64       `json:"ioOperationsTotal" orm:"io_operations_total" description:""` //
	IoBytesTotal      int64       `json:"ioBytesTotal"      orm:"io_bytes_total"      description:""` //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""` //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""` //
}
