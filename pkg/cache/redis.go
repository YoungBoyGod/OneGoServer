package cache

// 读取配置文件中和env中redis的配置
// 初始化redis
// 返回redis的client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Client Redis的客户端实例
var RedisClient *redis.Client

// RedisManager Redis管理器
type RedisManager struct {
	client   *redis.Client
	mu       sync.RWMutex
	isReady  bool
	lastPing time.Time
}

// 全局Redis管理器实例
var (
	redisManager *RedisManager
	redisOnce    sync.Once
)

// RedisStats Redis统计信息
type RedisStats struct {
	TotalConns  uint32        `json:"total_conns"`
	IdleConns   uint32        `json:"idle_conns"`
	StaleConns  uint32        `json:"stale_conns"`
	Hits        uint64        `json:"hits"`
	Misses      uint64        `json:"misses"`
	Timeouts    uint64        `json:"timeouts"`
	LastPing    time.Time     `json:"last_ping"`
	IsConnected bool          `json:"is_connected"`
	Uptime      time.Duration `json:"uptime"`
}

// InitRedis 初始化Redis连接
func InitRedis(ctx context.Context, cfg *config.RedisConfig) error {
	// 验证配置（在redisOnce.Do外部进行，以支持重试）
	if err := validateRedisConfig(cfg); err != nil {
		return fmt.Errorf("invalid redis config: %w", err)
	}

	var initErr error

	redisOnce.Do(func() {

		// 创建Redis客户端
		rdb := redis.NewClient(&redis.Options{
			Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password:        cfg.Password,
			DB:              cfg.DB,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			DialTimeout:     time.Duration(cfg.DialTimeout) * time.Second,
			ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
			ConnMaxIdleTime: time.Duration(cfg.IdleTimeout) * time.Minute,
		})

		// 测试连接
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if _, err := rdb.Ping(pingCtx).Result(); err != nil {
			initErr = fmt.Errorf("failed to ping redis: %w", err)
			return
		}

		// 初始化管理器
		redisManager = &RedisManager{
			client:   rdb,
			isReady:  true,
			lastPing: time.Now(),
		}

		pkglog.LogInfo("Redis initialized successfully",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.Int("db", cfg.DB))
	})

	return initErr
}

// InitRedisWithRetry 带重试的Redis初始化
func InitRedisWithRetry(ctx context.Context, cfg *config.RedisConfig, maxRetries int, retryInterval time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		// 重置once以允许重试
		if i > 0 {
			redisOnce = sync.Once{}
		}

		if err := InitRedis(ctx, cfg); err != nil {
			lastErr = err
			pkglog.LogWarn("Redis connection failed, retrying",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))

			if i < maxRetries-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(retryInterval):
				}
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to initialize redis after %d attempts: %w", maxRetries, lastErr)
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	if redisManager == nil {
		return nil
	}
	redisManager.mu.RLock()
	defer redisManager.mu.RUnlock()
	return redisManager.client
}

// IsReady 检查Redis是否就绪
func IsReady() bool {
	if redisManager == nil {
		return false
	}
	redisManager.mu.RLock()
	defer redisManager.mu.RUnlock()
	return redisManager.isReady
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisManager == nil {
		return nil
	}

	redisManager.mu.Lock()
	defer redisManager.mu.Unlock()

	if redisManager.client != nil {
		if err := redisManager.client.Close(); err != nil {
			return fmt.Errorf("failed to close redis client: %w", err)
		}
		redisManager.isReady = false
		pkglog.LogInfo("Redis connection closed successfully")
	}

	return nil
}

