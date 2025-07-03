package device

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
)

// DeviceBusiness 设备业务逻辑层
type DeviceBusiness struct {
	validator *DeviceValidator
}

// NewDeviceBusiness 创建设备业务逻辑实例
func NewDeviceBusiness() *DeviceBusiness {
	return &DeviceBusiness{
		validator: NewDeviceValidator(),
	}
}

// CalculateDeviceScore 计算设备评分
func (b *DeviceBusiness) CalculateDeviceScore(device *Device) int {
	if device == nil {
		return 0
	}

	// 基础分数（健康度）
	score := device.HealthScore

	// 在线状态加分
	switch device.Status {
	case DeviceStatusOnline:
		score += 20 // 在线设备优先
	case DeviceStatusMaintenance:
		score -= 30 // 维护中设备降级
	case DeviceStatusError:
		score -= 50 // 错误状态严重降级
	case DeviceStatusOffline:
		score -= 80 // 离线设备大幅降级
	}

	// 设备类型权重
	switch device.Type {
	case DeviceTypeGateway:
		score += 15 // 网关设备重要性高
	case DeviceTypeSensor:
		score += 10 // 传感器重要性中等
	case DeviceTypeCamera:
		score += 8 // 摄像头重要性中等
	case DeviceTypeActuator:
		score += 12 // 执行器重要性较高
	}

	// 成功率加分
	if device.TotalCommands > 0 {
		successRate := float64(device.SuccessCommands) / float64(device.TotalCommands)
		score += int(successRate * 20) // 成功率最高加20分
	}

	// 最近活跃度加分
	if device.LastSeen != nil {
		timeSinceLastSeen := time.Since(*device.LastSeen)
		if timeSinceLastSeen < time.Hour {
			score += 15 // 1小时内活跃
		} else if timeSinceLastSeen < 6*time.Hour {
			score += 10 // 6小时内活跃
		} else if timeSinceLastSeen < 24*time.Hour {
			score += 5 // 24小时内活跃
		}
	}

	// 确保分数在合理范围内
	if score < 0 {
		score = 0
	}
	if score > 200 {
		score = 200
	}

	return score
}

// CalculateHealthScore 计算设备健康度
func (b *DeviceBusiness) CalculateHealthScore(device *Device) int {
	if device == nil {
		return 0
	}

	baseScore := 100

	// 状态影响
	switch device.Status {
	case DeviceStatusOnline:
		// 在线状态不扣分
	case DeviceStatusMaintenance:
		baseScore -= 20 // 维护状态轻微扣分
	case DeviceStatusError:
		baseScore -= 60 // 错误状态大幅扣分
	case DeviceStatusOffline:
		baseScore -= 80 // 离线状态严重扣分
	}

	// 命令成功率影响
	if device.TotalCommands > 0 {
		failureRate := float64(device.FailedCommands) / float64(device.TotalCommands)
		baseScore -= int(failureRate * 30) // 失败率影响健康度
	}

	// 最后活跃时间影响
	if device.LastSeen != nil {
		timeSinceLastSeen := time.Since(*device.LastSeen)
		if timeSinceLastSeen > 24*time.Hour {
			baseScore -= 10 // 超过24小时未活跃
		}
		if timeSinceLastSeen > 7*24*time.Hour {
			baseScore -= 20 // 超过7天未活跃
		}
	}

	// 确保健康度在0-100范围内
	if baseScore < 0 {
		baseScore = 0
	}
	if baseScore > 100 {
		baseScore = 100
	}

	return baseScore
}

// DetermineDeviceStatus 确定设备状态
func (b *DeviceBusiness) DetermineDeviceStatus(device *Device, heartbeatData *DeviceHeartbeatRequest) string {
	if device == nil {
		return DeviceStatusOffline
	}

	// 如果有心跳数据，优先使用心跳状态
	if heartbeatData != nil {
		if b.validator.isValidDeviceStatus(heartbeatData.Status) {
			return heartbeatData.Status
		}
	}

	// 基于最后活跃时间判断
	if device.LastSeen != nil {
		timeSinceLastSeen := time.Since(*device.LastSeen)
		if timeSinceLastSeen > 5*time.Minute {
			return DeviceStatusOffline
		}
	}

	// 基于健康度判断
	if device.HealthScore < 30 {
		return DeviceStatusError
	}

	// 默认返回当前状态或在线
	if device.Status != "" {
		return device.Status
	}
	return DeviceStatusOnline
}

