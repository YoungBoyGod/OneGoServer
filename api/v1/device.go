package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设备接口
func Device(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Device!"})
}
