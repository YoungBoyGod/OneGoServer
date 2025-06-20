package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"learngo0619/internal/config"
	"learngo0619/internal/logger"
	"learngo0619/internal/server/api/handlers"
	"learngo0619/internal/server/api/routes"
	"learngo0619/internal/server/services"

	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	server *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start(verbose bool) error {
	// 初始化日志系统
	if err := logger.Init(s.config); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 初始化注册token和client列表
	if err := services.InitRegisterToken(); err != nil {
		return fmt.Errorf("failed to init register token: %w", err)
	}
	if err := services.InitClientList(); err != nil {
		return fmt.Errorf("failed to init client list: %w", err)
	}

	// 初始化客户端管理器
	handlers.InitClientManager()

	// 设置全局配置供处理器使用
	handlers.SetGlobalConfig(s.config)

	// 记录认证配置状态
	logger.Info("Predefined token authentication mode enabled",
		zap.String("token_configured", func() string {
			if s.config.Security.PredefinedToken != "" {
				return "yes"
			}
			return "no"
		}()),
		zap.String("token_prefix", func() string {
			if s.config.Security.PredefinedToken != "" {
				if len(s.config.Security.PredefinedToken) > 8 {
					return s.config.Security.PredefinedToken[:8] + "..."
				}
				return s.config.Security.PredefinedToken
			}
			return "none"
		}()),
	)

	// 注册清理函数
	defer func() {
		logger.Cleanup()
		handlers.StopClientManager() // 停止客户端管理器
	}()

	// 创建路由
	router := routes.SetupRouter(s.config)

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:    s.config.Server.Host + ":" + s.config.Server.Port,
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		s.printStartupInfo(verbose)

		logger.Info("Starting HTTP server",
			zap.String("address", s.server.Addr),
			zap.String("mode", s.config.Server.Mode),
		)

		// 写入PID文件
		if err := logger.WritePIDFile(s.config.App.Name); err != nil {
			logger.Warn("Failed to write PID file", zap.Error(err))
		}

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	return s.waitForShutdown()
}

// printStartupInfo prints server startup information
func (s *Server) printStartupInfo(verbose bool) {
	pid, sessionID, macAddress := logger.GetInstanceInfo()
	// 如果启动成功，就把这个pid写入到文件中,文件名就是程序名称.pid
	pidFile := fmt.Sprintf("%s.pid", s.config.App.Name)
	os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644)

	// 输出类似于用户期望的格式
	fmt.Printf("build running pid: %d\n", pid)

	fmt.Printf("session: {%s}\n", sessionID)
	fmt.Printf("server mac: %s\n", macAddress)
	fmt.Printf("🚀 服务器启动在 http://%s:%s\n", s.config.Server.Host, s.config.Server.Port)
	fmt.Printf("📋 模式: %s\n", s.config.Server.Mode)
	fmt.Printf("📱 应用: %s v%s\n", s.config.App.Name, s.config.App.Version)
	fmt.Printf("📝 日志级别: %s\n", s.config.Log.Level)
	fmt.Printf("📂 日志输出: %s\n", s.config.Log.Output)
	if s.config.Log.Output == "file" || s.config.Log.Output == "both" {
		fmt.Printf("📄 日志目录: %s\n", s.config.Log.Dir)
	}
	fmt.Println("按 Ctrl+C 优雅关闭服务器")
	fmt.Println(strings.Repeat("-", 50))

	if verbose {
		fmt.Printf("📍 详细配置信息:\n")
		fmt.Printf("   进程ID: %d\n", pid)
		fmt.Printf("   会话ID: %s\n", sessionID)
		fmt.Printf("   服务器MAC: %s\n", macAddress)
		fmt.Printf("   主机: %s\n", s.config.Server.Host)
		fmt.Printf("   端口: %s\n", s.config.Server.Port)
		fmt.Printf("   模式: %s\n", s.config.Server.Mode)
		fmt.Printf("   应用名: %s\n", s.config.App.Name)
		fmt.Printf("   版本: %s\n", s.config.App.Version)
		fmt.Printf("   日志级别: %s\n", s.config.Log.Level)
		fmt.Printf("   日志格式: %s\n", s.config.Log.Format)
		fmt.Printf("   日志输出: %s\n", s.config.Log.Output)
		fmt.Printf("   日志目录: %s\n", s.config.Log.Dir)
		fmt.Println(strings.Repeat("-", 50))
	}
}

// waitForShutdown waits for interrupt signal and performs graceful shutdown
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	fmt.Printf("\n🛑 收到信号: %v，开始优雅关闭服务器...\n", sig)

	// 优雅关闭的上下文，超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return fmt.Errorf("强制关闭服务器: %v", err)
	}

	// 删除PID文件
	if err := logger.RemovePIDFile(s.config.App.Name); err != nil {
		logger.Warn("Failed to remove PID file", zap.Error(err))
	}

	logger.Info("Server shutdown completed")
	fmt.Println("✅ 服务器已优雅关闭")
	fmt.Printf("🕐 关闭时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.config
}
