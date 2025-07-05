package task

import (
	"context"
	"math"
	"sort"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务分配相关业务逻辑
// ===============================

// AssignTaskToDevice 分配任务给设备
func (s *sTask) AssignTaskToDevice(ctx context.Context, taskId string, deviceIds []string, strategy string, force bool) (map[string]interface{}, error) {
	// 验证任务存在
	if err := s.validateTaskExists(ctx, taskId); err != nil {
		return nil, err
	}

	// 获取任务信息
	taskInfo, err := s.getTaskInfo(ctx, taskId)
	if err != nil {
		return nil, err
	}

	// 获取可用设备列表
	availableDevices, err := s.getAvailableDevices(ctx, deviceIds)
	if err != nil {
		return nil, err
	}

	// 根据策略选择最佳设备
	selectedDevice, assignmentScore, reason, err := s.selectBestDevice(ctx, taskInfo, availableDevices, strategy, force)
	if err != nil {
		return nil, err
	}

	// 执行任务分配
	deviceId := s.getStringValue(selectedDevice, "device_id", "")
	if err := s.executeTaskAssignment(ctx, taskId, deviceId); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"taskId":           taskId,
		"assignedDeviceId": deviceId,
		"assignmentScore":  assignmentScore,
		"reason":           reason,
		"status":           "assigned",
	}

	return result, nil
}

// selectBestDevice 选择最佳设备
func (s *sTask) selectBestDevice(ctx context.Context, taskInfo map[string]interface{}, devices []map[string]interface{}, strategy string, force bool) (map[string]interface{}, float64, string, error) {
	// 过滤兼容设备
	compatibleDevices := s.filterCompatibleDevices(taskInfo, devices)
	if len(compatibleDevices) == 0 {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有兼容的设备可用")
	}

	// 根据策略计算分配评分
	switch strategy {
	case "auto":
		return s.calculateAutoAssignment(compatibleDevices, taskInfo)
	case "manual":
		return s.calculateManualAssignment(compatibleDevices, taskInfo)
	case "load_balanced":
		return s.calculateLoadBalancedAssignment(compatibleDevices, taskInfo)
	case "priority":
		return s.calculatePriorityAssignment(compatibleDevices, taskInfo)
	default:
		return s.calculateAutoAssignment(compatibleDevices, taskInfo)
	}
}

// filterCompatibleDevices 过滤兼容设备
func (s *sTask) filterCompatibleDevices(taskInfo map[string]interface{}, devices []map[string]interface{}) []map[string]interface{} {
	var compatibleDevices []map[string]interface{}
	taskType := s.getStringValue(taskInfo, "type", "")
	requiredCapabilities := s.getStringValue(taskInfo, "required_capabilities", "")

	for _, device := range devices {
		deviceType := s.getStringValue(device, "type", "")
		deviceCapabilities := s.getStringValue(device, "capabilities", "")
		deviceStatus := s.getStringValue(device, "status", "")

		// 检查设备状态
		if deviceStatus != "online" && deviceStatus != "idle" {
			continue
		}

		// 检查设备类型兼容性
		if taskType != "" && deviceType != "" && taskType != deviceType {
			continue
		}

		// 检查设备能力兼容性
		if requiredCapabilities != "" && deviceCapabilities != "" {
			if !s.checkCapabilityCompatibility(requiredCapabilities, deviceCapabilities) {
				continue
			}
		}

		compatibleDevices = append(compatibleDevices, device)
	}

	return compatibleDevices
}

