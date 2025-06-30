package log

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"go.uber.org/zap"
)

// setupTestLogger 设置测试用的日志器
func setupTestLogger(t *testing.T) *config.LoggingConfig {
	cfg := &config.LoggingConfig{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/test.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	err := InitLoggerEnhanced(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize enhanced logger: %v", err)
	}

	return cfg
}

// TestInitLoggerEnhanced 测试增强版日志器初始化
func TestInitLoggerEnhanced(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := &config.LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "logs/test.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}

		err := InitLoggerEnhanced(cfg)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		cfg := &config.LoggingConfig{
			Level:    "invalid",
			Format:   "json",
			Output:   "stdout",
			FilePath: "logs/test.log",
			MaxSize:  -1, // 无效值
		}

		err := InitLoggerEnhanced(cfg)
		if err == nil {
			t.Error("Expected error for invalid config, got nil")
		}
	})
}

// TestStructuredLogHelpers 测试结构化日志助手
func TestStructuredLogHelpers(t *testing.T) {
	cfg := setupTestLogger(t)

	t.Run("LogUser", func(t *testing.T) {
		// 测试用户登录日志
		LogUser("login", 123,
			zap.String("ip", "192.168.1.100"),
			zap.String("user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"),
			zap.String("device", "Chrome Browser"),
			zap.String("location", "Shanghai, China"),
		)

		// 测试用户注销日志
		LogUser("logout", 123,
			zap.String("ip", "192.168.1.100"),
			zap.String("user_agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)"),
			zap.String("device", "Mobile Safari"),
			zap.String("session_duration", "2h30m"),
		)

		// 验证统计
		stats := GetStats()
		if stats.TotalLogs == 0 {
			t.Error("Expected logs to be recorded")
		}
	})

	t.Run("LogDBOperation", func(t *testing.T) {
		initialStats := GetStats()

		// 测试成功的数据库操作
		LogDBOperation("SELECT", "users", 150*time.Millisecond, nil)

		// 测试失败的数据库操作
		LogDBOperation("UPDATE", "users", 200*time.Millisecond, fmt.Errorf("connection timeout"))

		newStats := GetStats()
		if newStats.TotalLogs <= initialStats.TotalLogs {
			t.Error("Expected log count to increase")
		}
	})

	t.Run("LogAPICall", func(t *testing.T) {
		initialStats := GetStats()

		// 测试成功的API调用
		LogAPICall("GET", "/api/users", 200, 250*time.Millisecond, nil)

		// 测试失败的API调用
		LogAPICall("POST", "/api/users", 500, 100*time.Millisecond, fmt.Errorf("internal server error"))

		newStats := GetStats()
		if newStats.TotalLogs <= initialStats.TotalLogs {
			t.Error("Expected log count to increase")
		}
	})

	t.Run("LogSystemEvent", func(t *testing.T) {
		LogSystemEvent("server_start", "success",
			zap.Int("port", 8080),
			zap.String("mode", "production"),
			zap.String("version", "v1.0.0"),
		)

		// 验证统计
		stats := GetStats()
		if stats.InfoLogs == 0 {
			t.Error("Expected info logs to be recorded")
		}
	})

	// 清理
	Sync()
	_ = cfg
}

// TestContextLogger 测试上下文日志器
func TestContextLogger(t *testing.T) {
	setupTestLogger(t)

	t.Run("WithContext", func(t *testing.T) {
		contextLogger := WithContext(
			zap.String("request_id", "req-12345"),
			zap.String("user_id", "user123"),
		)

		if contextLogger == nil {
			t.Error("Expected context logger to be created")
		}

		// 测试不同级别的日志
		contextLogger.Info("Processing user request",
			zap.String("action", "get_profile"),
			zap.String("resource", "/api/profile"),
		)

		contextLogger.Warn("Potential issue detected")
		contextLogger.Error("Request failed")
		contextLogger.Debug("Debug information")
	})
}

