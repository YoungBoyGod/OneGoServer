package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// === 基础类型定义 ===

// LoggerType 日志类型
type LoggerType string

// 日志类型常量
const (
	LoggerTypeApp   LoggerType = "app"   // 应用主日志
	LoggerTypeHTTP  LoggerType = "http"  // HTTP访问日志
	LoggerTypeError LoggerType = "error" // 错误日志
)

// LoggerStats 日志统计信息
type LoggerStats struct {
	TotalLogs   int64     `json:"total_logs"`
	ErrorLogs   int64     `json:"error_logs"`
	WarnLogs    int64     `json:"warn_logs"`
	InfoLogs    int64     `json:"info_logs"`
	DebugLogs   int64     `json:"debug_logs"`
	LastLogTime time.Time `json:"last_log_time"`
	StartTime   time.Time `json:"start_time"`
}

// ContextLogger 上下文日志器
type ContextLogger struct {
	logger *zap.Logger
	fields []zap.Field
}

// === 全局变量 ===

var (
	// 日志器实例
	appLogger   *zap.Logger // 应用主日志
	httpLogger  *zap.Logger // HTTP访问日志
	errorLogger *zap.Logger // 错误日志

	// 统计和配置
	stats      = &LoggerStats{StartTime: time.Now()}
	statsMu    sync.RWMutex
	loggerMu   sync.RWMutex
	currentCfg *config.LoggingConfig
)

// === 配置验证和标准化 ===

// validateConfig 验证配置参数
func validateConfig(cfg *config.LoggingConfig) error {
	if cfg == nil {
		return fmt.Errorf("logging config cannot be nil")
	}

	if cfg.MaxSize <= 0 {
		return fmt.Errorf("max_size must be positive, got %d", cfg.MaxSize)
	}

	if cfg.MaxBackups < 0 {
		return fmt.Errorf("max_backups must be non-negative, got %d", cfg.MaxBackups)
	}

	if cfg.MaxAge < 0 {
		return fmt.Errorf("max_age must be non-negative, got %d", cfg.MaxAge)
	}

	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[cfg.Level] {
		return fmt.Errorf("invalid log level: %s, must be one of: debug, info, warn, error", cfg.Level)
	}

	validFormats := map[string]bool{
		"json": true, "console": true,
	}
	if !validFormats[cfg.Format] {
		return fmt.Errorf("invalid log format: %s, must be one of: json, console", cfg.Format)
	}

	validOutputs := map[string]bool{
		"stdout": true, "file": true, "both": true,
	}
	if !validOutputs[cfg.Output] {
		return fmt.Errorf("invalid output type: %s, must be one of: stdout, file, both", cfg.Output)
	}

	return nil
}

// normalizeConfig 标准化配置，补充默认值
func normalizeConfig(cfg *config.LoggingConfig) *config.LoggingConfig {
	normalized := *cfg

	if normalized.Level == "" {
		normalized.Level = "info"
	}
	if normalized.Format == "" {
		normalized.Format = "json"
	}
	if normalized.Output == "" {
		normalized.Output = "stdout"
	}
	if normalized.FilePath == "" {
		normalized.FilePath = "logs/app.log"
	}
	if normalized.MaxSize <= 0 {
		normalized.MaxSize = 100
	}
	if normalized.MaxBackups < 0 {
		normalized.MaxBackups = 3
	}
	if normalized.MaxAge < 0 {
		normalized.MaxAge = 30
	}

	return &normalized
}

// === 初始化函数 ===

// InitLogger 初始化所有日志器（保持向后兼容）
func InitLogger(cfg *config.LoggingConfig) error {
	var err error

	// 初始化应用主日志
	appLogger, err = createLogger(cfg, LoggerTypeApp)
	if err != nil {
		return fmt.Errorf("failed to create app logger: %w", err)
	}

	// 初始化HTTP访问日志
	httpLogger, err = createLogger(cfg, LoggerTypeHTTP)
	if err != nil {
		return fmt.Errorf("failed to create http logger: %w", err)
	}

	// 初始化错误日志
	errorLogger, err = createLogger(cfg, LoggerTypeError)
	if err != nil {
		return fmt.Errorf("failed to create error logger: %w", err)
	}

	return nil
}

