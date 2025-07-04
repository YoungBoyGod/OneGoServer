package service

import (
	"context"
	"fmt"
	"time"

	data "github.com/YoungBoyGod/OneGoServer/internal/data/repository"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/task"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"go.uber.org/zap"
)

// TaskService 任务服务接口
type TaskService interface {
	// 基础CRUD操作
	// 创建任务
	CreateTask(ctx context.Context, req *task.TaskCreateRequest) (*task.Task, error)
	// 获取任务
	GetTask(ctx context.Context, id int64) (*task.Task, error)
	// 更新任务
	UpdateTask(ctx context.Context, id int64, req *task.TaskUpdateRequest) (*task.Task, error)
	// 删除任务
	DeleteTask(ctx context.Context, id int64) error
	// 获取任务列表
	ListTasks(ctx context.Context, filter *task.TaskFilter, sort *task.TaskSortOption, pagination *task.PaginationOption) (*task.TaskListResponse, error)

	// 状态管理
	// 更新任务状态
	UpdateTaskStatus(ctx context.Context, id int64, status string) error
	// 获取任务状态
	GetTaskStatus(ctx context.Context, id int64) (string, error)

	// 执行控制
	// 执行任务
	ExecuteTask(ctx context.Context, id int64, req *task.TaskExecuteRequest) (*task.TaskExecution, error)
	// 取消任务
	CancelTask(ctx context.Context, id int64) error
	// 分发任务
	DispatchTask(ctx context.Context, id int64) error

	// 查询功能
	// 获取任务详情
	GetTaskDetail(ctx context.Context, id int64) (*task.Task, error)
	// 获取任务结果
	GetTaskResult(ctx context.Context, id int64) (*task.TaskExecution, error)
	// 获取任务执行记录
	GetTaskExecutions(ctx context.Context, id int64) ([]task.TaskExecution, error)
	// 获取任务统计
	GetTaskStats(ctx context.Context, filter *task.TaskFilter) (*task.TaskStatistics, error)
	// 查询任务
	QueryTasks(ctx context.Context, keyword string, filter *task.TaskFilter) ([]task.Task, error)
}

// taskServiceImpl 任务服务实现
type taskServiceImpl struct {
	taskRepo data.TaskRepository
	taskBiz  *task.TaskBusiness
	logger   *zap.Logger
}

// NewTaskService 创建任务服务实例
func NewTaskService(taskRepo data.TaskRepository) TaskService {
	return &taskServiceImpl{
		taskRepo: taskRepo,
		taskBiz:  task.NewTaskBusiness(),
		logger:   pkglog.GetAppLogger(nil), // 使用默认配置
	}
}

// CreateTask 创建任务
func (s *taskServiceImpl) CreateTask(ctx context.Context, req *task.TaskCreateRequest) (*task.Task, error) {
	// 参数验证
	if req.Name == "" {
		return nil, fmt.Errorf("任务名称不能为空")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("任务类型不能为空")
	}

	// 构建任务对象
	newTask := &task.Task{
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		Status:       task.TaskStatusPending,
		Priority:     5,   // 默认优先级
		Timeout:      300, // 默认5分钟超时
		MaxRetries:   3,   // 默认最大重试3次
		ExecuteTime:  req.ExecuteTime,
		ExecutorType: req.ExecutorType,
		ExecutorID:   req.ExecutorID,
		DeviceID:     req.DeviceID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 设置可选参数
	if req.Priority != nil {
		if *req.Priority < 1 || *req.Priority > 10 {
			return nil, fmt.Errorf("优先级必须在1-10之间")
		}
		newTask.Priority = *req.Priority
	}
	if req.Timeout != nil {
		newTask.Timeout = *req.Timeout
	}
	if req.MaxRetries != nil {
		newTask.MaxRetries = *req.MaxRetries
	}
	if req.Parameters != nil {
		newTask.Parameters = task.JSONB(req.Parameters)
	}

	// 业务逻辑：动态计算优先级
	newTask.Priority = s.taskBiz.CalculatePriority(newTask)

	// 创建任务
	if err := s.taskRepo.Create(ctx, newTask); err != nil {
		s.logger.Error("创建任务失败",
			zap.String("name", req.Name),
			zap.String("type", req.Type),
			zap.Error(err))
		return nil, fmt.Errorf("创建任务失败: %w", err)
	}

	s.logger.Info("任务创建成功",
		zap.Int64("task_id", newTask.ID),
		zap.String("name", newTask.Name),
		zap.String("type", newTask.Type),
		zap.Int("priority", newTask.Priority))

	return newTask, nil
}

// GetTask 获取任务
func (s *taskServiceImpl) GetTask(ctx context.Context, id int64) (*task.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("无效的任务ID: %d", id)
	}

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取任务失败",
			zap.Int64("task_id", id),
			zap.Error(err))
		return nil, err
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *taskServiceImpl) UpdateTask(ctx context.Context, id int64, req *task.TaskUpdateRequest) (*task.Task, error) {
	// 获取现有任务
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查任务状态是否允许更新
	if existingTask.Status == task.TaskStatusRunning {
		return nil, fmt.Errorf("正在执行的任务不能修改")
	}
	if existingTask.Status == task.TaskStatusCompleted {
		return nil, fmt.Errorf("已完成的任务不能修改")
	}

	// 更新字段
	hasChanges := false
	if req.Name != nil && *req.Name != existingTask.Name {
		existingTask.Name = *req.Name
		hasChanges = true
	}
	if req.Description != nil {
		existingTask.Description = req.Description
		hasChanges = true
	}
	if req.Priority != nil && *req.Priority != existingTask.Priority {
		if *req.Priority < 1 || *req.Priority > 10 {
			return nil, fmt.Errorf("优先级必须在1-10之间")
		}
		existingTask.Priority = *req.Priority
		hasChanges = true
	}
	if req.Timeout != nil {
		existingTask.Timeout = *req.Timeout
		hasChanges = true
	}
	if req.MaxRetries != nil {
		existingTask.MaxRetries = *req.MaxRetries
		hasChanges = true
	}
	if req.Parameters != nil {
		existingTask.Parameters = task.JSONB(req.Parameters)
		hasChanges = true
	}
	if req.ExecutorType != nil {
		existingTask.ExecutorType = req.ExecutorType
		hasChanges = true
	}
	if req.ExecutorID != nil {
		existingTask.ExecutorID = req.ExecutorID
		hasChanges = true
	}
	if req.DeviceID != nil {
		existingTask.DeviceID = req.DeviceID
		hasChanges = true
	}

	if !hasChanges {
		return existingTask, nil
	}

	// 保存更新
	existingTask.UpdatedAt = time.Now()
	if err := s.taskRepo.Update(ctx, existingTask); err != nil {
		s.logger.Error("更新任务失败",
			zap.Int64("task_id", id),
			zap.Error(err))
		return nil, fmt.Errorf("更新任务失败: %w", err)
	}

	s.logger.Info("任务更新成功",
		zap.Int64("task_id", id),
		zap.String("name", existingTask.Name))

	return existingTask, nil
}

