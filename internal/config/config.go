package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the main configuration struct
type Config struct {
	Server ServerConfig `yaml:"server"`
	App    AppConfig    `yaml:"app"`
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// AppConfig is the configuration for the app
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// LoadConfig reads config from file
func LoadConfig(filename string) (*Config, error) {
	// 1. 读取文件内容
	configData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 2. 创建配置对象
	var config Config
	// 3. 解析配置对象
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	// 4. 返回配置对象
	return &config, nil
}

// LoadConfigByEnv 根据环境变量加载配置
func LoadConfigByEnv(env string) (*Config, error) {
	// 如果环境为空，尝试从环境变量获取
	if env == "" {
		env = os.Getenv("APP_ENV")
	}

	// 如果仍为空，默认使用 "dev"
	if env == "" {
		env = "dev"
	}

	// 构建配置文件路径
	var filename string
	switch strings.ToLower(env) {
	case "dev", "development":
		filename = "config/server.dev.yaml"
	case "prod", "production":
		filename = "config/server.prod.yaml"
	case "test", "testing":
		filename = "config/server.test.yaml"
	default:
		// 如果是自定义环境，尝试查找对应文件
		filename = fmt.Sprintf("config/server.%s.yaml", env)
		// 如果文件不存在，回退到默认配置
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			filename = "config/server.yaml"
		}
	}

	return LoadConfig(filename)
}

// GetEnvironment 获取当前环境
func GetEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "dev"
	}
	return strings.ToLower(env)
}

// IsProduction 检查是否为生产环境
func IsProduction() bool {
	env := GetEnvironment()
	return env == "prod" || env == "production"
}

// IsDevelopment 检查是否为开发环境
func IsDevelopment() bool {
	env := GetEnvironment()
	return env == "dev" || env == "development"
}

// ApplyEnvOverrides 应用环境变量覆盖
func (c *Config) ApplyEnvOverrides() {
	// 服务器配置覆盖
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		c.Server.Port = port
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		c.Server.Mode = mode
	}

	// 应用配置覆盖
	if name := os.Getenv("APP_NAME"); name != "" {
		c.App.Name = name
	}
	if version := os.Getenv("APP_VERSION"); version != "" {
		c.App.Version = version
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.App.Name == "" {
		return fmt.Errorf("app name is required")
	}
	if c.App.Version == "" {
		return fmt.Errorf("app version is required")
	}
	return nil
}
