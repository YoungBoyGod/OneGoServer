// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package device

import (
	"context"

	"OneGfServer/api/device/v1"
)

type IDeviceV1 interface {
	RegisterDevice(ctx context.Context, req *v1.RegisterDeviceReq) (res *v1.RegisterDeviceRes, err error)
	GetDeviceList(ctx context.Context, req *v1.GetDeviceListReq) (res *v1.GetDeviceListRes, err error)
	ManageDeviceWhitelist(ctx context.Context, req *v1.ManageDeviceWhitelistReq) (res *v1.ManageDeviceWhitelistRes, err error)
	DeleteDevice(ctx context.Context, req *v1.DeleteDeviceReq) (res *v1.DeleteDeviceRes, err error)
	GetDeviceStatus(ctx context.Context, req *v1.GetDeviceStatusReq) (res *v1.GetDeviceStatusRes, err error)
	UpdateDeviceStatus(ctx context.Context, req *v1.UpdateDeviceStatusReq) (res *v1.UpdateDeviceStatusRes, err error)
	SendDeviceCommand(ctx context.Context, req *v1.SendDeviceCommandReq) (res *v1.SendDeviceCommandRes, err error)
	GetDeviceHeartbeat(ctx context.Context, req *v1.GetDeviceHeartbeatReq) (res *v1.GetDeviceHeartbeatRes, err error)
	UpdateDeviceHeartbeat(ctx context.Context, req *v1.UpdateDeviceHeartbeatReq) (res *v1.UpdateDeviceHeartbeatRes, err error)
	GetDeviceDetail(ctx context.Context, req *v1.GetDeviceDetailReq) (res *v1.GetDeviceDetailRes, err error)
	UpdateDeviceInfo(ctx context.Context, req *v1.UpdateDeviceInfoReq) (res *v1.UpdateDeviceInfoRes, err error)
	GetDeviceTaskList(ctx context.Context, req *v1.GetDeviceTaskListReq) (res *v1.GetDeviceTaskListRes, err error)
	GetDeviceTaskQueue(ctx context.Context, req *v1.GetDeviceTaskQueueReq) (res *v1.GetDeviceTaskQueueRes, err error)
	GetDeviceTaskDetail(ctx context.Context, req *v1.GetDeviceTaskDetailReq) (res *v1.GetDeviceTaskDetailRes, err error)
	GetDeviceAlertList(ctx context.Context, req *v1.GetDeviceAlertListReq) (res *v1.GetDeviceAlertListRes, err error)
	GetDeviceAlertDetail(ctx context.Context, req *v1.GetDeviceAlertDetailReq) (res *v1.GetDeviceAlertDetailRes, err error)
	UpdateDeviceAlert(ctx context.Context, req *v1.UpdateDeviceAlertReq) (res *v1.UpdateDeviceAlertRes, err error)
	GetDeviceLogList(ctx context.Context, req *v1.GetDeviceLogListReq) (res *v1.GetDeviceLogListRes, err error)
	GetDeviceLogDetail(ctx context.Context, req *v1.GetDeviceLogDetailReq) (res *v1.GetDeviceLogDetailRes, err error)
	ClearDeviceLogs(ctx context.Context, req *v1.ClearDeviceLogsReq) (res *v1.ClearDeviceLogsRes, err error)
	GetDeviceLoadStatus(ctx context.Context, req *v1.GetDeviceLoadStatusReq) (res *v1.GetDeviceLoadStatusRes, err error)
	GetDeviceLoadHistory(ctx context.Context, req *v1.GetDeviceLoadHistoryReq) (res *v1.GetDeviceLoadHistoryRes, err error)
	UpdateDeviceLoadConfig(ctx context.Context, req *v1.UpdateDeviceLoadConfigReq) (res *v1.UpdateDeviceLoadConfigRes, err error)
	GetDeviceQueueHistory(ctx context.Context, req *v1.GetDeviceQueueHistoryReq) (res *v1.GetDeviceQueueHistoryRes, err error)
	GetDeviceCommandHistory(ctx context.Context, req *v1.GetDeviceCommandHistoryReq) (res *v1.GetDeviceCommandHistoryRes, err error)
	GetDeviceCommandDetail(ctx context.Context, req *v1.GetDeviceCommandDetailReq) (res *v1.GetDeviceCommandDetailRes, err error)
	BatchOperateDevices(ctx context.Context, req *v1.BatchOperateDevicesReq) (res *v1.BatchOperateDevicesRes, err error)
	GetBatchOperationStatus(ctx context.Context, req *v1.GetBatchOperationStatusReq) (res *v1.GetBatchOperationStatusRes, err error)
	GetDeviceConfig(ctx context.Context, req *v1.GetDeviceConfigReq) (res *v1.GetDeviceConfigRes, err error)
	UpdateDeviceConfig(ctx context.Context, req *v1.UpdateDeviceConfigReq) (res *v1.UpdateDeviceConfigRes, err error)
	GetDeviceConfigHistory(ctx context.Context, req *v1.GetDeviceConfigHistoryReq) (res *v1.GetDeviceConfigHistoryRes, err error)
	GetDeviceStatistics(ctx context.Context, req *v1.GetDeviceStatisticsReq) (res *v1.GetDeviceStatisticsRes, err error)
	GetDevicePerformanceReport(ctx context.Context, req *v1.GetDevicePerformanceReportReq) (res *v1.GetDevicePerformanceReportRes, err error)
	CalculateDeviceLoadScore(ctx context.Context, req *v1.CalculateDeviceLoadScoreReq) (res *v1.CalculateDeviceLoadScoreRes, err error)
	GetDeviceLoadMetrics(ctx context.Context, req *v1.GetDeviceLoadMetricsReq) (res *v1.GetDeviceLoadMetricsRes, err error)
	OptimizeDeviceLoad(ctx context.Context, req *v1.OptimizeDeviceLoadReq) (res *v1.OptimizeDeviceLoadRes, err error)
	GetDeviceLoadHistory(ctx context.Context, req *v1.GetDeviceLoadHistoryReq) (res *v1.GetDeviceLoadHistoryRes, err error)
	SetDeviceLoadThreshold(ctx context.Context, req *v1.SetDeviceLoadThresholdReq) (res *v1.SetDeviceLoadThresholdRes, err error)
	GetDeviceLoadThreshold(ctx context.Context, req *v1.GetDeviceLoadThresholdReq) (res *v1.GetDeviceLoadThresholdRes, err error)
}
