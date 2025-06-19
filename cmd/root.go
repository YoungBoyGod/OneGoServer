package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	// 全局配置文件路径
	ConfigFile string
)

// RootCmd 根命令
var RootCmd = &cobra.Command{
	Use:   "OneGoTask",
	Short: "OneGoTask 服务管理工具",
	Long:  "OneGoTask 服务端和客户端管理工具，提供启动、管理和监控功能",
}

// Execute 执行根命令
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("❌ 命令执行失败: %v", err)
	}
}

func init() {
	// 全局标志 - 配置文件路径
	RootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "config/server.yaml", "指定配置文件路径")

	// 添加子命令
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(clientCmd)
}
