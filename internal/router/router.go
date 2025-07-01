package router

import (
	"net/http"

	"github.com/YoungBoyGod/OneGoServer/internal/controller"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由配置
func InitRouter() *gin.Engine {
	// 创建Gin引擎
	r := gin.Default()

	// 添加全局中间件
	r.Use(corsMiddleware())    // CORS中间件
	r.Use(loggingMiddleware()) // 日志中间件
	r.Use(gin.Recovery())      // 恢复中间件

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
	}

	return r
}

// registerTaskRoutes 注册任务相关路由
func registerTaskRoutes(rg *gin.RouterGroup) {
	taskController := controller.NewTaskController()

	// 任务路由组
	tasks := rg.Group("/tasks")
	{
		// 基础CRUD操作
		tasks.POST("", taskController.CreateTask)       // POST /api/v1/tasks
		tasks.GET("", taskController.GetTasks)          // GET /api/v1/tasks
		tasks.GET("/:id", taskController.GetTask)       // GET /api/v1/tasks/:id
		tasks.PUT("/:id", taskController.UpdateTask)    // PUT /api/v1/tasks/:id
		tasks.DELETE("/:id", taskController.DeleteTask) // DELETE /api/v1/tasks/:id

		// 任务状态和执行操作
		tasks.GET("/:id/status", taskController.GetTaskStatus) // GET /api/v1/tasks/:id/status
		tasks.POST("/:id/execute", taskController.ExecuteTask) // POST /api/v1/tasks/:id/execute
	}
}

// registerDeviceRoutes 注册设备相关路由
func registerDeviceRoutes(rg *gin.RouterGroup) {
	deviceController := controller.NewDeviceController()

	// 设备路由组
	devices := rg.Group("/devices")
	{
		// 基础CRUD操作
		devices.POST("", deviceController.RegisterDevice)     // POST /api/v1/devices
		devices.GET("", deviceController.GetDevices)          // GET /api/v1/devices
		devices.GET("/:id", deviceController.GetDevice)       // GET /api/v1/devices/:id
		devices.PUT("/:id", deviceController.UpdateDevice)    // PUT /api/v1/devices/:id
		devices.DELETE("/:id", deviceController.DeleteDevice) // DELETE /api/v1/devices/:id

		// 设备状态和控制操作
		devices.GET("/:id/status", deviceController.GetDeviceStatus)           // GET /api/v1/devices/:id/status
		devices.POST("/:id/command", deviceController.SendCommand)             // POST /api/v1/devices/:id/command
		devices.GET("/:id/heartbeat", deviceController.GetDeviceHeartbeat)     // GET /api/v1/devices/:id/heartbeat
		devices.POST("/:id/heartbeat", deviceController.UpdateDeviceHeartbeat) // POST /api/v1/devices/:id/heartbeat
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

// authMiddleware 认证中间件 (占位)
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现JWT认证逻辑
		// 1. 从Header获取Authorization
		// 2. 验证JWT Token
		// 3. 解析用户信息
		// 4. 设置用户上下文
		c.Next()
	}
}

// rateLimitMiddleware 限流中间件 (占位)
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现API限流逻辑
		// 1. 获取客户端IP
		// 2. 检查请求频率
		// 3. 实施限流策略
		c.Next()
	}
}
