package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Dependencies 路由依赖
type Dependencies struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// 临时处理器函数 - 用于占位，避免编译错误
func placeholder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "API endpoint under development",
		"path":    c.FullPath(),
		"method":  c.Request.Method,
	})
}

// NewRouter 创建路由，接受依赖注入
func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.Default()

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 健康检查 - 可以检查数据库和Redis连接状态
		v1.GET("/health", func(c *gin.Context) {
			status := "ok"
			checks := make(map[string]string)

			// 检查数据库连接
			if deps.DB != nil {
				if sqlDB, err := deps.DB.DB(); err == nil && sqlDB.Ping() == nil {
					checks["database"] = "healthy"
				} else {
					checks["database"] = "unhealthy"
					status = "degraded"
				}
			} else {
				checks["database"] = "not_configured"
			}

			// 检查Redis连接
			if deps.Redis != nil {
				if deps.Redis.Ping(c.Request.Context()).Err() == nil {
					checks["redis"] = "healthy"
				} else {
					checks["redis"] = "unhealthy"
					status = "degraded"
				}
			} else {
				checks["redis"] = "not_configured"
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  status,
				"message": "OneGoServer is running",
				"checks":  checks,
			})
		})

		// 认证路由组
		auth := v1.Group("/auth")
		{
			auth.POST("/register", placeholder) // 用户注册
			auth.POST("/login", placeholder)    // 用户登录
			auth.POST("/logout", placeholder)   // 用户登出
		}

		// 用户管理路由组
		users := v1.Group("/users")
		{
			users.GET("", placeholder)              // 获取用户列表
			users.POST("", placeholder)             // 创建用户
			users.GET("/:id", placeholder)          // 获取用户详情
			users.PUT("/:id", placeholder)          // 更新用户信息
			users.DELETE("/:id", placeholder)       // 删除用户
			users.PUT("/:id/password", placeholder) // 修改密码
		}

		// 设备管理路由组
		devices := v1.Group("/devices")
		{
			devices.GET("", placeholder)               // 获取设备列表
			devices.POST("/register", placeholder)     // 注册设备
			devices.POST("/online", placeholder)       // 设备上线
			devices.POST("/offline", placeholder)      // 设备下线
			devices.POST("/meta", placeholder)         // 设备元数据
			devices.GET("/stats", placeholder)         // 获取设备统计信息
			devices.GET("/:id", placeholder)           // 获取设备详情
			devices.PUT("/:id", placeholder)           // 更新设备信息
			devices.DELETE("/:id", placeholder)        // 删除设备
			devices.PUT("/:id/heartbeat", placeholder) // 设备心跳
			devices.PUT("/:id/status", placeholder)    // 更新设备状态
		}

		// 任务管理路由组
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", placeholder)             // 获取任务列表
			tasks.POST("", placeholder)            // 创建任务
			tasks.GET("/:id", placeholder)         // 获取任务详情
			tasks.PUT("/:id", placeholder)         // 更新任务信息
			tasks.DELETE("/:id", placeholder)      // 删除任务
			tasks.POST("/:id/cancel", placeholder) // 取消任务
			tasks.POST("/:id/retry", placeholder)  // 重试任务
			tasks.GET("/:id/logs", placeholder)    // 获取任务日志
		}

		// 告警管理路由组 (预留)
		alerts := v1.Group("/alerts")
		{
			alerts.GET("", placeholder)                  // 获取告警列表
			alerts.POST("", placeholder)                 // 创建告警
			alerts.GET("/stats", placeholder)            // 获取告警统计
			alerts.GET("/:id", placeholder)              // 获取告警详情
			alerts.POST("/:id/acknowledge", placeholder) // 确认告警
			alerts.POST("/:id/resolve", placeholder)     // 解决告警
			alerts.POST("/:id/close", placeholder)       // 关闭告警
		}
	}

	return router
}
