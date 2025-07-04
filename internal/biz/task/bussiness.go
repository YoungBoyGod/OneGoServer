package task

import (
	"context"
	"errors"
	"time"

	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
)

// TaskBusiness 任务业务逻辑层
type TaskBusiness struct {
	validator *TaskValidator
}

// NewTaskBusiness 创建任务业务逻辑实例
func NewTaskBusiness() *TaskBusiness {
	return &TaskBusiness{
		validator: NewTaskValidator(),
	}
}

// ValidateTask 验证任务（兼容性方法）
func (b *TaskBusiness) ValidateTask(ctx context.Context, task *Task) error {
	return b.validator.ValidateTaskForExecution(task)
}

// CalculateTaskScore 计算任务分数
func (b *TaskBusiness) CalculateTaskScore(task *Task) int {
	if task == nil {
		return 0
	}

	// 基础分数
	score := task.Priority * 10

	// 紧急任务加分
	if task.IsUrgent {
		score += 50
	}

	// 重试次数影响分数（重试越多分数越低）
	if task.RetryCount > 0 {
		score -= task.RetryCount * 5
	}

	// 超时任务加分
	if task.ExecuteTime != nil && task.ExecuteTime.Before(time.Now()) {
		score += 30
	}

	// 任务类型权重
	switch task.Type {
	case TaskTypeBackup:
		score += 20 // 备份任务优先级较高
	case TaskTypeSync:
		score += 15 // 同步任务中等优先级
	case TaskTypeMonitor:
		score += 10 // 监控任务较低优先级
	case TaskTypeCustom:
		score += 5 // 自定义任务最低优先级
	}

	// 确保分数为正数
	if score < 0 {
		score = 1
	}

	return score
}

// DetermineNextStatus 确定下一个状态
func (b *TaskBusiness) DetermineNextStatus(task *Task) string {
	if task == nil {
		return TaskStatusPending
	}

	switch task.Status {
	case TaskStatusPending:
		return TaskStatusQueued
	case TaskStatusQueued:
		return TaskStatusAssigning
	case TaskStatusAssigning:
		return TaskStatusAssigned
	case TaskStatusAssigned:
		return TaskStatusDispatching
	case TaskStatusDispatching:
		return TaskStatusRunning
	case TaskStatusRunning:
		// 根据执行结果确定
		if task.RetryCount >= task.MaxRetries {
			return TaskStatusFailed
		}
		return TaskStatusCompleted
	case TaskStatusFailed:
		if b.ShouldRetry(task) {
			return TaskStatusQueued
		}
		return TaskStatusFailed
	default:
		return task.Status
	}
}

// CanTransitionTo 检查是否可以转换到目标状态
func (b *TaskBusiness) CanTransitionTo(task *Task, targetStatus string) bool {
	if task == nil {
		return false
	}

	err := b.validator.ValidateStatusTransition(task.Status, targetStatus)
	return err == nil
}

// CalculatePriority 动态计算优先级
func (b *TaskBusiness) CalculatePriority(task *Task) int {
	if task == nil {
		return 5
	}

	priority := task.Priority

	// 根据等待时间调整优先级
	if task.CreatedAt.Before(time.Now().Add(-1 * time.Hour)) {
		priority += 1 // 等待超过1小时的任务优先级+1
	}
	if task.CreatedAt.Before(time.Now().Add(-6 * time.Hour)) {
		priority += 2 // 等待超过6小时的任务优先级+2
	}

	// 重试次数影响优先级
	if task.RetryCount > 0 && task.RetryCount < task.MaxRetries {
		priority -= task.RetryCount // 重试次数越多优先级越低
	}

	// 紧急任务最高优先级
	if task.IsUrgent {
		priority = 10
	}

	// 限制优先级范围
	if priority > 10 {
		priority = 10
	}
	if priority < 1 {
		priority = 1
	}

	return priority
}

// ShouldRetry 判断是否应该重试
func (b *TaskBusiness) ShouldRetry(task *Task) bool {
	if task == nil {
		return false
	}

	// 检查重试次数
	if task.RetryCount >= task.MaxRetries {
		return false
	}

	// 已取消或已完成的任务不重试
	if task.Status == TaskStatusCanceled || task.Status == TaskStatusCompleted {
		return false
	}

	// 只有失败的任务才考虑重试
	if task.Status != TaskStatusFailed {
		return false
	}

	// 检查错误类型（如果有特定的错误不允许重试）
	if task.ErrorMessage != nil {
		errorMsg := *task.ErrorMessage
		// 某些致命错误不重试
		fatalErrors := []string{
			"权限不足",
			"参数错误",
			"资源不存在",
			"配置错误",
		}
		for _, fatalError := range fatalErrors {
			if contains(errorMsg, fatalError) {
				return false
			}
		}
	}

	return true
}