// Ping Redis健康检查
func Ping(ctx context.Context) error {
	client := GetRedisClient()
	if client == nil {
		return errors.New("redis client not initialized")
	}

	startTime := time.Now()

	if _, err := client.Ping(ctx).Result(); err != nil {
		pkglog.LogRedisOperation("PING", "", time.Since(startTime), err)
		return fmt.Errorf("redis ping failed: %w", err)
	}

	// 更新最后ping时间
	if redisManager != nil {
		redisManager.mu.Lock()
		redisManager.lastPing = time.Now()
		redisManager.mu.Unlock()
	}

	pkglog.LogRedisOperation("PING", "", time.Since(startTime), nil)
	return nil
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	client := GetRedisClient()
	if client == nil {
		return errors.New("redis client not initialized")
	}

	startTime := time.Now()

	if err := client.Set(ctx, key, value, expiration).Err(); err != nil {
		pkglog.LogRedisOperation("SET", key, time.Since(startTime), err)
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("SET", key, time.Since(startTime), nil)
	return nil
}

// Get 获取键值对
func Get(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	if client == nil {
		return "", errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.Get(ctx, key).Result()
	if err != nil {
		pkglog.LogRedisOperation("GET", key, time.Since(startTime), err)
		if err == redis.Nil {
			return "", fmt.Errorf("key %s not found", key)
		}
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("GET", key, time.Since(startTime), nil)
	return result, nil
}

// Del 删除键
func Del(ctx context.Context, keys ...string) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.Del(ctx, keys...).Result()
	if err != nil {
		pkglog.LogRedisOperation("DEL", fmt.Sprintf("%v", keys), time.Since(startTime), err)
		return 0, fmt.Errorf("failed to delete keys %v: %w", keys, err)
	}

	pkglog.LogRedisOperation("DEL", fmt.Sprintf("%v", keys), time.Since(startTime), nil)
	return result, nil
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.Exists(ctx, keys...).Result()
	if err != nil {
		pkglog.LogRedisOperation("EXISTS", fmt.Sprintf("%v", keys), time.Since(startTime), err)
		return 0, fmt.Errorf("failed to check existence of keys %v: %w", keys, err)
	}

	pkglog.LogRedisOperation("EXISTS", fmt.Sprintf("%v", keys), time.Since(startTime), nil)
	return result, nil
}

// SetObject 设置对象（JSON序列化）
func SetObject(ctx context.Context, key string, obj interface{}, expiration time.Duration) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %w", err)
	}

	return Set(ctx, key, data, expiration)
}

// GetObject 获取对象（JSON反序列化）
func GetObject(ctx context.Context, key string, obj interface{}) error {
	data, err := Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), obj); err != nil {
		return fmt.Errorf("failed to unmarshal object: %w", err)
	}

	return nil
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	client := GetRedisClient()
	if client == nil {
		return errors.New("redis client not initialized")
	}

	startTime := time.Now()

	if err := client.Expire(ctx, key, expiration).Err(); err != nil {
		pkglog.LogRedisOperation("EXPIRE", key, time.Since(startTime), err)
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("EXPIRE", key, time.Since(startTime), nil)
	return nil
}

// TTL 获取过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.TTL(ctx, key).Result()
	if err != nil {
		pkglog.LogRedisOperation("TTL", key, time.Since(startTime), err)
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("TTL", key, time.Since(startTime), nil)
	return result, nil
}

