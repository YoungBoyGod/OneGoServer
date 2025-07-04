package queue

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/consts"
)

// QueueValidator 队列验证器
type QueueValidator struct{}

// NewQueueValidator 创建队列验证器实例
func NewQueueValidator() *QueueValidator {
	return &QueueValidator{}
}

// ValidateAddToQueue 验证添加任务到队列请求
func (v *QueueValidator) ValidateAddToQueue(req *DeviceTaskQueue) error {
	if req == nil {
		return errors.New("队列请求不能为空")
	}

	// 验证设备ID
	if req.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	// 验证设备ESN
	if strings.TrimSpace(req.DeviceESN) == "" {
		return errors.New("设备ESN不能为空")
	}
	if len(req.DeviceESN) > 100 {
		return errors.New("设备ESN不能超过100个字符")
	}

	// 验证任务ID
	if req.TaskID <= 0 {
		return errors.New("任务ID必须大于0")
	}

	// 验证队列优先级
	if req.QueuePriority < consts.QueuePriorityLowest || req.QueuePriority > consts.QueuePriorityHighest {
		return fmt.Errorf("队列优先级必须在%d-%d之间", consts.QueuePriorityLowest, consts.QueuePriorityHighest)
	}

	// 验证队列位置
	if req.QueuePosition < 0 {
		return errors.New("队列位置不能为负数")
	}

	// 验证重试配置
	if req.MaxRetryCount < consts.QueueRetryMin || req.MaxRetryCount > consts.QueueRetryMax {
		return fmt.Errorf("最大重试次数必须在%d-%d之间", consts.QueueRetryMin, consts.QueueRetryMax)
	}

	if req.CurrentRetry < 0 || req.CurrentRetry > req.MaxRetryCount {
		return errors.New("当前重试次数不能超过最大重试次数")
	}

	// 验证超时配置
	if req.TimeoutSeconds < consts.QueueTimeoutMin || req.TimeoutSeconds > consts.QueueTimeoutMax {
		return fmt.Errorf("超时时间必须在%d-%d秒之间", consts.QueueTimeoutMin, consts.QueueTimeoutMax)
	}

	// 验证预估开始时间（如果提供）
	if req.EstimatedStartTime != nil && req.EstimatedStartTime.Before(time.Now()) {
		return errors.New("预估开始时间不能早于当前时间")
	}

	// 验证预估执行时长
	if req.EstimatedDuration != nil && *req.EstimatedDuration <= 0 {
		return errors.New("预估执行时长必须大于0")
	}

	return nil
}

// ValidatePriorityChange 验证优先级变更请求
func (v *QueueValidator) ValidatePriorityChange(req *QueuePriorityChangeRequest) error {
	if req == nil {
		return errors.New("优先级变更请求不能为空")
	}

	if req.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	if req.TaskID <= 0 {
		return errors.New("任务ID必须大于0")
	}

	if req.NewPriority < consts.QueuePriorityLowest || req.NewPriority > consts.QueuePriorityHighest {
		return fmt.Errorf("新优先级必须在%d-%d之间", consts.QueuePriorityLowest, consts.QueuePriorityHighest)
	}

	if len(req.Reason) > 255 {
		return errors.New("变更原因不能超过255个字符")
	}

	return nil
}

// ValidatePositionChange 验证位置变更请求
func (v *QueueValidator) ValidatePositionChange(req *QueuePositionChangeRequest) error {
	if req == nil {
		return errors.New("位置变更请求不能为空")
	}

	if req.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	if req.TaskID <= 0 {
		return errors.New("任务ID必须大于0")
	}

	if req.NewPosition < 0 {
		return errors.New("新位置不能为负数")
	}

	if len(req.Reason) > 255 {
		return errors.New("变更原因不能超过255个字符")
	}

	return nil
}

// ValidateBatchOperation 验证批量操作请求
func (v *QueueValidator) ValidateBatchOperation(req *BatchQueueOperation) error {
	if req == nil {
		return errors.New("批量操作请求不能为空")
	}

	if req.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	if len(req.TaskIDs) == 0 {
		return errors.New("任务ID列表不能为空")
	}

	if len(req.TaskIDs) > 100 {
		return errors.New("批量操作最多支持100个任务")
	}

	// 验证任务ID
	for i, taskID := range req.TaskIDs {
		if taskID <= 0 {
			return fmt.Errorf("第%d个任务ID必须大于0", i+1)
		}
	}

	// 验证操作类型
	if !v.isValidOperation(req.Operation) {
		return fmt.Errorf("无效的操作类型: %s", req.Operation)
	}

	// 验证优先级变更
	if req.Operation == QueueOperationPriorityChange {
		if req.NewPriority == nil {
			return errors.New("优先级变更操作必须提供新优先级")
		}
		if *req.NewPriority < consts.QueuePriorityLowest || *req.NewPriority > consts.QueuePriorityHighest {
			return fmt.Errorf("新优先级必须在%d-%d之间", consts.QueuePriorityLowest, consts.QueuePriorityHighest)
		}
	}

	// 验证位置变更
	if req.Operation == QueueOperationPositionChange {
		if len(req.NewPositions) != len(req.TaskIDs) {
			return errors.New("位置变更操作中新位置数量必须与任务数量一致")
		}
		for i, pos := range req.NewPositions {
			if pos < 0 {
				return fmt.Errorf("第%d个新位置不能为负数", i+1)
			}
		}
	}

	if req.Reason != nil && len(*req.Reason) > 255 {
		return errors.New("操作原因不能超过255个字符")
	}

	return nil
}