// CanTransitionTo 检查设备是否可以转换到目标状态
func (b *DeviceBusiness) CanTransitionTo(device *Device, targetStatus string) bool {
	if device == nil {
		return false
	}

	err := b.validator.ValidateStatusTransition(device.Status, targetStatus)
	return err == nil
}

// CalculateUptime 计算设备运行时间
func (b *DeviceBusiness) CalculateUptime(device *Device) float64 {
	if device == nil || device.LastSeen == nil {
		return 0
	}

	// 如果设备当前在线，计算从创建到现在的运行时间
	if device.Status == DeviceStatusOnline {
		totalTime := time.Since(device.CreatedAt)
		return totalTime.Hours()
	}

	// 否则返回已记录的运行时间
	return device.UptimeHours
}

// IsDeviceHealthy 检查设备是否健康
func (b *DeviceBusiness) IsDeviceHealthy(device *Device) bool {
	if device == nil {
		return false
	}

	// 健康度检查
	if device.HealthScore < 50 {
		return false
	}

	// 状态检查
	if device.Status == DeviceStatusError || device.Status == DeviceStatusOffline {
		return false
	}

	// 最近活跃性检查
	if device.LastSeen != nil {
		timeSinceLastSeen := time.Since(*device.LastSeen)
		if timeSinceLastSeen > 30*time.Minute {
			return false
		}
	}

	return true
}

// ShouldSendHeartbeat 判断是否应该发送心跳
func (b *DeviceBusiness) ShouldSendHeartbeat(device *Device) bool {
	if device == nil {
		return false
	}

	// 离线设备不需要心跳
	if device.Status == DeviceStatusOffline {
		return false
	}

	// 检查最后心跳时间
	if device.LastSeen == nil {
		return true // 没有心跳记录，需要发送
	}

	// 超过心跳间隔时间需要发送
	heartbeatInterval := b.getHeartbeatInterval(device)
	return time.Since(*device.LastSeen) >= heartbeatInterval
}

// GetExpectedHeartbeatInterval 获取设备心跳间隔
func (b *DeviceBusiness) GetExpectedHeartbeatInterval(device *Device) time.Duration {
	return b.getHeartbeatInterval(device)
}

// BuildDeviceCommand 构建设备命令
func (b *DeviceBusiness) BuildDeviceCommand(device *Device, commandType string, commandData map[string]interface{}, createdBy *int64) *DeviceCommand {
	if device == nil {
		return nil
	}

	return &DeviceCommand{
		DeviceID:    device.ID,
		CommandID:   utils.GenerateExecutionID(), // 复用Task的ID生成器
		CommandType: commandType,
		CommandData: JSONB(commandData),
		Status:      "pending",
		SentTime:    time.Now(),
		CreatedBy:   createdBy,
	}
}

// BuildDeviceHeartbeat 构建设备心跳记录
func (b *DeviceBusiness) BuildDeviceHeartbeat(device *Device, heartbeatReq *DeviceHeartbeatRequest) *DeviceHeartbeat {
	if device == nil || heartbeatReq == nil {
		return nil
	}

	heartbeat := &DeviceHeartbeat{
		DeviceID:      device.ID,
		HeartbeatTime: time.Now(),
		Status:        heartbeatReq.Status,
		ResponseTime:  heartbeatReq.ResponseTime,
	}

	if heartbeatReq.Metrics != nil {
		metrics := JSONB(heartbeatReq.Metrics)
		heartbeat.Metrics = &metrics
	}

	if heartbeatReq.SystemInfo != nil {
		systemInfo := JSONB(heartbeatReq.SystemInfo)
		heartbeat.SystemInfo = &systemInfo
	}

	return heartbeat
}

// BuildDeviceLog 构建设备日志记录
func (b *DeviceBusiness) BuildDeviceLog(device *Device, logReq *DeviceLogRequest) *DeviceLog {
	if device == nil || logReq == nil {
		return nil
	}

	log := &DeviceLog{
		DeviceID:      device.ID,
		LogTime:       time.Now(),
		Level:         logReq.Level,
		Category:      logReq.Category,
		Message:       logReq.Message,
		Source:        logReq.Source,
		CorrelationID: logReq.CorrelationID,
	}

	if logReq.Details != nil {
		details := JSONB(logReq.Details)
		log.Details = &details
	}

	return log
}

// CalculateCommandSuccessRate 计算命令成功率
func (b *DeviceBusiness) CalculateCommandSuccessRate(device *Device) float64 {
	if device == nil || device.TotalCommands == 0 {
		return 0
	}

	return float64(device.SuccessCommands) / float64(device.TotalCommands) * 100
}

