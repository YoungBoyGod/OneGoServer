package queue

import (
	"context"

	queue "OneGfServer/internal/model/queue"
)

// ===============================
// 队列任务管理业务逻辑
// ===============================

// ValidateTaskEnqueue 验证任务入队
func (s *sQueue) ValidateTaskEnqueue(ctx context.Context, input *queue.ValidateTaskEnqueueInput) (*queue.ValidateTaskEnqueueOutput, error) {
	var errors []string

	// 验证队列数据
	if input.QueueData == nil {
		errors = append(errors, "队列数据不能为空")
		return &queue.ValidateTaskEnqueueOutput{
			IsValid: false,
			Message: "验证失败",
			Errors:  errors,
		}, nil
	}

	// 验证任务数据
	taskValidationInput := &queue.ValidateTaskDataInput{
		TaskData: input.TaskData,
	}
	taskValidationOutput := s.validateTaskData(taskValidationInput)
	if !taskValidationOutput.IsValid {
		errors = append(errors, taskValidationOutput.Errors...)
	}

	// 检查队列容量
	if currentLength, ok := input.QueueData["current_length"].(int); ok {
		if capacity, ok := input.QueueData["capacity"].(int); ok {
			if currentLength >= capacity {
				errors = append(errors, "队列已满，无法添加新任务")
			}
		}
	}

	// 检查队列状态
	if status, ok := input.QueueData["status"].(string); ok {
		if status != "running" && status != "paused" {
			errors = append(errors, "队列状态不允许添加任务")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateTaskEnqueueOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}, nil
}

// validateTaskData 验证任务数据
func (s *sQueue) validateTaskData(input *queue.ValidateTaskDataInput) *queue.ValidateTaskDataOutput {
	var errors []string

	// 检查任务ID
	if taskId, ok := input.TaskData["task_id"].(string); ok {
		if taskId == "" {
			errors = append(errors, "任务ID不能为空")
		}
	} else {
		errors = append(errors, "任务ID不能为空")
	}

	// 检查任务类型
	if taskType, ok := input.TaskData["type"].(string); ok {
		if taskType == "" {
			errors = append(errors, "任务类型不能为空")
		}
	} else {
		errors = append(errors, "任务类型不能为空")
	}

	// 检查超时时间
	if timeout, ok := input.TaskData["timeout"].(int); ok {
		if timeout <= 0 || timeout > 3600 {
			errors = append(errors, "任务超时时间必须在1-3600秒之间")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateTaskDataOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}
