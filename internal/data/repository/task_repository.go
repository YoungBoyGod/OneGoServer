package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/task"
	"gorm.io/gorm"
)

// TaskRepository 任务数据访问层接口
type TaskRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, task *task.Task) error
	GetByID(ctx context.Context, id int64) (*task.Task, error)
	Update(ctx context.Context, task *task.Task) error
	Delete(ctx context.Context, id int64) error

	// 查询操作
	List(ctx context.Context, filter *task.TaskFilter, sort *task.TaskSortOption, pagination *task.PaginationOption) ([]task.Task, int64, error)
	GetByDeviceID(ctx context.Context, deviceID int64) ([]task.Task, error)
	GetByStatus(ctx context.Context, status []string) ([]task.Task, error)

	// 状态管理
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateRetryCount(ctx context.Context, id int64, count int) error

	// 统计查询
	GetStatistics(ctx context.Context, filter *task.TaskFilter) (*task.TaskStatistics, error)
	CountByStatus(ctx context.Context) (map[string]int64, error)
	CountByType(ctx context.Context) (map[string]int64, error)

	// 执行记录相关
	CreateExecution(ctx context.Context, execution *task.TaskExecution) error
	GetExecutions(ctx context.Context, taskID int64) ([]task.TaskExecution, error)
	GetExecutionByID(ctx context.Context, executionID string) (*task.TaskExecution, error)
	UpdateExecution(ctx context.Context, execution *task.TaskExecution) error
}

// taskRepositoryImpl 任务Repository实现
type taskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务Repository实例
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepositoryImpl{
		db: db,
	}
}

// Create 创建任务
func (r *taskRepositoryImpl) Create(ctx context.Context, t *task.Task) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("创建任务失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取任务
func (r *taskRepositoryImpl) GetByID(ctx context.Context, id int64) (*task.Task, error) {
	var t task.Task
	if err := r.db.WithContext(ctx).
		Preload("Device").
		Preload("Executions").
		First(&t, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("任务不存在: %d", id)
		}
		return nil, fmt.Errorf("获取任务失败: %w", err)
	}
	return &t, nil
}

// Update 更新任务
func (r *taskRepositoryImpl) Update(ctx context.Context, t *task.Task) error {
	if err := r.db.WithContext(ctx).Save(t).Error; err != nil {
		return fmt.Errorf("更新任务失败: %w", err)
	}
	return nil
}

// Delete 删除任务
func (r *taskRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&task.Task{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除任务失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("任务不存在: %d", id)
	}
	return nil
}

// List 获取任务列表（支持过滤、排序、分页）
func (r *taskRepositoryImpl) List(ctx context.Context, filter *task.TaskFilter, sort *task.TaskSortOption, pagination *task.PaginationOption) ([]task.Task, int64, error) {
	var tasks []task.Task
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&task.Task{})

	// 应用过滤条件
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取任务总数失败: %w", err)
	}

	// 应用排序
	query = r.applySort(query, sort)

	// 应用分页
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Size
		query = query.Offset(offset).Limit(pagination.Size)
	}

	// 预加载关联数据
	query = query.Preload("Device").Preload("Executions")

	// 执行查询
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("获取任务列表失败: %w", err)
	}

	return tasks, total, nil
}

// GetByDeviceID 根据设备ID获取任务
func (r *taskRepositoryImpl) GetByDeviceID(ctx context.Context, deviceID int64) ([]task.Task, error) {
	var tasks []task.Task
	if err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Preload("Device").
		Preload("Executions").
		Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("根据设备ID获取任务失败: %w", err)
	}
	return tasks, nil
}

// GetByStatus 根据状态获取任务
func (r *taskRepositoryImpl) GetByStatus(ctx context.Context, status []string) ([]task.Task, error) {
	var tasks []task.Task
	if err := r.db.WithContext(ctx).
		Where("status IN ?", status).
		Preload("Device").
		Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("根据状态获取任务失败: %w", err)
	}
	return tasks, nil
}

