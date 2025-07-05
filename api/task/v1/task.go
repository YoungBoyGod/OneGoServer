package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

/*
任务管理：
	1. 创建任务
	2. 获取任务列表
	3. 获取任务详情
	4. 更新任务
	5. 删除任务
	6. 取消任务
	7. 重试任务

任务优先级管理：
	1. 获取任务优先级列表
	2. 更新任务优先级
	3. 删除任务优先级

任务状态管理：
	1. 获取任务状态列表
	2. 更新任务状态

任务日志管理：
	1. 获取任务日志列表
	2. 获取任务日志详情
	3. 删除任务日志


任务队列管理：
	1. 获取任务队列列表
	2. 获取任务队列详情
	3. 删除任务队列

任务调度管理：
	1. 获取任务调度列表
	2. 更新任务调度
	3. 删除任务调度

*/

// ===============================
// 1. 任务基础管理 API (5个)
// ===============================

// CreateTaskReq 创建任务请求
type CreateTaskReq struct {
	g.Meta      `path:"/task/create" method:"post" tags:"任务管理" summary:"创建任务"`
	TaskName    string                 `json:"task_name" v:"required|length:1,100#任务名称不能为空|任务名称长度为1-100字符"`
	TaskType    string                 `json:"task_type" v:"required|in:backup,sync,monitor,custom#任务类型不能为空"`
	Description string                 `json:"description,omitempty" v:"max-length:500#任务描述最大500字符"`
	Priority    int                    `json:"priority" d:"5" v:"between:1,10#任务优先级为1-10"`
	DeviceId    string                 `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	Config      map[string]interface{} `json:"config,omitempty"`
	ScheduleAt  *gtime.Time            `json:"schedule_at,omitempty"`
	DeadlineAt  *gtime.Time            `json:"deadline_at,omitempty"`
	RetryCount  int                    `json:"retry_count" d:"3" v:"between:0,10#重试次数为0-10"`
}

type CreateTaskRes struct {
	TaskId  string `json:"task_id"`
	Message string `json:"message"`
}

// GetTaskListReq 获取任务列表请求
type GetTaskListReq struct {
	g.Meta    `path:"/task/list" method:"get" tags:"任务管理" summary:"获取任务列表"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	Status    string      `json:"status,omitempty" v:"in:pending,running,completed,failed,cancelled#状态值无效"`
	TaskType  string      `json:"task_type,omitempty" v:"in:backup,sync,monitor,custom#任务类型无效"`
	DeviceId  string      `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	Priority  int         `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	Keyword   string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
	SortBy    string      `json:"sort_by" d:"created_at" v:"in:created_at,updated_at,priority,deadline_at#排序字段无效"`
	SortOrder string      `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"`
}

type GetTaskListRes struct {
	List  []TaskInfo `json:"list"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

// GetTaskDetailReq 获取任务详情请求
type GetTaskDetailReq struct {
	g.Meta `path:"/task/{taskId}" method:"get" tags:"任务管理" summary:"获取任务详情"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
}

type GetTaskDetailRes struct {
	TaskDetail TaskDetailInfo `json:"task_detail"`
}

// UpdateTaskReq 更新任务请求
type UpdateTaskReq struct {
	g.Meta      `path:"/task/{taskId}" method:"put" tags:"任务管理" summary:"更新任务"`
	TaskId      string                 `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	TaskName    string                 `json:"task_name,omitempty" v:"length:1,100#任务名称长度为1-100字符"`
	Description string                 `json:"description,omitempty" v:"max-length:500#任务描述最大500字符"`
	Priority    int                    `json:"priority,omitempty" v:"between:1,10#任务优先级为1-10"`
	Config      map[string]interface{} `json:"config,omitempty"`
	ScheduleAt  *gtime.Time            `json:"schedule_at,omitempty"`
	DeadlineAt  *gtime.Time            `json:"deadline_at,omitempty"`
	RetryCount  int                    `json:"retry_count,omitempty" v:"between:0,10#重试次数为0-10"`
}

type UpdateTaskRes struct {
	Message string `json:"message"`
}

// DeleteTaskReq 删除任务请求
type DeleteTaskReq struct {
	g.Meta `path:"/task/{taskId}" method:"delete" tags:"任务管理" summary:"删除任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Force  bool   `json:"force,omitempty"`
}

