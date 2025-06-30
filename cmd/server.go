package cmd

import (
	"log"
	"strconv"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/YoungBoyGod/OneGoServer/internal/data/postgres"
	"github.com/YoungBoyGod/OneGoServer/internal/data/redis"
	"github.com/YoungBoyGod/OneGoServer/internal/server"
	"github.com/spf13/cobra"
)

var (
	// 服务器端口
	port string
	// 服务器运行模式
	mode string
)

// serverCmd server子命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动OneGoServer HTTP服务器",
	Long: `启动OneGoServer HTTP服务器，提供REST API接口。

该命令将启动以下服务组件：
  • HTTP服务器 (Gin)
  • PostgreSQL数据库连接
  • Redis缓存连接
  • API路由注册
  • 健康检查端点

服务器将监听HTTP请求并提供：
  • 用户管理API
  • 设备管理API  
  • 任务管理API
  • 告警系统API`,
	Example: `  # 使用默认配置启动服务器
  onego server

  # 指定配置文件启动
  onego server --config configs/config.production.yaml

  # 指定端口和运行模式
  onego server --port 9090 --mode release

  # 显示详细启动信息
  onego server --verbose`,
	Run: runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	// 显示启动信息
	if verbose {
		log.Printf("🚀 正在启动 OneGoServer...")
		log.Printf("📁 配置文件: %s", cfgFile)
		log.Printf("🌐 端口: %s", port)
		log.Printf("🔧 运行模式: %s", mode)
	}

	// 1. 加载配置
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 命令行参数覆盖配置文件
	if port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		} else {
			log.Fatalf("❌ 端口格式错误: %v", err)
		}
	}
	if mode != "" {
		cfg.Server.Mode = mode
	}

	if verbose {
		log.Printf("✅ 配置加载成功")
	}

	// 2. 初始化数据库连接
	if err := postgres.InitPostgreSQL(cfg); err != nil {
		log.Fatalf("❌ 数据库初始化失败: %v", err)
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Printf("⚠️ 数据库关闭失败: %v", err)
		}
	}()

	if verbose {
		log.Printf("✅ 数据库连接成功")
	}

	// 3. 初始化Redis连接
	if err := redis.InitRedis(cfg); err != nil {
		log.Fatalf("❌ Redis初始化失败: %v", err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("⚠️ Redis关闭失败: %v", err)
		}
	}()

	if verbose {
		log.Printf("✅ Redis连接成功")
	}

	// 4. 创建并启动服务器
	srv := server.NewServer(cfg, postgres.GetDB(), redis.GetClient())

	log.Printf("🚀 OneGoServer 启动成功!")
	log.Printf("🌐 服务地址: http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📊 健康检查: http://%s:%d/api/v1/health", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📚 API文档: http://%s:%d/api/v1", cfg.Server.Host, cfg.Server.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}

func init() {
	// 添加server子命令到根命令
	rootCmd.AddCommand(serverCmd)

	// server命令的标志
	serverCmd.Flags().StringVarP(&port, "port", "p", "", "服务器端口 (覆盖配置文件)")
	serverCmd.Flags().StringVarP(&mode, "mode", "m", "", "运行模式: debug/release/test (覆盖配置文件)")

	// 设置标志说明
	serverCmd.Flags().Lookup("port").Usage = "指定HTTP服务器监听端口"
	serverCmd.Flags().Lookup("mode").Usage = "指定Gin运行模式 (debug/release/test)"
}
