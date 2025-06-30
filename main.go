package main

import (
	"fmt"
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/YoungBoyGod/OneGoServer/pkg/sql"
)

func main() {
	fmt.Println("Hello, World!")

	// 加载配置文件
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Server started on port %d", cfg.Server.Port)

	// 打印配置文件内容
	log.Printf("Database Config: %+v", cfg.Database)
	log.Printf("Server Config: %+v", cfg.Server)

	// 校验配置参数
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	log.Printf("Database DSN: %s", cfg.Database.GetDsn())

	// 初始化增强版日志系统
	if err := pkglog.InitLoggerEnhanced(&cfg.Logging); err != nil {
		log.Fatalf("Failed to initialize enhanced logger: %v", err)
	}
	log.Printf("Enhanced Logger initialized successfully")

	// 初始化数据库
	sql.InitDB(cfg)

	// 使用配置的日志器记录信息
	logger := pkglog.GetAppLogger(&cfg.Logging)
	logger.Info("Application started successfully")
	logger.Info("Database connected successfully")

	// 同步日志
	pkglog.Sync()
	log.Printf("Application setup completed")
}
