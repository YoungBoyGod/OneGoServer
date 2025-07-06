package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	task "OneGfServer/internal/model/task"
)

// ===============================
// 任务调度业务逻辑
// ===============================

// ScheduleTask 调度任务
func (s *sTask) ScheduleTask(ctx context.Context, input *task.ScheduleTaskInput) (*task.ScheduleTaskOutput, error) {
	// 验证调度时间
	if input.ScheduleAt == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "调度时间不能为空")
	}

	// 检查任务是否存在
	taskStatusInput := &task.GetTaskStatusInput{TaskID: input.TaskID}
	taskStatusOutput := s.getTaskStatus(taskStatusInput)
	if taskStatusOutput.Status == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound, "任务不存在")
	}

	// 检查是否已有调度
	if s.isTaskScheduled(input.TaskID) {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "任务已有调度")
	}

	// 生成调度ID
	scheduleID := s.generateScheduleID()

	// 这里应该保存调度信息到数据库
	// 目前返回模拟结果

	return &task.ScheduleTaskOutput{
		TaskID:      input.TaskID,
		ScheduleID:  scheduleID,
		Status:      "scheduled",
		ScheduledAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:     "任务调度成功",
	}, nil
}

// CancelSchedule 取消调度
func (s *sTask) CancelSchedule(ctx context.Context, input *task.CancelScheduleInput) (*task.CancelScheduleOutput, error) {
	// 检查调度是否存在
	if !s.isScheduleExists(input.ScheduleID) {
		return nil, gerror.NewCode(gcode.CodeNotFound, "调度不存在")
	}

	// 这里应该更新数据库中的调度状态
	// 目前返回模拟结果

	return &task.CancelScheduleOutput{
		TaskID:      input.TaskID,
		ScheduleID:  input.ScheduleID,
		Status:      "cancelled",
		CancelledAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:     "调度取消成功",
	}, nil
}

// isTaskScheduled 检查任务是否已有调度（内部方法）
func (s *sTask) isTaskScheduled(taskID string) bool {
	// 这里应该检查任务调度状态
	// 目前返回模拟结果
	return false
}

// isScheduleExists 检查调度是否存在（内部方法）
func (s *sTask) isScheduleExists(scheduleID string) bool {
	// 这里应该检查调度是否存在
	// 目前返回模拟结果
	return true
}

// generateScheduleID 生成调度ID（内部方法）
func (s *sTask) generateScheduleID() string {
	return "schedule-" + gtime.Now().Format("20060102150405") + "-" + string(rune(gtime.Now().UnixNano()%1000))
}

// getTaskStatus 获取任务状态（内部方法）
func (s *sTask) getTaskStatus(input *task.GetTaskStatusInput) *task.GetTaskStatusOutput {
	// 这里应该从数据库获取任务状态
	// 目前返回模拟数据
	return &task.GetTaskStatusOutput{
		TaskID:    input.TaskID,
		Status:    "pending",
		Progress:  0.0,
		StartTime: "",
		EndTime:   "",
		Duration:  "",
		Details:   map[string]interface{}{},
	}
}
