package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// TaskValidator 任务验证器
type TaskValidator struct{}

// NewTaskValidator 创建任务验证器实例
func NewTaskValidator() *TaskValidator {
	return &TaskValidator{}
}

// ValidateCreate 验证创建任务请求
func (v *TaskValidator) ValidateCreate(req *TaskCreateRequest) error {
	if req == nil {
		return errors.New("创建请求不能为空")
	}

	// 必填字段验证
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("任务名称不能为空")
	}
	if len(req.Name) > 255 {
		return errors.New("任务名称不能超过255个字符")
	}

	if strings.TrimSpace(req.Type) == "" {
		return errors.New("任务类型不能为空")
	}
	if !v.isValidTaskType(req.Type) {
		return fmt.Errorf("无效的任务类型: %s", req.Type)
	}

	// 可选字段验证
	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 10 {
			return errors.New("任务优先级必须在1-10之间")
		}
	}

	if req.ExecuteTime != nil && req.ExecuteTime.Before(time.Now()) {
		return errors.New("执行时间不能早于当前时间")
	}

	if req.Timeout != nil && *req.Timeout <= 0 {
		return errors.New("超时时间必须大于0")
	}

	if req.MaxRetries != nil && (*req.MaxRetries < 0 || *req.MaxRetries > 10) {
		return errors.New("最大重试次数必须在0-10之间")
	}

	if req.ExecutorType != nil && !v.isValidExecutorType(*req.ExecutorType) {
		return fmt.Errorf("无效的执行器类型: %s", *req.ExecutorType)
	}

	return nil
}

// ValidateUpdate 验证更新任务请求
func (v *TaskValidator) ValidateUpdate(req *TaskUpdateRequest) error {
	if req == nil {
		return errors.New("更新请求不能为空")
	}

	// 至少要有一个字段需要更新
	if v.isEmptyUpdateRequest(req) {
		return errors.New("至少需要更新一个字段")
	}

	// 验证具体字段
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return errors.New("任务名称不能为空")
		}
		if len(*req.Name) > 255 {
			return errors.New("任务名称不能超过255个字符")
		}
	}

	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 10 {
			return errors.New("任务优先级必须在1-10之间")
		}
	}

	if req.ExecuteTime != nil && req.ExecuteTime.Before(time.Now()) {
		return errors.New("执行时间不能早于当前时间")
	}

	if req.Timeout != nil && *req.Timeout <= 0 {
		return errors.New("超时时间必须大于0")
	}

	if req.MaxRetries != nil && (*req.MaxRetries < 0 || *req.MaxRetries > 10) {
		return errors.New("最大重试次数必须在0-10之间")
	}

	if req.ExecutorType != nil && !v.isValidExecutorType(*req.ExecutorType) {
		return fmt.Errorf("无效的执行器类型: %s", *req.ExecutorType)
	}

	return nil
}

// ValidateExecute 验证执行任务请求
func (v *TaskValidator) ValidateExecute(req *TaskExecuteRequest) error {
	if req == nil {
		return errors.New("执行请求不能为空")
	}

	if req.ExecutorType != nil && !v.isValidExecutorType(*req.ExecutorType) {
		return fmt.Errorf("无效的执行器类型: %s", *req.ExecutorType)
	}

	if req.ExecutorID != nil && strings.TrimSpace(*req.ExecutorID) == "" {
		return errors.New("执行器ID不能为空字符串")
	}

	return nil
}

// ValidateStatusTransition 验证状态转换
func (v *TaskValidator) ValidateStatusTransition(from, to string) error {
	if !v.isValidTaskStatus(from) {
		return fmt.Errorf("无效的源状态: %s", from)
	}
	if !v.isValidTaskStatus(to) {
		return fmt.Errorf("无效的目标状态: %s", to)
	}

	// 定义允许的状态转换
	allowedTransitions := map[string][]string{
		TaskStatusPending: {
			TaskStatusQueued,
			TaskStatusCanceled,
		},
		TaskStatusQueued: {
			TaskStatusAssigning,
			TaskStatusCanceled,
		},
		TaskStatusAssigning: {
			TaskStatusAssigned,
			TaskStatusFailed,
			TaskStatusQueued, // 重新入队
		},
		TaskStatusAssigned: {
			TaskStatusDispatching,
			TaskStatusCanceled,
			TaskStatusQueued, // 重新分配
		},
		TaskStatusDispatching: {
			TaskStatusRunning,
			TaskStatusFailed,
		},
		TaskStatusRunning: {
			TaskStatusCompleted,
			TaskStatusFailed,
			TaskStatusCanceled,
		},
		TaskStatusCompleted: {
			// 已完成的任务通常不允许状态转换
		},
		TaskStatusFailed: {
			TaskStatusQueued, // 重试
			TaskStatusCanceled,
		},
		TaskStatusCanceled: {
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

// ValidateTaskForExecution 验证任务是否可以执行
func (v *TaskValidator) ValidateTaskForExecution(task *Task) error {
	if task == nil {
		return errors.New("任务不能为空")
	}

	if task.Status != TaskStatusAssigned {
		return fmt.Errorf("只有已分配的任务可以执行，当前状态: %s", task.Status)
	}

	if task.ExecuteTime != nil && task.ExecuteTime.After(time.Now()) {
		return errors.New("任务执行时间未到")
	}

	return nil
}

// 辅助方法

// isValidTaskType 检查任务类型是否有效
func (v *TaskValidator) isValidTaskType(taskType string) bool {
	validTypes := []string{
		TaskTypeBackup,
		TaskTypeSync,
		TaskTypeMonitor,
		TaskTypeCustom,
	}
	for _, valid := range validTypes {
		if valid == taskType {
			return true
		}
	}
	return false
}

// isValidExecutorType 检查执行器类型是否有效
func (v *TaskValidator) isValidExecutorType(executorType string) bool {
	validTypes := []string{
		ExecutorTypeLocal,
		ExecutorTypeRemote,
		ExecutorTypeContainer,
		ExecutorTypeLambda,
	}
	for _, valid := range validTypes {
		if valid == executorType {
			return true
		}
	}
	return false
}

// isValidTaskStatus 检查任务状态是否有效
func (v *TaskValidator) isValidTaskStatus(status string) bool {
	validStatuses := []string{
		TaskStatusPending,
		TaskStatusQueued,
		TaskStatusAssigning,
		TaskStatusAssigned,
		TaskStatusDispatching,
		TaskStatusRunning,
		TaskStatusCompleted,
		TaskStatusFailed,
		TaskStatusCanceled,
	}
	for _, valid := range validStatuses {
		if valid == status {
			return true
		}
	}
	return false
}

// isEmptyUpdateRequest 检查更新请求是否为空
func (v *TaskValidator) isEmptyUpdateRequest(req *TaskUpdateRequest) bool {
	return req.Name == nil &&
		req.Description == nil &&
		req.Priority == nil &&
		req.ExecuteTime == nil &&
		req.Timeout == nil &&
		req.MaxRetries == nil &&
		req.Parameters == nil &&
		req.ExecutorType == nil &&
		req.ExecutorID == nil &&
		req.DeviceID == nil
}
