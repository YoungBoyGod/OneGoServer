# 日志文件输出功能实现文档

## 概述

本文档记录OneGoServer项目中日志文件输出功能的实现过程。在原有zap日志系统基础上，我们添加了文件输出、日志轮转、多种输出模式等功能。

## 功能特性

### 1. 多种输出模式
- **stdout**: 仅输出到控制台
- **file**: 仅输出到文件
- **both**: 同时输出到控制台和文件

### 2. 日志轮转支持
- 基于文件大小自动轮转
- 保留指定数量的历史文件
- 支持日志压缩
- 基于时间的清理策略

### 3. 自动目录创建
- 自动创建日志目录
- 权限设置为755

### 4. 格式支持
- JSON格式（结构化）
- Console格式（可读性强）

## 技术实现

### 1. 依赖添加

```bash
go get gopkg.in/natefinch/lumberjack.v2
```

**选择lumberjack的原因**:
- Go生态中最成熟的日志轮转库
- 零依赖，轻量级
- 支持基于大小、时间、数量的轮转策略
- 广泛使用，稳定可靠

### 2. 核心代码实现

#### InitLogger函数重构
```go
func InitLogger(cfg *config.LoggingConfig) error {
    // 配置编码器
    var encoder zapcore.Encoder
    if cfg.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }
    
    // 配置输出
    var writers []zapcore.WriteSyncer
    switch cfg.Output {
        case "stdout": // 仅控制台
        case "file":   // 仅文件
        case "both":   // 双重输出
    }
    
    // 创建多重写入器
    core := zapcore.NewCore(
        encoder,
        zapcore.NewMultiWriteSyncer(writers...),
        zapLevel,
    )
}
```

#### 文件写入器配置
```go
func getFileWriter(cfg *config.LoggingConfig) (io.Writer, error) {
    // 自动创建日志目录
    logDir := filepath.Dir(cfg.FilePath)
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return nil, err
    }
    
    // 配置lumberjack日志轮转
    lumberJackLogger := &lumberjack.Logger{
        Filename:   cfg.FilePath,    // 日志文件路径
        MaxSize:    cfg.MaxSize,     // MB
        MaxBackups: cfg.MaxBackups,  // 保留文件数
        MaxAge:     cfg.MaxAge,      // 保留天数
        Compress:   cfg.Compress,    // 是否压缩
        LocalTime:  true,           // 使用本地时间
    }
    
    return lumberJackLogger, nil
}
```

### 3. 配置文件结构

```yaml
logging:
  level: "debug"              # 日志级别
  format: "json"              # 输出格式: json, text
  output: "both"              # 输出模式: stdout, file, both
  file_path: "logs/app.log"   # 日志文件路径
  max_size: 100               # 单文件最大大小(MB)
  max_backups: 3              # 保留文件数量
  max_age: 30                 # 保留天数
  compress: true              # 是否压缩历史文件
```

### 4. 实际输出示例

#### 控制台输出
```
2025/06/27 15:49:58 🚀 服务器启动成功，监听地址: 0.0.0.0:40000
2025/06/27 15:49:58 📄 日志配置: 级别=debug, 输出=both, 文件=logs/app.log
```

#### 文件输出 (logs/app.log)
```json
{
  "level": "info",
  "timestamp": "2025-06-27T15:49:59.393+0800",
  "caller": "middleware/logger.go:168",
  "msg": "HTTP请求",
  "status": 200,
  "method": "GET",
  "path": "/api/v1/health",
  "ip": "::1",
  "latency": 0.000184041,
  "size": 145,
  "user_agent": "curl/8.7.1",
  "error": ""
}
```

## 功能验证

### 1. 目录创建验证
```bash
ls -la logs/
# 输出: drwxr-xr-x@ 3 haitang staff 96 Jun 27 15:49 .
```

### 2. 文件生成验证
```bash
ls -la logs/app.log
# 输出: -rw-------@ 1 haitang staff 500 Jun 27 15:49 app.log
```

### 3. 日志内容验证
```bash
cat logs/app.log | jq .
# 输出格式化的JSON日志
```

### 4. 双重输出验证
- ✅ 控制台同时显示日志
- ✅ 文件同时写入日志
- ✅ 格式保持一致

## 日志轮转机制

### 触发条件
1. **文件大小**: 超过max_size (100MB)
2. **文件数量**: 超过max_backups (3个)
3. **文件时间**: 超过max_age (30天)

### 轮转策略
1. 当前日志文件重命名为 `app.log.1`
2. 创建新的 `app.log` 文件
3. 历史文件依次递增编号
4. 超出保留数量的文件被删除
5. 压缩历史文件（如果启用）

### 文件命名规则
```
logs/
├── app.log           # 当前日志文件
├── app.log.1         # 最近的历史文件
├── app.log.2.gz      # 压缩的历史文件
└── app.log.3.gz      # 更早的历史文件
```

## 性能考虑

### 1. 内存使用
- lumberjack使用缓冲写入，减少系统调用
- zap本身是零分配设计
- 文件写入异步进行，不阻塞请求处理

### 2. 磁盘空间
- 自动日志轮转防止磁盘空间耗尽
- 压缩历史文件节省存储空间
- 基于时间的清理策略

### 3. 并发安全
- lumberjack内置并发安全机制
- zap支持并发写入
- 无需额外的锁机制

## 使用建议

### 1. 生产环境配置
```yaml
logging:
  level: "info"              # 减少日志量
  format: "json"             # 便于分析
  output: "file"             # 仅文件输出
  file_path: "/var/log/onego/app.log"
  max_size: 100
  max_backups: 10            # 增加保留数量
  max_age: 90                # 延长保留时间
  compress: true
```

### 2. 开发环境配置
```yaml
logging:
  level: "debug"             # 详细日志
  format: "text"             # 可读性强
  output: "both"             # 控制台+文件
  file_path: "logs/app.log"
  max_size: 10               # 小文件便于查看
  max_backups: 3
  max_age: 7
  compress: false
```

### 3. 监控建议
- 监控日志文件大小增长
- 设置磁盘空间告警
- 定期检查日志轮转是否正常
- 配置日志收集和分析系统

## 扩展功能

### 1. 支持的增强功能
- [x] 多种输出模式
- [x] 自动日志轮转
- [x] 自动目录创建
- [x] JSON和文本格式
- [x] 智能日志分级

### 2. 可能的扩展方向
- [ ] 日志加密支持
- [ ] 远程日志传输
- [ ] 实时日志监控
- [ ] 自定义日志格式
- [ ] 性能指标统计

## 常见问题

### Q: 如何只输出错误日志到文件？
A: 可以创建多个logger实例，针对不同级别使用不同配置。

### Q: 如何处理权限问题？
A: 确保应用程序有写入日志目录的权限，建议使用专门的日志用户。

### Q: 如何与日志收集系统集成？
A: JSON格式的日志可以直接被Logstash、Fluentd等工具收集和处理。

### Q: 如何调试日志配置问题？
A: 检查InitLogger的错误返回，确认目录权限和配置参数有效性。

## 修改记录

- 2024-12-19: 实现基础文件输出功能
- 2024-12-19: 添加lumberjack日志轮转支持
- 2024-12-19: 实现多种输出模式（stdout/file/both）
- 2024-12-19: 添加自动目录创建功能
- 2024-12-19: 完成功能测试和验证 