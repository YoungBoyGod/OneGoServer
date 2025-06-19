package handlers

import (
	"net/http"
	"time"

	"learngo0619/internal/config"
	"learngo0619/internal/logger"

	"github.com/gin-gonic/gin"
)

// VersionHandler returns the application version
func VersionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":     cfg.App.Version,
			"app_name":    cfg.App.Name,
			"environment": config.GetEnvironment(),
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}
}

// APIStatusHandler returns the API status
func APIStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "active",
		"api_version": "v1",
		"timestamp":   time.Now().Format(time.RFC3339),
		"endpoints": []string{
			"/api/v1/server/status",
			"/api/v1/server/info",
			"/api/v1/server/client-info",
			"/api/v1/server/log-stats",
			// 客户端注册管理端点
			"/api/v1/clients/register",  // POST - 客户端注册
			"/api/v1/clients/heartbeat", // POST - 客户端心跳
			"/api/v1/clients",           // GET - 客户端列表
			"/api/v1/clients/online",    // GET - 在线客户端
			"/api/v1/clients/stats",     // GET - 客户端统计
			"/api/v1/clients/types",     // GET - 客户端类型
			"/api/v1/clients/:id",       // GET - 客户端详情
			"/api/v1/clients/:id",       // DELETE - 客户端注销
		},
	})
}

// APIInfoHandler returns detailed API information
func APIInfoHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"api": gin.H{
				"version":   "v1",
				"timestamp": time.Now().Format(time.RFC3339),
			},
			"server": gin.H{
				"host": cfg.Server.Host,
				"port": cfg.Server.Port,
				"mode": cfg.Server.Mode,
			},
			"app": gin.H{
				"name":    cfg.App.Name,
				"version": cfg.App.Version,
			},
		})
	}
}

// ClientInfoHandler 返回客户端和服务器信息（包含MAC地址信息）
func ClientInfoHandler(c *gin.Context) {
	// 获取服务器实例信息
	pid, sessionID, serverMac := logger.GetInstanceInfo()

	// 从中间件获取客户端信息
	requestID, _ := c.Get("request_id")
	clientFingerprint, _ := c.Get("client_fingerprint")

	c.JSON(http.StatusOK, gin.H{
		"timestamp": time.Now().Format(time.RFC3339),
		"server": gin.H{
			"pid":         pid,
			"session_id":  sessionID,
			"mac_address": serverMac, // 服务端真实MAC地址
		},
		"client": gin.H{
			"ip":              c.ClientIP(),
			"user_agent":      c.Request.UserAgent(),
			"fingerprint":     clientFingerprint, // 客户端指纹（MAC替代）
			"accept_language": c.GetHeader("Accept-Language"),
			"accept_encoding": c.GetHeader("Accept-Encoding"),
		},
		"request": gin.H{
			"id":     requestID,
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		},
		"note": "服务端可获取真实MAC地址，客户端MAC地址因安全限制无法获取，使用指纹替代",
	})
}

// LogStatsHandler 返回日志统计信息
func LogStatsHandler(c *gin.Context) {
	stats := logger.GetClientLoggerStats()

	c.JSON(http.StatusOK, gin.H{
		"timestamp":   time.Now().Format(time.RFC3339),
		"client_logs": stats,
		"description": "客户端日志分离统计信息",
	})
}
