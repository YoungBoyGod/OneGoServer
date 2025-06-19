package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Config is the main configuration struct
type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// AppConfig is the configuration for the app
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// LoadConfig reads config from file
func LoadConfig(filename string) (*Config, error) {
	// 1. 读取文件内容
	configData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 2. 创建配置对象
	var config Config
	// 3. 解析配置对象
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	// 4. 返回配置对象
	return &config, nil
}

// SetupRouter creates and configures the Gin router
func SetupRouter(config *Config) *gin.Engine {
	// 1. 判断config.Server.Mode是否为debug
	if config.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 2. 设置路由
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
