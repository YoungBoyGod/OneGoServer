package service

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/device"
	"github.com/YoungBoyGod/OneGoServer/internal/data/repository"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"go.uber.org/zap"
)

// DeviceService 设备服务接口
type DeviceService interface {
	// 基础CRUD操作
	// 注册设备
	RegisterDevice(ctx context.Context, req *device.DeviceCreateRequest) (*device.Device, error)
	// 获取设备
	GetDevice(ctx context.Context, id int64) (*device.Device, error)
	// 根据设备ID获取设备
	GetDeviceByDeviceID(ctx context.Context, deviceID string) (*device.Device, error)
	// 更新设备
	UpdateDevice(ctx context.Context, id int64, req *device.DeviceUpdateRequest) (*device.Device, error)
	// 删除设备
	DeleteDevice(ctx context.Context, id int64) error
	// 获取设备列表
	ListDevices(ctx context.Context, filter *device.DeviceFilter, sort *device.DeviceSortOption, pagination *device.PaginationOption) (*device.DeviceListResponse, error)

	// 设备状态管理
	// 设备上线
	DeviceOnline(ctx context.Context, deviceID string) error
	// 设备下线
	DeviceOffline(ctx context.Context, deviceID string) error
	// 更新设备状态
	UpdateDeviceStatus(ctx context.Context, deviceID string, status string) error
	// 获取设备状态
	GetDeviceStatus(ctx context.Context, deviceID string) (*device.DeviceStatusResponse, error)
	// 更新设备健康度
	UpdateHealthScore(ctx context.Context, deviceID string, score int) error
	// 批量更新设备状态
	BatchUpdateStatus(ctx context.Context, deviceIDs []string, status string) error

	// 心跳管理
	// 处理设备心跳
	ProcessHeartbeat(ctx context.Context, deviceID string, req *device.DeviceHeartbeatRequest) error
	// 获取设备心跳
	GetDeviceHeartbeat(ctx context.Context, deviceID string) (*device.DeviceHeartbeat, error)
	// 获取心跳历史
	GetHeartbeatHistory(ctx context.Context, deviceID string, hours int) ([]device.DeviceHeartbeat, error)
	// 检查离线设备
	CheckOfflineDevices(ctx context.Context, duration time.Duration) ([]device.Device, error)

	// 设备命令管理
	// 发送设备命令
	SendCommand(ctx context.Context, deviceID string, req *device.DeviceCommandRequest) (*device.DeviceCommand, error)
	// 获取设备命令
	GetCommand(ctx context.Context, commandID string) (*device.DeviceCommand, error)
	// 获取设备命令列表
	GetDeviceCommands(ctx context.Context, deviceID string, status []string) ([]device.DeviceCommand, error)
	// 更新命令状态
	UpdateCommandStatus(ctx context.Context, commandID string, status string) error
	// 更新命令响应
	UpdateCommandResponse(ctx context.Context, commandID string, response device.JSONB, errorMsg *string) error

	// 设备日志管理
	// 创建设备日志
	CreateDeviceLog(ctx context.Context, deviceID string, req *device.DeviceLogRequest) error
	// 获取设备日志
	GetDeviceLogs(ctx context.Context, deviceID string, filter *device.LogFilter) ([]device.DeviceLog, error)
	// 根据级别获取日志
	GetLogsByLevel(ctx context.Context, level string, limit int) ([]device.DeviceLog, error)

	// 查询和统计
	// 查询设备
	QueryDevices(ctx context.Context, keyword string, filter *device.DeviceFilter) ([]device.Device, error)
	// 获取设备统计
	GetDeviceStats(ctx context.Context, filter *device.DeviceFilter) (*device.DeviceStatistics, error)
	// 获取在线设备
	GetOnlineDevices(ctx context.Context) ([]device.Device, error)
	// 获取设备健康报告
	GetHealthReport(ctx context.Context) (map[string]interface{}, error)

	// 系统维护
	// 清理过期心跳
	CleanupOldHeartbeats(ctx context.Context, days int) error
	// 清理过期日志
	CleanupOldLogs(ctx context.Context, days int) error
}

