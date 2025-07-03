package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateExecutionID 生成唯一的执行ID
func GenerateExecutionID() string {
	// 使用时间戳和安全随机数生成唯一ID
	return fmt.Sprintf("exec_%s_%s",
		time.Now().Format("20060102_150405"),
		GenerateSecureRandomString(8))
}

// GenerateSecureRandomString 生成加密安全的随机字符串
func GenerateSecureRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)

	// 使用crypto/rand生成真正安全的随机数
	if _, err := rand.Read(b); err != nil {
		// 如果加密随机数生成失败，回退到时间戳方式
		return generateTimeBasedRandomString(length)
	}

	// 将随机字节映射到字符集
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}

	return string(b)
}

// generateTimeBasedRandomString 基于时间的随机字符串生成（回退方案）
func generateTimeBasedRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	baseTime := time.Now().UnixNano()

	for i := range b {
		seed := (baseTime + int64(i*1000)) % int64(len(charset))
		b[i] = charset[seed]
	}

	return string(b)
}

// GenerateTaskID 生成任务ID
func GenerateTaskID() string {
	return fmt.Sprintf("task_%s_%s",
		time.Now().Format("20060102"),
		GenerateSecureRandomString(10))
}

// GenerateDeviceID 生成设备ID
func GenerateDeviceID() string {
	return fmt.Sprintf("dev_%s_%s",
		time.Now().Format("20060102"),
		GenerateSecureRandomString(12))
}
