package router

import (
	"net/http"

	"github.com/YoungBoyGod/OneGoServer/internal/data/repository"
	"github.com/YoungBoyGod/OneGoServer/internal/service"
	"github.com/YoungBoyGod/OneGoServer/pkg/sql"

	"github.com/YoungBoyGod/OneGoServer/internal/controller"
	"github.com/YoungBoyGod/OneGoServer/internal/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由配置
func InitRouter() *gin.Engine {
	// 创建Gin引擎
	r := gin.Default()

	// 添加全局中间件
	r.Use(corsMiddleware())    // CORS中间件
	r.Use(loggingMiddleware()) // 日志中间件
	r.Use(middleware.TraceMiddleware())
	r.Use(middleware.ErrorRecovery())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(middleware.AuthMiddleware())

	// 健康检查端点
	r.GET("/health", healthCheck)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// API版本路由组
	v1 := r.Group("/api/v1")
	{
		// 注册Task路由
		registerTaskRoutes(v1)

		// 注册Device路由
		registerDeviceRoutes(v1)

		// 注册队列路由
		registerQueueRoutes(v1)
	}

	return r
}

// registerTaskRoutes 注册任务相关路由
func registerTaskRoutes(rg *gin.RouterGroup) {
	db := sql.GetDB()
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskController := controller.NewTaskController(taskService)

	// 任务路由组
	tasks := rg.Group("/tasks")
	{
		// 基础CRUD操作
		tasks.POST("", taskController.CreateTask)       // POST /api/v1/tasks
		tasks.GET("", taskController.GetTasks)          // GET /api/v1/tasks
		tasks.GET("/:id", taskController.GetTask)       // GET /api/v1/tasks/:id
		tasks.PUT("/:id", taskController.UpdateTask)    // PUT /api/v1/tasks/:id
		tasks.DELETE("/:id", taskController.DeleteTask) // DELETE /api/v1/tasks/:id

		// 任务状态管理
		tasks.GET("/:id/status", taskController.GetTaskStatus)    // GET /api/v1/tasks/:id/status
		tasks.PUT("/:id/status", taskController.UpdateTaskStatus) // PUT /api/v1/tasks/:id/status

		// 任务执行控制
		tasks.POST("/:id/execute", taskController.ExecuteTask)   // POST /api/v1/tasks/:id/execute
		tasks.POST("/:id/cancel", taskController.CancelTask)     // POST /api/v1/tasks/:id/cancel
		tasks.POST("/:id/dispatch", taskController.DispatchTask) // POST /api/v1/tasks/:id/dispatch

		// 任务信息查询
		tasks.GET("/query", taskController.QueryTask)                // GET /api/v1/tasks/query
		tasks.GET("/:id/detail", taskController.GetTaskDetail)       // GET /api/v1/tasks/:id/detail
		tasks.GET("/:id/result", taskController.GetTaskResult)       // GET /api/v1/tasks/:id/result
		tasks.GET("/:id/execution", taskController.GetTaskExecution) // GET /api/v1/tasks/:id/execution
		tasks.GET("/:id/stats", taskController.GetTaskStats)         // GET /api/v1/tasks/:id/stats

		// 批量任务控制
		tasks.POST("/batch/cancel", taskController.BatchCancelTasks)
		tasks.POST("/batch/dispatch", taskController.BatchDispatchTasks)
	}
}

// registerDeviceRoutes 注册设备相关路由
func registerDeviceRoutes(rg *gin.RouterGroup) {
	db := sql.GetDB()
	deviceRepo := repository.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)
	deviceController := controller.NewDeviceController(deviceService)

	// 设备路由组
	devices := rg.Group("/devices")
	{
		// 基础CRUD操作
		devices.POST("", deviceController.RegisterDevice)     // POST /api/v1/devices
		devices.GET("", deviceController.GetDevices)          // GET /api/v1/devices
		devices.GET("/:id", deviceController.GetDevice)       // GET /api/v1/devices/:id
		devices.PUT("/:id", deviceController.UpdateDevice)    // PUT /api/v1/devices/:id
		devices.DELETE("/:id", deviceController.DeleteDevice) // DELETE /api/v1/devices/:id

		// 设备状态管理
		devices.GET("/:id/status", deviceController.GetDeviceStatus) // GET /api/v1/devices/:id/status
		devices.POST("/online", deviceController.DeviceOnline)       // POST /api/v1/devices/online
		devices.POST("/offline", deviceController.DeviceOffline)     // POST /api/v1/devices/offline

		// 设备控制操作
		devices.POST("/:id/command", deviceController.SendCommand)             // POST /api/v1/devices/:id/command
		devices.GET("/:id/heartbeat", deviceController.GetDeviceHeartbeat)     // GET /api/v1/devices/:id/heartbeat
		devices.POST("/:id/heartbeat", deviceController.UpdateDeviceHeartbeat) // POST /api/v1/devices/:id/heartbeat

		// 设备信息查询
		devices.GET("/query", deviceController.QueryDevice)        // GET /api/v1/devices/query
		devices.GET("/:id/logs", deviceController.GetDeviceLogs)   // GET /api/v1/devices/:id/logs
		devices.GET("/:id/stats", deviceController.GetDeviceStats) // GET /api/v1/devices/:id/stats
	}
}

// registerQueueRoutes 注册队列相关路由
func registerQueueRoutes(rg *gin.RouterGroup) {
	db := sql.GetDB()
	taskRepo := repository.NewTaskRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	queueService := service.NewQueueService(taskRepo, deviceRepo)
	queueController := controller.NewQueueController(queueService)

	q := rg.Group("/device-queues")
	{
		q.POST("/enqueue", queueController.EnqueueTask)               // POST /device-queues/enqueue
		q.POST("/:device_id/dequeue", queueController.DequeueTask)    // POST /device-queues/:device_id/dequeue
		q.GET("/:device_id/metrics", queueController.GetQueueMetrics) // GET  /device-queues/:device_id/metrics
		q.POST("/:device_id/reorder", queueController.ReorderQueue)   // POST /device-queues/:device_id/reorder
	}
}

// healthCheck 健康检查处理函数
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"message":   "OneGo服务器运行正常",
		"timestamp": c.GetHeader("X-Request-ID"),
		"services": gin.H{
			"database": "connected", // TODO: 实际检查数据库连接
			"redis":    "connected", // TODO: 实际检查Redis连接
			"kafka":    "connected", // TODO: 实际检查Kafka连接
		},
	})
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从配置文件读取CORS设置
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// loggingMiddleware 日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// TODO: 集成统一日志系统
		return ""
	})
}
