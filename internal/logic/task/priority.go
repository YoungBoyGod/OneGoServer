package task

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务优先级计算相关业务逻辑
// ===============================

// CalculateTaskPriority 计算任务智能优先级
func (s *sTask) CalculateTaskPriority(ctx context.Context, taskData map[string]interface{}) int {
	baseScore := 0.0

	// 1. 基础优先级权重 (权重: 30%)
	if priority, ok := taskData["priority"].(int); ok && priority >= 1 && priority <= 10 {
		baseScore += float64(priority) * 10 * 0.3
	} else {
		baseScore += 50 * 0.3 // 默认中等优先级
	}

	// 2. 紧急程度权重 (权重: 25%)
	urgencyScore := s.calculateUrgencyScore(taskData)
	baseScore += urgencyScore * 0.25

	// 3. 业务重要性权重 (权重: 20%)
	businessScore := s.calculateBusinessImportanceScore(taskData)
	baseScore += businessScore * 0.20

	// 4. 截止时间影响 (权重: 15%)
	deadlineScore := s.calculateDeadlineScore(taskData)
	baseScore += deadlineScore * 0.15

	// 5. 资源需求权重 (权重: 10%)
	resourceScore := s.calculateResourceScore(taskData)
	baseScore += resourceScore * 0.10

	// 将评分转换为1-10的优先级
	priority := int(math.Round(baseScore / 10))
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}

	return priority
}

// calculateUrgencyScore 计算紧急程度评分
func (s *sTask) calculateUrgencyScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查任务类型的紧急性
	if taskType, ok := taskData["task_type"].(string); ok {
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
	if tags, ok := taskData["tags"].([]string); ok {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), "urgent") ||
				strings.Contains(strings.ToLower(tag), "emergency") {
				score = math.Max(score, 90)
			}
		}
	}

	return score
}

// calculateBusinessImportanceScore 计算业务重要性评分
func (s *sTask) calculateBusinessImportanceScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查任务类别
	if category, ok := taskData["category"].(string); ok {
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

	// 检查影响范围
	if impact, ok := taskData["impact_scope"].(string); ok {
		switch impact {
		case "global":
			score += 20
		case "regional":
			score += 10
		case "local":
			score += 5
		}
	}

	return math.Min(score, 100)
}

// calculateDeadlineScore 计算截止时间评分
func (s *sTask) calculateDeadlineScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	if deadlineAt, ok := taskData["deadline_at"].(*gtime.Time); ok && deadlineAt != nil {
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

	return score
}

// calculateResourceScore 计算资源需求评分
func (s *sTask) calculateResourceScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查预估执行时间
	if estimatedTime, ok := taskData["estimated_duration"].(int); ok {
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
	if complexity, ok := taskData["complexity"].(string); ok {
		switch complexity {
		case "low":
			score += 10
		case "medium":
			score += 0
		case "high":
			score -= 10
		}
	}

	return math.Max(score, 0)
}
