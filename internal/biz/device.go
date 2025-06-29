package biz

import (
	"errors"
	"time"
)

// DeviceRepo 设备仓储接口
type DeviceRepo interface {
	CreateDevice(device *Device) error
	GetDeviceByID(id uint) (*Device, error)
	GetDeviceByDeviceID(deviceID string) (*Device, error)
	UpdateDevice(device *Device) error
	DeleteDevice(id uint) error
	ListDevices(offset, limit int, filter *DeviceFilter) ([]*Device, int64, error)
	ListOfflineDevices(offlineThreshold time.Duration) ([]*Device, error)
}

// DeviceFilter 设备查询过滤器
type DeviceFilter struct {
	UserID *uint
	Status *DeviceStatus
	Type   string
	Name   string
}

// DeviceUsecase 设备用例
type DeviceUsecase struct {
	repo              DeviceRepo
	alertUsecase      *AlertUsecase
	heartbeatInterval time.Duration
	offlineTimeout    time.Duration
}

// NewDeviceUsecase 创建设备用例
func NewDeviceUsecase(repo DeviceRepo, alertUsecase *AlertUsecase) *DeviceUsecase {
	return &DeviceUsecase{
		repo:              repo,
		alertUsecase:      alertUsecase,
		heartbeatInterval: 30 * time.Second,
		offlineTimeout:    5 * time.Minute,
	}
}

