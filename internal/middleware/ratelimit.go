package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterMap 用于存储每个IP的限流器
var limiterMap = make(map[string]*rate.Limiter)
var limiterMutex sync.RWMutex

// RateLimit 限流中间件
func RateLimit(requestsPerMinute, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 获取或创建限流器
		limiter := getLimiter(ip, requestsPerMinute, burst)

		// 检查是否允许请求
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getLimiter 获取或创建IP对应的限流器
func getLimiter(ip string, requestsPerMinute, burst int) *rate.Limiter {
	limiterMutex.RLock()
	limiter, exists := limiterMap[ip]
	limiterMutex.RUnlock()

	if !exists {
		limiterMutex.Lock()
		// 双重检查
		if limiter, exists = limiterMap[ip]; !exists {
			// 创建新的限流器
			limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst)
			limiterMap[ip] = limiter
		}
		limiterMutex.Unlock()
	}

	return limiter
}

// cleanupLimiters 清理过期的限流器（可以定期调用）
func cleanupLimiters() {
	limiterMutex.Lock()
	defer limiterMutex.Unlock()

	// 清空所有限流器（简单实现，实际可以根据时间戳判断是否过期）
	limiterMap = make(map[string]*rate.Limiter)
}
