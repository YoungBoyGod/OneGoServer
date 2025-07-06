// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package task

import (
	"context"

	"OneGfServer/api/task/v1"
)

type ITaskV1 interface {
	CreateTask(ctx context.Context, req *v1.CreateTaskReq) (res *v1.CreateTaskRes, err error)
	GetTaskList(ctx context.Context, req *v1.GetTaskListReq) (res *v1.GetTaskListRes, err error)
	GetTaskDetail(ctx context.Context, req *v1.GetTaskDetailReq) (res *v1.GetTaskDetailRes, err error)
	UpdateTask(ctx context.Context, req *v1.UpdateTaskReq) (res *v1.UpdateTaskRes, err error)
	DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error)
	StartTask(ctx context.Context, req *v1.StartTaskReq) (res *v1.StartTaskRes, err error)
	StopTask(ctx context.Context, req *v1.StopTaskReq) (res *v1.StopTaskRes, err error)
	RestartTask(ctx context.Context, req *v1.RestartTaskReq) (res *v1.RestartTaskRes, err error)
	CancelTask(ctx context.Context, req *v1.CancelTaskReq) (res *v1.CancelTaskRes, err error)
	RetryTask(ctx context.Context, req *v1.RetryTaskReq) (res *v1.RetryTaskRes, err error)
	GetTaskStatus(ctx context.Context, req *v1.GetTaskStatusReq) (res *v1.GetTaskStatusRes, err error)
	UpdateTaskStatus(ctx context.Context, req *v1.UpdateTaskStatusReq) (res *v1.UpdateTaskStatusRes, err error)
	GetTaskPriorityList(ctx context.Context, req *v1.GetTaskPriorityListReq) (res *v1.GetTaskPriorityListRes, err error)
	UpdateTaskPriority(ctx context.Context, req *v1.UpdateTaskPriorityReq) (res *v1.UpdateTaskPriorityRes, err error)
	BatchUpdateTaskPriority(ctx context.Context, req *v1.BatchUpdateTaskPriorityReq) (res *v1.BatchUpdateTaskPriorityRes, err error)
	GetTaskLogList(ctx context.Context, req *v1.GetTaskLogListReq) (res *v1.GetTaskLogListRes, err error)
	GetTaskLogDetail(ctx context.Context, req *v1.GetTaskLogDetailReq) (res *v1.GetTaskLogDetailRes, err error)
	ClearTaskLogs(ctx context.Context, req *v1.ClearTaskLogsReq) (res *v1.ClearTaskLogsRes, err error)
	ExportTaskLogs(ctx context.Context, req *v1.ExportTaskLogsReq) (res *v1.ExportTaskLogsRes, err error)
	GetTaskQueueList(ctx context.Context, req *v1.GetTaskQueueListReq) (res *v1.GetTaskQueueListRes, err error)
	GetTaskQueueDetail(ctx context.Context, req *v1.GetTaskQueueDetailReq) (res *v1.GetTaskQueueDetailRes, err error)
	ManageTaskQueue(ctx context.Context, req *v1.ManageTaskQueueReq) (res *v1.ManageTaskQueueRes, err error)
	GetTaskQueueStats(ctx context.Context, req *v1.GetTaskQueueStatsReq) (res *v1.GetTaskQueueStatsRes, err error)
	GetTaskScheduleList(ctx context.Context, req *v1.GetTaskScheduleListReq) (res *v1.GetTaskScheduleListRes, err error)
	CreateTaskSchedule(ctx context.Context, req *v1.CreateTaskScheduleReq) (res *v1.CreateTaskScheduleRes, err error)
	UpdateTaskSchedule(ctx context.Context, req *v1.UpdateTaskScheduleReq) (res *v1.UpdateTaskScheduleRes, err error)
	DeleteTaskSchedule(ctx context.Context, req *v1.DeleteTaskScheduleReq) (res *v1.DeleteTaskScheduleRes, err error)
	GetTaskStatistics(ctx context.Context, req *v1.GetTaskStatisticsReq) (res *v1.GetTaskStatisticsRes, err error)
	GetTaskPerformanceReport(ctx context.Context, req *v1.GetTaskPerformanceReportReq) (res *v1.GetTaskPerformanceReportRes, err error)
	AssignTaskToDevice(ctx context.Context, req *v1.AssignTaskToDeviceReq) (res *v1.AssignTaskToDeviceRes, err error)
	GetTaskAssignmentOptions(ctx context.Context, req *v1.GetTaskAssignmentOptionsReq) (res *v1.GetTaskAssignmentOptionsRes, err error)
	BatchAssignTasks(ctx context.Context, req *v1.BatchAssignTasksReq) (res *v1.BatchAssignTasksRes, err error)
	GetTaskAssignmentHistory(ctx context.Context, req *v1.GetTaskAssignmentHistoryReq) (res *v1.GetTaskAssignmentHistoryRes, err error)
	OptimizeTaskAssignment(ctx context.Context, req *v1.OptimizeTaskAssignmentReq) (res *v1.OptimizeTaskAssignmentRes, err error)
	GetDeviceTaskQueue(ctx context.Context, req *v1.GetDeviceTaskQueueReq) (res *v1.GetDeviceTaskQueueRes, err error)
}
