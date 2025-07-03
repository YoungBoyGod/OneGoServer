package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/consts"
)

// QueueLifecycleManager 队列生命周期管理器
type QueueLifecycleManager struct {
	business  *QueueBusiness
	validator *QueueValidator
	mutex     sync.RWMutex

	// 状态监听器
	stateListeners map[string][]QueueStateChangeListener

	// 调度器控制
	schedulerRunning bool
	schedulerStop    chan bool

	// 队列配置缓存
	configCache map[int64]*DeviceQueueConfig
	configMutex sync.RWMutex
}

// QueueStateChangeListener 队列状态变更监听器
type QueueStateChangeListener func(ctx context.Context, queueItem *DeviceTaskQueue, oldStatus, newStatus string) error

// QueueMetrics 队列指标
type QueueMetrics struct {
	DeviceID         int64     `json:"device_id"`
	TotalTasks       int       `json:"total_tasks"`
	QueuedTasks      int       `json:"queued_tasks"`
	ExecutingTasks   int       `json:"executing_tasks"`
	CompletedTasks   int       `json:"completed_tasks"`
	FailedTasks      int       `json:"failed_tasks"`
	AverageWaitTime  float64   `json:"average_wait_time"`
	QueueHealthScore float64   `json:"queue_health_score"`
	LastUpdated      time.Time `json:"last_updated"`
}

// NewQueueLifecycleManager 创建队列生命周期管理器
func NewQueueLifecycleManager() *QueueLifecycleManager {
	return &QueueLifecycleManager{
		business:         NewQueueBusiness(),
		validator:        NewQueueValidator(),
		stateListeners:   make(map[string][]QueueStateChangeListener),
		schedulerRunning: false,
		schedulerStop:    make(chan bool),
		configCache:      make(map[int64]*DeviceQueueConfig),
	}
}

// TransitionQueueItem 转换队列项状态
func (m *QueueLifecycleManager) TransitionQueueItem(ctx context.Context, item *DeviceTaskQueue, newStatus string, reason *string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldStatus := item.Status

	// 验证状态转换
	if err := m.validator.ValidateStatusTransition(oldStatus, newStatus); err != nil {
		return fmt.Errorf("状态转换验证失败: %w", err)
	}

	// 记录原始状态
	originalItem := *item

	// 执行状态转换
	item.Status = newStatus
	item.UpdatedAt = time.Now()

	// 根据新状态更新相关字段
	switch newStatus {
	case DeviceQueueStatusExecuting:
		now := time.Now()
		item.ActualStartTime = &now
	case DeviceQueueStatusCompleted, DeviceQueueStatusFailed, DeviceQueueStatusCanceled:
		if item.ActualStartTime != nil {
			now := time.Now()
			item.ActualEndTime = &now
		}
	}

	// 触发状态变更监听器
	if err := m.triggerStateListeners(ctx, item, oldStatus, newStatus); err != nil {
		// 如果监听器失败，回滚状态
		*item = originalItem
		return fmt.Errorf("状态监听器执行失败: %w", err)
	}

	return nil
}

