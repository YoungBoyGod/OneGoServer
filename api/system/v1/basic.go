package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 基础健康检查相关API
// ===============================

// HealthCheckReq 健康检查请求
type HealthCheckReq struct {
	g.Meta `path:"/health" method:"get" tags:"系统管理" summary:"系统健康检查"`
}

// HealthCheckRes 健康检查响应
type HealthCheckRes struct {
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
	Checks    map[string]interface{} `json:"checks"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime"`
}
