package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义log的接口

var (
	appLogger   *zap.Logger // 应用主日志
	httpLogger  *zap.Logger // HTTP访问日志
	errorLogger *zap.Logger // 错误日志
)

// LoggerType 日志类型
type LoggerType string

const (
	LoggerTypeApp   LoggerType = "app"   // 应用主日志
	LoggerTypeHTTP  LoggerType = "http"  // HTTP访问日志
	LoggerTypeError LoggerType = "error" // 错误日志
)

// InitLogger 初始化所有日志器
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

// 生成带时间戳的日志文件路径
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

// 创建指定级别的logger
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
		// 请求状态码
		accessStatus := c.Writer.Status()
		// 请求体大小
		bodySize := c.Writer.Size()

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
		LogInfo(accessMethod, zap.String("path", accessPath), zap.Int("status", accessStatus), zap.Duration("duration", latency), zap.String("raw", raw))
	})
}
