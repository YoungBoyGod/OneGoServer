package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/device"

	"gorm.io/gorm"
)

// DeviceRepository 设备数据访问层接口
type DeviceRepository interface {
	// 基础CRUD操作
	Create(ctx context.Context, device *device.Device) error
	GetByID(ctx context.Context, id int64) (*device.Device, error)
	GetByDeviceID(ctx context.Context, deviceID string) (*device.Device, error)
	Update(ctx context.Context, device *device.Device) error
	Delete(ctx context.Context, id int64) error

	// 查询操作
	List(ctx context.Context, filter *device.DeviceFilter, sort *device.DeviceSortOption, pagination *device.PaginationOption) ([]device.Device, int64, error)
	GetByType(ctx context.Context, deviceType []string) ([]device.Device, error)
	GetByStatus(ctx context.Context, status []string) ([]device.Device, error)
	GetOnlineDevices(ctx context.Context) ([]device.Device, error)
	GetOfflineDevices(ctx context.Context, duration time.Duration) ([]device.Device, error)

	// 状态管理
	UpdateStatus(ctx context.Context, deviceID string, status string) error
	UpdateHealthScore(ctx context.Context, deviceID string, score int) error
	UpdateLastSeen(ctx context.Context, deviceID string, lastSeen time.Time) error
	BatchUpdateStatus(ctx context.Context, deviceIDs []string, status string) error

	// 统计查询
	GetStatistics(ctx context.Context, filter *device.DeviceFilter) (*device.DeviceStatistics, error)
	CountByStatus(ctx context.Context) (map[string]int64, error)
	CountByType(ctx context.Context) (map[string]int64, error)
	GetHealthReport(ctx context.Context) (map[string]interface{}, error)

	// 心跳管理
	CreateHeartbeat(ctx context.Context, heartbeat *device.DeviceHeartbeat) error
	GetLatestHeartbeat(ctx context.Context, deviceID int64) (*device.DeviceHeartbeat, error)
	GetHeartbeatHistory(ctx context.Context, deviceID int64, hours int) ([]device.DeviceHeartbeat, error)
	CleanupOldHeartbeats(ctx context.Context, days int) error

	// 日志管理
	CreateLog(ctx context.Context, log *device.DeviceLog) error
	GetLogs(ctx context.Context, deviceID int64, filter *device.LogFilter) ([]device.DeviceLog, error)
	GetLogsByLevel(ctx context.Context, level string, limit int) ([]device.DeviceLog, error)
	CleanupOldLogs(ctx context.Context, days int) error

	// 命令管理
	CreateCommand(ctx context.Context, command *device.DeviceCommand) error
	GetCommand(ctx context.Context, commandID string) (*device.DeviceCommand, error)
	GetDeviceCommands(ctx context.Context, deviceID int64, status []string) ([]device.DeviceCommand, error)
	UpdateCommandStatus(ctx context.Context, commandID string, status string) error
	UpdateCommandResponse(ctx context.Context, commandID string, response device.JSONB, errorMsg *string) error
}

// deviceRepositoryImpl 设备Repository实现
type deviceRepositoryImpl struct {
	db *gorm.DB
}

// NewDeviceRepository 创建设备Repository实例
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepositoryImpl{
		db: db,
	}
}

// Create 创建设备
func (r *deviceRepositoryImpl) Create(ctx context.Context, d *device.Device) error {
	if err := r.db.WithContext(ctx).Create(d).Error; err != nil {
		return fmt.Errorf("创建设备失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取设备
func (r *deviceRepositoryImpl) GetByID(ctx context.Context, id int64) (*device.Device, error) {
	var d device.Device
	if err := r.db.WithContext(ctx).
		Preload("Heartbeats", func(db *gorm.DB) *gorm.DB {
			return db.Order("heartbeat_time DESC").Limit(10)
		}).
		Preload("Logs", func(db *gorm.DB) *gorm.DB {
			return db.Order("log_time DESC").Limit(50)
		}).
		Preload("Commands", func(db *gorm.DB) *gorm.DB {
			return db.Order("sent_time DESC").Limit(20)
		}).
		First(&d, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("设备不存在: %d", id)
		}
		return nil, fmt.Errorf("获取设备失败: %w", err)
	}
	return &d, nil
}

// GetByDeviceID 根据设备ID获取设备
func (r *deviceRepositoryImpl) GetByDeviceID(ctx context.Context, deviceID string) (*device.Device, error) {
	var d device.Device
	if err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Preload("Heartbeats", func(db *gorm.DB) *gorm.DB {
			return db.Order("heartbeat_time DESC").Limit(5)
		}).
		First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("设备不存在: %s", deviceID)
		}
		return nil, fmt.Errorf("获取设备失败: %w", err)
	}
	return &d, nil
}

