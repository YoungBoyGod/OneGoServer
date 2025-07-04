package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorRecovery 统一错误恢复与JSON返回
func ErrorRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误", "error": err})
	})
}
