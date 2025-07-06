package device

import (
	"OneGfServer/internal/model/device"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备验证相关业务逻辑
// ===============================

// ValidateDeviceRegistration 验证设备注册
func (s *sDevice) ValidateDeviceRegistration(ctx context.Context, input *device.ValidateDeviceRegistrationInput) (*device.ValidateDeviceRegistrationOutput, error) {
	output := &device.ValidateDeviceRegistrationOutput{
		IsValid: true,
		Errors:  []string{},
	}

	// 验证设备名称
	if input.Name == "" {
		output.IsValid = false
		output.Errors = append(output.Errors, "设备名称不能为空")
	} else if !isValidDeviceName(input.Name) {
		output.IsValid = false
		output.Errors = append(output.Errors, "设备名称格式无效")
	}

	// 验证MAC地址
	if input.MacAddress == "" {
		output.IsValid = false
		output.Errors = append(output.Errors, "MAC地址不能为空")
	} else if !isValidMacAddress(input.MacAddress) {
		output.IsValid = false
		output.Errors = append(output.Errors, "MAC地址格式无效")
	}

	// 验证IP地址
	if input.IPAddress == "" {
		output.IsValid = false
		output.Errors = append(output.Errors, "IP地址不能为空")
	} else if !isValidIPAddress(input.IPAddress) {
		output.IsValid = false
		output.Errors = append(output.Errors, "IP地址格式无效")
	}

	// 验证设备类型
	if input.DeviceType == "" {
		output.IsValid = false
		output.Errors = append(output.Errors, "设备类型不能为空")
	} else if !s.isValidDeviceType(input.DeviceType) {
		output.IsValid = false
		output.Errors = append(output.Errors, "不支持的设备类型")
	}

	if output.IsValid {
		output.Message = "设备注册验证通过"
	} else {
		output.Message = "设备注册验证失败"
	}

	return output, nil
}

// ValidateDeviceConfiguration 验证设备配置
func (s *sDevice) ValidateDeviceConfiguration(ctx context.Context, input *device.ValidateDeviceConfigurationInput) (*device.ValidateDeviceConfigurationOutput, error) {
	output := &device.ValidateDeviceConfigurationOutput{
		IsValid: true,
		Errors:  []string{},
	}

	// 验证基础配置
	if input.BasicConfig != nil {
		if err := s.validateBasicConfig(input.BasicConfig); err != nil {
			output.IsValid = false
			output.Errors = append(output.Errors, err.Error())
		}
	}

	// 验证性能配置
	if input.PerformanceConfig != nil {
		if err := s.validatePerformanceConfig(input.PerformanceConfig); err != nil {
			output.IsValid = false
			output.Errors = append(output.Errors, err.Error())
		}
	}

	// 验证安全配置
	if input.SecurityConfig != nil {
		if err := s.validateSecurityConfig(input.SecurityConfig); err != nil {
			output.IsValid = false
			output.Errors = append(output.Errors, err.Error())
		}
	}

	if output.IsValid {
		output.Message = "设备配置验证通过"
	} else {
		output.Message = "设备配置验证失败"
	}

	return output, nil
}

// validateBasicConfig 验证基础配置
func (s *sDevice) validateBasicConfig(config *device.BasicConfig) error {
	// 验证设备名称
	if config.Name != "" && !isValidDeviceName(config.Name) {
		return gerror.NewCode(gcode.CodeValidationFailed, "设备名称格式无效")
	}

	// 验证描述
	if config.Description != "" && len(config.Description) > 500 {
		return gerror.NewCode(gcode.CodeValidationFailed, "设备描述不能超过500个字符")
	}

	// 验证标签
	if config.Tags != nil {
		for _, tag := range config.Tags {
			if len(tag) > 50 {
				return gerror.NewCode(gcode.CodeValidationFailed, "标签长度不能超过50个字符")
			}
		}
	}

	return nil
}

// validatePerformanceConfig 验证性能配置
func (s *sDevice) validatePerformanceConfig(config *device.PerformanceConfig) error {
	// 验证最大并发任务数
	if config.MaxConcurrentTasks <= 0 || config.MaxConcurrentTasks > 100 {
		return gerror.NewCode(gcode.CodeValidationFailed, "最大并发任务数必须在1-100之间")
	}

	// 验证CPU阈值
	if config.CPUThreshold <= 0 || config.CPUThreshold > 100 {
		return gerror.NewCode(gcode.CodeValidationFailed, "CPU阈值必须在1-100之间")
	}

	// 验证内存阈值
	if config.MemoryThreshold <= 0 || config.MemoryThreshold > 100 {
		return gerror.NewCode(gcode.CodeValidationFailed, "内存阈值必须在1-100之间")
	}

	return nil
}

// validateSecurityConfig 验证安全配置
func (s *sDevice) validateSecurityConfig(config *device.SecurityConfig) error {
	// 验证访问控制
	if config.AccessControl != nil {
		if err := s.validateAccessControl(config.AccessControl); err != nil {
			return err
		}
	}

	// 验证加密设置
	if config.Encryption != nil {
		if err := s.validateEncryptionSettings(config.Encryption); err != nil {
			return err
		}
	}

	return nil
}

// validateAccessControl 验证访问控制
func (s *sDevice) validateAccessControl(accessControl *device.AccessControl) error {
	// 验证允许的IP地址
	if accessControl.AllowedIPs != nil {
		for _, ip := range accessControl.AllowedIPs {
			if !isValidIPAddress(ip) {
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("无效的IP地址: %s", ip))
			}
		}
	}

	// 验证用户权限
	if accessControl.UserPermissions != nil {
		validPermissions := []string{"read", "write", "admin", "monitor"}
		for _, permission := range accessControl.UserPermissions {
			if !contains(validPermissions, permission) {
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("无效的权限: %s", permission))
			}
		}
	}

	return nil
}

