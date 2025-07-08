// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	system "OneGfServer/internal/model/system"
	"context"
)

type (
	ISystem interface {
		// GetSystemAlerts 获取系统告警
		GetSystemAlerts(ctx context.Context, input *system.GetSystemAlertsInput) (*system.GetSystemAlertsOutput, error)
		// AcknowledgeAlert 确认告警
		AcknowledgeAlert(ctx context.Context, input *system.AcknowledgeAlertInput) (*system.AcknowledgeAlertOutput, error)
		// ResolveAlert 解决告警
		ResolveAlert(ctx context.Context, input *system.ResolveAlertInput) (*system.ResolveAlertOutput, error)
		// GetSystemLogs 获取系统日志
		GetSystemLogs(ctx context.Context, input *system.GetSystemLogsInput) (*system.GetSystemLogsOutput, error)
		// GetSystemMetrics 获取系统指标
		GetSystemMetrics(ctx context.Context, input *system.GetSystemMetricsInput) (*system.GetSystemMetricsOutput, error)
		// GetSystemPerformance 获取系统性能
		GetSystemPerformance(ctx context.Context, input *system.GetSystemPerformanceInput) (*system.GetSystemPerformanceOutput, error)
		// GetSystemStatus 获取系统状态
		GetSystemStatus(ctx context.Context, input *system.GetSystemStatusInput) (*system.GetSystemStatusOutput, error)
		// GetSystemInfo 获取系统信息
		GetSystemInfo(ctx context.Context, input *system.GetSystemInfoInput) (*system.GetSystemInfoOutput, error)
	}
)

var (
	localSystem ISystem
)

func System() ISystem {
	if localSystem == nil {
		panic("implement not found for interface ISystem, forgot register?")
	}
	return localSystem
}

func RegisterSystem(i ISystem) {
	localSystem = i
}
