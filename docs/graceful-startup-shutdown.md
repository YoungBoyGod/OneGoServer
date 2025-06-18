# Gin 优雅启动和关闭指南

## 概述

本文档记录了 OneGoTask 项目中实现的优雅启动、关闭和重启功能。通过使用 Go 的 context 包和信号处理机制，实现了零停机时间的服务更新和部署。

## 核心特性

### 1. 优雅关闭 (Graceful Shutdown)
- **无连接丢失**: 等待现有请求完成后再关闭服务器
- **超时保护**: 设置最大等待时间，避免无限等待
- **信号监听**: 支持 SIGINT (Ctrl+C) 和 SIGTERM 信号

### 2. 优雅重启 (Graceful Restart)
- **热重启**: 无需停止服务即可重启服务器
- **配置热加载**: 重启时重新加载配置文件
- **零停机**: 新服务器启动后再关闭旧服务器

### 3. 多种触发方式
- **系统信号**: kill 命令触发
- **API 接口**: HTTP POST 请求触发
- **键盘快捷键**: Ctrl+C 优雅关闭

## 技术实现

### 1. 信号处理机制

```go
// 创建信号通道
quit := make(chan os.Signal, 1)
restart := make(chan os.Signal, 1)

// 监听系统信号
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
signal.Notify(restart, syscall.SIGUSR1)

// 等待信号
for {
    select {
    case <-quit:
        gracefulShutdown(srv)
        return
    case <-restart:
        gracefulRestart(srv)
    }
}
```

### 2. 优雅关闭实现

```go
func gracefulShutdown(srv *http.Server) {
    // 创建关闭上下文，设置超时时间
    ctx, cancel := context.WithTimeout(
        context.Background(), 
        time.Duration(config.Server.ShutdownTimeout)*time.Second
    )
    defer cancel()
    
    // 优雅关闭服务器
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("❌ 强制关闭服务器: %v", err)
        srv.Close() // 强制关闭
    } else {
        log.Printf("✅ 服务器已优雅关闭")
    }
}
```

### 3. 优雅重启实现

```go
func gracefulRestart(srv *http.Server) {
    // 1. 关闭当前服务器
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    srv.Shutdown(ctx)
    
    // 2. 重新加载配置
    loadConfig(configFile)
    mergeConfig()
    
    // 3. 创建新服务器
    router := setupRoutes()
    newSrv := &http.Server{
        Addr:    config.Server.Host + ":" + config.Server.Port,
        Handler: router,
        // ... 其他配置
    }
    
    // 4. 启动新服务器
    go newSrv.ListenAndServe()
    
    // 5. 更新服务器引用
    *srv = *newSrv
}
```

## 配置选项

### server.yaml 配置文件

```yaml
server:
  port: 30000
  mode: debug
  host: "0.0.0.0"
  shutdown_timeout: 30  # 优雅关闭超时时间(秒)
  read_timeout: 30      # 读取超时时间(秒)
  write_timeout: 30     # 写入超时时间(秒)
```

### 超时配置说明

- **shutdown_timeout**: 优雅关闭最大等待时间
- **read_timeout**: HTTP 请求读取超时时间
- **write_timeout**: HTTP 响应写入超时时间

## 使用方法

### 1. 启动服务器

```bash
# 启动服务器
./OneGoTask server

# 启动时会显示进程ID和操作指令
🛑 优雅关闭: Ctrl+C 或 kill -TERM 12345
🔄 优雅重启: kill -USR1 12345
```

### 2. 优雅关闭

```bash
# 方法1: 键盘快捷键
Ctrl+C

# 方法2: 发送SIGTERM信号
kill -TERM <PID>

# 方法3: 发送SIGINT信号
kill -INT <PID>
```

### 3. 优雅重启

```bash
# 方法1: 发送SIGUSR1信号
kill -USR1 <PID>

# 方法2: 通过API接口
curl -X POST http://localhost:30000/restart

# 方法3: 使用systemd (生产环境)
systemctl reload onegotask
```

### 4. 监控服务状态

```bash
# 检查服务健康状态
curl http://localhost:30000/health

# 查看配置信息
curl http://localhost:30000/config

# 测试服务响应
curl http://localhost:30000/ping
```

## API 接口

### POST /restart - 触发优雅重启

**请求示例:**
```bash
curl -X POST http://localhost:30000/restart
```

**响应示例:**
```json
{
  "message": "服务器将在处理完当前请求后重启",
  "status": "restarting"
}
```

### GET /config - 查看当前配置

