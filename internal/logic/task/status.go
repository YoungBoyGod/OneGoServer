package task

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	task "OneGfServer/internal/model/task"
)

// ===============================
// 任务状态流转相关业务逻辑
// ===============================

// ValidateTaskStatusTransition 验证任务状态转换是否合法
func (s *sTask) ValidateTaskStatusTransition(ctx context.Context, input *task.ValidateTaskStatusTransitionInput) (*task.ValidateTaskStatusTransitionOutput, error) {
	// 定义合法的状态转换规则
	validTransitions := map[string][]string{
		"pending":   {"running", "cancelled", "deleted"},
		"running":   {"completed", "failed", "paused", "cancelled"},
		"paused":    {"running", "cancelled"},
		"completed": {"deleted"},
		"failed":    {"retry", "cancelled", "deleted"},
		"retry":     {"pending", "cancelled"},
		"cancelled": {"deleted"},
		"deleted":   {}, // 删除状态不能转换到其他状态
	}

	allowedStatuses, exists := validTransitions[input.CurrentStatus]
	if !exists {
		return &task.ValidateTaskStatusTransitionOutput{
			IsValid: false,
			Message: "无效的当前状态",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "无效的当前状态")
	}

	// 检查目标状态是否在允许的转换列表中
	for _, status := range allowedStatuses {
		if status == input.TargetStatus {
			return &task.ValidateTaskStatusTransitionOutput{
				IsValid: true,
				Message: "状态转换有效",
			}, nil
		}
	}

	return &task.ValidateTaskStatusTransitionOutput{
		IsValid: false,
		Message: "不允许的状态转换",
	}, gerror.NewCode(gcode.CodeInvalidParameter, "不允许的状态转换")
}

// DetermineTaskStatus 根据任务数据自动确定任务状态
func (s *sTask) DetermineTaskStatus(ctx context.Context, input *task.DetermineTaskStatusInput) (*task.DetermineTaskStatusOutput, error) {
	// 检查是否已完成
	if completedAt, ok := input.TaskData["completed_at"].(*gtime.Time); ok && completedAt != nil {
		return &task.DetermineTaskStatusOutput{
			Status: "completed",
		}, nil
	}

	// 检查是否失败
	if failedAt, ok := input.TaskData["failed_at"].(*gtime.Time); ok && failedAt != nil {
		return &task.DetermineTaskStatusOutput{
			Status: "failed",
		}, nil
	}

	// 检查是否已取消
	if cancelledAt, ok := input.TaskData["cancelled_at"].(*gtime.Time); ok && cancelledAt != nil {
		return &task.DetermineTaskStatusOutput{
			Status: "cancelled",
		}, nil
	}

	// 检查是否正在执行
	if startedAt, ok := input.TaskData["started_at"].(*gtime.Time); ok && startedAt != nil {
		// 检查是否暂停
		if pausedAt, ok := input.TaskData["paused_at"].(*gtime.Time); ok && pausedAt != nil {
			return &task.DetermineTaskStatusOutput{
				Status: "paused",
			}, nil
		}
		return &task.DetermineTaskStatusOutput{
			Status: "running",
		}, nil
	}

	// 检查是否已删除
	if deletedAt, ok := input.TaskData["deleted_at"].(*gtime.Time); ok && deletedAt != nil {
		return &task.DetermineTaskStatusOutput{
			Status: "deleted",
		}, nil
	}

	// 默认状态
	return &task.DetermineTaskStatusOutput{
		Status: "pending",
	}, nil
}

// CanTransitionToStatus 检查是否可以转换到指定状态
func (s *sTask) CanTransitionToStatus(ctx context.Context, input *task.CanTransitionToStatusInput) (*task.CanTransitionToStatusOutput, error) {
	determineInput := &task.DetermineTaskStatusInput{
		TaskData: input.TaskData,
	}
	determineOutput := s.DetermineTaskStatus(ctx, determineInput)

	validateInput := &task.ValidateTaskStatusTransitionInput{
		CurrentStatus: determineOutput.Status,
		TargetStatus:  input.TargetStatus,
	}
	validateOutput := s.ValidateTaskStatusTransition(ctx, validateInput)

	return &task.CanTransitionToStatusOutput{
		CanTransition: validateOutput.IsValid,
		Reason:        validateOutput.Message,
	}, nil
}

// ===============================
// 任务状态业务逻辑
// ===============================

// GetTaskStatus 获取任务状态
func (s *sTask) GetTaskStatus(ctx context.Context, input *task.GetTaskStatusInput) (*task.GetTaskStatusOutput, error) {
	// 这里应该从数据库获取任务状态
	// 目前返回模拟数据
	startTime := gtime.Now().Add(-30 * time.Minute)
	endTime := gtime.Now()
	duration := endTime.Sub(startTime.Time)

	return &task.GetTaskStatusOutput{
		TaskID:    input.TaskID,
		Status:    "running",
		Progress:  65.5,
		StartTime: startTime.Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Duration:  duration.String(),
		Details: map[string]interface{}{
			"current_step":        "数据处理",
			"total_steps":         10,
			"current_step_number": 7,
			"estimated_remaining": "15分钟",
			"memory_usage":        "512MB",
			"cpu_usage":           "45%",
		},
	}, nil
}