// ValidateQueueConfig 验证队列配置
func (v *QueueValidator) ValidateQueueConfig(config *DeviceQueueConfig) error {
	if config == nil {
		return errors.New("队列配置不能为空")
	}

	if config.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	if strings.TrimSpace(config.DeviceESN) == "" {
		return errors.New("设备ESN不能为空")
	}

	// 验证队列大小
	if config.MaxQueueSize < 1 || config.MaxQueueSize > consts.QueueMaxSizeLimit {
		return fmt.Errorf("最大队列大小必须在1-%d之间", consts.QueueMaxSizeLimit)
	}

	// 验证并发任务数
	if config.MaxConcurrentTasks < 1 || config.MaxConcurrentTasks > consts.QueueConcurrentLimit {
		return fmt.Errorf("最大并发任务数必须在1-%d之间", consts.QueueConcurrentLimit)
	}

	// 验证调度策略
	if !v.isValidSchedulingStrategy(config.SchedulingStrategy) {
		return fmt.Errorf("无效的调度策略: %s", config.SchedulingStrategy)
	}

	// 验证时区
	if config.Timezone != "" {
		if _, err := time.LoadLocation(config.Timezone); err != nil {
			return fmt.Errorf("无效的时区: %s", config.Timezone)
		}
	}

	// 验证资源限制
	if config.MaxCPUUsage < 0 || config.MaxCPUUsage > 100 {
		return errors.New("最大CPU使用率必须在0-100之间")
	}

	if config.MaxMemoryUsage < 0 || config.MaxMemoryUsage > 100 {
		return errors.New("最大内存使用率必须在0-100之间")
	}

	if config.MinFreeDisk < 0 {
		return errors.New("最小剩余磁盘空间不能为负数")
	}

	// 验证通知配置
	if config.NotificationWebhook != nil && len(*config.NotificationWebhook) > 255 {
		return errors.New("通知webhook URL不能超过255个字符")
	}

	return nil
}

// ValidateFilter 验证查询过滤器
func (v *QueueValidator) ValidateFilter(filter *DeviceQueueFilter) error {
	if filter == nil {
		return nil
	}

	// 验证设备ID
	if filter.DeviceID != nil && *filter.DeviceID <= 0 {
		return errors.New("设备ID必须大于0")
	}

	// 验证任务ID
	if filter.TaskID != nil && *filter.TaskID <= 0 {
		return errors.New("任务ID必须大于0")
	}

	// 验证状态过滤
	for _, status := range filter.Status {
		if !v.isValidQueueStatus(status) {
			return fmt.Errorf("无效的队列状态: %s", status)
		}
	}

	// 验证优先级
	if filter.Priority != nil {
		if *filter.Priority < consts.QueuePriorityLowest || *filter.Priority > consts.QueuePriorityHighest {
			return fmt.Errorf("优先级过滤条件必须在%d-%d之间", consts.QueuePriorityLowest, consts.QueuePriorityHighest)
		}
	}

	// 验证时间范围
	if filter.CreatedAfter != nil && filter.CreatedBefore != nil {
		if filter.CreatedAfter.After(*filter.CreatedBefore) {
			return errors.New("创建开始时间不能晚于结束时间")
		}
	}

	return nil
}