**响应示例:**
```json
{
  "server": {
    "port": "30000",
    "mode": "debug",
    "host": "0.0.0.0",
    "shutdown_timeout": 30,
    "read_timeout": 30,
    "write_timeout": 30
  },
  "log": {
    "level": "info",
    "file": "logs/app.log"
  }
}
```

## 日志监控

### 启动日志
```
🚀 OneGoTask 服务器启动成功！
📍 监听地址: http://0.0.0.0:30000
🔍 健康检查: http://localhost:30000/health
📡 测试接口: http://localhost:30000/ping
⚙️  配置信息: http://localhost:30000/config
🔄 优雅重启: curl -X POST http://localhost:30000/restart
📂 配置文件: config/server.yaml
🎯 运行模式: debug
⏰ 超时配置: 关闭(30s) 读取(30s) 写入(30s)
🛑 优雅关闭: Ctrl+C 或 kill -TERM 12345
🔄 优雅重启: kill -USR1 12345
```

### 重启日志
```
🔄 收到重启信号，开始优雅重启...
⏳ 正在关闭旧服务器...
✅ 旧服务器已关闭，重新加载配置...
✅ 成功加载配置文件: config/server.yaml
📝 使用配置文件中的端口: 30000
🐛 Gin 运行在调试模式
🚀 启动新服务器...
✅ 服务器重启完成
```

### 关闭日志
```
🛑 收到关闭信号，开始优雅关闭...
⏳ 等待现有连接完成，最多等待 30 秒...
✅ 服务器已优雅关闭
```

## 生产环境部署

### 1. Systemd 服务配置

创建 `/etc/systemd/system/onegotask.service`:

```ini
[Unit]
Description=OneGoTask HTTP Server
After=network.target

[Service]
Type=simple
User=onegotask
Group=onegotask
WorkingDirectory=/opt/onegotask
ExecStart=/opt/onegotask/OneGoTask server
ExecReload=/bin/kill -USR1 $MAINPID
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### 2. 部署脚本

```bash
#!/bin/bash
# deploy.sh - 零停机部署脚本

# 1. 构建新版本
go build -o OneGoTask-new

# 2. 备份当前版本
cp OneGoTask OneGoTask-backup

# 3. 替换可执行文件
mv OneGoTask-new OneGoTask

# 4. 触发优雅重启
systemctl reload onegotask

# 5. 验证服务状态
sleep 5
curl -f http://localhost:30000/health || {
    echo "健康检查失败，回滚到备份版本"
    mv OneGoTask-backup OneGoTask
    systemctl reload onegotask
    exit 1
}

echo "部署成功完成"
```

## 优势与特点

### 1. 零停机时间
- 新连接在重启期间不会被拒绝
- 现有连接得到妥善处理
- 用户无感知的服务更新

### 2. 高可靠性
- 超时保护机制防止无限等待
- 失败时自动回退到强制关闭
- 详细的日志记录便于监控

### 3. 易于操作
- 多种触发方式满足不同场景
- 清晰的控制台输出和状态提示
- 与标准系统工具集成良好

### 4. 配置灵活
- 可调节的超时时间
- 支持配置热重载
- 适应不同的部署环境

## 注意事项

### 1. 信号处理
- SIGUSR1 专用于优雅重启
- SIGTERM/SIGINT 用于优雅关闭
- 避免使用 SIGKILL (kill -9)

### 2. 超时设置
- 根据应用特点调整超时时间
- 长连接应用需要更长的超时时间
- 考虑负载均衡器的超时设置

### 3. 资源管理
- 确保旧服务器完全关闭后再启动新服务器
- 监控内存和文件描述符使用情况
- 适当的日志轮转和清理策略

## 故障排除

### 1. 重启失败
```bash
# 检查配置文件语法
./OneGoTask server --config config/server.yaml --help

# 查看详细错误日志
tail -f logs/app.log
```

### 2. 端口占用
```bash
# 检查端口占用情况
netstat -tulnp | grep :30000
lsof -i :30000
```

### 3. 权限问题
```bash
# 检查文件权限
ls -la OneGoTask config/server.yaml

# 检查进程权限
ps aux | grep OneGoTask
```

## 扩展建议

1. **健康检查增强**: 添加更详细的健康状态检查
2. **指标监控**: 集成 Prometheus 监控指标
3. **滚动部署**: 支持多实例滚动更新
4. **配置验证**: 重启前验证配置文件有效性
5. **回滚机制**: 自动检测失败并回滚到稳定版本

## 修改记录

- **2024年**: 实现优雅启动和关闭功能
- **功能**: 信号处理、上下文管理、热重启
- **新增**: `/restart` API、超时配置、详细日志 