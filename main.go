package main

import (
	"context"
	"fmt"
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/pkg/cache"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/YoungBoyGod/OneGoServer/pkg/sql"
	"go.uber.org/zap"
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
	log.Printf("Redis Config: %+v", cfg.Redis)

	// 校验配置参数
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}
	log.Printf("Database DSN: %s", cfg.Database.GetDsn())

	// 初始化统一日志系统
	if err := pkglog.InitLoggerEnhanced(&cfg.Logging); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.Printf("Logger system initialized successfully")

	// 初始化数据库
	if err := sql.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 使用统一日志系统记录信息
	logger := pkglog.GetAppLogger(&cfg.Logging)
	logger.Info("Application started successfully")
	logger.Info("Database connected successfully")

	// 测试统一日志系统的增强功能
	pkglog.LogSystemEvent("application_start", "main",
		zap.String("version", cfg.Server.Version),
		zap.Int("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode))

	// 测试健康检查
	if err := pkglog.HealthCheck(); err != nil {
		logger.Warn("Logger health check failed", zap.Error(err))
	} else {
		logger.Info("Logger system health check passed")
	}

	// 获取日志统计信息
	stats := pkglog.GetStats()
	logger.Info("Logger statistics",
		zap.Int64("total_logs", stats.TotalLogs),
		zap.Int64("info_logs", stats.InfoLogs),
		zap.Int64("error_logs", stats.ErrorLogs))

	// 同步日志
	pkglog.Sync()
	log.Printf("Application setup completed successfully")

	// 初始化Redis
	if err := cache.InitRedis(context.Background(), &cfg.Redis); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 测试Redis连接
	if err := cache.Ping(context.Background()); err != nil {
		log.Fatalf("Redis health check failed: %v", err)
	}

	logger.Info("🎉 Redis Context增强功能已就绪！所有函数都支持Context参数")

	// 应用程序启动完成，可以开始处理业务逻辑
	// TODO: 添加HTTP服务器启动、API路由等业务逻辑
}
