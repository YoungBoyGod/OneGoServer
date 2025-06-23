package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"learngo0619/internal/server/models"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks       map[string]*models.Task       // 任务存储
	taskResults map[string]*models.TaskResult // 任务执行结果存储
	mu          sync.RWMutex                  // 读写锁
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:       make(map[string]*models.Task),
		taskResults: make(map[string]*models.TaskResult),
	}
}

// generateID 生成唯一ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreateTask 创建新任务
func (tm *TaskManager) CreateTask(req *models.CreateTaskRequest, createdBy string) (*models.Task, error) {
	// 验证脚本类型
	if !req.ScriptType.IsValid() {
		return nil, fmt.Errorf("无效的脚本类型: %s", req.ScriptType)
	}

	// 设置默认值
	if req.Priority <= 0 {
		req.Priority = 5 // 默认优先级
	}
	if req.Timeout <= 0 {
		req.Timeout = 300 // 默认5分钟超时
	}

	now := time.Now()
	task := &models.Task{
		ID:            generateID(),
		Name:          req.Name,
		Description:   req.Description,
		ScriptType:    req.ScriptType,
		ScriptContent: req.ScriptContent,
		TargetClients: req.TargetClients,
		Priority:      req.Priority,
		Timeout:       req.Timeout,
		Status:        models.TaskStatusPending,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedAt:     now,
		ScheduledAt:   req.ScheduledAt,
	}

	tm.mu.Lock()
	tm.tasks[task.ID] = task
	tm.mu.Unlock()

	return task, nil
}

// GetTask 根据ID获取任务
func (tm *TaskManager) GetTask(taskID string) (*models.Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不存在: %s", taskID)
	}

	return task, nil
}

// GetTasksForClient 获取指定客户端的待执行任务
func (tm *TaskManager) GetTasksForClient(clientID string) ([]*models.Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var availableTasks []*models.Task
	now := time.Now()

	for _, task := range tm.tasks {
		// 只返回待执行状态的任务
		if task.Status != models.TaskStatusPending {
			continue
		}

		// 检查是否已到计划执行时间
		if task.ScheduledAt != nil && task.ScheduledAt.After(now) {
			continue
		}

		// 检查目标客户端
		if len(task.TargetClients) == 0 {
			// 空目标列表表示所有客户端
			availableTasks = append(availableTasks, task)
		} else {
			// 检查是否在目标客户端列表中
			for _, targetID := range task.TargetClients {
				if targetID == clientID {
					availableTasks = append(availableTasks, task)
					break
				}
			}
		}
	}

	// 按优先级排序（优先级高的在前）
	sort.Slice(availableTasks, func(i, j int) bool {
		return availableTasks[i].Priority > availableTasks[j].Priority
	})

	return availableTasks, nil
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() ([]*models.Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}

	// 按创建时间排序（最新的在前）
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})

	return tasks, nil
}

// UpdateTaskStatus 更新任务状态
func (tm *TaskManager) UpdateTaskStatus(taskID string, req *models.UpdateTaskStatusRequest) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	// 验证状态转换的合法性
	if !req.Status.IsValid() {
		return fmt.Errorf("无效的任务状态: %s", req.Status)
	}

	// 更新任务状态
	task.Status = req.Status
	task.UpdatedAt = time.Now()

	return nil
}

// RecordTaskResult 记录任务执行结果
func (tm *TaskManager) RecordTaskResult(taskID, clientID string, req *models.UpdateTaskStatusRequest) (*models.TaskResult, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 检查任务是否存在
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不存在: %s", taskID)
	}

	// 计算执行时间
	var duration int64
	if req.StartTime != nil && req.EndTime != nil {
		duration = req.EndTime.Sub(*req.StartTime).Milliseconds()
	}

	// 创建执行结果记录
	result := &models.TaskResult{
		ID:       generateID(),
		TaskID:   taskID,
		ClientID: clientID,
		Status:   req.Status,
		Output:   req.Output,
		Error:    req.Error,
		Duration: duration,
	}

	if req.StartTime != nil {
		result.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		result.EndTime = *req.EndTime
	}

	// 存储执行结果
	tm.taskResults[result.ID] = result

	// 更新任务状态
	task.Status = req.Status
	task.UpdatedAt = time.Now()

	return result, nil
}

// GetTaskResults 获取任务的执行结果
func (tm *TaskManager) GetTaskResults(taskID string) ([]*models.TaskResult, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var results []*models.TaskResult
	for _, result := range tm.taskResults {
		if result.TaskID == taskID {
			results = append(results, result)
		}
	}

	// 按执行时间排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.After(results[j].StartTime)
	})

	return results, nil
}

// GetClientTaskResults 获取客户端的任务执行结果
func (tm *TaskManager) GetClientTaskResults(clientID string) ([]*models.TaskResult, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	var results []*models.TaskResult
	for _, result := range tm.taskResults {
		if result.ClientID == clientID {
			results = append(results, result)
		}
	}

	// 按执行时间排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime.After(results[j].StartTime)
	})

	return results, nil
}

// DeleteTask 删除任务
func (tm *TaskManager) DeleteTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tasks[taskID]; !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}

	delete(tm.tasks, taskID)

	// 同时删除相关的执行结果
	for id, result := range tm.taskResults {
		if result.TaskID == taskID {
			delete(tm.taskResults, id)
		}
	}

	return nil
}

// GetTaskStats 获取任务统计信息
func (tm *TaskManager) GetTaskStats() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	stats := map[string]interface{}{
		"total_tasks":     len(tm.tasks),
		"total_results":   len(tm.taskResults),
		"pending_tasks":   0,
		"running_tasks":   0,
		"completed_tasks": 0,
		"failed_tasks":    0,
		"timeout_tasks":   0,
	}

	// 统计各状态的任务数量
	for _, task := range tm.tasks {
		switch task.Status {
		case models.TaskStatusPending:
			stats["pending_tasks"] = stats["pending_tasks"].(int) + 1
		case models.TaskStatusRunning:
			stats["running_tasks"] = stats["running_tasks"].(int) + 1
		case models.TaskStatusCompleted:
			stats["completed_tasks"] = stats["completed_tasks"].(int) + 1
		case models.TaskStatusFailed:
			stats["failed_tasks"] = stats["failed_tasks"].(int) + 1
		case models.TaskStatusTimeout:
			stats["timeout_tasks"] = stats["timeout_tasks"].(int) + 1
		}
	}

	return stats
}

// 全局任务管理器实例
var globalTaskManager *TaskManager
var taskManagerOnce sync.Once

// GetTaskManager 获取全局任务管理器实例
func GetTaskManager() *TaskManager {
	taskManagerOnce.Do(func() {
		globalTaskManager = NewTaskManager()
	})
	return globalTaskManager
}
