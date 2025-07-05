package device

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sDevice struct{}

func New() *sDevice {
	return &sDevice{}
}

func init() {
	// service.RegisterDevice(New())
}

// ===============================
// 设备状态管理相关业务逻辑
// ===============================

// ValidateDeviceStatus 验证设备状态转换是否合法
func (s *sDevice) ValidateDeviceStatus(ctx context.Context, currentStatus, targetStatus string) error {
	// 定义合法的状态转换规则
	validTransitions := map[string][]string{
		"offline":     {"online", "maintenance", "deleted"},
		"online":      {"offline", "busy", "maintenance", "error"},
		"busy":        {"online", "offline", "error"},
		"maintenance": {"online", "offline"},
		"error":       {"offline", "maintenance"},
		"deleted":     {}, // 删除状态不能转换到其他状态
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("无效的当前状态: %s", currentStatus))
	}

	// 检查目标状态是否在允许的转换列表中
	for _, status := range allowedStatuses {
		if status == targetStatus {
			return nil
		}
	}

	return gerror.NewCode(gcode.CodeInvalidParameter,
		fmt.Sprintf("不允许从状态 %s 转换到 %s", currentStatus, targetStatus))
}

// CalculateDeviceHealthScore 计算设备健康度评分
func (s *sDevice) CalculateDeviceHealthScore(ctx context.Context, deviceData map[string]interface{}) float64 {
	var totalScore float64 = 100.0

	// 1. CPU使用率影响 (权重: 25%)
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			totalScore -= 25
		} else if cpuUsage > 70 {
			totalScore -= 15
		} else if cpuUsage > 50 {
			totalScore -= 5
		}
	}

	// 2. 内存使用率影响 (权重: 25%)
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage > 90 {
			totalScore -= 25
		} else if memUsage > 80 {
			totalScore -= 15
		} else if memUsage > 60 {
			totalScore -= 5
		}
	}

	// 3. 磁盘使用率影响 (权重: 20%)
	if diskUsage, ok := deviceData["disk_usage"].(float64); ok {
		if diskUsage > 95 {
			totalScore -= 20
		} else if diskUsage > 85 {
			totalScore -= 10
		} else if diskUsage > 70 {
			totalScore -= 3
		}
	}

	// 4. 网络连接状态影响 (权重: 15%)
	if networkStatus, ok := deviceData["network_status"].(string); ok {
		switch networkStatus {
		case "disconnected":
			totalScore -= 15
		case "unstable":
			totalScore -= 8
		case "slow":
			totalScore -= 3
		}
	}

	// 5. 最后心跳时间影响 (权重: 10%)
	if lastHeartbeat, ok := deviceData["last_heartbeat"].(*gtime.Time); ok && lastHeartbeat != nil {
		timeDiff := time.Since(lastHeartbeat.Time)
		if timeDiff > 10*time.Minute {
			totalScore -= 10
		} else if timeDiff > 5*time.Minute {
			totalScore -= 5
		} else if timeDiff > 2*time.Minute {
			totalScore -= 2
		}
	}

	// 6. 错误次数影响 (权重: 5%)
	if errorCount, ok := deviceData["error_count"].(int); ok {
		if errorCount > 10 {
			totalScore -= 5
		} else if errorCount > 5 {
			totalScore -= 3
		} else if errorCount > 2 {
			totalScore -= 1
		}
	}

	// 确保评分在0-100范围内
	if totalScore < 0 {
		totalScore = 0
	}
	if totalScore > 100 {
		totalScore = 100
	}

	return math.Round(totalScore*100) / 100
}

