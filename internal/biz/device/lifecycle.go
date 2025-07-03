package device

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// DeviceLifecycleManager 设备生命周期管理器
type DeviceLifecycleManager struct {
	business  *DeviceBusiness
	validator *DeviceValidator
	mutex     sync.RWMutex

	// 状态监听器
	stateListeners map[string][]DeviceStateChangeListener

	// 心跳监控
	heartbeatMonitors map[int64]*HeartbeatMonitor

	// 设备连接池
	connectionPool map[int64]*DeviceConnection

	// 设备分组
	deviceGroups map[string][]int64 // groupName -> deviceIDs
}

// DeviceStateChangeListener 设备状态变更监听器
type DeviceStateChangeListener func(ctx context.Context, device *Device, oldStatus, newStatus string) error

// HeartbeatMonitor 心跳监控器
type HeartbeatMonitor struct {
	DeviceID        int64
	LastHeartbeat   time.Time
	MissedHeartbeat int
	MaxMissed       int
	Interval        time.Duration
	IsActive        bool
	stopChan        chan struct{}
}

// DeviceConnection 设备连接信息
type DeviceConnection struct {
	DeviceID          int64
	Connected         bool
	LastConnected     time.Time
	ConnectionCount   int64
	FailedAttempts    int
	MaxFailedAttempts int
}

