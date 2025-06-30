package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// 定义配置结构体
type DatabaseConfig struct {
	Type            string `yaml:"type" mapstructure:"type"`
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	DBName          string `yaml:"dbname" mapstructure:"dbname"`
	Charset         string `yaml:"charset" mapstructure:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime string `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	AutoMigrate     bool   `yaml:"auto_migrate" mapstructure:"auto_migrate"`
	ParseTime       bool   `yaml:"parse_time" mapstructure:"parse_time"`
	Loc             string `yaml:"loc" mapstructure:"loc"`
}

type ServerConfig struct {
	Host        string `yaml:"host" mapstructure:"host"`
	Port        int    `yaml:"port" mapstructure:"port"`
	Mode        string `yaml:"mode" mapstructure:"mode"`
	Name        string `yaml:"name" mapstructure:"name"`
	Version     string `yaml:"version" mapstructure:"version"`
	Description string `yaml:"description" mapstructure:"description"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
}

// 读取配置文件
func LoadConfig(path string) (*Config, error) {
	// 首先加载.env文件
	if err := godotenv.Load(".env"); err != nil {
		// .env文件不存在时不报错，继续执行
	}
	log.Printf("env: %+v", os.Getenv("DB_HOST"))
	log.Printf("env: %+v", os.Getenv("DB_PORT"))
	log.Printf("env: %+v", os.Getenv("DB_USERNAME"))
	log.Printf("env: %+v", os.Getenv("DB_PASSWORD"))
	log.Printf("env: %+v", os.Getenv("DB_NAME"))
	// 读取配置文件中的环境变量
	bindEnvironmentVariables()

	// 从配置文件中读取
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 将配置文件中的配置绑定到Config结构体中
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func bindEnvironmentVariables() {
	// 数据库配置环境变量绑定
	viper.BindEnv("database.type", "DB_TYPE")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.username", "DB_USERNAME")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.charset", "DB_CHARSET")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	viper.BindEnv("database.auto_migrate", "DB_AUTO_MIGRATE")

	// 服务器配置环境变量绑定
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")
}

// 校验配置参数
func (c *Config) Validate() error {
	// 验证数据库配置
	if c.Database.Host == "" {
		return errors.New("database host is required")
	}
	// 验证服务器端口
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server port is required and must be between 0 and 65535")
	}
	return nil
}

// GetDsn 获取数据库连接字符串
func (c *DatabaseConfig) GetDsn() string {
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
