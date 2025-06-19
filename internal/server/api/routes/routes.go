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

	// 3. 添加自定义中间件
	router.Use(middleware.RequestID())   // 请求ID
	router.Use(middleware.ZapLogger())   // Zap日志
	router.Use(middleware.ZapRecovery()) // Zap恢复

	// 4. 注册路由
	registerRoutes(router, cfg)

	return router
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
	}
}
