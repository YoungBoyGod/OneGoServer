package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 队列接口
func Queue(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Queue!"})
}
