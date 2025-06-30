# Logger系统优化计划

## 当前系统分析

### 已有功能
✅ 基础日志记录功能  
✅ 多类型日志器支持  
✅ 配置文件驱动  
✅ 文件轮转支持  
✅ Gin中间件集成  
✅ 双重错误记录  

### 存在的优化机会

## 1. 配置验证与处理优化

### 问题
- 缺少配置参数验证
- 默认值处理不够完善
- 没有配置热重载功能

### 优化方案
```go
// 配置验证器
func validateConfig(cfg *config.LoggingConfig) error {
    if cfg.MaxSize <= 0 {
        return fmt.Errorf("max_size must be positive, got %d", cfg.MaxSize)
    }
    if cfg.MaxBackups < 0 {
        return fmt.Errorf("max_backups must be non-negative, got %d", cfg.MaxBackups)
    }
    // 更多验证...
}

// 配置标准化
func normalizeConfig(cfg *config.LoggingConfig) *config.LoggingConfig {
    normalized := *cfg
    if normalized.Level == "" {
        normalized.Level = "info"
    }
    if normalized.Format == "" {
        normalized.Format = "json"
    }
    return &normalized
}

// 配置热重载
func ReloadConfig(newCfg *config.LoggingConfig) error {
    // 验证新配置
    if err := validateConfig(newCfg); err != nil {
        return err
    }
    
    // 重新初始化logger
    return InitLogger(newCfg)
}
```

## 2. 性能优化

### 问题
- 每次调用都生成时间戳路径
- 缺少日志采样功能
- 同步写入可能影响性能

### 优化方案
```go
// 日志采样配置
type SamplingConfig struct {
    Initial    int           // 初始采样数
    Thereafter int           // 后续采样数
    Tick       time.Duration // 采样周期
}

// 性能监控
type LoggerMetrics struct {
    TotalLogs    int64
    ErrorLogs    int64
    DroppedLogs  int64
    WriteLatency time.Duration
}

// 异步日志写入
func createAsyncLogger(cfg *config.LoggingConfig, logType LoggerType) (*zap.Logger, error) {
    // 使用buffered writer提高性能
    // 支持批量写入
}

// 路径缓存
var pathCache = sync.Map{}

func getCachedLogPath(originPath string, logType LoggerType) string {
    key := fmt.Sprintf("%s_%s", originPath, logType)
    if cached, ok := pathCache.Load(key); ok {
        return cached.(string)
    }
    
    path := generateTimestampLogPath(originPath, logType)
    pathCache.Store(key, path)
    return path
}
```

## 3. 功能扩展

### 3.1 动态日志级别调整
```go
// 运行时调整日志级别
func SetLogLevel(level string) error {
    zapLevel, err := parseLogLevel(level)
    if err != nil {
        return err
    }
    
    // 动态调整所有logger的级别
    if appLogger != nil {
        appLogger = appLogger.WithOptions(zap.IncreaseLevel(zapLevel))
    }
    return nil
}

// 获取当前日志级别
func GetLogLevel() string {
    // 返回当前有效日志级别
}
```

### 3.2 结构化日志助手
```go
// 便捷的结构化日志方法
func LogUser(action string, userID int64, fields ...zap.Field) {
    allFields := append([]zap.Field{
        zap.String("type", "user_action"),
        zap.String("action", action),
        zap.Int64("user_id", userID),
        zap.Time("timestamp", time.Now()),
    }, fields...)
    
    LogInfo("User action", allFields...)
}

func LogDBOperation(operation, table string, duration time.Duration, err error) {
    fields := []zap.Field{
        zap.String("type", "db_operation"),
        zap.String("operation", operation),
        zap.String("table", table),
        zap.Duration("duration", duration),
    }
    
    if err != nil {
        fields = append(fields, zap.Error(err))
        LogError("Database operation failed", fields...)
    } else {
        LogInfo("Database operation completed", fields...)
    }
}

func LogAPICall(method, endpoint string, statusCode int, duration time.Duration) {
    LogInfo("API call",
        zap.String("type", "api_call"),
        zap.String("method", method),
        zap.String("endpoint", endpoint),
        zap.Int("status_code", statusCode),
        zap.Duration("duration", duration),
    )
}
```

