package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 任务接口
func Task(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Task!"})
}

// TaskCreateReq 任务创建请求
type TaskCreateReq struct {
	Name        string                 `json:"name" binding:"required,max=200"`
	Description string                 `json:"description" binding:"omitempty,max=1000"`
	Type        int                    `json:"type" binding:"required,min=1,max=4"`
	Priority    int                    `json:"priority" binding:"omitempty,min=1,max=4"`
	Command     string                 `json:"command" binding:"required,max=2000"`
	Parameters  map[string]interface{} `json:"parameters" binding:"omitempty"`
	DeviceID    uint                   `json:"device_id" binding:"required"`
	MaxRetries  int                    `json:"max_retries" binding:"omitempty,min=0,max=10"`
	Timeout     int                    `json:"timeout" binding:"omitempty,min=1"`
	ScheduledAt *time.Time             `json:"scheduled_at" binding:"omitempty"`
}

// TaskUpdateReq 任务更新请求
type TaskUpdateReq struct {
	Name        string                 `json:"name" binding:"omitempty,max=200"`
	Description string                 `json:"description" binding:"omitempty,max=1000"`
	Priority    int                    `json:"priority" binding:"omitempty,min=1,max=4"`
	Parameters  map[string]interface{} `json:"parameters" binding:"omitempty"`
	MaxRetries  int                    `json:"max_retries" binding:"omitempty,min=0,max=10"`
	Timeout     int                    `json:"timeout" binding:"omitempty,min=1"`
	ScheduledAt *time.Time             `json:"scheduled_at" binding:"omitempty"`
}

// TaskResp 任务响应
type TaskResp struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        int        `json:"type"`
	Status      int        `json:"status"`
	Priority    int        `json:"priority"`
	Command     string     `json:"command"`
	Parameters  string     `json:"parameters"`
	Result      string     `json:"result"`
	ErrorMsg    string     `json:"error_msg"`
	RetryCount  int        `json:"retry_count"`
	MaxRetries  int        `json:"max_retries"`
	Timeout     int        `json:"timeout"`
	UserID      uint       `json:"user_id"`
	DeviceID    uint       `json:"device_id"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskListReq 任务列表请求
type TaskListReq struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	UserID   *uint  `form:"user_id"`
	DeviceID *uint  `form:"device_id"`
	Status   *int   `form:"status"`
	Type     *int   `form:"type"`
	Priority *int   `form:"priority"`
	Name     string `form:"name"`
}

// TaskLogResp 任务日志响应
type TaskLogResp struct {
	TaskID    uint      `json:"task_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}
