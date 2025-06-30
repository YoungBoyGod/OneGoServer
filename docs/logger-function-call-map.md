# Logger函数调用关系映射表

## 函数分层结构

### 配置层
- `config.LoggingConfig` - 日志配置结构体

### 初始化层  
- `InitLogger(cfg)` - 主初始化函数
- `createLogger(cfg, logType)` - 核心创建函数
- `generateTimestampLogPath(path, logType)` - 生成文件路径
- `getFileWriter(cfg, logType)` - 创建文件写入器

### 存储层
- `appLogger` - 应用主日志器全局变量
- `httpLogger` - HTTP访问日志器全局变量  
- `errorLogger` - 错误日志器全局变量

### 获取器层
- `GetAppLogger(cfg)` - 获取应用日志器
- `GetHTTPLogger()` - 获取HTTP日志器
- `GetErrorLogger()` - 获取错误日志器

### 使用层
- `LogInfo(msg, fields...)` - 记录信息日志
- `LogWarn(msg, fields...)` - 记录警告日志
- `LogError(msg, fields...)` - 记录错误日志
- `LogDebug(msg, fields...)` - 记录调试日志
- `Logger()` - Gin中间件函数
- `Sync()` - 同步缓冲区

## 详细调用关系

### 初始化调用链
```
main.go
└── InitLogger(cfg)
    ├── createLogger(cfg, LoggerTypeApp) → appLogger
    │   ├── generateTimestampLogPath()
    │   ├── getFileWriter()
    │   └── zapcore.NewCore()
    ├── createLogger(cfg, LoggerTypeHTTP) → httpLogger
    │   ├── generateTimestampLogPath()
    │   ├── getFileWriter()
    │   └── zapcore.NewCore()
    └── createLogger(cfg, LoggerTypeError) → errorLogger
        ├── generateTimestampLogPath()
        ├── getFileWriter()
        └── zapcore.NewCore()
```

### 运行时调用链

#### 应用日志记录
```
业务代码
├── LogInfo() → appLogger.Info()
├── LogWarn() → appLogger.Warn()
├── LogError() → appLogger.Error() + errorLogger.Error()
└── LogDebug() → appLogger.Debug()
```

#### HTTP中间件调用
```
gin.Engine.Use(Logger())
└── gin.HandlerFunc
    ├── 收集请求信息
    ├── c.Next() [执行业务逻辑]
    ├── httpLogger.Info("HTTP请求", fields...)
    ├── if status >= 400 → errorLogger.Error("HTTP错误", fields...)
    └── LogInfo(method, fields...)
```

#### 获取器调用
```
GetAppLogger(cfg)
├── if appLogger == nil
│   └── InitLogger(cfg)
└── return appLogger

GetHTTPLogger() → return httpLogger
GetErrorLogger() → return errorLogger
```

## 函数依赖关系

### 核心依赖
- `InitLogger` → `createLogger` (依赖)
- `createLogger` → `generateTimestampLogPath` + `getFileWriter` (依赖)
- `GetAppLogger` → `InitLogger` (条件依赖)

### 使用依赖
- `LogXXX函数` → `appLogger/errorLogger` (依赖全局变量)
- `Logger中间件` → `httpLogger + errorLogger + LogInfo` (依赖)
- `Sync` → `所有logger实例` (依赖)

## 数据流向

### 配置数据流
```
config.yaml → config.LoggingConfig → InitLogger → createLogger → zap.Logger
```

### 日志数据流
```
业务代码 → LogXXX函数 → zap.Logger → 输出目标(console/file)
HTTP请求 → Logger中间件 → httpLogger/errorLogger → 输出目标
```

### 错误处理流
```
LogError → appLogger.Error + errorLogger.Error (双重记录)
HTTP错误 → httpLogger.Info + errorLogger.Error (双重记录)
```

## 生命周期管理

### 启动阶段
1. 加载配置 (`config.LoadConfig`)
2. 初始化日志器 (`InitLogger`)
3. 创建全局logger实例 (`createLogger`)
4. 应用可以使用日志功能

### 运行阶段
1. 业务代码调用 `LogXXX` 函数
2. HTTP请求通过 `Logger()` 中间件记录
3. 错误自动双重记录

### 退出阶段
1. 调用 `Sync()` 刷新缓冲区
2. 确保所有日志写入完成

## 线程安全性

### 全局变量
- `appLogger`, `httpLogger`, `errorLogger` 在初始化后只读
- zap.Logger 本身是线程安全的

### 并发调用
- 多个goroutine可以安全并发调用 `LogXXX` 函数
- Gin中间件在每个请求的goroutine中安全执行

## 错误恢复机制

### 初始化失败
```go
GetAppLogger(cfg) {
    if err := InitLogger(cfg); err != nil {
        // 使用默认配置重试
        InitLogger(defaultConfig)
    }
}
```

### 运行时保护
```go
LogInfo(msg, fields...) {
    if appLogger != nil {  // 空指针检查
        appLogger.Info(msg, fields...)
    }
}
``` 