// calculateAutoAssignment 自动分配策略
func (s *sTask) calculateAutoAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	for _, device := range devices {
		// 计算综合评分
		loadScore := s.calculateDeviceLoadScore(device)
		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
		availabilityScore := s.calculateAvailabilityScore(device)
		performanceScore := s.calculatePerformanceScore(device)

		// 加权综合评分
		totalScore := loadScore*0.3 + compatibilityScore*0.3 + availabilityScore*0.2 + performanceScore*0.2

		if totalScore > bestScore {
			bestScore = totalScore
			bestDevice = device
			reason = "自动选择综合评分最高的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculateLoadBalancedAssignment 负载均衡分配策略
func (s *sTask) calculateLoadBalancedAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	for _, device := range devices {
		// 主要考虑负载评分
		loadScore := s.calculateDeviceLoadScore(device)
		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)

		// 负载越低分数越高
		balancedScore := (1.0-loadScore)*0.7 + compatibilityScore*0.3

		if balancedScore > bestScore {
			bestScore = balancedScore
			bestDevice = device
			reason = "负载均衡策略选择负载最低的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculatePriorityAssignment 优先级分配策略
func (s *sTask) calculatePriorityAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	taskPriority := s.getIntValue(taskInfo, "priority", 0)

	for _, device := range devices {
		// 主要考虑设备优先级和任务优先级匹配
		devicePriority := s.getIntValue(device, "priority", 0)
		priorityMatch := 1.0 - math.Abs(float64(taskPriority-devicePriority))/100.0

		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
		loadScore := s.calculateDeviceLoadScore(device)

		// 优先级匹配度权重最高
		totalScore := priorityMatch*0.5 + compatibilityScore*0.3 + (1.0-loadScore)*0.2

		if totalScore > bestScore {
			bestScore = totalScore
			bestDevice = device
			reason = "优先级策略选择优先级最匹配的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculateManualAssignment 手动分配策略
func (s *sTask) calculateManualAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	// 手动分配通常选择第一个可用设备
	if len(devices) > 0 {
		device := devices[0]
		score := s.calculateDeviceLoadScore(device)
		return device, score, "手动分配选择第一个可用设备", nil
	}

	return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有可用的设备")
}

// calculateDeviceLoadScore 计算设备负载评分
func (s *sTask) calculateDeviceLoadScore(device map[string]interface{}) float64 {
	cpuUsage := s.getFloatValue(device, "cpu_usage", 0.0)
	memoryUsage := s.getFloatValue(device, "memory_usage", 0.0)
	currentTasks := s.getIntValue(device, "current_tasks", 0)
	maxTasks := s.getIntValue(device, "max_tasks", 1)

	// 计算综合负载评分 (0-1, 越高表示负载越重)
	cpuScore := cpuUsage / 100.0
	memoryScore := memoryUsage / 100.0
	taskScore := float64(currentTasks) / float64(maxTasks)

	// 加权平均
	loadScore := cpuScore*0.4 + memoryScore*0.3 + taskScore*0.3
	return math.Min(loadScore, 1.0)
}

// calculateCompatibilityScore 计算兼容性评分
func (s *sTask) calculateCompatibilityScore(device map[string]interface{}, taskInfo map[string]interface{}) float64 {
	deviceType := s.getStringValue(device, "type", "")
	taskType := s.getStringValue(taskInfo, "type", "")
	deviceCapabilities := s.getStringValue(device, "capabilities", "")
	taskCapabilities := s.getStringValue(taskInfo, "required_capabilities", "")

	score := 0.0

	// 类型匹配
	if deviceType == taskType {
		score += 0.5
	}

	// 能力匹配
	if deviceCapabilities != "" && taskCapabilities != "" {
		if s.checkCapabilityCompatibility(taskCapabilities, deviceCapabilities) {
			score += 0.5
		}
	}

	return score
}

// calculateAvailabilityScore 计算可用性评分
func (s *sTask) calculateAvailabilityScore(device map[string]interface{}) float64 {
	status := s.getStringValue(device, "status", "")
	uptime := s.getFloatValue(device, "uptime_hours", 0.0)
	healthScore := s.getFloatValue(device, "health_score", 0.0)

	score := 0.0

	// 状态评分
	switch status {
	case "online":
		score += 0.4
	case "idle":
		score += 0.3
	case "busy":
		score += 0.2
	default:
		score += 0.0
	}

	// 运行时间评分
	if uptime > 24 {
		score += 0.3
	} else if uptime > 1 {
		score += 0.2
	} else {
		score += 0.1
	}

	// 健康评分
	score += healthScore / 100.0 * 0.3

	return math.Min(score, 1.0)
}

// calculatePerformanceScore 计算性能评分
func (s *sTask) calculatePerformanceScore(device map[string]interface{}) float64 {
	// 基于设备性能指标计算评分
	performanceScore := s.getFloatValue(device, "performance_score", 50.0)
	return performanceScore / 100.0
}

// checkCapabilityCompatibility 检查能力兼容性
func (s *sTask) checkCapabilityCompatibility(required, available string) bool {
	// 简单的字符串包含检查，实际应该更复杂的匹配逻辑
	return available != "" && (required == "" || available == required)
}

// GetTaskAssignmentOptions 获取任务分配选项
func (s *sTask) GetTaskAssignmentOptions(ctx context.Context, taskId string) (map[string]interface{}, error) {
	// 验证任务存在
	if err := s.validateTaskExists(ctx, taskId); err != nil {
		return nil, err
	}

	// 获取任务信息
	taskInfo, err := s.getTaskInfo(ctx, taskId)
	if err != nil {
		return nil, err
	}

	// 获取所有可用设备
	allDevices, err := s.getAllAvailableDevices(ctx)
	if err != nil {
		return nil, err
	}

	// 计算每个设备的分配选项
	var options []map[string]interface{}
	for _, device := range allDevices {
		option := s.calculateDeviceAssignmentOption(device, taskInfo)
		if option != nil {
			options = append(options, option)
		}
	}

	// 按分配评分排序
	sort.Slice(options, func(i, j int) bool {
		scoreI := s.getFloatValue(options[i], "assignmentScore", 0.0)
		scoreJ := s.getFloatValue(options[j], "assignmentScore", 0.0)
		return scoreI > scoreJ
	})

	result := map[string]interface{}{
		"taskId":  taskId,
		"options": options,
		"total":   len(options),
	}

	return result, nil
}

// calculateDeviceAssignmentOption 计算设备分配选项
func (s *sTask) calculateDeviceAssignmentOption(device map[string]interface{}, taskInfo map[string]interface{}) map[string]interface{} {
	deviceId := s.getStringValue(device, "device_id", "")
	if deviceId == "" {
		return nil
	}

	// 计算各项评分
	loadScore := s.calculateDeviceLoadScore(device)
	compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
	availabilityScore := s.calculateAvailabilityScore(device)
	performanceScore := s.calculatePerformanceScore(device)

	// 计算综合分配评分
	assignmentScore := loadScore*0.3 + compatibilityScore*0.3 + availabilityScore*0.2 + performanceScore*0.2

	// 生成推荐
	recommendation := s.generateAssignmentRecommendation(assignmentScore, loadScore, compatibilityScore)

	option := map[string]interface{}{
		"deviceId":           deviceId,
		"deviceName":         s.getStringValue(device, "name", ""),
		"deviceType":         s.getStringValue(device, "type", ""),
		"loadScore":          loadScore * 100,
		"compatibilityScore": compatibilityScore * 100,
		"assignmentScore":    assignmentScore * 100,
		"currentTasks":       s.getIntValue(device, "current_tasks", 0),
		"maxTasks":           s.getIntValue(device, "max_tasks", 1),
		"estimatedWaitTime":  s.calculateEstimatedWaitTime(device),
		"status":             s.getStringValue(device, "status", ""),
		"recommendation":     recommendation,
	}

	return option
}

// generateAssignmentRecommendation 生成分配推荐
func (s *sTask) generateAssignmentRecommendation(assignmentScore, loadScore, compatibilityScore float64) string {
	if assignmentScore > 0.8 {
		return "强烈推荐"
	} else if assignmentScore > 0.6 {
		return "推荐"
	} else if assignmentScore > 0.4 {
		return "一般"
	} else {
		return "不推荐"
	}
}

// calculateEstimatedWaitTime 计算预估等待时间
func (s *sTask) calculateEstimatedWaitTime(device map[string]interface{}) int {
	currentTasks := s.getIntValue(device, "current_tasks", 0)
	avgTaskDuration := s.getFloatValue(device, "avg_task_duration", 300.0) // 默认5分钟

	// 简单计算：当前任务数 * 平均任务时长
	estimatedTime := int(float64(currentTasks) * avgTaskDuration)
	return estimatedTime
}

// BatchAssignTasks 批量分配任务
func (s *sTask) BatchAssignTasks(ctx context.Context, assignments []map[string]interface{}, strategy string) (map[string]interface{}, error) {
	var results []map[string]interface{}
	success := 0
	failed := 0

	for _, assignment := range assignments {
		taskId := s.getStringValue(assignment, "taskId", "")
		deviceIds := s.getStringSliceValue(assignment, "deviceIds", []string{})

		if taskId == "" || len(deviceIds) == 0 {
			failed++
			results = append(results, map[string]interface{}{
				"taskId": taskId,
				"status": "failed",
				"error":  "任务ID或设备ID列表为空",
			})
			continue
		}

		// 执行单个任务分配
		result, err := s.AssignTaskToDevice(ctx, taskId, deviceIds, strategy, false)
		if err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"taskId": taskId,
				"status": "failed",
				"error":  err.Error(),
			})
		} else {
			success++
			results = append(results, result)
		}
	}

	return map[string]interface{}{
		"total":       len(assignments),
		"success":     success,
		"failed":      failed,
		"assignments": results,
	}, nil
}

// executeTaskAssignment 执行任务分配
func (s *sTask) executeTaskAssignment(ctx context.Context, taskId, deviceId string) error {
	// 这里应该执行实际的任务分配逻辑
	// 包括更新数据库、发送通知等

	g.Log().Info(ctx, "执行任务分配", g.Map{
		"task_id":   taskId,
		"device_id": deviceId,
		"time":      gtime.Now(),
	})

	return nil
}

// 工具方法
func (s *sTask) getStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

func (s *sTask) getFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key]; ok {
		if floatValue, ok := value.(float64); ok {
			return floatValue
		}
	}
	return defaultValue
}

