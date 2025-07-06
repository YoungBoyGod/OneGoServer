package task

import (
	"context"
	"math"
	"sort"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	task "OneGfServer/internal/model/task"
)

// ===============================
// 任务分配相关业务逻辑
// ===============================

// AssignTaskToDevice 将任务分配给设备
func (s *sTask) AssignTaskToDevice(ctx context.Context, taskId string, deviceIds []string, strategy string, force bool) (map[string]interface{}, error) {
	// 验证任务是否存在
	if err := s.validateTaskExists(ctx, taskId); err != nil {
		return nil, err
	}

	// 获取任务信息
	taskInfo, err := s.getTaskInfo(ctx, taskId)
	if err != nil {
		return nil, err
	}

	// 获取可用设备
	var devices []map[string]interface{}
	if len(deviceIds) > 0 {
		devices, err = s.getAvailableDevices(ctx, deviceIds)
	} else {
		devices, err = s.getAllAvailableDevices(ctx)
	}
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, gerror.NewCode(gcode.CodeResourceExhausted, "没有可用的设备")
	}

	// 选择最佳设备
	selectedDevice, score, reason, err := s.selectBestDevice(ctx, taskInfo, devices, strategy, force)
	if err != nil {
		return nil, err
	}

	// 执行任务分配
	if err := s.executeTaskAssignment(ctx, taskId, selectedDevice["device_id"].(string)); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"task_id":           taskId,
		"assigned_device":   selectedDevice,
		"assignment_score":  score,
		"assignment_reason": reason,
		"strategy":          strategy,
		"timestamp":         gtime.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// selectBestDevice 选择最佳设备
func (s *sTask) selectBestDevice(ctx context.Context, taskInfo map[string]interface{}, devices []map[string]interface{}, strategy string, force bool) (map[string]interface{}, float64, string, error) {
	// 过滤兼容设备
	compatibleDevices := s.filterCompatibleDevices(taskInfo, devices)
	if len(compatibleDevices) == 0 {
		return nil, 0, "", gerror.NewCode(gcode.CodeResourceExhausted, "没有兼容的设备")
	}

	// 根据策略选择设备
	var selectedDevice map[string]interface{}
	var score float64
	var reason string
	var err error

	switch strategy {
	case "auto":
		selectedDevice, score, reason, err = s.calculateAutoAssignment(compatibleDevices, taskInfo)
	case "load_balanced":
		selectedDevice, score, reason, err = s.calculateLoadBalancedAssignment(compatibleDevices, taskInfo)
	case "priority":
		selectedDevice, score, reason, err = s.calculatePriorityAssignment(compatibleDevices, taskInfo)
	case "manual":
		selectedDevice, score, reason, err = s.calculateManualAssignment(compatibleDevices, taskInfo)
	default:
		return nil, 0, "", gerror.NewCode(gcode.CodeInvalidParameter, "不支持的分配策略")
	}

	if err != nil {
		return nil, 0, "", err
	}

	return selectedDevice, score, reason, nil
}

// filterCompatibleDevices 过滤兼容设备
func (s *sTask) filterCompatibleDevices(taskInfo map[string]interface{}, devices []map[string]interface{}) []map[string]interface{} {
	var compatibleDevices []map[string]interface{}

	for _, device := range devices {
		// 检查设备状态
		if status, ok := device["status"].(string); ok && status != "online" {
			continue
		}

		// 检查设备能力
		if taskType, ok := taskInfo["task_type"].(string); ok {
			if deviceType, ok := device["device_type"].(string); ok {
				if !s.checkCapabilityCompatibility(taskType, deviceType) {
					continue
				}
			}
		}

		// 检查资源可用性
		if !s.checkResourceAvailability(device, taskInfo) {
			continue
		}

		compatibleDevices = append(compatibleDevices, device)
	}

	return compatibleDevices
}

// calculateAutoAssignment 自动分配策略
func (s *sTask) calculateAutoAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64

	for _, device := range devices {
		score := s.calculateDeviceAssignmentScore(device, taskInfo)
		if score > bestScore {
			bestScore = score
			bestDevice = device
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.New("无法找到合适的设备")
	}

	return bestDevice, bestScore, "自动选择最佳匹配设备", nil
}