// DetermineDeviceStatus 根据设备数据自动确定设备状态
func (s *sDevice) DetermineDeviceStatus(ctx context.Context, deviceData map[string]interface{}) string {
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)

	// 检查是否离线
	if lastHeartbeat, ok := deviceData["last_heartbeat"].(*gtime.Time); ok && lastHeartbeat != nil {
		if time.Since(lastHeartbeat.Time) > 15*time.Minute {
			return "offline"
		}
	}

	// 检查是否处于维护状态
	if maintenanceMode, ok := deviceData["maintenance_mode"].(bool); ok && maintenanceMode {
		return "maintenance"
	}

	// 根据健康度评分确定状态
	if healthScore < 30 {
		return "error"
	} else if healthScore < 70 {
		// 检查是否正在执行任务
		if taskCount, ok := deviceData["running_task_count"].(int); ok && taskCount > 0 {
			return "busy"
		}
		return "online"
	} else {
		// 检查是否正在执行任务
		if taskCount, ok := deviceData["running_task_count"].(int); ok && taskCount > 0 {
			return "busy"
		}
		return "online"
	}
}

// CanAcceptNewTask 判断设备是否可以接受新任务
func (s *sDevice) CanAcceptNewTask(ctx context.Context, deviceData map[string]interface{}) bool {
	// 检查设备状态
	status, ok := deviceData["status"].(string)
	if !ok || status != "online" {
		return false
	}

	// 检查健康度
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)
	if healthScore < 50 {
		return false
	}

	// 检查任务负载
	if currentTasks, ok := deviceData["running_task_count"].(int); ok {
		if maxTasks, ok := deviceData["max_concurrent_tasks"].(int); ok {
			return currentTasks < maxTasks
		}
		// 默认最大并发任务数为5
		return currentTasks < 5
	}

	return true
}

// ===============================
// 设备验证相关业务逻辑
// ===============================

// ValidateDeviceRegistration 验证设备注册数据
func (s *sDevice) ValidateDeviceRegistration(ctx context.Context, deviceData map[string]interface{}) error {
	// 验证必需字段
	requiredFields := []string{"device_name", "device_type", "device_mac"}
	for _, field := range requiredFields {
		if value, ok := deviceData[field]; !ok || value == nil || value == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("缺少必需字段: %s", field))
		}
	}

	// 验证设备名称格式
	if deviceName, ok := deviceData["device_name"].(string); ok {
		if len(deviceName) < 3 || len(deviceName) > 50 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "设备名称长度必须在3-50字符之间")
		}
		if !isValidDeviceName(deviceName) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "设备名称包含无效字符")
		}
	}

	// 验证MAC地址格式
	if deviceMac, ok := deviceData["device_mac"].(string); ok {
		if !isValidMacAddress(deviceMac) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "MAC地址格式不正确")
		}
	}

	// 验证设备类型
	if deviceType, ok := deviceData["device_type"].(string); ok {
		validTypes := []string{"server", "workstation", "mobile", "iot", "embedded"}
		if !contains(validTypes, deviceType) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "无效的设备类型")
		}
	}

	// 验证IP地址格式（如果提供）
	if ipAddress, ok := deviceData["ip_address"].(string); ok && ipAddress != "" {
		if !isValidIPAddress(ipAddress) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "IP地址格式不正确")
		}
	}

	return nil
}

// ValidateDeviceConfiguration 验证设备配置数据
func (s *sDevice) ValidateDeviceConfiguration(ctx context.Context, configData map[string]interface{}) error {
	// 验证最大并发任务数
	if maxTasks, ok := configData["max_concurrent_tasks"]; ok {
		if taskCount := gconv.Int(maxTasks); taskCount < 1 || taskCount > 100 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "最大并发任务数必须在1-100之间")
		}
	}

	// 验证心跳间隔
	if heartbeatInterval, ok := configData["heartbeat_interval"]; ok {
		if interval := gconv.Int(heartbeatInterval); interval < 30 || interval > 3600 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "心跳间隔必须在30-3600秒之间")
		}
	}

	// 验证资源限制
	if resourceLimits, ok := configData["resource_limits"].(map[string]interface{}); ok {
		if cpuLimit, exists := resourceLimits["cpu"]; exists {
			if limit := gconv.Float64(cpuLimit); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "CPU限制必须在0-100之间")
			}
		}
		if memLimit, exists := resourceLimits["memory"]; exists {
			if limit := gconv.Float64(memLimit); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "内存限制必须在0-100之间")
			}
		}
	}

	return nil
}

