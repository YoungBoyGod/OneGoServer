package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/go-redis/redis/v8"
)

// Client Redis客户端实例
var Client *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	// Redis配置
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Redis.IdleTimeout) * time.Minute,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	Client = rdb
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Health 检查Redis健康状态
func Health() error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return Client.Ping(ctx).Err()
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	return Client.Get(ctx, key).Result()
}

// GetObject 获取对象
func GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := Get(ctx, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// SetObject 设置对象
func SetObject(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return Set(ctx, key, data, expiration)
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis client not initialized")
	}

	count, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return Client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return Client.TTL(ctx, key).Result()
}

// Incr 原子递增
func Incr(ctx context.Context, key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return Client.Incr(ctx, key).Result()
}

// IncrBy 原子递增指定值
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return Client.IncrBy(ctx, key, value).Result()
}

// LPush 列表左推
func LPush(ctx context.Context, key string, values ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return Client.LPush(ctx, key, values...).Err()
}

// RPop 列表右弹
func RPop(ctx context.Context, key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	return Client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	if Client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	return Client.LLen(ctx, key).Result()
}

// SAdd 集合添加
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	if Client == nil {
		return fmt.Errorf("redis client not initialized")
	}

	return Client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	if Client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	return Client.SMembers(ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	if Client == nil {
		return false, fmt.Errorf("redis client not initialized")
	}

	return Client.SIsMember(ctx, key, member).Result()
}
