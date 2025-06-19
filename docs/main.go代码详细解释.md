# main.go 代码详细解释

## 代码概述
这是一个Go语言编写的Web服务器程序，实现了配置加载、HTTP服务启动和优雅关闭功能。

## 逐行代码解释

### 1. 包声明和导入
```go
package main

import (
	"context"     // 用于创建上下文，控制超时
	"log"         // 日志记录
	"net/http"    // HTTP服务器
	"os"          // 操作系统接口
	"os/signal"   // 信号处理
	"syscall"     // 系统调用常量
	"time"        // 时间操作
)
```

**说明**：
- `package main` 表示这是一个可执行程序
- 导入了处理HTTP服务、信号监听、上下文控制等必需的包

### 2. 配置加载函数
```go
func LoadConfig(filename string) (*Config, error) {
	config, err := ReadConfig(filename)  // 调用ReadConfig读取配置文件
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)  // 读取失败时终止程序
	}
	return config, nil
}
```

**说明**：
- 封装了配置文件读取逻辑
- 使用`log.Fatalf`在配置读取失败时立即终止程序
- 返回配置对象指针

### 3. 主函数开始
```go
func main() {
	config, err := LoadConfig("config/server.yaml")  // 加载配置文件
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := setupRouter(config)  // 创建路由器
```

**说明**：
- 程序入口点
- 加载`config/server.yaml`配置文件
- 调用`setupRouter`函数创建Gin路由器

### 4. HTTP服务器创建
```go
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Port,  // 服务器地址
		Handler: router,  // 路由处理器
	}
```

**说明**：
- 创建标准库的HTTP服务器实例
- `Addr`字段设置监听地址（从配置文件读取）
- `Handler`字段设置请求处理器（Gin路由器）

### 5. 异步启动服务器
```go
	// 在goroutine中启动服务器，以便它不会阻塞
	go func() {
		log.Printf("Server starting on %s:%s", config.Server.Host, config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %s\n", err)
		}
	}()
```

**说明**：
- 使用`go func()`创建新的goroutine（协程）
- 在后台运行HTTP服务器，不阻塞主线程
- `ListenAndServe()`开始监听和处理HTTP请求
- 只有在非正常关闭的错误时才会记录致命错误

### 6. 信号监听设置
```go
	// 等待中断信号来优雅地关闭服务器，超时时间为5秒
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT  
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")
```

**说明**：
- `make(chan os.Signal, 1)`创建信号通道，缓冲区大小为1
- `signal.Notify()`注册信号监听：
  - `SIGINT`：Ctrl+C发送的中断信号
  - `SIGTERM`：kill命令发送的终止信号
- `<-quit`阻塞等待信号
- 收到信号后开始关闭流程

### 7. 优雅关闭实现
```go
	// 优雅关闭的上下文，超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
```

**说明**：
- `context.WithTimeout()`创建5秒超时的上下文
- `defer cancel()`确保上下文资源被释放
- `srv.Shutdown(ctx)`优雅关闭服务器：
  - 停止接受新请求
  - 等待现有请求完成
  - 如果5秒内未完成，强制关闭
- 关闭成功后记录退出日志

## 程序执行流程

### 阶段1：初始化
1. 加载配置文件 `config/server.yaml`
2. 创建Gin路由器
3. 创建HTTP服务器实例

### 阶段2：服务启动
1. 在新协程中启动HTTP服务器
2. 主协程设置信号监听
3. 等待中断信号

### 阶段3：优雅关闭
1. 接收到中断信号
2. 创建5秒超时上下文
3. 调用服务器Shutdown方法
4. 等待现有请求完成或超时
5. 程序退出

## 关键技术概念

### 1. Goroutine（协程）
```go
go func() {
    // 异步执行的代码
}()
```
- Go语言的轻量级线程
- 允许并发执行多个任务
- 这里用于让HTTP服务器在后台运行

### 2. Channel（通道）
```go
quit := make(chan os.Signal, 1)
<-quit  // 阻塞等待
```
- Go语言的并发通信机制
- 用于在协程之间传递数据
- 这里用于接收操作系统信号

### 3. Context（上下文）
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```
- 控制操作的截止时间、取消信号等
- `WithTimeout`创建带超时的上下文
- 用于控制优雅关闭的最大等待时间

### 4. 信号处理
```go
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
```
- 监听操作系统信号
- `SIGINT`：通常由Ctrl+C触发
- `SIGTERM`：通常由kill命令触发

## 优雅关闭的好处

### 1. 避免数据丢失
- 等待正在处理的请求完成
- 避免中断数据库事务

### 2. 提供更好的用户体验
- 客户端不会收到连接重置错误
- 响应时间更稳定

### 3. 资源清理
- 正确释放网络连接
- 清理临时资源

## 实际运行效果

从air工具的输出可以看到：
```
Server starting on 0.0.0.0:8080
[GIN-debug] GET    /health    --> main.setupRouter.func3 (3 handlers)
```

这表明：
1. 服务器成功启动在0.0.0.0:8080
2. 路由注册成功
3. 可以正常处理HTTP请求
4. 支持热重载（文件修改后自动重启）

## 总结

这段代码实现了一个**生产级别的HTTP服务器**，具备：
- ✅ 配置文件支持
- ✅ 并发处理能力
- ✅ 优雅关闭机制
- ✅ 完善的错误处理
- ✅ 信号监听支持

这是现代Go Web应用的标准实现模式，保证了服务的稳定性和可靠性。

## 创建时间
${new Date().toLocaleString('zh-CN')} 