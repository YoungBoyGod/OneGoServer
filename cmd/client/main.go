package main

import (
	"log"

	"OneGoTask/pkg/client"
	"OneGoTask/pkg/config"

	"github.com/spf13/cobra"
)

var (
	serverURL    string
	outputFormat string
	timeout      int
	configFile   string
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "OneGoTask 客户端工具",
	Long:  "OneGoTask 客户端工具，用于与服务器进行交互",
}

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "测试服务器连接",
	Long:  "发送ping请求到服务器，测试连接状态",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("ping")
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "检查服务器健康状态",
	Long:  "获取服务器的健康状态信息",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("health")
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "获取服务器配置",
	Long:  "获取服务器的当前配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("config")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取服务器状态",
	Long:  "获取服务器的运行状态信息",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("status")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "获取服务器版本",
	Long:  "获取服务器的版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("version")
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "重启服务器",
	Long:  "发送重启请求到服务器",
	Run: func(cmd *cobra.Command, args []string) {
		runClientCommand("restart")
	},
}

func init() {
	// 添加子命令
	clientCmd.AddCommand(pingCmd)
	clientCmd.AddCommand(healthCmd)
	clientCmd.AddCommand(configCmd)
	clientCmd.AddCommand(statusCmd)
	clientCmd.AddCommand(versionCmd)
	clientCmd.AddCommand(restartCmd)

	// 为 client 命令添加全局标志
	clientCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "", "服务器地址 (例如: http://localhost:30000)")
	clientCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "输出格式 (json, yaml, table)")
	clientCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 0, "请求超时时间(秒)")
	clientCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/server.yaml", "指定配置文件路径")
}

func runClientCommand(command string) {
	// 加载配置文件
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

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

func main() {
	// 执行 Cobra 命令
	if err := clientCmd.Execute(); err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}
