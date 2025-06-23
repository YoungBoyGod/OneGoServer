package handlers

import (
	"net/http"
	"strconv"

	"learngo0619/internal/logger"
	"learngo0619/internal/server/models"
	"learngo0619/internal/server/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateTaskHandler 创建任务处理器
func CreateTaskHandler(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid create task request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 获取创建者信息（可以从认证信息中获取，这里先使用默认值）
	createdBy := "admin" // TODO: 从认证信息中获取

	// 创建任务
	taskManager := services.GetTaskManager()
	task, err := taskManager.CreateTask(&req, createdBy)
	if err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建任务失败: " + err.Error(),
		})
		return
	}

	logger.Info("Task created successfully",
		zap.String("task_id", task.ID),
		zap.String("task_name", task.Name),
		zap.String("script_type", string(task.ScriptType)),
		zap.String("created_by", createdBy),
	)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "任务创建成功",
		"task":    task,
	})
}

// GetTasksForClientHandler 获取客户端待执行任务处理器
func GetTasksForClientHandler(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户端ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	tasks, err := taskManager.GetTasksForClient(clientID)
	if err != nil {
		logger.Error("Failed to get tasks for client",
			zap.String("client_id", clientID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取任务失败: " + err.Error(),
		})
		return
	}

	// Debug模式下记录详细信息
	serverConfig := getServerConfig()
	if serverConfig != nil && serverConfig.Server.Mode == "debug" {
		logger.Info("Debug: Tasks retrieved for client",
			zap.String("client_id", clientID),
			zap.Int("task_count", len(tasks)),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tasks":   tasks,
		"count":   len(tasks),
	})
}

// GetAllTasksHandler 获取所有任务处理器
func GetAllTasksHandler(c *gin.Context) {
	taskManager := services.GetTaskManager()
	tasks, err := taskManager.GetAllTasks()
	if err != nil {
		logger.Error("Failed to get all tasks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取任务列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tasks":   tasks,
		"total":   len(tasks),
	})
}

// GetTaskHandler 获取单个任务详情处理器
func GetTaskHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "任务ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	task, err := taskManager.GetTask(taskID)
	if err != nil {
		logger.Error("Failed to get task",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"task":    task,
	})
}

// UpdateTaskStatusHandler 更新任务状态处理器
func UpdateTaskStatusHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "任务ID不能为空",
		})
		return
	}

	var req models.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid update task status request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 获取客户端ID（从请求头或查询参数获取）
	clientID := c.GetHeader("X-Client-ID")
	if clientID == "" {
		clientID = c.Query("client_id")
	}

	taskManager := services.GetTaskManager()

	// 如果有执行结果相关信息，记录到结果表
	if req.StartTime != nil || req.EndTime != nil || req.Output != "" || req.Error != "" {
		result, err := taskManager.RecordTaskResult(taskID, clientID, &req)
		if err != nil {
			logger.Error("Failed to record task result",
				zap.String("task_id", taskID),
				zap.String("client_id", clientID),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "记录任务结果失败: " + err.Error(),
			})
			return
		}

		logger.Info("Task result recorded",
			zap.String("task_id", taskID),
			zap.String("client_id", clientID),
			zap.String("status", string(req.Status)),
			zap.String("result_id", result.ID),
		)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "任务结果记录成功",
			"result":  result,
		})
		return
	}

	// 只更新任务状态
	err := taskManager.UpdateTaskStatus(taskID, &req)
	if err != nil {
		logger.Error("Failed to update task status",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新任务状态失败: " + err.Error(),
		})
		return
	}

	logger.Info("Task status updated",
		zap.String("task_id", taskID),
		zap.String("client_id", clientID),
		zap.String("status", string(req.Status)),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "任务状态更新成功",
	})
}

// GetTaskResultsHandler 获取任务执行结果处理器
func GetTaskResultsHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "任务ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	results, err := taskManager.GetTaskResults(taskID)
	if err != nil {
		logger.Error("Failed to get task results",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取任务结果失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"results": results,
		"count":   len(results),
	})
}

// GetClientTaskResultsHandler 获取客户端任务执行结果处理器
func GetClientTaskResultsHandler(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户端ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	results, err := taskManager.GetClientTaskResults(clientID)
	if err != nil {
		logger.Error("Failed to get client task results",
			zap.String("client_id", clientID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取客户端任务结果失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"results": results,
		"count":   len(results),
	})
}

// DeleteTaskHandler 删除任务处理器
func DeleteTaskHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "任务ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	err := taskManager.DeleteTask(taskID)
	if err != nil {
		logger.Error("Failed to delete task",
			zap.String("task_id", taskID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除任务失败: " + err.Error(),
		})
		return
	}

	logger.Info("Task deleted successfully",
		zap.String("task_id", taskID),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "任务删除成功",
	})
}

// GetTaskStatsHandler 获取任务统计信息处理器
func GetTaskStatsHandler(c *gin.Context) {
	taskManager := services.GetTaskManager()
	stats := taskManager.GetTaskStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}

// TaskPollHandler 客户端轮询任务处理器（简化版本，获取单个任务）
func TaskPollHandler(c *gin.Context) {
	clientID := c.GetHeader("X-Client-ID")
	if clientID == "" {
		clientID = c.Query("client_id")
	}

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户端ID不能为空",
		})
		return
	}

	taskManager := services.GetTaskManager()
	tasks, err := taskManager.GetTasksForClient(clientID)
	if err != nil {
		logger.Error("Failed to poll tasks for client",
			zap.String("client_id", clientID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "轮询任务失败: " + err.Error(),
		})
		return
	}

	// 限制返回的任务数量（客户端一次只处理一个任务）
	limit := 1
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var returnTasks []*models.Task
	if len(tasks) > 0 {
		if len(tasks) < limit {
			returnTasks = tasks
		} else {
			returnTasks = tasks[:limit]
		}

		// 将返回的任务标记为运行中
		for _, task := range returnTasks {
			updateReq := &models.UpdateTaskStatusRequest{
				Status: models.TaskStatusRunning,
			}
			taskManager.UpdateTaskStatus(task.ID, updateReq)
		}
	}

	// Debug模式下记录轮询信息
	serverConfig := getServerConfig()
	if serverConfig != nil && serverConfig.Server.Mode == "debug" {
		logger.Info("Debug: Client task polling",
			zap.String("client_id", clientID),
			zap.Int("available_tasks", len(tasks)),
			zap.Int("returned_tasks", len(returnTasks)),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"tasks":    returnTasks,
		"count":    len(returnTasks),
		"has_more": len(tasks) > len(returnTasks),
	})
}
