package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 任务接口
func Task(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, Task!"})
}