// InitLoggerEnhanced 增强版日志初始化
func InitLoggerEnhanced(cfg *config.LoggingConfig) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	// 验证配置
	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 标准化配置
	normalizedCfg := normalizeConfig(cfg)

	// 尝试创建logger，失败时使用重试机制
	if err := createLoggersWithRetry(normalizedCfg, 3); err != nil {
		return fmt.Errorf("failed to initialize loggers: %w", err)
	}

	// 保存当前配置
	currentCfg = normalizedCfg

	// 重置统计信息
	resetStats()

	LogInfo("Logger system initialized successfully",
		zap.String("level", normalizedCfg.Level),
		zap.String("format", normalizedCfg.Format),
		zap.String("output", normalizedCfg.Output))

	return nil
}

// createLoggersWithRetry 带重试机制的logger创建
func createLoggersWithRetry(cfg *config.LoggingConfig, maxRetries int) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := InitLogger(cfg); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(time.Millisecond * time.Duration((i+1)*100)) // 递增延迟
				continue
			}
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// === 文件处理函数 ===

// generateTimestampLogPath 生成带时间戳的日志文件路径
func generateTimestampLogPath(originPath string, logType LoggerType) string {
	// 解析原始路径
	dir := filepath.Dir(originPath)
	ext := filepath.Ext(originPath)
	// 获取时间戳
	timestamp := time.Now().Format("20060102150405")
	// 根据日志类型，获取后缀
	var filename string
	switch logType {
	case LoggerTypeApp:
		filename = fmt.Sprintf("app_%s%s", timestamp, ext)
	case LoggerTypeHTTP:
		filename = fmt.Sprintf("http_%s%s", timestamp, ext)
	case LoggerTypeError:
		filename = fmt.Sprintf("error_%s%s", timestamp, ext)
	default:
		filename = fmt.Sprintf("app_%s%s", timestamp, ext)
	}
	newPath := filepath.Join(dir, filename)
	return newPath
}

// getFileWriter 获取文件写入器，支持日志轮转
func getFileWriter(cfg *config.LoggingConfig, logType LoggerType) (io.Writer, error) {
	// 生成带时间戳的日志文件路径
	timestampFilePath := generateTimestampLogPath(cfg.FilePath, logType)

	// 确保日志目录存在
	logDir := filepath.Dir(timestampFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 配置lumberjack进行日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   timestampFilePath,
		MaxSize:    cfg.MaxSize,    // MB
		MaxBackups: cfg.MaxBackups, // 保留的旧文件数量
		MaxAge:     cfg.MaxAge,     // 天数
		Compress:   cfg.Compress,   // 是否压缩
		LocalTime:  true,           // 使用本地时间
	}

	return lumberJackLogger, nil
}

