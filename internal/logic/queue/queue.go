package queue

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

type sQueue struct{}

func New() *sQueue {
	return &sQueue{}
}

func init() {
	// service.RegisterQueue(New())
}

// ===============================
// 队列基础管理业务逻辑
// ===============================

// ValidateQueueCreation 验证队列创建
func (s *sQueue) ValidateQueueCreation(ctx context.Context, queueData map[string]interface{}) error {
	// 验证队列名称
	if name, ok := queueData["name"].(string); ok {
		if err := s.validateQueueName(name); err != nil {
			return err
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "队列名称不能为空")
	}

	// 验证队列类型
	if queueType, ok := queueData["type"].(string); ok {
		if err := s.validateQueueType(queueType); err != nil {
			return err
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "队列类型不能为空")
	}

	// 验证容量配置
	if capacity, ok := queueData["capacity"].(int); ok {
		if capacity <= 0 || capacity > 10000 {
			return gerror.NewCode(gcode.CodeValidationFailed, "队列容量必须在1-10000之间")
		}
	}

	// 验证优先级配置
	if maxPriority, ok := queueData["max_priority"].(int); ok {
		if maxPriority <= 0 || maxPriority > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "最大优先级必须在1-100之间")
		}
	}

	return nil
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
	if !strings.MatchString(`^[a-zA-Z0-9_-]+$`, name) {
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

// ===============================
// 队列操作控制业务逻辑
// ===============================

// ValidateQueueOperation 验证队列操作
func (s *sQueue) ValidateQueueOperation(ctx context.Context, queueData map[string]interface{}, operation string) error {
	// 检查队列状态
	if status, ok := queueData["status"].(string); ok {
		switch operation {
		case "start":
			if status == "running" {
				return gerror.NewCode(gcode.CodeInvalidOperation, "队列已在运行状态")
			}
		case "stop":
			if status == "stopped" {
				return gerror.NewCode(gcode.CodeInvalidOperation, "队列已在停止状态")
			}
		case "pause":
			if status == "paused" {
				return gerror.NewCode(gcode.CodeInvalidOperation, "队列已在暂停状态")
			}
		case "resume":
			if status != "paused" {
				return gerror.NewCode(gcode.CodeInvalidOperation, "只有暂停状态的队列才能恢复")
			}
		}
	}

	// 检查队列健康状态
	if healthScore, ok := queueData["health_score"].(float64); ok {
		if healthScore < 30 && operation == "start" {
			return gerror.NewCode(gcode.CodeInvalidOperation, "队列健康度过低，无法启动")
		}
	}

	return nil
}

// CalculateQueueHealthScore 计算队列健康度评分
func (s *sQueue) CalculateQueueHealthScore(ctx context.Context, queueData map[string]interface{}) float64 {
	score := 100.0

	// 1. 队列状态评分 (权重: 30%)
	statusScore := s.calculateStatusScore(queueData)
	score += statusScore * 0.30

	// 2. 性能指标评分 (权重: 25%)
	performanceScore := s.calculatePerformanceScore(queueData)
	score += performanceScore * 0.25

	// 3. 错误率评分 (权重: 20%)
	errorScore := s.calculateErrorScore(queueData)
	score += errorScore * 0.20

	// 4. 资源使用评分 (权重: 15%)
	resourceScore := s.calculateResourceScore(queueData)
	score += resourceScore * 0.15

	// 5. 响应时间评分 (权重: 10%)
	responseScore := s.calculateResponseScore(queueData)
	score += responseScore * 0.10

	return math.Max(0, math.Min(score, 100))
}

// calculateStatusScore 计算状态评分
func (s *sQueue) calculateStatusScore(queueData map[string]interface{}) float64 {
	status, _ := queueData["status"].(string)
	switch status {
	case "running":
		return 100
	case "paused":
		return 70
	case "stopped":
		return 50
	case "error":
		return 20
	default:
		return 30
	}
}

// calculatePerformanceScore 计算性能评分
func (s *sQueue) calculatePerformanceScore(queueData map[string]interface{}) float64 {
	score := 100.0

	// 检查处理速率
	if rate, ok := queueData["processing_rate"].(float64); ok {
		if rate < 10 {
			score -= 30
		} else if rate < 50 {
			score -= 15
		}
	}

	// 检查队列长度
	if length, ok := queueData["current_length"].(int); ok {
		if capacity, ok := queueData["capacity"].(int); ok && capacity > 0 {
			utilization := float64(length) / float64(capacity)
			if utilization > 0.9 {
				score -= 40
			} else if utilization > 0.7 {
				score -= 20
			}
		}
	}

	return math.Max(0, score)
}

// calculateErrorScore 计算错误率评分
func (s *sQueue) calculateErrorScore(queueData map[string]interface{}) float64 {
	score := 100.0

	if errorRate, ok := queueData["error_rate"].(float64); ok {
		if errorRate > 0.1 { // 错误率超过10%
			score -= 60
		} else if errorRate > 0.05 { // 错误率超过5%
			score -= 30
		} else if errorRate > 0.01 { // 错误率超过1%
			score -= 10
		}
	}

	return math.Max(0, score)
}

// calculateResourceScore 计算资源使用评分
func (s *sQueue) calculateResourceScore(queueData map[string]interface{}) float64 {
	score := 100.0

	// 检查内存使用
	if memoryUsage, ok := queueData["memory_usage"].(float64); ok {
		if memoryUsage > 0.9 {
			score -= 40
		} else if memoryUsage > 0.7 {
			score -= 20
		}
	}

	// 检查CPU使用
	if cpuUsage, ok := queueData["cpu_usage"].(float64); ok {
		if cpuUsage > 0.9 {
			score -= 30
		} else if cpuUsage > 0.7 {
			score -= 15
		}
	}

	return math.Max(0, score)
}

// calculateResponseScore 计算响应时间评分
func (s *sQueue) calculateResponseScore(queueData map[string]interface{}) float64 {
	score := 100.0

	if avgResponseTime, ok := queueData["avg_response_time"].(float64); ok {
		if avgResponseTime > 5000 { // 超过5秒
			score -= 50
		} else if avgResponseTime > 2000 { // 超过2秒
			score -= 30
		} else if avgResponseTime > 1000 { // 超过1秒
			score -= 15
		}
	}

	return math.Max(0, score)
}

// ===============================
// 队列任务管理业务逻辑
// ===============================

// ValidateTaskEnqueue 验证任务入队
func (s *sQueue) ValidateTaskEnqueue(ctx context.Context, queueData map[string]interface{}, taskData map[string]interface{}) error {
	// 检查队列状态
	if status, ok := queueData["status"].(string); ok {
		if status != "running" && status != "paused" {
			return gerror.NewCode(gcode.CodeInvalidOperation, "队列未运行，无法入队任务")
		}
	}

	// 检查队列容量
	if currentLength, ok := queueData["current_length"].(int); ok {
		if capacity, ok := queueData["capacity"].(int); ok {
			if currentLength >= capacity {
				return gerror.NewCode(gcode.CodeResourceExhausted, "队列已满，无法入队新任务")
			}
		}
	}

	// 验证任务数据
	if err := s.validateTaskData(taskData); err != nil {
		return err
	}

	// 检查任务优先级
	if taskPriority, ok := taskData["priority"].(int); ok {
		if maxPriority, ok := queueData["max_priority"].(int); ok {
			if taskPriority > maxPriority {
				return gerror.NewCode(gcode.CodeValidationFailed,
					fmt.Sprintf("任务优先级超过队列最大优先级限制: %d", maxPriority))
			}
		}
	}

	return nil
}

// validateTaskData 验证任务数据
func (s *sQueue) validateTaskData(taskData map[string]interface{}) error {
	// 检查任务ID
	if taskId, ok := taskData["task_id"].(string); ok {
		if taskId == "" {
			return gerror.NewCode(gcode.CodeValidationFailed, "任务ID不能为空")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "任务ID不能为空")
	}

	// 检查任务类型
	if taskType, ok := taskData["type"].(string); ok {
		if taskType == "" {
			return gerror.NewCode(gcode.CodeValidationFailed, "任务类型不能为空")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "任务类型不能为空")
	}

	// 检查超时时间
	if timeout, ok := taskData["timeout"].(int); ok {
		if timeout <= 0 || timeout > 3600 {
			return gerror.NewCode(gcode.CodeValidationFailed, "任务超时时间必须在1-3600秒之间")
		}
	}

	return nil
}

// ===============================
// 队列排序算法业务逻辑
// ===============================

// SortQueueTasks 队列任务排序
func (s *sQueue) SortQueueTasks(ctx context.Context, queueType string, tasks []map[string]interface{}) []map[string]interface{} {
	switch queueType {
	case "fifo":
		return s.sortFIFO(tasks)
	case "lifo":
		return s.sortLIFO(tasks)
	case "priority":
		return s.sortByPriority(tasks)
	case "round_robin":
		return s.sortRoundRobin(tasks)
	case "weighted":
		return s.sortByWeight(tasks)
	default:
		return s.sortFIFO(tasks) // 默认FIFO
	}
}

// sortFIFO 先进先出排序
func (s *sQueue) sortFIFO(tasks []map[string]interface{}) []map[string]interface{} {
	sort.Slice(tasks, func(i, j int) bool {
		timeI, _ := tasks[i]["created_at"].(*gtime.Time)
		timeJ, _ := tasks[j]["created_at"].(*gtime.Time)

		if timeI == nil || timeJ == nil {
			return false
		}

		return timeI.Before(timeJ.Time)
	})
	return tasks
}

// sortLIFO 后进先出排序
func (s *sQueue) sortLIFO(tasks []map[string]interface{}) []map[string]interface{} {
	sort.Slice(tasks, func(i, j int) bool {
		timeI, _ := tasks[i]["created_at"].(*gtime.Time)
		timeJ, _ := tasks[j]["created_at"].(*gtime.Time)

		if timeI == nil || timeJ == nil {
			return false
		}

		return timeI.After(timeJ.Time)
	})
	return tasks
}

// sortByPriority 按优先级排序
func (s *sQueue) sortByPriority(tasks []map[string]interface{}) []map[string]interface{} {
	sort.Slice(tasks, func(i, j int) bool {
		priorityI, _ := tasks[i]["priority"].(int)
		priorityJ, _ := tasks[j]["priority"].(int)

		// 优先级高的在前
		if priorityI != priorityJ {
			return priorityI > priorityJ
		}

		// 优先级相同时，按创建时间排序
		timeI, _ := tasks[i]["created_at"].(*gtime.Time)
		timeJ, _ := tasks[j]["created_at"].(*gtime.Time)

		if timeI == nil || timeJ == nil {
			return false
		}

		return timeI.Before(timeJ.Time)
	})
	return tasks
}

// sortRoundRobin 轮询排序
func (s *sQueue) sortRoundRobin(tasks []map[string]interface{}) []map[string]interface{} {
	// 按任务类型分组
	typeGroups := make(map[string][]map[string]interface{})
	for _, task := range tasks {
		taskType, _ := task["type"].(string)
		typeGroups[taskType] = append(typeGroups[taskType], task)
	}

	// 轮询排序
	var result []map[string]interface{}
	indices := make(map[string]int)

	for len(result) < len(tasks) {
		for taskType, typeTasks := range typeGroups {
			if indices[taskType] < len(typeTasks) {
				result = append(result, typeTasks[indices[taskType]])
				indices[taskType]++
			}
		}
	}

	return result
}

// sortByWeight 按权重排序
func (s *sQueue) sortByWeight(tasks []map[string]interface{}) []map[string]interface{} {
	sort.Slice(tasks, func(i, j int) bool {
		weightI := s.calculateTaskWeight(tasks[i])
		weightJ := s.calculateTaskWeight(tasks[j])

		return weightI > weightJ
	})
	return tasks
}

// calculateTaskWeight 计算任务权重
func (s *sQueue) calculateTaskWeight(task map[string]interface{}) float64 {
	weight := 0.0

	// 基础权重
	if priority, ok := task["priority"].(int); ok {
		weight += float64(priority) * 0.4
	}

	// 紧急程度权重
	if isUrgent, ok := task["is_urgent"].(bool); ok && isUrgent {
		weight += 30.0
	}

	// 等待时间权重
	if createdAt, ok := task["created_at"].(*gtime.Time); ok && createdAt != nil {
		waitTime := time.Since(createdAt.Time).Seconds()
		weight += math.Min(waitTime/60, 20) // 最多20分权重
	}

	// 业务重要性权重
	if importance, ok := task["business_importance"].(int); ok {
		weight += float64(importance) * 0.3
	}

	return weight
}

// ===============================
// 负载均衡业务逻辑
// ===============================

// CalculateLoadBalance 计算负载均衡
func (s *sQueue) CalculateLoadBalance(ctx context.Context, queues []map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 计算总体负载
	totalLoad := 0.0
	activeQueues := 0

	for _, queue := range queues {
		if status, _ := queue["status"].(string); status == "running" {
			if load, ok := queue["current_load"].(float64); ok {
				totalLoad += load
				activeQueues++
			}
		}
	}

	if activeQueues == 0 {
		result["balanced"] = true
		result["recommendation"] = "无活跃队列"
		return result
	}

	// 计算平均负载
	avgLoad := totalLoad / float64(activeQueues)

	// 检查负载分布
	unbalancedQueues := 0
	var recommendations []string

	for _, queue := range queues {
		if status, _ := queue["status"].(string); status == "running" {
			if load, ok := queue["current_load"].(float64); ok {
				deviation := math.Abs(load-avgLoad) / avgLoad
				if deviation > 0.3 { // 负载偏差超过30%
					unbalancedQueues++
					queueName, _ := queue["name"].(string)
					if load > avgLoad {
						recommendations = append(recommendations,
							fmt.Sprintf("队列 %s 负载过高，建议减少任务分配", queueName))
					} else {
						recommendations = append(recommendations,
							fmt.Sprintf("队列 %s 负载较低，可以增加任务分配", queueName))
					}
				}
			}
		}
	}

	result["balanced"] = unbalancedQueues == 0
	result["average_load"] = avgLoad
	result["unbalanced_queues"] = unbalancedQueues
	result["recommendations"] = recommendations

	return result
}

// SelectOptimalQueue 选择最优队列
func (s *sQueue) SelectOptimalQueue(ctx context.Context, taskData map[string]interface{}, availableQueues []map[string]interface{}) (map[string]interface{}, error) {
	if len(availableQueues) == 0 {
		return nil, gerror.NewCode(gcode.CodeResourceExhausted, "没有可用的队列")
	}

	var bestQueue map[string]interface{}
	bestScore := -1.0

	for _, queue := range availableQueues {
		score := s.calculateQueueScore(queue, taskData)
		if score > bestScore {
			bestScore = score
			bestQueue = queue
		}
	}

	if bestQueue == nil {
		return nil, gerror.NewCode(gcode.CodeResourceExhausted, "没有合适的队列")
	}

	return bestQueue, nil
}

// calculateQueueScore 计算队列评分
func (s *sQueue) calculateQueueScore(queue map[string]interface{}, taskData map[string]interface{}) float64 {
	score := 0.0

	// 1. 队列健康度 (权重: 30%)
	if healthScore, ok := queue["health_score"].(float64); ok {
		score += healthScore * 0.30
	}

	// 2. 负载情况 (权重: 25%)
	if currentLoad, ok := queue["current_load"].(float64); ok {
		// 负载越低分数越高
		loadScore := 100 - currentLoad
		score += loadScore * 0.25
	}

	// 3. 响应时间 (权重: 20%)
	if avgResponseTime, ok := queue["avg_response_time"].(float64); ok {
		// 响应时间越短分数越高
		responseScore := math.Max(0, 100-avgResponseTime/10)
		score += responseScore * 0.20
	}

	// 4. 错误率 (权重: 15%)
	if errorRate, ok := queue["error_rate"].(float64); ok {
		// 错误率越低分数越高
		errorScore := 100 - errorRate*100
		score += errorScore * 0.15
	}

	// 5. 队列类型匹配度 (权重: 10%)
	if queueType, ok := queue["type"].(string); ok {
		if taskType, ok := taskData["type"].(string); ok {
			if s.isTypeCompatible(queueType, taskType) {
				score += 100 * 0.10
			}
		}
	}

	return score
}

// isTypeCompatible 检查类型兼容性
func (s *sQueue) isTypeCompatible(queueType, taskType string) bool {
	// 这里可以根据实际业务逻辑定义类型兼容性规则
	compatibilityMap := map[string][]string{
		"priority": {"high_priority", "normal", "low_priority"},
		"fifo":     {"normal", "batch"},
		"lifo":     {"normal", "batch"},
		"weighted": {"high_priority", "normal", "low_priority"},
	}

	if compatibleTypes, exists := compatibilityMap[queueType]; exists {
		for _, compatibleType := range compatibleTypes {
			if taskType == compatibleType {
				return true
			}
		}
	}

	return false
}

// ===============================
// 监控统计业务逻辑
// ===============================

// CalculateQueueStatistics 计算队列统计信息
func (s *sQueue) CalculateQueueStatistics(ctx context.Context, queueData map[string]interface{}, historicalData []map[string]interface{}) map[string]interface{} {
	stats := make(map[string]interface{})

	// 基础统计
	stats["current_length"] = queueData["current_length"]
	stats["capacity"] = queueData["capacity"]
	stats["utilization_rate"] = s.calculateUtilizationRate(queueData)

	// 性能统计
	stats["avg_processing_time"] = s.calculateAverageProcessingTime(historicalData)
	stats["throughput"] = s.calculateThroughput(historicalData)
	stats["error_rate"] = s.calculateErrorRate(historicalData)

	// 趋势分析
	stats["trend_analysis"] = s.analyzeQueueTrend(historicalData)

	// 预测分析
	stats["prediction"] = s.predictQueueBehavior(historicalData)

	return stats
}

// calculateUtilizationRate 计算队列利用率
func (s *sQueue) calculateUtilizationRate(queueData map[string]interface{}) float64 {
	currentLength, _ := queueData["current_length"].(int)
	capacity, _ := queueData["capacity"].(int)

	if capacity <= 0 {
		return 0.0
	}

	return float64(currentLength) / float64(capacity) * 100
}

// calculateAverageProcessingTime 计算平均处理时间
func (s *sQueue) calculateAverageProcessingTime(historicalData []map[string]interface{}) float64 {
	if len(historicalData) == 0 {
		return 0.0
	}

	totalTime := 0.0
	count := 0

	for _, data := range historicalData {
		if processingTime, ok := data["processing_time"].(float64); ok {
			totalTime += processingTime
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return totalTime / float64(count)
}

// calculateThroughput 计算吞吐量
func (s *sQueue) calculateThroughput(historicalData []map[string]interface{}) float64 {
	if len(historicalData) == 0 {
		return 0.0
	}

	// 计算最近1小时的处理任务数
	oneHourAgo := time.Now().Add(-time.Hour)
	processedCount := 0

	for _, data := range historicalData {
		if completedAt, ok := data["completed_at"].(*gtime.Time); ok && completedAt != nil {
			if completedAt.After(oneHourAgo) {
				processedCount++
			}
		}
	}

	return float64(processedCount) // 每小时处理的任务数
}

// calculateErrorRate 计算错误率
func (s *sQueue) calculateErrorRate(historicalData []map[string]interface{}) float64 {
	if len(historicalData) == 0 {
		return 0.0
	}

	totalTasks := len(historicalData)
	errorTasks := 0

	for _, data := range historicalData {
		if status, ok := data["status"].(string); ok && status == "failed" {
			errorTasks++
		}
	}

	return float64(errorTasks) / float64(totalTasks) * 100
}

// analyzeQueueTrend 分析队列趋势
func (s *sQueue) analyzeQueueTrend(historicalData []map[string]interface{}) map[string]interface{} {
	trend := make(map[string]interface{})

	if len(historicalData) < 2 {
		trend["trend"] = "insufficient_data"
		return trend
	}

	// 按时间排序
	sort.Slice(historicalData, func(i, j int) bool {
		timeI, _ := historicalData[i]["created_at"].(*gtime.Time)
		timeJ, _ := historicalData[j]["created_at"].(*gtime.Time)

		if timeI == nil || timeJ == nil {
			return false
		}

		return timeI.Before(timeJ.Time)
	})

	// 计算趋势
	recentData := historicalData[len(historicalData)/2:]
	olderData := historicalData[:len(historicalData)/2]

	recentAvg := s.calculateAverageProcessingTime(recentData)
	olderAvg := s.calculateAverageProcessingTime(olderData)

	if recentAvg > olderAvg*1.1 {
		trend["trend"] = "increasing"
		trend["change_percentage"] = (recentAvg - olderAvg) / olderAvg * 100
	} else if recentAvg < olderAvg*0.9 {
		trend["trend"] = "decreasing"
		trend["change_percentage"] = (olderAvg - recentAvg) / olderAvg * 100
	} else {
		trend["trend"] = "stable"
		trend["change_percentage"] = 0.0
	}

	return trend
}

// predictQueueBehavior 预测队列行为
func (s *sQueue) predictQueueBehavior(historicalData []map[string]interface{}) map[string]interface{} {
	prediction := make(map[string]interface{})

	if len(historicalData) < 10 {
		prediction["prediction"] = "insufficient_data"
		return prediction
	}

	// 简单的线性回归预测
	// 这里可以实现更复杂的预测算法

	// 预测未来1小时的负载
	currentLoad := s.calculateAverageProcessingTime(historicalData)
	predictedLoad := currentLoad * 1.05 // 假设增长5%

	prediction["predicted_load"] = predictedLoad
	prediction["confidence"] = 0.75
	prediction["timeframe"] = "1_hour"

	return prediction
}

// ===============================
// 队列配置管理业务逻辑
// ===============================

// ValidateQueueConfiguration 验证队列配置
func (s *sQueue) ValidateQueueConfiguration(ctx context.Context, config map[string]interface{}) error {
	// 验证基础配置
	if err := s.validateBasicConfig(config); err != nil {
		return err
	}

	// 验证性能配置
	if err := s.validatePerformanceConfig(config); err != nil {
		return err
	}

	// 验证安全配置
	if err := s.validateSecurityConfig(config); err != nil {
		return err
	}

	return nil
}

// validateBasicConfig 验证基础配置
func (s *sQueue) validateBasicConfig(config map[string]interface{}) error {
	// 验证队列名称
	if name, ok := config["name"].(string); ok {
		if err := s.validateQueueName(name); err != nil {
			return err
		}
	}

	// 验证容量配置
	if capacity, ok := config["capacity"].(int); ok {
		if capacity <= 0 || capacity > 100000 {
			return gerror.NewCode(gcode.CodeValidationFailed, "队列容量必须在1-100000之间")
		}
	}

	return nil
}

// validatePerformanceConfig 验证性能配置
func (s *sQueue) validatePerformanceConfig(config map[string]interface{}) error {
	// 验证并发数
	if concurrency, ok := config["max_concurrency"].(int); ok {
		if concurrency <= 0 || concurrency > 1000 {
			return gerror.NewCode(gcode.CodeValidationFailed, "最大并发数必须在1-1000之间")
		}
	}

	// 验证超时时间
	if timeout, ok := config["default_timeout"].(int); ok {
		if timeout <= 0 || timeout > 3600 {
			return gerror.NewCode(gcode.CodeValidationFailed, "默认超时时间必须在1-3600秒之间")
		}
	}

	return nil
}

// validateSecurityConfig 验证安全配置
func (s *sQueue) validateSecurityConfig(config map[string]interface{}) error {
	// 验证访问控制
	if accessControl, ok := config["access_control"].(map[string]interface{}); ok {
		if err := s.validateAccessControl(accessControl); err != nil {
			return err
		}
	}

	return nil
}

// validateAccessControl 验证访问控制配置
func (s *sQueue) validateAccessControl(accessControl map[string]interface{}) error {
	// 验证允许的用户组
	if allowedGroups, ok := accessControl["allowed_groups"].([]interface{}); ok {
		for _, group := range allowedGroups {
			if groupStr, ok := group.(string); ok {
				if groupStr == "" {
					return gerror.NewCode(gcode.CodeValidationFailed, "用户组名称不能为空")
				}
			}
		}
	}

	return nil
}

// ===============================
// 工具方法
// ===============================

// GetQueueInstance 获取队列实例（保留原有方法）
func (s *sQueue) GetQueueInstance(ctx context.Context) (interface{}, error) {
	// 这里可以返回具体的队列客户端实例
	// 目前返回nil，实际使用时需要根据具体需求实现
	return nil, nil
}
