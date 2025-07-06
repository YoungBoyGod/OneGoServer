package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 任务日志管理 API (4个)
// ===============================

// GetTaskLogListReq 获取任务日志列表请求
type GetTaskLogListReq struct {
	g.Meta    `path:"/task/{taskId}/logs" method:"get" tags:"任务日志" summary:"获取任务日志列表"`
	TaskId    string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	Page      int         `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size      int         `json:"size" d:"20" v:"between:1,100#每页数量为1-100"`
	Level     string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	Keyword   string      `json:"keyword,omitempty" v:"max-length:100#关键词最大100字符"`
}

type GetTaskLogListRes struct {
	List  []TaskLogInfo `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
}

// GetTaskLogDetailReq 获取任务日志详情请求
type GetTaskLogDetailReq struct {
	g.Meta `path:"/task/{taskId}/log/{logId}" method:"get" tags:"任务日志" summary:"获取任务日志详情"`
	TaskId string `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	LogId  string `json:"log_id" v:"required|max-length:50#日志ID不能为空"`
}

type GetTaskLogDetailRes struct {
	LogDetail TaskLogDetailInfo `json:"log_detail"`
}

// ClearTaskLogsReq 清空任务日志请求
type ClearTaskLogsReq struct {
	g.Meta     `path:"/task/{taskId}/logs" method:"delete" tags:"任务日志" summary:"清空任务日志"`
	TaskId     string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	BeforeTime *gtime.Time `json:"before_time,omitempty"`
	Level      string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
}

type ClearTaskLogsRes struct {
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// ExportTaskLogsReq 导出任务日志请求
type ExportTaskLogsReq struct {
	g.Meta    `path:"/task/{taskId}/logs/export" method:"post" tags:"任务日志" summary:"导出任务日志"`
	TaskId    string      `json:"task_id" v:"required|max-length:50#任务ID不能为空"`
	StartTime *gtime.Time `json:"start_time,omitempty"`
	EndTime   *gtime.Time `json:"end_time,omitempty"`
	Level     string      `json:"level,omitempty" v:"in:debug,info,warn,error#日志级别无效"`
	Format    string      `json:"format" d:"csv" v:"in:csv,json,txt#导出格式无效"`
}

type ExportTaskLogsRes struct {
	ExportId    string `json:"export_id"`
	DownloadUrl string `json:"download_url"`
	Message     string `json:"message"`
}