// createLogger 创建指定级别的logger
func createLogger(cfg *config.LoggingConfig, logType LoggerType) (*zap.Logger, error) {
	// 创建zap core
	var zapLevel zapcore.Level
	switch cfg.Level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置输出
	var writers []zapcore.WriteSyncer
	switch cfg.Output {
	case "stdout":
		writers = append(writers, zapcore.AddSync(os.Stdout))
	case "file":
		fileWriter, err := getFileWriter(cfg, logType)
		if err != nil {
			return nil, err
		}
		writers = append(writers, zapcore.AddSync(fileWriter))
	case "both":
		fileWriter, err := getFileWriter(cfg, logType)
		if err != nil {
			return nil, err
		}
		writers = append(writers, zapcore.AddSync(os.Stdout))
		writers = append(writers, zapcore.AddSync(fileWriter))
	default:
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 针对错误日志，只记录错误级别以上的日志
	if logType == LoggerTypeError {
		zapLevel = zapcore.ErrorLevel
	}

	// 创建core
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		zapLevel,
	)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// === 获取日志器实例 ===

// GetAppLogger 获取应用主日志器实例
func GetAppLogger(cfg *config.LoggingConfig) *zap.Logger {
	if appLogger == nil {
		// 使用传入的配置初始化日志器
		if err := InitLogger(cfg); err != nil {
			// 如果初始化失败，使用默认配置
			defaultConfig := &config.LoggingConfig{
				Level:      "info",
				Format:     "json",
				Output:     "stdout",
				FilePath:   "logs/app.log",
				MaxSize:    100,
				MaxBackups: 3,
				MaxAge:     30,
				Compress:   true,
			}
			InitLogger(defaultConfig)
		}
	}
	return appLogger
}

// GetHTTPLogger 获取HTTP访问日志器实例
func GetHTTPLogger() *zap.Logger {
	return httpLogger
}

// GetErrorLogger 获取错误日志器实例
func GetErrorLogger() *zap.Logger {
	return errorLogger
}

// === 基础日志记录函数 ===

// LogInfo 记录应用信息级别日志
func LogInfo(msg string, fields ...zap.Field) {
	if appLogger != nil {
		appLogger.Info(msg, fields...)
	}
}

// LogWarn 记录应用警告级别日志
func LogWarn(msg string, fields ...zap.Field) {
	if appLogger != nil {
		appLogger.Warn(msg, fields...)
	}
}

// LogError 记录应用错误级别日志
func LogError(msg string, fields ...zap.Field) {
	if appLogger != nil {
		appLogger.Error(msg, fields...)
		// 同时记录到错误日志文件
		if errorLogger != nil {
			errorLogger.Error(msg, fields...)
		}
	}
}

// LogDebug 记录应用调试级别日志
func LogDebug(msg string, fields ...zap.Field) {
	if appLogger != nil {
		appLogger.Debug(msg, fields...)
	}
}

// Sync 刷新所有日志缓冲区
func Sync() {
	if appLogger != nil {
		_ = appLogger.Sync()
	}
	if httpLogger != nil {
		_ = httpLogger.Sync()
	}
	if errorLogger != nil {
		_ = errorLogger.Sync()
	}
}

// === 统计功能 ===

// recordLogEvent 记录日志事件
func recordLogEvent(level zapcore.Level) {
	atomic.AddInt64(&stats.TotalLogs, 1)

	switch level {
	case zapcore.DebugLevel:
		atomic.AddInt64(&stats.DebugLogs, 1)
	case zapcore.InfoLevel:
		atomic.AddInt64(&stats.InfoLogs, 1)
	case zapcore.WarnLevel:
		atomic.AddInt64(&stats.WarnLogs, 1)
	case zapcore.ErrorLevel:
		atomic.AddInt64(&stats.ErrorLogs, 1)
	}

	statsMu.Lock()
	stats.LastLogTime = time.Now()
	statsMu.Unlock()
}

// GetStats 获取日志统计信息
func GetStats() LoggerStats {
	statsMu.RLock()
	defer statsMu.RUnlock()
	return LoggerStats{
		TotalLogs:   atomic.LoadInt64(&stats.TotalLogs),
		ErrorLogs:   atomic.LoadInt64(&stats.ErrorLogs),
		WarnLogs:    atomic.LoadInt64(&stats.WarnLogs),
		InfoLogs:    atomic.LoadInt64(&stats.InfoLogs),
		DebugLogs:   atomic.LoadInt64(&stats.DebugLogs),
		LastLogTime: stats.LastLogTime,
		StartTime:   stats.StartTime,
	}
}

// resetStats 重置统计信息
func resetStats() {
	atomic.StoreInt64(&stats.TotalLogs, 0)
	atomic.StoreInt64(&stats.ErrorLogs, 0)
	atomic.StoreInt64(&stats.WarnLogs, 0)
	atomic.StoreInt64(&stats.InfoLogs, 0)
	atomic.StoreInt64(&stats.DebugLogs, 0)

	statsMu.Lock()
	stats.StartTime = time.Now()
	stats.LastLogTime = time.Time{}
	statsMu.Unlock()
}

// === 增强的日志函数 ===

// LogInfoWithStats 记录信息日志并更新统计
func LogInfoWithStats(msg string, fields ...zap.Field) {
	LogInfo(msg, fields...)
	recordLogEvent(zapcore.InfoLevel)
}

// LogWarnWithStats 记录警告日志并更新统计
func LogWarnWithStats(msg string, fields ...zap.Field) {
	LogWarn(msg, fields...)
	recordLogEvent(zapcore.WarnLevel)
}

// LogErrorWithStats 记录错误日志并更新统计
func LogErrorWithStats(msg string, fields ...zap.Field) {
	LogError(msg, fields...)
	recordLogEvent(zapcore.ErrorLevel)
}

// LogDebugWithStats 记录调试日志并更新统计
func LogDebugWithStats(msg string, fields ...zap.Field) {
	LogDebug(msg, fields...)
	recordLogEvent(zapcore.DebugLevel)
}

// === 结构化日志助手 ===

// LogUser 记录用户相关操作
func LogUser(action string, userID int64, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("type", "user_action"),
		zap.String("action", action),
		zap.Int64("user_id", userID),
		zap.Time("timestamp", time.Now()),
	}, fields...)

	LogInfoWithStats("User action", allFields...)
}