// ValidateStatusTransition 验证状态转换
func (v *QueueValidator) ValidateStatusTransition(from, to string) error {
	if !v.isValidQueueStatus(from) {
		return fmt.Errorf("无效的源状态: %s", from)
	}
	if !v.isValidQueueStatus(to) {
		return fmt.Errorf("无效的目标状态: %s", to)
	}

	// 定义允许的状态转换
	allowedTransitions := map[string][]string{
		DeviceQueueStatusQueued: {
			DeviceQueueStatusExecuting,
			DeviceQueueStatusPaused,
			DeviceQueueStatusCanceled,
		},
		DeviceQueueStatusExecuting: {
			DeviceQueueStatusCompleted,
			DeviceQueueStatusFailed,
			DeviceQueueStatusPaused,
			DeviceQueueStatusCanceled,
		},
		DeviceQueueStatusPaused: {
			DeviceQueueStatusQueued,
			DeviceQueueStatusExecuting,
			DeviceQueueStatusCanceled,
		},
		DeviceQueueStatusCompleted: {
			// 已完成的任务通常不允许状态转换
		},
		DeviceQueueStatusFailed: {
			DeviceQueueStatusQueued, // 重试
			DeviceQueueStatusCanceled,
		},
		DeviceQueueStatusCanceled: {
			// 已取消的任务通常不允许状态转换
		},
	}

	validTargets, exists := allowedTransitions[from]
	if !exists {
		return fmt.Errorf("状态 %s 不支持任何转换", from)
	}

	for _, validTarget := range validTargets {
		if validTarget == to {
			return nil
		}
	}

	return fmt.Errorf("不允许从状态 %s 转换到 %s", from, to)
}

// 辅助验证方法

// isValidOperation 检查操作类型是否有效
func (v *QueueValidator) isValidOperation(operation string) bool {
	validOperations := []string{
		QueueOperationAdd,
		QueueOperationRemove,
		QueueOperationPriorityChange,
		QueueOperationPositionChange,
		QueueOperationStart,
		QueueOperationPause,
		QueueOperationResume,
		QueueOperationCancel,
	}
	for _, valid := range validOperations {
		if valid == operation {
			return true
		}
	}
	return false
}

// isValidSchedulingStrategy 检查调度策略是否有效
func (v *QueueValidator) isValidSchedulingStrategy(strategy string) bool {
	validStrategies := []string{
		SchedulingStrategyPriorityFirst,
		SchedulingStrategyFIFO,
		SchedulingStrategyLIFO,
		SchedulingStrategyWeighted,
	}
	for _, valid := range validStrategies {
		if valid == strategy {
			return true
		}
	}
	return false
}

// isValidQueueStatus 检查队列状态是否有效
func (v *QueueValidator) isValidQueueStatus(status string) bool {
	validStatuses := []string{
		DeviceQueueStatusQueued,
		DeviceQueueStatusExecuting,
		DeviceQueueStatusPaused,
		DeviceQueueStatusCompleted,
		DeviceQueueStatusFailed,
		DeviceQueueStatusCanceled,
	}
	for _, valid := range validStatuses {
		if valid == status {
			return true
		}
	}
	return false
}

// isValidOperationSource 检查操作来源是否有效
func (v *QueueValidator) isValidOperationSource(source string) bool {
	validSources := []string{
		OperationSourceManual,
		OperationSourceSystem,
		OperationSourceAPI,
		OperationSourceScheduler,
	}
	for _, valid := range validSources {
		if valid == source {
			return true
		}
	}
	return false
}

// ValidateQueueItem 验证队列项是否可以执行
func (v *QueueValidator) ValidateQueueItem(item *DeviceTaskQueue) error {
	if item == nil {
		return errors.New("队列项不能为空")
	}

	if item.Status != DeviceQueueStatusQueued {
		return fmt.Errorf("只有排队中的任务可以执行，当前状态: %s", item.Status)
	}

	if item.EstimatedStartTime != nil && item.EstimatedStartTime.After(time.Now()) {
		return errors.New("任务预估开始时间未到")
	}

	if item.CurrentRetry >= item.MaxRetryCount {
		return errors.New("任务已达到最大重试次数")
	}

	return nil
}

// ValidateWorkingHours 验证是否在工作时间内
func (v *QueueValidator) ValidateWorkingHours(config *DeviceQueueConfig) error {
	if config == nil {
		return nil // 没有配置时间窗口，允许全天执行
	}

	if config.WorkStartTime == nil || config.WorkEndTime == nil {
		return nil // 没有设置工作时间窗口
	}

	now := time.Now()
	if config.Timezone != "" {
		if loc, err := time.LoadLocation(config.Timezone); err == nil {
			now = now.In(loc)
		}
	}

	currentTime := now.Format("15:04:05")
	startTime := config.WorkStartTime.Format("15:04:05")
	endTime := config.WorkEndTime.Format("15:04:05")

	if startTime <= endTime {
		// 同一天内的时间窗口
		if currentTime < startTime || currentTime > endTime {
			return fmt.Errorf("当前时间 %s 不在工作时间窗口 %s-%s 内", currentTime, startTime, endTime)
		}
	} else {
		// 跨天的时间窗口
		if currentTime < startTime && currentTime > endTime {
			return fmt.Errorf("当前时间 %s 不在工作时间窗口 %s-%s 内", currentTime, startTime, endTime)
		}
	}

	return nil
}
