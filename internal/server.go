package internal

import (
	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/gin-gonic/gin"
)

func InitServer(cfg *config.Config) (*gin.Engine, error) {
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	return engine, nil
}
