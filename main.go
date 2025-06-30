package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// 测试Redis基本操作（Context支持）
	ctx := context.Background()

	// 测试SET操作
	if err := cache.Set(ctx, "test_key", "Hello Redis with Context!", time.Minute*5); err != nil {
		logger.Error("Redis SET failed", zap.Error(err))
	} else {
		logger.Info("Redis SET operation successful")
	}

	// 测试GET操作
	if value, err := cache.Get(ctx, "test_key"); err != nil {
		logger.Error("Redis GET failed", zap.Error(err))
	} else {
		logger.Info("Redis GET operation successful", zap.String("value", value))
	}

	// 测试对象序列化操作
	testObj := map[string]interface{}{
		"message":   "Context增强测试成功",
		"timestamp": time.Now().Unix(),
		"features":  []string{"context支持", "超时控制", "取消机制"},
	}

	if err := cache.SetObject(ctx, "test_object", testObj, time.Minute*5); err != nil {
		logger.Error("Redis SetObject failed", zap.Error(err))
	} else {
		logger.Info("Redis SetObject operation successful")
	}

	// 获取对象
	var retrievedObj map[string]interface{}
	if err := cache.GetObject(ctx, "test_object", &retrievedObj); err != nil {
		logger.Error("Redis GetObject failed", zap.Error(err))
	} else {
		logger.Info("Redis GetObject operation successful", zap.Any("object", retrievedObj))
	}

	// 测试Redis统计信息
	if stats, err := cache.GetRedisStats(ctx); err != nil {
		logger.Error("Redis GetStats failed", zap.Error(err))
	} else {
		logger.Info("Redis statistics",
			zap.Uint32("total_conns", stats.TotalConns),
			zap.Uint32("idle_conns", stats.IdleConns),
			zap.Bool("is_connected", stats.IsConnected))
	}

	logger.Info("🎉 Redis Context增强功能测试完成！")
}