// calculateLoadBalancedAssignment 负载均衡分配策略
func (s *sTask) calculateLoadBalancedAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestLoadScore float64

	for _, device := range devices {
		loadScore := s.calculateDeviceLoadScore(device)
		if bestDevice == nil || loadScore < bestLoadScore {
			bestLoadScore = loadScore
			bestDevice = device
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.New("无法找到合适的设备")
	}

	return bestDevice, bestLoadScore, "选择负载最低的设备", nil
}

// calculatePriorityAssignment 优先级分配策略
func (s *sTask) calculatePriorityAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	// 按设备优先级排序
	sort.Slice(devices, func(i, j int) bool {
		priorityI := s.getIntValue(devices[i], "priority", 5)
		priorityJ := s.getIntValue(devices[j], "priority", 5)
		return priorityI > priorityJ
	})

	if len(devices) == 0 {
		return nil, 0, "", gerror.New("没有可用的设备")
	}

	selectedDevice := devices[0]
	priority := s.getIntValue(selectedDevice, "priority", 5)

	return selectedDevice, float64(priority), "按设备优先级分配", nil
}

// calculateManualAssignment 手动分配策略
func (s *sTask) calculateManualAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	// 手动分配通常由用户指定设备
	// 这里返回第一个可用设备
	if len(devices) == 0 {
		return nil, 0, "", gerror.New("没有可用的设备")
	}

	return devices[0], 100.0, "手动分配", nil
}

// calculateDeviceAssignmentScore 计算设备分配评分
func (s *sTask) calculateDeviceAssignmentScore(device map[string]interface{}, taskInfo map[string]interface{}) float64 {
	score := 0.0

	// 负载评分 (权重: 40%)
	loadScore := s.calculateDeviceLoadScore(device)
	score += loadScore * 0.4

	// 兼容性评分 (权重: 30%)
	compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
	score += compatibilityScore * 0.3

	// 可用性评分 (权重: 20%)
	availabilityScore := s.calculateAvailabilityScore(device)
	score += availabilityScore * 0.2

	// 性能评分 (权重: 10%)
	performanceScore := s.calculatePerformanceScore(device)
	score += performanceScore * 0.1

	return score
}

// calculateDeviceLoadScore 计算设备负载评分
func (s *sTask) calculateDeviceLoadScore(device map[string]interface{}) float64 {
	score := 100.0

	// CPU负载影响
	if cpuUsage, ok := device["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			score -= 40
		} else if cpuUsage > 70 {
			score -= 20
		} else if cpuUsage > 50 {
			score -= 10
		}
	}

	// 内存负载影响
	if memUsage, ok := device["memory_usage"].(float64); ok {
		if memUsage > 90 {
			score -= 30
		} else if memUsage > 80 {
			score -= 15
		} else if memUsage > 60 {
			score -= 5
		}
	}

	// 任务负载影响
	if currentTasks, ok := device["running_task_count"].(int); ok {
		if maxTasks, ok := device["max_concurrent_tasks"].(int); ok && maxTasks > 0 {
			utilization := float64(currentTasks) / float64(maxTasks)
			if utilization > 0.9 {
				score -= 30
			} else if utilization > 0.7 {
				score -= 15
			} else if utilization > 0.5 {
				score -= 5
			}
		}
	}

	return math.Max(0, score)
}

// calculateCompatibilityScore 计算兼容性评分
func (s *sTask) calculateCompatibilityScore(device map[string]interface{}, taskInfo map[string]interface{}) float64 {
	score := 50.0

	// 检查设备类型兼容性
	if taskType, ok := taskInfo["task_type"].(string); ok {
		if deviceType, ok := device["device_type"].(string); ok {
			if s.checkCapabilityCompatibility(taskType, deviceType) {
				score += 30
			}
		}
	}

	// 检查资源需求兼容性
	if s.checkResourceCompatibility(device, taskInfo) {
		score += 20
	}

	return math.Min(score, 100)
}

