package logger

import (
	"crypto/md5"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"learngo0619/internal/config"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// 全局日志实例
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger

	// 应用实例信息
	ProcessID  int
	SessionID  string
	MacAddress string // 服务端MAC地址

	// 客户端日志映射
	clientLoggers map[string]*zap.Logger
	clientMutex   sync.RWMutex

	// 配置实例
	globalConfig *config.Config
)

// ClientLogConfig 客户端日志配置
type ClientLogConfig struct {
	MaxClients    int  `yaml:"max_clients"`    // 最大客户端数量
	EnableCleanup bool `yaml:"enable_cleanup"` // 是否启用自动清理
	CleanupHours  int  `yaml:"cleanup_hours"`  // 清理时间(小时)
}

// 默认客户端日志配置
var defaultClientLogConfig = ClientLogConfig{
	MaxClients:    50,   // 最多50个客户端
	EnableCleanup: true, // 启用自动清理
	CleanupHours:  24,   // 24小时未使用的日志文件清理
}

// Init 初始化日志系统
func Init(cfg *config.Config) error {
	// 生成应用实例信息
	ProcessID = os.Getpid()
	SessionID = generateSessionID()
	MacAddress = getMacAddress()

	// 保存配置
	globalConfig = cfg

	// 初始化客户端日志映射
	clientLoggers = make(map[string]*zap.Logger)

	logger, err := NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 设置全局实例
	Logger = logger
	Sugar = logger.Sugar()

	// 替换zap的全局logger
	zap.ReplaceGlobals(logger)

	return nil
}

// generateSessionID 生成唯一会话ID
func generateSessionID() string {
	return uuid.New().String()
}

// getMacAddress 获取服务端MAC地址
func getMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}

	for _, iface := range interfaces {
		// 跳过回环接口和非激活接口
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			mac := iface.HardwareAddr.String()
			if mac != "" {
				return mac
			}
		}
	}
	return "unknown"
}

// generateClientID 生成客户端标识符
func generateClientID(ip, userAgent string) string {
	// 使用IP和User-Agent的MD5哈希作为客户端ID
	data := fmt.Sprintf("%s|%s", ip, userAgent)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:8] // 取前8位
}

// getClientFileName 获取客户端日志文件名
func getClientFileName(clientID, ip string) string {
	// 文件名格式: client_<clientID>_<safeIP>.log
	safeIP := strings.ReplaceAll(ip, ".", "_")
	safeIP = strings.ReplaceAll(safeIP, ":", "_")
	return fmt.Sprintf("client_%s_%s.log", clientID, safeIP)
}

// GetOrCreateClientLogger 获取或创建客户端专用日志器
func GetOrCreateClientLogger(ip, userAgent string) *zap.Logger {
	clientID := generateClientID(ip, userAgent)

	clientMutex.RLock()
	if logger, exists := clientLoggers[clientID]; exists {
		clientMutex.RUnlock()
		return logger
	}
	clientMutex.RUnlock()

	// 需要创建新的客户端日志器
	clientMutex.Lock()
	defer clientMutex.Unlock()

	// 双重检查
	if logger, exists := clientLoggers[clientID]; exists {
		return logger
	}

	// 检查客户端数量限制
	if len(clientLoggers) >= defaultClientLogConfig.MaxClients {
		// 返回主日志器而不是创建新的
		return Logger
	}

	// 创建客户端日志器
	clientLogger := createClientLogger(clientID, ip, userAgent)
	clientLoggers[clientID] = clientLogger

	// 记录新客户端创建
	Logger.Info("Created new client logger",
		zap.String("client_id", clientID),
		zap.String("client_ip", ip),
		zap.String("user_agent", userAgent),
		zap.Int("total_clients", len(clientLoggers)),
	)

	return clientLogger
}

// createClientLogger 创建客户端专用日志器
func createClientLogger(clientID, ip, userAgent string) *zap.Logger {
	if globalConfig == nil {
		return Logger
	}

	// 确保客户端日志目录存在
	clientLogDir := filepath.Join(globalConfig.Log.Dir, "clients")
	os.MkdirAll(clientLogDir, 0755)

	// 客户端日志文件
	fileName := getClientFileName(clientID, ip)
	clientLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(clientLogDir, fileName),
		MaxSize:    globalConfig.Log.MaxSize / 2, // 客户端日志文件稍小
		MaxAge:     globalConfig.Log.MaxAge,
		MaxBackups: globalConfig.Log.MaxBackups / 2,
		Compress:   globalConfig.Log.Compress,
	}

	// 创建编码器
	encoder := getEncoder(globalConfig.Log.Format)
	level := getLogLevel(globalConfig.Log.Level)

	// 创建核心
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(clientLogWriter),
		level,
	)

	// 创建日志器并添加客户端特定字段
	clientLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	clientLogger = clientLogger.With(
		zap.String("app", globalConfig.App.Name),
		zap.String("version", globalConfig.App.Version),
		zap.String("env", config.GetEnvironment()),
		zap.Int("pid", ProcessID),
		zap.String("session_id", SessionID),
		zap.String("server_mac", MacAddress),
		zap.String("client_id", clientID),
		zap.String("client_ip", ip),
		zap.String("client_ua_hash", fmt.Sprintf("%x", md5.Sum([]byte(userAgent)))[:8]),
	)

	return clientLogger
}

// GetInstanceInfo 获取实例信息
func GetInstanceInfo() (int, string, string) {
	return ProcessID, SessionID, MacAddress
}

// GetInstanceInfoLegacy 获取实例信息(兼容旧版本)
func GetInstanceInfoLegacy() (int, string) {
	return ProcessID, SessionID
}