// LogDBOperation 记录数据库操作
func LogDBOperation(operation, table string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("type", "db_operation"),
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		LogErrorWithStats("Database operation failed", fields...)
	} else {
		LogInfoWithStats("Database operation completed", fields...)
	}
}

// LogRedisOperation 记录Redis操作
func LogRedisOperation(operation, key string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("type", "redis_operation"),
		zap.String("operation", operation),
		zap.String("key", key),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		LogErrorWithStats("Redis operation failed", fields...)
	} else {
		LogInfoWithStats("Redis operation completed", fields...)
	}
}

// LogKafkaOperation 记录Kafka操作
func LogKafkaOperation(operation, topic string, duration time.Duration, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("type", "kafka_operation"),
		zap.String("operation", operation),
		zap.String("topic", topic),
		zap.Duration("duration", duration),
	}

	// 合并额外字段
	allFields := append(baseFields, fields...)

	LogInfoWithStats("Kafka operation completed", allFields...)
}

// LogAPICall 记录API调用
func LogAPICall(method, endpoint string, statusCode int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("type", "api_call"),
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		LogErrorWithStats("API call failed", fields...)
	} else if statusCode >= 400 {
		LogWarnWithStats("API call returned error status", fields...)
	} else {
		LogInfoWithStats("API call completed", fields...)
	}
}

// LogSystemEvent 记录系统事件
func LogSystemEvent(event string, component string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("type", "system_event"),
		zap.String("event", event),
		zap.String("component", component),
		zap.Time("timestamp", time.Now()),
	}, fields...)

	LogInfoWithStats("System event", allFields...)
}

// === 上下文日志 ===

// WithContext 创建带上下文的日志器
func WithContext(fields ...zap.Field) *ContextLogger {
	return &ContextLogger{
		logger: GetAppLogger(currentCfg),
		fields: fields,
	}
}

// Info 记录信息级别日志
func (cl *ContextLogger) Info(msg string, fields ...zap.Field) {
	if cl.logger == nil {
		return
	}
	allFields := append(cl.fields, fields...)
	cl.logger.Info(msg, allFields...)
	recordLogEvent(zapcore.InfoLevel)
}

// Warn 记录警告级别日志
func (cl *ContextLogger) Warn(msg string, fields ...zap.Field) {
	if cl.logger == nil {
		return
	}
	allFields := append(cl.fields, fields...)
	cl.logger.Warn(msg, allFields...)
	recordLogEvent(zapcore.WarnLevel)
}

// Error 记录错误级别日志
func (cl *ContextLogger) Error(msg string, fields ...zap.Field) {
	if cl.logger == nil {
		return
	}
	allFields := append(cl.fields, fields...)
	cl.logger.Error(msg, allFields...)

	// 同时记录到错误日志
	if errorLogger := GetErrorLogger(); errorLogger != nil {
		errorLogger.Error(msg, allFields...)
	}
	recordLogEvent(zapcore.ErrorLevel)
}

// Debug 记录调试级别日志
func (cl *ContextLogger) Debug(msg string, fields ...zap.Field) {
	if cl.logger == nil {
		return
	}
	allFields := append(cl.fields, fields...)
	cl.logger.Debug(msg, allFields...)
	recordLogEvent(zapcore.DebugLevel)
}

// === 配置管理 ===

