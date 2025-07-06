package task

import "context"

// ===============================
// Task模块服务接口定义
// ===============================

// ITask 任务服务接口
type ITask interface {
	// 基础任务管理
	CreateTask(ctx context.Context, input *CreateTaskInput) (*CreateTaskOutput, error)
	GetTask(ctx context.Context, input *GetTaskInput) (*GetTaskOutput, error)
	UpdateTask(ctx context.Context, input *UpdateTaskInput) (*UpdateTaskOutput, error)
	DeleteTask(ctx context.Context, input *DeleteTaskInput) (*DeleteTaskOutput, error)
	ListTasks(ctx context.Context, input *ListTasksInput) (*ListTasksOutput, error)

	// 任务控制
	StartTask(ctx context.Context, input *StartTaskInput) (*StartTaskOutput, error)
	StopTask(ctx context.Context, input *StopTaskInput) (*StopTaskOutput, error)
	PauseTask(ctx context.Context, input *PauseTaskInput) (*PauseTaskOutput, error)
	ResumeTask(ctx context.Context, input *ResumeTaskInput) (*ResumeTaskOutput, error)

	// 任务分配
	AssignTask(ctx context.Context, input *AssignTaskInput) (*AssignTaskOutput, error)
	UnassignTask(ctx context.Context, input *UnassignTaskInput) (*UnassignTaskOutput, error)

	// 任务状态
	GetTaskStatus(ctx context.Context, input *GetTaskStatusInput) (*GetTaskStatusOutput, error)

	// 任务优先级
	UpdateTaskPriority(ctx context.Context, input *UpdateTaskPriorityInput) (*UpdateTaskPriorityOutput, error)

	// 任务统计
	GetTaskStatistics(ctx context.Context, input *GetTaskStatisticsInput) (*GetTaskStatisticsOutput, error)

	// 任务日志
	GetTaskLogs(ctx context.Context, input *GetTaskLogsInput) (*GetTaskLogsOutput, error)

	// 任务调度
	ScheduleTask(ctx context.Context, input *ScheduleTaskInput) (*ScheduleTaskOutput, error)
	CancelSchedule(ctx context.Context, input *CancelScheduleInput) (*CancelScheduleOutput, error)
}