// CalculateEstimatedDuration 计算预估执行时长
func (b *TaskBusiness) CalculateEstimatedDuration(task *Task) time.Duration {
	if task == nil {
		return 5 * time.Minute // 默认5分钟
	}

	// 基础时长根据任务类型
	baseDuration := map[string]time.Duration{
		TaskTypeBackup:  30 * time.Minute,
		TaskTypeSync:    10 * time.Minute,
		TaskTypeMonitor: 5 * time.Minute,
		TaskTypeCustom:  15 * time.Minute,
	}

	duration, exists := baseDuration[task.Type]
	if !exists {
		duration = 10 * time.Minute
	}

	// 根据参数调整
	if task.Parameters != nil {
		if size, ok := task.Parameters["data_size"]; ok {
			if sizeFloat, ok := size.(float64); ok && sizeFloat > 1000 {
				duration *= 2 // 大数据量任务时间翻倍
			}
		}
	}

	return duration
}

// GetTaskLifecycleSteps 获取任务生命周期步骤
func (b *TaskBusiness) GetTaskLifecycleSteps(task *Task) []string {
	if task == nil {
		return []string{}
	}

	steps := []string{
		TaskStatusPending,
		TaskStatusQueued,
		TaskStatusAssigning,
		TaskStatusAssigned,
		TaskStatusDispatching,
		TaskStatusRunning,
	}

	// 根据当前状态添加可能的结束状态
	switch task.Status {
	case TaskStatusRunning:
		if b.ShouldRetry(task) {
			steps = append(steps, TaskStatusCompleted, TaskStatusFailed)
		} else {
			steps = append(steps, TaskStatusCompleted, TaskStatusCanceled)
		}
	default:
		steps = append(steps, TaskStatusCompleted)
	}

	return steps
}

// BuildTaskExecution 构建任务执行记录
func (b *TaskBusiness) BuildTaskExecution(task *Task, executorType, executorID string) *TaskExecution {
	if task == nil {
		return nil
	}

	execution := &TaskExecution{
		TaskID:      task.ID,
		ExecutionID: utils.GenerateExecutionID(),
		Status:      ExecutionStatusStarted,
		StartTime:   time.Now(),
		ExecutorInfo: &JSONB{
			"executor_type": executorType,
			"executor_id":   executorID,
			"start_time":    time.Now(),
		},
	}

	// 如果任务关联了设备
	if task.DeviceID != nil && task.Device != nil {
		execution.DeviceESN = &task.Device.DeviceID
	}

	return execution
}

// ValidateTaskTransition 验证任务状态转换
func (b *TaskBusiness) ValidateTaskTransition(task *Task, targetStatus string) error {
	if task == nil {
		return errors.New("任务不能为空")
	}

	return b.validator.ValidateStatusTransition(task.Status, targetStatus)
}

// CalculateTaskDelay 计算任务延迟时间
func (b *TaskBusiness) CalculateTaskDelay(task *Task) time.Duration {
	if task == nil || task.ExecuteTime == nil {
		return 0
	}

	if task.ExecuteTime.Before(time.Now()) {
		return time.Since(*task.ExecuteTime)
	}

	return 0
}

// IsTaskOverdue 检查任务是否超时
func (b *TaskBusiness) IsTaskOverdue(task *Task) bool {
	if task == nil {
		return false
	}

	// 检查执行时间是否超时
	if task.ExecuteTime != nil && task.ExecuteTime.Before(time.Now().Add(-time.Duration(task.Timeout)*time.Second)) {
		return true
	}

	// 检查任务创建后是否超过了合理的等待时间
	maxWaitTime := 24 * time.Hour // 最大等待24小时
	if task.CreatedAt.Before(time.Now().Add(-maxWaitTime)) {
		return true
	}

	return false
}

// GetTaskPriorityLevel 获取任务优先级等级描述
func (b *TaskBusiness) GetTaskPriorityLevel(priority int) string {
	switch {
	case priority >= 9:
		return "紧急"
	case priority >= 7:
		return "高"
	case priority >= 5:
		return "中"
	case priority >= 3:
		return "低"
	default:
		return "很低"
	}
}

// 辅助函数
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			str[:len(substr)] == substr ||
			str[len(str)-len(substr):] == substr ||
			containsMiddle(str, substr))
}

func containsMiddle(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
