package device

import (
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
func (s *sDevice) HandleDeviceRegistration(ctx context.Context, deviceData map[string]interface{}) (map[string]interface{}, error) {
	// 验证设备注册数据
	if err := s.ValidateDeviceRegistration(ctx, deviceData); err != nil {
		return nil, err
	}

	// 生成设备ID（如果未提供）
	if _, ok := deviceData["device_id"]; !ok {
		deviceData["device_id"] = s.generateDeviceId()
	}

	// 设置注册时间
	deviceData["registered_at"] = gtime.Now().Format("2006-01-02 15:04:05")
	deviceData["last_heartbeat"] = gtime.Now().Format("2006-01-02 15:04:05")
	deviceData["status"] = "online"

	// 初始化设备指标
	deviceData["cpu_usage"] = 0.0
	deviceData["memory_usage"] = 0.0
	deviceData["disk_usage"] = 0.0
	deviceData["network_latency"] = 0.0
	deviceData["running_task_count"] = 0
	deviceData["error_count"] = 0

	// 这里应该保存到数据库
	// 目前返回处理后的数据
	return deviceData, nil
}

// HandleDeviceHeartbeat 处理设备心跳
func (s *sDevice) HandleDeviceHeartbeat(ctx context.Context, deviceId string, heartbeatData map[string]interface{}) (map[string]interface{}, error) {
	// 验证设备ID
	if deviceId == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "设备ID不能为空")
	}

	// 更新心跳时间
	heartbeatData["last_heartbeat"] = gtime.Now().Format("2006-01-02 15:04:05")
	heartbeatData["device_id"] = deviceId

	// 更新设备状态
	if status, ok := heartbeatData["status"].(string); ok {
		heartbeatData["status"] = status
	} else {
		heartbeatData["status"] = "online"
	}

	// 计算健康度评分
	healthScore := s.CalculateDeviceHealthScore(ctx, heartbeatData)
	heartbeatData["health_score"] = healthScore

	// 自动确定设备状态
	determinedStatus := s.DetermineDeviceStatus(ctx, heartbeatData)
	heartbeatData["determined_status"] = determinedStatus

	// 这里应该更新数据库
	// 目前返回处理后的数据
	return heartbeatData, nil
}

// HandleDeviceDeactivation 处理设备停用
func (s *sDevice) HandleDeviceDeactivation(ctx context.Context, deviceId string, reason string) error {
	// 验证设备ID
	if deviceId == "" {
		return gerror.NewCode(gcode.CodeValidationFailed, "设备ID不能为空")
	}

	// 验证停用原因
	if reason == "" {
		return gerror.NewCode(gcode.CodeValidationFailed, "停用原因不能为空")
	}

	// 这里应该更新设备状态为停用
	// 目前只是验证参数
	return nil
}

// CalculateTaskAssignmentScore 计算任务分配评分
func (s *sDevice) CalculateTaskAssignmentScore(ctx context.Context, deviceData map[string]interface{}, taskData map[string]interface{}) float64 {
	score := 0.0

	// 1. 设备健康度评分 (权重: 30%)
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)
	score += healthScore * 0.30

	// 2. 设备负载评分 (权重: 25%)
	loadScore := s.calculateDeviceLoadScore(deviceData)
	score += loadScore * 0.25

	// 3. 设备能力匹配评分 (权重: 20%)
	capabilityScore := s.calculateCapabilityMatchScore(deviceData, taskData)
	score += capabilityScore * 0.20

	// 4. 设备可用性评分 (权重: 15%)
	availabilityScore := s.calculateAvailabilityScore(deviceData)
	score += availabilityScore * 0.15

	// 5. 设备性能评分 (权重: 10%)
	performanceScore := s.calculatePerformanceScore(deviceData)
	score += performanceScore * 0.10

	return score
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
	score := 50.0 // 默认分数

	// 检查设备类型匹配
	if deviceType, ok := deviceData["device_type"].(string); ok {
		if taskType, ok := taskData["task_type"].(string); ok {
			if s.isTypeCompatible(deviceType, taskType) {
				score += 30
			}
		}
	}

	// 检查资源需求匹配
	if taskCPU, ok := taskData["required_cpu"].(float64); ok {
		if deviceCPU, ok := deviceData["cpu_usage"].(float64); ok {
			if deviceCPU+taskCPU <= 100 {
				score += 20
			} else {
				score -= 20
			}
		}
	}

	return math.Max(0, math.Min(score, 100))
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
func (s *sDevice) isTypeCompatible(deviceType, taskType string) bool {
	// 定义兼容性映射
	compatibilityMap := map[string][]string{
		"server":      {"compute", "storage", "network"},
		"workstation": {"compute", "display", "interactive"},
		"mobile":      {"mobile", "location", "sensor"},
		"iot":         {"sensor", "monitor", "control"},
		"edge":        {"compute", "network", "local"},
	}

	if compatibleTypes, exists := compatibilityMap[deviceType]; exists {
		for _, compatibleType := range compatibleTypes {
			if compatibleType == taskType {
				return true
			}
		}
	}

	return false
}

// generateDeviceId 生成设备ID
func (s *sDevice) generateDeviceId() string {
	// 这里应该使用更复杂的ID生成算法
	// 目前使用时间戳作为简单示例
	return fmt.Sprintf("device_%d", time.Now().UnixNano())
}