// Update 更新设备
func (r *deviceRepositoryImpl) Update(ctx context.Context, d *device.Device) error {
	if err := r.db.WithContext(ctx).Save(d).Error; err != nil {
		return fmt.Errorf("更新设备失败: %w", err)
	}
	return nil
}

// Delete 删除设备
func (r *deviceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&device.Device{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除设备失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("设备不存在: %d", id)
	}
	return nil
}

// List 获取设备列表（支持过滤、排序、分页）
func (r *deviceRepositoryImpl) List(ctx context.Context, filter *device.DeviceFilter, sort *device.DeviceSortOption, pagination *device.PaginationOption) ([]device.Device, int64, error) {
	var devices []device.Device
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(&device.Device{})

	// 应用过滤条件
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取设备总数失败: %w", err)
	}

	// 应用排序
	query = r.applySort(query, sort)

	// 应用分页
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Size
		query = query.Offset(offset).Limit(pagination.Size)
	}

	// 预加载最新心跳
	query = query.Preload("Heartbeats", func(db *gorm.DB) *gorm.DB {
		return db.Order("heartbeat_time DESC").Limit(1)
	})

	// 执行查询
	if err := query.Find(&devices).Error; err != nil {
		return nil, 0, fmt.Errorf("获取设备列表失败: %w", err)
	}

	return devices, total, nil
}

// GetByType 根据设备类型获取设备
func (r *deviceRepositoryImpl) GetByType(ctx context.Context, deviceType []string) ([]device.Device, error) {
	var devices []device.Device
	if err := r.db.WithContext(ctx).
		Where("type IN ?", deviceType).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("根据类型获取设备失败: %w", err)
	}
	return devices, nil
}

// GetByStatus 根据状态获取设备
func (r *deviceRepositoryImpl) GetByStatus(ctx context.Context, status []string) ([]device.Device, error) {
	var devices []device.Device
	if err := r.db.WithContext(ctx).
		Where("status IN ?", status).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("根据状态获取设备失败: %w", err)
	}
	return devices, nil
}

// GetOnlineDevices 获取在线设备
func (r *deviceRepositoryImpl) GetOnlineDevices(ctx context.Context) ([]device.Device, error) {
	var devices []device.Device
	if err := r.db.WithContext(ctx).
		Where("status = ?", device.DeviceStatusOnline).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("获取在线设备失败: %w", err)
	}
	return devices, nil
}

// GetOfflineDevices 获取长时间离线的设备
func (r *deviceRepositoryImpl) GetOfflineDevices(ctx context.Context, duration time.Duration) ([]device.Device, error) {
	var devices []device.Device
	threshold := time.Now().Add(-duration)

	if err := r.db.WithContext(ctx).
		Where("status = ? AND (last_seen IS NULL OR last_seen < ?)",
			device.DeviceStatusOffline, threshold).
		Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("获取离线设备失败: %w", err)
	}
	return devices, nil
}

