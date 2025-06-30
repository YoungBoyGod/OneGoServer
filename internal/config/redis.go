package config

import "fmt"

// 获取redis的连接字符串
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
