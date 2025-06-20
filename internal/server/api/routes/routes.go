package routes

import (
	"learngo0619/internal/config"
	"learngo0619/internal/logger"
	"learngo0619/internal/server/api/handlers"
	"learngo0619/internal/server/api/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		}
	}

	// 记录认证模式
	logger.Info("Predefined token authentication mode enabled",
		zap.String("token_configured", func() string {
			if cfg.Server.Security.PredefinedToken != "" {
				return "yes"
			}
			return "no"
		}()),
	)
}
