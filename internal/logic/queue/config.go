package queue

import (
	"context"

	queue "OneGfServer/internal/model/queue"
)

// ===============================
// 队列配置管理业务逻辑
// ===============================

// ValidateQueueConfiguration 验证队列配置
func (s *sQueue) ValidateQueueConfiguration(ctx context.Context, input *queue.ValidateQueueConfigurationInput) (*queue.ValidateQueueConfigurationOutput, error) {
	var errors []string

	// 验证基础配置
	basicConfigInput := &queue.ValidateBasicConfigInput{
		Config: input.Config,
	}
	basicConfigOutput := s.validateBasicConfig(basicConfigInput)
	if !basicConfigOutput.IsValid {
		errors = append(errors, basicConfigOutput.Errors...)
	}

	// 验证性能配置
	performanceConfigInput := &queue.ValidatePerformanceConfigInput{
		Config: input.Config,
	}
	performanceConfigOutput := s.validatePerformanceConfig(performanceConfigInput)
	if !performanceConfigOutput.IsValid {
		errors = append(errors, performanceConfigOutput.Errors...)
	}

	// 验证安全配置
	securityConfigInput := &queue.ValidateSecurityConfigInput{
		Config: input.Config,
	}
	securityConfigOutput := s.validateSecurityConfig(securityConfigInput)
	if !securityConfigOutput.IsValid {
		errors = append(errors, securityConfigOutput.Errors...)
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateQueueConfigurationOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}, nil
}

// validateBasicConfig 验证基础配置
func (s *sQueue) validateBasicConfig(input *queue.ValidateBasicConfigInput) *queue.ValidateBasicConfigOutput {
	var errors []string

	// 检查队列名称
	if name, ok := input.Config["name"].(string); ok {
		if name == "" {
			errors = append(errors, "队列名称不能为空")
		}
	} else {
		errors = append(errors, "队列名称不能为空")
	}

	// 检查队列类型
	if queueType, ok := input.Config["type"].(string); ok {
		validTypes := []string{"fifo", "lifo", "priority", "round_robin", "weighted"}
		found := false
		for _, validType := range validTypes {
			if queueType == validType {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, "不支持的队列类型")
		}
	} else {
		errors = append(errors, "队列类型不能为空")
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateBasicConfigOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}

// validatePerformanceConfig 验证性能配置
func (s *sQueue) validatePerformanceConfig(input *queue.ValidatePerformanceConfigInput) *queue.ValidatePerformanceConfigOutput {
	var errors []string

	// 检查容量配置
	if capacity, ok := input.Config["capacity"].(int); ok {
		if capacity <= 0 || capacity > 10000 {
			errors = append(errors, "队列容量必须在1-10000之间")
		}
	}

	// 检查并发配置
	if concurrency, ok := input.Config["max_concurrency"].(int); ok {
		if concurrency <= 0 || concurrency > 100 {
			errors = append(errors, "最大并发数必须在1-100之间")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidatePerformanceConfigOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}

// validateSecurityConfig 验证安全配置
func (s *sQueue) validateSecurityConfig(input *queue.ValidateSecurityConfigInput) *queue.ValidateSecurityConfigOutput {
	var errors []string

	// 检查访问控制配置
	if accessControl, ok := input.Config["access_control"].(map[string]interface{}); ok {
		accessControlInput := &queue.ValidateAccessControlInput{
			AccessControl: accessControl,
		}
		accessControlOutput := s.validateAccessControl(accessControlInput)
		if !accessControlOutput.IsValid {
			errors = append(errors, accessControlOutput.Errors...)
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateSecurityConfigOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}

// validateAccessControl 验证访问控制
func (s *sQueue) validateAccessControl(input *queue.ValidateAccessControlInput) *queue.ValidateAccessControlOutput {
	var errors []string

	// 检查权限配置
	if permissions, ok := input.AccessControl["permissions"].([]string); ok {
		if len(permissions) == 0 {
			errors = append(errors, "访问权限不能为空")
		}
	} else {
		errors = append(errors, "访问权限配置无效")
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateAccessControlOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}
