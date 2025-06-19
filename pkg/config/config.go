package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 主配置结构体
type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
	Client ClientConfig `yaml:"client"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port            string `yaml:"port"`
	Mode            string `yaml:"mode"`
	Host            string `yaml:"host"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // 优雅关闭超时时间(秒)
	ReadTimeout     int    `yaml:"read_timeout"`     // 读取超时时间(秒)
	WriteTimeout    int    `yaml:"write_timeout"`    // 写入超时时间(秒)
}

// ClientConfig 客户端配置
type ClientConfig struct {
	ServerURL    string `yaml:"server_url"`    // 服务器地址
	Timeout      int    `yaml:"timeout"`       // 请求超时时间(秒)
	RetryCount   int    `yaml:"retry_count"`   // 重试次数
	RetryDelay   int    `yaml:"retry_delay"`   // 重试延迟(秒)
	OutputFormat string `yaml:"output_format"` // 输出格式: json, table, yaml
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	config := &Config{}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("⚠️ 配置文件不存在: %s，使用默认配置", configPath)
		// 使用默认配置
		config = &Config{
			Server: ServerConfig{
				Port:            "30000",
				Mode:            "release",
				Host:            "0.0.0.0",
				ShutdownTimeout: 30,
				ReadTimeout:     30,
				WriteTimeout:    30,
			},
			Client: ClientConfig{
				ServerURL:    "http://localhost:30000",
				Timeout:      30,
				RetryCount:   3,
				RetryDelay:   1,
				OutputFormat: "json",
			},
			Log: LogConfig{
				Level: "info",
				File:  "logs/app.log",
			},
		}
		return config, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 YAML
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置默认值
	config.setDefaults()

	log.Printf("✅ 成功加载配置文件: %s", configPath)
	return config, nil
}

// setDefaults 设置默认值
func (c *Config) setDefaults() {
	// Server 默认值
	if c.Server.ShutdownTimeout == 0 {
		c.Server.ShutdownTimeout = 30
	}
	if c.Server.ReadTimeout == 0 {
		c.Server.ReadTimeout = 30
	}
	if c.Server.WriteTimeout == 0 {
		c.Server.WriteTimeout = 30
	}
	if c.Server.Port == "" {
		c.Server.Port = "30000"
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "release"
	}

	// Client 默认值
	if c.Client.ServerURL == "" {
		c.Client.ServerURL = "http://localhost:30000"
	}
	if c.Client.Timeout == 0 {
		c.Client.Timeout = 30
	}
	if c.Client.RetryCount == 0 {
		c.Client.RetryCount = 3
	}
	if c.Client.RetryDelay == 0 {
		c.Client.RetryDelay = 1
	}
	if c.Client.OutputFormat == "" {
		c.Client.OutputFormat = "json"
	}

	// Log 默认值
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.File == "" {
		c.Log.File = "logs/app.log"
	}
}

// MergeServerFlags 合并服务器命令行参数
func (c *Config) MergeServerFlags(port string) {
	// 命令行参数优先级高于配置文件
	if port != "" {
		c.Server.Port = port
		log.Printf("🔧 使用命令行指定的端口: %s", port)
	} else {
		log.Printf("📝 使用配置文件中的端口: %s", c.Server.Port)
	}

	// 验证端口号
	if portNum, err := strconv.Atoi(c.Server.Port); err != nil || portNum <= 0 || portNum > 65535 {
		log.Printf("⚠️ 无效的端口号: %s，使用默认端口: 30000", c.Server.Port)
		c.Server.Port = "30000"
	}
}

// MergeClientFlags 合并客户端命令行参数
func (c *Config) MergeClientFlags(serverURL, outputFormat string, timeout int) {
	if serverURL != "" {
		c.Client.ServerURL = serverURL
		log.Printf("🔧 使用命令行指定的服务器地址: %s", serverURL)
	}

	if outputFormat != "" {
		c.Client.OutputFormat = outputFormat
		log.Printf("🔧 使用命令行指定的输出格式: %s", outputFormat)
	}

	if timeout > 0 {
		c.Client.Timeout = timeout
		log.Printf("🔧 使用命令行指定的超时时间: %ds", timeout)
	}
}