// ReloadConfig 热重载配置
func ReloadConfig(newCfg *config.LoggingConfig) error {
	if err := validateConfig(newCfg); err != nil {
		return fmt.Errorf("invalid new config: %w", err)
	}

	LogInfo("Reloading logger configuration...")

	if err := InitLoggerEnhanced(newCfg); err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	LogInfo("Logger configuration reloaded successfully")
	return nil
}

// GetCurrentConfig 获取当前配置
func GetCurrentConfig() *config.LoggingConfig {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	if currentCfg == nil {
		return nil
	}
	// 返回副本避免外部修改
	cfg := *currentCfg
	return &cfg
}

// === 健康检查 ===

// HealthCheck 执行日志系统健康检查
func HealthCheck() error {
	// 检查logger是否初始化
	if GetAppLogger(currentCfg) == nil {
		return fmt.Errorf("app logger not initialized")
	}

	if GetHTTPLogger() == nil {
		return fmt.Errorf("http logger not initialized")
	}

	if GetErrorLogger() == nil {
		return fmt.Errorf("error logger not initialized")
	}

	// 检查错误率
	stats := GetStats()
	if stats.TotalLogs > 100 {
		errorRate := float64(stats.ErrorLogs) / float64(stats.TotalLogs)
		if errorRate > 0.1 { // 错误率超过10%
			return fmt.Errorf("high error rate detected: %.2f%% (%d/%d)",
				errorRate*100, stats.ErrorLogs, stats.TotalLogs)
		}
	}

	// 测试写入
	testMsg := fmt.Sprintf("Health check test at %s", time.Now().Format(time.RFC3339))
	LogDebug(testMsg)

	return nil
}

// GetLoggerStatus 获取logger状态信息
func GetLoggerStatus() map[string]interface{} {
	stats := GetStats()
	cfg := GetCurrentConfig()

	status := map[string]interface{}{
		"initialized": GetAppLogger(currentCfg) != nil,
		"config":      cfg,
		"stats":       stats,
		"uptime":      time.Since(stats.StartTime).String(),
	}

	if stats.TotalLogs > 0 {
		status["error_rate"] = float64(stats.ErrorLogs) / float64(stats.TotalLogs)
		status["logs_per_second"] = float64(stats.TotalLogs) / time.Since(stats.StartTime).Seconds()
	}

	return status
}

// === Gin 中间件 ===

// Logger Gin框架的HTTP日志中间件
func Logger() gin.HandlerFunc {
	// 如果logger未初始化，使用默认配置
	if httpLogger == nil {
		defaultConfig := &config.LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "logs/http.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}
		InitLogger(defaultConfig)
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()
		// 访问路径
		accessPath := c.Request.URL.Path
		// 请求方法
		accessMethod := c.Request.Method
		// 请求参数
		raw := c.Request.URL.RawQuery
		// 请求头
		userAgent := c.Request.UserAgent()
		// 请求IP
		accessIP := c.ClientIP()
		realIP := c.GetHeader("X-Real-IP")

		if raw != "" {
			accessPath = accessPath + "?" + raw
		}

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(startTime)
		// 请求状态码
		accessStatus := c.Writer.Status()
		// 请求体大小
		bodySize := c.Writer.Size()
		// 获取错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 记录HTTP访问日志
		httpLogger.Info("HTTP请求",
			zap.Int("status", accessStatus),
			zap.String("method", accessMethod),
			zap.String("path", accessPath),
			zap.String("accessIP", accessIP),
			zap.Duration("latency", latency),
			zap.Int("size", bodySize),
			zap.String("user_agent", userAgent),
			zap.String("error", errorMessage),
			zap.String("realIP", realIP),
		)

		// 如果是错误状态码，同时记录到错误日志
		if accessStatus >= 400 && errorLogger != nil {
			errorLogger.Error("HTTP错误",
				zap.Int("status", accessStatus),
				zap.String("method", accessMethod),
				zap.String("path", accessPath),
				zap.String("accessIP", accessIP),
				zap.Duration("latency", latency),
				zap.String("user_agent", userAgent),
				zap.String("error", errorMessage),
				zap.String("realIP", realIP),
			)
		}

		// 记录请求日志
		LogInfo(accessMethod,
			zap.String("path", accessPath),
			zap.Int("status", accessStatus),
			zap.Duration("duration", latency),
			zap.String("raw", raw))
	})
}