type DeleteTaskRes struct {
	Message string `json:"message"`
}

// ===============================
// 2. 任务执行控制 API (5个)
// ===============================

// StartTaskReq 启动任务请求
type StartTaskReq struct {
	g.Meta `path:"/task/{taskId}/start" method:"post" tags:"任务执行" summary:"启动任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Force  bool   `json:"force,omitempty"`
}

type StartTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}

// StopTaskReq 停止任务请求
type StopTaskReq struct {
	g.Meta `path:"/task/{taskId}/stop" method:"post" tags:"任务执行" summary:"停止任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#停止原因最大200字符"`
}

type StopTaskRes struct {
	Message string `json:"message"`
}

// RestartTaskReq 重启任务请求
type RestartTaskReq struct {
	g.Meta `path:"/task/{taskId}/restart" method:"post" tags:"任务执行" summary:"重启任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#重启原因最大200字符"`
}

type RestartTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}

// CancelTaskReq 取消任务请求
type CancelTaskReq struct {
	g.Meta `path:"/task/{taskId}/cancel" method:"post" tags:"任务执行" summary:"取消任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#取消原因最大200字符"`
}

type CancelTaskRes struct {
	Message string `json:"message"`
}

// RetryTaskReq 重试任务请求
type RetryTaskReq struct {
	g.Meta `path:"/task/{taskId}/retry" method:"post" tags:"任务执行" summary:"重试任务"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Reason string `json:"reason,omitempty" v:"max-length:200#重试原因最大200字符"`
}

type RetryTaskRes struct {
	ExecutionId string `json:"execution_id"`
	Message     string `json:"message"`
}

// ===============================
// 3. 任务状态管理 API (2个)
// ===============================

// GetTaskStatusReq 获取任务状态请求
type GetTaskStatusReq struct {
	g.Meta `path:"/task/{taskId}/status" method:"get" tags:"任务状态" summary:"获取任务状态"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
}

type GetTaskStatusRes struct {
	TaskId         string      `json:"task_id"`
	Status         string      `json:"status"`
	Progress       float64     `json:"progress"`
	ExecutionId    string      `json:"execution_id"`
	StartTime      *gtime.Time `json:"start_time"`
	EndTime        *gtime.Time `json:"end_time"`
	Duration       int64       `json:"duration"`
	ErrorMessage   string      `json:"error_message"`
	ExecutionCount int         `json:"execution_count"`
	LastCheckTime  *gtime.Time `json:"last_check_time"`
}

// UpdateTaskStatusReq 更新任务状态请求
type UpdateTaskStatusReq struct {
	g.Meta       `path:"/task/{taskId}/status" method:"put" tags:"任务状态" summary:"更新任务状态"`
	TaskId       string  `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Status       string  `json:"status" v:"required|in:pending,running,completed,failed,cancelled#状态值无效"`
	Progress     float64 `json:"progress,omitempty" v:"between:0,100#进度为0-100"`
	ErrorMessage string  `json:"error_message,omitempty" v:"max-length:1000#错误信息最大1000字符"`
}

type UpdateTaskStatusRes struct {
	Message string `json:"message"`
}

// ===============================
// 4. 任务数据结构定义
// ===============================

