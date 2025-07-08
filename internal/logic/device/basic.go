package device

import (
	"OneGfServer/internal/model/device"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备基础业务逻辑
// ===============================

// HandleDeviceRegistration 处理设备注册
func (s *sDevice) HandleDeviceRegistration(ctx context.Context, input *device.HandleDeviceRegistrationInput) (*device.HandleDeviceRegistrationOutput, error) {
	// output := &device.HandleDeviceRegistrationOutput{
	// 	DeviceData: make(map[string]interface{}),
	// 	IsSuccess:  false,
	// }

	// // 验证设备注册数据
	// validateInput := &device.ValidateDeviceRegistrationInput{
	// 	Name:       s.getStringValue(input.DeviceData, "name", ""),
	// 	MacAddress: s.getStringValue(input.DeviceData, "mac_address", ""),
	// 	IPAddress:  s.getStringValue(input.DeviceData, "ip_address", ""),
	// 	DeviceType: s.getStringValue(input.DeviceData, "device_type", ""),
	// 	Model:      s.getStringValue(input.DeviceData, "model", ""),
	// 	Protocol:   s.getStringValue(input.DeviceData, "protocol", ""),
	// 	Port:       s.getIntValue(input.DeviceData, "port", 0),
	// }

	// validateOutput, err := s.ValidateDeviceRegistration(ctx, validateInput)
	// if err != nil {
	// 	return nil, err
	// }

	// if !validateOutput.IsValid {
	// 	output.Message = "设备注册验证失败: " + validateOutput.Message
	// 	return output, nil
	// }

	// // 生成设备ID（如果未提供）
	// if _, ok := input.DeviceData["device_id"]; !ok {
	// 	input.DeviceData["device_id"] = s.generateDeviceId()
	// }

	// // 设置注册时间
	// input.DeviceData["registered_at"] = gtime.Now().Format("2006-01-02 15:04:05")
	// input.DeviceData["last_heartbeat"] = gtime.Now().Format("2006-01-02 15:04:05")
	// input.DeviceData["status"] = "online"

	// // 初始化设备指标
	// input.DeviceData["cpu_usage"] = 0.0
	// input.DeviceData["memory_usage"] = 0.0
	// input.DeviceData["disk_usage"] = 0.0
	// input.DeviceData["network_latency"] = 0.0
	// input.DeviceData["running_task_count"] = 0
	// input.DeviceData["error_count"] = 0

	// // 复制处理后的数据
	// for key, value := range input.DeviceData {
	// 	output.DeviceData[key] = value
	// }

	// output.DeviceID = s.getStringValue(input.DeviceData, "device_id", "")
	// output.Message = "设备注册成功"
	// output.IsSuccess = true

	// return output, nil
	return nil, nil
}

// HandleDeviceHeartbeat 处理设备心跳
func (s *sDevice) HandleDeviceHeartbeat(ctx context.Context, input *device.HandleDeviceHeartbeatInput) (*device.HandleDeviceHeartbeatOutput, error) {
	output := &device.HandleDeviceHeartbeatOutput{
		DeviceID:      input.DeviceID,
		HeartbeatData: make(map[string]interface{}),
	}

	// 验证设备ID
	if input.DeviceID == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "设备ID不能为空")
	}

	// 更新心跳时间
	input.HeartbeatData["last_heartbeat"] = gtime.Now().Format("2006-01-02 15:04:05")
	input.HeartbeatData["device_id"] = input.DeviceID

	// 更新设备状态
	if status, ok := input.HeartbeatData["status"].(string); ok {
		input.HeartbeatData["status"] = status
	} else {
		input.HeartbeatData["status"] = "online"
	}

	// 计算健康度评分
	healthInput := &device.CalculateDeviceHealthScoreInput{
		DeviceData: input.HeartbeatData,
	}
	healthOutput, err := s.CalculateDeviceHealthScore(ctx, healthInput)
	if err != nil {
		return nil, err
	}

	// 自动确定设备状态
	statusInput := &device.DetermineDeviceStatusInput{
		DeviceData: input.HeartbeatData,
	}
	statusOutput, err := s.DetermineDeviceStatus(ctx, statusInput)
	if err != nil {
		return nil, err
	}

	// 复制处理后的数据
	for key, value := range input.HeartbeatData {
		output.HeartbeatData[key] = value
	}

	output.HealthScore = healthOutput.HealthScore
	output.Status = statusOutput.Status
	output.Message = "设备心跳处理成功"

	return output, nil
}

// HandleDeviceDeactivation 处理设备停用
func (s *sDevice) HandleDeviceDeactivation(ctx context.Context, input *device.HandleDeviceDeactivationInput) (*device.HandleDeviceDeactivationOutput, error) {
	output := &device.HandleDeviceDeactivationOutput{
		DeviceID: input.DeviceID,
	}

	// 验证设备ID
	if input.DeviceID == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "设备ID不能为空")
	}

	// 验证停用原因
	if input.Reason == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "停用原因不能为空")
	}

	// 这里应该更新设备状态为停用
	// 目前只是验证参数
	output.Message = "设备停用处理成功"
	output.IsSuccess = true

	return output, nil
}