func (s *sTask) getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

func (s *sTask) getStringSliceValue(data map[string]interface{}, key string, defaultValue []string) []string {
	if value, ok := data[key]; ok {
		if sliceValue, ok := value.([]string); ok {
			return sliceValue
		}
	}
	return defaultValue
}

// 模拟方法（实际应该从数据库获取）
func (s *sTask) validateTaskExists(ctx context.Context, taskId string) error {
	// 这里应该验证任务是否存在
	return nil
}

func (s *sTask) getTaskInfo(ctx context.Context, taskId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取任务信息
	return map[string]interface{}{
		"task_id":               taskId,
		"type":                  "data_processing",
		"priority":              5,
		"required_capabilities": "high_performance",
	}, nil
}

func (s *sTask) getAvailableDevices(ctx context.Context, deviceIds []string) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取设备信息
	var devices []map[string]interface{}
	for _, deviceId := range deviceIds {
		device := map[string]interface{}{
			"device_id":         deviceId,
			"name":              "device-" + deviceId,
			"type":              "data_processing",
			"status":            "online",
			"cpu_usage":         45.5,
			"memory_usage":      60.2,
			"current_tasks":     2,
			"max_tasks":         10,
			"capabilities":      "high_performance",
			"priority":          5,
			"uptime_hours":      168.5,
			"health_score":      85.0,
			"performance_score": 78.0,
			"avg_task_duration": 300.0,
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *sTask) getAllAvailableDevices(ctx context.Context) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取所有可用设备
	return []map[string]interface{}{
		{
			"device_id":         "dev001",
			"name":              "device-001",
			"type":              "data_processing",
			"status":            "online",
			"cpu_usage":         45.5,
			"memory_usage":      60.2,
			"current_tasks":     2,
			"max_tasks":         10,
			"capabilities":      "high_performance",
			"priority":          5,
			"uptime_hours":      168.5,
			"health_score":      85.0,
			"performance_score": 78.0,
			"avg_task_duration": 300.0,
		},
		{
			"device_id":         "dev002",
			"name":              "device-002",
			"type":              "data_processing",
			"status":            "idle",
			"cpu_usage":         25.0,
			"memory_usage":      40.0,
			"current_tasks":     0,
			"max_tasks":         8,
			"capabilities":      "high_performance",
			"priority":          3,
			"uptime_hours":      120.0,
			"health_score":      90.0,
			"performance_score": 85.0,
			"avg_task_duration": 250.0,
		},
	}, nil
}
