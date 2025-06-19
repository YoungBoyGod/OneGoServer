package middleware

import (
	"crypto/md5"
	"fmt"
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
