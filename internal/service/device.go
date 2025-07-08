// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	device "OneGfServer/internal/model/device"
	"context"
)

type (
	IDevice interface {
		// HandleDeviceRegistration 处理设备注册
		HandleDeviceRegistration(ctx context.Context, input *device.HandleDeviceRegistrationInput) (*device.HandleDeviceRegistrationOutput, error)
		// HandleDeviceHeartbeat 处理设备心跳
		HandleDeviceHeartbeat(ctx context.Context, input *device.HandleDeviceHeartbeatInput) (*device.HandleDeviceHeartbeatOutput, error)
		// HandleDeviceDeactivation 处理设备停用
		HandleDeviceDeactivation(ctx context.Context, input *device.HandleDeviceDeactivationInput) (*device.HandleDeviceDeactivationOutput, error)
		// CalculateTaskAssignmentScore 计算任务分配评分
		CalculateTaskAssignmentScore(ctx context.Context, input *device.CalculateTaskAssignmentScoreInput) (*device.CalculateTaskAssignmentScoreOutput, error)
		// CalculateDeviceHealthScore 计算设备健康度评分
		CalculateDeviceHealthScore(ctx context.Context, input *device.CalculateDeviceHealthScoreInput) (*device.CalculateDeviceHealthScoreOutput, error)
		// CheckDeviceHealth 检查设备健康状态
		CheckDeviceHealth(ctx context.Context, input *device.CheckDeviceHealthInput) (*device.CheckDeviceHealthOutput, error)
		// CalculateDeviceLoadScore 计算设备负载评分
		CalculateDeviceLoadScore(ctx context.Context, input *device.CalculateDeviceLoadScoreInput) (*device.CalculateDeviceLoadScoreOutput, error)
		// OptimizeDeviceLoad 优化设备负载
		OptimizeDeviceLoad(ctx context.Context, input *device.OptimizeDeviceLoadInput) (*device.OptimizeDeviceLoadOutput, error)
		// SetDeviceLoadThreshold 设置设备负载阈值
		SetDeviceLoadThreshold(ctx context.Context, input *device.SetDeviceLoadThresholdInput) (*device.SetDeviceLoadThresholdOutput, error)
		// GetDeviceLoadThreshold 获取设备负载阈值
		GetDeviceLoadThreshold(ctx context.Context, input *device.GetDeviceLoadThresholdInput) (*device.GetDeviceLoadThresholdOutput, error)
		// GetDeviceLoadMetrics 获取设备负载指标
		GetDeviceLoadMetrics(ctx context.Context, input *device.GetDeviceLoadMetricsInput) (*device.GetDeviceLoadMetricsOutput, error)
		// AnalyzeDevicePerformance 分析设备性能
		AnalyzeDevicePerformance(ctx context.Context, input *device.AnalyzeDevicePerformanceInput) (*device.AnalyzeDevicePerformanceOutput, error)
		// ValidateDeviceStatus 验证设备状态转换是否合法
		ValidateDeviceStatus(ctx context.Context, input *device.ValidateDeviceStatusInput) (*device.ValidateDeviceStatusOutput, error)
		// DetermineDeviceStatus 根据设备数据自动确定设备状态
		DetermineDeviceStatus(ctx context.Context, input *device.DetermineDeviceStatusInput) (*device.DetermineDeviceStatusOutput, error)
		// CanAcceptNewTask 判断设备是否可以接受新任务
		CanAcceptNewTask(ctx context.Context, input *device.CanAcceptNewTaskInput) (*device.CanAcceptNewTaskOutput, error)
		// ValidateDeviceRegistration 验证设备注册
		ValidateDeviceRegistration(ctx context.Context, input *device.ValidateDeviceRegistrationInput) (*device.ValidateDeviceRegistrationOutput, error)
		// ValidateDeviceConfiguration 验证设备配置
		ValidateDeviceConfiguration(ctx context.Context, input *device.ValidateDeviceConfigurationInput) (*device.ValidateDeviceConfigurationOutput, error)
		// ValidateDeviceData 验证设备数据
		ValidateDeviceData(ctx context.Context, input *device.ValidateDeviceDataInput) (*device.ValidateDeviceDataOutput, error)
		// ProcessDeviceData 处理设备数据
		ProcessDeviceData(ctx context.Context, input *device.ProcessDeviceDataInput) (*device.ProcessDeviceDataOutput, error)
	}
)

var (
	localDevice IDevice
)

func Device() IDevice {
	if localDevice == nil {
		panic("implement not found for interface IDevice, forgot register?")
	}
	return localDevice
}

func RegisterDevice(i IDevice) {
	localDevice = i
}