// TestLoggerStats 测试日志统计功能
func TestLoggerStats(t *testing.T) {
	setupTestLogger(t)

	t.Run("StatsTracking", func(t *testing.T) {
		initialStats := GetStats()

		// 使用增强版日志函数记录一些日志
		LogInfoWithStats("Test info message")
		LogWarnWithStats("Test warn message")
		LogErrorWithStats("Test error message")
		LogDebugWithStats("Test debug message")

		newStats := GetStats()

		// 验证统计增加
		if newStats.TotalLogs <= initialStats.TotalLogs {
			t.Errorf("Expected total logs to increase. Initial: %d, New: %d",
				initialStats.TotalLogs, newStats.TotalLogs)
		}
		if newStats.LastLogTime.Before(initialStats.LastLogTime) {
			t.Error("Expected last log time to be updated")
		}

		// 验证各级别日志计数
		if newStats.InfoLogs <= initialStats.InfoLogs {
			t.Error("Expected info logs to increase")
		}
		if newStats.WarnLogs <= initialStats.WarnLogs {
			t.Error("Expected warn logs to increase")
		}
		if newStats.ErrorLogs <= initialStats.ErrorLogs {
			t.Error("Expected error logs to increase")
		}
		if newStats.DebugLogs <= initialStats.DebugLogs {
			t.Error("Expected debug logs to increase")
		}
	})
}

// TestHealthCheck 测试健康检查功能
func TestHealthCheck(t *testing.T) {
	setupTestLogger(t)

	t.Run("HealthCheck", func(t *testing.T) {
		err := HealthCheck()
		if err != nil {
			t.Errorf("Health check failed: %v", err)
		}
	})

	t.Run("LoggerStatus", func(t *testing.T) {
		status := GetLoggerStatus()

		// 验证状态信息
		if !status["initialized"].(bool) {
			t.Error("Expected logger to be initialized")
		}

		if status["config"] == nil {
			t.Error("Expected config to be present in status")
		}

		if status["stats"] == nil {
			t.Error("Expected stats to be present in status")
		}
	})
}

// TestHTTPRequestLogging 测试HTTP请求日志记录
func TestHTTPRequestLogging(t *testing.T) {
	setupTestLogger(t)

	testCases := []struct {
		name      string
		method    string
		path      string
		ip        string
		userAgent string
		userID    int64
		status    int
		duration  time.Duration
	}{
		{
			name:      "ChromeLogin",
			method:    "POST",
			path:      "/api/auth/login",
			ip:        "192.168.1.100",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			userID:    123,
			status:    200,
			duration:  150 * time.Millisecond,
		},
		{
			name:      "MobileProfile",
			method:    "GET",
			path:      "/api/users/profile",
			ip:        "10.0.0.45",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1",
			userID:    456,
			status:    200,
			duration:  85 * time.Millisecond,
		},
		{
			name:      "PostmanAPI",
			method:    "POST",
			path:      "/api/orders",
			ip:        "203.208.60.1",
			userAgent: "PostmanRuntime/7.28.4",
			userID:    789,
			status:    201,
			duration:  320 * time.Millisecond,
		},
		{
			name:      "CurlError",
			method:    "GET",
			path:      "/api/admin/users",
			ip:        "172.16.0.1",
			userAgent: "curl/7.68.0",
			userID:    1,
			status:    403,
			duration:  25 * time.Millisecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 使用上下文日志器记录完整的请求信息
			contextLogger := WithContext(
				zap.String("request_id", fmt.Sprintf("test-req-%s", tc.name)),
				zap.String("method", tc.method),
				zap.String("path", tc.path),
				zap.String("client_ip", tc.ip),
				zap.String("user_agent", tc.userAgent),
				zap.Int64("user_id", tc.userID),
			)

			// 请求开始
			contextLogger.Info("HTTP request started")

			// 记录用户操作（如果是认证相关的请求）
			if tc.path == "/api/auth/login" {
				LogUser("login_attempt", tc.userID,
					zap.String("client_ip", tc.ip),
					zap.String("user_agent", tc.userAgent),
					zap.String("browser", extractBrowser(tc.userAgent)),
					zap.String("platform", extractPlatform(tc.userAgent)),
				)
			}

			// API调用记录
			var err error
			if tc.status >= 400 {
				err = fmt.Errorf("request failed with status %d", tc.status)
			}
			LogAPICall(tc.method, tc.path, tc.status, tc.duration, err)

			// 请求完成
			if tc.status >= 400 {
				contextLogger.Error("HTTP request failed",
					zap.Int("status_code", tc.status),
					zap.Duration("duration", tc.duration),
				)
			} else {
				contextLogger.Info("HTTP request completed",
					zap.Int("status_code", tc.status),
					zap.Duration("duration", tc.duration),
				)
			}

			// 验证浏览器识别
			browser := extractBrowser(tc.userAgent)
			if browser == "" {
				t.Error("Expected browser to be extracted from user agent")
			}

			// 验证平台识别
			platform := extractPlatform(tc.userAgent)
			if platform == "" {
				t.Error("Expected platform to be extracted from user agent")
			}

			t.Logf("IP: %s, Browser: %s, Platform: %s, Status: %d",
				tc.ip, browser, platform, tc.status)
		})
	}
}

