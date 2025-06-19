package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"OneGoTask/pkg/config"

	"github.com/gin-gonic/gin"
)

// Server HTTP服务器结构体
type Server struct {
	config *config.Config
	srv    *http.Server
	router *gin.Engine
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start 启动服务器
func (s *Server) Start() {
	// 设置路由
	s.setupRoutes()

	// 创建 HTTP 服务器
	s.srv = &http.Server{
		Addr:         s.config.Server.Host + ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	// 打印启动信息
	s.printStartupInfo()

	// 创建信号通道
	quit := make(chan os.Signal, 1)
	restart := make(chan os.Signal, 1)

	// 监听系统信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(restart, syscall.SIGUSR1)

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("🌟 服务器开始监听...")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	// 等待信号
	for {
		select {
		case <-quit:
			log.Printf("🛑 收到关闭信号，开始优雅关闭...")
			s.gracefulShutdown()
			return
		case <-restart:
			log.Printf("🔄 收到重启信号，开始优雅重启...")
			s.gracefulRestart()
		}
	}
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 根据配置设置 Gin 模式
	switch s.config.Server.Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
		log.Printf("🐛 Gin 运行在调试模式")
	case "test":
		gin.SetMode(gin.TestMode)
		log.Printf("🧪 Gin 运行在测试模式")
	default:
		gin.SetMode(gin.ReleaseMode)
		log.Printf("🚀 Gin 运行在发布模式")
	}

	// 创建 Gin 路由器
	s.router = gin.Default()

	// 添加基本路由
	s.router.GET("/ping", s.handlePing)
	s.router.GET("/health", s.handleHealth)
	s.router.GET("/config", s.handleConfig)
	s.router.POST("/restart", s.handleRestart)

	// 添加客户端相关路由
	api := s.router.Group("/api/v1")
	{
		api.GET("/status", s.handleStatus)
		api.GET("/version", s.handleVersion)
	}
}

// handlePing ping接口
func (s *Server) handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"status":  "success",
		"config": gin.H{
			"port": s.config.Server.Port,
			"mode": s.config.Server.Mode,
		},
		"timestamp": time.Now().Unix(),
	})
}

// handleHealth 健康检查接口
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "OneGoTask",
		"version": "1.0.0",
		"config": gin.H{
			"server_mode": s.config.Server.Mode,
			"log_level":   s.config.Log.Level,
		},
		"uptime":    time.Now().Unix(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// handleConfig 配置信息接口
func (s *Server) handleConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"server": gin.H{
			"port":             s.config.Server.Port,
			"mode":             s.config.Server.Mode,
			"host":             s.config.Server.Host,
			"shutdown_timeout": s.config.Server.ShutdownTimeout,
			"read_timeout":     s.config.Server.ReadTimeout,
			"write_timeout":    s.config.Server.WriteTimeout,
		},
		"client": gin.H{
			"server_url":    s.config.Client.ServerURL,
			"timeout":       s.config.Client.Timeout,
			"retry_count":   s.config.Client.RetryCount,
			"retry_delay":   s.config.Client.RetryDelay,
			"output_format": s.config.Client.OutputFormat,
		},
		"log": gin.H{
			"level": s.config.Log.Level,
			"file":  s.config.Log.File,
		},
	})
}

// handleRestart 重启接口
func (s *Server) handleRestart(c *gin.Context) {
	log.Printf("🔄 收到重启请求")
	c.JSON(200, gin.H{
		"message": "服务器将在处理完当前请求后重启",
		"status":  "restarting",
	})

	// 发送重启信号
	go func() {
		time.Sleep(1 * time.Second) // 给响应时间
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGUSR1)
	}()
}

// handleStatus 状态接口
func (s *Server) handleStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "running",
		"server": gin.H{
			"pid":  os.Getpid(),
			"mode": s.config.Server.Mode,
			"port": s.config.Server.Port,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// handleVersion 版本接口
func (s *Server) handleVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"version":    "1.0.0",
		"build_time": "2024",
		"go_version": "1.23+",
		"service":    "OneGoTask",
	})
}

// gracefulShutdown 优雅关闭
func (s *Server) gracefulShutdown() {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	log.Printf("⏳ 等待现有连接完成，最多等待 %d 秒...", s.config.Server.ShutdownTimeout)

	// 关闭服务器
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Printf("❌ 强制关闭服务器: %v", err)
		s.srv.Close()
	} else {
		log.Printf("✅ 服务器已优雅关闭")
	}
}

// gracefulRestart 优雅重启
func (s *Server) gracefulRestart() {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	log.Printf("⏳ 正在关闭旧服务器...")

	// 关闭当前服务器
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Printf("❌ 关闭旧服务器失败: %v", err)
		return
	}

	log.Printf("✅ 旧服务器已关闭，重新加载配置...")

	// 重新加载配置（这里需要配置文件路径，暂时使用默认）
	newConfig, err := config.Load("config/server.yaml")
	if err != nil {
		log.Printf("❌ 重新加载配置失败: %v", err)
		return
	}

	// 更新配置
	s.config = newConfig

	// 重新设置路由
	s.setupRoutes()

	// 创建新的服务器
	s.srv = &http.Server{
		Addr:         s.config.Server.Host + ":" + s.config.Server.Port,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	log.Printf("🚀 启动新服务器...")
	s.printStartupInfo()

	// 启动新服务器
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 新服务器启动失败: %v", err)
		}
	}()

	log.Printf("✅ 服务器重启完成")
}

// printStartupInfo 打印启动信息
func (s *Server) printStartupInfo() {
	fmt.Printf("\n🚀 OneGoTask 服务器启动成功！\n")
	fmt.Printf("📍 监听地址: http://%s:%s\n", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("🔍 健康检查: http://localhost:%s/health\n", s.config.Server.Port)
	fmt.Printf("📡 测试接口: http://localhost:%s/ping\n", s.config.Server.Port)
	fmt.Printf("⚙️  配置信息: http://localhost:%s/config\n", s.config.Server.Port)
	fmt.Printf("📊 服务状态: http://localhost:%s/api/v1/status\n", s.config.Server.Port)
	fmt.Printf("🔄 优雅重启: curl -X POST http://localhost:%s/restart\n", s.config.Server.Port)
	fmt.Printf("📂 配置文件: config/server.yaml\n")
	fmt.Printf("🎯 运行模式: %s\n", s.config.Server.Mode)
	fmt.Printf("⏰ 超时配置: 关闭(%ds) 读取(%ds) 写入(%ds)\n",
		s.config.Server.ShutdownTimeout, s.config.Server.ReadTimeout, s.config.Server.WriteTimeout)
	fmt.Printf("🛑 优雅关闭: Ctrl+C 或 kill -TERM %d\n", os.Getpid())
	fmt.Printf("🔄 优雅重启: kill -USR1 %d\n\n", os.Getpid())
}