// UpdateStatus 更新设备状态
func (r *deviceRepositoryImpl) UpdateStatus(ctx context.Context, deviceID string, status string) error {
	result := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Where("device_id = ?", deviceID).
		Updates(map[string]interface{}{
			"status":    status,
			"last_seen": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("更新设备状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}
	return nil
}

// UpdateHealthScore 更新设备健康度
func (r *deviceRepositoryImpl) UpdateHealthScore(ctx context.Context, deviceID string, score int) error {
	result := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Where("device_id = ?", deviceID).
		Update("health_score", score)

	if result.Error != nil {
		return fmt.Errorf("更新健康度失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}
	return nil
}

// UpdateLastSeen 更新设备最后活跃时间
func (r *deviceRepositoryImpl) UpdateLastSeen(ctx context.Context, deviceID string, lastSeen time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Where("device_id = ?", deviceID).
		Update("last_seen", lastSeen)

	if result.Error != nil {
		return fmt.Errorf("更新最后活跃时间失败: %w", result.Error)
	}
	return nil
}

// BatchUpdateStatus 批量更新设备状态
func (r *deviceRepositoryImpl) BatchUpdateStatus(ctx context.Context, deviceIDs []string, status string) error {
	result := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Where("device_id IN ?", deviceIDs).
		Updates(map[string]interface{}{
			"status":    status,
			"last_seen": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("批量更新设备状态失败: %w", result.Error)
	}
	return nil
}

// GetStatistics 获取设备统计信息
func (r *deviceRepositoryImpl) GetStatistics(ctx context.Context, filter *device.DeviceFilter) (*device.DeviceStatistics, error) {
	stats := &device.DeviceStatistics{
		ByStatus: make(map[string]int64),
		ByType:   make(map[string]int64),
	}

	// 构建基础查询
	query := r.db.WithContext(ctx).Model(&device.Device{})
	query = r.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("获取设备总数失败: %w", err)
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
		if sc.Status == device.DeviceStatusOnline {
			stats.Online = sc.Count
		} else if sc.Status == device.DeviceStatusOffline {
			stats.Offline = sc.Count
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

	// 平均健康度
	var avgHealth struct {
		Avg float64
	}
	if err := query.Select("AVG(health_score) as avg").
		Find(&avgHealth).Error; err != nil {
		return nil, fmt.Errorf("计算平均健康度失败: %w", err)
	}
	stats.AvgHealth = avgHealth.Avg

	// 最近24小时活跃设备
	if err := query.Where("last_seen >= ?", time.Now().Add(-24*time.Hour)).
		Count(&stats.RecentActive).Error; err != nil {
		return nil, fmt.Errorf("获取最近活跃设备数失败: %w", err)
	}

	return stats, nil
}

// CountByStatus 按状态统计设备数量
func (r *deviceRepositoryImpl) CountByStatus(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("按状态统计设备失败: %w", err)
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// CountByType 按类型统计设备数量
func (r *deviceRepositoryImpl) CountByType(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Type  string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("按类型统计设备失败: %w", err)
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Type] = r.Count
	}
	return counts, nil
}

// GetHealthReport 获取设备健康报告
func (r *deviceRepositoryImpl) GetHealthReport(ctx context.Context) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// 健康度分布
	var healthDistribution []struct {
		Range string
		Count int64
	}

	if err := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Select(`
			CASE 
				WHEN health_score >= 90 THEN 'excellent'
				WHEN health_score >= 70 THEN 'good'
				WHEN health_score >= 50 THEN 'fair'
				ELSE 'poor'
			END as range,
			COUNT(*) as count
		`).
		Group("range").
		Find(&healthDistribution).Error; err != nil {
		return nil, fmt.Errorf("获取健康度分布失败: %w", err)
	}

	report["health_distribution"] = healthDistribution

	// 平均运行时间
	var avgUptime struct {
		Avg float64
	}
	if err := r.db.WithContext(ctx).
		Model(&device.Device{}).
		Select("AVG(uptime_hours) as avg").
		Find(&avgUptime).Error; err != nil {
		return nil, fmt.Errorf("计算平均运行时间失败: %w", err)
	}
	report["avg_uptime_hours"] = avgUptime.Avg

	return report, nil
}

// CreateHeartbeat 创建设备心跳记录
func (r *deviceRepositoryImpl) CreateHeartbeat(ctx context.Context, heartbeat *device.DeviceHeartbeat) error {
	if err := r.db.WithContext(ctx).Create(heartbeat).Error; err != nil {
		return fmt.Errorf("创建心跳记录失败: %w", err)
	}
	return nil
}

// GetLatestHeartbeat 获取设备最新心跳
func (r *deviceRepositoryImpl) GetLatestHeartbeat(ctx context.Context, deviceID int64) (*device.DeviceHeartbeat, error) {
	var heartbeat device.DeviceHeartbeat
	if err := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("heartbeat_time DESC").
		First(&heartbeat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("心跳记录不存在")
		}
		return nil, fmt.Errorf("获取最新心跳失败: %w", err)
	}
	return &heartbeat, nil
}

// GetHeartbeatHistory 获取设备心跳历史
func (r *deviceRepositoryImpl) GetHeartbeatHistory(ctx context.Context, deviceID int64, hours int) ([]device.DeviceHeartbeat, error) {
	var heartbeats []device.DeviceHeartbeat
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	if err := r.db.WithContext(ctx).
		Where("device_id = ? AND heartbeat_time >= ?", deviceID, startTime).
		Order("heartbeat_time DESC").
		Find(&heartbeats).Error; err != nil {
		return nil, fmt.Errorf("获取心跳历史失败: %w", err)
	}
	return heartbeats, nil
}

// CleanupOldHeartbeats 清理旧的心跳记录
func (r *deviceRepositoryImpl) CleanupOldHeartbeats(ctx context.Context, days int) error {
	cutoffTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour)

	result := r.db.WithContext(ctx).
		Where("heartbeat_time < ?", cutoffTime).
		Delete(&device.DeviceHeartbeat{})

	if result.Error != nil {
		return fmt.Errorf("清理旧心跳记录失败: %w", result.Error)
	}
	return nil
}

// CreateLog 创建设备日志
func (r *deviceRepositoryImpl) CreateLog(ctx context.Context, log *device.DeviceLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("创建设备日志失败: %w", err)
	}
	return nil
}

// LogFilter 日志过滤器
type LogFilter struct {
	Level     *string
	Category  *string
	StartTime *time.Time
	EndTime   *time.Time
	Keyword   *string
}

// GetLogs 获取设备日志
func (r *deviceRepositoryImpl) GetLogs(ctx context.Context, deviceID int64, filter *LogFilter) ([]device.DeviceLog, error) {
	var logs []device.DeviceLog

	query := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("log_time DESC")

	// 应用过滤条件
	if filter != nil {
		if filter.Level != nil {
			query = query.Where("level = ?", *filter.Level)
		}
		if filter.Category != nil {
			query = query.Where("category = ?", *filter.Category)
		}
		if filter.StartTime != nil {
			query = query.Where("log_time >= ?", *filter.StartTime)
		}
		if filter.EndTime != nil {
			query = query.Where("log_time <= ?", *filter.EndTime)
		}
		if filter.Keyword != nil && *filter.Keyword != "" {
			keyword := "%" + strings.ToLower(*filter.Keyword) + "%"
			query = query.Where("LOWER(message) LIKE ?", keyword)
		}
	}

	if err := query.Limit(1000).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("获取设备日志失败: %w", err)
	}
	return logs, nil
}

// GetLogsByLevel 根据日志级别获取日志
func (r *deviceRepositoryImpl) GetLogsByLevel(ctx context.Context, level string, limit int) ([]device.DeviceLog, error) {
	var logs []device.DeviceLog
	if err := r.db.WithContext(ctx).
		Where("level = ?", level).
		Order("log_time DESC").
		Limit(limit).
		Preload("Device").
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("根据级别获取日志失败: %w", err)
	}
	return logs, nil
}

