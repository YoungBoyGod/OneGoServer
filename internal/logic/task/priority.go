package task

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	task "OneGfServer/internal/model/task"
)

// ===============================
// 任务优先级计算相关业务逻辑
// ===============================

// CalculateTaskPriority 计算任务智能优先级
func (s *sTask) CalculateTaskPriority(ctx context.Context, input *task.CalculateTaskPriorityInput) (*task.CalculateTaskPriorityOutput, error) {
	// 计算各个维度的评分
	urgencyInput := &task.CalculateUrgencyScoreInput{
		TaskData: input.TaskData,
	}
	urgencyOutput := s.calculateUrgencyScore(urgencyInput)
	urgencyScore := urgencyOutput.Score

	businessInput := &task.CalculateBusinessImportanceScoreInput{
		TaskData: input.TaskData,
	}
	businessOutput := s.calculateBusinessImportanceScore(businessInput)
	businessScore := businessOutput.Score

	deadlineInput := &task.CalculateDeadlineScoreInput{
		TaskData: input.TaskData,
	}
	deadlineOutput := s.calculateDeadlineScore(deadlineInput)
	deadlineScore := deadlineOutput.Score

	resourceInput := &task.CalculateResourceScoreInput{
		TaskData: input.TaskData,
	}
	resourceOutput := s.calculateResourceScore(resourceInput)
	resourceScore := resourceOutput.Score

	// 综合计算优先级 (1-10)
	totalScore := urgencyScore*0.3 + businessScore*0.3 + deadlineScore*0.25 + resourceScore*0.15
	priority := int(math.Round(totalScore))

	// 确保优先级在1-10范围内
	if priority < 1 {
		priority = 1
	} else if priority > 10 {
		priority = 10
	}

	return &task.CalculateTaskPriorityOutput{
		Priority: priority,
	}, nil
}

// calculateUrgencyScore 计算紧急程度评分
func (s *sTask) calculateUrgencyScore(input *task.CalculateUrgencyScoreInput) *task.CalculateUrgencyScoreOutput {
	score := 0.0

	// 检查任务类型的紧急性
	if taskType, ok := input.TaskData["task_type"].(string); ok {
		switch taskType {
		case "emergency":
			score = 100
		case "urgent":
			score = 80
		case "normal":
			score = 50
		case "low":
			score = 20
		}
	}

	// 检查是否有紧急标签
	if tags, ok := input.TaskData["tags"].([]string); ok {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), "urgent") ||
				strings.Contains(strings.ToLower(tag), "emergency") {
				score = math.Max(score, 90)
			}
		}
	}

	// 检查是否有紧急标记
	if urgent, ok := input.TaskData["is_urgent"].(bool); ok && urgent {
		score += 30
	}

	// 检查创建时间
	if createdAt, ok := input.TaskData["created_at"].(*gtime.Time); ok && createdAt != nil {
		hoursSinceCreation := time.Since(createdAt.Time).Hours()
		if hoursSinceCreation > 24 {
			score += 20
		} else if hoursSinceCreation > 12 {
			score += 15
		} else if hoursSinceCreation > 6 {
			score += 10
		}
	}

	// 检查依赖任务状态
	if dependencies, ok := input.TaskData["dependencies"].([]string); ok {
		if len(dependencies) == 0 {
			score += 10 // 无依赖的任务优先级更高
		}
	}

	return &task.CalculateUrgencyScoreOutput{
		Score: math.Min(score, 100),
	}
}

// calculateBusinessImportanceScore 计算业务重要性评分
func (s *sTask) calculateBusinessImportanceScore(input *task.CalculateBusinessImportanceScoreInput) *task.CalculateBusinessImportanceScoreOutput {
	score := 0.0

	// 检查任务类别
	if category, ok := input.TaskData["category"].(string); ok {
		switch category {
		case "critical":
			score = 100
		case "important":
			score = 80
		case "normal":
			score = 50
		case "minor":
			score = 30
		}
	}

	// 检查业务类型
	if businessType, ok := input.TaskData["business_type"].(string); ok {
		switch businessType {
		case "critical":
			score += 40
		case "high":
			score += 30
		case "medium":
			score += 20
		case "low":
			score += 10
		}
	}

	// 检查影响范围
	if impact, ok := input.TaskData["impact_scope"].(string); ok {
		switch impact {
		case "global":
			score += 20
		case "regional":
			score += 10
		case "local":
			score += 5
		}
	}

	// 检查影响范围
	if impactScope, ok := input.TaskData["impact_scope"].(string); ok {
		switch impactScope {
		case "system":
			score += 30
		case "department":
			score += 20
		case "team":
			score += 15
		case "individual":
			score += 10
		}
	}

	// 检查客户影响
	if customerImpact, ok := input.TaskData["customer_impact"].(bool); ok && customerImpact {
		score += 20
	}

	// 检查收益影响
	if revenueImpact, ok := input.TaskData["revenue_impact"].(float64); ok {
		if revenueImpact > 1000000 {
			score += 20
		} else if revenueImpact > 100000 {
			score += 15
		} else if revenueImpact > 10000 {
			score += 10
		}
	}

	return &task.CalculateBusinessImportanceScoreOutput{
		Score: math.Min(score, 100),
	}
}