// deviceServiceImpl 设备服务实现
type deviceServiceImpl struct {
	deviceRepo repository.DeviceRepository
	deviceBiz  *device.DeviceBusiness
	logger     *zap.Logger
}

// NewDeviceService 创建设备服务实例
func NewDeviceService(deviceRepo repository.DeviceRepository) DeviceService {
	return &deviceServiceImpl{
		deviceRepo: deviceRepo,
		deviceBiz:  device.NewDeviceBusiness(),
		logger:     pkglog.GetAppLogger(nil),
	}
}

// RegisterDevice 注册设备
func (s *deviceServiceImpl) RegisterDevice(ctx context.Context, req *device.DeviceCreateRequest) (*device.Device, error) {
	// 参数验证
	if req.DeviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("设备名称不能为空")
	}
	if req.Type == "" {
		return nil, fmt.Errorf("设备类型不能为空")
	}

	// 检查设备ID是否已存在
	existingDevice, err := s.deviceRepo.GetByDeviceID(ctx, req.DeviceID)
	if err == nil && existingDevice != nil {
		return nil, fmt.Errorf("设备ID已存在: %s", req.DeviceID)
	}

	// 构建设备对象
	newDevice := &device.Device{
		DeviceID:     req.DeviceID,
		Name:         req.Name,
		Type:         req.Type,
		Model:        req.Model,
		Manufacturer: req.Manufacturer,
		Status:       device.DeviceStatusOffline,
		HealthScore:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 设置可选字段
	if req.IPAddress != nil {
		ip := net.ParseIP(*req.IPAddress)
		if ip == nil {
			return nil, fmt.Errorf("无效的IP地址: %s", *req.IPAddress)
		}
		newDevice.IPAddress = &ip
	}
	if req.Port != nil {
		newDevice.Port = req.Port
	}
	if req.Protocol != nil {
		newDevice.Protocol = req.Protocol
	}
	if req.Endpoint != nil {
		newDevice.Endpoint = req.Endpoint
	}
	if req.AuthType != nil {
		newDevice.AuthType = req.AuthType
	}
	if req.Credentials != nil {
		credentials := device.JSONB(req.Credentials)
		newDevice.Credentials = &credentials
	}
	if req.Config != nil {
		config := device.JSONB(req.Config)
		newDevice.Config = &config
	}
	if req.Capabilities != nil {
		capabilities := device.JSONB(req.Capabilities)
		newDevice.Capabilities = &capabilities
	}
	if req.Metadata != nil {
		metadata := device.JSONB(req.Metadata)
		newDevice.Metadata = &metadata
	}

	// 业务逻辑：计算健康度与状态
	newDevice.HealthScore = s.deviceBiz.CalculateHealthScore(newDevice)
	newDevice.Status = s.deviceBiz.DetermineDeviceStatus(newDevice, nil)

	// 创建设备
	if err := s.deviceRepo.Create(ctx, newDevice); err != nil {
		s.logger.Error("注册设备失败",
			zap.String("device_id", req.DeviceID),
			zap.String("name", req.Name),
			zap.Error(err))
		return nil, fmt.Errorf("注册设备失败: %w", err)
	}

	s.logger.Info("设备注册成功",
		zap.Int64("id", newDevice.ID),
		zap.String("device_id", newDevice.DeviceID),
		zap.String("name", newDevice.Name),
		zap.String("type", newDevice.Type))

	return newDevice, nil
}

// GetDevice 获取设备
func (s *deviceServiceImpl) GetDevice(ctx context.Context, id int64) (*device.Device, error) {
	if id <= 0 {
		return nil, fmt.Errorf("无效的设备ID: %d", id)
	}

	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取设备失败",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, err
	}

	return device, nil
}

// GetDeviceByDeviceID 根据设备ID获取设备
func (s *deviceServiceImpl) GetDeviceByDeviceID(ctx context.Context, deviceID string) (*device.Device, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}

	device, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		s.logger.Error("根据设备ID获取设备失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return nil, err
	}

	return device, nil
}

