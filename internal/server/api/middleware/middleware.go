package middleware

import (
	"crypto/md5"
	"fmt"
	"time"

	"learngo0619/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start).Seconds()

		// 获取客户端信息
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("status", statusCode),
			zap.Float64("latency", latency),
			zap.String("time", time.Now().Format(time.RFC3339)),
		}

		// 根据状态码选择日志级别并记录到客户端专用日志
		if statusCode >= 400 {
			logger.ErrorForClient(clientIP, userAgent, "HTTP Request", fields...)
		} else if statusCode >= 300 {
			logger.WarnForClient(clientIP, userAgent, "HTTP Request", fields...)
		} else {
			logger.InfoForClient(clientIP, userAgent, "HTTP Request", fields...)
		}
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				clientIP := c.ClientIP()
				userAgent := c.Request.UserAgent()

				// 记录panic到客户端专用日志
				logger.ErrorForClient(clientIP, userAgent, "Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				c.JSON(500, gin.H{
					"error": "Internal Server Error",
					"code":  500,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RequestID 为每个请求添加唯一ID和客户端指纹
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 获取会话ID
		_, sessionID, serverMac := logger.GetInstanceInfo()

		// 生成客户端指纹（作为客户端MAC的替代）
		clientFingerprint := generateClientFingerprint(c)

		c.Set("request_id", requestID)
		c.Set("session_id", sessionID)
		c.Set("server_mac", serverMac)
		c.Set("client_fingerprint", clientFingerprint)

		// 设置响应头
		c.Header("X-Request-ID", requestID)
		c.Header("X-Session-ID", sessionID)
		c.Header("X-Server-MAC", serverMac)
		c.Header("X-Client-Fingerprint", clientFingerprint)

		c.Next()
	}
}

// generateClientFingerprint 生成客户端指纹（替代MAC地址）
func generateClientFingerprint(c *gin.Context) string {
	// 组合客户端特征信息
	fingerprint := fmt.Sprintf("%s|%s|%s|%s",
		c.ClientIP(),
		c.Request.UserAgent(),
		c.GetHeader("Accept-Language"),
		c.GetHeader("Accept-Encoding"),
	)

	// 生成MD5哈希作为指纹
	hash := md5.Sum([]byte(fingerprint))
	return fmt.Sprintf("%x", hash)[:16] // 取前16位
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(b)
}
