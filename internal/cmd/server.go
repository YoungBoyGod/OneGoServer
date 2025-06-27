package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	apiv1 "github.com/YoungBoyGod/OneGoServer/api/v1"
	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/controller"
	"github.com/YoungBoyGod/OneGoServer/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server 服务器结构
type Server struct {
	config     *config.Config
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer 创建新的服务器实例
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 初始化zap日志器
	if err := middleware.InitLogger(&s.config.Logging); err != nil {
		return fmt.Errorf("初始化日志器失败: %w", err)
	}

	// 获取应用日志器
	appLogger := middleware.GetAppLogger()

	// 设置Gin模式
	gin.SetMode(s.config.Server.Mode)

	// 创建Gin引擎
	s.engine = gin.New()

	// 添加中间件
	s.setupMiddleware()

	// 设置路由
	s.setupRoutes()

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 使用zap记录启动日志
	appLogger.Info("🚀 服务器启动成功",
		zap.String("address", addr),
		zap.String("name", s.config.Server.Name),
		zap.String("version", s.config.Server.Version),
		zap.String("mode", s.config.Server.Mode),
	)

	// 显示日志配置信息
	actualLogPath := s.config.Logging.FilePath
	if s.config.Logging.Output == "file" || s.config.Logging.Output == "both" {
		// 生成实际的日志文件路径用于显示
		if strings.Contains(actualLogPath, "{timestamp}") {
			timestamp := time.Now().Format("20060102150405")
			actualLogPath = strings.Replace(actualLogPath, "{timestamp}", timestamp, -1)
		}
	}

	appLogger.Info("📄 日志配置",
		zap.String("level", s.config.Logging.Level),
		zap.String("output", s.config.Logging.Output),
		zap.String("app_file", strings.Replace(actualLogPath, "app_", "app_", 1)),
		zap.String("http_file", strings.Replace(actualLogPath, "app_", "http_", 1)),
		zap.String("error_file", strings.Replace(actualLogPath, "app_", "error_", 1)),
	)

	// 启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("❌ 服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	s.gracefulShutdown()

	return nil
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware() {
	// 恢复中间件
	s.engine.Use(gin.Recovery())

	// 自定义日志中间件
	s.engine.Use(middleware.Logger())

	// CORS中间件
	s.engine.Use(middleware.CORS())

	// 限流中间件（如果启用）
	if s.config.Security.RateLimit.Enabled {
		s.engine.Use(middleware.RateLimit(
			s.config.Security.RateLimit.RequestsPerMinute,
			s.config.Security.RateLimit.Burst,
		))
	}
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 创建控制器实例
	healthController := controller.NewHealthController(s.config)

	// API版本分组
	v1 := s.engine.Group("/api/v1")
	{
		// 健康检查 - 使用控制器方式（如果有Check方法）或直接使用API函数
		v1.GET("/health", apiv1.HealthCheck)

		// 设备相关路由
		v1.GET("/device", apiv1.Device)
		// 客户端相关路由
		v1.GET("/client", apiv1.Client)
		// 任务相关路由
		v1.GET("/task", apiv1.Task)
		// 队列相关路由
		v1.GET("/queue", apiv1.Queue)

		// 用户相关路由（暂时注释，后续实现）
		// userGroup := v1.Group("/users")
		// {
		// 	userGroup.POST("/register", userController.Register)
		// 	userGroup.POST("/login", userController.Login)
		// }
	}

	// 避免未使用的变量警告
	_ = healthController
}

// gracefulShutdown 优雅关闭
func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger := middleware.GetAppLogger()
	appLogger.Info("🛑 正在关闭服务器...")

	// 同步zap日志
	middleware.Sync()

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		appLogger.Error("❌ 服务器强制关闭", zap.Error(err))
	}

	appLogger.Info("✅ 服务器已安全关闭")
}
