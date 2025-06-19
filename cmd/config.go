package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"learngo0619/server"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd 代表config命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置文件管理",
	Long: `配置文件管理工具，支持以下操作:

• 验证配置文件格式
• 显示当前配置
• 生成示例配置文件

子命令:
  validate    验证配置文件
  show        显示配置内容
  generate    生成示例配置文件`,
}

// validateCmd 验证配置文件
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证配置文件格式",
	Long:  `验证指定的配置文件格式是否正确，检查必需字段是否存在。`,
	Run:   validateConfig,
}

// showCmd 显示配置文件内容
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "显示配置文件内容",
	Long:  `以格式化的方式显示配置文件内容，支持JSON和YAML格式输出。`,
	Run:   showConfig,
}

// generateCmd 生成示例配置文件
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成示例配置文件",
	Long:  `生成一个包含所有选项的示例配置文件。`,
	Run:   generateConfig,
}

func validateConfig(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")

	fmt.Printf("🔍 验证配置文件: %s\n", configFile)

	config, err := server.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("❌ 配置文件验证失败: %v\n", err)
		return
	}

	// 基本验证
	if config.Server.Host == "" {
		fmt.Println("❌ 服务器主机地址不能为空")
		return
	}

	if config.Server.Port == "" {
		fmt.Println("❌ 服务器端口不能为空")
		return
	}

	if config.App.Name == "" {
		fmt.Println("❌ 应用名称不能为空")
		return
	}

	fmt.Println("✅ 配置文件格式正确")
	fmt.Printf("📋 服务器: %s:%s\n", config.Server.Host, config.Server.Port)
	fmt.Printf("📱 应用: %s v%s\n", config.App.Name, config.App.Version)
}

func showConfig(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")
	format, _ := cmd.Flags().GetString("format")

	config, err := server.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("📄 配置文件: %s\n", configFile)
	fmt.Println(strings.Repeat("-", 40))

	switch format {
	case "json":
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal config to JSON: %v", err)
		}
		fmt.Println(string(data))
	case "yaml":
		data, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("Failed to marshal config to YAML: %v", err)
		}
		fmt.Println(string(data))
	default:
		fmt.Printf("服务器配置:\n")
		fmt.Printf("  主机: %s\n", config.Server.Host)
		fmt.Printf("  端口: %s\n", config.Server.Port)
		fmt.Printf("  模式: %s\n", config.Server.Mode)
		fmt.Printf("\n应用配置:\n")
		fmt.Printf("  名称: %s\n", config.App.Name)
		fmt.Printf("  版本: %s\n", config.App.Version)
	}
}

func generateConfig(cmd *cobra.Command, args []string) {
	output, _ := cmd.Flags().GetString("output")

	// 示例配置
	exampleConfig := &server.Config{
		Server: server.ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
			Mode: "debug",
		},
		App: server.AppConfig{
			Name:    "myapp",
			Version: "1.0.0",
		},
	}

	data, err := yaml.Marshal(exampleConfig)
	if err != nil {
		log.Fatalf("Failed to generate config: %v", err)
	}

	if output != "" {
		// 写入文件
		err := os.WriteFile(output, data, 0644)
		if err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
		fmt.Printf("✅ 示例配置文件已生成: %s\n", output)
	} else {
		// 输出到控制台
		fmt.Println("# 示例配置文件")
		fmt.Println(string(data))
	}
}

func init() {
	// 为show命令添加格式选项
	showCmd.Flags().StringP("format", "f", "pretty", "输出格式 (pretty, json, yaml)")

	// 为generate命令添加输出选项
	generateCmd.Flags().StringP("output", "o", "", "输出文件路径 (默认输出到控制台)")

	// 添加子命令到config命令
	configCmd.AddCommand(validateCmd)
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(generateCmd)

	// 将config命令添加到根命令
	rootCmd.AddCommand(configCmd)
}