// TestUserAgentParsing 测试用户代理解析功能
func TestUserAgentParsing(t *testing.T) {
	testCases := []struct {
		userAgent      string
		expectBrowser  string
		expectPlatform string
	}{
		{
			userAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			expectBrowser:  "Chrome",
			expectPlatform: "macOS",
		},
		{
			userAgent:      "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Mobile/15E148 Safari/604.1",
			expectBrowser:  "Safari",
			expectPlatform: "iOS",
		},
		{
			userAgent:      "PostmanRuntime/7.28.4",
			expectBrowser:  "Postman",
			expectPlatform: "Unknown",
		},
		{
			userAgent:      "curl/7.68.0",
			expectBrowser:  "cURL",
			expectPlatform: "Unknown",
		},
		{
			userAgent:      "",
			expectBrowser:  "Unknown",
			expectPlatform: "Unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expectBrowser, func(t *testing.T) {
			browser := extractBrowser(tc.userAgent)
			platform := extractPlatform(tc.userAgent)

			if browser != tc.expectBrowser {
				t.Errorf("Expected browser %s, got %s", tc.expectBrowser, browser)
			}

			if platform != tc.expectPlatform {
				t.Errorf("Expected platform %s, got %s", tc.expectPlatform, platform)
			}
		})
	}
}

// TestConfigReload 测试配置热重载
func TestConfigReload(t *testing.T) {
	setupTestLogger(t)

	t.Run("ReloadConfig", func(t *testing.T) {
		newCfg := &config.LoggingConfig{
			Level:      "warn",
			Format:     "console",
			Output:     "stdout",
			FilePath:   "logs/test-reload.log",
			MaxSize:    50,
			MaxBackups: 5,
			MaxAge:     15,
			Compress:   false,
		}

		err := ReloadConfig(newCfg)
		if err != nil {
			t.Errorf("Failed to reload config: %v", err)
		}

		// 验证配置已更新
		currentCfg := GetCurrentConfig()
		if currentCfg.Level != "warn" {
			t.Errorf("Expected level to be 'warn', got '%s'", currentCfg.Level)
		}
	})
}

// 从User-Agent中提取浏览器信息
func extractBrowser(userAgent string) string {
	ua := userAgent
	if len(ua) == 0 {
		return "Unknown"
	}

	browsers := map[string]string{
		"Chrome":         "Chrome",
		"Safari":         "Safari",
		"Firefox":        "Firefox",
		"Edge":           "Edge",
		"PostmanRuntime": "Postman",
		"curl":           "cURL",
	}

	for key, browser := range browsers {
		if strings.Contains(ua, key) {
			return browser
		}
	}

	return "Unknown"
}

// 从User-Agent中提取平台信息
func extractPlatform(userAgent string) string {
	ua := userAgent
	if len(ua) == 0 {
		return "Unknown"
	}

	platforms := map[string]string{
		"Macintosh": "macOS",
		"Windows":   "Windows",
		"iPhone":    "iOS",
		"Android":   "Android",
		"Linux":     "Linux",
	}

	for key, platform := range platforms {
		if strings.Contains(ua, key) {
			return platform
		}
	}

	return "Unknown"
}