// UpdateDevice 更新设备
func (s *deviceServiceImpl) UpdateDevice(ctx context.Context, id int64, req *device.DeviceUpdateRequest) (*device.Device, error) {
	// 获取现有设备
	existingDevice, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	hasChanges := false
	if req.Name != nil && *req.Name != existingDevice.Name {
		existingDevice.Name = *req.Name
		hasChanges = true
	}
	if req.Model != nil {
		existingDevice.Model = req.Model
		hasChanges = true
	}
	if req.Manufacturer != nil {
		existingDevice.Manufacturer = req.Manufacturer
		hasChanges = true
	}
	if req.IPAddress != nil {
		ip := net.ParseIP(*req.IPAddress)
		if ip == nil {
			return nil, fmt.Errorf("无效的IP地址: %s", *req.IPAddress)
		}
		existingDevice.IPAddress = &ip
		hasChanges = true
	}
	if req.Port != nil {
		existingDevice.Port = req.Port
		hasChanges = true
	}
	if req.Protocol != nil {
		existingDevice.Protocol = req.Protocol
		hasChanges = true
	}
	if req.Endpoint != nil {
		existingDevice.Endpoint = req.Endpoint
		hasChanges = true
	}
	if req.AuthType != nil {
		existingDevice.AuthType = req.AuthType
		hasChanges = true
	}
	if req.Credentials != nil {
		credentials := device.JSONB(req.Credentials)
		existingDevice.Credentials = &credentials
		hasChanges = true
	}
	if req.Config != nil {
		config := device.JSONB(req.Config)
		existingDevice.Config = &config
		hasChanges = true
	}
	if req.Capabilities != nil {
		capabilities := device.JSONB(req.Capabilities)
		existingDevice.Capabilities = &capabilities
		hasChanges = true
	}
	if req.Metadata != nil {
		metadata := device.JSONB(req.Metadata)
		existingDevice.Metadata = &metadata
		hasChanges = true
	}

	if !hasChanges {
		s.logger.Info("设备信息无变化，跳过更新",
			zap.Int64("id", id))
		return existingDevice, nil
	}

	// 更新时间戳
	existingDevice.UpdatedAt = time.Now()

	// 保存更新
	if err := s.deviceRepo.Update(ctx, existingDevice); err != nil {
		s.logger.Error("更新设备失败",
			zap.Int64("id", id),
			zap.Error(err))
		return nil, fmt.Errorf("更新设备失败: %w", err)
	}

	s.logger.Info("设备更新成功",
		zap.Int64("id", existingDevice.ID),
		zap.String("device_id", existingDevice.DeviceID),
		zap.String("name", existingDevice.Name))

	return existingDevice, nil
}

// DeleteDevice 删除设备
func (s *deviceServiceImpl) DeleteDevice(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("无效的设备ID: %d", id)
	}

	// 检查设备是否存在
	existingDevice, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查设备状态，如果在线则不允许删除
	if existingDevice.Status == device.DeviceStatusOnline {
		return fmt.Errorf("在线设备不能删除，请先下线设备")
	}

	// 删除设备
	if err := s.deviceRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除设备失败",
			zap.Int64("id", id),
			zap.String("device_id", existingDevice.DeviceID),
			zap.Error(err))
		return fmt.Errorf("删除设备失败: %w", err)
	}

	s.logger.Info("设备删除成功",
		zap.Int64("id", id),
		zap.String("device_id", existingDevice.DeviceID),
		zap.String("name", existingDevice.Name))

	return nil
}

