package task

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// TaskLifecycleManager 任务生命周期管理器
type TaskLifecycleManager struct {
	business  *TaskBusiness
	validator *TaskValidator
	mutex     sync.RWMutex

	// 状态监听器
	stateListeners map[string][]StateChangeListener

	// 任务依赖关系
	dependencies map[int64][]int64 // taskID -> 依赖的taskID列表
	dependents   map[int64][]int64 // taskID -> 依赖它的taskID列表
}

// StateChangeListener 状态变更监听器
type StateChangeListener func(ctx context.Context, task *Task, oldStatus, newStatus string) error

// TaskDependency 任务依赖关系
type TaskDependency struct {
	TaskID          int64     `json:"task_id"`
	DependentTaskID int64     `json:"dependent_task_id"`
	DependencyType  string    `json:"dependency_type"` // "before", "after", "parallel"
	CreatedAt       time.Time `json:"created_at"`
}

// 依赖类型常量
const (
	DependencyTypeBefore   = "before"   // 前置依赖
	DependencyTypeAfter    = "after"    // 后置依赖
	DependencyTypeParallel = "parallel" // 并行依赖
)

// NewTaskLifecycleManager 创建任务生命周期管理器
func NewTaskLifecycleManager() *TaskLifecycleManager {
	return &TaskLifecycleManager{
		business:       NewTaskBusiness(),
		validator:      NewTaskValidator(),
		stateListeners: make(map[string][]StateChangeListener),
		dependencies:   make(map[int64][]int64),
		dependents:     make(map[int64][]int64),
	}
}

// TransitionTask 转换任务状态
func (m *TaskLifecycleManager) TransitionTask(ctx context.Context, task *Task, targetStatus string) error {
	if task == nil {
		return errors.New("任务不能为空")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldStatus := task.Status

	// 验证状态转换
	if err := m.validator.ValidateStatusTransition(oldStatus, targetStatus); err != nil {
		return fmt.Errorf("状态转换验证失败: %w", err)
	}

	// 检查依赖关系
	if err := m.checkDependencies(ctx, task, targetStatus); err != nil {
		return fmt.Errorf("依赖检查失败: %w", err)
	}

	// 执行状态转换
	task.Status = targetStatus
	task.UpdatedAt = time.Now()

	// 触发状态变更监听器
	if err := m.triggerStateListeners(ctx, task, oldStatus, targetStatus); err != nil {
		// 回滚状态
		task.Status = oldStatus
		return fmt.Errorf("状态监听器执行失败: %w", err)
	}

	// 处理依赖任务
	if err := m.processDependentTasks(ctx, task, targetStatus); err != nil {
		// 记录错误但不回滚，因为当前任务状态已经成功转换
		// 这里可以记录日志或添加到错误队列中稍后处理
		return fmt.Errorf("处理依赖任务失败: %w", err)
	}

	return nil
}

// AddDependency 添加任务依赖关系
func (m *TaskLifecycleManager) AddDependency(taskID, dependentTaskID int64, dependencyType string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 验证依赖类型
	if !isValidDependencyType(dependencyType) {
		return fmt.Errorf("无效的依赖类型: %s", dependencyType)
	}

	// 防止循环依赖
	if m.wouldCreateCycle(taskID, dependentTaskID) {
		return errors.New("添加依赖会创建循环依赖")
	}

	// 添加依赖关系
	if m.dependencies[dependentTaskID] == nil {
		m.dependencies[dependentTaskID] = make([]int64, 0)
	}
	m.dependencies[dependentTaskID] = append(m.dependencies[dependentTaskID], taskID)

	if m.dependents[taskID] == nil {
		m.dependents[taskID] = make([]int64, 0)
	}
	m.dependents[taskID] = append(m.dependents[taskID], dependentTaskID)

	return nil
}

// RemoveDependency 移除任务依赖关系
func (m *TaskLifecycleManager) RemoveDependency(taskID, dependentTaskID int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 从dependencies中移除
	if deps := m.dependencies[dependentTaskID]; deps != nil {
		m.dependencies[dependentTaskID] = removeFromSlice(deps, taskID)
	}

	// 从dependents中移除
	if deps := m.dependents[taskID]; deps != nil {
		m.dependents[taskID] = removeFromSlice(deps, dependentTaskID)
	}

	return nil
}

// GetDependencies 获取任务的依赖列表
func (m *TaskLifecycleManager) GetDependencies(taskID int64) []int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if deps := m.dependencies[taskID]; deps != nil {
		result := make([]int64, len(deps))
		copy(result, deps)
		return result
	}
	return []int64{}
}

// GetDependents 获取依赖此任务的任务列表
func (m *TaskLifecycleManager) GetDependents(taskID int64) []int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if deps := m.dependents[taskID]; deps != nil {
		result := make([]int64, len(deps))
		copy(result, deps)
		return result
	}
	return []int64{}
}

// CanExecute 检查任务是否可以执行（所有依赖都已完成）
func (m *TaskLifecycleManager) CanExecute(ctx context.Context, task *Task, getTaskFunc func(int64) (*Task, error)) (bool, error) {
	if task == nil {
		return false, errors.New("任务不能为空")
	}

	dependencies := m.GetDependencies(task.ID)
	if len(dependencies) == 0 {
		return true, nil // 没有依赖，可以执行
	}

	// 检查所有依赖任务的状态
	for _, depTaskID := range dependencies {
		depTask, err := getTaskFunc(depTaskID)
		if err != nil {
			return false, fmt.Errorf("获取依赖任务 %d 失败: %w", depTaskID, err)
		}

		if depTask == nil {
			return false, fmt.Errorf("依赖任务 %d 不存在", depTaskID)
		}

		// 依赖任务必须已完成
		if depTask.Status != TaskStatusCompleted {
			return false, nil
		}
	}

	return true, nil
}

