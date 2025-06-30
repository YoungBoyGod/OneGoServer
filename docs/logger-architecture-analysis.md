# Logger架构分析与函数调用关系

## 概述
OneGoServer的日志系统基于zap构建，支持多种日志类型、灵活的输出配置和自动轮转功能。

## 核心组件架构

### 1. 全局变量
```go
var (
    appLogger   *zap.Logger // 应用主日志
    httpLogger  *zap.Logger // HTTP访问日志  
    errorLogger *zap.Logger // 错误日志
)
```

### 2. 日志类型定义
```go
type LoggerType string

const (
    LoggerTypeApp   LoggerType = "app"   // 应用主日志
    LoggerTypeHTTP  LoggerType = "http"  // HTTP访问日志
    LoggerTypeError LoggerType = "error" // 错误日志
)
```

## 函数调用关系图

### 初始化阶段
```
main.go: InitLogger(cfg)
    ├── createLogger(cfg, LoggerTypeApp) → appLogger
    ├── createLogger(cfg, LoggerTypeHTTP) → httpLogger
    └── createLogger(cfg, LoggerTypeError) → errorLogger
```

### createLogger函数调用链
```
createLogger(cfg, logType)
    ├── 解析日志级别 (cfg.Level → zapcore.Level)
    ├── 配置编码器 (EncoderConfig)
    ├── 配置输出方式
    │   ├── stdout → zapcore.AddSync(os.Stdout)
    │   ├── file → getFileWriter(cfg, logType)
    │   └── both → stdout + file
    ├── getFileWriter(cfg, logType)
    │   ├── generateTimestampLogPath(cfg.FilePath, logType)
    │   ├── os.MkdirAll(logDir, 0755)
    │   └── &lumberjack.Logger{...}
    ├── zapcore.NewCore(encoder, writers, level)
    └── zap.New(core, options...)
```

## 详细函数分析

### 1. 初始化函数

#### InitLogger(cfg *config.LoggingConfig) error
**作用**：系统启动时初始化所有日志器
**调用顺序**：
1. `createLogger(cfg, LoggerTypeApp)` → 设置appLogger
2. `createLogger(cfg, LoggerTypeHTTP)` → 设置httpLogger  
3. `createLogger(cfg, LoggerTypeError)` → 设置errorLogger

**错误处理**：任何一个logger创建失败都会返回错误

#### createLogger(cfg, logType) (*zap.Logger, error)
**作用**：根据配置和类型创建具体的zap.Logger实例
**处理流程**：
1. **级别解析**：cfg.Level → zapcore.Level
2. **编码器配置**：根据cfg.Format选择JSON或Console编码器
3. **输出配置**：根据cfg.Output配置writer
4. **特殊处理**：ErrorLogger强制使用ErrorLevel
5. **Core创建**：zapcore.NewCore(encoder, writers, level)
6. **Logger创建**：zap.New(core, zap.AddCaller(), zap.AddStacktrace())

### 2. 辅助函数

#### generateTimestampLogPath(originPath, logType) string
**作用**：为不同类型的日志生成带时间戳的文件路径
**命名规则**：
- App日志：`app_20060102150405.log`
- HTTP日志：`http_20060102150405.log`
- Error日志：`error_20060102150405.log`

#### getFileWriter(cfg, logType) (io.Writer, error)
**作用**：创建支持轮转的文件写入器
**配置参数**：
- MaxSize：单文件最大大小(MB)
- MaxBackups：保留的旧文件数量
- MaxAge：文件保留天数
- Compress：是否压缩旧文件

### 3. 获取器函数

#### GetAppLogger(cfg) *zap.Logger
**懒加载机制**：
```go
if appLogger == nil {
    InitLogger(cfg) // 首次调用时初始化
}
return appLogger
```

#### GetHTTPLogger() *zap.Logger
**直接返回**：返回已初始化的httpLogger

#### GetErrorLogger() *zap.Logger
**直接返回**：返回已初始化的errorLogger

### 4. 便捷日志函数

#### LogInfo/LogWarn/LogDebug(msg, fields...)
**作用**：封装appLogger的对应方法
**空指针保护**：检查appLogger != nil

#### LogError(msg, fields...)
**双重记录**：
1. 记录到appLogger
2. 同时记录到errorLogger（错误专用）

### 5. Gin中间件

#### Logger() gin.HandlerFunc
**功能**：HTTP请求日志中间件
**记录信息**：
- 请求方法、路径、参数
- 响应状态码、大小、延迟
- 客户端IP、User-Agent
- 错误信息

**双重记录机制**：
- 所有请求 → httpLogger
- 错误请求(status >= 400) → errorLogger

### 6. 维护函数

#### Sync()
**作用**：刷新所有日志器的缓冲区
**调用场景**：程序退出前，确保所有日志写入磁盘

## 运行流程分析

### 系统启动流程
```
1. main.go 加载配置
2. 调用 pkglog.InitLogger(&cfg.Logging)
3. 创建三个类型的logger实例
4. 设置全局变量
5. 应用就绪，可以使用日志功能
```

### HTTP请求处理流程
```
1. gin.Engine.Use(pkglog.Logger())
2. 请求到达，记录开始时间
3. 执行 c.Next() 处理业务逻辑
4. 计算响应时间和收集响应信息
5. 写入httpLogger
6. 如果是错误状态，同时写入errorLogger
7. 调用LogInfo记录简要信息
```

### 应用日志记录流程
```
1. 业务代码调用 pkglog.LogInfo/LogError等
2. 检查对应的logger是否初始化
3. 调用zap.Logger的对应方法
4. zap根据配置输出到控制台/文件
5. 错误日志额外写入errorLogger
```

## 配置驱动特性

### 输出方式
- **stdout**：仅控制台输出
- **file**：仅文件输出  
- **both**：同时输出到控制台和文件

### 日志级别
- **debug**：开发环境，输出所有日志
- **info**：生产环境标准级别
- **warn**：仅警告和错误
- **error**：仅错误日志

### 格式选择
- **json**：结构化JSON格式，便于日志分析
- **console**：人类可读的控制台格式

## 错误处理策略

### 初始化错误
- 配置错误：返回详细错误信息
- 文件创建失败：创建目录后重试
- Logger创建失败：使用默认配置备用

### 运行时错误
- 空指针保护：所有日志函数检查logger != nil
- 静默失败：Sync()函数忽略错误（使用_接收）
- 默认配置：中间件在logger未初始化时自动初始化

## 性能考虑

### 懒加载
- GetAppLogger使用懒加载，首次调用时才初始化
- 避免不必要的资源消耗

### 缓冲机制
- zap内置缓冲机制，提高写入性能
- Sync()函数确保缓冲区及时刷新

### 文件轮转
- lumberjack自动处理文件轮转
- 避免单个日志文件过大影响性能 