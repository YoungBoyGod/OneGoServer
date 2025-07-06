package task

// ===============================
// Task模块逻辑层Input/Output结构体定义
// ===============================

// CreateTaskInput 创建任务输入
type CreateTaskInput struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Timeout     int                    `json:"timeout"`
	RetryCount  int                    `json:"retry_count"`
}

// CreateTaskOutput 创建任务输出
type CreateTaskOutput struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message"`
}

// GetTaskInput 获取任务输入
type GetTaskInput struct {
	TaskID string `json:"task_id"`
}

// GetTaskOutput 获取任务输出
type GetTaskOutput struct {
	Task map[string]interface{} `json:"task"`
}

// UpdateTaskInput 更新任务输入
type UpdateTaskInput struct {
	TaskID      string                 `json:"task_id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Timeout     int                    `json:"timeout"`
	RetryCount  int                    `json:"retry_count"`
}

// UpdateTaskOutput 更新任务输出
type UpdateTaskOutput struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
	Message   string `json:"message"`
}

// DeleteTaskInput 删除任务输入
type DeleteTaskInput struct {
	TaskID string `json:"task_id"`
}

// DeleteTaskOutput 删除任务输出
type DeleteTaskOutput struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ListTasksInput 任务列表输入
type ListTasksInput struct {
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
}

// ListTasksOutput 任务列表输出
type ListTasksOutput struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
}

// StartTaskInput 启动任务输入
type StartTaskInput struct {
	TaskID string `json:"task_id"`
}

// StartTaskOutput 启动任务输出
type StartTaskOutput struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
	Message   string `json:"message"`
}

// StopTaskInput 停止任务输入
type StopTaskInput struct {
	TaskID string `json:"task_id"`
}

// StopTaskOutput 停止任务输出
type StopTaskOutput struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	StoppedAt string `json:"stopped_at"`
	Message   string `json:"message"`
}

// PauseTaskInput 暂停任务输入
type PauseTaskInput struct {
	TaskID string `json:"task_id"`
}

// PauseTaskOutput 暂停任务输出
type PauseTaskOutput struct {
	TaskID   string `json:"task_id"`
	Status   string `json:"status"`
	PausedAt string `json:"paused_at"`
	Message  string `json:"message"`
}

// ResumeTaskInput 恢复任务输入
type ResumeTaskInput struct {
	TaskID string `json:"task_id"`
}

// ResumeTaskOutput 恢复任务输出
type ResumeTaskOutput struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	ResumedAt string `json:"resumed_at"`
	Message   string `json:"message"`
}

// AssignTaskInput 分配任务输入
type AssignTaskInput struct {
	TaskID   string `json:"task_id"`
	DeviceID string `json:"device_id"`
}

// AssignTaskOutput 分配任务输出
type AssignTaskOutput struct {
	TaskID     string `json:"task_id"`
	DeviceID   string `json:"device_id"`
	Status     string `json:"status"`
	AssignedAt string `json:"assigned_at"`
	Message    string `json:"message"`
}

// UnassignTaskInput 取消分配任务输入
type UnassignTaskInput struct {
	TaskID string `json:"task_id"`
}

// UnassignTaskOutput 取消分配任务输出
type UnassignTaskOutput struct {
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	UnassignedAt string `json:"unassigned_at"`
	Message      string `json:"message"`
}

// GetTaskStatusInput 获取任务状态输入
type GetTaskStatusInput struct {
	TaskID string `json:"task_id"`
}

// GetTaskStatusOutput 获取任务状态输出
type GetTaskStatusOutput struct {
	TaskID    string                 `json:"task_id"`
	Status    string                 `json:"status"`
	Progress  float64                `json:"progress"`
	StartTime string                 `json:"start_time"`
	EndTime   string                 `json:"end_time"`
	Duration  string                 `json:"duration"`
	Details   map[string]interface{} `json:"details"`
}

// UpdateTaskPriorityInput 更新任务优先级输入
type UpdateTaskPriorityInput struct {
	TaskID   string `json:"task_id"`
	Priority int    `json:"priority"`
}

// UpdateTaskPriorityOutput 更新任务优先级输出
type UpdateTaskPriorityOutput struct {
	TaskID    string `json:"task_id"`
	Priority  int    `json:"priority"`
	UpdatedAt string `json:"updated_at"`
	Message   string `json:"message"`
}

// GetTaskStatisticsInput 获取任务统计输入
type GetTaskStatisticsInput struct {
	Period string `json:"period"`
}

// GetTaskStatisticsOutput 获取任务统计输出
type GetTaskStatisticsOutput struct {
	Period     string                 `json:"period"`
	Statistics map[string]interface{} `json:"statistics"`
}

// GetTaskLogsInput 获取任务日志输入
type GetTaskLogsInput struct {
	TaskID    string `json:"task_id"`
	Level     string `json:"level"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

// GetTaskLogsOutput 获取任务日志输出
type GetTaskLogsOutput struct {
	TaskID string                   `json:"task_id"`
	List   []map[string]interface{} `json:"list"`
	Total  int                      `json:"total"`
	Page   int                      `json:"page"`
	Size   int                      `json:"size"`
}

// ScheduleTaskInput 调度任务输入
type ScheduleTaskInput struct {
	TaskID     string `json:"task_id"`
	ScheduleAt string `json:"schedule_at"`
	Recurring  bool   `json:"recurring"`
	CronExpr   string `json:"cron_expr"`
}

// ScheduleTaskOutput 调度任务输出
type ScheduleTaskOutput struct {
	TaskID      string `json:"task_id"`
	ScheduleID  string `json:"schedule_id"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduled_at"`
	Message     string `json:"message"`
}

// CancelScheduleInput 取消调度输入
type CancelScheduleInput struct {
	TaskID     string `json:"task_id"`
	ScheduleID string `json:"schedule_id"`
}

// CancelScheduleOutput 取消调度输出
type CancelScheduleOutput struct {
	TaskID      string `json:"task_id"`
	ScheduleID  string `json:"schedule_id"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelled_at"`
	Message     string `json:"message"`
}

// 辅助方法相关结构体

// ValidateTaskInput 验证任务输入
type ValidateTaskInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// ValidateTaskOutput 验证任务输出
type ValidateTaskOutput struct {
	IsValid bool     `json:"is_valid"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// CalculateTaskPriorityInput 计算任务优先级输入
type CalculateTaskPriorityInput struct {
	TaskType     string                 `json:"task_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	UserPriority int                    `json:"user_priority"`
}

// CalculateTaskPriorityOutput 计算任务优先级输出
type CalculateTaskPriorityOutput struct {
	Priority     int    `json:"priority"`
	Reason       string `json:"reason"`
	CalculatedAt string `json:"calculated_at"`
}

// CheckTaskDependenciesInput 检查任务依赖输入
type CheckTaskDependenciesInput struct {
	TaskID string `json:"task_id"`
}

// CheckTaskDependenciesOutput 检查任务依赖输出
type CheckTaskDependenciesOutput struct {
	TaskID       string   `json:"task_id"`
	Dependencies []string `json:"dependencies"`
	IsReady      bool     `json:"is_ready"`
	BlockedBy    []string `json:"blocked_by"`
}

// GetTaskExecutionHistoryInput 获取任务执行历史输入
type GetTaskExecutionHistoryInput struct {
	TaskID string `json:"task_id"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

// GetTaskExecutionHistoryOutput 获取任务执行历史输出
type GetTaskExecutionHistoryOutput struct {
	TaskID string                   `json:"task_id"`
	List   []map[string]interface{} `json:"list"`
	Total  int                      `json:"total"`
	Page   int                      `json:"page"`
	Size   int                      `json:"size"`
}

// GenerateTaskReportInput 生成任务报告输入
type GenerateTaskReportInput struct {
	TaskID string `json:"task_id"`
	Format string `json:"format"`
}

// GenerateTaskReportOutput 生成任务报告输出
type GenerateTaskReportOutput struct {
	TaskID   string `json:"task_id"`
	ReportID string `json:"report_id"`
	Format   string `json:"format"`
	URL      string `json:"url"`
	Message  string `json:"message"`
}
