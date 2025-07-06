package device

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	device "OneGfServer/internal/model/device"
)

// ===============================
// 设备状态管理相关业务逻辑
// ===============================

// ValidateDeviceStatus 验证设备状态转换是否合法
func (s *sDevice) ValidateDeviceStatus(ctx context.Context, input *device.ValidateDeviceStatusInput) (*device.ValidateDeviceStatusOutput, error) {
	// 定义合法的状态转换规则
	validTransitions := map[string][]string{
		"offline":     {"online", "maintenance", "deleted"},
		"online":      {"offline", "busy", "maintenance", "error"},
		"busy":        {"online", "offline", "error"},
		"maintenance": {"online", "offline"},
		"error":       {"offline", "maintenance"},
		"deleted":     {}, // 删除状态不能转换到其他状态
	}

	allowedStatuses, exists := validTransitions[input.CurrentStatus]
	if !exists {
		return &device.ValidateDeviceStatusOutput{
			IsValid: false,
			Message: fmt.Sprintf("无效的当前状态: %s", input.CurrentStatus),
		}, gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("无效的当前状态: %s", input.CurrentStatus))
	}

	// 检查目标状态是否在允许的转换列表中
	for _, status := range allowedStatuses {
		if status == input.TargetStatus {
			return &device.ValidateDeviceStatusOutput{
				IsValid: true,
				Message: "状态转换有效",
			}, nil
		}
	}

	return &device.ValidateDeviceStatusOutput{
			IsValid: false,
			Message: fmt.Sprintf("不允许从状态 %s 转换到 %s", input.CurrentStatus, input.TargetStatus),
		}, gerror.NewCode(gcode.CodeInvalidParameter,
			fmt.Sprintf("不允许从状态 %s 转换到 %s", input.CurrentStatus, input.TargetStatus))
}

// DetermineDeviceStatus 根据设备数据自动确定设备状态
func (s *sDevice) DetermineDeviceStatus(ctx context.Context, input *device.DetermineDeviceStatusInput) (*device.DetermineDeviceStatusOutput, error) {
	healthScore, err := s.CalculateDeviceHealthScore(ctx, &device.CalculateDeviceHealthScoreInput{
		DeviceData: input.DeviceData,
	})
	if err != nil {
		return nil, err
	}

	status := ""
	reason := ""

	// 检查是否离线
	if lastHeartbeat, ok := input.DeviceData["last_heartbeat"].(*gtime.Time); ok && lastHeartbeat != nil {
		if time.Since(lastHeartbeat.Time) > 15*time.Minute {
			status = "offline"
			reason = "设备心跳超时"
		}
	}

	// 检查是否处于维护状态
	if maintenanceMode, ok := input.DeviceData["maintenance_mode"].(bool); ok && maintenanceMode {
		status = "maintenance"
		reason = "设备处于维护模式"
	}

	// 根据健康度评分确定状态
	if status == "" {
		if healthScore.HealthScore < 30 {
			status = "error"
			reason = "设备健康度评分过低"
		} else if healthScore.HealthScore < 70 {
			// 检查是否正在执行任务
			if taskCount, ok := input.DeviceData["running_task_count"].(int); ok && taskCount > 0 {
				status = "busy"
				reason = "设备正在执行任务"
			} else {
				status = "online"
				reason = "设备在线但健康度一般"
			}
		} else {
			// 检查是否正在执行任务
			if taskCount, ok := input.DeviceData["running_task_count"].(int); ok && taskCount > 0 {
				status = "busy"
				reason = "设备正在执行任务"
			} else {
				status = "online"
				reason = "设备在线且健康度良好"
			}
		}
	}

	return &device.DetermineDeviceStatusOutput{
		Status:      status,
		Reason:      reason,
		HealthScore: healthScore.HealthScore,
	}, nil
}

// CanAcceptNewTask 判断设备是否可以接受新任务
func (s *sDevice) CanAcceptNewTask(ctx context.Context, input *device.CanAcceptNewTaskInput) (*device.CanAcceptNewTaskOutput, error) {
	// 检查设备状态
	status, ok := input.DeviceData["status"].(string)
	if !ok || status != "online" {
		return &device.CanAcceptNewTaskOutput{
			CanAccept: false,
			Reason:    "设备状态不是在线状态",
			Score:     0,
		}, nil
	}

	// 检查健康度
	healthScore, err := s.CalculateDeviceHealthScore(ctx, &device.CalculateDeviceHealthScoreInput{
		DeviceData: input.DeviceData,
	})
	if err != nil {
		return nil, err
	}

	if healthScore.HealthScore < 50 {
		return &device.CanAcceptNewTaskOutput{
			CanAccept: false,
			Reason:    "设备健康度评分过低",
			Score:     healthScore.HealthScore,
		}, nil
	}

	// 检查任务负载
	score := healthScore.HealthScore
	if currentTasks, ok := input.DeviceData["running_task_count"].(int); ok {
		if maxTasks, ok := input.DeviceData["max_concurrent_tasks"].(int); ok {
			if currentTasks >= maxTasks {
				return &device.CanAcceptNewTaskOutput{
					CanAccept: false,
					Reason:    "设备任务负载已满",
					Score:     score * 0.5, // 负载满时评分减半
				}, nil
			}
			// 根据负载情况调整评分
			loadRatio := float64(currentTasks) / float64(maxTasks)
			score = score * (1 - loadRatio*0.3) // 负载越高评分越低
		} else {
			// 默认最大并发任务数为5
			if currentTasks >= 5 {
				return &device.CanAcceptNewTaskOutput{
					CanAccept: false,
					Reason:    "设备任务负载已满",
					Score:     score * 0.5,
				}, nil
			}
			loadRatio := float64(currentTasks) / 5.0
			score = score * (1 - loadRatio*0.3)
		}
	}

	return &device.CanAcceptNewTaskOutput{
		CanAccept: true,
		Reason:    "设备可以接受新任务",
		Score:     score,
	}, nil
}