// ===============================
// 设备生命周期管理
// ===============================

// HandleDeviceRegistration 处理设备注册业务逻辑
func (s *sDevice) HandleDeviceRegistration(ctx context.Context, deviceData map[string]interface{}) (map[string]interface{}, error) {
	// 验证注册数据
	if err := s.ValidateDeviceRegistration(ctx, deviceData); err != nil {
		return nil, err
	}

	// 设置默认值
	result := make(map[string]interface{})
	for k, v := range deviceData {
		result[k] = v
	}

	// 设置默认状态
	result["status"] = "online"

	// 计算初始健康度
	result["health_score"] = s.CalculateDeviceHealthScore(ctx, result)

	// 设置注册时间
	result["registered_at"] = gtime.Now()
	result["last_heartbeat"] = gtime.Now()

	// 设置默认配置
	if _, exists := result["max_concurrent_tasks"]; !exists {
		result["max_concurrent_tasks"] = 5
	}
	if _, exists := result["heartbeat_interval"]; !exists {
		result["heartbeat_interval"] = 60 // 60秒
	}

	g.Log().Info(ctx, "设备注册成功", g.Map{
		"device_name":  result["device_name"],
		"device_mac":   result["device_mac"],
		"status":       result["status"],
		"health_score": result["health_score"],
	})

	return result, nil
}

// HandleDeviceHeartbeat 处理设备心跳业务逻辑
func (s *sDevice) HandleDeviceHeartbeat(ctx context.Context, deviceId string, heartbeatData map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 更新心跳时间
	result["last_heartbeat"] = gtime.Now()

	// 更新设备状态信息
	if heartbeatData != nil {
		for k, v := range heartbeatData {
			result[k] = v
		}
	}

	// 重新计算健康度
	result["health_score"] = s.CalculateDeviceHealthScore(ctx, result)

	// 自动确定设备状态
	newStatus := s.DetermineDeviceStatus(ctx, result)
	result["status"] = newStatus

	// 更新最后活跃时间
	result["last_active_at"] = gtime.Now()

	g.Log().Debug(ctx, "处理设备心跳", g.Map{
		"device_id":    deviceId,
		"new_status":   newStatus,
		"health_score": result["health_score"],
	})

	return result, nil
}

// HandleDeviceDeactivation 处理设备停用业务逻辑
func (s *sDevice) HandleDeviceDeactivation(ctx context.Context, deviceId string, reason string) error {
	// 记录停用操作
	g.Log().Info(ctx, "设备停用", g.Map{
		"device_id": deviceId,
		"reason":    reason,
		"timestamp": gtime.Now(),
	})

	// 这里可以添加额外的业务逻辑，比如：
	// 1. 停止所有正在运行的任务
	// 2. 清理缓存数据
	// 3. 发送通知

	return nil
}

// ===============================
// 设备任务管理相关业务逻辑
// ===============================

// CalculateTaskAssignmentScore 计算任务分配评分
func (s *sDevice) CalculateTaskAssignmentScore(ctx context.Context, deviceData map[string]interface{}, taskData map[string]interface{}) float64 {
	score := 0.0

	// 基础健康度评分 (权重: 40%)
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)
	score += healthScore * 0.4

	// 任务负载评分 (权重: 30%)
	if currentTasks, ok := deviceData["running_task_count"].(int); ok {
		if maxTasks, ok := deviceData["max_concurrent_tasks"].(int); ok && maxTasks > 0 {
			loadRatio := float64(currentTasks) / float64(maxTasks)
			loadScore := (1.0 - loadRatio) * 100
			score += loadScore * 0.3
		}
	}

	// 设备类型匹配度 (权重: 20%)
	if deviceType, ok := deviceData["device_type"].(string); ok {
		if taskType, ok := taskData["preferred_device_type"].(string); ok {
			if deviceType == taskType {
				score += 100 * 0.2
			} else {
				score += 50 * 0.2 // 部分匹配
			}
		} else {
			score += 75 * 0.2 // 无偏好，给予中等评分
		}
	}

	// 地理位置因素 (权重: 10%)
	if deviceLocation, ok := deviceData["location"].(string); ok {
		if taskLocation, ok := taskData["preferred_location"].(string); ok {
			if deviceLocation == taskLocation {
				score += 100 * 0.1
			} else {
				score += 30 * 0.1 // 不同位置，给予较低评分
			}
		} else {
			score += 50 * 0.1 // 无位置偏好
		}
	}

	return math.Round(score*100) / 100
}

