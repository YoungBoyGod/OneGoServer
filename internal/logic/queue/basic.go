package queue

import (
	"context"
	"regexp"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	queue "OneGfServer/internal/model/queue"
)

// ===============================
// 队列基础管理业务逻辑
// ===============================

// ValidateQueueCreation 验证队列创建
func (s *sQueue) ValidateQueueCreation(ctx context.Context, input *queue.ValidateQueueCreationInput) (*queue.ValidateQueueCreationOutput, error) {
	var errors []string

	// 验证队列名称
	if name, ok := input.QueueData["name"].(string); ok {
		if err := s.validateQueueName(name); err != nil {
			errors = append(errors, err.Error())
		}
	} else {
		errors = append(errors, "队列名称不能为空")
	}

	// 验证队列类型
	if queueType, ok := input.QueueData["type"].(string); ok {
		if err := s.validateQueueType(queueType); err != nil {
			errors = append(errors, err.Error())
		}
	} else {
		errors = append(errors, "队列类型不能为空")
	}

	// 验证容量配置
	if capacity, ok := input.QueueData["capacity"].(int); ok {
		if capacity <= 0 || capacity > 10000 {
			errors = append(errors, "队列容量必须在1-10000之间")
		}
	}

	// 验证优先级配置
	if maxPriority, ok := input.QueueData["max_priority"].(int); ok {
		if maxPriority <= 0 || maxPriority > 100 {
			errors = append(errors, "最大优先级必须在1-100之间")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateQueueCreationOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}, nil
}

// validateQueueName 验证队列名称
func (s *sQueue) validateQueueName(name string) error {
	if name == "" {
		return gerror.NewCode(gcode.CodeValidationFailed, "队列名称不能为空")
	}

	if len(name) > 50 {
		return gerror.NewCode(gcode.CodeValidationFailed, "队列名称长度不能超过50个字符")
	}

	// 检查名称格式
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return gerror.NewCode(gcode.CodeValidationFailed, "队列名称只能包含字母、数字、下划线和连字符")
	}

	return nil
}

// validateQueueType 验证队列类型
func (s *sQueue) validateQueueType(queueType string) error {
	validTypes := []string{"fifo", "lifo", "priority", "round_robin", "weighted"}
	for _, validType := range validTypes {
		if queueType == validType {
			return nil
		}
	}
	return gerror.NewCode(gcode.CodeValidationFailed, "不支持的队列类型")
}
