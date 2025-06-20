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
	Log    LogConfig    `yaml:"log"` // 新增日志配置
}

// ServerConfig is the configuration for the server
type ServerConfig struct {
	Host     string       `yaml:"host"`
	Port     string       `yaml:"port"`
	Mode     string       `yaml:"mode"`
	Server   ServerInfo   `yaml:"server"`
	Database DatabaseInfo `yaml:"database"`
	Log      LogInfo      `yaml:"log"`
	Security SecurityInfo `yaml:"security"` // 新增安全配置
}

// AppConfig is the configuration for the app
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// LogConfig is the configuration for logging
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别: debug, info, warn, error
	Format     string `yaml:"format"`      // 输出格式: json, console
	Output     string `yaml:"output"`      // 输出方式: stdout, file, both
	Dir        string `yaml:"dir"`         // 日志目录
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxAge     int    `yaml:"max_age"`     // 日志文件最长保存时间(天)
	MaxBackups int    `yaml:"max_backups"` // 最大备份文件数量
	Compress   bool   `yaml:"compress"`    // 是否压缩归档
}

// ServerInfo 服务器信息
type ServerInfo struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseInfo 数据库信息
type DatabaseInfo struct {
	// Add any necessary fields for database configuration
}

// LogInfo 日志信息
type LogInfo struct {
	// Add any necessary fields for log configuration
}

// SecurityInfo 安全配置（简化版）
type SecurityInfo struct {
	PredefinedToken string `yaml:"predefined_token"` // 预定义注册Token（唯一认证方式）
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
	// 4. 设置默认值
	config.setDefaults()
	// 5. 返回配置对象
	return &config, nil
}

// setDefaults 设置默认配置值
func (c *Config) setDefaults() {
	// 日志默认配置
	if c.Log.Level == "" {
		if c.Server.Mode == "release" {
			c.Log.Level = "info"
		} else {
			c.Log.Level = "debug"
		}
	}
	if c.Log.Format == "" {
		if c.Server.Mode == "release" {
			c.Log.Format = "json"
		} else {
			c.Log.Format = "console"
		}
	}
	if c.Log.Output == "" {
		if c.Server.Mode == "release" {
			c.Log.Output = "file"
		} else {
			c.Log.Output = "both"
		}
	}
	if c.Log.Dir == "" {
		c.Log.Dir = "logs"
	}
	if c.Log.MaxSize == 0 {
		c.Log.MaxSize = 100 // 100MB
	}
	if c.Log.MaxAge == 0 {
		c.Log.MaxAge = 30 // 30天
	}
	if c.Log.MaxBackups == 0 {
		c.Log.MaxBackups = 10 // 10个备份
	}
	// 生产环境默认压缩
	if c.Server.Mode == "release" {
		c.Log.Compress = true
	}
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

	// 日志配置覆盖
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.Log.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		c.Log.Format = format
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		c.Log.Output = output
	}
	if dir := os.Getenv("LOG_DIR"); dir != "" {
		c.Log.Dir = dir
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

	// 验证日志配置
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLevels, c.Log.Level) {
		return fmt.Errorf("invalid log level: %s, must be one of %v", c.Log.Level, validLevels)
	}

	validFormats := []string{"json", "console"}
	if !contains(validFormats, c.Log.Format) {
		return fmt.Errorf("invalid log format: %s, must be one of %v", c.Log.Format, validFormats)
	}

	validOutputs := []string{"stdout", "file", "both"}
	if !contains(validOutputs, c.Log.Output) {
		return fmt.Errorf("invalid log output: %s, must be one of %v", c.Log.Output, validOutputs)
	}

	return nil
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