// ===============================
// 设备监控和分析
// ===============================

// AnalyzeDevicePerformance 分析设备性能表现
func (s *sDevice) AnalyzeDevicePerformance(ctx context.Context, performanceData []map[string]interface{}) map[string]interface{} {
	if len(performanceData) == 0 {
		return map[string]interface{}{
			"status":  "insufficient_data",
			"message": "性能数据不足",
		}
	}

	analysis := make(map[string]interface{})

	// 计算平均值
	var totalCPU, totalMemory, totalDisk float64
	var healthScores []float64

	for _, data := range performanceData {
		if cpu, ok := data["cpu_usage"].(float64); ok {
			totalCPU += cpu
		}
		if mem, ok := data["memory_usage"].(float64); ok {
			totalMemory += mem
		}
		if disk, ok := data["disk_usage"].(float64); ok {
			totalDisk += disk
		}

		healthScore := s.CalculateDeviceHealthScore(ctx, data)
		healthScores = append(healthScores, healthScore)
	}

	count := float64(len(performanceData))
	analysis["avg_cpu_usage"] = math.Round((totalCPU/count)*100) / 100
	analysis["avg_memory_usage"] = math.Round((totalMemory/count)*100) / 100
	analysis["avg_disk_usage"] = math.Round((totalDisk/count)*100) / 100

	// 计算健康度趋势
	if len(healthScores) > 0 {
		avgHealth := 0.0
		for _, score := range healthScores {
			avgHealth += score
		}
		avgHealth /= float64(len(healthScores))
		analysis["avg_health_score"] = math.Round(avgHealth*100) / 100

		// 判断性能趋势
		if len(healthScores) >= 2 {
			recent := healthScores[len(healthScores)-1]
			previous := healthScores[len(healthScores)-2]
			if recent > previous+5 {
				analysis["trend"] = "improving"
			} else if recent < previous-5 {
				analysis["trend"] = "declining"
			} else {
				analysis["trend"] = "stable"
			}
		}
	}

	// 生成建议
	recommendations := s.generatePerformanceRecommendations(analysis)
	analysis["recommendations"] = recommendations

	return analysis
}

// ===============================
// 辅助函数
// ===============================

// generatePerformanceRecommendations 生成性能优化建议
func (s *sDevice) generatePerformanceRecommendations(analysis map[string]interface{}) []string {
	var recommendations []string

	if avgCPU, ok := analysis["avg_cpu_usage"].(float64); ok && avgCPU > 80 {
		recommendations = append(recommendations, "CPU使用率较高，建议优化进程或增加处理能力")
	}

	if avgMem, ok := analysis["avg_memory_usage"].(float64); ok && avgMem > 85 {
		recommendations = append(recommendations, "内存使用率较高，建议清理缓存或增加内存")
	}

	if avgDisk, ok := analysis["avg_disk_usage"].(float64); ok && avgDisk > 90 {
		recommendations = append(recommendations, "磁盘空间不足，建议清理文件或扩容磁盘")
	}

	if trend, ok := analysis["trend"].(string); ok && trend == "declining" {
		recommendations = append(recommendations, "设备性能呈下降趋势，建议进行维护检查")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "设备运行状态良好，保持当前配置")
	}

	return recommendations
}

