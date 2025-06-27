package controller

import (
	"github.com/YoungBoyGod/OneGoServer/internal/config"
)

// HealthController 健康检查控制器
type HealthController struct {
	config *config.Config
}

// NewHealthController 创建健康检查控制器
func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{
		config: cfg,
	}
}
