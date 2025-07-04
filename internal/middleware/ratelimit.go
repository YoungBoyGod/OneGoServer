package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenBucket struct {
	capacity int
	tokens   int
	mutex    sync.Mutex
	rate     time.Duration
}

func newBucket(capacity int, rate time.Duration) *tokenBucket {
	b := &tokenBucket{capacity: capacity, tokens: capacity, rate: rate}
	go b.refill()
	return b
}

func (b *tokenBucket) refill() {
	ticker := time.NewTicker(b.rate)
	for range ticker.C {
		b.mutex.Lock()
		if b.tokens < b.capacity {
			b.tokens++
		}
		b.mutex.Unlock()
	}
}

func (b *tokenBucket) Allow() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// RateLimitMiddleware 简单速率限制中间件（全局）
func RateLimitMiddleware() gin.HandlerFunc {
	bucket := newBucket(100, time.Millisecond*10) // 100 QPS 全局
	return func(c *gin.Context) {
		if !bucket.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "请求过多"})
			return
		}
		c.Next()
	}
}