// isValidDeviceName 验证设备名称格式
func isValidDeviceName(name string) bool {
	if len(name) == 0 {
		return false
	}
	// 允许字母、数字、连字符和下划线
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}

// isValidMacAddress 验证MAC地址格式
func isValidMacAddress(mac string) bool {
	// 简单的MAC地址格式验证 (XX:XX:XX:XX:XX:XX)
	if len(mac) != 17 {
		return false
	}
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		return false
	}
	for _, part := range parts {
		if len(part) != 2 {
			return false
		}
		for _, char := range part {
			if !((char >= '0' && char <= '9') ||
				(char >= 'A' && char <= 'F') ||
				(char >= 'a' && char <= 'f')) {
				return false
			}
		}
	}
	return true
}

// isValidIPAddress 验证IP地址格式
func isValidIPAddress(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if num := gconv.Int(part); num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ===============================
// 设备配置管理相关业务逻辑
// ===============================

// ValidateDeviceConfig 验证设备配置
func (s *sDevice) ValidateDeviceConfig(ctx context.Context, config map[string]interface{}) error {
	// 验证配置格式
	if config == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "设备配置不能为空")
	}

	// 验证心跳间隔
	if heartbeatInterval, ok := config["heartbeat_interval"]; ok {
		if interval := gconv.Int(heartbeatInterval); interval < 30 || interval > 3600 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "心跳间隔必须在30-3600秒之间")
		}
	}

	// 验证最大并发任务数
	if maxTasks, ok := config["max_concurrent_tasks"]; ok {
		if taskCount := gconv.Int(maxTasks); taskCount < 1 || taskCount > 100 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "最大并发任务数必须在1-100之间")
		}
	}

	// 验证资源限制
	if resourceLimits, ok := config["resource_limits"].(map[string]interface{}); ok {
		if cpuLimit, exists := resourceLimits["cpu"]; exists {
			if limit := gconv.Float64(cpuLimit); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "CPU限制必须在0-100之间")
			}
		}
		if memLimit, exists := resourceLimits["memory"]; exists {
			if limit := gconv.Float64(memLimit); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "内存限制必须在0-100之间")
			}
		}
	}

	return nil
}

// GenerateDeviceConfig 生成设备默认配置
func (s *sDevice) GenerateDeviceConfig(ctx context.Context, deviceType string) map[string]interface{} {
	config := make(map[string]interface{})

	// 根据设备类型设置默认配置
	switch deviceType {
	case "sensor":
		config["heartbeat_interval"] = 60
		config["max_concurrent_tasks"] = 3
		config["resource_limits"] = map[string]interface{}{
			"cpu":    80.0,
			"memory": 70.0,
			"disk":   85.0,
		}
	case "camera":
		config["heartbeat_interval"] = 30
		config["max_concurrent_tasks"] = 5
		config["resource_limits"] = map[string]interface{}{
			"cpu":    90.0,
			"memory": 80.0,
			"disk":   90.0,
		}
	case "actuator":
		config["heartbeat_interval"] = 120
		config["max_concurrent_tasks"] = 2
		config["resource_limits"] = map[string]interface{}{
			"cpu":    60.0,
			"memory": 50.0,
			"disk":   70.0,
		}
	case "gateway":
		config["heartbeat_interval"] = 45
		config["max_concurrent_tasks"] = 10
		config["resource_limits"] = map[string]interface{}{
			"cpu":    85.0,
			"memory": 75.0,
			"disk":   80.0,
		}
	default:
		config["heartbeat_interval"] = 60
		config["max_concurrent_tasks"] = 5
		config["resource_limits"] = map[string]interface{}{
			"cpu":    75.0,
			"memory": 65.0,
			"disk":   75.0,
		}
	}

	// 通用配置
	config["retry_count"] = 3
	config["timeout_seconds"] = 300
	config["enable_auto_recovery"] = true
	config["enable_performance_monitoring"] = true

	return config
}

