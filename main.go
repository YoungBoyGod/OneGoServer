package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/router"
	"github.com/YoungBoyGod/OneGoServer/pkg/cache"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/YoungBoyGod/OneGoServer/pkg/queue"
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

	// 增加kafka初始化
	if err := queue.InitKafka(context.Background(), &cfg.Kafka); err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	logger.Info("Kafka initialized successfully")

	// 测试Kafka生产者和消费者功能
	testKafkaProducerConsumer(context.Background(), logger)

	// 应用程序启动完成，可以开始处理业务逻辑
	// TODO: 添加HTTP服务器启动、API路由等业务逻辑
	router := router.InitRouter()
	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))

	// 执行数据库迁移
	// gormDB := sql.GetDB()
	// if gormDB == nil {
	// 	log.Fatalf("Failed to get database connection: database not initialized")
	// }

	// db, err := gormDB.DB()
	// if err != nil {
	// 	log.Fatalf("Failed to get underlying sql.DB: %v", err)
	// }
	// if err := utils.RunMigrations(db); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }
	// logger.Info("Database migrations completed successfully")
}

// testKafkaProducerConsumer 测试Kafka生产者和消费者功能
func testKafkaProducerConsumer(ctx context.Context, logger *zap.Logger) {
	logger.Info("🚀 开始测试Kafka生产者和消费者功能")

	// 定义测试Topic
	testTopic := "onego-test-topic"

	// 测试数据
	testMessages := []struct {
		key   string
		value interface{}
	}{
		{"user_action", map[string]interface{}{
			"user_id":    "12345",
			"action":     "login",
			"timestamp":  time.Now().Unix(),
			"ip_address": "192.168.1.100",
		}},
		{"system_event", map[string]interface{}{
			"event_type": "server_start",
			"server_id":  "onego-server-001",
			"timestamp":  time.Now().Unix(),
			"version":    "v0.1.0",
		}},
		{"task_update", map[string]interface{}{
			"task_id":      "task_001",
			"status":       "completed",
			"user_id":      "12345",
			"completed_at": time.Now().Unix(),
		}},
	}

	// 1. 测试生产者 - 发送消息
	logger.Info("📤 开始测试生产者功能")
	for i, msg := range testMessages {
		kafkaMsg := &queue.KafkaMessage{
			Topic: testTopic,
			Key:   msg.key,
			Value: msg.value,
			Headers: map[string]string{
				"source":       "onego-server",
				"message_id":   fmt.Sprintf("msg_%d_%d", i+1, time.Now().Unix()),
				"content_type": "application/json",
			},
		}

		logger.Info("发送消息",
			zap.String("topic", testTopic),
			zap.String("key", msg.key),
			zap.Any("value", msg.value))

		if err := queue.SendMessage(ctx, kafkaMsg); err != nil {
			logger.Error("发送消息失败", zap.Error(err))
		} else {
			logger.Info("✅ 消息发送成功",
				zap.String("key", msg.key))
		}

		// 短暂延迟，避免消息发送过快
		time.Sleep(100 * time.Millisecond)
	}

	// 2. 等待一下确保消息已发送
	time.Sleep(1 * time.Second)

	// 3. 测试消费者 - 消费消息
	logger.Info("📥 开始测试消费者功能")

	// 定义消息处理函数
	messageHandler := func(ctx context.Context, message *queue.KafkaMessage) error {
		logger.Info("✅ 收到消息",
			zap.String("topic", message.Topic),
			zap.String("key", message.Key),
			zap.Any("value", message.Value),
			zap.Any("headers", message.Headers),
			zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset),
			zap.Time("timestamp", message.Timestamp))

		// 模拟消息处理逻辑
		switch message.Key {
		case "user_action":
			logger.Info("🔐 处理用户行为数据", zap.Any("data", message.Value))
		case "system_event":
			logger.Info("⚙️ 处理系统事件", zap.Any("data", message.Value))
		case "task_update":
			logger.Info("📋 处理任务更新", zap.Any("data", message.Value))
		default:
			logger.Info("📝 处理通用消息", zap.Any("data", message.Value))
		}

		return nil
	}

	// 创建消费者上下文（设置超时避免无限等待）
	consumerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 从最新的offset开始消费（offset设为-1表示最新）
	logger.Info("开始消费消息",
		zap.String("topic", testTopic),
		zap.String("timeout", "5秒"))

	if err := queue.ConsumeMessages(consumerCtx, testTopic, 0, -1, messageHandler); err != nil {
		// 如果是超时错误，这是正常的（消费完成）
		if consumerCtx.Err() == context.DeadlineExceeded {
			logger.Info("⏰ 消费者测试完成（已超时，这是正常的）")
		} else {
			logger.Error("消费消息时出错", zap.Error(err))
		}
	}

	// 4. 获取Kafka统计信息
	if stats, err := queue.GetKafkaStats(ctx); err != nil {
		logger.Error("获取Kafka统计信息失败", zap.Error(err))
	} else {
		logger.Info("📊 Kafka统计信息",
			zap.Int64("发送消息数", stats.ProducerStats.MessagesSent),
			zap.Int64("成功消息数", stats.ProducerStats.MessagesSucceeded),
			zap.Int64("失败消息数", stats.ProducerStats.MessagesFailed),
			zap.Int64("接收消息数", stats.ConsumerStats.MessagesReceived),
			zap.Int64("处理消息数", stats.ConsumerStats.MessagesProcessed),
			zap.Bool("连接状态", stats.IsConnected))
	}

	// 5. 测试获取Topic列表
	if topics, err := queue.GetTopics(ctx); err != nil {
		logger.Error("获取Topic列表失败", zap.Error(err))
	} else {
		logger.Info("📋 当前Topics列表", zap.Strings("topics", topics))
	}

	logger.Info("🎉 Kafka生产者和消费者测试完成！")
}
