package v1

import (
	"time"
)

// AlertCreateReq 告警创建请求
type AlertCreateReq struct {
	Title    string                 `json:"title" binding:"required,max=200"`
	Message  string                 `json:"message" binding:"required,max=1000"`
	Type     int                    `json:"type" binding:"required,min=1,max=5"`
	Level    int                    `json:"level" binding:"required,min=1,max=4"`
	Source   string                 `json:"source" binding:"omitempty,max=100"`
	DeviceID uint                   `json:"device_id" binding:"required"`
	Metadata map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// AlertResp 告警响应
type AlertResp struct {
	ID         uint       `json:"id"`
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	Type       int        `json:"type"`
	Level      int        `json:"level"`
	Status     int        `json:"status"`
	Source     string     `json:"source"`
	Metadata   string     `json:"metadata"`
	DeviceID   uint       `json:"device_id"`
	UserID     uint       `json:"user_id"`
	AckedBy    *uint      `json:"acked_by"`
	AckedAt    *time.Time `json:"acked_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// AlertListReq 告警列表请求
type AlertListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	DeviceID *uint  `form:"device_id"`
	UserID   *uint  `form:"user_id"`
	Type     *int   `form:"type"`
	Level    *int   `form:"level"`
	Status   *int   `form:"status"`
	Source   string `form:"source"`
}

// AlertAckReq 告警确认请求
type AlertAckReq struct {
	UserID uint `json:"user_id" binding:"required"`
}

// AlertStatsResp 告警统计响应
type AlertStatsResp struct {
	Total        int64 `json:"total"`
	Open         int64 `json:"open"`
	Acknowledged int64 `json:"acknowledged"`
	Resolved     int64 `json:"resolved"`
	Closed       int64 `json:"closed"`
	Info         int64 `json:"info"`
	Warning      int64 `json:"warning"`
	Error        int64 `json:"error"`
	Critical     int64 `json:"critical"`
}
