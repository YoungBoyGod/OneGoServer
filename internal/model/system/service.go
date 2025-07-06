package system

import "context"

// ===============================
// System模块服务接口定义
// ===============================

// ISystem 系统服务接口
type ISystem interface {
	// 健康检查
	HealthCheck(ctx context.Context, input *HealthCheckInput) (*HealthCheckOutput, error)

	// 系统指标
	GetSystemMetrics(ctx context.Context, input *GetSystemMetricsInput) (*GetSystemMetricsOutput, error)

	// 系统状态
	GetSystemStatus(ctx context.Context, input *GetSystemStatusInput) (*GetSystemStatusOutput, error)

	// 系统日志
	GetSystemLogs(ctx context.Context, input *GetSystemLogsInput) (*GetSystemLogsOutput, error)

	// 系统告警
	GetSystemAlerts(ctx context.Context, input *GetSystemAlertsInput) (*GetSystemAlertsOutput, error)
	AcknowledgeAlert(ctx context.Context, input *AcknowledgeAlertInput) (*AcknowledgeAlertOutput, error)
	ResolveAlert(ctx context.Context, input *ResolveAlertInput) (*ResolveAlertOutput, error)

	// 系统性能
	GetSystemPerformance(ctx context.Context, input *GetSystemPerformanceInput) (*GetSystemPerformanceOutput, error)
}
