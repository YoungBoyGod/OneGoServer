package log

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

var (
	stats      = &LoggerStats{StartTime: time.Now()}
	statsMu    sync.RWMutex
	loggerMu   sync.RWMutex
	currentCfg *config.LoggingConfig
)

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
