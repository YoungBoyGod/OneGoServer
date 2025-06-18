package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// 配置结构体
type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
}

type ServerConfig struct {
	Port            string `yaml:"port"`
	Mode            string `yaml:"mode"`
	Host            string `yaml:"host"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // 优雅关闭超时时间(秒)
	ReadTimeout     int    `yaml:"read_timeout"`     // 读取超时时间(秒)
	WriteTimeout    int    `yaml:"write_timeout"`    // 写入超时时间(秒)
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

var (
	port       string
	configFile string
	config     Config
)

var rootCmd = &cobra.Command{
	Use:   "OneGoTask",
	Short: "OneGoTask 服务管理工具",
	Long:  "OneGoTask 服务端管理工具，提供启动、停止等功能",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 HTTP 服务器",
	Long:  "启动 OneGoTask HTTP 服务器，提供 API 接口，支持优雅启动和关闭",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	// 添加 server 子命令到根命令
	rootCmd.AddCommand(serverCmd)

	// 为 server 命令添加标志
	serverCmd.Flags().StringVarP(&port, "port", "p", "", "指定服务器监听端口")
	serverCmd.Flags().StringVarP(&configFile, "config", "c", "config/server.yaml", "指定配置文件路径")
}

// 加载配置文件
func loadConfig(configPath string) error {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("⚠️ 配置文件不存在: %s，使用默认配置", configPath)
		// 使用默认配置
		config = Config{
			Server: ServerConfig{
				Port:            "30000",
				Mode:            "release",
				Host:            "0.0.0.0",
				ShutdownTimeout: 30,
				ReadTimeout:     30,
				WriteTimeout:    30,
			},
			Log: LogConfig{
				Level: "info",
				File:  "logs/app.log",
			},
		}
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	if config.Server.ShutdownTimeout == 0 {
		config.Server.ShutdownTimeout = 30
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}

	log.Printf("✅ 成功加载配置文件: %s", configPath)
	return nil
}

// 合并命令行参数和配置文件
func mergeConfig() {
	// 命令行参数优先级高于配置文件
	if port != "" {
		config.Server.Port = port
		log.Printf("🔧 使用命令行指定的端口: %s", port)
	} else {
		log.Printf("📝 使用配置文件中的端口: %s", config.Server.Port)
	}

	// 验证端口号
	if portNum, err := strconv.Atoi(config.Server.Port); err != nil || portNum <= 0 || portNum > 65535 {
		log.Printf("⚠️ 无效的端口号: %s，使用默认端口: 30000", config.Server.Port)
		config.Server.Port = "30000"
	}
}

// 设置路由
func setupRoutes() *gin.Engine {
	// 根据配置设置 Gin 模式
	switch config.Server.Mode {
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
	router := gin.Default()

	// 添加基本路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "success",
			"config": gin.H{
				"port": config.Server.Port,
				"mode": config.Server.Mode,
			},
			"timestamp": time.Now().Unix(),
		})
	})

	// 添加健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "OneGoTask",
			"version": "1.0.0",
			"config": gin.H{
				"server_mode": config.Server.Mode,
				"log_level":   config.Log.Level,
			},
			"uptime":    time.Now().Unix(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 添加配置信息路由
	router.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"server": gin.H{
				"port":             config.Server.Port,
				"mode":             config.Server.Mode,
				"host":             config.Server.Host,
				"shutdown_timeout": config.Server.ShutdownTimeout,
				"read_timeout":     config.Server.ReadTimeout,
				"write_timeout":    config.Server.WriteTimeout,
			},
			"log": gin.H{
				"level": config.Log.Level,
				"file":  config.Log.File,
			},
		})
	})

	// 添加优雅重启接口
	router.POST("/restart", func(c *gin.Context) {
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
	})

	return router
}

func startServer() {
	// 加载配置文件
	if err := loadConfig(configFile); err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 合并命令行参数
	mergeConfig()

	// 设置路由
	router := setupRoutes()

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	// 打印启动信息
	printStartupInfo()

	// 创建信号通道
	quit := make(chan os.Signal, 1)
	restart := make(chan os.Signal, 1)

	// 监听系统信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(restart, syscall.SIGUSR1)

	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("🌟 服务器开始监听...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	// 等待信号
	for {
		select {
		case <-quit:
			log.Printf("🛑 收到关闭信号，开始优雅关闭...")
			gracefulShutdown(srv)
			return
		case <-restart:
			log.Printf("🔄 收到重启信号，开始优雅重启...")
			gracefulRestart(srv)
		}
	}
}

// 优雅关闭
func gracefulShutdown(srv *http.Server) {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	log.Printf("⏳ 等待现有连接完成，最多等待 %d 秒...", config.Server.ShutdownTimeout)

	// 关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ 强制关闭服务器: %v", err)
		srv.Close()
	} else {
		log.Printf("✅ 服务器已优雅关闭")
	}
}

// 优雅重启
func gracefulRestart(srv *http.Server) {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	log.Printf("⏳ 正在关闭旧服务器...")

	// 关闭当前服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ 关闭旧服务器失败: %v", err)
		return
	}

	log.Printf("✅ 旧服务器已关闭，重新加载配置...")

	// 重新加载配置
	if err := loadConfig(configFile); err != nil {
		log.Printf("❌ 重新加载配置失败: %v", err)
		return
	}

	// 重新合并配置
	mergeConfig()

	// 重新设置路由
	router := setupRoutes()

	// 创建新的服务器
	newSrv := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	log.Printf("�� 启动新服务器...")
	printStartupInfo()

	// 启动新服务器
	go func() {
		if err := newSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 新服务器启动失败: %v", err)
		}
	}()

	// 更新服务器引用
	*srv = *newSrv

	log.Printf("✅ 服务器重启完成")
}

// 打印启动信息
func printStartupInfo() {
	fmt.Printf("\n🚀 OneGoTask 服务器启动成功！\n")
	fmt.Printf("📍 监听地址: http://%s:%s\n", config.Server.Host, config.Server.Port)
	fmt.Printf("🔍 健康检查: http://localhost:%s/health\n", config.Server.Port)
	fmt.Printf("📡 测试接口: http://localhost:%s/ping\n", config.Server.Port)
	fmt.Printf("⚙️  配置信息: http://localhost:%s/config\n", config.Server.Port)
	fmt.Printf("🔄 优雅重启: curl -X POST http://localhost:%s/restart\n", config.Server.Port)
	fmt.Printf("📂 配置文件: %s\n", configFile)
	fmt.Printf("🎯 运行模式: %s\n", config.Server.Mode)
	fmt.Printf("⏰ 超时配置: 关闭(%ds) 读取(%ds) 写入(%ds)\n",
		config.Server.ShutdownTimeout, config.Server.ReadTimeout, config.Server.WriteTimeout)
	fmt.Printf("🛑 优雅关闭: Ctrl+C 或 kill -TERM %d\n", os.Getpid())
	fmt.Printf("🔄 优雅重启: kill -USR1 %d\n\n", os.Getpid())
}

func main() {
	// 执行 Cobra 命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}
