# YAML 配置集成指南

## 概述

本文档记录了 OneGoTask 项目中 YAML 配置文件读取功能的实现，支持配置文件优先、命令行参数补充的配置机制。

## 配置优先级

配置参数的优先级从高到低：

1. **命令行参数** (最高优先级)
2. **配置文件参数**  
3. **默认配置** (最低优先级)

## 配置结构

### 配置文件结构 (`config/server.yaml`)

```yaml
server:
  port: "30000"        # 服务器端口
  mode: "debug"        # 运行模式: debug, release, test
  host: "0.0.0.0"      # 监听地址

log:
  level: "info"        # 日志级别: debug, info, warn, error
  file: "logs/app.log" # 日志文件路径
```

### Go 结构体定义

```go
type Config struct {
    Server ServerConfig `yaml:"server"`
    Log    LogConfig    `yaml:"log"`
}

type ServerConfig struct {
    Port string `yaml:"port"`
    Mode string `yaml:"mode"`
    Host string `yaml:"host"`
}

type LogConfig struct {
    Level string `yaml:"level"`
    File  string `yaml:"file"`
}
```

## 新增功能

### 1. 配置加载函数

#### `loadConfig(configPath string) error`
- 检查配置文件是否存在
- 读取并解析 YAML 文件
- 如果文件不存在，使用默认配置
- 提供详细的错误处理和日志输出

```go
func loadConfig(configPath string) error {
    // 检查配置文件是否存在
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        log.Printf("⚠️ 配置文件不存在: %s，使用默认配置", configPath)
        // 使用默认配置
        config = Config{...}
        return nil
    }
    
    // 读取和解析配置文件
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("读取配置文件失败: %v", err)
    }
    
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return fmt.Errorf("解析配置文件失败: %v", err)
    }
    
    return nil
}
```

### 2. 配置合并函数

#### `mergeConfig()`
- 处理命令行参数覆盖配置文件参数的逻辑
- 验证参数有效性
- 提供详细的配置来源日志

```go
func mergeConfig() {
    // 命令行参数优先级高于配置文件
    if port != "" {
        config.Server.Port = port
        log.Printf("🔧 使用命令行指定的端口: %s", port)
    } else {
        log.Printf("📝 使用配置文件中的端口: %s", config.Server.Port)
    }
    
    // 验证端口号有效性
    if portNum, err := strconv.Atoi(config.Server.Port); err != nil || portNum <= 0 || portNum > 65535 {
        log.Printf("⚠️ 无效的端口号: %s，使用默认端口: 30000", config.Server.Port)
        config.Server.Port = "30000"
    }
}
```

### 3. 增强的命令行参数

```bash
# 新增配置文件参数
./OneGoTask server --config custom-config.yaml
./OneGoTask server -c custom-config.yaml

# 端口参数现在可以覆盖配置文件
./OneGoTask server --port 8080
./OneGoTask server -p 8080

# 组合使用
./OneGoTask server -c custom-config.yaml -p 8080
```

### 4. 新增 API 接口

#### `/config` - 配置信息接口
```json
{
  "server": {
    "port": "30000",
    "mode": "debug",
    "host": "0.0.0.0"
  },
  "log": {
    "level": "info",
    "file": "logs/app.log"
  }
}
```

#### 增强的 `/ping` 接口
```json
{
  "message": "pong",
  "status": "success",
  "config": {
    "port": "30000",
    "mode": "debug"
  }
}
```

#### 增强的 `/health` 接口
```json
{
  "status": "healthy",
  "service": "OneGoTask",
  "version": "1.0.0",
  "config": {
    "server_mode": "debug",
    "log_level": "info"
  }
}
```

## 运行模式

根据配置文件中的 `server.mode` 设置：

- **debug**: 调试模式，提供详细的请求日志
- **release**: 发布模式，优化性能，简化日志
- **test**: 测试模式，用于单元测试

## 使用示例

### 1. 使用默认配置
```bash
./OneGoTask server
```

### 2. 指定配置文件
```bash
./OneGoTask server --config config/production.yaml
```

### 3. 覆盖配置参数
```bash
# 使用配置文件，但覆盖端口
./OneGoTask server -c config/server.yaml -p 8080
```

### 4. 测试配置
```bash
# 查看当前配置
curl http://localhost:30000/config

# 测试服务
curl http://localhost:30000/ping
curl http://localhost:30000/health
```

## 错误处理

### 配置文件错误
- 文件不存在：自动使用默认配置
- 文件格式错误：显示详细错误信息并退出
- 权限错误：显示权限错误信息

### 参数验证
- 端口号验证：1-65535 范围
- 无效参数：自动回退到默认值
- 详细的警告和错误日志

## 启动信息

程序启动时会显示详细的配置信息：

```
🚀 OneGoTask 服务器启动成功！
📍 监听地址: http://0.0.0.0:30000
🔍 健康检查: http://localhost:30000/health
📡 测试接口: http://localhost:30000/ping
⚙️  配置信息: http://localhost:30000/config
📂 配置文件: config/server.yaml
🎯 运行模式: debug
```

## 优势

1. **灵活配置**: 支持配置文件和命令行参数
2. **优先级明确**: 命令行 > 配置文件 > 默认值
3. **容错性强**: 配置文件缺失时自动使用默认配置
4. **实时查看**: 通过 API 接口查看当前配置
5. **详细日志**: 清晰的配置来源和状态信息

## 扩展建议

1. 支持环境变量配置
2. 添加配置热重载功能
3. 支持多环境配置文件
4. 添加配置验证规则
5. 支持加密敏感配置项

## 修改记录

- **2024年**: 初始实现 YAML 配置读取功能
- **功能**: 配置文件优先级、命令行覆盖、默认配置回退
- **新增**: `/config` API 接口，增强的日志和错误处理 