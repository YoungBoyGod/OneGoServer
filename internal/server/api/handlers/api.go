package handlers

import (
	"time"

	"learngo0619/internal/config"

	"github.com/gin-gonic/gin"
)

// VersionHandler 版本信息处理器
func VersionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":     cfg.App.Version,
			"app_name":    cfg.App.Name,
			"server_mode": cfg.Server.Mode,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}
}

// APIStatusHandler API状态处理器
func APIStatusHandler(c *gin.Context) {
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

// APIInfoHandler API信息处理器
func APIInfoHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"app": map[string]interface{}{
				"name":    cfg.App.Name,
				"version": cfg.App.Version,
			},
			"server": map[string]interface{}{
				"host": cfg.Server.Host,
				"port": cfg.Server.Port,
				"mode": cfg.Server.Mode,
			},
			"api": map[string]interface{}{
				"version":   "v1",
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	}
}