// CalculateTaskAssignmentScore 计算任务分配评分
func (s *sDevice) CalculateTaskAssignmentScore(ctx context.Context, input *device.CalculateTaskAssignmentScoreInput) (*device.CalculateTaskAssignmentScoreOutput, error) {
	output := &device.CalculateTaskAssignmentScoreOutput{
		Score:      0.0,
		Components: make(map[string]interface{}),
	}

	score := 0.0

	// 1. 设备健康度评分 (权重: 30%)
	healthInput := &device.CalculateDeviceHealthScoreInput{
		DeviceData: input.DeviceData,
	}
	healthOutput, err := s.CalculateDeviceHealthScore(ctx, healthInput)
	if err != nil {
		return nil, err
	}
	score += healthOutput.HealthScore * 0.30
	output.Components["health"] = map[string]interface{}{
		"score":  healthOutput.HealthScore,
		"weight": 0.30,
	}

	// 2. 设备负载评分 (权重: 25%)
	loadScore := s.calculateDeviceLoadScore(input.DeviceData)
	score += loadScore * 0.25
	output.Components["load"] = map[string]interface{}{
		"score":  loadScore,
		"weight": 0.25,
	}

	// 3. 设备能力匹配评分 (权重: 20%)
	capabilityScore := s.calculateCapabilityMatchScore(input.DeviceData, input.TaskData)
	score += capabilityScore * 0.20
	output.Components["capability"] = map[string]interface{}{
		"score":  capabilityScore,
		"weight": 0.20,
	}

	// 4. 设备可用性评分 (权重: 15%)
	availabilityScore := s.calculateAvailabilityScore(input.DeviceData)
	score += availabilityScore * 0.15
	output.Components["availability"] = map[string]interface{}{
		"score":  availabilityScore,
		"weight": 0.15,
	}

	// 5. 设备性能评分 (权重: 10%)
	performanceScore := s.calculatePerformanceScore(input.DeviceData)
	score += performanceScore * 0.10
	output.Components["performance"] = map[string]interface{}{
		"score":  performanceScore,
		"weight": 0.10,
	}

	output.Score = math.Round(score*100) / 100

	// 生成分配建议
	if output.Score >= 80 {
		output.Recommendation = "强烈推荐分配"
	} else if output.Score >= 60 {
		output.Recommendation = "推荐分配"
	} else if output.Score >= 40 {
		output.Recommendation = "谨慎分配"
	} else {
		output.Recommendation = "不推荐分配"
	}

	return output, nil
}

// calculateDeviceLoadScore 计算设备负载评分
func (s *sDevice) calculateDeviceLoadScore(deviceData map[string]interface{}) float64 {
	score := 100.0

	// CPU负载影响
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			score -= 40
		} else if cpuUsage > 70 {
			score -= 20
		} else if cpuUsage > 50 {
			score -= 10
		}
	}

	// 内存负载影响
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage > 90 {
			score -= 30
		} else if memUsage > 80 {
			score -= 15
		} else if memUsage > 60 {
			score -= 5
		}
	}

	// 任务负载影响
	if currentTasks, ok := deviceData["running_task_count"].(int); ok {
		if maxTasks, ok := deviceData["max_concurrent_tasks"].(int); ok && maxTasks > 0 {
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

// calculateCapabilityMatchScore 计算能力匹配评分
func (s *sDevice) calculateCapabilityMatchScore(deviceData, taskData map[string]interface{}) float64 {
	score := 50.0

	// 检查设备类型兼容性
	if taskType, ok := taskData["task_type"].(string); ok {
		if deviceType, ok := deviceData["device_type"].(string); ok {
			if s.isTypeCompatible(taskType, deviceType) {
				score += 30
			}
		}
	}

	// 检查资源需求兼容性
	if s.checkResourceCompatibility(deviceData, taskData) {
		score += 20
	}

	return math.Min(score, 100)
}

// calculateAvailabilityScore 计算可用性评分
func (s *sDevice) calculateAvailabilityScore(deviceData map[string]interface{}) float64 {
	score := 100.0

	// 检查设备状态
	if status, ok := deviceData["status"].(string); ok {
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
	if maintenanceMode, ok := deviceData["maintenance_mode"].(bool); ok && maintenanceMode {
		score *= 0.5
	}

	return score
}

// calculatePerformanceScore 计算性能评分
func (s *sDevice) calculatePerformanceScore(deviceData map[string]interface{}) float64 {
	score := 100.0

	// 基于错误次数调整评分
	if errorCount, ok := deviceData["error_count"].(int); ok {
		if errorCount > 10 {
			score -= 40
		} else if errorCount > 5 {
			score -= 20
		} else if errorCount > 2 {
			score -= 10
		}
	}

	// 基于网络延迟调整评分
	if latency, ok := deviceData["network_latency"].(float64); ok {
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

// isTypeCompatible 检查类型兼容性
func (s *sDevice) isTypeCompatible(taskType, deviceType string) bool {
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

	if compatibleTypes, exists := compatibilityMap[taskType]; exists {
		for _, compatibleType := range compatibleTypes {
			if compatibleType == deviceType {
				return true
			}
		}
	}

	return false
}

// checkResourceCompatibility 检查资源兼容性
func (s *sDevice) checkResourceCompatibility(deviceData, taskData map[string]interface{}) bool {
	// 检查是否有足够的资源余量
	if taskCPU, ok := taskData["required_cpu"].(float64); ok {
		if deviceCPU, ok := deviceData["cpu_usage"].(float64); ok {
			if deviceCPU+taskCPU > 90 { // 留10%余量
				return false
			}
		}
	}

	if taskMem, ok := taskData["required_memory"].(float64); ok {
		if deviceMem, ok := deviceData["memory_usage"].(float64); ok {
			if deviceMem+taskMem > 85 { // 留15%余量
				return false
			}
		}
	}

	return true
}

// generateDeviceId 生成设备ID
func (s *sDevice) generateDeviceId() string {
	// 这里应该使用更复杂的ID生成算法
	// 目前使用时间戳作为简单示例
	return fmt.Sprintf("device_%d", time.Now().UnixNano())
}

// 辅助方法
func (s *sDevice) getStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

func (s *sDevice) getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(int); ok {
		return value
	}
	return defaultValue
}

func (s *sDevice) getFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key].(float64); ok {
		return value
	}
	return defaultValue
}
