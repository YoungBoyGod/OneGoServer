package device

import (
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
func (s *sDevice) ValidateDeviceRegistration(ctx context.Context, deviceData map[string]interface{}) error {
	// 验证设备名称
	if name, ok := deviceData["name"].(string); ok {
		if !isValidDeviceName(name) {
			return gerror.NewCode(gcode.CodeValidationFailed, "设备名称格式无效")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "设备名称不能为空")
	}

	// 验证MAC地址
	if macAddress, ok := deviceData["mac_address"].(string); ok {
		if !isValidMacAddress(macAddress) {
			return gerror.NewCode(gcode.CodeValidationFailed, "MAC地址格式无效")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "MAC地址不能为空")
	}

	// 验证IP地址
	if ipAddress, ok := deviceData["ip_address"].(string); ok {
		if !isValidIPAddress(ipAddress) {
			return gerror.NewCode(gcode.CodeValidationFailed, "IP地址格式无效")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "IP地址不能为空")
	}

	// 验证设备类型
	if deviceType, ok := deviceData["device_type"].(string); ok {
		if !s.isValidDeviceType(deviceType) {
			return gerror.NewCode(gcode.CodeValidationFailed, "不支持的设备类型")
		}
	} else {
		return gerror.NewCode(gcode.CodeValidationFailed, "设备类型不能为空")
	}

	return nil
}

// ValidateDeviceConfiguration 验证设备配置
func (s *sDevice) ValidateDeviceConfiguration(ctx context.Context, configData map[string]interface{}) error {
	// 验证基础配置
	if err := s.validateBasicConfig(configData); err != nil {
		return err
	}

	// 验证性能配置
	if err := s.validatePerformanceConfig(configData); err != nil {
		return err
	}

	// 验证安全配置
	if err := s.validateSecurityConfig(configData); err != nil {
		return err
	}

	return nil
}

// validateBasicConfig 验证基础配置
func (s *sDevice) validateBasicConfig(config map[string]interface{}) error {
	// 验证设备名称
	if name, ok := config["name"].(string); ok {
		if !isValidDeviceName(name) {
			return gerror.NewCode(gcode.CodeValidationFailed, "设备名称格式无效")
		}
	}

	// 验证描述
	if description, ok := config["description"].(string); ok {
		if len(description) > 500 {
			return gerror.NewCode(gcode.CodeValidationFailed, "设备描述不能超过500个字符")
		}
	}

	// 验证标签
	if tags, ok := config["tags"].([]string); ok {
		for _, tag := range tags {
			if len(tag) > 50 {
				return gerror.NewCode(gcode.CodeValidationFailed, "标签长度不能超过50个字符")
			}
		}
	}

	return nil
}

// validatePerformanceConfig 验证性能配置
func (s *sDevice) validatePerformanceConfig(config map[string]interface{}) error {
	// 验证最大并发任务数
	if maxTasks, ok := config["max_concurrent_tasks"].(int); ok {
		if maxTasks <= 0 || maxTasks > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "最大并发任务数必须在1-100之间")
		}
	}

	// 验证CPU阈值
	if cpuThreshold, ok := config["cpu_threshold"].(float64); ok {
		if cpuThreshold <= 0 || cpuThreshold > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "CPU阈值必须在1-100之间")
		}
	}

	// 验证内存阈值
	if memThreshold, ok := config["memory_threshold"].(float64); ok {
		if memThreshold <= 0 || memThreshold > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "内存阈值必须在1-100之间")
		}
	}

	return nil
}

// validateSecurityConfig 验证安全配置
func (s *sDevice) validateSecurityConfig(config map[string]interface{}) error {
	// 验证访问控制
	if accessControl, ok := config["access_control"].(map[string]interface{}); ok {
		if err := s.validateAccessControl(accessControl); err != nil {
			return err
		}
	}

	// 验证加密设置
	if encryption, ok := config["encryption"].(map[string]interface{}); ok {
		if err := s.validateEncryptionSettings(encryption); err != nil {
			return err
		}
	}

	return nil
}

// validateAccessControl 验证访问控制
func (s *sDevice) validateAccessControl(accessControl map[string]interface{}) error {
	// 验证允许的IP地址
	if allowedIPs, ok := accessControl["allowed_ips"].([]string); ok {
		for _, ip := range allowedIPs {
			if !isValidIPAddress(ip) {
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("无效的IP地址: %s", ip))
			}
		}
	}

	// 验证用户权限
	if userPermissions, ok := accessControl["user_permissions"].([]string); ok {
		validPermissions := []string{"read", "write", "admin", "monitor"}
		for _, permission := range userPermissions {
			if !contains(validPermissions, permission) {
				return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("无效的权限: %s", permission))
			}
		}
	}

	return nil
}

// validateEncryptionSettings 验证加密设置
func (s *sDevice) validateEncryptionSettings(encryption map[string]interface{}) error {
	// 验证加密算法
	if algorithm, ok := encryption["algorithm"].(string); ok {
		validAlgorithms := []string{"AES-256", "AES-128", "3DES"}
		if !contains(validAlgorithms, algorithm) {
			return gerror.NewCode(gcode.CodeValidationFailed, "不支持的加密算法")
		}
	}

	// 验证密钥长度
	if keyLength, ok := encryption["key_length"].(int); ok {
		if keyLength < 128 || keyLength > 256 {
			return gerror.NewCode(gcode.CodeValidationFailed, "密钥长度必须在128-256之间")
		}
	}

	return nil
}

// ValidateDeviceData 验证设备数据
func (s *sDevice) ValidateDeviceData(ctx context.Context, deviceData map[string]interface{}) error {
	// 验证必需字段
	requiredFields := []string{"device_id", "name", "status"}
	for _, field := range requiredFields {
		if value, ok := deviceData[field]; !ok || value == "" {
			return gerror.NewCode(gcode.CodeValidationFailed, fmt.Sprintf("必需字段 %s 不能为空", field))
		}
	}

	// 验证状态值
	if status, ok := deviceData["status"].(string); ok {
		validStatuses := []string{"online", "offline", "busy", "maintenance", "error", "deleted"}
		if !contains(validStatuses, status) {
			return gerror.NewCode(gcode.CodeValidationFailed, "无效的设备状态")
		}
	}

	// 验证数值字段
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage < 0 || cpuUsage > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "CPU使用率必须在0-100之间")
		}
	}

	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage < 0 || memUsage > 100 {
			return gerror.NewCode(gcode.CodeValidationFailed, "内存使用率必须在0-100之间")
		}
	}

	return nil
}

// ProcessDeviceData 处理设备数据
func (s *sDevice) ProcessDeviceData(ctx context.Context, rawData map[string]interface{}) (map[string]interface{}, error) {
	// 验证原始数据
	if err := s.ValidateDeviceData(ctx, rawData); err != nil {
		return nil, err
	}

	// 数据清洗和转换
	processedData := make(map[string]interface{})
	for key, value := range rawData {
		processedData[key] = value
	}

	// 添加处理时间戳
	processedData["processed_at"] = gtime.Now().Format("2006-01-02 15:04:05")

	// 计算派生字段
	if cpuUsage, ok := rawData["cpu_usage"].(float64); ok {
		processedData["cpu_status"] = s.getCPUStatus(cpuUsage)
	}

	if memUsage, ok := rawData["memory_usage"].(float64); ok {
		processedData["memory_status"] = s.getMemoryStatus(memUsage)
	}

	return processedData, nil
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