// ListDevices 获取设备列表
func (s *deviceServiceImpl) ListDevices(ctx context.Context, filter *device.DeviceFilter, sort *device.DeviceSortOption, pagination *device.PaginationOption) (*device.DeviceListResponse, error) {
	// 设置默认分页
	if pagination == nil {
		pagination = &device.PaginationOption{
			Page: 1,
			Size: 20,
		}
	}

	// 设置默认排序
	if sort == nil {
		sort = &device.DeviceSortOption{
			Field: "created_at",
			Order: "desc",
		}
	}

	// 获取设备列表
	devices, total, err := s.deviceRepo.List(ctx, filter, sort, pagination)
	if err != nil {
		s.logger.Error("获取设备列表失败",
			zap.Error(err))
		return nil, fmt.Errorf("获取设备列表失败: %w", err)
	}

	// 计算分页信息
	totalPages := (int(total) + pagination.Size - 1) / pagination.Size
	paginationInfo := device.PaginationInfo{
		Page:       pagination.Page,
		Size:       pagination.Size,
		Total:      total,
		TotalPages: totalPages,
	}

	// 获取统计信息（可选）
	var statistics *device.DeviceStatistics
	if filter != nil {
		stats, err := s.deviceRepo.GetStatistics(ctx, filter)
		if err != nil {
			s.logger.Warn("获取设备统计失败，继续返回列表",
				zap.Error(err))
		} else {
			statistics = stats
		}
	}

	response := &device.DeviceListResponse{
		Devices:    devices,
		Pagination: paginationInfo,
		Statistics: statistics,
	}

	s.logger.Info("获取设备列表成功",
		zap.Int("count", len(devices)),
		zap.Int64("total", total),
		zap.Int("page", pagination.Page))

	return response, nil
}

// DeviceOnline 设备上线
func (s *deviceServiceImpl) DeviceOnline(ctx context.Context, deviceID string) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}

	// 更新设备状态为在线
	if err := s.deviceRepo.UpdateStatus(ctx, deviceID, device.DeviceStatusOnline); err != nil {
		s.logger.Error("设备上线失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return fmt.Errorf("设备上线失败: %w", err)
	}

	// 更新最后活跃时间
	if err := s.deviceRepo.UpdateLastSeen(ctx, deviceID, time.Now()); err != nil {
		s.logger.Warn("更新设备最后活跃时间失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
	}

	s.logger.Info("设备上线成功",
		zap.String("device_id", deviceID))

	return nil
}

// DeviceOffline 设备下线
func (s *deviceServiceImpl) DeviceOffline(ctx context.Context, deviceID string) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}

	// 更新设备状态为离线
	if err := s.deviceRepo.UpdateStatus(ctx, deviceID, device.DeviceStatusOffline); err != nil {
		s.logger.Error("设备下线失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return fmt.Errorf("设备下线失败: %w", err)
	}

	s.logger.Info("设备下线成功",
		zap.String("device_id", deviceID))

	return nil
}

// UpdateDeviceStatus 更新设备状态
func (s *deviceServiceImpl) UpdateDeviceStatus(ctx context.Context, deviceID string, status string) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}

	// 获取现有设备并校验状态转换
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return err
	}
	if !s.deviceBiz.CanTransitionTo(dev, status) {
		return fmt.Errorf("不允许的状态转换: %s -> %s", dev.Status, status)
	}

	// 验证状态值
	validStatuses := []string{
		device.DeviceStatusOnline,
		device.DeviceStatusOffline,
		device.DeviceStatusMaintenance,
		device.DeviceStatusError,
	}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的设备状态: %s", status)
	}

	// 更新设备状态
	if err := s.deviceRepo.UpdateStatus(ctx, deviceID, status); err != nil {
		s.logger.Error("更新设备状态失败",
			zap.String("device_id", deviceID),
			zap.String("status", status),
			zap.Error(err))
		return fmt.Errorf("更新设备状态失败: %w", err)
	}

	// 如果状态变为在线，更新最后活跃时间
	if status == device.DeviceStatusOnline {
		if err := s.deviceRepo.UpdateLastSeen(ctx, deviceID, time.Now()); err != nil {
			s.logger.Warn("更新设备最后活跃时间失败",
				zap.String("device_id", deviceID),
				zap.Error(err))
		}
	}

	s.logger.Info("设备状态更新成功",
		zap.String("device_id", deviceID),
		zap.String("status", status))

	return nil
}

