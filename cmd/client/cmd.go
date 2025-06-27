package client

import (
	"log"

	"OneGoTask/pkg/client"
	"OneGoTask/pkg/config"

	"github.com/spf13/cobra"
)

// GetClientCmd 返回客户端命令
func GetClientCmd() *cobra.Command {
	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "客户端工具",
		Long:  "OneGoTask 客户端工具，用于与服务器进行交互",
	}

	// ping 命令
	var pingCmd = &cobra.Command{
		Use:   "ping",
		Short: "测试服务器连接",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("ping", clientCmd)
		},
	}

	// health 命令
	var healthCmd = &cobra.Command{
		Use:   "health",
		Short: "检查服务器健康状态",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("health", clientCmd)
		},
	}

	// config 命令
	var configShowCmd = &cobra.Command{
		Use:   "config",
		Short: "获取服务器配置",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("config", clientCmd)
		},
	}

	// status 命令
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "获取服务器状态",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("status", clientCmd)
		},
	}

	// version 命令
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "获取服务器版本",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("version", clientCmd)
		},
	}

	// restart 命令
	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "重启服务器",
		Run: func(cmd *cobra.Command, args []string) {
			runClientCommand("restart", clientCmd)
		},
	}

	// 添加客户端子命令
	clientCmd.AddCommand(pingCmd)
	clientCmd.AddCommand(healthCmd)
	clientCmd.AddCommand(configShowCmd)
	clientCmd.AddCommand(statusCmd)
	clientCmd.AddCommand(versionCmd)
	clientCmd.AddCommand(restartCmd)

	// 为客户端命令添加全局标志
	clientCmd.PersistentFlags().StringP("server", "s", "", "服务器地址")
	clientCmd.PersistentFlags().StringP("output", "o", "", "输出格式 (json, yaml, table)")
	clientCmd.PersistentFlags().IntP("timeout", "t", 0, "请求超时时间(秒)")

	return clientCmd
}

// runClientCommand 执行客户端命令的函数
func runClientCommand(command string, cmd *cobra.Command) {
	// 获取配置文件路径
	configFile, _ := cmd.Flags().GetString("config")

	// 加载配置文件
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 获取客户端标志
	serverURL, _ := cmd.PersistentFlags().GetString("server")
	outputFormat, _ := cmd.PersistentFlags().GetString("output")
	timeout, _ := cmd.PersistentFlags().GetInt("timeout")

	// 合并命令行参数
	cfg.MergeClientFlags(serverURL, outputFormat, timeout)

	// 创建客户端实例
	cli := client.New(cfg)

	// 执行对应的命令
	switch command {
	case "ping":
		err = cli.Ping()
	case "health":
		err = cli.Health()
	case "config":
		err = cli.Config()
	case "status":
		err = cli.Status()
	case "version":
		err = cli.Version()
	case "restart":
		err = cli.Restart()
	default:
		log.Fatalf("❌ 未知命令: %s", command)
	}

	if err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}
