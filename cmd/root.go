package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// 版本信息
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

// rootCmd 代表应用程序的基础命令
var rootCmd = &cobra.Command{
	Use:   "learngo0619",
	Short: "一个现代化的Go Web应用程序",
	Long: `learngo0619 是一个功能完整的Go Web应用程序，具备以下特性:

• HTTP服务器支持
• 优雅关闭机制
• 配置文件管理
• 热重载开发模式
• RESTful API接口

使用不同的子命令来执行不同的操作。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 如果没有指定子命令，显示帮助信息
		cmd.Help()
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令时出错: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 添加全局标志
	rootCmd.PersistentFlags().StringP("config", "c", "config/server.yaml", "配置文件路径")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出模式")
}
