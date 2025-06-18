package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	srv := &http.Server{
		Addr:    ":30000",
		Handler: router,
	}
	srv.ListenAndServe() // 监听并在 0.0.0.0:30000 上启动服务
}