// GetDeviceStatus 获取设备状态
func (s *deviceServiceImpl) GetDeviceStatus(ctx context.Context, deviceID string) (*device.DeviceStatusResponse, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 判断设备是否在线（基于最后活跃时间）
	online := false
	if dev.Status == device.DeviceStatusOnline && dev.LastSeen != nil {
		// 如果最后活跃时间在5分钟内，认为在线
		if time.Since(*dev.LastSeen) <= 5*time.Minute {
			online = true
		}
	}

	response := &device.DeviceStatusResponse{
		DeviceID:    dev.DeviceID,
		Status:      dev.Status,
		HealthScore: dev.HealthScore,
		LastSeen:    dev.LastSeen,
		Uptime:      dev.UptimeHours,
		Online:      online,
	}

	return response, nil
}

// UpdateHealthScore 更新设备健康度
func (s *deviceServiceImpl) UpdateHealthScore(ctx context.Context, deviceID string, score int) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}
	if score < 0 || score > 100 {
		return fmt.Errorf("健康度必须在0-100之间")
	}

	if err := s.deviceRepo.UpdateHealthScore(ctx, deviceID, score); err != nil {
		s.logger.Error("更新设备健康度失败",
			zap.String("device_id", deviceID),
			zap.Int("score", score),
			zap.Error(err))
		return fmt.Errorf("更新设备健康度失败: %w", err)
	}

	s.logger.Info("设备健康度更新成功",
		zap.String("device_id", deviceID),
		zap.Int("score", score))

	return nil
}

// BatchUpdateStatus 批量更新设备状态
func (s *deviceServiceImpl) BatchUpdateStatus(ctx context.Context, deviceIDs []string, status string) error {
	if len(deviceIDs) == 0 {
		return fmt.Errorf("设备ID列表不能为空")
	}

	// 验证状态值
	validStatuses := []string{
		device.DeviceStatusOnline,
		device.DeviceStatusOffline,
		device.DeviceStatusMaintenance,
		device.DeviceStatusError,
	}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的设备状态: %s", status)
	}

	if err := s.deviceRepo.BatchUpdateStatus(ctx, deviceIDs, status); err != nil {
		s.logger.Error("批量更新设备状态失败",
			zap.Strings("device_ids", deviceIDs),
			zap.String("status", status),
			zap.Error(err))
		return fmt.Errorf("批量更新设备状态失败: %w", err)
	}

	s.logger.Info("批量更新设备状态成功",
		zap.Int("count", len(deviceIDs)),
		zap.String("status", status))

	return nil
}

