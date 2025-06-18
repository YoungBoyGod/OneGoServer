package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	port string
)

var rootCmd = &cobra.Command{
	Use:   "OneGoTask",
	Short: "OneGoTask 服务管理工具",
	Long:  "OneGoTask 服务端管理工具，提供启动、停止等功能",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 HTTP 服务器",
	Long:  "启动 OneGoTask HTTP 服务器，提供 API 接口",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	// 添加 server 子命令到根命令
	rootCmd.AddCommand(serverCmd)

	// 为 server 命令添加端口标志
	serverCmd.Flags().StringVarP(&port, "port", "p", "30000", "指定服务器监听端口")
}

func startServer() {
	// 配置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由器
	router := gin.Default()

	// 添加基本路由
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "success",
		})
	})

	// 添加健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "OneGoTask",
		})
	})

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("🚀 OneGoTask 服务器启动成功！\n")
	fmt.Printf("📍 监听地址: http://localhost:%s\n", port)
	fmt.Printf("🔍 健康检查: http://localhost:%s/health\n", port)
	fmt.Printf("📡 测试接口: http://localhost:%s/ping\n", port)

	// 启动服务器
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}

func main() {
	// 执行 Cobra 命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}
