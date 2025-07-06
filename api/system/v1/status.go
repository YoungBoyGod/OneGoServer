package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 系统状态查询相关API
// ===============================

// GetSystemStatusReq 获取系统状态请求
type GetSystemStatusReq struct {
	g.Meta `path:"/status" method:"get" tags:"系统管理" summary:"获取系统状态"`
}

// GetSystemStatusRes 获取系统状态响应
type GetSystemStatusRes struct {
	Status      string `json:"status"`
	Version     string `json:"version"`
	Uptime      string `json:"uptime"`
	StartTime   string `json:"startTime"`
	Environment string `json:"environment"`
	Database    struct {
		Status         string `json:"status"`
		Connections    int    `json:"connections"`
		MaxConnections int    `json:"maxConnections"`
	} `json:"database"`
	Cache struct {
		Status string `json:"status"`
		Size   int64  `json:"size"`
		Hits   int64  `json:"hits"`
		Misses int64  `json:"misses"`
	} `json:"cache"`
	Devices struct {
		Total   int `json:"total"`
		Online  int `json:"online"`
		Offline int `json:"offline"`
		Error   int `json:"error"`
	} `json:"devices"`
	Tasks struct {
		Total     int `json:"total"`
		Running   int `json:"running"`
		Pending   int `json:"pending"`
		Completed int `json:"completed"`
		Failed    int `json:"failed"`
	} `json:"tasks"`
	Users struct {
		Total  int `json:"total"`
		Active int `json:"active"`
		Online int `json:"online"`
	} `json:"users"`
}
