package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// HomeHandler 首页处理器
func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "Welcome to learngo0619!",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// PingHandler ping处理器
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// HealthHandler 健康检查处理器
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "ok",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    "running",
	})
}
