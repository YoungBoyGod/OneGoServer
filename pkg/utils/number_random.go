package utils

import (
	"time"
)

// generateRandomString 生成随机字符串 - 改进版本使用更好的随机源
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	// 使用当前纳秒时间加索引作为种子，增加随机性
	baseTime := time.Now().UnixNano()
	for i := range b {
		// 结合索引和时间变化增加随机性
		seed := (baseTime + int64(i*1000)) % int64(len(charset))
		b[i] = charset[seed]
	}
	return string(b)
}
