package main

import (
	"log"

	"OneGoTask/pkg/config"
	"OneGoTask/pkg/server"

	"github.com/spf13/cobra"
)

var (
	port       string
	configFile string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 OneGoTask HTTP 服务器",
	Long:  "启动 OneGoTask HTTP 服务器，提供 API 接口，支持优雅启动和关闭",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	// 为 server 命令添加标志
	serverCmd.Flags().StringVarP(&port, "port", "p", "", "指定服务器监听端口")
	serverCmd.Flags().StringVarP(&configFile, "config", "c", "config/server.yaml", "指定配置文件路径")
}

func startServer() {
	// 加载配置文件
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 合并命令行参数
	cfg.MergeServerFlags(port)

	// 创建服务器实例
	srv := server.New(cfg)

	// 启动服务器
	srv.Start()
}

func main() {
	// 执行 Cobra 命令
	if err := serverCmd.Execute(); err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}
