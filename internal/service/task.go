// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	task "OneGfServer/internal/model/task"
	"context"
)

type (
	ITask interface {
		// AssignTaskToDevice 分配任务到设备
		AssignTaskToDevice(ctx context.Context, input *task.AssignTaskToDeviceInput) (*task.AssignTaskToDeviceOutput, error)
		// AssignTask 分配任务
		AssignTask(ctx context.Context, input *task.AssignTaskInput) (*task.AssignTaskOutput, error)
		// UnassignTask 取消分配任务
		UnassignTask(ctx context.Context, input *task.UnassignTaskInput) (*task.UnassignTaskOutput, error)
		// CreateTask 创建任务
		CreateTask(ctx context.Context, input *task.CreateTaskInput) (*task.CreateTaskOutput, error)
		// GetTask 获取任务
		GetTask(ctx context.Context, input *task.GetTaskInput) (*task.GetTaskOutput, error)
		// UpdateTask 更新任务
		UpdateTask(ctx context.Context, input *task.UpdateTaskInput) (*task.UpdateTaskOutput, error)
		// DeleteTask 删除任务
		DeleteTask(ctx context.Context, input *task.DeleteTaskInput) (*task.DeleteTaskOutput, error)
		// ListTasks 任务列表
		ListTasks(ctx context.Context, input *task.ListTasksInput) (*task.ListTasksOutput, error)
		// StartTask 启动任务
		StartTask(ctx context.Context, input *task.StartTaskInput) (*task.StartTaskOutput, error)
		// StopTask 停止任务
		StopTask(ctx context.Context, input *task.StopTaskInput) (*task.StopTaskOutput, error)
		// PauseTask 暂停任务
		PauseTask(ctx context.Context, input *task.PauseTaskInput) (*task.PauseTaskOutput, error)
		// ResumeTask 恢复任务
		ResumeTask(ctx context.Context, input *task.ResumeTaskInput) (*task.ResumeTaskOutput, error)
		// GetTaskLogs 获取任务日志
		GetTaskLogs(ctx context.Context, input *task.GetTaskLogsInput) (*task.GetTaskLogsOutput, error)
		// CalculateTaskPriority 计算任务智能优先级
		CalculateTaskPriority(ctx context.Context, input *task.CalculateTaskPriorityInput) (*task.CalculateTaskPriorityOutput, error)
		// UpdateTaskPriority 更新任务优先级
		UpdateTaskPriority(ctx context.Context, input *task.UpdateTaskPriorityInput) (*task.UpdateTaskPriorityOutput, error)
		// ScheduleTask 调度任务
		ScheduleTask(ctx context.Context, input *task.ScheduleTaskInput) (*task.ScheduleTaskOutput, error)
		// CancelSchedule 取消调度
		CancelSchedule(ctx context.Context, input *task.CancelScheduleInput) (*task.CancelScheduleOutput, error)
		// GetTaskStatistics 获取任务统计
		GetTaskStatistics(ctx context.Context, input *task.GetTaskStatisticsInput) (*task.GetTaskStatisticsOutput, error)
		// ValidateTaskStatusTransition 验证任务状态转换是否合法
		ValidateTaskStatusTransition(ctx context.Context, input *task.ValidateTaskStatusTransitionInput) (*task.ValidateTaskStatusTransitionOutput, error)
		// DetermineTaskStatus 根据任务数据自动确定任务状态
		DetermineTaskStatus(ctx context.Context, input *task.DetermineTaskStatusInput) (*task.DetermineTaskStatusOutput, error)
		// CanTransitionToStatus 检查是否可以转换到指定状态
		CanTransitionToStatus(ctx context.Context, input *task.CanTransitionToStatusInput) (*task.CanTransitionToStatusOutput, error)
		// GetTaskStatus 获取任务状态
		GetTaskStatus(ctx context.Context, input *task.GetTaskStatusInput) (*task.GetTaskStatusOutput, error)
	}
)

var (
	localTask ITask
)

func Task() ITask {
	if localTask == nil {
		panic("implement not found for interface ITask, forgot register?")
	}
	return localTask
}

func RegisterTask(i ITask) {
	localTask = i
}
