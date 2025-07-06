// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package queue

import (
	"context"

	"OneGfServer/api/queue/v1"
)

type IQueueV1 interface {
	CreateQueue(ctx context.Context, req *v1.CreateQueueReq) (res *v1.CreateQueueRes, err error)
	GetQueueList(ctx context.Context, req *v1.GetQueueListReq) (res *v1.GetQueueListRes, err error)
	GetQueueDetail(ctx context.Context, req *v1.GetQueueDetailReq) (res *v1.GetQueueDetailRes, err error)
	UpdateQueue(ctx context.Context, req *v1.UpdateQueueReq) (res *v1.UpdateQueueRes, err error)
	DeleteQueue(ctx context.Context, req *v1.DeleteQueueReq) (res *v1.DeleteQueueRes, err error)
	GetQueueConfig(ctx context.Context, req *v1.GetQueueConfigReq) (res *v1.GetQueueConfigRes, err error)
	UpdateQueueConfig(ctx context.Context, req *v1.UpdateQueueConfigReq) (res *v1.UpdateQueueConfigRes, err error)
	GetQueueConfigHistory(ctx context.Context, req *v1.GetQueueConfigHistoryReq) (res *v1.GetQueueConfigHistoryRes, err error)
	StartQueue(ctx context.Context, req *v1.StartQueueReq) (res *v1.StartQueueRes, err error)
	PauseQueue(ctx context.Context, req *v1.PauseQueueReq) (res *v1.PauseQueueRes, err error)
	ResumeQueue(ctx context.Context, req *v1.ResumeQueueReq) (res *v1.ResumeQueueRes, err error)
	StopQueue(ctx context.Context, req *v1.StopQueueReq) (res *v1.StopQueueRes, err error)
	ClearQueue(ctx context.Context, req *v1.ClearQueueReq) (res *v1.ClearQueueRes, err error)
	ResetQueue(ctx context.Context, req *v1.ResetQueueReq) (res *v1.ResetQueueRes, err error)
	GetQueueStatus(ctx context.Context, req *v1.GetQueueStatusReq) (res *v1.GetQueueStatusRes, err error)
	GetQueueMetrics(ctx context.Context, req *v1.GetQueueMetricsReq) (res *v1.GetQueueMetricsRes, err error)
	GetQueueStatistics(ctx context.Context, req *v1.GetQueueStatisticsReq) (res *v1.GetQueueStatisticsRes, err error)
	GetQueuePerformanceReport(ctx context.Context, req *v1.GetQueuePerformanceReportReq) (res *v1.GetQueuePerformanceReportRes, err error)
	EnqueueTask(ctx context.Context, req *v1.EnqueueTaskReq) (res *v1.EnqueueTaskRes, err error)
	DequeueTask(ctx context.Context, req *v1.DequeueTaskReq) (res *v1.DequeueTaskRes, err error)
	GetQueueTasks(ctx context.Context, req *v1.GetQueueTasksReq) (res *v1.GetQueueTasksRes, err error)
	ReorderQueueTasks(ctx context.Context, req *v1.ReorderQueueTasksReq) (res *v1.ReorderQueueTasksRes, err error)
	RemoveQueueTask(ctx context.Context, req *v1.RemoveQueueTaskReq) (res *v1.RemoveQueueTaskRes, err error)
	UpdateTaskPriority(ctx context.Context, req *v1.UpdateTaskPriorityReq) (res *v1.UpdateTaskPriorityRes, err error)
}
