package middleware

import (
	"time"

	"learngo0619/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger 返回使用zap的Gin日志中间件
func ZapLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 使用zap记录访问日志
		fields := []zap.Field{
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("query", param.Request.URL.RawQuery),
			zap.String("ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("time", param.TimeStamp.Format(time.RFC3339)),
		}

		if param.ErrorMessage != "" {
			fields = append(fields, zap.String("error", param.ErrorMessage))
		}

		// 根据状态码决定日志级别
		switch {
		case param.StatusCode >= 500:
			logger.Error("HTTP Request", fields...)
		case param.StatusCode >= 400:
			logger.Warn("HTTP Request", fields...)
		default:
			logger.Info("HTTP Request", fields...)
		}

		return ""
	})
}

// ZapRecovery 返回使用zap的恢复中间件
func ZapRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered",
				zap.String("error", err),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
		}
		c.AbortWithStatus(500)
	})
}

// RequestID 为每个请求添加唯一ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// 获取会话ID
		_, sessionID := logger.GetInstanceInfo()

		c.Set("request_id", requestID)
		c.Set("session_id", sessionID)
		c.Header("X-Request-ID", requestID)
		c.Header("X-Session-ID", sessionID)
		c.Next()
	}
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