// ===============================
// 设备监控和告警相关业务逻辑
// ===============================

// CheckDeviceHealth 检查设备健康状态
func (s *sDevice) CheckDeviceHealth(ctx context.Context, deviceData map[string]interface{}) map[string]interface{} {
	healthInfo := make(map[string]interface{})

	// 计算健康度评分
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)
	healthInfo["health_score"] = healthScore

	// 确定健康状态
	if healthScore >= 90 {
		healthInfo["status"] = "excellent"
	} else if healthScore >= 70 {
		healthInfo["status"] = "good"
	} else if healthScore >= 50 {
		healthInfo["status"] = "fair"
	} else {
		healthInfo["status"] = "poor"
	}

	// 检查各项指标
	healthInfo["checks"] = s.performHealthChecks(ctx, deviceData)

	// 生成告警信息
	alerts := s.generateHealthAlerts(ctx, deviceData, healthScore)
	healthInfo["alerts"] = alerts

	return healthInfo
}

// performHealthChecks 执行健康检查
func (s *sDevice) performHealthChecks(ctx context.Context, deviceData map[string]interface{}) map[string]interface{} {
	checks := make(map[string]interface{})

	// CPU检查
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		checks["cpu"] = map[string]interface{}{
			"value":   cpuUsage,
			"status":  s.getCheckStatus(cpuUsage, 80, 90),
			"message": s.getCPUMessage(cpuUsage),
		}
	}

	// 内存检查
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		checks["memory"] = map[string]interface{}{
			"value":   memUsage,
			"status":  s.getCheckStatus(memUsage, 85, 95),
			"message": s.getMemoryMessage(memUsage),
		}
	}

	// 磁盘检查
	if diskUsage, ok := deviceData["disk_usage"].(float64); ok {
		checks["disk"] = map[string]interface{}{
			"value":   diskUsage,
			"status":  s.getCheckStatus(diskUsage, 90, 95),
			"message": s.getDiskMessage(diskUsage),
		}
	}

	// 网络检查
	if networkStatus, ok := deviceData["network_status"].(string); ok {
		checks["network"] = map[string]interface{}{
			"value":   networkStatus,
			"status":  s.getNetworkStatus(networkStatus),
			"message": s.getNetworkMessage(networkStatus),
		}
	}

	return checks
}

// generateHealthAlerts 生成健康告警
func (s *sDevice) generateHealthAlerts(ctx context.Context, deviceData map[string]interface{}, healthScore float64) []map[string]interface{} {
	var alerts []map[string]interface{}

	// 健康度告警
	if healthScore < 30 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "critical",
			"type":    "health_score",
			"message": "设备健康度严重不足，需要立即检查",
		})
	} else if healthScore < 50 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "warning",
			"type":    "health_score",
			"message": "设备健康度偏低，建议进行维护",
		})
	}

	// CPU告警
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage > 95 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "critical",
				"type":    "cpu_usage",
				"message": "CPU使用率过高，可能导致系统不稳定",
			})
		} else if cpuUsage > 85 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "warning",
				"type":    "cpu_usage",
				"message": "CPU使用率较高，建议优化进程",
			})
		}
	}

	// 内存告警
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage > 95 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "critical",
				"type":    "memory_usage",
				"message": "内存使用率过高，可能导致系统崩溃",
			})
		} else if memUsage > 85 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "warning",
				"type":    "memory_usage",
				"message": "内存使用率较高，建议清理缓存",
			})
		}
	}

	// 磁盘告警
	if diskUsage, ok := deviceData["disk_usage"].(float64); ok {
		if diskUsage > 95 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "critical",
				"type":    "disk_usage",
				"message": "磁盘空间严重不足，需要立即清理",
			})
		} else if diskUsage > 90 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "warning",
				"type":    "disk_usage",
				"message": "磁盘空间不足，建议清理文件",
			})
		}
	}

	return alerts
}

