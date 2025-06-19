package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"learngo0619/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	formatFlag string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置文件管理",
	Long: `配置文件管理工具，支持配置验证、显示和生成。

可用的子命令:
  validate  验证配置文件格式
  show      显示当前配置内容  
  generate  生成示例配置文件`,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证配置文件",
	Long:  `验证指定的配置文件格式是否正确`,
	Run:   runValidate,
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "显示配置内容",
	Long:  `显示当前配置文件的内容，支持YAML和JSON格式`,
	Run:   runShow,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成示例配置",
	Long:  `生成一个示例配置文件模板`,
	Run:   runGenerate,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(validateCmd)
	configCmd.AddCommand(showCmd)
	configCmd.AddCommand(generateCmd)

	showCmd.Flags().StringVarP(&formatFlag, "format", "f", "yaml", "输出格式 (yaml|json)")
}

func runValidate(cmd *cobra.Command, args []string) {
	cfgFile, _ := cmd.Flags().GetString("config")
	if cfgFile == "" {
		cfgFile = "config/server.yaml"
	}

	_, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("❌ 配置文件验证失败: %v", err)
	}

	fmt.Printf("✅ 配置文件 %s 验证通过\n", cfgFile)
}

func runShow(cmd *cobra.Command, args []string) {
	cfgFile, _ := cmd.Flags().GetString("config")
	if cfgFile == "" {
		cfgFile = "config/server.yaml"
	}

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Fatalf("❌ 读取配置文件失败: %v", err)
	}

	switch formatFlag {
	case "json":
		jsonData, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			log.Fatalf("❌ 转换JSON失败: %v", err)
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatalf("❌ 转换YAML失败: %v", err)
		}
		fmt.Print(string(yamlData))
	default:
		log.Fatalf("❌ 不支持的格式: %s", formatFlag)
	}
}

func runGenerate(cmd *cobra.Command, args []string) {
	sampleConfig := &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
			Mode: "debug",
		},
		App: config.AppConfig{
			Name:    "myapp",
			Version: "1.0.0",
		},
	}

	yamlData, err := yaml.Marshal(sampleConfig)
	if err != nil {
		log.Fatalf("❌ 生成配置失败: %v", err)
	}

	fmt.Println("# 示例配置文件")
	fmt.Print(string(yamlData))
}
