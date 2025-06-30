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

type RedisConfig struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Password     string `yaml:"password" mapstructure:"password"`
	DB           int    `yaml:"db" mapstructure:"db"`
	PoolSize     int    `yaml:"pool_size" mapstructure:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
	DialTimeout  int    `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout  int    `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" mapstructure:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

type KafkaConfig struct {
	Brokers                   []string `yaml:"brokers" mapstructure:"brokers"`
	ClientID                  string   `yaml:"client_id" mapstructure:"client_id"`
	Version                   string   `yaml:"version" mapstructure:"version"`
	Username                  string   `yaml:"username" mapstructure:"username"`
	Password                  string   `yaml:"password" mapstructure:"password"`
	EnableSASL                bool     `yaml:"enable_sasl" mapstructure:"enable_sasl"`
	SASLMechanism             string   `yaml:"sasl_mechanism" mapstructure:"sasl_mechanism"`
	EnableTLS                 bool     `yaml:"enable_tls" mapstructure:"enable_tls"`
	ProducerReturnSuccesses   bool     `yaml:"producer_return_successes" mapstructure:"producer_return_successes"`
	ProducerReturnErrors      bool     `yaml:"producer_return_errors" mapstructure:"producer_return_errors"`
	ProducerRequiredAcks      int      `yaml:"producer_required_acks" mapstructure:"producer_required_acks"`
	ProducerRetryMax          int      `yaml:"producer_retry_max" mapstructure:"producer_retry_max"`
	ProducerMaxMessageBytes   int      `yaml:"producer_max_message_bytes" mapstructure:"producer_max_message_bytes"`
	ConsumerGroupID           string   `yaml:"consumer_group_id" mapstructure:"consumer_group_id"`
	ConsumerOffsetInitial     string   `yaml:"consumer_offset_initial" mapstructure:"consumer_offset_initial"`
	ConsumerSessionTimeout    int      `yaml:"consumer_session_timeout" mapstructure:"consumer_session_timeout"`
	ConsumerHeartbeatInterval int      `yaml:"consumer_heartbeat_interval" mapstructure:"consumer_heartbeat_interval"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Logging  LoggingConfig  `yaml:"logging" mapstructure:"logging"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka" mapstructure:"kafka"`
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

	// Redis配置环境变量绑定
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("redis.pool_size", "REDIS_POOL_SIZE")
	viper.BindEnv("redis.min_idle_conns", "REDIS_MIN_IDLE_CONNS")
	viper.BindEnv("redis.dial_timeout", "REDIS_DIAL_TIMEOUT")
	viper.BindEnv("redis.read_timeout", "REDIS_READ_TIMEOUT")
	viper.BindEnv("redis.write_timeout", "REDIS_WRITE_TIMEOUT")
	viper.BindEnv("redis.idle_timeout", "REDIS_IDLE_TIMEOUT")

	// Kafka配置环境变量绑定
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.client_id", "KAFKA_CLIENT_ID")
	viper.BindEnv("kafka.version", "KAFKA_VERSION")
	viper.BindEnv("kafka.username", "KAFKA_USERNAME")
	viper.BindEnv("kafka.password", "KAFKA_PASSWORD")
	viper.BindEnv("kafka.enable_sasl", "KAFKA_ENABLE_SASL")
	viper.BindEnv("kafka.sasl_mechanism", "KAFKA_SASL_MECHANISM")
	viper.BindEnv("kafka.enable_tls", "KAFKA_ENABLE_TLS")
	viper.BindEnv("kafka.producer_return_successes", "KAFKA_PRODUCER_RETURN_SUCCESSES")
	viper.BindEnv("kafka.producer_return_errors", "KAFKA_PRODUCER_RETURN_ERRORS")
	viper.BindEnv("kafka.producer_required_acks", "KAFKA_PRODUCER_REQUIRED_ACKS")
	viper.BindEnv("kafka.producer_retry_max", "KAFKA_PRODUCER_RETRY_MAX")
	viper.BindEnv("kafka.producer_max_message_bytes", "KAFKA_PRODUCER_MAX_MESSAGE_BYTES")
	viper.BindEnv("kafka.consumer_group_id", "KAFKA_CONSUMER_GROUP_ID")
	viper.BindEnv("kafka.consumer_offset_initial", "KAFKA_CONSUMER_OFFSET_INITIAL")
	viper.BindEnv("kafka.consumer_session_timeout", "KAFKA_CONSUMER_SESSION_TIMEOUT")
	viper.BindEnv("kafka.consumer_heartbeat_interval", "KAFKA_CONSUMER_HEARTBEAT_INTERVAL")
}

// GetDsn 获取数据库连接字符串
func (d *DatabaseConfig) GetDsn() string {
	switch d.Type {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Password, d.DBName)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			d.Username, d.Password, d.Host, d.Port, d.DBName, d.Charset, d.ParseTime, d.Loc)
	default:
		return ""
	}
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
