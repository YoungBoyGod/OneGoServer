package queue

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/consts"
	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
)

// QueueBusiness 队列业务逻辑处理器
type QueueBusiness struct {
	validator *QueueValidator
}

// NewQueueBusiness 创建队列业务逻辑处理器实例
func NewQueueBusiness() *QueueBusiness {
	return &QueueBusiness{
		validator: NewQueueValidator(),
	}
}

// CalculateQueuePosition 计算新任务在队列中的位置
func (b *QueueBusiness) CalculateQueuePosition(existingQueue []DeviceTaskQueue, newTask *DeviceTaskQueue, strategy string) int {
	if len(existingQueue) == 0 {
		return 1
	}

	switch strategy {
	case SchedulingStrategyPriorityFirst:
		return b.calculatePositionByPriority(existingQueue, newTask)
	case SchedulingStrategyFIFO:
		return len(existingQueue) + 1 // 添加到队列末尾
	case SchedulingStrategyLIFO:
		return 1 // 添加到队列开头
	case SchedulingStrategyWeighted:
		return b.calculatePositionByWeight(existingQueue, newTask)
	default:
		return b.calculatePositionByPriority(existingQueue, newTask)
	}
}

// calculatePositionByPriority 根据优先级计算位置
func (b *QueueBusiness) calculatePositionByPriority(existingQueue []DeviceTaskQueue, newTask *DeviceTaskQueue) int {
	for i, item := range existingQueue {
		if newTask.QueuePriority > item.QueuePriority {
			return i + 1
		}
		// 相同优先级时，按创建时间排序
		if newTask.QueuePriority == item.QueuePriority && newTask.CreatedAt.Before(item.CreatedAt) {
			return i + 1
		}
	}
	return len(existingQueue) + 1
}

// calculatePositionByWeight 根据权重计算位置
func (b *QueueBusiness) calculatePositionByWeight(existingQueue []DeviceTaskQueue, newTask *DeviceTaskQueue) int {
	newWeight := b.calculateTaskWeight(newTask)

	for i, item := range existingQueue {
		itemWeight := b.calculateTaskWeight(&item)
		if newWeight > itemWeight {
			return i + 1
		}
	}
	return len(existingQueue) + 1
}

// calculateTaskWeight 计算任务权重
func (b *QueueBusiness) calculateTaskWeight(task *DeviceTaskQueue) float64 {
	weight := float64(task.QueuePriority) * 10.0 // 基础优先级权重

	// 手动设置优先级的任务权重更高
	if task.IsManualPriority {
		weight += 20.0
	}

	// 重试次数影响权重
	if task.CurrentRetry > 0 {
		weight += float64(task.CurrentRetry) * 5.0
	}

	// 等待时间影响权重
	waitTime := time.Since(task.CreatedAt).Hours()
	weight += waitTime * 2.0

	// 预估执行时长影响权重（短任务优先）
	if task.EstimatedDuration != nil {
		durationHours := float64(*task.EstimatedDuration) / 3600.0
		weight -= durationHours * 1.0
	}

	return weight
}

// ReorderQueue 重新排序队列
func (b *QueueBusiness) ReorderQueue(queue []DeviceTaskQueue, strategy string) []DeviceTaskQueue {
	sortedQueue := make([]DeviceTaskQueue, len(queue))
	copy(sortedQueue, queue)

	switch strategy {
	case SchedulingStrategyPriorityFirst:
		sort.Slice(sortedQueue, func(i, j int) bool {
			// 优先级高的在前
			if sortedQueue[i].QueuePriority != sortedQueue[j].QueuePriority {
				return sortedQueue[i].QueuePriority > sortedQueue[j].QueuePriority
			}
			// 优先级相同时，创建时间早的在前
			return sortedQueue[i].CreatedAt.Before(sortedQueue[j].CreatedAt)
		})
	case SchedulingStrategyFIFO:
		sort.Slice(sortedQueue, func(i, j int) bool {
			return sortedQueue[i].CreatedAt.Before(sortedQueue[j].CreatedAt)
		})
	case SchedulingStrategyLIFO:
		sort.Slice(sortedQueue, func(i, j int) bool {
			return sortedQueue[i].CreatedAt.After(sortedQueue[j].CreatedAt)
		})
	case SchedulingStrategyWeighted:
		sort.Slice(sortedQueue, func(i, j int) bool {
			weightI := b.calculateTaskWeight(&sortedQueue[i])
			weightJ := b.calculateTaskWeight(&sortedQueue[j])
			return weightI > weightJ
		})
	}

	// 更新位置
	for i := range sortedQueue {
		sortedQueue[i].QueuePosition = i + 1
	}

	return sortedQueue
}