// RegisterDevice 注册设备
func (uc *DeviceUsecase) RegisterDevice(userID uint, deviceID, name, deviceType, ip, mac, os, version string) (*Device, error) {
	// 验证输入
	if err := uc.validateDeviceInput(deviceID, name); err != nil {
		return nil, err
	}

	// 检查设备ID是否已存在
	if _, err := uc.repo.GetDeviceByDeviceID(deviceID); err == nil {
		return nil, errors.New("设备ID已存在")
	}

	// 创建设备实体
	device := &Device{
		DeviceID:   deviceID,
		Name:       name,
		Type:       deviceType,
		Status:     DeviceStatusOnline,
		IP:         ip,
		MAC:        mac,
		OS:         os,
		Version:    version,
		UserID:     userID,
		LastSeenAt: &time.Time{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 设置初始心跳时间
	now := time.Now()
	device.LastSeenAt = &now

	// 保存到数据库
	if err := uc.repo.CreateDevice(device); err != nil {
		return nil, err
	}

	return device, nil
}

// GetDevice 获取设备信息
func (uc *DeviceUsecase) GetDevice(id uint) (*Device, error) {
	return uc.repo.GetDeviceByID(id)
}

// GetDeviceByDeviceID 通过设备ID获取设备
func (uc *DeviceUsecase) GetDeviceByDeviceID(deviceID string) (*Device, error) {
	return uc.repo.GetDeviceByDeviceID(deviceID)
}

// UpdateDevice 更新设备信息
func (uc *DeviceUsecase) UpdateDevice(id uint, updates map[string]interface{}) (*Device, error) {
	device, err := uc.repo.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if name, ok := updates["name"].(string); ok {
		device.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		device.Description = description
	}
	if deviceType, ok := updates["type"].(string); ok {
		device.Type = deviceType
	}
	if status, ok := updates["status"].(DeviceStatus); ok {
		device.Status = status
	}
	if ip, ok := updates["ip"].(string); ok {
		device.IP = ip
	}
	if version, ok := updates["version"].(string); ok {
		device.Version = version
	}

	device.UpdatedAt = time.Now()

	if err := uc.repo.UpdateDevice(device); err != nil {
		return nil, err
	}

	return device, nil
}

// ProcessHeartbeat 处理设备心跳
func (uc *DeviceUsecase) ProcessHeartbeat(deviceID string, metadata map[string]interface{}) error {
	device, err := uc.repo.GetDeviceByDeviceID(deviceID)
	if err != nil {
		return errors.New("设备不存在")
	}

	// 更新最后心跳时间
	now := time.Now()
	device.LastSeenAt = &now
	device.UpdatedAt = now

	// 如果设备之前是离线状态，现在变为在线
	if device.Status == DeviceStatusOffline {
		device.Status = DeviceStatusOnline

		// 创建设备上线告警
		uc.alertUsecase.CreateAlert(AlertTypeDevice, AlertLevelInfo,
			"设备上线",
			"设备 "+device.Name+" 已重新上线",
			device.ID, device.UserID, nil)
	}

	// 更新设备信息（如果心跳包含额外信息）
	if ip, ok := metadata["ip"].(string); ok && ip != "" {
		device.IP = ip
	}
	if version, ok := metadata["version"].(string); ok && version != "" {
		device.Version = version
	}
	if os, ok := metadata["os"].(string); ok && os != "" {
		device.OS = os
	}

	return uc.repo.UpdateDevice(device)
}

// MarkDeviceOffline 标记设备离线
func (uc *DeviceUsecase) MarkDeviceOffline(deviceID string) error {
	device, err := uc.repo.GetDeviceByDeviceID(deviceID)
	if err != nil {
		return err
	}

	if device.Status != DeviceStatusOffline {
		device.Status = DeviceStatusOffline
		device.UpdatedAt = time.Now()

		// 创建设备离线告警
		uc.alertUsecase.CreateAlert(AlertTypeDevice, AlertLevelWarning,
			"设备离线",
			"设备 "+device.Name+" 已离线",
			device.ID, device.UserID, nil)

		return uc.repo.UpdateDevice(device)
	}

	return nil
}

// SetDeviceStatus 设置设备状态
func (uc *DeviceUsecase) SetDeviceStatus(id uint, status DeviceStatus) error {
	device, err := uc.repo.GetDeviceByID(id)
	if err != nil {
		return err
	}

	oldStatus := device.Status
	device.Status = status
	device.UpdatedAt = time.Now()

	if err := uc.repo.UpdateDevice(device); err != nil {
		return err
	}

	// 如果状态发生变化，创建相应告警
	if oldStatus != status {
		var level AlertLevel
		var message string

		switch status {
		case DeviceStatusOnline:
			level = AlertLevelInfo
			message = "设备状态变更为在线"
		case DeviceStatusOffline:
			level = AlertLevelWarning
			message = "设备状态变更为离线"
		case DeviceStatusMaintenance:
			level = AlertLevelInfo
			message = "设备进入维护模式"
		case DeviceStatusError:
			level = AlertLevelError
			message = "设备状态异常"
		}

		uc.alertUsecase.CreateAlert(AlertTypeDevice, level,
			"设备状态变更",
			message,
			device.ID, device.UserID, nil)
	}

	return nil
}

// DeleteDevice 删除设备
func (uc *DeviceUsecase) DeleteDevice(id uint) error {
	// TODO: 检查是否有正在运行的任务
	// TODO: 清理相关数据
	return uc.repo.DeleteDevice(id)
}

// ListDevices 设备列表
func (uc *DeviceUsecase) ListDevices(page, pageSize int, filter *DeviceFilter) ([]*Device, int64, error) {
	offset := (page - 1) * pageSize
	return uc.repo.ListDevices(offset, pageSize, filter)
}

// ListUserDevices 获取用户的设备列表
func (uc *DeviceUsecase) ListUserDevices(userID uint, page, pageSize int) ([]*Device, int64, error) {
	filter := &DeviceFilter{UserID: &userID}
	offset := (page - 1) * pageSize
	return uc.repo.ListDevices(offset, pageSize, filter)
}

// CheckOfflineDevices 检查离线设备
func (uc *DeviceUsecase) CheckOfflineDevices() error {
	devices, err := uc.repo.ListOfflineDevices(uc.offlineTimeout)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if device.Status == DeviceStatusOnline {
			uc.MarkDeviceOffline(device.DeviceID)
		}
	}

	return nil
}

// GetDeviceStats 获取设备统计信息
func (uc *DeviceUsecase) GetDeviceStats(userID *uint) (map[string]int64, error) {
	filter := &DeviceFilter{}
	if userID != nil {
		filter.UserID = userID
	}

	// 获取全部设备进行统计
	devices, _, err := uc.repo.ListDevices(0, 10000, filter)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	stats["total"] = int64(len(devices))
	stats["online"] = 0
	stats["offline"] = 0
	stats["maintenance"] = 0
	stats["error"] = 0

	for _, device := range devices {
		switch device.Status {
		case DeviceStatusOnline:
			stats["online"]++
		case DeviceStatusOffline:
			stats["offline"]++
		case DeviceStatusMaintenance:
			stats["maintenance"]++
		case DeviceStatusError:
			stats["error"]++
		}
	}

	return stats, nil
}

// 私有方法

// validateDeviceInput 验证设备输入
func (uc *DeviceUsecase) validateDeviceInput(deviceID, name string) error {
	if len(deviceID) == 0 {
		return errors.New("设备ID不能为空")
	}
	if len(name) == 0 {
		return errors.New("设备名称不能为空")
	}
	return nil
}
