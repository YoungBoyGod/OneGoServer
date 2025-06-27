package server

import (
	"log"

	"OneGoTask/pkg/config"
	"OneGoTask/pkg/server"

	"github.com/spf13/cobra"
)

// GetServerCmd 返回服务器命令
func GetServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "启动 HTTP 服务器",
		Long:  "启动 OneGoTask HTTP 服务器，提供 API 接口，支持优雅启动和关闭",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果直接执行 server 命令，默认启动服务器
			port, _ := cmd.Flags().GetString("port")
			configFile, _ := cmd.Flags().GetString("config")
			startServer(port, configFile)
		},
	}

	var serverStartCmd = &cobra.Command{
		Use:   "start",
		Short: "启动服务器",
		Long:  "启动 OneGoTask HTTP 服务器",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")
			configFile, _ := cmd.Flags().GetString("config")
			startServer(port, configFile)
		},
	}

	// 添加子命令
	serverCmd.AddCommand(serverStartCmd)

	// 为主命令和子命令都添加端口参数
	serverCmd.Flags().StringP("port", "p", "", "指定服务器监听端口")
	serverStartCmd.Flags().StringP("port", "p", "", "指定服务器监听端口")

	return serverCmd
}

// startServer 启动服务器的函数
func startServer(port, configFile string) {
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
