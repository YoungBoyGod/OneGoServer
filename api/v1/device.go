package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DeviceRegisterReq 设备注册请求
type DeviceRegisterReq struct {
	DeviceID    string `json:"device_id" binding:"required,max=100"`
	Name        string `json:"name" binding:"required,max=100"`
	Type        string `json:"type" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	IP          string `json:"ip" binding:"omitempty,ip"`
	MAC         string `json:"mac" binding:"omitempty,max=17"`
	OS          string `json:"os" binding:"omitempty,max=100"`
	Version     string `json:"version" binding:"omitempty,max=50"`
}

// DeviceUpdateReq 设备更新请求
type DeviceUpdateReq struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Type        string `json:"type" binding:"omitempty,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Version     string `json:"version" binding:"omitempty,max=50"`
}

// DeviceHeartbeatReq 设备心跳请求
type DeviceHeartbeatReq struct {
	IP       string                 `json:"ip" binding:"omitempty,ip"`
	Version  string                 `json:"version" binding:"omitempty"`
	OS       string                 `json:"os" binding:"omitempty"`
	Metadata map[string]interface{} `json:"metadata" binding:"omitempty"`
}

// DeviceStatusUpdateReq 设备状态更新请求
type DeviceStatusUpdateReq struct {
	Status int `json:"status" binding:"required,min=0,max=3"`
}

// DeviceResp 设备响应
type DeviceResp struct {
	ID          uint       `json:"id"`
	DeviceID    string     `json:"device_id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	IP          string     `json:"ip"`
	MAC         string     `json:"mac"`
	OS          string     `json:"os"`
	Version     string     `json:"version"`
	UserID      uint       `json:"user_id"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DeviceListReq 设备列表请求
type DeviceListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID   *uint  `form:"user_id"`
	Status   *int   `form:"status"`
	Type     string `form:"type"`
	Name     string `form:"name"`
}

// DeviceStatsResp 设备统计响应
type DeviceStatsResp struct {
	Total       int64 `json:"total"`
	Online      int64 `json:"online"`
	Offline     int64 `json:"offline"`
	Maintenance int64 `json:"maintenance"`
	Error       int64 `json:"error"`
}

// 设备接口
func Device(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Device!"})
}
