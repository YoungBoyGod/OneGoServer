package device

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

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
