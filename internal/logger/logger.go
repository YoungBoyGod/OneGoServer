package logger

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

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
)

// Init 初始化日志系统
func Init(cfg *config.Config) error {
	// 生成应用实例信息
	ProcessID = os.Getpid()
	SessionID = generateSessionID()
	MacAddress = getMacAddress()

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

// GetInstanceInfo 获取实例信息
func GetInstanceInfo() (int, string, string) {
	return ProcessID, SessionID, MacAddress
}

// GetInstanceInfoLegacy 获取实例信息(兼容旧版本)
func GetInstanceInfoLegacy() (int, string) {
	return ProcessID, SessionID
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
	return os.MkdirAll(archiveDir, 0755)
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

// Cleanup 清理资源
func Cleanup() {
	if Logger != nil {
		Logger.Sync()
	}
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
