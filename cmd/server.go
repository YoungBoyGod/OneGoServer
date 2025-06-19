package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"learngo0619/server"

	"github.com/spf13/cobra"
)

// serverCmd 代表server命令
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP Web服务器",
	Long: `启动HTTP Web服务器，提供RESTful API服务。

服务器支持以下功能:
• 优雅关闭机制
• 配置文件热重载
• RESTful API端点
• 健康检查接口

示例:
  learngo0619 server                    # 使用默认配置启动
  learngo0619 server -c custom.yaml     # 使用自定义配置文件
  learngo0619 server -v                 # 启用详细日志`,
	Run: runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	// 获取配置文件路径
	configFile, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// 获取主机和端口覆盖选项
	hostOverride, _ := cmd.Flags().GetString("host")
	portOverride, _ := cmd.Flags().GetString("port")

	if verbose {
		log.Printf("使用配置文件: %s", configFile)
	}

	// 加载配置
	config, err := server.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 应用命令行覆盖选项
	if hostOverride != "" {
		config.Server.Host = hostOverride
	}
	if portOverride != "" {
		config.Server.Port = portOverride
	}

	if verbose {
		log.Printf("配置加载成功: %+v", config)
	}

	// 创建路由
	router := server.SetupRouter(config)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Port,
		Handler: router,
	}

	// 在goroutine中启动服务器，以便它不会阻塞
	go func() {
		fmt.Printf("🚀 服务器启动在 http://%s:%s\n", config.Server.Host, config.Server.Port)
		fmt.Printf("📋 模式: %s\n", config.Server.Mode)
		fmt.Printf("📱 应用: %s v%s\n", config.App.Name, config.App.Version)
		fmt.Println("按 Ctrl+C 优雅关闭服务器")
		fmt.Println(strings.Repeat("-", 50))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	fmt.Printf("\n🛑 收到信号: %v，开始优雅关闭服务器...\n", sig)

	// 优雅关闭的上下文，超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("强制关闭服务器:", err)
	}

	fmt.Println("✅ 服务器已优雅关闭")
	fmt.Printf("🕐 关闭时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func init() {
	// 添加server特定的标志
	serverCmd.Flags().StringP("host", "H", "", "服务器主机地址 (覆盖配置文件)")
	serverCmd.Flags().StringP("port", "p", "", "服务器端口 (覆盖配置文件)")

	// 将server命令添加到根命令
	rootCmd.AddCommand(serverCmd)
}
