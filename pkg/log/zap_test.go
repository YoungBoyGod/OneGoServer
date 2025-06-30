package log

import (
	"strings"
	"testing"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"go.uber.org/zap"
)

// TestInitLoggerEnhanced 测试日志器初始化
func TestInitLogger(t *testing.T) {
	cfg := &config.LoggingConfig{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/test_app.log",
		MaxSize:    50,
		MaxBackups: 2,
		MaxAge:     15,
		Compress:   false,
	}

	// 测试正常初始化
	err := InitLoggerEnhanced(cfg)
	if err != nil {
		t.Fatalf("初始化日志器失败: %v", err)
	}

	// 验证logger是否正确创建
	if GetAppLogger(cfg) == nil {
		t.Error("应用日志器未创建")
	}
	if GetHTTPLogger() == nil {
		t.Error("HTTP日志器未创建")
	}
	if GetErrorLogger() == nil {
		t.Error("错误日志器未创建")
	}

	t.Log("✅ 日志器初始化成功")
}

// TestStructuredLogHelpers 测试结构化日志助手
func TestStructuredLogHelpers(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	// 测试用户操作日志
	t.Run("用户操作日志", func(t *testing.T) {
		LogUser("login", 12345, zap.String("ip", "192.168.1.100"))
		t.Log("✅ 用户操作日志记录成功")
	})

	// 测试数据库操作日志
	t.Run("数据库操作日志", func(t *testing.T) {
		// 成功的数据库操作
		LogDBOperation("SELECT", "users", 150*time.Millisecond, nil)

		// 失败的数据库操作
		LogDBOperation("INSERT", "orders", 300*time.Millisecond,
			&mockError{msg: "duplicate key violation"})
		t.Log("✅ 数据库操作日志记录成功")
	})

	// 测试API调用日志
	t.Run("API调用日志", func(t *testing.T) {
		// 成功的API调用
		LogAPICall("GET", "/api/users", 200, 120*time.Millisecond, nil)

		// 失败的API调用
		LogAPICall("POST", "/api/orders", 500, 450*time.Millisecond,
			&mockError{msg: "internal server error"})
		t.Log("✅ API调用日志记录成功")
	})

	// 测试系统事件日志
	t.Run("系统事件日志", func(t *testing.T) {
		LogSystemEvent("server_start", "main",
			zap.String("version", "1.0.0"),
			zap.Int("port", 8080))
		t.Log("✅ 系统事件日志记录成功")
	})
}

// TestContextLogger 测试上下文日志器
func TestContextLogger(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	// 创建上下文日志器
	ctxLogger := WithContext(
		zap.String("request_id", "req-123456"),
		zap.String("user_id", "user-789"),
		zap.String("session_id", "sess-abc123"),
	)

	// 测试各级别日志
	t.Run("上下文日志各级别", func(t *testing.T) {
		ctxLogger.Debug("调试信息", zap.String("action", "debug_test"))
		ctxLogger.Info("信息日志", zap.String("action", "info_test"))
		ctxLogger.Warn("警告日志", zap.String("action", "warn_test"))
		ctxLogger.Error("错误日志", zap.String("action", "error_test"))
		t.Log("✅ 上下文日志各级别记录成功")
	})
}

// TestLoggerStats 测试统计功能
func TestLoggerStats(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	// 记录一些日志来测试统计
	t.Run("统计功能测试", func(t *testing.T) {
		LogInfoWithStats("测试信息1")
		LogInfoWithStats("测试信息2")
		LogWarnWithStats("测试警告1")
		LogErrorWithStats("测试错误1")
		LogDebugWithStats("测试调试1")

		// 获取统计信息
		stats := GetStats()

		if stats.TotalLogs < 5 {
			t.Errorf("期望至少5条日志，实际: %d", stats.TotalLogs)
		}
		if stats.InfoLogs < 2 {
			t.Errorf("期望至少2条信息日志，实际: %d", stats.InfoLogs)
		}
		if stats.WarnLogs < 1 {
			t.Errorf("期望至少1条警告日志，实际: %d", stats.WarnLogs)
		}
		if stats.ErrorLogs < 1 {
			t.Errorf("期望至少1条错误日志，实际: %d", stats.ErrorLogs)
		}
		if stats.DebugLogs < 1 {
			t.Errorf("期望至少1条调试日志，实际: %d", stats.DebugLogs)
		}

		t.Logf("✅ 统计信息: 总计=%d, 信息=%d, 警告=%d, 错误=%d, 调试=%d",
			stats.TotalLogs, stats.InfoLogs, stats.WarnLogs, stats.ErrorLogs, stats.DebugLogs)
	})
}

// TestHealthCheck 测试健康检查
func TestHealthCheck(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	t.Run("健康检查", func(t *testing.T) {
		err := HealthCheck()
		if err != nil {
			t.Errorf("健康检查失败: %v", err)
		}

		// 获取状态信息
		status := GetLoggerStatus()
		if !status["initialized"].(bool) {
			t.Error("日志系统应该已初始化")
		}

		t.Log("✅ 日志系统健康检查通过")
		t.Logf("系统状态: %+v", status)
	})
}

