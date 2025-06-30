package cache

// 读取配置文件中和env中redis的配置
// 初始化redis
// 返回redis的client

import (
	"github.com/redis/go-redis/v9"
)

// Client Redis的客户端实例
var RedisClient *redis.Client

func InitRedis(cfg *config.RedisConfig) *redis.Client {

}
