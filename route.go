package main

import "github.com/gin-gonic/gin"

// setupRouter is the function to setup the router
func setupRouter(config *Config) *gin.Engine {
	//  1. 判断config.Server.Mode是否为debug
	if config.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//  2. 设置路由
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
			"status":  "healthy",
			"time":    "2025-01-19",
		})
	})
	router.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": config.App.Version})
	})
	return router
}
