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
	env     string // 新增环境参数
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动HTTP Web服务器",
	Long: `启动一个HTTP Web服务器，支持：
• RESTful API接口
• 优雅关闭机制
• 健康检查接口
• 多环境配置管理
• 热重载开发模式

环境配置:
  --env dev     使用开发环境配置 (config/server.dev.yaml)
  --env prod    使用生产环境配置 (config/server.prod.yaml)
  --env test    使用测试环境配置 (config/server.test.yaml)
  
环境变量:
  APP_ENV=prod  设置环境 (优先级低于 --env 参数)
  SERVER_PORT   覆盖服务器端口
  SERVER_HOST   覆盖服务器主机地址`,
	Run: runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&port, "port", "p", "", "服务器端口 (覆盖配置文件)")
	serverCmd.Flags().StringVar(&host, "host", "", "服务器主机地址 (覆盖配置文件)")
	serverCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "详细输出模式")
	serverCmd.Flags().StringVarP(&env, "env", "e", "", "环境 (dev|prod|test)")
}

func runServer(cmd *cobra.Command, args []string) {
	var cfg *config.Config
	var err error

	// 优先级: --config > --env > APP_ENV > 默认(dev)
	if configFile, _ := cmd.Flags().GetString("config"); configFile != "config/server.yaml" {
		// 如果明确指定了配置文件，直接使用
		cfg, err = config.LoadConfig(configFile)
		if err != nil {
			log.Fatalf("❌ 加载配置文件失败: %v", err)
		}
	} else {
		// 使用环境配置
		cfg, err = config.LoadConfigByEnv(env)
		if err != nil {
			log.Fatalf("❌ 加载环境配置失败: %v", err)
		}
	}

	// 应用环境变量覆盖
	cfg.ApplyEnvOverrides()

	// 从命令行参数覆盖配置
	if port != "" {
		cfg.Server.Port = port
	}
	if host != "" {
		cfg.Server.Host = host
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ 配置验证失败: %v", err)
	}

	// 获取全局verbose标志
	verboseGlobal, _ := cmd.Flags().GetBool("verbose")

	if verboseGlobal || verbose {
		fmt.Printf("🔧 正在启动服务器...\n")
		fmt.Printf("🌍 当前环境: %s\n", config.GetEnvironment())
		if config.IsProduction() {
			fmt.Printf("⚠️  生产环境模式\n")
		} else {
			fmt.Printf("🔧 开发环境模式\n")
		}
		fmt.Printf("📝 配置信息:\n")
		fmt.Printf("   主机: %s\n", cfg.Server.Host)
		fmt.Printf("   端口: %s\n", cfg.Server.Port)
		fmt.Printf("   模式: %s\n", cfg.Server.Mode)
		fmt.Printf("   应用: %s v%s\n", cfg.App.Name, cfg.App.Version)
	}

	// 创建和启动服务器
	srv := server.NewServer(cfg)
	if err := srv.Start(verboseGlobal || verbose); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
