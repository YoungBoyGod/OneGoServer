package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthCheckRes struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckRes{
		Status:  "ok",
		Version: "1.0.0",
		Message: "Hello, World!",
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	})
}
