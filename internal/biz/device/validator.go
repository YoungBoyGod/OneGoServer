package device

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// DeviceValidator 设备验证器
type DeviceValidator struct{}

// NewDeviceValidator 创建设备验证器实例
func NewDeviceValidator() *DeviceValidator {
	return &DeviceValidator{}
}

// ValidateCreate 验证创建设备请求
func (v *DeviceValidator) ValidateCreate(req *DeviceCreateRequest) error {
	if req == nil {
		return errors.New("创建请求不能为空")
	}

	// 必填字段验证
	if strings.TrimSpace(req.DeviceID) == "" {
		return errors.New("设备ID不能为空")
	}
	if len(req.DeviceID) > 100 {
		return errors.New("设备ID不能超过100个字符")
	}
	if !v.isValidDeviceIDFormat(req.DeviceID) {
		return errors.New("设备ID格式无效，只能包含字母、数字、连字符和下划线")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("设备名称不能为空")
	}
	if len(req.Name) > 255 {
		return errors.New("设备名称不能超过255个字符")
	}

	if strings.TrimSpace(req.Type) == "" {
		return errors.New("设备类型不能为空")
	}
	if !v.isValidDeviceType(req.Type) {
		return fmt.Errorf("无效的设备类型: %s", req.Type)
	}

	// 可选字段验证
	if req.Model != nil && len(*req.Model) > 100 {
		return errors.New("设备型号不能超过100个字符")
	}

	if req.Manufacturer != nil && len(*req.Manufacturer) > 100 {
		return errors.New("制造商不能超过100个字符")
	}

	if req.IPAddress != nil {
		if err := v.validateIPAddress(*req.IPAddress); err != nil {
			return fmt.Errorf("IP地址验证失败: %w", err)
		}
	}

	if req.Port != nil {
		if err := v.validatePort(*req.Port); err != nil {
			return fmt.Errorf("端口验证失败: %w", err)
		}
	}

	if req.Protocol != nil && !v.isValidProtocol(*req.Protocol) {
		return fmt.Errorf("无效的协议类型: %s", *req.Protocol)
	}

	if req.Endpoint != nil && len(*req.Endpoint) > 500 {
		return errors.New("端点地址不能超过500个字符")
	}

	if req.AuthType != nil && !v.isValidAuthType(*req.AuthType) {
		return fmt.Errorf("无效的认证类型: %s", *req.AuthType)
	}

	return nil
}

// ValidateUpdate 验证更新设备请求
func (v *DeviceValidator) ValidateUpdate(req *DeviceUpdateRequest) error {
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
			return errors.New("设备名称不能为空")
		}
		if len(*req.Name) > 255 {
			return errors.New("设备名称不能超过255个字符")
		}
	}

	if req.Model != nil && len(*req.Model) > 100 {
		return errors.New("设备型号不能超过100个字符")
	}

	if req.Manufacturer != nil && len(*req.Manufacturer) > 100 {
		return errors.New("制造商不能超过100个字符")
	}

	if req.IPAddress != nil {
		if err := v.validateIPAddress(*req.IPAddress); err != nil {
			return fmt.Errorf("IP地址验证失败: %w", err)
		}
	}

	if req.Port != nil {
		if err := v.validatePort(*req.Port); err != nil {
			return fmt.Errorf("端口验证失败: %w", err)
		}
	}

	if req.Protocol != nil && !v.isValidProtocol(*req.Protocol) {
		return fmt.Errorf("无效的协议类型: %s", *req.Protocol)
	}

	if req.Endpoint != nil && len(*req.Endpoint) > 500 {
		return errors.New("端点地址不能超过500个字符")
	}

	if req.AuthType != nil && !v.isValidAuthType(*req.AuthType) {
		return fmt.Errorf("无效的认证类型: %s", *req.AuthType)
	}

	return nil
}

// ValidateCommand 验证设备命令请求
func (v *DeviceValidator) ValidateCommand(req *DeviceCommandRequest) error {
	if req == nil {
		return errors.New("命令请求不能为空")
	}

	if strings.TrimSpace(req.CommandType) == "" {
		return errors.New("命令类型不能为空")
	}

	if !v.isValidCommandType(req.CommandType) {
		return fmt.Errorf("无效的命令类型: %s", req.CommandType)
	}

	if req.CommandData == nil || len(req.CommandData) == 0 {
		return errors.New("命令数据不能为空")
	}

	return nil
}

// ValidateHeartbeat 验证设备心跳请求
func (v *DeviceValidator) ValidateHeartbeat(req *DeviceHeartbeatRequest) error {
	if req == nil {
		return errors.New("心跳请求不能为空")
	}

	if strings.TrimSpace(req.Status) == "" {
		return errors.New("心跳状态不能为空")
	}

	if !v.isValidDeviceStatus(req.Status) {
		return fmt.Errorf("无效的设备状态: %s", req.Status)
	}

	if req.ResponseTime != nil && *req.ResponseTime < 0 {
		return errors.New("响应时间不能为负数")
	}

	return nil
}

// ValidateLog 验证设备日志请求
func (v *DeviceValidator) ValidateLog(req *DeviceLogRequest) error {
	if req == nil {
		return errors.New("日志请求不能为空")
	}

	if strings.TrimSpace(req.Level) == "" {
		return errors.New("日志级别不能为空")
	}

	if !v.isValidLogLevel(req.Level) {
		return fmt.Errorf("无效的日志级别: %s", req.Level)
	}

	if strings.TrimSpace(req.Message) == "" {
		return errors.New("日志消息不能为空")
	}

	if len(req.Message) > 10000 {
		return errors.New("日志消息不能超过10000个字符")
	}

	if req.CorrelationID != nil && len(*req.CorrelationID) > 100 {
		return errors.New("关联ID不能超过100个字符")
	}

	return nil
}

