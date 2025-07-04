package controller

import (
	"net/http"
	"strconv"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/task"
	"github.com/YoungBoyGod/OneGoServer/internal/service"
	"github.com/gin-gonic/gin"
)

/*
TaskController 任务控制器

已实现的API端点列表：
======================

基础CRUD操作：
- POST   /api/v1/tasks              创建任务
- GET    /api/v1/tasks              获取任务列表(支持分页：page, size, status)
- GET    /api/v1/tasks/{id}         获取任务详情
- PUT    /api/v1/tasks/{id}         更新任务
- DELETE /api/v1/tasks/{id}         删除任务

任务状态管理：
- GET    /api/v1/tasks/{id}/status  获取任务状态
- PUT    /api/v1/tasks/{id}/status  修改任务状态

任务执行控制：
- POST   /api/v1/tasks/{id}/execute    执行任务
- POST   /api/v1/tasks/{id}/cancel     取消任务
- POST   /api/v1/tasks/{id}/dispatch   分发任务

任务信息查询：
- GET    /api/v1/tasks/query           查询任务
- GET    /api/v1/tasks/{id}/detail     获取任务详情
- GET    /api/v1/tasks/{id}/result     获取任务结果
- GET    /api/v1/tasks/{id}/execution  获取任务执行记录
- GET    /api/v1/tasks/{id}/stats      获取任务统计

总计：14个API端点
状态：连接了真实的TaskService业务逻辑
*/

// TaskController 任务控制器
type TaskController struct {
	taskService service.TaskService
}

// NewTaskController 创建任务控制器实例
func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

// NewTaskControllerFallback 提供无 service 的占位，防止其他包误用
// Deprecated: 请使用 NewTaskController(service) 注入依赖
func NewTaskControllerFallback() *TaskController {
	return &TaskController{}
}

// CreateTask 创建任务
// @Summary 创建新任务
// @Description 创建一个新的任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body task.TaskCreateRequest true "任务信息"
// @Success 201 {object} object "创建成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks [post]
func (tc *TaskController) CreateTask(c *gin.Context) {
	var req task.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用Service层创建任务
	createdTask, err := tc.taskService.CreateTask(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "任务创建成功",
		"data":    createdTask,
	})
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param task body task.TaskUpdateRequest true "任务更新信息"
// @Success 200 {object} object "更新成功"
// @Failure 400 {object} object "请求参数错误"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id} [put]
func (tc *TaskController) UpdateTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	var req task.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用Service层更新任务
	updatedTask, err := tc.taskService.UpdateTask(c.Request.Context(), taskID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务更新成功",
		"data":    updatedTask,
	})
}

// GetTask 获取单个任务
// @Summary 获取任务详情
// @Description 根据ID获取任务详细信息
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "任务详情"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id} [get]
func (tc *TaskController) GetTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层获取任务
	foundTask, err := tc.taskService.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "获取任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取任务成功",
		"data":    foundTask,
	})
}

// GetTasks 获取任务列表
// @Summary 获取任务列表
// @Description 获取任务列表，支持分页和筛选
// @Tags tasks
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param status query string false "任务状态"
// @Param type query string false "任务类型"
// @Success 200 {object} object "任务列表"
// @Failure 400 {object} object "请求参数错误"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks [get]
func (tc *TaskController) GetTasks(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	// 构建过滤条件
	filter := &task.TaskFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = []string{status}
	}
	if taskType := c.Query("type"); taskType != "" {
		filter.Type = []string{taskType}
	}

	// 构建分页选项
	pagination := &task.PaginationOption{
		Page: page,
		Size: size,
	}

	// 调用Service层获取任务列表
	response, err := tc.taskService.ListTasks(c.Request.Context(), filter, nil, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取任务列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取任务列表成功",
		"data":    response,
	})
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 根据ID删除任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "删除成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id} [delete]
func (tc *TaskController) DeleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层删除任务
	if err := tc.taskService.DeleteTask(c.Request.Context(), taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务删除成功",
		"data":    nil,
	})
}

// GetTaskStatus 获取任务状态
// @Summary 获取任务状态
// @Description 获取任务当前执行状态
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "任务状态"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/status [get]
func (tc *TaskController) GetTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层获取任务状态
	status, err := tc.taskService.GetTaskStatus(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "获取任务状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取任务状态成功",
		"data": gin.H{
			"task_id": taskID,
			"status":  status,
		},
	})
}