// CleanupOldLogs 清理旧的日志记录
func (r *deviceRepositoryImpl) CleanupOldLogs(ctx context.Context, days int) error {
	cutoffTime := time.Now().Add(-time.Duration(days) * 24 * time.Hour)

	result := r.db.WithContext(ctx).
		Where("log_time < ?", cutoffTime).
		Delete(&device.DeviceLog{})

	if result.Error != nil {
		return fmt.Errorf("清理旧日志记录失败: %w", result.Error)
	}
	return nil
}

// CreateCommand 创建设备命令
func (r *deviceRepositoryImpl) CreateCommand(ctx context.Context, command *device.DeviceCommand) error {
	if err := r.db.WithContext(ctx).Create(command).Error; err != nil {
		return fmt.Errorf("创建设备命令失败: %w", err)
	}
	return nil
}

// GetCommand 获取设备命令
func (r *deviceRepositoryImpl) GetCommand(ctx context.Context, commandID string) (*device.DeviceCommand, error) {
	var command device.DeviceCommand
	if err := r.db.WithContext(ctx).
		Where("command_id = ?", commandID).
		Preload("Device").
		First(&command).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("命令不存在: %s", commandID)
		}
		return nil, fmt.Errorf("获取设备命令失败: %w", err)
	}
	return &command, nil
}