// LPop 列表左弹
func LPop(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	if client == nil {
		return "", errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.LPop(ctx, key).Result()
	if err != nil {
		pkglog.LogRedisOperation("LPOP", key, time.Since(startTime), err)
		if err == redis.Nil {
			return "", fmt.Errorf("list %s is empty", key)
		}
		return "", fmt.Errorf("failed to lpop from key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("LPOP", key, time.Since(startTime), nil)
	return result, nil
}

// RPop 列表右弹
func RPop(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	if client == nil {
		return "", errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.RPop(ctx, key).Result()
	if err != nil {
		pkglog.LogRedisOperation("RPOP", key, time.Since(startTime), err)
		if err == redis.Nil {
			return "", fmt.Errorf("list %s is empty", key)
		}
		return "", fmt.Errorf("failed to rpop from key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("RPOP", key, time.Since(startTime), nil)
	return result, nil
}

// LPush 列表左推
func LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.LPush(ctx, key, values...).Result()
	if err != nil {
		pkglog.LogRedisOperation("LPUSH", key, time.Since(startTime), err)
		return 0, fmt.Errorf("failed to lpush to key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("LPUSH", key, time.Since(startTime), nil)
	return result, nil
}

// RPush 列表右推
func RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	client := GetRedisClient()
	if client == nil {
		return 0, errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.RPush(ctx, key, values...).Result()
	if err != nil {
		pkglog.LogRedisOperation("RPUSH", key, time.Since(startTime), err)
		return 0, fmt.Errorf("failed to rpush to key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("RPUSH", key, time.Since(startTime), nil)
	return result, nil
}

// HSet 哈希设置
func HSet(ctx context.Context, key string, values ...interface{}) error {
	client := GetRedisClient()
	if client == nil {
		return errors.New("redis client not initialized")
	}

	startTime := time.Now()

	if err := client.HSet(ctx, key, values...).Err(); err != nil {
		pkglog.LogRedisOperation("HSET", key, time.Since(startTime), err)
		return fmt.Errorf("failed to hset key %s: %w", key, err)
	}

	pkglog.LogRedisOperation("HSET", key, time.Since(startTime), nil)
	return nil
}

// HGet 哈希获取
func HGet(ctx context.Context, key, field string) (string, error) {
	client := GetRedisClient()
	if client == nil {
		return "", errors.New("redis client not initialized")
	}

	startTime := time.Now()

	result, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		pkglog.LogRedisOperation("HGET", key, time.Since(startTime), err)
		if err == redis.Nil {
			return "", fmt.Errorf("field %s not found in hash %s", field, key)
		}
		return "", fmt.Errorf("failed to hget key %s field %s: %w", key, field, err)
	}

	pkglog.LogRedisOperation("HGET", key, time.Since(startTime), nil)
	return result, nil
}

// GetRedisStats 获取Redis统计信息
func GetRedisStats(ctx context.Context) (*RedisStats, error) {
	client := GetRedisClient()
	if client == nil {
		return nil, errors.New("redis client not initialized")
	}

	// 获取连接池统计
	poolStats := client.PoolStats()

	// 测试连接状态
	isConnected := true
	if err := Ping(ctx); err != nil {
		isConnected = false
	}

	stats := &RedisStats{
		TotalConns:  poolStats.TotalConns,
		IdleConns:   poolStats.IdleConns,
		StaleConns:  poolStats.StaleConns,
		Hits:        uint64(poolStats.Hits),
		Misses:      uint64(poolStats.Misses),
		Timeouts:    uint64(poolStats.Timeouts),
		IsConnected: isConnected,
	}

	if redisManager != nil {
		redisManager.mu.RLock()
		stats.LastPing = redisManager.lastPing
		stats.Uptime = time.Since(redisManager.lastPing)
		redisManager.mu.RUnlock()
	}

	return stats, nil
}

// validateRedisConfig 验证Redis配置
func validateRedisConfig(cfg *config.RedisConfig) error {
	if cfg.Host == "" {
		return errors.New("redis host is required")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("redis port must be between 1 and 65535")
	}

	if cfg.DB < 0 || cfg.DB > 15 {
		return errors.New("redis db must be between 0 and 15")
	}

	if cfg.PoolSize <= 0 {
		return errors.New("redis pool size must be positive")
	}

	if cfg.MinIdleConns < 0 {
		return errors.New("redis min idle conns cannot be negative")
	}

	if cfg.MinIdleConns > cfg.PoolSize {
		return errors.New("redis min idle conns cannot exceed pool size")
	}

	if cfg.DialTimeout <= 0 {
		return errors.New("redis dial timeout must be positive")
	}

	if cfg.ReadTimeout <= 0 {
		return errors.New("redis read timeout must be positive")
	}

	if cfg.WriteTimeout <= 0 {
		return errors.New("redis write timeout must be positive")
	}

	if cfg.IdleTimeout <= 0 {
		return errors.New("redis idle timeout must be positive")
	}

	return nil
}
