package controller

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/OneGoServer/internal/service"
	"github.com/gin-gonic/gin"
)

// QueueController 设备队列控制器
type QueueController struct {
	queueService service.QueueService
}

// NewQueueController 构造函数
func NewQueueController(queueService service.QueueService) *QueueController {
	return &QueueController{queueService: queueService}
}

// EnqueueTask 将任务加入队列
func (qc *QueueController) EnqueueTask(c *gin.Context) {
	var req struct {
		DeviceID int64 `json:"device_id" binding:"required"`
		TaskID   int64 `json:"task_id" binding:"required"`
		Priority int   `json:"priority" binding:"required,min=1,max=10"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := qc.queueService.EnqueueTask(c.Request.Context(), req.DeviceID, req.TaskID, req.Priority); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "入队成功"})
}

// DequeueTask 取出下一个任务
func (qc *QueueController) DequeueTask(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效DeviceID"})
		return
	}
	task, err := qc.queueService.DequeueTask(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "出队成功", "data": task})
}

// GetQueueMetrics 获取队列指标
func (qc *QueueController) GetQueueMetrics(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效DeviceID"})
		return
	}
	metrics, err := qc.queueService.GetQueueMetrics(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": metrics})
}

// ReorderQueue 重新排序队列
func (qc *QueueController) ReorderQueue(c *gin.Context) {
	deviceIDStr := c.Param("device_id")
	deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效DeviceID"})
		return
	}
	var body struct {
		Strategy string `json:"strategy" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := qc.queueService.ReorderQueue(c.Request.Context(), deviceID, body.Strategy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "重排成功"})
}
