package controller

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/device"
	"github.com/YoungBoyGod/OneGoServer/internal/service"
	"github.com/gin-gonic/gin"
)

/*
DeviceController 设备控制器

已实现的API端点列表：
======================

基础CRUD操作：
- POST   /api/v1/devices              注册设备
- GET    /api/v1/devices              获取设备列表(支持分页：page, size, status, type)
- GET    /api/v1/devices/{id}         获取设备详情
- PUT    /api/v1/devices/{id}         更新设备
- DELETE /api/v1/devices/{id}         删除设备

设备状态管理：
- GET    /api/v1/devices/{id}/status     获取设备状态
- POST   /api/v1/devices/online         设备上线
- POST   /api/v1/devices/offline        设备下线

设备控制操作：
- POST   /api/v1/devices/{id}/command      发送设备命令
- GET    /api/v1/devices/{id}/heartbeat   获取设备心跳
- POST   /api/v1/devices/{id}/heartbeat   更新设备心跳

设备信息查询：
- GET    /api/v1/devices/query          查询设备
- GET    /api/v1/devices/{id}/logs      获取设备日志
- GET    /api/v1/devices/{id}/stats     获取设备统计

总计：14个API端点
状态：所有端点已创建占位函数，待实现具体业务逻辑

待扩展功能：
- 设备认证机制
- 设备告警管理
- 设备固件升级
- 设备分组管理
*/

// DeviceController 设备控制器
type DeviceController struct {
	deviceService service.DeviceService
}

// NewDeviceController 创建设备控制器实例
func NewDeviceController(deviceService service.DeviceService) *DeviceController {
	return &DeviceController{deviceService: deviceService}
}

// RegisterDevice 注册设备
// @Summary 注册新设备
// @Description 注册一个新的设备到系统
// @Tags devices
// @Accept json
// @Produce json
// @Param device body object true "设备信息"
// @Success 201 {object} object "注册成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices [post]
func (dc *DeviceController) RegisterDevice(c *gin.Context) {
	var req device.DeviceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}

	dev, err := dc.deviceService.RegisterDevice(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "data": dev})
}

// DeviceOnline 设备上线
// @Summary 设备上线
// @Description 设备上线
// @Tags devices
// @Accept json
// @Produce json
// @Param device body object true "设备信息"
// @Success 200 {object} object "设备上线成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/online [post]
func (dc *DeviceController) DeviceOnline(c *gin.Context) {
	var body struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}
	if err := dc.deviceService.DeviceOnline(c.Request.Context(), body.DeviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "设备上线失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设备已上线", "device_id": body.DeviceID})
}

// DeviceOffline 设备下线
// @Summary 设备下线
// @Description 设备下线
// @Tags devices
// @Accept json
// @Produce json
// @Param device body object true "设备信息"
// @Success 200 {object} object "设备下线成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/offline [post]
func (dc *DeviceController) DeviceOffline(c *gin.Context) {
	var body struct {
		DeviceID string `json:"device_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
		return
	}
	if err := dc.deviceService.DeviceOffline(c.Request.Context(), body.DeviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "设备下线失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设备已下线", "device_id": body.DeviceID})
}

// GetDeviceHeartbeat 获取设备心跳
// @Summary 获取设备心跳
// @Description 获取设备心跳记录
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "心跳信息"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/heartbeat [get]
func (dc *DeviceController) GetDeviceHeartbeat(c *gin.Context) {
	deviceID := c.Param("id")
	hb, err := dc.deviceService.GetDeviceHeartbeat(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": hb})
}

// UpdateDeviceHeartbeat 更新设备心跳
// @Summary 更新设备心跳
// @Description 设备上报心跳信息
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "心跳更新成功"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/heartbeat [post]
func (dc *DeviceController) UpdateDeviceHeartbeat(c *gin.Context) {
	deviceID := c.Param("id")
	var req device.DeviceHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := dc.deviceService.ProcessHeartbeat(c.Request.Context(), deviceID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "心跳已更新"})
}

// GetDeviceStatus 获取设备状态
// @Summary 获取设备状态
// @Description 获取设备当前运行状态
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "设备状态"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/status [get]
func (dc *DeviceController) GetDeviceStatus(c *gin.Context) {
	deviceID := c.Param("id")
	status, err := dc.deviceService.GetDeviceStatus(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": status})
}

// GetDeviceLogs 获取设备日志
// @Summary 获取设备日志
// @Description 获取设备日志
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "设备日志"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/logs [get]
func (dc *DeviceController) GetDeviceLogs(c *gin.Context) {
	deviceID := c.Param("id")
	logs, err := dc.deviceService.GetDeviceLogs(c.Request.Context(), deviceID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": logs})
}

// GetDeviceStats 获取设备统计
// @Summary 获取设备统计
// @Description 获取设备统计信息
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "设备统计"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/stats [get]
func (dc *DeviceController) GetDeviceStats(c *gin.Context) {
	filter := &device.DeviceFilter{}
	stats, err := dc.deviceService.GetDeviceStats(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": stats})
}

// GetDevices 获取设备列表
// @Summary 获取设备列表
// @Description 获取设备列表，支持分页和筛选
// @Tags devices
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param status query string false "设备状态"
// @Param type query string false "设备类型"
// @Success 200 {object} object "设备列表"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices [get]
func (dc *DeviceController) GetDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	filter := &device.DeviceFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = []string{status}
	}
	if dType := c.Query("type"); dType != "" {
		filter.Type = []string{dType}
	}

	resp, err := dc.deviceService.ListDevices(c.Request.Context(), filter, nil, &device.PaginationOption{Page: page, Size: size})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": resp})
}

// GetDevice 获取单个设备详情
// @Summary 获取设备详情
// @Description 根据ID获取设备详细信息
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "设备详情"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id} [get]
func (dc *DeviceController) GetDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效ID"})
		return
	}
	dev, err := dc.deviceService.GetDevice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取成功", "data": dev})
}

// UpdateDevice 更新设备
// @Summary 更新设备
// @Description 更新设备信息
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Param device body object true "设备更新信息"
// @Success 200 {object} object "更新成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id} [put]
func (dc *DeviceController) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效ID"})
		return
	}
	var req device.DeviceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dev, err := dc.deviceService.UpdateDevice(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": dev})
}

// DeleteDevice 删除设备
// @Summary 删除设备
// @Description 根据ID删除设备
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Success 200 {object} object "删除成功"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id} [delete]
func (dc *DeviceController) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效ID"})
		return
	}
	if err := dc.deviceService.DeleteDevice(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// SendCommand 发送设备命令
// @Summary 发送设备命令
// @Description 向设备发送控制命令
// @Tags devices
// @Accept json
// @Produce json
// @Param id path string true "设备ID"
// @Param command body object true "命令信息"
// @Success 200 {object} object "命令发送成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/{id}/command [post]
func (dc *DeviceController) SendCommand(c *gin.Context) {
	deviceID := c.Param("id")
	var req device.DeviceCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	cmd, err := dc.deviceService.SendCommand(c.Request.Context(), deviceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "命令已发送", "data": cmd})
}

// QueryDevice 查询设备
// @Summary 查询设备
// @Description 根据条件查询设备
// @Tags devices
// @Accept json
// @Produce json
// @Success 200 {object} object "设备详情"
// @Failure 404 {object} object "设备不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/devices/query [get]
func (dc *DeviceController) QueryDevice(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少keyword"})
		return
	}
	devices, err := dc.deviceService.QueryDevices(c.Request.Context(), keyword, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "data": devices})
}
