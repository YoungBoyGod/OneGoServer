package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Server 服务器结构
type Server struct {
	config     *config.Config
	engine     *gin.Engine
	httpServer *http.Server
	db         *gorm.DB
	redis      *redis.Client
}

// NewServer 创建新的服务器实例
func NewServer(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *Server {
	return &Server{
		config: cfg,
		db:     db,
		redis:  rdb,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 设置Gin模式
	gin.SetMode(s.config.Server.Mode)

	// 创建Gin引擎并设置路由
	s.engine = router.NewRouter(&router.Dependencies{
		DB:    s.db,
		Redis: s.redis,
	})

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, strconv.Itoa(s.config.Server.Port))
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 记录启动日志
	fmt.Printf("🚀 服务器启动成功 - 地址: %s, 模式: %s\n", addr, s.config.Server.Mode)

	// 启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ 服务器启动失败: %v\n", err)
		}
	}()

	// 等待中断信号
	s.gracefulShutdown()

	return nil
}

// gracefulShutdown 优雅关闭
func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 正在关闭服务器...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("❌ 服务器强制关闭: %v\n", err)
	}

	fmt.Println("✅ 服务器已安全关闭")
}