// DeleteTask 删除任务
func (s *taskServiceImpl) DeleteTask(ctx context.Context, id int64) error {
	// 获取任务检查状态
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否可以删除
	if existingTask.Status == task.TaskStatusRunning {
		return fmt.Errorf("正在执行的任务不能删除")
	}

	// 执行删除
	if err := s.taskRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除任务失败",
			zap.Int64("task_id", id),
			zap.Error(err))
		return fmt.Errorf("删除任务失败: %w", err)
	}

	s.logger.Info("任务删除成功",
		zap.Int64("task_id", id),
		zap.String("name", existingTask.Name))

	return nil
}

// ListTasks 获取任务列表
func (s *taskServiceImpl) ListTasks(ctx context.Context, filter *task.TaskFilter, sort *task.TaskSortOption, pagination *task.PaginationOption) (*task.TaskListResponse, error) {
	// 设置默认分页
	if pagination == nil {
		pagination = &task.PaginationOption{Page: 1, Size: 20}
	}
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Size <= 0 {
		pagination.Size = 20
	}
	if pagination.Size > 100 {
		pagination.Size = 100 // 限制最大页面大小
	}

	// 设置默认排序
	if sort == nil {
		sort = &task.TaskSortOption{Field: "created_at", Order: "desc"}
	}

	// 查询任务列表
	tasks, total, err := s.taskRepo.List(ctx, filter, sort, pagination)
	if err != nil {
		s.logger.Error("获取任务列表失败", zap.Error(err))
		return nil, fmt.Errorf("获取任务列表失败: %w", err)
	}

	// 计算分页信息
	totalPages := int((total + int64(pagination.Size) - 1) / int64(pagination.Size))

	response := &task.TaskListResponse{
		Tasks: tasks,
		Pagination: task.PaginationInfo{
			Page:       pagination.Page,
			Size:       pagination.Size,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

// UpdateTaskStatus 更新任务状态
func (s *taskServiceImpl) UpdateTaskStatus(ctx context.Context, id int64, status string) error {
	// 验证状态
	validStatuses := []string{
		task.TaskStatusPending,
		task.TaskStatusQueued,
		task.TaskStatusAssigning,
		task.TaskStatusAssigned,
		task.TaskStatusDispatching,
		task.TaskStatusRunning,
		task.TaskStatusCompleted,
		task.TaskStatusFailed,
		task.TaskStatusCanceled,
	}

	validStatus := false
	for _, validS := range validStatuses {
		if status == validS {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return fmt.Errorf("无效的任务状态: %s", status)
	}

	// 业务逻辑校验状态转换
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !s.taskBiz.CanTransitionTo(existingTask, status) {
		return fmt.Errorf("不允许的状态转换: %s -> %s", existingTask.Status, status)
	}

	// 更新状态
	if err := s.taskRepo.UpdateStatus(ctx, id, status); err != nil {
		s.logger.Error("更新任务状态失败",
			zap.Int64("task_id", id),
			zap.String("status", status),
			zap.Error(err))
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	s.logger.Info("任务状态更新成功",
		zap.Int64("task_id", id),
		zap.String("status", status))

	return nil
}

// GetTaskStatus 获取任务状态
func (s *taskServiceImpl) GetTaskStatus(ctx context.Context, id int64) (string, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return task.Status, nil
}

// ExecuteTask 执行任务
func (s *taskServiceImpl) ExecuteTask(ctx context.Context, id int64, req *task.TaskExecuteRequest) (*task.TaskExecution, error) {
	// 获取任务
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查任务状态
	if existingTask.Status == task.TaskStatusRunning {
		return nil, fmt.Errorf("任务已在执行中")
	}
	if existingTask.Status == task.TaskStatusCompleted {
		return nil, fmt.Errorf("任务已完成")
	}

	// 创建执行记录
	execution := &task.TaskExecution{
		TaskID:    id,
		Status:    task.ExecutionStatusStarted,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
	}

	// 设置执行器信息
	if req != nil {
		if req.ExecutorType != nil {
			existingTask.ExecutorType = req.ExecutorType
		}
		if req.ExecutorID != nil {
			existingTask.ExecutorID = req.ExecutorID
		}
		if req.DeviceID != nil {
			existingTask.DeviceID = req.DeviceID
		}
	}

	// 更新任务状态为运行中
	existingTask.Status = task.TaskStatusRunning
	if err := s.taskRepo.Update(ctx, existingTask); err != nil {
		return nil, fmt.Errorf("更新任务状态失败: %w", err)
	}

	// 创建执行记录
	if err := s.taskRepo.CreateExecution(ctx, execution); err != nil {
		return nil, fmt.Errorf("创建执行记录失败: %w", err)
	}

	s.logger.Info("任务开始执行",
		zap.Int64("task_id", id),
		zap.String("execution_id", execution.ExecutionID))

	return execution, nil
}

// CancelTask 取消任务
func (s *taskServiceImpl) CancelTask(ctx context.Context, id int64) error {
	// 获取任务
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否可以取消
	if existingTask.Status == task.TaskStatusCompleted {
		return fmt.Errorf("已完成的任务不能取消")
	}
	if existingTask.Status == task.TaskStatusCanceled {
		return fmt.Errorf("任务已被取消")
	}

	// 更新状态
	if err := s.taskRepo.UpdateStatus(ctx, id, task.TaskStatusCanceled); err != nil {
		return fmt.Errorf("取消任务失败: %w", err)
	}

	s.logger.Info("任务已取消",
		zap.Int64("task_id", id),
		zap.String("name", existingTask.Name))

	return nil
}

// DispatchTask 分发任务
func (s *taskServiceImpl) DispatchTask(ctx context.Context, id int64) error {
	// 获取任务
	existingTask, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查任务状态
	if existingTask.Status != task.TaskStatusPending {
		return fmt.Errorf("只有待处理的任务才能分发")
	}

	// 更新状态为队列中
	if err := s.taskRepo.UpdateStatus(ctx, id, task.TaskStatusQueued); err != nil {
		return fmt.Errorf("分发任务失败: %w", err)
	}

	s.logger.Info("任务已分发",
		zap.Int64("task_id", id),
		zap.String("name", existingTask.Name))

	return nil
}

// GetTaskDetail 获取任务详情（包含执行记录）
func (s *taskServiceImpl) GetTaskDetail(ctx context.Context, id int64) (*task.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

// GetTaskResult 获取任务结果
func (s *taskServiceImpl) GetTaskResult(ctx context.Context, id int64) (*task.TaskExecution, error) {
	executions, err := s.taskRepo.GetExecutions(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(executions) == 0 {
		return nil, fmt.Errorf("任务还未执行")
	}

	// 返回最新的执行记录
	return &executions[len(executions)-1], nil
}

// GetTaskExecutions 获取任务执行记录
func (s *taskServiceImpl) GetTaskExecutions(ctx context.Context, id int64) ([]task.TaskExecution, error) {
	return s.taskRepo.GetExecutions(ctx, id)
}

// GetTaskStats 获取任务统计
func (s *taskServiceImpl) GetTaskStats(ctx context.Context, filter *task.TaskFilter) (*task.TaskStatistics, error) {
	return s.taskRepo.GetStatistics(ctx, filter)
}

// QueryTasks 查询任务
func (s *taskServiceImpl) QueryTasks(ctx context.Context, keyword string, filter *task.TaskFilter) ([]task.Task, error) {
	// 如果提供了关键字，添加到过滤条件中
	if keyword != "" {
		if filter == nil {
			filter = &task.TaskFilter{}
		}
		filter.Keyword = &keyword
	}

	// 使用List方法查询
	tasks, _, err := s.taskRepo.List(ctx, filter, nil, nil)
	return tasks, err
}
