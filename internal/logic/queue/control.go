package queue

import (
	"context"
	"math"

	queue "OneGfServer/internal/model/queue"
)

// ===============================
// 队列操作控制业务逻辑
// ===============================

// ValidateQueueOperation 验证队列操作
func (s *sQueue) ValidateQueueOperation(ctx context.Context, input *queue.ValidateQueueOperationInput) (*queue.ValidateQueueOperationOutput, error) {
	var errors []string

	// 检查队列状态
	if status, ok := input.QueueData["status"].(string); ok {
		switch input.Operation {
		case "start":
			if status == "running" {
				errors = append(errors, "队列已在运行状态")
			}
		case "stop":
			if status == "stopped" {
				errors = append(errors, "队列已在停止状态")
			}
		case "pause":
			if status == "paused" {
				errors = append(errors, "队列已在暂停状态")
			}
		case "resume":
			if status != "paused" {
				errors = append(errors, "只有暂停状态的队列才能恢复")
			}
		}
	}

	// 检查队列健康状态
	if healthScore, ok := input.QueueData["health_score"].(float64); ok {
		if healthScore < 30 && input.Operation == "start" {
			errors = append(errors, "队列健康度过低，无法启动")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &queue.ValidateQueueOperationOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}, nil
}

// CalculateQueueHealthScore 计算队列健康度评分
func (s *sQueue) CalculateQueueHealthScore(ctx context.Context, input *queue.CalculateQueueHealthScoreInput) (*queue.CalculateQueueHealthScoreOutput, error) {
	score := 100.0
	components := make(map[string]interface{})

	// 1. 队列状态评分 (权重: 30%)
	statusScoreInput := &queue.CalculateStatusScoreInput{
		QueueData: input.QueueData,
	}
	statusScoreOutput := s.calculateStatusScore(statusScoreInput)
	statusScore := statusScoreOutput.Score
	score += statusScore * 0.30
	components["status_score"] = statusScore

	// 2. 性能指标评分 (权重: 25%)
	performanceScoreInput := &queue.CalculatePerformanceScoreInput{
		QueueData: input.QueueData,
	}
	performanceScoreOutput := s.calculatePerformanceScore(performanceScoreInput)
	performanceScore := performanceScoreOutput.Score
	score += performanceScore * 0.25
	components["performance_score"] = performanceScore

	// 3. 错误率评分 (权重: 20%)
	errorScoreInput := &queue.CalculateErrorScoreInput{
		QueueData: input.QueueData,
	}
	errorScoreOutput := s.calculateErrorScore(errorScoreInput)
	errorScore := errorScoreOutput.Score
	score += errorScore * 0.20
	components["error_score"] = errorScore

	// 4. 资源使用评分 (权重: 15%)
	resourceScoreInput := &queue.CalculateResourceScoreInput{
		QueueData: input.QueueData,
	}
	resourceScoreOutput := s.calculateResourceScore(resourceScoreInput)
	resourceScore := resourceScoreOutput.Score
	score += resourceScore * 0.15
	components["resource_score"] = resourceScore

	// 5. 响应时间评分 (权重: 10%)
	responseScoreInput := &queue.CalculateResponseScoreInput{
		QueueData: input.QueueData,
	}
	responseScoreOutput := s.calculateResponseScore(responseScoreInput)
	responseScore := responseScoreOutput.Score
	score += responseScore * 0.10
	components["response_score"] = responseScore

	finalScore := math.Max(0, math.Min(score, 100))

	return &queue.CalculateQueueHealthScoreOutput{
		HealthScore: finalScore,
		Components:  components,
	}, nil
}

// calculateStatusScore 计算状态评分
func (s *sQueue) calculateStatusScore(input *queue.CalculateStatusScoreInput) *queue.CalculateStatusScoreOutput {
	status, _ := input.QueueData["status"].(string)
	switch status {
	case "running":
		return &queue.CalculateStatusScoreOutput{Score: 100}
	case "paused":
		return &queue.CalculateStatusScoreOutput{Score: 70}
	case "stopped":
		return &queue.CalculateStatusScoreOutput{Score: 50}
	case "error":
		return &queue.CalculateStatusScoreOutput{Score: 20}
	default:
		return &queue.CalculateStatusScoreOutput{Score: 30}
	}
}

// calculatePerformanceScore 计算性能评分
func (s *sQueue) calculatePerformanceScore(input *queue.CalculatePerformanceScoreInput) *queue.CalculatePerformanceScoreOutput {
	score := 100.0

	// 检查处理速率
	if rate, ok := input.QueueData["processing_rate"].(float64); ok {
		if rate < 10 {
			score -= 30
		} else if rate < 50 {
			score -= 15
		}
	}

	// 检查队列长度
	if length, ok := input.QueueData["current_length"].(int); ok {
		if capacity, ok := input.QueueData["capacity"].(int); ok && capacity > 0 {
			utilization := float64(length) / float64(capacity)
			if utilization > 0.9 {
				score -= 40
			} else if utilization > 0.7 {
				score -= 20
			}
		}
	}

	return &queue.CalculatePerformanceScoreOutput{Score: math.Max(0, score)}
}

// calculateErrorScore 计算错误率评分
func (s *sQueue) calculateErrorScore(input *queue.CalculateErrorScoreInput) *queue.CalculateErrorScoreOutput {
	score := 100.0

	// 检查错误率
	if errorRate, ok := input.QueueData["error_rate"].(float64); ok {
		if errorRate > 0.1 {
			score -= 40
		} else if errorRate > 0.05 {
			score -= 20
		} else if errorRate > 0.01 {
			score -= 10
		}
	}

	// 检查失败任务数
	if failedTasks, ok := input.QueueData["failed_tasks"].(int); ok {
		if totalTasks, ok := input.QueueData["total_tasks"].(int); ok && totalTasks > 0 {
			failureRate := float64(failedTasks) / float64(totalTasks)
			if failureRate > 0.2 {
				score -= 30
			} else if failureRate > 0.1 {
				score -= 15
			}
		}
	}

	return &queue.CalculateErrorScoreOutput{Score: math.Max(0, score)}
}

// calculateResourceScore 计算资源使用评分
func (s *sQueue) calculateResourceScore(input *queue.CalculateResourceScoreInput) *queue.CalculateResourceScoreOutput {
	score := 100.0

	// 检查CPU使用率
	if cpuUsage, ok := input.QueueData["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			score -= 30
		} else if cpuUsage > 70 {
			score -= 15
		}
	}

	// 检查内存使用率
	if memoryUsage, ok := input.QueueData["memory_usage"].(float64); ok {
		if memoryUsage > 90 {
			score -= 25
		} else if memoryUsage > 80 {
			score -= 12
		}
	}

	// 检查磁盘使用率
	if diskUsage, ok := input.QueueData["disk_usage"].(float64); ok {
		if diskUsage > 95 {
			score -= 20
		} else if diskUsage > 85 {
			score -= 10
		}
	}

	return &queue.CalculateResourceScoreOutput{Score: math.Max(0, score)}
}

// calculateResponseScore 计算响应时间评分
func (s *sQueue) calculateResponseScore(input *queue.CalculateResponseScoreInput) *queue.CalculateResponseScoreOutput {
	score := 100.0

	// 检查平均响应时间
	if avgResponseTime, ok := input.QueueData["avg_response_time"].(float64); ok {
		if avgResponseTime > 5000 {
			score -= 40
		} else if avgResponseTime > 2000 {
			score -= 20
		} else if avgResponseTime > 1000 {
			score -= 10
		}
	}

	// 检查最大响应时间
	if maxResponseTime, ok := input.QueueData["max_response_time"].(float64); ok {
		if maxResponseTime > 30000 {
			score -= 30
		} else if maxResponseTime > 15000 {
			score -= 15
		}
	}

	return &queue.CalculateResponseScoreOutput{Score: math.Max(0, score)}
}