// TaskInfo 任务基本信息
type TaskInfo struct {
	TaskId         string      `json:"task_id"`
	TaskName       string      `json:"task_name"`
	TaskType       string      `json:"task_type"`
	Description    string      `json:"description"`
	Status         string      `json:"status"`
	Priority       int         `json:"priority"`
	DeviceId       string      `json:"device_id"`
	DeviceName     string      `json:"device_name"`
	Progress       float64     `json:"progress"`
	ExecutionId    string      `json:"execution_id"`
	CreatedAt      *gtime.Time `json:"created_at"`
	UpdatedAt      *gtime.Time `json:"updated_at"`
	ScheduleAt     *gtime.Time `json:"schedule_at"`
	DeadlineAt     *gtime.Time `json:"deadline_at"`
	StartTime      *gtime.Time `json:"start_time"`
	EndTime        *gtime.Time `json:"end_time"`
	Duration       int64       `json:"duration"`
	RetryCount     int         `json:"retry_count"`
	ExecutionCount int         `json:"execution_count"`
}

// TaskDetailInfo 任务详细信息
type TaskDetailInfo struct {
	TaskInfo
	Config           map[string]interface{} `json:"config"`
	ErrorMessage     string                 `json:"error_message"`
	ExecutionHistory []TaskExecutionInfo    `json:"execution_history"`
	DependsOn        []string               `json:"depends_on"`
	DependedBy       []string               `json:"depended_by"`
	ResourceUsage    TaskResourceUsage      `json:"resource_usage"`
	Tags             []string               `json:"tags"`
	CreatedBy        string                 `json:"created_by"`
	UpdatedBy        string                 `json:"updated_by"`
	LastCheckTime    *gtime.Time            `json:"last_check_time"`
	NextScheduleTime *gtime.Time            `json:"next_schedule_time"`
}

// TaskExecutionInfo 任务执行信息
type TaskExecutionInfo struct {
	ExecutionId   string            `json:"execution_id"`
	Status        string            `json:"status"`
	StartTime     *gtime.Time       `json:"start_time"`
	EndTime       *gtime.Time       `json:"end_time"`
	Duration      int64             `json:"duration"`
	ErrorMessage  string            `json:"error_message"`
	Progress      float64           `json:"progress"`
	ExecutedBy    string            `json:"executed_by"`
	ResourceUsage TaskResourceUsage `json:"resource_usage"`
}

// TaskResourceUsage 任务资源使用情况
type TaskResourceUsage struct {
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskIO      float64 `json:"disk_io"`
	NetworkIO   float64 `json:"network_io"`
}

// ===============================
// 5. 任务优先级管理 API (3个)
// ===============================

// GetTaskPriorityListReq 获取任务优先级列表请求
type GetTaskPriorityListReq struct {
	g.Meta `path:"/task/priority/list" method:"get" tags:"任务优先级" summary:"获取任务优先级列表"`
}

type GetTaskPriorityListRes struct {
	PriorityList []TaskPriorityInfo `json:"priority_list"`
}

// UpdateTaskPriorityReq 更新任务优先级请求
type UpdateTaskPriorityReq struct {
	g.Meta   `path:"/task/{taskId}/priority" method:"put" tags:"任务优先级" summary:"更新任务优先级"`
	TaskId   string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Priority int    `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type UpdateTaskPriorityRes struct {
	Message string `json:"message"`
}

// BatchUpdateTaskPriorityReq 批量更新任务优先级请求
type BatchUpdateTaskPriorityReq struct {
	g.Meta   `path:"/task/priority/batch" method:"put" tags:"任务优先级" summary:"批量更新任务优先级"`
	TaskIds  []string `json:"task_ids" v:"required|max:100#任务ID列表不能为空|最多100个任务"`
	Priority int      `json:"priority" v:"required|between:1,10#任务优先级为1-10"`
	Reason   string   `json:"reason,omitempty" v:"max-length:200#调整原因最大200字符"`
}

type BatchUpdateTaskPriorityRes struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedTasks  []string `json:"failed_tasks"`
	Message      string   `json:"message"`
}

// ===============================
// 6. 任务日志管理 API (4个)
// ===============================

// GetTaskLogListReq 获取任务日志列表请求
type GetTaskLogListReq struct {
	g.Meta    `path:"/task/{taskId}/logs" method:"get" tags:"任务日志" summary:"获取任务日志列表"`
	TaskId    string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	Level     string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	Keyword   string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
}

type GetTaskLogListRes struct {
	List  []TaskLogInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// GetTaskLogDetailReq 获取任务日志详情请求
type GetTaskLogDetailReq struct {
	g.Meta `path:"/task/{taskId}/log/{logId}" method:"get" tags:"任务日志" summary:"获取任务日志详情"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	LogId  string `json:"log_id" v:"required|max-length:50#日志ID不能为空"`
}

