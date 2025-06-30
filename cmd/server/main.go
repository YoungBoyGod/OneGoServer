package main

import (
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/data/postgres"
	"github.com/YoungBoyGod/OneGoServer/internal/data/redis"
	"github.com/YoungBoyGod/OneGoServer/internal/router"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 初始化数据库连接
	db, err := postgres.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	defer postgres.CloseDB(db)

	// 3. 初始化Redis连接
	rdb, err := redis.InitRedis(cfg)
	if err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}
	defer redis.CloseRedis(rdb)

	// 4. 初始化Kafka生产者
	producer, err := kafka.InitProducer(cfg)
	if err != nil {
		log.Fatalf("❌ Kafka生产者初始化失败: %v", err)
	}
	defer kafka.CloseProducer(producer)

	// 5. 启动Kafka消费者 (可选，用于异步任务处理)
	consumer, err := kafka.InitConsumer(cfg)
	if err != nil {
		log.Printf("⚠️ Kafka消费者初始化失败: %v", err)
	} else {
		defer kafka.CloseConsumer(consumer)
	}

	// 6. 初始化并启动HTTP服务器
	r := router.InitRouter(cfg, db, rdb, producer)

	log.Printf("🚀 服务器启动中... 监听端口: %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