// GetClientLoggerStats 获取客户端日志器统计信息
func GetClientLoggerStats() map[string]interface{} {
	clientMutex.RLock()
	defer clientMutex.RUnlock()

	return map[string]interface{}{
		"total_clients":   len(clientLoggers),
		"max_clients":     defaultClientLogConfig.MaxClients,
		"cleanup_enabled": defaultClientLogConfig.EnableCleanup,
	}
}

// LogForClient 为特定客户端记录日志
func LogForClient(ip, userAgent, msg string, level zapcore.Level, fields ...zap.Field) {
	// 主日志记录
	switch level {
	case zapcore.DebugLevel:
		Logger.Debug(msg, fields...)
	case zapcore.InfoLevel:
		Logger.Info(msg, fields...)
	case zapcore.WarnLevel:
		Logger.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		Logger.Error(msg, fields...)
	case zapcore.FatalLevel:
		Logger.Fatal(msg, fields...)
	}

	// 客户端专用日志记录
	clientLogger := GetOrCreateClientLogger(ip, userAgent)
	if clientLogger != Logger { // 避免重复记录
		switch level {
		case zapcore.DebugLevel:
			clientLogger.Debug(msg, fields...)
		case zapcore.InfoLevel:
			clientLogger.Info(msg, fields...)
		case zapcore.WarnLevel:
			clientLogger.Warn(msg, fields...)
		case zapcore.ErrorLevel:
			clientLogger.Error(msg, fields...)
		case zapcore.FatalLevel:
			clientLogger.Fatal(msg, fields...)
		}
	}
}

// 便捷方法 - 客户端特定日志
func InfoForClient(ip, userAgent, msg string, fields ...zap.Field) {
	LogForClient(ip, userAgent, msg, zapcore.InfoLevel, fields...)
}

func WarnForClient(ip, userAgent, msg string, fields ...zap.Field) {
	LogForClient(ip, userAgent, msg, zapcore.WarnLevel, fields...)
}

func ErrorForClient(ip, userAgent, msg string, fields ...zap.Field) {
	LogForClient(ip, userAgent, msg, zapcore.ErrorLevel, fields...)
}

func DebugForClient(ip, userAgent, msg string, fields ...zap.Field) {
	LogForClient(ip, userAgent, msg, zapcore.DebugLevel, fields...)
}

// NewLogger 创建新的日志实例
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	// 确保日志目录存在
	if err := ensureLogDir(cfg.Log.Dir); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 配置输出
	writers := getWriters(cfg)

	// 创建core
	core := zapcore.NewTee(writers...)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 添加字段 (包含PID、会话ID和MAC地址)
	logger = logger.With(
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", config.GetEnvironment()),
		zap.Int("pid", ProcessID),
		zap.String("session_id", SessionID),
		zap.String("server_mac", MacAddress), // 服务端MAC地址
	)

	return logger, nil
}

// ensureLogDir 确保日志目录存在
func ensureLogDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// 创建归档目录
	archiveDir := filepath.Join(dir, "archives")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return err
	}
	// 创建客户端日志目录
	clientDir := filepath.Join(dir, "clients")
	return os.MkdirAll(clientDir, 0755)
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoder 获取编码器
func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch strings.ToLower(format) {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// getWriters 获取输出写入器
func getWriters(cfg *config.Config) []zapcore.Core {
	level := getLogLevel(cfg.Log.Level)
	encoder := getEncoder(cfg.Log.Format)

	var cores []zapcore.Core

	switch strings.ToLower(cfg.Log.Output) {
	case "stdout":
		// 仅控制台输出
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))

	case "file":
		// 仅文件输出
		cores = append(cores, getFileCores(cfg, level)...)

	case "both":
		// 控制台和文件双输出
		consoleEncoder := getEncoder("console")
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
		cores = append(cores, getFileCores(cfg, level)...)

	default:
		// 默认输出到控制台
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	return cores
}

// getFileCores 获取文件输出cores
func getFileCores(cfg *config.Config, level zapcore.Level) []zapcore.Core {
	var cores []zapcore.Core

	// 通用日志文件 (所有级别)
	allLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, "app.log"),
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}

	// 错误日志文件 (仅错误级别以上)
	errorLogWriter := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Log.Dir, "error.log"),
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}

	encoder := getEncoder(cfg.Log.Format)

	// 所有日志文件
	cores = append(cores, zapcore.NewCore(
		encoder,
		zapcore.AddSync(allLogWriter),
		level,
	))

	// 错误日志文件（仅错误级别）
	cores = append(cores, zapcore.NewCore(
		encoder,
		zapcore.AddSync(errorLogWriter),
		zapcore.ErrorLevel,
	))

	return cores
}

// CleanupInactiveClientLoggers 清理非活跃的客户端日志器
func CleanupInactiveClientLoggers() int {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	cleaned := 0
	for clientID, clientLogger := range clientLoggers {
		// 同步并清理客户端日志器
		if clientLogger != nil {
			clientLogger.Sync()
			delete(clientLoggers, clientID)
			cleaned++
		}
	}

	Logger.Info("Cleaned up inactive client loggers",
		zap.Int("cleaned_count", cleaned),
		zap.Int("remaining_clients", len(clientLoggers)),
	)

	return cleaned
}

// Cleanup 清理资源
func Cleanup() {
	if Logger != nil {
		Logger.Sync()
	}

	// 清理所有客户端日志器
	CleanupInactiveClientLoggers()
}

// 便捷方法
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Fatal(msg, fields...)
	}
}

// Sugar便捷方法
func Debugf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Debugf(template, args...)
	}
}

func Infof(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Infof(template, args...)
	}
}

func Warnf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Warnf(template, args...)
	}
}

func Errorf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Errorf(template, args...)
	}
}

func Fatalf(template string, args ...interface{}) {
	if Sugar != nil {
		Sugar.Fatalf(template, args...)
	}
}