// GetDeviceCommands 获取设备的命令列表
func (r *deviceRepositoryImpl) GetDeviceCommands(ctx context.Context, deviceID int64, status []string) ([]device.DeviceCommand, error) {
	var commands []device.DeviceCommand

	query := r.db.WithContext(ctx).
		Where("device_id = ?", deviceID).
		Order("sent_time DESC")

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	if err := query.Find(&commands).Error; err != nil {
		return nil, fmt.Errorf("获取设备命令失败: %w", err)
	}
	return commands, nil
}

// UpdateCommandStatus 更新命令状态
func (r *deviceRepositoryImpl) UpdateCommandStatus(ctx context.Context, commandID string, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 根据状态设置时间字段
	now := time.Now()
	switch status {
	case device.CommandStatusExecuted:
		updates["executed_time"] = now
	case device.CommandStatusCompleted, device.CommandStatusFailed:
		updates["completed_time"] = now
	}

	result := r.db.WithContext(ctx).
		Model(&device.DeviceCommand{}).
		Where("command_id = ?", commandID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("更新命令状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("命令不存在: %s", commandID)
	}
	return nil
}

// UpdateCommandResponse 更新命令响应
func (r *deviceRepositoryImpl) UpdateCommandResponse(ctx context.Context, commandID string, response device.JSONB, errorMsg *string) error {
	updates := map[string]interface{}{
		"response_data":  response,
		"completed_time": time.Now(),
	}

	if errorMsg != nil {
		updates["error_message"] = *errorMsg
		updates["status"] = device.CommandStatusFailed
	} else {
		updates["status"] = device.CommandStatusCompleted
	}

	result := r.db.WithContext(ctx).
		Model(&device.DeviceCommand{}).
		Where("command_id = ?", commandID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("更新命令响应失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("命令不存在: %s", commandID)
	}
	return nil
}

// applyFilter 应用过滤条件
func (r *deviceRepositoryImpl) applyFilter(query *gorm.DB, filter *device.DeviceFilter) *gorm.DB {
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

	// 制造商过滤
	if filter.Manufacturer != nil {
		query = query.Where("manufacturer = ?", *filter.Manufacturer)
	}

	// 型号过滤
	if filter.Model != nil {
		query = query.Where("model = ?", *filter.Model)
	}

	// 协议过滤
	if filter.Protocol != nil {
		query = query.Where("protocol = ?", *filter.Protocol)
	}

	// 健康度范围过滤
	if filter.MinHealth != nil {
		query = query.Where("health_score >= ?", *filter.MinHealth)
	}
	if filter.MaxHealth != nil {
		query = query.Where("health_score <= ?", *filter.MaxHealth)
	}

	// 最后活跃时间范围过滤
	if filter.LastSeenFrom != nil {
		query = query.Where("last_seen >= ?", *filter.LastSeenFrom)
	}
	if filter.LastSeenTo != nil {
		query = query.Where("last_seen <= ?", *filter.LastSeenTo)
	}

	// 关键词搜索
	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := "%" + strings.ToLower(*filter.Keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(device_id) LIKE ?", keyword, keyword)
	}

	return query
}

// applySort 应用排序
func (r *deviceRepositoryImpl) applySort(query *gorm.DB, sort *device.DeviceSortOption) *gorm.DB {
	if sort == nil {
		// 默认按最后活跃时间倒序
		return query.Order("last_seen DESC NULLS LAST")
	}

	// 验证排序字段
	validFields := map[string]bool{
		"id":           true,
		"name":         true,
		"status":       true,
		"health_score": true,
		"last_seen":    true,
		"created_at":   true,
	}

	if !validFields[sort.Field] {
		sort.Field = "last_seen"
	}

	// 验证排序方向
	if sort.Order != "asc" && sort.Order != "desc" {
		sort.Order = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sort.Field, sort.Order)
	if sort.Field == "last_seen" {
		orderClause += " NULLS LAST"
	}

	return query.Order(orderClause)
}