// calculateDeadlineScore 计算截止时间评分
func (s *sTask) calculateDeadlineScore(input *task.CalculateDeadlineScoreInput) *task.CalculateDeadlineScoreOutput {
	score := 0.0

	if deadlineAt, ok := input.TaskData["deadline_at"].(*gtime.Time); ok && deadlineAt != nil {
		now := time.Now()
		timeToDeadline := deadlineAt.Time.Sub(now)

		if timeToDeadline < 0 {
			// 已过期
			score = 100
		} else if timeToDeadline < 1*time.Hour {
			// 1小时内截止
			score = 95
		} else if timeToDeadline < 6*time.Hour {
			// 6小时内截止
			score = 85
		} else if timeToDeadline < 24*time.Hour {
			// 24小时内截止
			score = 75
		} else if timeToDeadline < 7*24*time.Hour {
			// 一周内截止
			score = 60
		} else {
			// 超过一周
			score = 30
		}
	}

	// 检查是否有时间约束
	if timeConstraint, ok := input.TaskData["time_constraint"].(string); ok {
		switch timeConstraint {
		case "strict":
			score += 20
		case "flexible":
			score += 10
		}
	}

	return &task.CalculateDeadlineScoreOutput{
		Score: math.Min(score, 100),
	}
}

// calculateResourceScore 计算资源需求评分
func (s *sTask) calculateResourceScore(input *task.CalculateResourceScoreInput) *task.CalculateResourceScoreOutput {
	score := 0.0

	// 检查预估执行时间
	if estimatedTime, ok := input.TaskData["estimated_duration"].(int); ok {
		if estimatedTime < 300 { // 5分钟以内
			score = 80 // 快速任务优先级高
		} else if estimatedTime < 1800 { // 30分钟以内
			score = 60
		} else if estimatedTime < 3600 { // 1小时以内
			score = 40
		} else {
			score = 20 // 长时间任务优先级低
		}
	}

	// 检查资源复杂度
	if complexity, ok := input.TaskData["complexity"].(string); ok {
		switch complexity {
		case "low":
			score += 10
		case "medium":
			score += 0
		case "high":
			score -= 10
		}
	}

	// 检查资源需求
	if resourceRequirements, ok := input.TaskData["resource_requirements"].(map[string]interface{}); ok {
		// CPU需求
		if cpuReq, ok := resourceRequirements["cpu"].(float64); ok {
			if cpuReq > 8 {
				score += 20
			} else if cpuReq > 4 {
				score += 15
			} else if cpuReq > 2 {
				score += 10
			}
		}

		// 内存需求
		if memoryReq, ok := resourceRequirements["memory"].(float64); ok {
			if memoryReq > 16 {
				score += 20
			} else if memoryReq > 8 {
				score += 15
			} else if memoryReq > 4 {
				score += 10
			}
		}

		// 存储需求
		if storageReq, ok := resourceRequirements["storage"].(float64); ok {
			if storageReq > 1000 {
				score += 15
			} else if storageReq > 500 {
				score += 10
			} else if storageReq > 100 {
				score += 5
			}
		}
	}

	// 检查资源可用性
	if resourceAvailability, ok := input.TaskData["resource_availability"].(string); ok {
		switch resourceAvailability {
		case "scarce":
			score += 25
		case "limited":
			score += 15
		case "available":
			score += 5
		}
	}

	// 检查资源成本
	if resourceCost, ok := input.TaskData["resource_cost"].(float64); ok {
		if resourceCost > 1000 {
			score += 20
		} else if resourceCost > 500 {
			score += 15
		} else if resourceCost > 100 {
			score += 10
		}
	}

	return &task.CalculateResourceScoreOutput{
		Score: math.Min(score, 100),
	}
}

// ===============================
// 任务优先级业务逻辑
// ===============================

// UpdateTaskPriority 更新任务优先级
func (s *sTask) UpdateTaskPriority(ctx context.Context, input *task.UpdateTaskPriorityInput) (*task.UpdateTaskPriorityOutput, error) {
	// 验证优先级范围
	if input.NewPriority < 1 || input.NewPriority > 10 {
		return &task.UpdateTaskPriorityOutput{
			Success: false,
			Message: "优先级必须在1-10范围内",
		}, gerror.New("优先级必须在1-10范围内")
	}

	// 检查任务是否存在
	taskStatusInput := &task.GetTaskStatusInput{TaskID: input.TaskID}
	taskStatusOutput := s.getTaskStatus(taskStatusInput)
	if taskStatusOutput.Status == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound, "任务不存在")
	}

	// 检查任务状态
	if taskStatusOutput.Status == "completed" {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "已完成的任务无法修改优先级")
	}

	// 这里应该更新数据库中的任务优先级
	// 目前返回模拟结果
	return &task.UpdateTaskPriorityOutput{
		Success:     true,
		Message:     "任务优先级更新成功",
		TaskID:      input.TaskID,
		OldPriority: 5, // 模拟旧优先级
		NewPriority: input.NewPriority,
		UpdatedAt:   gtime.Now().Format("2006-01-02 15:04:05"),
	}, nil
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