// AddStateListener 添加状态变更监听器
func (m *TaskLifecycleManager) AddStateListener(status string, listener StateChangeListener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.stateListeners[status] == nil {
		m.stateListeners[status] = make([]StateChangeListener, 0)
	}
	m.stateListeners[status] = append(m.stateListeners[status], listener)
}

// GetTaskLifecycleInfo 获取任务生命周期信息
func (m *TaskLifecycleManager) GetTaskLifecycleInfo(task *Task) *TaskLifecycleInfo {
	if task == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return &TaskLifecycleInfo{
		TaskID:           task.ID,
		CurrentStatus:    task.Status,
		PossibleStatuses: m.getPossibleNextStatuses(task),
		Dependencies:     m.GetDependencies(task.ID),
		Dependents:       m.GetDependents(task.ID),
		CanExecute:       len(m.GetDependencies(task.ID)) == 0, // 简化检查
		EstimatedSteps:   m.business.GetTaskLifecycleSteps(task),
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}

// TaskLifecycleInfo 任务生命周期信息
type TaskLifecycleInfo struct {
	TaskID           int64     `json:"task_id"`
	CurrentStatus    string    `json:"current_status"`
	PossibleStatuses []string  `json:"possible_statuses"`
	Dependencies     []int64   `json:"dependencies"`
	Dependents       []int64   `json:"dependents"`
	CanExecute       bool      `json:"can_execute"`
	EstimatedSteps   []string  `json:"estimated_steps"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// 私有方法

// checkDependencies 检查依赖关系
func (m *TaskLifecycleManager) checkDependencies(ctx context.Context, task *Task, targetStatus string) error {
	// 这里可以添加更复杂的依赖检查逻辑
	// 例如：检查前置任务是否完成、并行任务是否可以同时运行等
	return nil
}

// triggerStateListeners 触发状态监听器
func (m *TaskLifecycleManager) triggerStateListeners(ctx context.Context, task *Task, oldStatus, newStatus string) error {
	// 触发新状态的监听器
	if listeners := m.stateListeners[newStatus]; listeners != nil {
		for _, listener := range listeners {
			if err := listener(ctx, task, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	// 触发通用监听器（如果有的话）
	if listeners := m.stateListeners["*"]; listeners != nil {
		for _, listener := range listeners {
			if err := listener(ctx, task, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

// processDependentTasks 处理依赖任务
func (m *TaskLifecycleManager) processDependentTasks(ctx context.Context, task *Task, newStatus string) error {
	// 当任务完成时，检查是否可以触发依赖它的任务
	if newStatus == TaskStatusCompleted {
		dependents := m.GetDependents(task.ID)
		for _, depTaskID := range dependents {
			// 这里可以发送信号通知依赖任务检查是否可以开始执行
			// 具体实现可能需要与任务调度器集成
			_ = depTaskID // 暂时忽略未使用的变量，后续与调度器集成时会使用
		}
	}

	return nil
}

// wouldCreateCycle 检查是否会创建循环依赖
func (m *TaskLifecycleManager) wouldCreateCycle(taskID, dependentTaskID int64) bool {
	visited := make(map[int64]bool)
	return m.hasCycleDFS(dependentTaskID, taskID, visited)
}

// hasCycleDFS 深度优先搜索检查循环依赖
func (m *TaskLifecycleManager) hasCycleDFS(current, target int64, visited map[int64]bool) bool {
	if current == target {
		return true
	}

	if visited[current] {
		return false
	}

	visited[current] = true

	if deps := m.dependents[current]; deps != nil {
		for _, dep := range deps {
			if m.hasCycleDFS(dep, target, visited) {
				return true
			}
		}
	}

	return false
}

// getPossibleNextStatuses 获取可能的下一步状态
func (m *TaskLifecycleManager) getPossibleNextStatuses(task *Task) []string {
	// 基于当前状态返回可能的下一步状态
	switch task.Status {
	case TaskStatusPending:
		return []string{TaskStatusQueued, TaskStatusCanceled}
	case TaskStatusQueued:
		return []string{TaskStatusAssigning, TaskStatusCanceled}
	case TaskStatusAssigning:
		return []string{TaskStatusAssigned, TaskStatusFailed, TaskStatusQueued}
	case TaskStatusAssigned:
		return []string{TaskStatusDispatching, TaskStatusCanceled, TaskStatusQueued}
	case TaskStatusDispatching:
		return []string{TaskStatusRunning, TaskStatusFailed}
	case TaskStatusRunning:
		return []string{TaskStatusCompleted, TaskStatusFailed, TaskStatusCanceled}
	case TaskStatusFailed:
		if m.business.ShouldRetry(task) {
			return []string{TaskStatusQueued, TaskStatusCanceled}
		}
		return []string{TaskStatusCanceled}
	default:
		return []string{}
	}
}

// 辅助函数

// isValidDependencyType 验证依赖类型
func isValidDependencyType(dependencyType string) bool {
	validTypes := []string{
		DependencyTypeBefore,
		DependencyTypeAfter,
		DependencyTypeParallel,
	}
	for _, valid := range validTypes {
		if valid == dependencyType {
			return true
		}
	}
	return false
}

// removeFromSlice 从切片中移除元素
func removeFromSlice(slice []int64, element int64) []int64 {
	result := make([]int64, 0, len(slice))
	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}
	return result
}