// UpdateTaskStatus 修改任务状态
// @Summary 修改任务状态
// @Description 修改任务状态
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param status body object true "任务状态"
// @Success 200 {object} object "修改成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/status [put]
func (tc *TaskController) UpdateTaskStatus(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 调用Service层更新任务状态
	if err := tc.taskService.UpdateTaskStatus(c.Request.Context(), taskID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新任务状态失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务状态更新成功",
		"data": gin.H{
			"task_id": taskID,
			"status":  req.Status,
		},
	})
}

// ExecuteTask 执行任务
// @Summary 执行任务
// @Description 手动触发任务执行
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body task.TaskExecuteRequest false "执行参数"
// @Success 200 {object} object "执行成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/execute [post]
func (tc *TaskController) ExecuteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	var req task.TaskExecuteRequest
	// 执行参数是可选的
	c.ShouldBindJSON(&req)

	// 调用Service层执行任务
	execution, err := tc.taskService.ExecuteTask(c.Request.Context(), taskID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "执行任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务执行成功",
		"data":    execution,
	})
}

// CancelTask 取消任务
// @Summary 取消任务
// @Description 取消任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "取消成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/cancel [post]
func (tc *TaskController) CancelTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层取消任务
	if err := tc.taskService.CancelTask(c.Request.Context(), taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "取消任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务取消成功",
		"data": gin.H{
			"task_id": taskID,
		},
	})
}

// QueryTask 查询任务
// @Summary 查询任务
// @Description 根据关键字查询任务
// @Tags tasks
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键字"
// @Param status query string false "任务状态"
// @Param type query string false "任务类型"
// @Success 200 {object} object "查询成功"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/query [get]
func (tc *TaskController) QueryTask(c *gin.Context) {
	keyword := c.Query("keyword")

	// 构建过滤条件
	filter := &task.TaskFilter{}
	if status := c.Query("status"); status != "" {
		filter.Status = []string{status}
	}
	if taskType := c.Query("type"); taskType != "" {
		filter.Type = []string{taskType}
	}

	// 调用Service层查询任务
	tasks, err := tc.taskService.QueryTasks(c.Request.Context(), keyword, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "查询任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询任务成功",
		"data":    tasks,
	})
}

// GetTaskDetail 任务详情
// @Summary 查询任务详情
// @Description 查询任务详情（包含执行记录）
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "查询成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/detail [get]
func (tc *TaskController) GetTaskDetail(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层获取任务详情
	taskDetail, err := tc.taskService.GetTaskDetail(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "获取任务详情失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取任务详情成功",
		"data":    taskDetail,
	})
}

// GetTaskResult 任务结果查询
// @Summary 查询任务结果
// @Description 查询任务执行结果
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "查询成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/result [get]
func (tc *TaskController) GetTaskResult(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层获取任务结果
	result, err := tc.taskService.GetTaskResult(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "获取任务结果失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取任务结果成功",
		"data":    result,
	})
}

// DispatchTask 任务分发
// @Summary 分发任务
// @Description 将任务分发到执行队列
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "分发成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/dispatch [post]
func (tc *TaskController) DispatchTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层分发任务
	if err := tc.taskService.DispatchTask(c.Request.Context(), taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "分发任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "任务分发成功",
		"data": gin.H{
			"task_id": taskID,
		},
	})
}

// GetTaskExecution 任务执行记录查询
// @Summary 查询任务执行记录
// @Description 查询任务的所有执行记录
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "查询成功"
// @Failure 404 {object} object "任务不存在"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/execution [get]
func (tc *TaskController) GetTaskExecution(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的任务ID",
			"data":    nil,
		})
		return
	}

	// 调用Service层获取执行记录
	executions, err := tc.taskService.GetTaskExecutions(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "获取执行记录失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取执行记录成功",
		"data":    executions,
	})
}

// GetTaskStats 任务统计
// @Summary 查询任务统计
// @Description 查询任务统计信息
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} object "查询成功"
// @Failure 500 {object} object "服务器内部错误"
// @Router /api/v1/tasks/{id}/stats [get]
func (tc *TaskController) GetTaskStats(c *gin.Context) {
	// 构建过滤条件（可以根据查询参数扩展）
	filter := &task.TaskFilter{}

	// 调用Service层获取统计信息
	stats, err := tc.taskService.GetTaskStats(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取统计信息失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取统计信息成功",
		"data":    stats,
	})
}
