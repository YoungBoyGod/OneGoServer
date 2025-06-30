package server

import (
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/data/postgres"
	"github.com/YoungBoyGod/OneGoServer/internal/data/redis"
	"github.com/YoungBoyGod/OneGoServer/internal/server"
)

func InitServer(cfg *config.Config) {
	// 1. 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 初始化数据库连接
	if err := postgres.InitPostgreSQL(cfg); err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Printf("⚠️ 数据库关闭失败: %v", err)
		}
	}()

	// 3. 初始化Redis连接
	if err := redis.InitRedis(cfg); err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("⚠️ Redis关闭失败: %v", err)
		}
	}()

	// 4. 创建并启动服务器
	srv := server.NewServer(cfg, postgres.GetDB(), redis.GetClient())

	log.Printf("🚀 OneGoServer 启动中...")
	if err := srv.Start(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
