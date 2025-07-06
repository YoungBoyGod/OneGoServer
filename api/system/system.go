// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package system

import (
	"context"

	"OneGfServer/api/system/v1"
)

type ISystemV1 interface {
	GetSystemAlerts(ctx context.Context, req *v1.GetSystemAlertsReq) (res *v1.GetSystemAlertsRes, err error)
	AcknowledgeAlert(ctx context.Context, req *v1.AcknowledgeAlertReq) (res *v1.AcknowledgeAlertRes, err error)
	ResolveAlert(ctx context.Context, req *v1.ResolveAlertReq) (res *v1.ResolveAlertRes, err error)
	HealthCheck(ctx context.Context, req *v1.HealthCheckReq) (res *v1.HealthCheckRes, err error)
	GetSystemLogs(ctx context.Context, req *v1.GetSystemLogsReq) (res *v1.GetSystemLogsRes, err error)
	GetSystemMetrics(ctx context.Context, req *v1.GetSystemMetricsReq) (res *v1.GetSystemMetricsRes, err error)
	GetSystemPerformance(ctx context.Context, req *v1.GetSystemPerformanceReq) (res *v1.GetSystemPerformanceRes, err error)
	GetSystemStatus(ctx context.Context, req *v1.GetSystemStatusReq) (res *v1.GetSystemStatusRes, err error)
}