// CalculateEstimatedStartTime 计算预估开始时间
func (b *QueueBusiness) CalculateEstimatedStartTime(queue []DeviceTaskQueue, position int) *time.Time {
	if position <= 1 {
		now := time.Now()
		return &now
	}

	currentTime := time.Now()

	// 计算前面任务的总执行时间
	for i := 0; i < position-1 && i < len(queue); i++ {
		if queue[i].EstimatedDuration != nil {
			currentTime = currentTime.Add(time.Duration(*queue[i].EstimatedDuration) * time.Second)
		} else {
			// 如果没有预估时长，使用默认值
			currentTime = currentTime.Add(30 * time.Minute)
		}
	}

	return &currentTime
}

// ValidateQueueOperation 验证队列操作是否允许
func (b *QueueBusiness) ValidateQueueOperation(item *DeviceTaskQueue, operation string) error {
	switch operation {
	case QueueOperationStart:
		if item.Status != DeviceQueueStatusQueued {
			return fmt.Errorf("只有排队状态的任务可以开始执行，当前状态: %s", item.Status)
		}
	case QueueOperationPause:
		if item.Status != DeviceQueueStatusExecuting {
			return fmt.Errorf("只有执行中的任务可以暂停，当前状态: %s", item.Status)
		}
	case QueueOperationResume:
		if item.Status != DeviceQueueStatusPaused {
			return fmt.Errorf("只有暂停状态的任务可以恢复，当前状态: %s", item.Status)
		}
	case QueueOperationCancel:
		if item.Status == DeviceQueueStatusCompleted || item.Status == DeviceQueueStatusCanceled {
			return fmt.Errorf("已完成或已取消的任务不能再次取消，当前状态: %s", item.Status)
		}
	case QueueOperationRemove:
		if item.Status == DeviceQueueStatusExecuting {
			return fmt.Errorf("执行中的任务不能移除，请先暂停任务，当前状态: %s", item.Status)
		}
	}
	return nil
}

// CalculateQueueMetrics 计算队列指标
func (b *QueueBusiness) CalculateQueueMetrics(queue []DeviceTaskQueue) map[string]interface{} {
	metrics := make(map[string]interface{})

	if len(queue) == 0 {
		return metrics
	}

	// 基础统计
	total := len(queue)
	statusCount := make(map[string]int)
	priorityCount := make(map[string]int)
	var totalWaitTime float64
	var totalDuration int64

	now := time.Now()

	for _, item := range queue {
		// 状态统计
		statusCount[item.Status]++

		// 优先级统计
		priorityKey := fmt.Sprintf("priority_%d", item.QueuePriority)
		priorityCount[priorityKey]++

		// 等待时间统计
		if item.Status == DeviceQueueStatusQueued {
			waitTime := now.Sub(item.CreatedAt).Hours()
			totalWaitTime += waitTime
		}

		// 预估执行时长统计
		if item.EstimatedDuration != nil {
			totalDuration += int64(*item.EstimatedDuration)
		}
	}

	metrics["total_tasks"] = total
	metrics["status_distribution"] = statusCount
	metrics["priority_distribution"] = priorityCount

	// 平均等待时间
	queuedCount := statusCount[DeviceQueueStatusQueued]
	if queuedCount > 0 {
		metrics["avg_wait_time_hours"] = totalWaitTime / float64(queuedCount)
	} else {
		metrics["avg_wait_time_hours"] = 0.0
	}

	// 总预估执行时间
	metrics["total_estimated_duration_seconds"] = totalDuration
	metrics["total_estimated_duration_hours"] = float64(totalDuration) / 3600.0

	// 队列健康度评分
	metrics["queue_health_score"] = b.calculateQueueHealthScore(queue)

	return metrics
}

