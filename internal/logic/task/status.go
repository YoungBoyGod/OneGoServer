package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务状态流转相关业务逻辑
// ===============================

// ValidateTaskStatusTransition 验证任务状态转换是否合法
func (s *sTask) ValidateTaskStatusTransition(ctx context.Context, currentStatus, targetStatus string) error {
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

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return gerror.NewCode(gcode.CodeInvalidParameter, "无效的当前状态")
	}

	// 检查目标状态是否在允许的转换列表中
	for _, status := range allowedStatuses {
		if status == targetStatus {
			return nil
		}
	}

	return gerror.NewCode(gcode.CodeInvalidParameter,
		"不允许的状态转换")
}

// DetermineTaskStatus 根据任务数据自动确定任务状态
func (s *sTask) DetermineTaskStatus(ctx context.Context, taskData map[string]interface{}) string {
	// 检查是否已完成
	if completedAt, ok := taskData["completed_at"].(*gtime.Time); ok && completedAt != nil {
		return "completed"
	}

	// 检查是否失败
	if failedAt, ok := taskData["failed_at"].(*gtime.Time); ok && failedAt != nil {
		return "failed"
	}

	// 检查是否已取消
	if cancelledAt, ok := taskData["cancelled_at"].(*gtime.Time); ok && cancelledAt != nil {
		return "cancelled"
	}

	// 检查是否正在执行
	if startedAt, ok := taskData["started_at"].(*gtime.Time); ok && startedAt != nil {
		// 检查是否暂停
		if pausedAt, ok := taskData["paused_at"].(*gtime.Time); ok && pausedAt != nil {
			return "paused"
		}
		return "running"
	}

	// 检查是否已删除
	if deletedAt, ok := taskData["deleted_at"].(*gtime.Time); ok && deletedAt != nil {
		return "deleted"
	}

	// 默认状态
	return "pending"
}

// CanTransitionToStatus 检查是否可以转换到指定状态
func (s *sTask) CanTransitionToStatus(ctx context.Context, taskData map[string]interface{}, targetStatus string) bool {
	currentStatus := s.DetermineTaskStatus(ctx, taskData)
	return s.ValidateTaskStatusTransition(ctx, currentStatus, targetStatus) == nil
}