// validateEncryptionSettings 验证加密设置
func (s *sDevice) validateEncryptionSettings(encryption *device.Encryption) error {
	// 验证加密算法
	if encryption.Algorithm != "" {
		validAlgorithms := []string{"AES-256", "AES-128", "3DES"}
		if !contains(validAlgorithms, encryption.Algorithm) {
			return gerror.NewCode(gcode.CodeValidationFailed, "不支持的加密算法")
		}
	}

	// 验证密钥长度
	if encryption.KeyLength < 128 || encryption.KeyLength > 256 {
		return gerror.NewCode(gcode.CodeValidationFailed, "密钥长度必须在128-256之间")
	}

	return nil
}

// ValidateDeviceData 验证设备数据
func (s *sDevice) ValidateDeviceData(ctx context.Context, input *device.ValidateDeviceDataInput) (*device.ValidateDeviceDataOutput, error) {
	output := &device.ValidateDeviceDataOutput{
		IsValid: true,
		Errors:  []string{},
	}

	// 验证必需字段
	requiredFields := []string{"device_id", "name", "status"}
	for _, field := range requiredFields {
		if value, ok := input.DeviceData[field]; !ok || value == "" {
			output.IsValid = false
			output.Errors = append(output.Errors, fmt.Sprintf("必需字段 %s 不能为空", field))
		}
	}

	// 验证状态值
	if status, ok := input.DeviceData["status"].(string); ok {
		validStatuses := []string{"online", "offline", "busy", "maintenance", "error", "deleted"}
		if !contains(validStatuses, status) {
			output.IsValid = false
			output.Errors = append(output.Errors, "无效的设备状态")
		}
	}

	// 验证数值字段
	if cpuUsage, ok := input.DeviceData["cpu_usage"].(float64); ok {
		if cpuUsage < 0 || cpuUsage > 100 {
			output.IsValid = false
			output.Errors = append(output.Errors, "CPU使用率必须在0-100之间")
		}
	}

	if memUsage, ok := input.DeviceData["memory_usage"].(float64); ok {
		if memUsage < 0 || memUsage > 100 {
			output.IsValid = false
			output.Errors = append(output.Errors, "内存使用率必须在0-100之间")
		}
	}

	if output.IsValid {
		output.Message = "设备数据验证通过"
	} else {
		output.Message = "设备数据验证失败"
	}

	return output, nil
}

// ProcessDeviceData 处理设备数据
func (s *sDevice) ProcessDeviceData(ctx context.Context, input *device.ProcessDeviceDataInput) (*device.ProcessDeviceDataOutput, error) {
	output := &device.ProcessDeviceDataOutput{
		ProcessedData: make(map[string]interface{}),
		IsValid:       true,
	}

	// 验证原始数据
	validateInput := &device.ValidateDeviceDataInput{
		DeviceData: input.RawData,
	}
	validateOutput, err := s.ValidateDeviceData(ctx, validateInput)
	if err != nil {
		return nil, err
	}

	if !validateOutput.IsValid {
		output.IsValid = false
		output.Message = "设备数据验证失败"
		return output, nil
	}

	// 数据清洗和转换
	for key, value := range input.RawData {
		output.ProcessedData[key] = value
	}

	// 添加处理时间戳
	output.ProcessedData["processed_at"] = gtime.Now().Format("2006-01-02 15:04:05")

	// 计算派生字段
	if cpuUsage, ok := input.RawData["cpu_usage"].(float64); ok {
		output.ProcessedData["cpu_status"] = s.getCPUStatus(cpuUsage)
	}

	if memUsage, ok := input.RawData["memory_usage"].(float64); ok {
		output.ProcessedData["memory_status"] = s.getMemoryStatus(memUsage)
	}

	output.Message = "设备数据处理完成"
	return output, nil
}

// getCPUStatus 获取CPU状态
func (s *sDevice) getCPUStatus(cpuUsage float64) string {
	if cpuUsage > 90 {
		return "critical"
	} else if cpuUsage > 70 {
		return "warning"
	} else {
		return "normal"
	}
}

// getMemoryStatus 获取内存状态
func (s *sDevice) getMemoryStatus(memUsage float64) string {
	if memUsage > 90 {
		return "critical"
	} else if memUsage > 80 {
		return "warning"
	} else {
		return "normal"
	}
}

// isValidDeviceType 验证设备类型
func (s *sDevice) isValidDeviceType(deviceType string) bool {
	validTypes := []string{"server", "workstation", "mobile", "iot", "edge"}
	return contains(validTypes, deviceType)
}

// 工具函数
func isValidDeviceName(name string) bool {
	if name == "" || len(name) > 100 {
		return false
	}
	// 只允许字母、数字、下划线、连字符和空格
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-\s]+$`, name)
	return matched
}

func isValidMacAddress(mac string) bool {
	// 验证MAC地址格式 (xx:xx:xx:xx:xx:xx 或 xx-xx-xx-xx-xx-xx)
	matched, _ := regexp.MatchString(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, mac)
	return matched
}

func isValidIPAddress(ip string) bool {
	// 简单的IP地址验证
	matched, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
	if !matched {
		return false
	}

	// 检查每个段是否在0-255范围内
	parts := strings.Split(ip, ".")
	for _, part := range parts {
		if len(part) > 1 && part[0] == '0' {
			return false // 不允许前导零
		}
		if num := 0; num < 0 || num > 255 {
			return false
		}
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
