package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(config *Config) *gin.Engine {
	// 1. 设置Gin模式
	if config.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. 创建路由器
	router := gin.Default()

	// 3. 注册路由
	registerRoutes(router, config)

	return router
}

// registerRoutes 注册所有路由
func registerRoutes(router *gin.Engine, config *Config) {
	// 基础路由
	router.GET("/", homeHandler)
	router.GET("/ping", pingHandler)
	router.GET("/health", healthHandler)
	router.GET("/version", versionHandler(config))

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", apiStatusHandler)
		v1.GET("/info", apiInfoHandler(config))
	}
}

// homeHandler 首页处理器
func homeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "Welcome to learngo0619!",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// pingHandler ping处理器
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// healthHandler 健康检查处理器
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "ok",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	})
}

// versionHandler 版本信息处理器
func versionHandler(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":     config.App.Version,
			"app_name":    config.App.Name,
			"server_mode": config.Server.Mode,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}
}

// apiStatusHandler API状态处理器
func apiStatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"api_version": "v1",
		"status":      "active",
		"endpoints": []string{
			"/api/v1/status",
			"/api/v1/info",
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// apiInfoHandler API信息处理器
func apiInfoHandler(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": map[string]interface{}{
				"name":    config.App.Name,
				"version": config.App.Version,
			},
			"server": map[string]interface{}{
				"host": config.Server.Host,
				"port": config.Server.Port,
				"mode": config.Server.Mode,
			},
			"api": map[string]interface{}{
				"version":   "v1",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	}
}