// ValidateStatusTransition 验证设备状态转换
func (v *DeviceValidator) ValidateStatusTransition(from, to string) error {
	if !v.isValidDeviceStatus(from) {
		return fmt.Errorf("无效的源状态: %s", from)
	}
	if !v.isValidDeviceStatus(to) {
		return fmt.Errorf("无效的目标状态: %s", to)
	}

	// 定义允许的状态转换
	allowedTransitions := map[string][]string{
		DeviceStatusOffline: {
			DeviceStatusOnline,
			DeviceStatusMaintenance,
		},
		DeviceStatusOnline: {
			DeviceStatusOffline,
			DeviceStatusMaintenance,
			DeviceStatusError,
		},
		DeviceStatusMaintenance: {
			DeviceStatusOnline,
			DeviceStatusOffline,
		},
		DeviceStatusError: {
			DeviceStatusOnline,
			DeviceStatusOffline,
			DeviceStatusMaintenance,
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

// ValidateFilter 验证设备过滤条件
func (v *DeviceValidator) ValidateFilter(filter *DeviceFilter) error {
	if filter == nil {
		return nil
	}

	// 验证状态过滤
	for _, status := range filter.Status {
		if !v.isValidDeviceStatus(status) {
			return fmt.Errorf("无效的状态过滤条件: %s", status)
		}
	}

	// 验证类型过滤
	for _, deviceType := range filter.Type {
		if !v.isValidDeviceType(deviceType) {
			return fmt.Errorf("无效的类型过滤条件: %s", deviceType)
		}
	}

	// 验证健康度范围
	if filter.MinHealth != nil && (*filter.MinHealth < 0 || *filter.MinHealth > 100) {
		return errors.New("最小健康度必须在0-100之间")
	}
	if filter.MaxHealth != nil && (*filter.MaxHealth < 0 || *filter.MaxHealth > 100) {
		return errors.New("最大健康度必须在0-100之间")
	}
	if filter.MinHealth != nil && filter.MaxHealth != nil && *filter.MinHealth > *filter.MaxHealth {
		return errors.New("最小健康度不能大于最大健康度")
	}

	// 验证时间范围
	if filter.LastSeenFrom != nil && filter.LastSeenTo != nil {
		if filter.LastSeenFrom.After(*filter.LastSeenTo) {
			return errors.New("起始时间不能晚于结束时间")
		}
	}

	return nil
}

// 辅助验证方法

// isValidDeviceIDFormat 检查设备ID格式
func (v *DeviceValidator) isValidDeviceIDFormat(deviceID string) bool {
	// 只允许字母、数字、连字符和下划线
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, deviceID)
	return matched
}

// isValidDeviceType 检查设备类型是否有效
func (v *DeviceValidator) isValidDeviceType(deviceType string) bool {
	validTypes := []string{
		DeviceTypeSensor,
		DeviceTypeCamera,
		DeviceTypeActuator,
		DeviceTypeGateway,
	}
	for _, valid := range validTypes {
		if valid == deviceType {
			return true
		}
	}
	return false
}

// isValidDeviceStatus 检查设备状态是否有效
func (v *DeviceValidator) isValidDeviceStatus(status string) bool {
	validStatuses := []string{
		DeviceStatusOnline,
		DeviceStatusOffline,
		DeviceStatusMaintenance,
		DeviceStatusError,
	}
	for _, valid := range validStatuses {
		if valid == status {
			return true
		}
	}
	return false
}

// isValidProtocol 检查协议类型是否有效
func (v *DeviceValidator) isValidProtocol(protocol string) bool {
	validProtocols := []string{
		"http", "https", "tcp", "udp", "mqtt", "coap", "modbus", "bacnet",
	}
	for _, valid := range validProtocols {
		if valid == strings.ToLower(protocol) {
			return true
		}
	}
	return false
}

// isValidAuthType 检查认证类型是否有效
func (v *DeviceValidator) isValidAuthType(authType string) bool {
	validTypes := []string{
		"none", "basic", "bearer", "oauth2", "apikey", "certificate",
	}
	for _, valid := range validTypes {
		if valid == strings.ToLower(authType) {
			return true
		}
	}
	return false
}

// isValidCommandType 检查命令类型是否有效
func (v *DeviceValidator) isValidCommandType(commandType string) bool {
	validTypes := []string{
		"start", "stop", "restart", "status", "configure", "reset", "update", "sync",
	}
	for _, valid := range validTypes {
		if valid == strings.ToLower(commandType) {
			return true
		}
	}
	return false
}

// isValidLogLevel 检查日志级别是否有效
func (v *DeviceValidator) isValidLogLevel(level string) bool {
	validLevels := []string{
		"debug", "info", "warn", "error", "fatal",
	}
	for _, valid := range validLevels {
		if valid == strings.ToLower(level) {
			return true
		}
	}
	return false
}

// validateIPAddress 验证IP地址格式
func (v *DeviceValidator) validateIPAddress(ipStr string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return errors.New("无效的IP地址格式")
	}
	return nil
}

// validatePort 验证端口号
func (v *DeviceValidator) validatePort(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("端口号必须在1-65535之间")
	}
	return nil
}

// isEmptyUpdateRequest 检查更新请求是否为空
func (v *DeviceValidator) isEmptyUpdateRequest(req *DeviceUpdateRequest) bool {
	return req.Name == nil &&
		req.Model == nil &&
		req.Manufacturer == nil &&
		req.IPAddress == nil &&
		req.Port == nil &&
		req.Protocol == nil &&
		req.Endpoint == nil &&
		req.AuthType == nil &&
		req.Credentials == nil &&
		req.Config == nil &&
		req.Capabilities == nil &&
		req.Metadata == nil
}