// UpdateStatus 更新任务状态
func (r *taskRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status string) error {
	result := r.db.WithContext(ctx).
		Model(&task.Task{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新任务状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("任务不存在: %d", id)
	}
	return nil
}

// UpdateRetryCount 更新重试次数
func (r *taskRepositoryImpl) UpdateRetryCount(ctx context.Context, id int64, count int) error {
	result := r.db.WithContext(ctx).
		Model(&task.Task{}).
		Where("id = ?", id).
		Update("retry_count", count)

	if result.Error != nil {
		return fmt.Errorf("更新重试次数失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("任务不存在: %d", id)
	}
	return nil
}

// GetStatistics 获取任务统计信息
func (r *taskRepositoryImpl) GetStatistics(ctx context.Context, filter *task.TaskFilter) (*task.TaskStatistics, error) {
	stats := &task.TaskStatistics{
		ByStatus: make(map[string]int64),
		ByType:   make(map[string]int64),
	}

	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&task.Task{})
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("获取任务总数失败: %w", err)
	}

	// 按状态统计
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := query.Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("按状态统计失败: %w", err)
	}
	for _, sc := range statusCounts {
		stats.ByStatus[sc.Status] = sc.Count
		if sc.Status == task.TaskStatusCompleted {
			stats.Success = sc.Count
		} else if sc.Status == task.TaskStatusFailed {
			stats.Failed = sc.Count
		}
	}

	// 按类型统计
	var typeCounts []struct {
		Type  string
		Count int64
	}
	if err := query.Select("type, COUNT(*) as count").
		Group("type").
		Find(&typeCounts).Error; err != nil {
		return nil, fmt.Errorf("按类型统计失败: %w", err)
	}
	for _, tc := range typeCounts {
		stats.ByType[tc.Type] = tc.Count
	}

	// 最近24小时创建的任务数
	if err := query.Where("created_at >= ?", time.Now().Add(-24*time.Hour)).
		Count(&stats.Recent24h).Error; err != nil {
		return nil, fmt.Errorf("获取最近24小时任务数失败: %w", err)
	}

	return stats, nil
}

// CountByStatus 按状态统计任务数量
func (r *taskRepositoryImpl) CountByStatus(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&task.Task{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("按状态统计任务失败: %w", err)
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// CountByType 按类型统计任务数量
func (r *taskRepositoryImpl) CountByType(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&task.Task{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("按类型统计任务失败: %w", err)
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Type] = r.Count
	}
	return counts, nil
}

// CreateExecution 创建任务执行记录
func (r *taskRepositoryImpl) CreateExecution(ctx context.Context, execution *task.TaskExecution) error {
	if err := r.db.WithContext(ctx).Create(execution).Error; err != nil {
		return fmt.Errorf("创建执行记录失败: %w", err)
	}
	return nil
}

// GetExecutions 获取任务的执行记录
func (r *taskRepositoryImpl) GetExecutions(ctx context.Context, taskID int64) ([]task.TaskExecution, error) {
	var executions []task.TaskExecution
	if err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Order("start_time DESC").
		Find(&executions).Error; err != nil {
		return nil, fmt.Errorf("获取执行记录失败: %w", err)
	}
	return executions, nil
}

// GetExecutionByID 根据执行ID获取执行记录
func (r *taskRepositoryImpl) GetExecutionByID(ctx context.Context, executionID string) (*task.TaskExecution, error) {
	var execution task.TaskExecution
	if err := r.db.WithContext(ctx).
		Where("execution_id = ?", executionID).
		Preload("Task").
		First(&execution).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("执行记录不存在: %s", executionID)
		}
		return nil, fmt.Errorf("获取执行记录失败: %w", err)
	}
	return &execution, nil
}

// UpdateExecution 更新任务执行记录
func (r *taskRepositoryImpl) UpdateExecution(ctx context.Context, execution *task.TaskExecution) error {
	if err := r.db.WithContext(ctx).Save(execution).Error; err != nil {
		return fmt.Errorf("更新执行记录失败: %w", err)
	}
	return nil
}

// applyFilter 应用过滤条件
func (r *taskRepositoryImpl) applyFilter(query *gorm.DB, filter *task.TaskFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	// 状态过滤
	if len(filter.Status) > 0 {
		query = query.Where("status IN ?", filter.Status)
	}

	// 类型过滤
	if len(filter.Type) > 0 {
		query = query.Where("type IN ?", filter.Type)
	}

	// 优先级过滤
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	// 设备ID过滤
	if filter.DeviceID != nil {
		query = query.Where("device_id = ?", *filter.DeviceID)
	}

	// 创建者过滤
	if filter.CreatedBy != nil {
		query = query.Where("created_by = ?", *filter.CreatedBy)
	}

	// 时间范围过滤
	if filter.StartTime != nil {
		query = query.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("created_at <= ?", *filter.EndTime)
	}

	// 关键词搜索
	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := "%" + strings.ToLower(*filter.Keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	return query
}

// applySort 应用排序
func (r *taskRepositoryImpl) applySort(query *gorm.DB, sort *task.TaskSortOption) *gorm.DB {
	if sort == nil {
		// 默认按创建时间倒序
		return query.Order("created_at DESC")
	}

	// 验证排序字段
	validFields := map[string]bool{
		"id":         true,
		"name":       true,
		"status":     true,
		"priority":   true,
		"created_at": true,
		"updated_at": true,
	}

	if !validFields[sort.Field] {
		sort.Field = "created_at"
	}

	// 验证排序方向
	if sort.Order != "asc" && sort.Order != "desc" {
		sort.Order = "desc"
	}

	return query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Order))
}
