package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简易 JWT 校验（示例）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "缺少或无效的令牌"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		// TODO: 解析 token 校验，此处仅示例
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "无效令牌"})
			return
		}
		c.Set("user_id", "demo")
		c.Next()
	}
}