// ProcessHeartbeat 处理设备心跳
func (s *deviceServiceImpl) ProcessHeartbeat(ctx context.Context, deviceID string, req *device.DeviceHeartbeatRequest) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}
	if req.Status == "" {
		return fmt.Errorf("心跳状态不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}

	// 创建心跳记录
	heartbeat := &device.DeviceHeartbeat{
		DeviceID:      dev.ID,
		HeartbeatTime: time.Now(),
		Status:        req.Status,
		ResponseTime:  req.ResponseTime,
	}

	if req.Metrics != nil {
		metrics := device.JSONB(req.Metrics)
		heartbeat.Metrics = &metrics
	}
	if req.SystemInfo != nil {
		systemInfo := device.JSONB(req.SystemInfo)
		heartbeat.SystemInfo = &systemInfo
	}

	// 设置IP地址（从设备信息中获取）
	if dev.IPAddress != nil {
		heartbeat.IPAddress = dev.IPAddress
	}

	// 创建心跳记录
	if err := s.deviceRepo.CreateHeartbeat(ctx, heartbeat); err != nil {
		s.logger.Error("创建心跳记录失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return fmt.Errorf("创建心跳记录失败: %w", err)
	}

	// 根据业务逻辑更新设备状态与健康度
	newStatus := s.deviceBiz.DetermineDeviceStatus(dev, req)
	if newStatus != dev.Status {
		_ = s.deviceRepo.UpdateStatus(ctx, deviceID, newStatus)
	}
	newScore := s.deviceBiz.CalculateHealthScore(dev)
	_ = s.deviceRepo.UpdateHealthScore(ctx, deviceID, newScore)

	// 更新设备状态和最后活跃时间
	if err := s.deviceRepo.UpdateStatus(ctx, deviceID, req.Status); err != nil {
		s.logger.Warn("更新设备状态失败",
			zap.String("device_id", deviceID),
			zap.String("status", req.Status),
			zap.Error(err))
	}

	if err := s.deviceRepo.UpdateLastSeen(ctx, deviceID, time.Now()); err != nil {
		s.logger.Warn("更新设备最后活跃时间失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
	}

	s.logger.Info("处理设备心跳成功",
		zap.String("device_id", deviceID),
		zap.String("status", req.Status))

	return nil
}

// GetDeviceHeartbeat 获取设备心跳
func (s *deviceServiceImpl) GetDeviceHeartbeat(ctx context.Context, deviceID string) (*device.DeviceHeartbeat, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 获取最新心跳
	heartbeat, err := s.deviceRepo.GetLatestHeartbeat(ctx, dev.ID)
	if err != nil {
		s.logger.Error("获取设备心跳失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return nil, err
	}

	return heartbeat, nil
}

// GetHeartbeatHistory 获取心跳历史
func (s *deviceServiceImpl) GetHeartbeatHistory(ctx context.Context, deviceID string, hours int) ([]device.DeviceHeartbeat, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}
	if hours <= 0 {
		hours = 24 // 默认24小时
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 获取心跳历史
	heartbeats, err := s.deviceRepo.GetHeartbeatHistory(ctx, dev.ID, hours)
	if err != nil {
		s.logger.Error("获取心跳历史失败",
			zap.String("device_id", deviceID),
			zap.Int("hours", hours),
			zap.Error(err))
		return nil, err
	}

	return heartbeats, nil
}

// CheckOfflineDevices 检查离线设备
func (s *deviceServiceImpl) CheckOfflineDevices(ctx context.Context, duration time.Duration) ([]device.Device, error) {
	if duration <= 0 {
		duration = 5 * time.Minute // 默认5分钟
	}

	devices, err := s.deviceRepo.GetOfflineDevices(ctx, duration)
	if err != nil {
		s.logger.Error("检查离线设备失败",
			zap.Duration("duration", duration),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("检查离线设备完成",
		zap.Int("offline_count", len(devices)),
		zap.Duration("duration", duration))

	return devices, nil
}

// SendCommand 发送设备命令
func (s *deviceServiceImpl) SendCommand(ctx context.Context, deviceID string, req *device.DeviceCommandRequest) (*device.DeviceCommand, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}
	if req.CommandType == "" {
		return nil, fmt.Errorf("命令类型不能为空")
	}
	if req.CommandData == nil {
		return nil, fmt.Errorf("命令数据不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 检查设备是否在线
	if dev.Status != device.DeviceStatusOnline {
		return nil, fmt.Errorf("设备不在线，无法发送命令")
	}

	// 创建命令对象
	command := &device.DeviceCommand{
		DeviceID:    dev.ID,
		CommandType: req.CommandType,
		CommandData: device.JSONB(req.CommandData),
		Status:      "pending",
		SentTime:    time.Now(),
	}

	// 创建命令记录
	if err := s.deviceRepo.CreateCommand(ctx, command); err != nil {
		s.logger.Error("创建设备命令失败",
			zap.String("device_id", deviceID),
			zap.String("command_type", req.CommandType),
			zap.Error(err))
		return nil, fmt.Errorf("创建设备命令失败: %w", err)
	}

	s.logger.Info("设备命令发送成功",
		zap.String("device_id", deviceID),
		zap.String("command_id", command.CommandID),
		zap.String("command_type", req.CommandType))

	return command, nil
}

// GetCommand 获取设备命令
func (s *deviceServiceImpl) GetCommand(ctx context.Context, commandID string) (*device.DeviceCommand, error) {
	if commandID == "" {
		return nil, fmt.Errorf("命令ID不能为空")
	}

	command, err := s.deviceRepo.GetCommand(ctx, commandID)
	if err != nil {
		s.logger.Error("获取设备命令失败",
			zap.String("command_id", commandID),
			zap.Error(err))
		return nil, err
	}

	return command, nil
}

// GetDeviceCommands 获取设备命令列表
func (s *deviceServiceImpl) GetDeviceCommands(ctx context.Context, deviceID string, status []string) ([]device.DeviceCommand, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 获取命令列表
	commands, err := s.deviceRepo.GetDeviceCommands(ctx, dev.ID, status)
	if err != nil {
		s.logger.Error("获取设备命令列表失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return nil, err
	}

	return commands, nil
}

// UpdateCommandStatus 更新命令状态
func (s *deviceServiceImpl) UpdateCommandStatus(ctx context.Context, commandID string, status string) error {
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	// 验证状态值
	validStatuses := []string{"pending", "executing", "completed", "failed", "timeout"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的命令状态: %s", status)
	}

	if err := s.deviceRepo.UpdateCommandStatus(ctx, commandID, status); err != nil {
		s.logger.Error("更新命令状态失败",
			zap.String("command_id", commandID),
			zap.String("status", status),
			zap.Error(err))
		return fmt.Errorf("更新命令状态失败: %w", err)
	}

	s.logger.Info("命令状态更新成功",
		zap.String("command_id", commandID),
		zap.String("status", status))

	return nil
}

// UpdateCommandResponse 更新命令响应
func (s *deviceServiceImpl) UpdateCommandResponse(ctx context.Context, commandID string, response device.JSONB, errorMsg *string) error {
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	if err := s.deviceRepo.UpdateCommandResponse(ctx, commandID, response, errorMsg); err != nil {
		s.logger.Error("更新命令响应失败",
			zap.String("command_id", commandID),
			zap.Error(err))
		return fmt.Errorf("更新命令响应失败: %w", err)
	}

	s.logger.Info("命令响应更新成功",
		zap.String("command_id", commandID))

	return nil
}

// CreateDeviceLog 创建设备日志
func (s *deviceServiceImpl) CreateDeviceLog(ctx context.Context, deviceID string, req *device.DeviceLogRequest) error {
	if deviceID == "" {
		return fmt.Errorf("设备ID不能为空")
	}
	if req.Level == "" {
		return fmt.Errorf("日志级别不能为空")
	}
	if req.Message == "" {
		return fmt.Errorf("日志消息不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}

	// 验证日志级别
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	isValid := false
	for _, validLevel := range validLevels {
		if req.Level == validLevel {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的日志级别: %s", req.Level)
	}

	// 创建日志对象
	log := &device.DeviceLog{
		DeviceID:      dev.ID,
		LogTime:       time.Now(),
		Level:         req.Level,
		Category:      req.Category,
		Message:       req.Message,
		Source:        req.Source,
		CorrelationID: req.CorrelationID,
	}

	if req.Details != nil {
		details := device.JSONB(req.Details)
		log.Details = &details
	}

	// 创建日志记录
	if err := s.deviceRepo.CreateLog(ctx, log); err != nil {
		s.logger.Error("创建设备日志失败",
			zap.String("device_id", deviceID),
			zap.String("level", req.Level),
			zap.Error(err))
		return fmt.Errorf("创建设备日志失败: %w", err)
	}

	s.logger.Info("设备日志创建成功",
		zap.String("device_id", deviceID),
		zap.String("level", req.Level),
		zap.String("message", req.Message))

	return nil
}

// GetDeviceLogs 获取设备日志
func (s *deviceServiceImpl) GetDeviceLogs(ctx context.Context, deviceID string, filter *device.LogFilter) ([]device.DeviceLog, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("设备ID不能为空")
	}

	// 获取设备信息
	dev, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	// 获取日志列表
	logs, err := s.deviceRepo.GetLogs(ctx, dev.ID, filter)
	if err != nil {
		s.logger.Error("获取设备日志失败",
			zap.String("device_id", deviceID),
			zap.Error(err))
		return nil, err
	}

	return logs, nil
}

// GetLogsByLevel 根据级别获取日志
func (s *deviceServiceImpl) GetLogsByLevel(ctx context.Context, level string, limit int) ([]device.DeviceLog, error) {
	if level == "" {
		return nil, fmt.Errorf("日志级别不能为空")
	}
	if limit <= 0 {
		limit = 100 // 默认100条
	}

	logs, err := s.deviceRepo.GetLogsByLevel(ctx, level, limit)
	if err != nil {
		s.logger.Error("根据级别获取日志失败",
			zap.String("level", level),
			zap.Int("limit", limit),
			zap.Error(err))
		return nil, err
	}

	return logs, nil
}

// QueryDevices 查询设备
func (s *deviceServiceImpl) QueryDevices(ctx context.Context, keyword string, filter *device.DeviceFilter) ([]device.Device, error) {
	if keyword == "" {
		return nil, fmt.Errorf("查询关键字不能为空")
	}

	// 将关键字添加到过滤条件中
	if filter == nil {
		filter = &device.DeviceFilter{}
	}
	filter.Keyword = &keyword

	// 使用列表查询（无分页）
	devices, _, err := s.deviceRepo.List(ctx, filter, nil, nil)
	if err != nil {
		s.logger.Error("查询设备失败",
			zap.String("keyword", keyword),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("查询设备完成",
		zap.String("keyword", keyword),
		zap.Int("count", len(devices)))

	return devices, nil
}

// GetDeviceStats 获取设备统计
func (s *deviceServiceImpl) GetDeviceStats(ctx context.Context, filter *device.DeviceFilter) (*device.DeviceStatistics, error) {
	stats, err := s.deviceRepo.GetStatistics(ctx, filter)
	if err != nil {
		s.logger.Error("获取设备统计失败",
			zap.Error(err))
		return nil, err
	}

	return stats, nil
}

// GetOnlineDevices 获取在线设备
func (s *deviceServiceImpl) GetOnlineDevices(ctx context.Context) ([]device.Device, error) {
	devices, err := s.deviceRepo.GetOnlineDevices(ctx)
	if err != nil {
		s.logger.Error("获取在线设备失败",
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("获取在线设备完成",
		zap.Int("count", len(devices)))

	return devices, nil
}

// GetHealthReport 获取设备健康报告
func (s *deviceServiceImpl) GetHealthReport(ctx context.Context) (map[string]interface{}, error) {
	report, err := s.deviceRepo.GetHealthReport(ctx)
	if err != nil {
		s.logger.Error("获取设备健康报告失败",
			zap.Error(err))
		return nil, err
	}

	return report, nil
}

// CleanupOldHeartbeats 清理过期心跳
func (s *deviceServiceImpl) CleanupOldHeartbeats(ctx context.Context, days int) error {
	if days <= 0 {
		days = 30 // 默认30天
	}

	if err := s.deviceRepo.CleanupOldHeartbeats(ctx, days); err != nil {
		s.logger.Error("清理过期心跳失败",
			zap.Int("days", days),
			zap.Error(err))
		return fmt.Errorf("清理过期心跳失败: %w", err)
	}

	s.logger.Info("清理过期心跳成功",
		zap.Int("days", days))

	return nil
}

// CleanupOldLogs 清理过期日志
func (s *deviceServiceImpl) CleanupOldLogs(ctx context.Context, days int) error {
	if days <= 0 {
		days = 90 // 默认90天
	}

	if err := s.deviceRepo.CleanupOldLogs(ctx, days); err != nil {
		s.logger.Error("清理过期日志失败",
			zap.Int("days", days),
			zap.Error(err))
		return fmt.Errorf("清理过期日志失败: %w", err)
	}

	s.logger.Info("清理过期日志成功",
		zap.Int("days", days))

	return nil
}

// Service 服务容器
type Service struct {
	TaskRepository   repository.TaskRepository
	DeviceRepository repository.DeviceRepository
}