// calculateAvailabilityScore 计算可用性评分
func (s *sTask) calculateAvailabilityScore(device map[string]interface{}) float64 {
	score := 100.0

	// 检查设备状态
	if status, ok := device["status"].(string); ok {
		switch status {
		case "online":
			score = 100
		case "busy":
			score = 80
		case "maintenance":
			score = 30
		case "error":
			score = 20
		case "offline":
			score = 0
		default:
			score = 50
		}
	}

	// 检查维护模式
	if maintenanceMode, ok := device["maintenance_mode"].(bool); ok && maintenanceMode {
		score *= 0.5
	}

	return score
}

// calculatePerformanceScore 计算性能评分
func (s *sTask) calculatePerformanceScore(device map[string]interface{}) float64 {
	score := 100.0

	// 基于错误次数调整评分
	if errorCount, ok := device["error_count"].(int); ok {
		if errorCount > 10 {
			score -= 40
		} else if errorCount > 5 {
			score -= 20
		} else if errorCount > 2 {
			score -= 10
		}
	}

	// 基于网络延迟调整评分
	if latency, ok := device["network_latency"].(float64); ok {
		if latency > 200 {
			score -= 30
		} else if latency > 100 {
			score -= 15
		} else if latency > 50 {
			score -= 5
		}
	}

	return math.Max(0, score)
}

// checkCapabilityCompatibility 检查能力兼容性
func (s *sTask) checkCapabilityCompatibility(required, available string) bool {
	// 定义兼容性映射
	compatibilityMap := map[string][]string{
		"compute":     {"server", "workstation", "edge"},
		"storage":     {"server"},
		"network":     {"server", "edge"},
		"display":     {"workstation"},
		"interactive": {"workstation"},
		"mobile":      {"mobile"},
		"location":    {"mobile", "iot"},
		"sensor":      {"iot", "mobile"},
		"monitor":     {"iot", "edge"},
		"control":     {"iot", "edge"},
		"local":       {"edge"},
	}

	if compatibleTypes, exists := compatibilityMap[required]; exists {
		for _, compatibleType := range compatibleTypes {
			if compatibleType == available {
				return true
			}
		}
	}

	return false
}

// checkResourceAvailability 检查资源可用性
func (s *sTask) checkResourceAvailability(device map[string]interface{}, taskInfo map[string]interface{}) bool {
	// 检查CPU资源
	if taskCPU, ok := taskInfo["required_cpu"].(float64); ok {
		if deviceCPU, ok := device["cpu_usage"].(float64); ok {
			if deviceCPU+taskCPU > 100 {
				return false
			}
		}
	}

	// 检查内存资源
	if taskMem, ok := taskInfo["required_memory"].(float64); ok {
		if deviceMem, ok := device["memory_usage"].(float64); ok {
			if deviceMem+taskMem > 100 {
				return false
			}
		}
	}

	// 检查任务数量限制
	if currentTasks, ok := device["running_task_count"].(int); ok {
		if maxTasks, ok := device["max_concurrent_tasks"].(int); ok {
			if currentTasks >= maxTasks {
				return false
			}
		}
	}

	return true
}

// checkResourceCompatibility 检查资源兼容性
func (s *sTask) checkResourceCompatibility(device map[string]interface{}, taskInfo map[string]interface{}) bool {
	// 检查是否有足够的资源余量
	if taskCPU, ok := taskInfo["required_cpu"].(float64); ok {
		if deviceCPU, ok := device["cpu_usage"].(float64); ok {
			if deviceCPU+taskCPU > 90 { // 留10%余量
				return false
			}
		}
	}

	if taskMem, ok := taskInfo["required_memory"].(float64); ok {
		if deviceMem, ok := device["memory_usage"].(float64); ok {
			if deviceMem+taskMem > 85 { // 留15%余量
				return false
			}
		}
	}

	return true
}

// 辅助方法
func (s *sTask) getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(int); ok {
		return value
	}
	return defaultValue
}

func (s *sTask) validateTaskExists(ctx context.Context, taskId string) error {
	// 这里应该检查任务是否存在
	if taskId == "" {
		return gerror.New("任务ID不能为空")
	}
	return nil
}

func (s *sTask) getTaskInfo(ctx context.Context, taskId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取任务信息
	return map[string]interface{}{
		"task_id":         taskId,
		"task_type":       "compute",
		"required_cpu":    20.0,
		"required_memory": 30.0,
	}, nil
}