// triggerStateListeners 触发状态变更监听器
func (m *QueueLifecycleManager) triggerStateListeners(ctx context.Context, item *DeviceTaskQueue, oldStatus, newStatus string) error {
	// 触发通用监听器
	if listeners, exists := m.stateListeners["*"]; exists {
		for _, listener := range listeners {
			if err := listener(ctx, item, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	// 触发特定状态监听器
	if listeners, exists := m.stateListeners[newStatus]; exists {
		for _, listener := range listeners {
			if err := listener(ctx, item, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

// AddStateListener 添加状态变更监听器
func (m *QueueLifecycleManager) AddStateListener(status string, listener QueueStateChangeListener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.stateListeners[status] == nil {
		m.stateListeners[status] = make([]QueueStateChangeListener, 0)
	}
	m.stateListeners[status] = append(m.stateListeners[status], listener)
}

// RemoveStateListener 移除状态变更监听器
func (m *QueueLifecycleManager) RemoveStateListener(status string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.stateListeners, status)
}

// ScheduleQueueExecution 调度队列执行
func (m *QueueLifecycleManager) ScheduleQueueExecution(ctx context.Context, deviceID int64, queue []DeviceTaskQueue) error {
	config := m.getDeviceQueueConfig(deviceID)
	if config == nil {
		return errors.New("设备队列配置不存在")
	}

	// 验证工作时间
	if err := m.validator.ValidateWorkingHours(config); err != nil {
		return fmt.Errorf("不在工作时间内: %w", err)
	}

	// 检查并发限制
	executingCount := m.countExecutingTasks(queue)
	if executingCount >= config.MaxConcurrentTasks {
		return fmt.Errorf("已达到最大并发任务数限制: %d", config.MaxConcurrentTasks)
	}

	// 获取可执行的任务
	readyTasks := m.getReadyTasks(queue, config.MaxConcurrentTasks-executingCount)

	// 执行任务
	for _, task := range readyTasks {
		if err := m.executeQueueItem(ctx, &task); err != nil {
			// 记录错误但继续执行其他任务
			continue
		}
	}

	return nil
}

// getReadyTasks 获取准备执行的任务
func (m *QueueLifecycleManager) getReadyTasks(queue []DeviceTaskQueue, maxCount int) []DeviceTaskQueue {
	readyTasks := make([]DeviceTaskQueue, 0, maxCount)

	for _, task := range queue {
		if len(readyTasks) >= maxCount {
			break
		}

		// 只处理排队状态的任务
		if task.Status != DeviceQueueStatusQueued {
			continue
		}

		// 检查依赖关系
		if !m.business.CheckDependencies(&task, queue) {
			continue
		}

		// 检查预估开始时间
		if task.EstimatedStartTime != nil && task.EstimatedStartTime.After(time.Now()) {
			continue
		}

		readyTasks = append(readyTasks, task)
	}

	return readyTasks
}

// countExecutingTasks 计算正在执行的任务数量
func (m *QueueLifecycleManager) countExecutingTasks(queue []DeviceTaskQueue) int {
	count := 0
	for _, task := range queue {
		if task.Status == DeviceQueueStatusExecuting {
			count++
		}
	}
	return count
}

// executeQueueItem 执行队列项
func (m *QueueLifecycleManager) executeQueueItem(ctx context.Context, item *DeviceTaskQueue) error {
	// 验证任务是否可以执行
	if err := m.validator.ValidateQueueItem(item); err != nil {
		return err
	}

	// 转换为执行状态
	return m.TransitionQueueItem(ctx, item, DeviceQueueStatusExecuting, nil)
}

// StartScheduler 启动调度器
func (m *QueueLifecycleManager) StartScheduler(ctx context.Context) error {
	m.mutex.Lock()
	if m.schedulerRunning {
		m.mutex.Unlock()
		return errors.New("调度器已在运行")
	}
	m.schedulerRunning = true
	m.mutex.Unlock()

	go m.runScheduler(ctx)
	return nil
}

// StopScheduler 停止调度器
func (m *QueueLifecycleManager) StopScheduler() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.schedulerRunning {
		return errors.New("调度器未运行")
	}

	m.schedulerStop <- true
	m.schedulerRunning = false
	return nil
}

// runScheduler 运行调度器
func (m *QueueLifecycleManager) runScheduler(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(consts.QueueMonitorInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.schedulerStop:
			return
		case <-ticker.C:
			// TODO: 在这里添加具体的调度逻辑
			// 例如：从数据库获取所有设备队列，然后调度执行
			m.performScheduling(ctx)
		}
	}
}

// performScheduling 执行调度逻辑
func (m *QueueLifecycleManager) performScheduling(ctx context.Context) {
	// TODO: 实现具体的调度逻辑
	// 1. 获取所有活跃设备的队列
	// 2. 为每个设备调度队列执行
	// 3. 处理超时任务
	// 4. 清理已完成的任务
}

// ProcessTaskTimeout 处理任务超时
func (m *QueueLifecycleManager) ProcessTaskTimeout(ctx context.Context, item *DeviceTaskQueue) error {
	if item.Status != DeviceQueueStatusExecuting {
		return nil // 只处理执行中的任务
	}

	if item.ActualStartTime == nil {
		return nil // 没有开始时间，无法判断超时
	}

	// 检查是否超时
	elapsed := time.Since(*item.ActualStartTime)
	timeout := time.Duration(item.TimeoutSeconds) * time.Second

	if elapsed <= timeout {
		return nil // 未超时
	}

	// 处理超时
	reason := fmt.Sprintf("任务执行超时: 执行时间%v超过限制%v", elapsed, timeout)

	if m.business.ShouldRetryTask(item) {
		// 增加重试次数
		item.CurrentRetry++
		// 重新入队
		return m.TransitionQueueItem(ctx, item, DeviceQueueStatusQueued, &reason)
	} else {
		// 标记为失败
		return m.TransitionQueueItem(ctx, item, DeviceQueueStatusFailed, &reason)
	}
}

// OptimizeQueue 优化队列
func (m *QueueLifecycleManager) OptimizeQueue(ctx context.Context, deviceID int64, queue []DeviceTaskQueue) ([]DeviceTaskQueue, error) {
	config := m.getDeviceQueueConfig(deviceID)
	if config == nil {
		return queue, errors.New("设备队列配置不存在")
	}

	// 使用业务逻辑优化队列顺序
	optimized := m.business.OptimizeQueueOrder(queue, config)

	return optimized, nil
}

// CalculateQueueMetrics 计算队列指标
func (m *QueueLifecycleManager) CalculateQueueMetrics(deviceID int64, queue []DeviceTaskQueue) *QueueMetrics {
	metrics := &QueueMetrics{
		DeviceID:    deviceID,
		LastUpdated: time.Now(),
	}

	if len(queue) == 0 {
		return metrics
	}

	// 统计任务状态
	for _, task := range queue {
		metrics.TotalTasks++
		switch task.Status {
		case DeviceQueueStatusQueued:
			metrics.QueuedTasks++
		case DeviceQueueStatusExecuting:
			metrics.ExecutingTasks++
		case DeviceQueueStatusCompleted:
			metrics.CompletedTasks++
		case DeviceQueueStatusFailed:
			metrics.FailedTasks++
		}
	}

	// 计算平均等待时间
	totalWaitTime := 0.0
	queuedCount := 0
	now := time.Now()

	for _, task := range queue {
		if task.Status == DeviceQueueStatusQueued {
			waitTime := now.Sub(task.CreatedAt).Hours()
			totalWaitTime += waitTime
			queuedCount++
		}
	}

	if queuedCount > 0 {
		metrics.AverageWaitTime = totalWaitTime / float64(queuedCount)
	}

	// 计算队列健康度
	businessMetrics := m.business.CalculateQueueMetrics(queue)
	if score, exists := businessMetrics["queue_health_score"]; exists {
		if scoreFloat, ok := score.(float64); ok {
			metrics.QueueHealthScore = scoreFloat
		}
	}

	return metrics
}

// BatchUpdateQueueItems 批量更新队列项
func (m *QueueLifecycleManager) BatchUpdateQueueItems(ctx context.Context, operation *BatchQueueOperation) error {
	// 验证批量操作请求
	if err := m.validator.ValidateBatchOperation(operation); err != nil {
		return fmt.Errorf("批量操作验证失败: %w", err)
	}

	// TODO: 这里需要与数据库交互，获取要操作的队列项
	// 由于没有repository层，这里只做逻辑验证

	// 记录批量操作历史
	for _, taskID := range operation.TaskIDs {
		history := m.business.GenerateOperationHistory(
			operation.DeviceID,
			"", // 需要从数据库获取设备ESN
			&taskID,
			operation.Operation,
			nil, // 批量操作暂时不记录操作人
			operation.Reason,
		)
		history.IsBatchOperation = true
		history.BatchID = operation.BatchID

		// TODO: 保存操作历史到数据库
	}

	return nil
}

// PauseQueue 暂停设备队列
func (m *QueueLifecycleManager) PauseQueue(ctx context.Context, deviceID int64, reason *string) error {
	// TODO: 实现队列暂停逻辑
	// 1. 获取设备的所有执行中任务
	// 2. 将它们转换为暂停状态
	// 3. 记录操作历史

	return nil
}

// ResumeQueue 恢复设备队列
func (m *QueueLifecycleManager) ResumeQueue(ctx context.Context, deviceID int64, reason *string) error {
	// TODO: 实现队列恢复逻辑
	// 1. 获取设备的所有暂停任务
	// 2. 将它们转换为排队状态
	// 3. 重新调度执行
	// 4. 记录操作历史

	return nil
}

// ClearCompletedTasks 清理已完成的任务
func (m *QueueLifecycleManager) ClearCompletedTasks(ctx context.Context, deviceID int64, olderThan time.Time) (int, error) {
	// TODO: 实现清理已完成任务的逻辑
	// 1. 查找指定时间之前完成的任务
	// 2. 移动到历史表或删除
	// 3. 返回清理的任务数量

	return 0, nil
}

// GetQueueStatus 获取队列状态摘要
func (m *QueueLifecycleManager) GetQueueStatus(deviceID int64) (map[string]interface{}, error) {
	// TODO: 实现获取队列状态的逻辑
	// 1. 获取设备队列
	// 2. 计算各种统计信息
	// 3. 返回状态摘要

	status := map[string]interface{}{
		"device_id":         deviceID,
		"scheduler_running": m.schedulerRunning,
		"last_updated":      time.Now(),
	}

	return status, nil
}

// 配置管理相关方法

// getDeviceQueueConfig 获取设备队列配置
func (m *QueueLifecycleManager) getDeviceQueueConfig(deviceID int64) *DeviceQueueConfig {
	m.configMutex.RLock()
	config, exists := m.configCache[deviceID]
	m.configMutex.RUnlock()

	if exists {
		return config
	}

	// TODO: 从数据库加载配置
	// 这里返回默认配置
	defaultConfig := &DeviceQueueConfig{
		DeviceID:           deviceID,
		MaxQueueSize:       consts.QueueMaxSizeDefault,
		MaxConcurrentTasks: consts.QueueConcurrentDefault,
		AutoStartTasks:     true,
		PriorityScheduling: true,
		SchedulingStrategy: SchedulingStrategyPriorityFirst,
		LoadBalancing:      true,
		Timezone:           "UTC",
		MaxCPUUsage:        80.0,
		MaxMemoryUsage:     80.0,
		MinFreeDisk:        1073741824, // 1GB
	}

	// 缓存配置
	m.updateConfigCache(deviceID, defaultConfig)

	return defaultConfig
}

// updateConfigCache 更新配置缓存
func (m *QueueLifecycleManager) updateConfigCache(deviceID int64, config *DeviceQueueConfig) {
	m.configMutex.Lock()
	defer m.configMutex.Unlock()

	m.configCache[deviceID] = config
}

// clearConfigCache 清理配置缓存
func (m *QueueLifecycleManager) clearConfigCache(deviceID int64) {
	m.configMutex.Lock()
	defer m.configMutex.Unlock()

	delete(m.configCache, deviceID)
}

// IsSchedulerRunning 检查调度器是否运行
func (m *QueueLifecycleManager) IsSchedulerRunning() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.schedulerRunning
}
