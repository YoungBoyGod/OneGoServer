package routes

import (
	"learngo0619/internal/config"
	"learngo0619/internal/server/api/handlers"
	"learngo0619/internal/server/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 1. 设置Gin模式
	if cfg.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. 创建路由器 (不使用默认中间件)
	router := gin.New()

	// 3. 注册路由
	RegisterRoutes(router, cfg)

	return router
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(engine *gin.Engine, cfg *config.Config) {
	// 注册全局中间件
	engine.Use(middleware.LoggerMiddleware())   // 使用新的客户端日志中间件
	engine.Use(middleware.RecoveryMiddleware()) // 使用新的恢复中间件
	engine.Use(middleware.CORSMiddleware())     // 跨域中间件
	engine.Use(middleware.SecurityMiddleware()) // 安全中间件

	// 增加favicon.ico静态文件路由
	engine.StaticFile("/favicon.ico", "favicon.ico")

	// 注册具体路由
	registerRoutes(engine, cfg)
}

// registerRoutes 注册所有路由
func registerRoutes(router *gin.Engine, cfg *config.Config) {
	// 基础路由
	router.GET("/", handlers.HomeHandler)
	router.GET("/ping", handlers.PingHandler)
	router.GET("/health", handlers.HealthHandler)
	router.GET("/version", handlers.VersionHandler(cfg))

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", handlers.APIStatusHandler)
		v1.GET("/info", handlers.APIInfoHandler(cfg))
		v1.GET("/client-info", handlers.ClientInfoHandler)
		v1.GET("/log-stats", handlers.LogStatsHandler) // 日志统计端点

		// 客户端相关接口（简化版 - 仅预定义Token认证）
		clients := v1.Group("/clients")
		{
			clients.POST("/register", handlers.ClientRegisterHandler)   // 注册客户端（需要预定义Token）
			clients.POST("/heartbeat", handlers.ClientHeartbeatHandler) // 客户端心跳
			clients.GET("", handlers.ClientListHandler)                 // 客户端列表
			clients.GET("/online", handlers.ClientOnlineHandler)        // 在线客户端
			clients.GET("/stats", handlers.ClientStatsHandler)          // 客户端统计
			clients.GET("/types", handlers.ClientTypesHandler)          // 支持的客户端类型
			clients.GET("/:id", handlers.ClientDetailHandler)           // 客户端详情
			clients.DELETE("/:id", handlers.ClientUnregisterHandler)    // 客户端注销

			// 客户端任务相关接口
			clients.GET("/:clientId/tasks", handlers.GetTasksForClientHandler)           // 获取客户端待执行任务
			clients.GET("/:clientId/task-results", handlers.GetClientTaskResultsHandler) // 获取客户端任务执行结果
		}

		// 任务管理接口（管理员使用）
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", handlers.CreateTaskHandler)                     // 创建任务
			tasks.GET("", handlers.GetAllTasksHandler)                     // 获取所有任务
			tasks.GET("/stats", handlers.GetTaskStatsHandler)              // 获取任务统计信息
			tasks.GET("/:taskId", handlers.GetTaskHandler)                 // 获取任务详情
			tasks.PUT("/:taskId/status", handlers.UpdateTaskStatusHandler) // 更新任务状态
			tasks.GET("/:taskId/results", handlers.GetTaskResultsHandler)  // 获取任务执行结果
			tasks.DELETE("/:taskId", handlers.DeleteTaskHandler)           // 删除任务
		}

		// 客户端任务轮询接口（客户端使用）
		v1.GET("/poll-tasks", handlers.TaskPollHandler) // 客户端轮询任务
	}
}
