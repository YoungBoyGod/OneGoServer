package cmd

import (
	"fmt"
	"log"

	"learngo0619/internal/config"
	"learngo0619/internal/server"

	"github.com/spf13/cobra"
)

var (
	port    string
	host    string
	verbose bool
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP Web服务器",
	Long: `启动一个HTTP Web服务器，支持：
• RESTful API接口
• 优雅关闭机制
• 健康检查接口
• 配置文件管理
• 热重载开发模式`,
	Run: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&port, "port", "p", "", "服务器端口 (覆盖配置文件)")
	serverCmd.Flags().StringVar(&host, "host", "", "服务器主机地址 (覆盖配置文件)")
	serverCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "详细输出模式")
}

func runServer(cmd *cobra.Command, args []string) {
	// 获取全局配置文件路径
	cfgFile, _ := cmd.Flags().GetString("config")
	if cfgFile == "" {
		cfgFile = "config/server.yaml"
	}

	// 获取全局verbose标志
	verboseGlobal, _ := cmd.Flags().GetBool("verbose")

	// 加载配置
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 从命令行参数覆盖配置
	if port != "" {
		cfg.Server.Port = port
	}
	if host != "" {
		cfg.Server.Host = host
	}

	if verboseGlobal || verbose {
		fmt.Printf("🔧 正在启动服务器...\n")
		fmt.Printf("📝 配置文件: %s\n", cfgFile)
	}

	// 创建和启动服务器
	srv := server.NewServer(cfg)
	if err := srv.Start(verboseGlobal || verbose); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
