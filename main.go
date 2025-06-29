package main

import (
	"log"

	"github.com/YoungBoyGod/OneGoServer/cmd"
	"github.com/YoungBoyGod/OneGoServer/internal/config"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config/server.yaml")
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 创建服务器实例
	server := cmd.NewServer(cfg)

	// 3. 启动服务器
	if err := server.Start(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