// ===============================
// 设备数据验证和处理
// ===============================

// ValidateDeviceData 验证设备数据完整性
func (s *sDevice) ValidateDeviceData(ctx context.Context, deviceData map[string]interface{}) error {
	// 验证必需字段
	requiredFields := []string{"device_id", "name", "type", "status"}
	for _, field := range requiredFields {
		if value, ok := deviceData[field]; !ok || value == nil || value == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("缺少必需字段: %s", field))
		}
	}

	// 验证设备类型
	if deviceType, ok := deviceData["type"].(string); ok {
		validTypes := []string{"sensor", "camera", "actuator", "gateway"}
		if !contains(validTypes, deviceType) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "无效的设备类型")
		}
	}

	// 验证设备状态
	if status, ok := deviceData["status"].(string); ok {
		validStatuses := []string{"online", "offline", "maintenance", "error", "busy"}
		if !contains(validStatuses, status) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "无效的设备状态")
		}
	}

	// 验证IP地址
	if ipAddress, ok := deviceData["ip_address"].(string); ok && ipAddress != "" {
		if !isValidIPAddress(ipAddress) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "IP地址格式不正确")
		}
	}

	return nil
}

// ProcessDeviceData 处理设备数据
func (s *sDevice) ProcessDeviceData(ctx context.Context, rawData map[string]interface{}) (map[string]interface{}, error) {
	// 验证数据
	if err := s.ValidateDeviceData(ctx, rawData); err != nil {
		return nil, err
	}

	// 处理数据
	processedData := make(map[string]interface{})
	for k, v := range rawData {
		processedData[k] = v
	}

	// 计算健康度
	processedData["health_score"] = s.CalculateDeviceHealthScore(ctx, processedData)

	// 确定状态
	processedData["status"] = s.DetermineDeviceStatus(ctx, processedData)

	// 添加时间戳
	processedData["processed_at"] = gtime.Now()

	return processedData, nil
}

// ===============================
// 辅助方法
// ===============================

// getCheckStatus 获取检查状态
func (s *sDevice) getCheckStatus(value, warningThreshold, criticalThreshold float64) string {
	if value >= criticalThreshold {
		return "critical"
	} else if value >= warningThreshold {
		return "warning"
	}
	return "normal"
}

// getCPUMessage 获取CPU消息
func (s *sDevice) getCPUMessage(cpuUsage float64) string {
	if cpuUsage > 95 {
		return "CPU使用率过高，需要立即处理"
	} else if cpuUsage > 85 {
		return "CPU使用率较高，建议优化"
	}
	return "CPU使用率正常"
}

// getMemoryMessage 获取内存消息
func (s *sDevice) getMemoryMessage(memUsage float64) string {
	if memUsage > 95 {
		return "内存使用率过高，可能导致系统崩溃"
	} else if memUsage > 85 {
		return "内存使用率较高，建议清理缓存"
	}
	return "内存使用率正常"
}

// getDiskMessage 获取磁盘消息
func (s *sDevice) getDiskMessage(diskUsage float64) string {
	if diskUsage > 95 {
		return "磁盘空间严重不足，需要立即清理"
	} else if diskUsage > 90 {
		return "磁盘空间不足，建议清理文件"
	}
	return "磁盘空间充足"
}

// getNetworkStatus 获取网络状态
func (s *sDevice) getNetworkStatus(networkStatus string) string {
	switch networkStatus {
	case "connected":
		return "normal"
	case "unstable":
		return "warning"
	case "disconnected":
		return "critical"
	default:
		return "unknown"
	}
}

// getNetworkMessage 获取网络消息
func (s *sDevice) getNetworkMessage(networkStatus string) string {
	switch networkStatus {
	case "connected":
		return "网络连接正常"
	case "unstable":
		return "网络连接不稳定"
	case "disconnected":
		return "网络连接断开"
	default:
		return "网络状态未知"
	}
}
