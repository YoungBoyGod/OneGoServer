package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// 第一步：定义配置结构体
// 这个结构体告诉程序：YAML文件里有什么内容
type Config struct {
	Server ServerConfig `yaml:"server"` // yaml:"server" 表示对应YAML中的server部分
	App    AppConfig    `yaml:"app"`    // yaml:"app" 表示对应YAML中的app部分
}

type ServerConfig struct {
	Port string `yaml:"port"` // 服务器端口
	Mode string `yaml:"mode"` // 运行模式：debug 或 release
}

type AppConfig struct {
	Name    string `yaml:"name"`    // 应用名称
	Version string `yaml:"version"` // 应用版本
}

// 第二步：读取配置文件的函数
func LoadConfig(filename string) (*Config, error) {
	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %v", filename, err)
	}

	// 创建配置对象
	var config Config

	// 将YAML内容解析到结构体中
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("无法解析YAML配置: %v", err)
	}

	fmt.Printf("✅ 成功加载配置文件: %s\n", filename)
	return &config, nil
}

// 第三步：创建Gin服务器
func CreateServer(config *Config) *gin.Engine {
	// 根据配置设置Gin模式
	if config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("🚀 Gin运行在发布模式")
	} else {
		gin.SetMode(gin.DebugMode)
		fmt.Println("🐛 Gin运行在调试模式")
	}

	// 创建Gin路由器
	router := gin.Default()

	// 添加首页路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "欢迎使用我的第一个Web应用！",
			"app":     config.App.Name,
			"version": config.App.Version,
			"tip":     "这是一个使用Cobra+Gin+YAML构建的示例",
		})
	})

	// 添加ping路由 - 用于测试服务器是否正常工作
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "服务器正常运行",
		})
	})

	// 添加配置信息路由 - 查看当前配置
	router.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"server": config.Server,
			"app":    config.App,
		})
	})

	// 添加帮助路由 - 显示可用的API
	router.GET("/help", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "可用的API接口",
			"apis": []map[string]string{
				{"path": "/", "description": "首页，显示应用信息"},
				{"path": "/ping", "description": "测试服务器状态"},
				{"path": "/config", "description": "查看配置信息"},
				{"path": "/help", "description": "显示帮助信息"},
			},
		})
	})

	return router
}

// 第四步：全局变量存储命令行参数
var configFile string

// 第五步：定义Cobra命令
var rootCmd = &cobra.Command{
	Use:   "simple-app",                  // 程序名称
	Short: "我的第一个Web应用",                  // 简短描述
	Long:  "这是一个学习Cobra+Gin+YAML的简单示例应用", // 详细描述
}

var serverCmd = &cobra.Command{
	Use:   "server",                // 子命令名称
	Short: "启动Web服务器",              // 简短描述
	Long:  "启动HTTP服务器，提供Web API接口", // 详细描述
	Run: func(cmd *cobra.Command, args []string) {
		// 当用户执行 "./simple-app server" 时，这里的代码会运行
		startServer()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Simple App v1.0.0")
		fmt.Println("使用 Cobra + Gin + YAML 构建")
	},
}

// 第六步：初始化函数 - 程序启动时自动执行
func init() {
	// 添加子命令到根命令
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)

	// 为server命令添加配置文件参数
	// 用户可以使用 --config 或 -c 来指定配置文件
	serverCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "配置文件路径")
}

// 第七步：启动服务器的函数
func startServer() {
	fmt.Println("🚀 正在启动服务器...")

	// 1. 加载配置文件
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 2. 创建Gin服务器
	router := CreateServer(config)

	// 3. 打印启动信息
	fmt.Printf("\n🌟 %s 启动成功！\n", config.App.Name)
	fmt.Printf("📍 访问地址: http://localhost:%s\n", config.Server.Port)
	fmt.Printf("📡 测试接口: http://localhost:%s/ping\n", config.Server.Port)
	fmt.Printf("📋 帮助信息: http://localhost:%s/help\n", config.Server.Port)
	fmt.Printf("⚙️  配置信息: http://localhost:%s/config\n", config.Server.Port)
	fmt.Printf("🛑 按 Ctrl+C 停止服务器\n\n")

	// 4. 启动HTTP服务器
	err = router.Run(":" + config.Server.Port)
	if err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}

// 第八步：主函数 - 程序的入口点
func main() {
	// 执行Cobra命令解析
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