### 3.3 上下文日志支持
```go
type ContextLogger struct {
    logger *zap.Logger
    fields []zap.Field
}

func WithContext(fields ...zap.Field) *ContextLogger {
    return &ContextLogger{
        logger: appLogger,
        fields: fields,
    }
}

func (cl *ContextLogger) Info(msg string, fields ...zap.Field) {
    allFields := append(cl.fields, fields...)
    cl.logger.Info(msg, allFields...)
}

func (cl *ContextLogger) Error(msg string, fields ...zap.Field) {
    allFields := append(cl.fields, fields...)
    cl.logger.Error(msg, allFields...)
    if errorLogger != nil {
        errorLogger.Error(msg, allFields...)
    }
}

// 使用示例：
// userLogger := WithContext(zap.Int64("user_id", 12345), zap.String("session", "abc123"))
// userLogger.Info("User logged in")
```

### 3.4 钩子函数支持
```go
type LogHook func(entry zapcore.Entry, fields []zapcore.Field)

var logHooks []LogHook

func AddHook(hook LogHook) {
    logHooks = append(logHooks, hook)
}

func RemoveAllHooks() {
    logHooks = []LogHook{}
}

// 在createLogger中集成钩子
func createLoggerWithHooks(cfg *config.LoggingConfig, logType LoggerType) (*zap.Logger, error) {
    // 创建带钩子的core
    core := &hookCore{
        Core:  zapcore.NewCore(...),
        hooks: logHooks,
    }
    return zap.New(core), nil
}
```

## 4. 错误处理优化

### 4.1 优雅降级
```go
func safeCreateLogger(cfg *config.LoggingConfig, logType LoggerType) *zap.Logger {
    logger, err := createLogger(cfg, logType)
    if err != nil {
        // 降级到最简单的控制台日志
        config := zap.NewDevelopmentConfig()
        fallbackLogger, _ := config.Build()
        LogWarn("Failed to create logger, using fallback", 
               zap.String("type", string(logType)), 
               zap.Error(err))
        return fallbackLogger
    }
    return logger
}
```

### 4.2 错误重试机制
```go
func createLoggerWithRetry(cfg *config.LoggingConfig, logType LoggerType, maxRetries int) (*zap.Logger, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        logger, err := createLogger(cfg, logType)
        if err == nil {
            return logger, nil
        }
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
    }
    return nil, fmt.Errorf("failed to create logger after %d retries: %w", maxRetries, lastErr)
}
```

## 5. 监控和统计

### 5.1 性能监控
```go
type LoggerStats struct {
    mu           sync.RWMutex
    totalLogs    int64
    errorLogs    int64
    writeLatency time.Duration
    lastError    error
    lastErrorAt  time.Time
}

var stats = &LoggerStats{}

func GetStats() LoggerStats {
    stats.mu.RLock()
    defer stats.mu.RUnlock()
    return *stats
}

func recordLogEvent(level zapcore.Level, duration time.Duration) {
    stats.mu.Lock()
    stats.totalLogs++
    if level >= zapcore.ErrorLevel {
        stats.errorLogs++
    }
    stats.writeLatency = duration
    stats.mu.Unlock()
}
```

### 5.2 健康检查
```go
func HealthCheck() error {
    // 检查日志器是否正常
    if appLogger == nil {
        return fmt.Errorf("app logger not initialized")
    }
    
    // 检查文件写入权限
    testMsg := "health check test"
    LogInfo(testMsg)
    
    // 检查错误率
    stats := GetStats()
    if stats.totalLogs > 100 && float64(stats.errorLogs)/float64(stats.totalLogs) > 0.1 {
        return fmt.Errorf("high error rate: %d/%d", stats.errorLogs, stats.totalLogs)
    }
    
    return nil
}
```

## 6. 调试和开发支持

### 6.1 调试模式
```go
func EnableDebugMode() {
    SetLogLevel("debug")
    // 添加更详细的调用信息
    if appLogger != nil {
        appLogger = appLogger.WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))
    }
}

func DisableDebugMode() {
    SetLogLevel("info")
    // 移除详细信息以提高性能
}
```

### 6.2 日志回放功能
```go
func StartLogRecording() {
    // 开始记录日志到内存缓冲区用于回放
}

func StopLogRecording() []LogEntry {
    // 停止记录并返回缓冲的日志
}

func ReplayLogs(entries []LogEntry) {
    // 重放记录的日志
}
```

## 实施优先级

### 高优先级 🔴
1. 配置验证和标准化
2. 错误处理优化
3. 性能基础优化

### 中优先级 🟡
1. 结构化日志助手
2. 动态日志级别
3. 监控统计

### 低优先级 🟢
1. 钩子函数支持
2. 日志回放功能
3. 高级异步特性

## 向后兼容性

所有优化都将保持与现有API的兼容性，新功能通过可选参数或新函数提供。 