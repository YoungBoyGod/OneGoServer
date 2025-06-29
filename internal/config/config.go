package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Name        string `yaml:"name" mapstructure:"name"`
	Mode        string `yaml:"mode" mapstructure:"mode"`
	Version     string `yaml:"version" mapstructure:"version"`
	Description string `yaml:"description" mapstructure:"description"`
	Author      string `yaml:"author" mapstructure:"author"`
	Email       string `yaml:"email" mapstructure:"email"`
	Url         string `yaml:"url" mapstructure:"url"`
	Copyright   string `yaml:"copyright" mapstructure:"copyright"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string        `yaml:"type" mapstructure:"type"`
	Host            string        `yaml:"host" mapstructure:"host"`
	Port            int           `yaml:"port" mapstructure:"port"`
	Username        string        `yaml:"username" mapstructure:"username"`
	Password        string        `yaml:"password" mapstructure:"password"`
	DBName          string        `yaml:"dbname" mapstructure:"dbname"`
	Charset         string        `yaml:"charset" mapstructure:"charset"`
	ParseTime       bool          `yaml:"parse_time" mapstructure:"parse_time"`
	Loc             string        `yaml:"loc" mapstructure:"loc"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`
	Format     string `yaml:"format" mapstructure:"format"`
	Output     string `yaml:"output" mapstructure:"output"`
	FilePath   string `yaml:"file_path" mapstructure:"file_path"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"`
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"`
	Compress   bool   `yaml:"compress" mapstructure:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret" mapstructure:"secret"`
	ExpireHours int    `yaml:"expire_hours" mapstructure:"expire_hours"`
	Issuer      string `yaml:"issuer" mapstructure:"issuer"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins  []string `yaml:"allow_origins" mapstructure:"allow_origins"`
	AllowMethods  []string `yaml:"allow_methods" mapstructure:"allow_methods"`
	AllowHeaders  []string `yaml:"allow_headers" mapstructure:"allow_headers"`
	ExposeHeaders []string `yaml:"expose_headers" mapstructure:"expose_headers"`
	MaxAge        int      `yaml:"max_age" mapstructure:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled" mapstructure:"enabled"`
	RequestsPerMinute int  `yaml:"requests_per_minute" mapstructure:"requests_per_minute"`
	Burst             int  `yaml:"burst" mapstructure:"burst"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	CORS      CORSConfig      `yaml:"cors" mapstructure:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit" mapstructure:"rate_limit"`
}

// UserBusinessConfig 用户业务配置
type UserBusinessConfig struct {
	PasswordMinLength      int `yaml:"password_min_length" mapstructure:"password_min_length"`
	UsernameMinLength      int `yaml:"username_min_length" mapstructure:"username_min_length"`
	MaxLoginAttempts       int `yaml:"max_login_attempts" mapstructure:"max_login_attempts"`
	AccountLockoutDuration int `yaml:"account_lockout_duration" mapstructure:"account_lockout_duration"`
}

// BusinessConfig 业务配置
type BusinessConfig struct {
	User UserBusinessConfig `yaml:"user" mapstructure:"user"`
}

// Config 总配置
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Logging  LoggingConfig  `yaml:"logging" mapstructure:"logging"`
	JWT      JWTConfig      `yaml:"jwt" mapstructure:"jwt"`
	Security SecurityConfig `yaml:"security" mapstructure:"security"`
	Business BusinessConfig `yaml:"business" mapstructure:"business"`
}

// LoadConfig 加载配置
func LoadConfig(path string) (*Config, error) {
	// 加载.env文件（如果存在）
	if err := godotenv.Load(); err != nil {
		// .env文件不存在或加载失败时不报错，使用系统环境变量
		fmt.Println("Warning: .env file not found or failed to load, using system environment variables")
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 自动读取环境变量
	viper.AutomaticEnv()

	// 设置环境变量前缀
	viper.SetEnvPrefix("ONEGO")

	// 绑定环境变量
	bindEnvironmentVariables()

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖敏感配置
	overrideFromEnv(&cfg)

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// bindEnvironmentVariables 绑定环境变量
func bindEnvironmentVariables() {
	// 数据库相关
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.username", "DB_USERNAME")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")

	// JWT相关
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	// 服务器相关（可选）
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
}

// overrideFromEnv 从环境变量覆盖敏感配置
func overrideFromEnv(cfg *Config) {
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p := viper.GetInt("DB_PORT"); p != 0 {
			cfg.Database.Port = p
		}
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		cfg.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}

	// 可选的服务器配置
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := viper.GetInt("SERVER_PORT"); port != 0 {
		cfg.Server.Port = port
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		cfg.Server.Mode = mode
	}
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "dev")

	// 数据库默认配置
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.dbname", "onego_server")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	// 日志默认配置
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// JWT默认配置
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.issuer", "OneGoServer")

	// 业务默认配置
	viper.SetDefault("business.user.password_min_length", 6)
	viper.SetDefault("business.user.username_min_length", 3)
}

// validate 验证配置
func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	if c.Business.User.PasswordMinLength < 4 {
		return fmt.Errorf("password minimum length cannot be less than 4")
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	switch c.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
			c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			c.Host, c.Username, c.Password, c.DBName, c.Port, c.Loc)
	case "sqlite":
		if c.DBName == ":memory:" {
			return ":memory:"
		}
		return c.DBName + ".db"
	default:
		return ""
	}
}