type GetTaskLogDetailRes struct {
	LogDetail TaskLogDetailInfo `json:"log_detail"`
}

// ClearTaskLogsReq 清空任务日志请求
type ClearTaskLogsReq struct {
	g.Meta     `path:"/task/{taskId}/logs" method:"delete" tags:"任务日志" summary:"清空任务日志"`
	TaskId     string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	BeforeTime *gtime.Time `json:"before_time,omitempty"`
	Level      string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
}

type ClearTaskLogsRes struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ExportTaskLogsReq 导出任务日志请求
type ExportTaskLogsReq struct {
	g.Meta    `path:"/task/{taskId}/logs/export" method:"post" tags:"任务日志" summary:"导出任务日志"`
	TaskId    string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	Level     string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
	Format    string      `json:"format" d:"csv" v:"in:csv,json,txt#导出格式无效"`
}

type ExportTaskLogsRes struct {
	ExportId    string `json:"export_id"`
	DownloadUrl string `json:"download_url"`
	Message     string `json:"message"`
}

// ===============================
// 7. 任务队列管理 API (4个)
// ===============================

// GetTaskQueueListReq 获取任务队列列表请求
type GetTaskQueueListReq struct {
	g.Meta    `path:"/task/queue/list" method:"get" tags:"任务队列" summary:"获取任务队列列表"`
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	QueueType string `json:"queue_type,omitempty" v:"in:pending,running,priority#队列类型无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
}

type GetTaskQueueListRes struct {
	List  []TaskQueueInfo `json:"list"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

// GetTaskQueueDetailReq 获取任务队列详情请求
type GetTaskQueueDetailReq struct {
	g.Meta  `path:"/task/queue/{queueId}" method:"get" tags:"任务队列" summary:"获取任务队列详情"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetTaskQueueDetailRes struct {
	QueueDetail TaskQueueDetailInfo `json:"queue_detail"`
}

// ManageTaskQueueReq 管理任务队列请求
type ManageTaskQueueReq struct {
	g.Meta   `path:"/task/queue/{queueId}/manage" method:"post" tags:"任务队列" summary:"管理任务队列"`
	QueueId  string   `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Action   string   `json:"action" v:"required|in:pause,resume,clear,reorder#操作类型无效"`
	TaskIds  []string `json:"task_ids,omitempty"`
	Priority int      `json:"priority,omitempty" v:"between:1,10#优先级为1-10"`
}

type ManageTaskQueueRes struct {
	Message string `json:"message"`
}

// GetTaskQueueStatsReq 获取任务队列统计请求
type GetTaskQueueStatsReq struct {
	g.Meta    `path:"/task/queue/stats" method:"get" tags:"任务队列" summary:"获取任务队列统计"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
}

type GetTaskQueueStatsRes struct {
	QueueStats TaskQueueStats `json:"queue_stats"`
}

// ===============================
// 8. 任务调度管理 API (4个)
// ===============================

// GetTaskScheduleListReq 获取任务调度列表请求
type GetTaskScheduleListReq struct {
	g.Meta       `path:"/task/schedule/list" method:"get" tags:"任务调度" summary:"获取任务调度列表"`
	Page         int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size         int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`
	ScheduleType string `json:"schedule_type,omitempty" v:"in:once,recurring,cron#调度类型无效"`
	Status       string `json:"status,omitempty" v:"in:active,inactive,paused#状态无效"`
}

type GetTaskScheduleListRes struct {
	List  []TaskScheduleInfo `json:"list"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

// CreateTaskScheduleReq 创建任务调度请求
type CreateTaskScheduleReq struct {
	g.Meta       `path:"/task/schedule/create" method:"post" tags:"任务调度" summary:"创建任务调度"`
	ScheduleName string                 `json:"schedule_name" v:"required|length:1,100#调度名称不能为空"`
	TaskTemplate map[string]interface{} `json:"task_template" v:"required#任务模板不能为空"`
	ScheduleType string                 `json:"schedule_type" v:"required|in:once,recurring,cron#调度类型无效"`
	CronExpr     string                 `json:"cron_expr,omitempty" v:"max-length:100#Cron表达式最大100字符"`
	StartTime    *gtime.Time            `json:"start_time,omitempty"`
	EndTime      *gtime.Time            `json:"end_time,omitempty"`
	Interval     int                    `json:"interval,omitempty" v:"min:1#间隔时间最小为1"`
	Enabled      bool                   `json:"enabled" d:"true"`
}

type CreateTaskScheduleRes struct {
	ScheduleId string `json:"schedule_id"`
	Message    string `json:"message"`
}

// UpdateTaskScheduleReq 更新任务调度请求
type UpdateTaskScheduleReq struct {
	g.Meta       `path:"/task/schedule/{scheduleId}" method:"put" tags:"任务调度" summary:"更新任务调度"`
	ScheduleId   string                 `json:"schedule_id" v:"required|max-length:50#调度ID不能为空"`
	ScheduleName string                 `json:"schedule_name,omitempty" v:"length:1,100#调度名称长度为1-100字符"`
	TaskTemplate map[string]interface{} `json:"task_template,omitempty"`
	CronExpr     string                 `json:"cron_expr,omitempty" v:"max-length:100#Cron表达式最大100字符"`
	StartTime    *gtime.Time            `json:"start_time,omitempty"`
	EndTime      *gtime.Time            `json:"end_time,omitempty"`
	Interval     int                    `json:"interval,omitempty" v:"min:1#间隔时间最小为1"`
	Enabled      bool                   `json:"enabled,omitempty"`
}

type UpdateTaskScheduleRes struct {
	Message string `json:"message"`
}

// DeleteTaskScheduleReq 删除任务调度请求
type DeleteTaskScheduleReq struct {
	g.Meta     `path:"/task/schedule/{scheduleId}" method:"delete" tags:"任务调度" summary:"删除任务调度"`
	ScheduleId string `json:"schedule_id" v:"required|max-length:50#调度ID不能为空"`
	Force      bool   `json:"force,omitempty"`
}

type DeleteTaskScheduleRes struct {
	Message string `json:"message"`
}

// ===============================
// 9. 任务统计分析 API (2个)
// ===============================

// GetTaskStatisticsReq 获取任务统计信息请求
type GetTaskStatisticsReq struct {
	g.Meta    `path:"/task/statistics" method:"get" tags:"任务统计" summary:"获取任务统计信息"`
	TimeRange string `json:"time_range" d:"24h" v:"in:1h,6h,24h,7d,30d#时间范围无效"`
	TaskType  string `json:"task_type,omitempty" v:"in:backup,sync,monitor,custom#任务类型无效"`
	DeviceId  string `json:"device_id,omitempty" v:"max-length:50#设备ID最大50字符"`
	GroupBy   string `json:"group_by" d:"status" v:"in:status,type,device,priority#分组字段无效"`
}

type GetTaskStatisticsRes struct {
	Statistics TaskStatistics `json:"statistics"`
}

// GetTaskPerformanceReportReq 获取任务性能报告请求
type GetTaskPerformanceReportReq struct {
	g.Meta    `path:"/task/{taskId}/performance" method:"get" tags:"任务统计" summary:"获取任务性能报告"`
	TaskId    string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	TimeRange string `json:"time_range" d:"7d" v:"in:24h,7d,30d#时间范围无效"`
}

type GetTaskPerformanceReportRes struct {
	PerformanceReport TaskPerformanceReport `json:"performance_report"`
}

// ===============================
// 10. 扩展数据结构定义
// ===============================

// TaskPriorityInfo 任务优先级信息
type TaskPriorityInfo struct {
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TaskCount   int    `json:"task_count"`
	Color       string `json:"color"`
}

// TaskLogInfo 任务日志信息
type TaskLogInfo struct {
	LogId       string      `json:"log_id"`
	TaskId      string      `json:"task_id"`
	ExecutionId string      `json:"execution_id"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	Timestamp   *gtime.Time `json:"timestamp"`
	Source      string      `json:"source"`
}

// TaskLogDetailInfo 任务日志详细信息
type TaskLogDetailInfo struct {
	TaskLogInfo
	StackTrace string                 `json:"stack_trace"`
	Context    map[string]interface{} `json:"context"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// TaskQueueInfo 任务队列信息
type TaskQueueInfo struct {
	QueueId   string      `json:"queue_id"`
	QueueName string      `json:"queue_name"`
	QueueType string      `json:"queue_type"`
	DeviceId  string      `json:"device_id"`
	TaskCount int         `json:"task_count"`
	Status    string      `json:"status"`
	Priority  int         `json:"priority"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

// TaskQueueDetailInfo 任务队列详细信息
type TaskQueueDetailInfo struct {
	TaskQueueInfo
	Tasks           []TaskInfo  `json:"tasks"`
	MaxConcurrency  int         `json:"max_concurrency"`
	CurrentRunning  int         `json:"current_running"`
	ProcessingSpeed float64     `json:"processing_speed"`
	AverageWaitTime float64     `json:"average_wait_time"`
	LastProcessTime *gtime.Time `json:"last_process_time"`
}

// TaskQueueStats 任务队列统计信息
type TaskQueueStats struct {
	TotalQueues     int     `json:"total_queues"`
	TotalTasks      int     `json:"total_tasks"`
	PendingTasks    int     `json:"pending_tasks"`
	RunningTasks    int     `json:"running_tasks"`
	CompletedTasks  int     `json:"completed_tasks"`
	FailedTasks     int     `json:"failed_tasks"`
	AverageWaitTime float64 `json:"average_wait_time"`
	Throughput      float64 `json:"throughput"`
}

// TaskScheduleInfo 任务调度信息
type TaskScheduleInfo struct {
	ScheduleId   string      `json:"schedule_id"`
	ScheduleName string      `json:"schedule_name"`
	ScheduleType string      `json:"schedule_type"`
	CronExpr     string      `json:"cron_expr"`
	Status       string      `json:"status"`
	Enabled      bool        `json:"enabled"`
	TaskCount    int         `json:"task_count"`
	LastRun      *gtime.Time `json:"last_run"`
	NextRun      *gtime.Time `json:"next_run"`
	CreatedAt    *gtime.Time `json:"created_at"`
	UpdatedAt    *gtime.Time `json:"updated_at"`
}

// TaskStatistics 任务统计信息
type TaskStatistics struct {
	TotalTasks           int               `json:"total_tasks"`
	StatusDistribution   map[string]int    `json:"status_distribution"`
	TypeDistribution     map[string]int    `json:"type_distribution"`
	PriorityDistribution map[string]int    `json:"priority_distribution"`
	DeviceDistribution   map[string]int    `json:"device_distribution"`
	SuccessRate          float64           `json:"success_rate"`
	AverageExecutionTime float64           `json:"average_execution_time"`
	TrendData            []TaskTrendPoint  `json:"trend_data"`
	TopFailedTasks       []TaskFailureInfo `json:"top_failed_tasks"`
}

// TaskTrendPoint 任务趋势数据点
type TaskTrendPoint struct {
	Timestamp    *gtime.Time `json:"timestamp"`
	TaskCount    int         `json:"task_count"`
	SuccessCount int         `json:"success_count"`
	FailureCount int         `json:"failure_count"`
	SuccessRate  float64     `json:"success_rate"`
	AvgDuration  float64     `json:"avg_duration"`
}

// TaskFailureInfo 任务失败信息
type TaskFailureInfo struct {
	TaskId       string  `json:"task_id"`
	TaskName     string  `json:"task_name"`
	FailureCount int     `json:"failure_count"`
	FailureRate  float64 `json:"failure_rate"`
	LastError    string  `json:"last_error"`
}

// TaskPerformanceReport 任务性能报告
type TaskPerformanceReport struct {
	TaskId             string                 `json:"task_id"`
	TaskName           string                 `json:"task_name"`
	OverallScore       float64                `json:"overall_score"`
	ExecutionMetrics   TaskExecutionMetrics   `json:"execution_metrics"`
	ResourceMetrics    TaskResourceMetrics    `json:"resource_metrics"`
	ReliabilityMetrics TaskReliabilityMetrics `json:"reliability_metrics"`
	PerformanceTrend   []TaskTrendPoint       `json:"performance_trend"`
	Recommendations    []string               `json:"recommendations"`
}

// TaskExecutionMetrics 任务执行指标
type TaskExecutionMetrics struct {
	TotalExecutions      int     `json:"total_executions"`
	SuccessfulExecutions int     `json:"successful_executions"`
	FailedExecutions     int     `json:"failed_executions"`
	SuccessRate          float64 `json:"success_rate"`
	AverageExecutionTime float64 `json:"average_execution_time"`
	MinExecutionTime     float64 `json:"min_execution_time"`
	MaxExecutionTime     float64 `json:"max_execution_time"`
}

// TaskResourceMetrics 任务资源指标
type TaskResourceMetrics struct {
	AverageCpuUsage    float64 `json:"average_cpu_usage"`
	PeakCpuUsage       float64 `json:"peak_cpu_usage"`
	AverageMemoryUsage float64 `json:"average_memory_usage"`
	PeakMemoryUsage    float64 `json:"peak_memory_usage"`
	TotalDiskIO        float64 `json:"total_disk_io"`
	TotalNetworkIO     float64 `json:"total_network_io"`
}

// TaskReliabilityMetrics 任务可靠性指标
type TaskReliabilityMetrics struct {
	MTBF             float64 `json:"mtbf"` // Mean Time Between Failures
	MTTR             float64 `json:"mttr"` // Mean Time To Recovery
	AvailabilityRate float64 `json:"availability_rate"`
	ErrorRate        float64 `json:"error_rate"`
	RetrySuccessRate float64 `json:"retry_success_rate"`
}

/*
任务管理：
	1. 创建任务 ✅
	2. 获取任务列表 ✅
	3. 获取任务详情 ✅
	4. 更新任务 ✅
	5. 删除任务 ✅
	6. 取消任务 ✅
	7. 重试任务 ✅

任务执行控制：
	1. 启动任务 ✅
	2. 停止任务 ✅
	3. 重启任务 ✅
	4. 取消任务 ✅
	5. 重试任务 ✅

任务优先级管理：
	1. 获取任务优先级列表 ✅
	2. 更新任务优先级 ✅
	3. 批量更新任务优先级 ✅

任务状态管理：
	1. 获取任务状态 ✅
	2. 更新任务状态 ✅

任务日志管理：
	1. 获取任务日志列表 ✅
	2. 获取任务日志详情 ✅
	3. 清空任务日志 ✅
	4. 导出任务日志 ✅

任务队列管理：
	1. 获取任务队列列表 ✅
	2. 获取任务队列详情 ✅
	3. 管理任务队列 ✅
	4. 获取任务队列统计 ✅

任务调度管理：
	1. 获取任务调度列表 ✅
	2. 创建任务调度 ✅
	3. 更新任务调度 ✅
	4. 删除任务调度 ✅

任务统计分析：
	1. 获取任务统计信息 ✅
	2. 获取任务性能报告 ✅
*/