// calculateQueueHealthScore 计算队列健康度评分
func (b *QueueBusiness) calculateQueueHealthScore(queue []DeviceTaskQueue) float64 {
	if len(queue) == 0 {
		return 100.0
	}

	score := 100.0

	// 失败任务比例扣分
	failedCount := 0
	completedCount := 0
	for _, item := range queue {
		if item.Status == DeviceQueueStatusFailed {
			failedCount++
		} else if item.Status == DeviceQueueStatusCompleted {
			completedCount++
		}
	}

	if failedCount+completedCount > 0 {
		failureRate := float64(failedCount) / float64(failedCount+completedCount)
		score -= failureRate * 30.0 // 失败率每10%扣3分
	}

	// 长时间等待的任务扣分
	now := time.Now()
	longWaitCount := 0
	for _, item := range queue {
		if item.Status == DeviceQueueStatusQueued {
			waitHours := now.Sub(item.CreatedAt).Hours()
			if waitHours > 24 { // 等待超过24小时
				longWaitCount++
			}
		}
	}

	if longWaitCount > 0 {
		longWaitRate := float64(longWaitCount) / float64(len(queue))
		score -= longWaitRate * 20.0
	}

	// 重试次数过多的任务扣分
	highRetryCount := 0
	for _, item := range queue {
		if item.CurrentRetry >= item.MaxRetryCount/2 { // 重试次数超过一半
			highRetryCount++
		}
	}

	if highRetryCount > 0 {
		highRetryRate := float64(highRetryCount) / float64(len(queue))
		score -= highRetryRate * 15.0
	}

	return math.Max(score, 0.0)
}

// OptimizeQueueOrder 优化队列顺序
func (b *QueueBusiness) OptimizeQueueOrder(queue []DeviceTaskQueue, config *DeviceQueueConfig) []DeviceTaskQueue {
	if len(queue) <= 1 {
		return queue
	}

	optimizedQueue := make([]DeviceTaskQueue, len(queue))
	copy(optimizedQueue, queue)

	// 根据配置的调度策略优化
	if config != nil && config.PriorityScheduling {
		optimizedQueue = b.ReorderQueue(optimizedQueue, config.SchedulingStrategy)
	}

	// 负载均衡优化
	if config != nil && config.LoadBalancing {
		optimizedQueue = b.applyLoadBalancing(optimizedQueue)
	}

	return optimizedQueue
}

// applyLoadBalancing 应用负载均衡
func (b *QueueBusiness) applyLoadBalancing(queue []DeviceTaskQueue) []DeviceTaskQueue {
	// 将长任务和短任务交替排列，避免连续执行长任务
	shortTasks := make([]DeviceTaskQueue, 0)
	longTasks := make([]DeviceTaskQueue, 0)

	avgDuration := b.calculateAverageDuration(queue)

	for _, task := range queue {
		duration := float64(consts.DefaultTimeoutSeconds) // 默认值
		if task.EstimatedDuration != nil {
			duration = float64(*task.EstimatedDuration)
		}

		if duration <= avgDuration {
			shortTasks = append(shortTasks, task)
		} else {
			longTasks = append(longTasks, task)
		}
	}

	// 交替合并
	balanced := make([]DeviceTaskQueue, 0, len(queue))
	maxLen := len(shortTasks)
	if len(longTasks) > maxLen {
		maxLen = len(longTasks)
	}

	for i := 0; i < maxLen; i++ {
		if i < len(shortTasks) {
			balanced = append(balanced, shortTasks[i])
		}
		if i < len(longTasks) {
			balanced = append(balanced, longTasks[i])
		}
	}

	// 更新位置
	for i := range balanced {
		balanced[i].QueuePosition = i + 1
	}

	return balanced
}

// calculateAverageDuration 计算平均执行时长
func (b *QueueBusiness) calculateAverageDuration(queue []DeviceTaskQueue) float64 {
	total := 0.0
	count := 0

	for _, task := range queue {
		if task.EstimatedDuration != nil {
			total += float64(*task.EstimatedDuration)
			count++
		}
	}

	if count == 0 {
		return float64(consts.DefaultTimeoutSeconds)
	}

	return total / float64(count)
}