// GetDeviceCapabilities 获取设备能力列表
func (b *DeviceBusiness) GetDeviceCapabilities(device *Device) []string {
	if device == nil || device.Capabilities == nil {
		return []string{}
	}

	capabilities := []string{}
	capMap := map[string]interface{}(*device.Capabilities)

	for capability, enabled := range capMap {
		if enabledBool, ok := enabled.(bool); ok && enabledBool {
			capabilities = append(capabilities, capability)
		}
	}

	return capabilities
}

// ValidateDeviceConnection 验证设备连接
func (b *DeviceBusiness) ValidateDeviceConnection(device *Device) error {
	if device == nil {
		return errors.New("设备不能为空")
	}

	// 检查必要的连接信息
	if device.IPAddress == nil && device.Endpoint == nil {
		return errors.New("设备必须提供IP地址或端点信息")
	}

	// 检查协议支持
	if device.Protocol != nil && !b.validator.isValidProtocol(*device.Protocol) {
		return fmt.Errorf("不支持的协议: %s", *device.Protocol)
	}

	// 检查认证信息
	if device.AuthType != nil && *device.AuthType != "none" {
		if device.Credentials == nil {
			return errors.New("需要认证的设备必须提供认证信息")
		}
	}

	return nil
}

// EstimateCommandDuration 预估命令执行时长
func (b *DeviceBusiness) EstimateCommandDuration(device *Device, commandType string) time.Duration {
	if device == nil {
		return 30 * time.Second // 默认30秒
	}

	// 基于命令类型预估
	baseDuration := map[string]time.Duration{
		"status":    5 * time.Second,
		"start":     10 * time.Second,
		"stop":      10 * time.Second,
		"restart":   30 * time.Second,
		"configure": 60 * time.Second,
		"reset":     45 * time.Second,
		"update":    300 * time.Second, // 5分钟
		"sync":      120 * time.Second, // 2分钟
	}

	duration, exists := baseDuration[commandType]
	if !exists {
		duration = 30 * time.Second
	}

	// 基于设备性能调整
	if device.HealthScore < 50 {
		duration = time.Duration(float64(duration) * 1.5) // 健康度低的设备执行时间更长
	}

	return duration
}

// GetDeviceStatusDescription 获取设备状态描述
func (b *DeviceBusiness) GetDeviceStatusDescription(status string) string {
	descriptions := map[string]string{
		DeviceStatusOnline:      "设备在线，运行正常",
		DeviceStatusOffline:     "设备离线，无法连接",
		DeviceStatusMaintenance: "设备维护中，暂停服务",
		DeviceStatusError:       "设备异常，需要处理",
	}

	if desc, exists := descriptions[status]; exists {
		return desc
	}
	return "未知状态"
}

// CalculateDeviceLoad 计算设备负载
func (b *DeviceBusiness) CalculateDeviceLoad(device *Device) float64 {
	if device == nil {
		return 0
	}

	// 基于命令处理情况计算负载
	if device.TotalCommands == 0 {
		return 0
	}

	// 计算失败率
	failureRate := float64(device.FailedCommands) / float64(device.TotalCommands)

	// 计算基础负载（基于健康度）
	baseLoad := (100 - float64(device.HealthScore)) / 100

	// 综合负载计算
	totalLoad := baseLoad + failureRate*0.5

	// 确保负载在0-1范围内
	return math.Min(totalLoad, 1.0)
}

// IsDeviceOverloaded 检查设备是否过载
func (b *DeviceBusiness) IsDeviceOverloaded(device *Device) bool {
	load := b.CalculateDeviceLoad(device)
	return load > 0.8 // 负载超过80%认为过载
}

// 私有辅助方法

// getHeartbeatInterval 获取心跳间隔
func (b *DeviceBusiness) getHeartbeatInterval(device *Device) time.Duration {
	if device == nil {
		return 5 * time.Minute // 默认5分钟
	}

	// 基于设备类型设置心跳间隔
	intervals := map[string]time.Duration{
		DeviceTypeSensor:   2 * time.Minute, // 传感器2分钟
		DeviceTypeCamera:   3 * time.Minute, // 摄像头3分钟
		DeviceTypeActuator: 1 * time.Minute, // 执行器1分钟
		DeviceTypeGateway:  5 * time.Minute, // 网关5分钟
	}

	if interval, exists := intervals[device.Type]; exists {
		return interval
	}

	return 5 * time.Minute // 默认5分钟
}
