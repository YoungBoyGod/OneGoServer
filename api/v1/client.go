package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 客户端接口
func Client(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Client!"})
}
