package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// 配置文件路径
	cfgFile string
	// 详细输出
	verbose bool
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "onego",
	Short: "OneGoServer - 分布式任务管理平台",
	Long: `OneGoServer 是一个使用Go语言开发的分布式任务管理平台，
采用DDD(领域驱动设计)架构，支持设备管理、任务调度和告警系统。

技术栈：
  • Web框架: Gin
  • 数据库: PostgreSQL + GORM  
  • 缓存: Redis
  • 消息队列: Kafka
  • 认证: JWT + bcrypt
  • 架构: DDD + Clean Architecture`,
	Example: `  # 启动服务器
  onego server

  # 启动服务器并指定配置文件
  onego server --config configs/config.production.yaml

  # 查看版本信息
  onego version

  # 显示详细帮助信息
  onego --help`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "配置文件路径")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细输出信息")

	// 设置版本信息
	rootCmd.Version = "v0.1.0"
	rootCmd.SetVersionTemplate(`{{.Version}}`)
}