func (s *sTask) getAvailableDevices(ctx context.Context, deviceIds []string) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取指定设备信息
	var devices []map[string]interface{}
	for _, deviceId := range deviceIds {
		devices = append(devices, map[string]interface{}{
			"device_id":            deviceId,
			"status":               "online",
			"cpu_usage":            45.0,
			"memory_usage":         62.0,
			"running_task_count":   2,
			"max_concurrent_tasks": 5,
		})
	}
	return devices, nil
}

func (s *sTask) getAllAvailableDevices(ctx context.Context) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取所有可用设备
	return []map[string]interface{}{
		{
			"device_id":            "device_1",
			"status":               "online",
			"cpu_usage":            45.0,
			"memory_usage":         62.0,
			"running_task_count":   2,
			"max_concurrent_tasks": 5,
		},
		{
			"device_id":            "device_2",
			"status":               "online",
			"cpu_usage":            30.0,
			"memory_usage":         50.0,
			"running_task_count":   1,
			"max_concurrent_tasks": 5,
		},
	}, nil
}

func (s *sTask) executeTaskAssignment(ctx context.Context, taskId, deviceId string) error {
	// 这里应该执行实际的任务分配操作
	// 更新数据库中的任务分配信息
	return nil
}

// ===============================
// 任务分配业务逻辑
// ===============================

// AssignTask 分配任务
func (s *sTask) AssignTask(ctx context.Context, input *task.AssignTaskInput) (*task.AssignTaskOutput, error) {
	// 检查任务状态
	taskStatusInput := &task.GetTaskStatusInput{TaskID: input.TaskID}
	taskStatusOutput := s.getTaskStatus(taskStatusInput)

	if taskStatusOutput.Status == "running" || taskStatusOutput.Status == "completed" {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "任务已在运行或已完成，无法重新分配")
	}

	// 检查设备可用性
	if !s.isDeviceAvailable(input.DeviceID) {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "设备不可用")
	}

	// 检查任务是否已被分配
	if s.isTaskAssigned(input.TaskID) {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "任务已被分配")
	}

	// 这里应该更新数据库中的任务分配
	// 目前返回模拟结果

	return &task.AssignTaskOutput{
		TaskID:     input.TaskID,
		DeviceID:   input.DeviceID,
		Status:     "assigned",
		AssignedAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:    "任务分配成功",
	}, nil
}

// UnassignTask 取消分配任务
func (s *sTask) UnassignTask(ctx context.Context, input *task.UnassignTaskInput) (*task.UnassignTaskOutput, error) {
	// 检查任务状态
	taskStatusInput := &task.GetTaskStatusInput{TaskID: input.TaskID}
	taskStatusOutput := s.getTaskStatus(taskStatusInput)

	if taskStatusOutput.Status == "running" {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "运行中的任务无法取消分配")
	}

	// 检查任务是否已被分配
	if !s.isTaskAssigned(input.TaskID) {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation, "任务未被分配")
	}

	// 这里应该更新数据库中的任务分配
	// 目前返回模拟结果

	return &task.UnassignTaskOutput{
		TaskID:       input.TaskID,
		Status:       "unassigned",
		UnassignedAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:      "任务取消分配成功",
	}, nil
}

// isDeviceAvailable 检查设备是否可用（内部方法）
func (s *sTask) isDeviceAvailable(deviceID string) bool {
	// 这里应该检查设备状态
	// 目前返回模拟结果
	return true
}

// isTaskAssigned 检查任务是否已被分配（内部方法）
func (s *sTask) isTaskAssigned(taskID string) bool {
	// 这里应该检查任务分配状态
	// 目前返回模拟结果
	return false
}

// getTaskStatus 获取任务状态（内部方法）
func (s *sTask) getTaskStatus(input *task.GetTaskStatusInput) *task.GetTaskStatusOutput {
	// 这里应该从数据库获取任务状态
	// 目前返回模拟数据
	return &task.GetTaskStatusOutput{
		TaskID:    input.TaskID,
		Status:    "pending",
		Progress:  0.0,
		StartTime: "",
		EndTime:   "",
		Duration:  "",
		Details:   map[string]interface{}{},
	}
}