// ParseDependencyIds 解析依赖任务ID列表
func (b *QueueBusiness) ParseDependencyIds(dependencyStr *string) []int64 {
	if dependencyStr == nil || *dependencyStr == "" {
		return []int64{}
	}

	ids := make([]int64, 0)
	parts := strings.Split(*dependencyStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if id, err := strconv.ParseInt(part, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

// FormatDependencyIds 格式化依赖任务ID列表
func (b *QueueBusiness) FormatDependencyIds(ids []int64) *string {
	if len(ids) == 0 {
		return nil
	}

	strIds := make([]string, len(ids))
	for i, id := range ids {
		strIds[i] = strconv.FormatInt(id, 10)
	}

	result := strings.Join(strIds, ",")
	return &result
}

// CheckDependencies 检查任务依赖是否满足
func (b *QueueBusiness) CheckDependencies(task *DeviceTaskQueue, allTasks []DeviceTaskQueue) bool {
	dependencyIds := b.ParseDependencyIds(task.DependsOnTaskIds)
	if len(dependencyIds) == 0 {
		return true // 没有依赖，可以执行
	}

	// 创建任务状态映射
	taskStatusMap := make(map[int64]string)
	for _, t := range allTasks {
		taskStatusMap[t.TaskID] = t.Status
	}

	// 检查所有依赖任务是否已完成
	for _, depId := range dependencyIds {
		status, exists := taskStatusMap[depId]
		if !exists || status != DeviceQueueStatusCompleted {
			return false
		}
	}

	return true
}

// GenerateOperationHistory 生成操作历史记录
func (b *QueueBusiness) GenerateOperationHistory(deviceID int64, deviceESN string, taskID *int64, operation string, operationBy *int64, reason *string) *DeviceQueueOperationHistory {
	batchID := utils.GenerateExecutionID()

	history := &DeviceQueueOperationHistory{
		DeviceID:         deviceID,
		DeviceESN:        deviceESN,
		TaskID:           taskID,
		OperationType:    operation,
		OperationBy:      operationBy,
		OperationTime:    time.Now(),
		Reason:           reason,
		OperationSource:  OperationSourceAPI,
		BatchID:          &batchID,
		IsBatchOperation: false,
	}

	return history
}

// CalculateQueueEfficiency 计算队列执行效率
func (b *QueueBusiness) CalculateQueueEfficiency(queue []DeviceTaskQueue) float64 {
	if len(queue) == 0 {
		return 0.0
	}

	totalTasks := len(queue)
	completedTasks := 0
	totalDuration := 0.0
	actualDuration := 0.0

	for _, task := range queue {
		if task.Status == DeviceQueueStatusCompleted {
			completedTasks++

			// 计算预估时长和实际时长
			if task.EstimatedDuration != nil {
				totalDuration += float64(*task.EstimatedDuration)
			}

			if task.ActualStartTime != nil && task.ActualEndTime != nil {
				actual := task.ActualEndTime.Sub(*task.ActualStartTime).Seconds()
				actualDuration += actual
			}
		}
	}

	if completedTasks == 0 {
		return 0.0
	}

	// 完成率
	completionRate := float64(completedTasks) / float64(totalTasks)

	// 时间效率（实际时间/预估时间）
	timeEfficiency := 1.0
	if totalDuration > 0 && actualDuration > 0 {
		timeEfficiency = totalDuration / actualDuration
		if timeEfficiency > 1.0 {
			timeEfficiency = 1.0 // 实际时间短于预估时间，效率为100%
		}
	}

	// 综合效率 = 完成率 * 时间效率
	return completionRate * timeEfficiency * 100.0
}

// ShouldRetryTask 判断任务是否应该重试
func (b *QueueBusiness) ShouldRetryTask(task *DeviceTaskQueue) bool {
	if task.Status != DeviceQueueStatusFailed {
		return false
	}

	if task.CurrentRetry >= task.MaxRetryCount {
		return false
	}

	// TODO: 可以根据失败原因判断是否值得重试

	return true
}

// EstimateQueueCompletion 估算队列完成时间
func (b *QueueBusiness) EstimateQueueCompletion(queue []DeviceTaskQueue, concurrency int) *time.Time {
	if len(queue) == 0 {
		now := time.Now()
		return &now
	}

	// 只计算未完成的任务
	pendingTasks := make([]DeviceTaskQueue, 0)
	for _, task := range queue {
		if task.Status == DeviceQueueStatusQueued || task.Status == DeviceQueueStatusExecuting {
			pendingTasks = append(pendingTasks, task)
		}
	}

	if len(pendingTasks) == 0 {
		now := time.Now()
		return &now
	}

	// 计算总执行时间
	totalDuration := 0
	for _, task := range pendingTasks {
		if task.EstimatedDuration != nil {
			totalDuration += *task.EstimatedDuration
		} else {
			totalDuration += consts.DefaultTimeoutSeconds
		}
	}

	// 考虑并发执行
	if concurrency > 1 {
		totalDuration = totalDuration / concurrency
	}

	estimatedTime := time.Now().Add(time.Duration(totalDuration) * time.Second)
	return &estimatedTime
}