// DeviceLifecycleInfo 设备生命周期信息
type DeviceLifecycleInfo struct {
	DeviceID           int64      `json:"device_id"`
	CurrentStatus      string     `json:"current_status"`
	PossibleStatuses   []string   `json:"possible_statuses"`
	HealthScore        int        `json:"health_score"`
	LastSeen           *time.Time `json:"last_seen,omitempty"`
	ConnectionStatus   string     `json:"connection_status"`
	HeartbeatInterval  string     `json:"heartbeat_interval"`
	NextHeartbeat      *time.Time `json:"next_heartbeat,omitempty"`
	UpTime             float64    `json:"uptime_hours"`
	CommandSuccessRate float64    `json:"command_success_rate"`
	IsHealthy          bool       `json:"is_healthy"`
	IsOverloaded       bool       `json:"is_overloaded"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// NewDeviceLifecycleManager 创建设备生命周期管理器
func NewDeviceLifecycleManager() *DeviceLifecycleManager {
	return &DeviceLifecycleManager{
		business:          NewDeviceBusiness(),
		validator:         NewDeviceValidator(),
		stateListeners:    make(map[string][]DeviceStateChangeListener),
		heartbeatMonitors: make(map[int64]*HeartbeatMonitor),
		connectionPool:    make(map[int64]*DeviceConnection),
		deviceGroups:      make(map[string][]int64),
	}
}

// TransitionDevice 转换设备状态
func (m *DeviceLifecycleManager) TransitionDevice(ctx context.Context, device *Device, targetStatus string) error {
	if device == nil {
		return errors.New("设备不能为空")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	oldStatus := device.Status

	// 验证状态转换
	if err := m.validator.ValidateStatusTransition(oldStatus, targetStatus); err != nil {
		return fmt.Errorf("状态转换验证失败: %w", err)
	}

	// 执行状态转换
	device.Status = targetStatus
	device.UpdatedAt = time.Now()

	// 更新健康度
	device.HealthScore = m.business.CalculateHealthScore(device)

	// 触发状态变更监听器
	if err := m.triggerStateListeners(ctx, device, oldStatus, targetStatus); err != nil {
		// 回滚状态
		device.Status = oldStatus
		return fmt.Errorf("状态监听器执行失败: %w", err)
	}

	// 更新连接状态
	m.updateConnectionStatus(device, targetStatus)

	// 更新心跳监控
	m.updateHeartbeatMonitor(device, targetStatus)

	return nil
}

// StartHeartbeatMonitoring 启动心跳监控
func (m *DeviceLifecycleManager) StartHeartbeatMonitoring(device *Device) error {
	if device == nil {
		return errors.New("设备不能为空")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已存在监控器
	if monitor, exists := m.heartbeatMonitors[device.ID]; exists && monitor.IsActive {
		return nil // 已经在监控中
	}

	// 创建新的心跳监控器
	monitor := &HeartbeatMonitor{
		DeviceID:        device.ID,
		LastHeartbeat:   time.Now(),
		MissedHeartbeat: 0,
		MaxMissed:       3, // 连续3次心跳丢失后认为设备离线
		Interval:        m.business.GetExpectedHeartbeatInterval(device),
		IsActive:        true,
		stopChan:        make(chan struct{}),
	}

	m.heartbeatMonitors[device.ID] = monitor

	// 启动监控协程
	go m.runHeartbeatMonitor(monitor)

	return nil
}

// StopHeartbeatMonitoring 停止心跳监控
func (m *DeviceLifecycleManager) StopHeartbeatMonitoring(deviceID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if monitor, exists := m.heartbeatMonitors[deviceID]; exists {
		monitor.IsActive = false
		close(monitor.stopChan)
		delete(m.heartbeatMonitors, deviceID)
	}
}

// ProcessHeartbeat 处理设备心跳
func (m *DeviceLifecycleManager) ProcessHeartbeat(ctx context.Context, device *Device, heartbeatReq *DeviceHeartbeatRequest) error {
	if device == nil || heartbeatReq == nil {
		return errors.New("设备或心跳数据不能为空")
	}

	// 验证心跳数据
	if err := m.validator.ValidateHeartbeat(heartbeatReq); err != nil {
		return fmt.Errorf("心跳数据验证失败: %w", err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 更新设备信息
	device.LastSeen = &[]time.Time{time.Now()}[0]
	device.UpdatedAt = time.Now()

	// 更新设备状态
	newStatus := m.business.DetermineDeviceStatus(device, heartbeatReq)
	if newStatus != device.Status {
		if err := m.TransitionDevice(ctx, device, newStatus); err != nil {
			return fmt.Errorf("状态转换失败: %w", err)
		}
	}

	// 更新心跳监控
	if monitor, exists := m.heartbeatMonitors[device.ID]; exists {
		monitor.LastHeartbeat = time.Now()
		monitor.MissedHeartbeat = 0
	}

	// 更新连接状态
	if conn, exists := m.connectionPool[device.ID]; exists {
		conn.Connected = true
		conn.LastConnected = time.Now()
		conn.FailedAttempts = 0
	}

	return nil
}

// AddDeviceToGroup 将设备添加到分组
func (m *DeviceLifecycleManager) AddDeviceToGroup(deviceID int64, groupName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.deviceGroups[groupName] == nil {
		m.deviceGroups[groupName] = make([]int64, 0)
	}

	// 检查设备是否已在分组中
	for _, id := range m.deviceGroups[groupName] {
		if id == deviceID {
			return nil // 已存在
		}
	}

	m.deviceGroups[groupName] = append(m.deviceGroups[groupName], deviceID)
	return nil
}

// RemoveDeviceFromGroup 从分组中移除设备
func (m *DeviceLifecycleManager) RemoveDeviceFromGroup(deviceID int64, groupName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	devices := m.deviceGroups[groupName]
	if devices == nil {
		return nil
	}

	// 移除设备
	newDevices := make([]int64, 0, len(devices))
	for _, id := range devices {
		if id != deviceID {
			newDevices = append(newDevices, id)
		}
	}

	m.deviceGroups[groupName] = newDevices
	return nil
}

// GetDevicesByGroup 获取分组中的设备列表
func (m *DeviceLifecycleManager) GetDevicesByGroup(groupName string) []int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if devices := m.deviceGroups[groupName]; devices != nil {
		result := make([]int64, len(devices))
		copy(result, devices)
		return result
	}
	return []int64{}
}

// AddStateListener 添加状态变更监听器
func (m *DeviceLifecycleManager) AddStateListener(status string, listener DeviceStateChangeListener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.stateListeners[status] == nil {
		m.stateListeners[status] = make([]DeviceStateChangeListener, 0)
	}
	m.stateListeners[status] = append(m.stateListeners[status], listener)
}

// GetDeviceLifecycleInfo 获取设备生命周期信息
func (m *DeviceLifecycleManager) GetDeviceLifecycleInfo(device *Device) *DeviceLifecycleInfo {
	if device == nil {
		return nil
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	info := &DeviceLifecycleInfo{
		DeviceID:           device.ID,
		CurrentStatus:      device.Status,
		PossibleStatuses:   m.getPossibleNextStatuses(device),
		HealthScore:        device.HealthScore,
		LastSeen:           device.LastSeen,
		ConnectionStatus:   m.getConnectionStatus(device.ID),
		HeartbeatInterval:  m.business.GetExpectedHeartbeatInterval(device).String(),
		UpTime:             m.business.CalculateUptime(device),
		CommandSuccessRate: m.business.CalculateCommandSuccessRate(device),
		IsHealthy:          m.business.IsDeviceHealthy(device),
		IsOverloaded:       m.business.IsDeviceOverloaded(device),
		CreatedAt:          device.CreatedAt,
		UpdatedAt:          device.UpdatedAt,
	}

	// 计算下次心跳时间
	if monitor, exists := m.heartbeatMonitors[device.ID]; exists && monitor.IsActive {
		nextHeartbeat := monitor.LastHeartbeat.Add(monitor.Interval)
		info.NextHeartbeat = &nextHeartbeat
	}

	return info
}

// GetDeviceGroupStatistics 获取设备分组统计
func (m *DeviceLifecycleManager) GetDeviceGroupStatistics() map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]int)
	for groupName, devices := range m.deviceGroups {
		stats[groupName] = len(devices)
	}
	return stats
}

// GetHeartbeatStatistics 获取心跳监控统计
func (m *DeviceLifecycleManager) GetHeartbeatStatistics() map[string]int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := map[string]int{
		"total_monitored":   len(m.heartbeatMonitors),
		"active_monitors":   0,
		"missed_heartbeats": 0,
	}

	for _, monitor := range m.heartbeatMonitors {
		if monitor.IsActive {
			stats["active_monitors"]++
		}
		stats["missed_heartbeats"] += monitor.MissedHeartbeat
	}

	return stats
}

// 私有方法

// triggerStateListeners 触发状态监听器
func (m *DeviceLifecycleManager) triggerStateListeners(ctx context.Context, device *Device, oldStatus, newStatus string) error {
	// 触发新状态的监听器
	if listeners := m.stateListeners[newStatus]; listeners != nil {
		for _, listener := range listeners {
			if err := listener(ctx, device, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	// 触发通用监听器
	if listeners := m.stateListeners["*"]; listeners != nil {
		for _, listener := range listeners {
			if err := listener(ctx, device, oldStatus, newStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateConnectionStatus 更新连接状态
func (m *DeviceLifecycleManager) updateConnectionStatus(device *Device, newStatus string) {
	conn, exists := m.connectionPool[device.ID]
	if !exists {
		conn = &DeviceConnection{
			DeviceID:          device.ID,
			MaxFailedAttempts: 3,
		}
		m.connectionPool[device.ID] = conn
	}

	switch newStatus {
	case DeviceStatusOnline:
		conn.Connected = true
		conn.LastConnected = time.Now()
		conn.ConnectionCount++
		conn.FailedAttempts = 0
	case DeviceStatusOffline, DeviceStatusError:
		conn.Connected = false
		conn.FailedAttempts++
	}
}

// updateHeartbeatMonitor 更新心跳监控
func (m *DeviceLifecycleManager) updateHeartbeatMonitor(device *Device, newStatus string) {
	monitor, exists := m.heartbeatMonitors[device.ID]
	if !exists {
		return
	}

	switch newStatus {
	case DeviceStatusOffline:
		monitor.IsActive = false
	case DeviceStatusOnline:
		monitor.IsActive = true
		monitor.MissedHeartbeat = 0
		monitor.LastHeartbeat = time.Now()
	}
}

// runHeartbeatMonitor 运行心跳监控
func (m *DeviceLifecycleManager) runHeartbeatMonitor(monitor *HeartbeatMonitor) {
	ticker := time.NewTicker(monitor.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-monitor.stopChan:
			return
		case <-ticker.C:
			if !monitor.IsActive {
				continue
			}

			// 检查心跳超时
			if time.Since(monitor.LastHeartbeat) > monitor.Interval {
				monitor.MissedHeartbeat++

				// 如果连续丢失心跳超过阈值，标记设备为离线
				if monitor.MissedHeartbeat >= monitor.MaxMissed {
					// 这里可以触发设备离线事件
					// 注意：实际实现中可能需要通过回调或事件系统来处理
					monitor.IsActive = false
				}
			}
		}
	}
}

// getPossibleNextStatuses 获取可能的下一步状态
func (m *DeviceLifecycleManager) getPossibleNextStatuses(device *Device) []string {
	switch device.Status {
	case DeviceStatusOffline:
		return []string{DeviceStatusOnline, DeviceStatusMaintenance}
	case DeviceStatusOnline:
		return []string{DeviceStatusOffline, DeviceStatusMaintenance, DeviceStatusError}
	case DeviceStatusMaintenance:
		return []string{DeviceStatusOnline, DeviceStatusOffline}
	case DeviceStatusError:
		return []string{DeviceStatusOnline, DeviceStatusOffline, DeviceStatusMaintenance}
	default:
		return []string{}
	}
}

// getConnectionStatus 获取连接状态描述
func (m *DeviceLifecycleManager) getConnectionStatus(deviceID int64) string {
	conn, exists := m.connectionPool[deviceID]
	if !exists {
		return "未知"
	}

	if conn.Connected {
		return "已连接"
	}

	if conn.FailedAttempts >= conn.MaxFailedAttempts {
		return "连接失败"
	}

	return "未连接"
}
