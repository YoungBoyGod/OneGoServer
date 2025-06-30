package log

import (
	"errors"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"go.uber.org/zap"
)

// ExampleUsage 展示增强版logger的使用方法
func ExampleUsage() {
	// 1. 使用增强版初始化
	cfg := &config.LoggingConfig{
		Level:      "info",
		Format:     "json",
		Output:     "both",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	// 初始化日志系统（带验证和重试）
	if err := InitLoggerEnhanced(cfg); err != nil {
		panic(err)
	}

	// 2. 基础日志记录（带统计）
	LogInfoWithStats("Application started")
	LogWarnWithStats("This is a warning message")
	LogErrorWithStats("This is an error message")

	// 3. 结构化日志助手

	// 用户操作日志
	LogUser("login", 12345,
		zap.String("ip", "192.168.1.1"),
		zap.String("user_agent", "Mozilla/5.0"))

	// 数据库操作日志
	LogDBOperation("SELECT", "users", time.Millisecond*150, nil)
	LogDBOperation("INSERT", "orders", time.Millisecond*300,
		errors.New("duplicate key"))

	// API调用日志
	LogAPICall("GET", "/api/users", 200, time.Millisecond*120, nil)
	LogAPICall("POST", "/api/orders", 500, time.Millisecond*200,
		errors.New("internal server error"))

	// 系统事件日志
	LogSystemEvent("service_start", "user-service",
		zap.String("version", "v1.2.3"),
		zap.Int("port", 8080))

	// 4. 上下文日志
	userLogger := WithContext(
		zap.Int64("user_id", 12345),
		zap.String("session_id", "sess_abc123"),
		zap.String("request_id", "req_xyz789"))

	userLogger.Info("User performed action",
		zap.String("action", "create_order"),
		zap.Float64("amount", 99.99))

	userLogger.Error("Operation failed",
		zap.String("operation", "payment"),
		zap.Error(errors.New("payment gateway timeout")))

	// 5. 配置热重载
	newCfg := &config.LoggingConfig{
		Level:      "debug", // 改为debug级别
		Format:     "json",
		Output:     "both",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	if err := ReloadConfig(newCfg); err != nil {
		LogErrorWithStats("Failed to reload config", zap.Error(err))
	}

	// 6. 统计信息
	stats := GetStats()
	LogInfoWithStats("Logger statistics",
		zap.Int64("total_logs", stats.TotalLogs),
		zap.Int64("error_logs", stats.ErrorLogs),
		zap.String("uptime", time.Since(stats.StartTime).String()))

	// 7. 健康检查
	if err := HealthCheck(); err != nil {
		LogErrorWithStats("Logger health check failed", zap.Error(err))
	} else {
		LogInfoWithStats("Logger health check passed")
	}

	// 8. 获取状态信息
	status := GetLoggerStatus()
	LogInfoWithStats("Logger status", zap.Any("status", status))

	// 9. 确保日志写入
	Sync()
}

// ExampleContextualLogging 展示上下文日志的高级用法
func ExampleContextualLogging() {
	// 为特定请求创建上下文logger
	requestLogger := WithContext(
		zap.String("request_id", "req_123456"),
		zap.String("user_id", "user_789"),
		zap.String("ip", "192.168.1.100"),
		zap.Time("request_start", time.Now()))

	// 在整个请求处理过程中使用这个logger
	requestLogger.Info("Request received",
		zap.String("method", "POST"),
		zap.String("path", "/api/orders"))

	// 模拟数据库操作
	start := time.Now()
	// ... 数据库操作 ...
	requestLogger.Info("Database query completed",
		zap.Duration("duration", time.Since(start)),
		zap.String("table", "orders"))

	// 模拟外部API调用
	start = time.Now()
	// ... API调用 ...
	requestLogger.Info("External API call completed",
		zap.Duration("duration", time.Since(start)),
		zap.String("service", "payment-gateway"))

	requestLogger.Info("Request completed",
		zap.Int("status_code", 201),
		zap.Duration("total_duration", time.Since(time.Now())))
}

// ExampleErrorHandling 展示错误处理的最佳实践
func ExampleErrorHandling() {
	// 1. 配置验证错误
	invalidCfg := &config.LoggingConfig{
		Level:      "invalid_level", // 无效级别
		MaxSize:    -1,              // 无效值
		MaxBackups: -5,              // 无效值
	}

	if err := InitLoggerEnhanced(invalidCfg); err != nil {
		// 这里会捕获配置验证错误
		LogErrorWithStats("Invalid configuration", zap.Error(err))

		// 使用默认配置作为备用
		defaultCfg := &config.LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}

		if err := InitLoggerEnhanced(defaultCfg); err != nil {
			panic("Failed to initialize logger with default config: " + err.Error())
		}
	}

	// 2. 记录不同类型的错误

	// 业务逻辑错误
	LogUser("purchase_failed", 12345,
		zap.Error(errors.New("insufficient balance")),
		zap.Float64("attempted_amount", 199.99),
		zap.Float64("current_balance", 50.00))

	// 系统错误
	LogSystemEvent("database_connection_failed", "user-service",
		zap.Error(errors.New("connection timeout")),
		zap.String("database", "postgresql"),
		zap.String("host", "db.example.com"))

	// API错误
	LogAPICall("GET", "/api/external/rates", 503, time.Second*30,
		errors.New("service unavailable"))
}