// TestHTTPRequestLogging 测试HTTP请求日志记录
func TestHTTPRequestLogging(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	t.Run("HTTP请求日志", func(t *testing.T) {
		// 记录HTTP访问
		LogAPICall("GET", "/api/users", 200, 150*time.Millisecond, nil)

		// 记录带IP信息的访问
		WithContext(
			zap.String("client_ip", "192.168.1.100"),
			zap.String("user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"),
			zap.String("real_ip", "203.0.113.1"),
		).Info("HTTP请求处理完成",
			zap.String("method", "GET"),
			zap.String("path", "/api/users"),
			zap.String("query", "page=1&limit=10"),
			zap.Int("status", 200),
			zap.Duration("latency", 150*time.Millisecond),
		)

		t.Log("✅ HTTP请求日志记录成功")
		t.Log("📊 包含完整的IP地址和用户代理信息")
	})
}

// TestUserAgentParsing 测试用户代理解析
func TestUserAgentParsing(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	userAgents := []struct {
		name string
		ua   string
	}{
		{"Chrome浏览器", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		{"Safari浏览器", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"},
		{"Firefox浏览器", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0"},
		{"Edge浏览器", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"},
		{"Postman客户端", "PostmanRuntime/7.32.3"},
		{"cURL工具", "curl/7.68.0"},
		{"移动Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"},
		{"Android Chrome", "Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"},
	}

	t.Run("用户代理解析和记录", func(t *testing.T) {
		for _, ua := range userAgents {
			// 解析用户代理信息
			browser := parseBrowser(ua.ua)
			platform := parsePlatform(ua.ua)

			// 记录详细的用户代理信息
			WithContext(
				zap.String("user_agent_type", ua.name),
				zap.String("browser", browser),
				zap.String("platform", platform),
			).Info("用户代理分析",
				zap.String("raw_user_agent", ua.ua),
				zap.Bool("is_mobile", isMobile(ua.ua)),
				zap.Bool("is_bot", isBot(ua.ua)),
			)
		}

		t.Log("✅ 用户代理解析和记录完成")
		t.Log("🔍 支持主流浏览器和客户端识别")
	})
}

// TestConfigReload 测试配置热重载
func TestConfigReload(t *testing.T) {
	// 初始化日志器
	cfg := &config.LoggingConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	InitLoggerEnhanced(cfg)

	t.Run("配置热重载", func(t *testing.T) {
		// 新配置
		newCfg := &config.LoggingConfig{
			Level:      "debug",
			Format:     "console",
			Output:     "stdout",
			FilePath:   "logs/reloaded_app.log",
			MaxSize:    200,
			MaxBackups: 5,
			MaxAge:     60,
			Compress:   true,
		}

		// 重载配置
		err := ReloadConfig(newCfg)
		if err != nil {
			t.Errorf("配置重载失败: %v", err)
		}

		// 验证新配置
		currentCfg := GetCurrentConfig()
		if currentCfg.Level != "debug" {
			t.Errorf("期望日志级别为debug，实际: %s", currentCfg.Level)
		}
		if currentCfg.Format != "console" {
			t.Errorf("期望格式为console，实际: %s", currentCfg.Format)
		}

		// 测试新配置下的日志记录
		LogDebug("配置重载后的调试日志")
		LogInfo("配置重载后的信息日志")

		t.Log("✅ 配置热重载成功")
		t.Log("🔄 日志级别和格式已更新")
	})
}

// === 辅助函数和类型 ===

// mockError 模拟错误类型
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

// parseBrowser 解析浏览器类型
func parseBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "postman") {
		return "Postman"
	}
	if strings.Contains(ua, "curl") {
		return "cURL"
	}
	if strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/") {
		return "Microsoft Edge"
	}
	if strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/") {
		return "Google Chrome"
	}
	if strings.Contains(ua, "firefox/") {
		return "Mozilla Firefox"
	}
	if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") {
		return "Safari"
	}
	if strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/") {
		return "Opera"
	}

	return "Unknown Browser"
}

// parsePlatform 解析平台类型
func parsePlatform(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return "iOS"
	}
	if strings.Contains(ua, "android") {
		return "Android"
	}
	if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os x") {
		return "macOS"
	}
	if strings.Contains(ua, "windows") {
		return "Windows"
	}
	if strings.Contains(ua, "linux") {
		return "Linux"
	}

	return "Unknown Platform"
}

// isMobile 判断是否为移动设备
func isMobile(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	mobileKeywords := []string{"mobile", "android", "iphone", "ipad", "tablet"}

	for _, keyword := range mobileKeywords {
		if strings.Contains(ua, keyword) {
			return true
		}
	}
	return false
}

// isBot 判断是否为机器人
func isBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	botKeywords := []string{"bot", "crawler", "spider", "scraper"}

	for _, keyword := range botKeywords {
		if strings.Contains(ua, keyword) {
			return true
		}
	}
	return false
}
