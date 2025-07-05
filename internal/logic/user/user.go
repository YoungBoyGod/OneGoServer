package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sUser struct{}

func New() *sUser {
	return &sUser{}
}

func init() {
	// service.RegisterUser(New())
}

// ===============================
// 用户认证授权相关业务逻辑
// ===============================

// ValidateUserLogin 验证用户登录
func (s *sUser) ValidateUserLogin(ctx context.Context, username, password string, userData map[string]interface{}) error {
	// 检查用户状态
	if status, ok := userData["status"].(string); ok {
		switch status {
		case "inactive":
			return gerror.NewCode(gcode.CodeNotAuthorized, "用户账户已被禁用")
		case "locked":
			return gerror.NewCode(gcode.CodeNotAuthorized, "用户账户已被锁定")
		case "deleted":
			return gerror.NewCode(gcode.CodeNotAuthorized, "用户账户不存在")
		}
	}

	// 检查账户锁定状态
	if lockedUntil, ok := userData["account_locked_until"].(*gtime.Time); ok && lockedUntil != nil {
		if time.Now().Before(lockedUntil.Time) {
			remainingTime := lockedUntil.Time.Sub(time.Now())
			return gerror.NewCode(gcode.CodeNotAuthorized,
				fmt.Sprintf("账户已锁定，剩余时间: %d分钟", int(remainingTime.Minutes())))
		}
	}

	// 验证密码
	if hashedPassword, ok := userData["password"].(string); ok {
		if !s.VerifyPassword(password, hashedPassword) {
			// 增加失败登录次数
			s.handleFailedLogin(ctx, userData)
			return gerror.NewCode(gcode.CodeNotAuthorized, "用户名或密码错误")
		}
	} else {
		return gerror.NewCode(gcode.CodeInternalError, "用户密码数据异常")
	}

	// 检查密码是否需要修改
	if forceChange, ok := userData["force_password_change"].(bool); ok && forceChange {
		return gerror.NewCode(gcode.CodeValidationFailed, "需要修改密码后才能登录")
	}

	return nil
}

// HashPassword 密码加密
func (s *sUser) HashPassword(password string) string {
	// 生成随机盐值
	salt := s.generateSalt()
	// 使用盐值和密码生成hash
	hash := sha256.Sum256([]byte(password + salt))
	return salt + ":" + hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码
func (s *sUser) VerifyPassword(password, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	storedHash := parts[1]

	// 使用相同的盐值计算hash
	hash := sha256.Sum256([]byte(password + salt))
	calculatedHash := hex.EncodeToString(hash[:])

	return calculatedHash == storedHash
}

// generateSalt 生成随机盐值
func (s *sUser) generateSalt() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// handleFailedLogin 处理登录失败
func (s *sUser) handleFailedLogin(ctx context.Context, userData map[string]interface{}) {
	failedAttempts := 0
	if attempts, ok := userData["failed_login_attempts"].(int); ok {
		failedAttempts = attempts
	}

	failedAttempts++
	userData["failed_login_attempts"] = failedAttempts
	userData["last_failed_login"] = gtime.Now()

	// 如果失败次数过多，锁定账户
	if failedAttempts >= 5 {
		lockDuration := time.Duration(math.Pow(2, float64(failedAttempts-5))) * time.Hour
		if lockDuration > 24*time.Hour {
			lockDuration = 24 * time.Hour // 最多锁定24小时
		}
		userData["account_locked_until"] = gtime.New(time.Now().Add(lockDuration))

		g.Log().Warning(ctx, "用户账户因多次登录失败被锁定", g.Map{
			"failed_attempts": failedAttempts,
			"lock_duration":   lockDuration.String(),
		})
	}
}

// CalculateUserRiskScore 计算用户风险评分
func (s *sUser) CalculateUserRiskScore(ctx context.Context, userData map[string]interface{}, loginData map[string]interface{}) float64 {
	riskScore := 0.0

	// 1. 登录频率异常 (权重: 25%)
	loginRisk := s.calculateLoginRiskScore(userData, loginData)
	riskScore += loginRisk * 0.25

	// 2. 地理位置异常 (权重: 20%)
	locationRisk := s.calculateLocationRiskScore(userData, loginData)
	riskScore += locationRisk * 0.20

	// 3. 设备异常 (权重: 20%)
	deviceRisk := s.calculateDeviceRiskScore(userData, loginData)
	riskScore += deviceRisk * 0.20

	// 4. 时间异常 (权重: 15%)
	timeRisk := s.calculateTimeRiskScore(userData, loginData)
	riskScore += timeRisk * 0.15

	// 5. 行为异常 (权重: 10%)
	behaviorRisk := s.calculateBehaviorRiskScore(userData, loginData)
	riskScore += behaviorRisk * 0.10

	// 6. 历史安全事件 (权重: 10%)
	historyRisk := s.calculateHistoryRiskScore(userData)
	riskScore += historyRisk * 0.10

	return math.Min(riskScore, 100)
}

// calculateLoginRiskScore 计算登录风险评分
func (s *sUser) calculateLoginRiskScore(userData, loginData map[string]interface{}) float64 {
	score := 0.0

	// 检查失败登录次数
	if failedAttempts, ok := userData["failed_login_attempts"].(int); ok {
		if failedAttempts > 3 {
			score += float64(failedAttempts) * 10
		}
	}

	// 检查登录频率
	if loginCount, ok := userData["recent_login_count"].(int); ok {
		if loginCount > 20 { // 最近短时间内登录次数过多
			score += 30
		}
	}

	return math.Min(score, 100)
}

// calculateLocationRiskScore 计算地理位置风险评分
func (s *sUser) calculateLocationRiskScore(userData, loginData map[string]interface{}) float64 {
	score := 0.0

	currentIP, _ := loginData["ip_address"].(string)
	lastIP, _ := userData["last_login_ip"].(string)

	if currentIP != "" && lastIP != "" && currentIP != lastIP {
		// 简单的IP变化检测，实际应该使用地理位置API
		if !s.isSameNetwork(currentIP, lastIP) {
			score += 40
		}
	}

	// 检查已知的恶意IP
	if s.isKnownMaliciousIP(currentIP) {
		score += 60
	}

	return math.Min(score, 100)
}

// calculateDeviceRiskScore 计算设备风险评分
func (s *sUser) calculateDeviceRiskScore(userData, loginData map[string]interface{}) float64 {
	score := 0.0

	currentUA, _ := loginData["user_agent"].(string)
	lastUA, _ := userData["last_user_agent"].(string)

	if currentUA != "" && lastUA != "" && currentUA != lastUA {
		score += 20 // 设备变化
	}

	// 检查是否为移动设备
	if s.isMobileDevice(currentUA) {
		score += 10 // 移动设备风险稍高
	}

	return math.Min(score, 100)
}

// calculateTimeRiskScore 计算时间异常风险评分
func (s *sUser) calculateTimeRiskScore(userData, loginData map[string]interface{}) float64 {
	score := 0.0

	now := time.Now()
	hour := now.Hour()

	// 检查是否在异常时间登录（深夜或凌晨）
	if hour < 6 || hour > 23 {
		score += 15
	}

	// 检查是否为周末
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		score += 5
	}

	return score
}

// calculateBehaviorRiskScore 计算行为异常风险评分
func (s *sUser) calculateBehaviorRiskScore(userData, loginData map[string]interface{}) float64 {
	score := 0.0

	// 检查登录模式变化
	if s.isAbnormalLoginPattern(userData, loginData) {
		score += 25
	}

	return score
}

// calculateHistoryRiskScore 计算历史安全事件风险评分
func (s *sUser) calculateHistoryRiskScore(userData map[string]interface{}) float64 {
	score := 0.0

	// 检查最近的安全事件
	if securityEvents, ok := userData["recent_security_events"].(int); ok {
		score += float64(securityEvents) * 15
	}

	return math.Min(score, 100)
}

// ===============================
// 用户验证相关业务逻辑
// ===============================

// ValidateUserRegistration 验证用户注册数据
func (s *sUser) ValidateUserRegistration(ctx context.Context, userData map[string]interface{}) error {
	// 验证必需字段
	requiredFields := []string{"username", "password", "email"}
	for _, field := range requiredFields {
		if value, ok := userData[field]; !ok || value == nil || value == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("缺少必需字段: %s", field))
		}
	}

	// 验证用户名格式
	if username, ok := userData["username"].(string); ok {
		if err := s.ValidateUsername(username); err != nil {
			return err
		}
	}

	// 验证密码强度
	if password, ok := userData["password"].(string); ok {
		if err := s.ValidatePasswordStrength(password); err != nil {
			return err
		}
	}

	// 验证邮箱格式
	if email, ok := userData["email"].(string); ok {
		if err := s.ValidateEmail(email); err != nil {
			return err
		}
	}

	// 验证手机号格式（如果提供）
	if phone, ok := userData["phone"].(string); ok && phone != "" {
		if err := s.ValidatePhone(phone); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUsername 验证用户名格式
func (s *sUser) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "用户名长度必须在3-20字符之间")
	}

	// 用户名只能包含字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return gerror.NewCode(gcode.CodeInvalidParameter, "用户名只能包含字母、数字和下划线")
	}

	// 不能以数字开头
	if username[0] >= '0' && username[0] <= '9' {
		return gerror.NewCode(gcode.CodeInvalidParameter, "用户名不能以数字开头")
	}

	return nil
}

// ValidatePasswordStrength 验证密码强度
func (s *sUser) ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "密码长度不能少于6个字符")
	}

	if len(password) > 32 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "密码长度不能超过32个字符")
	}

	// 检查密码复杂度
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\?]`).MatchString(password)

	complexityCount := 0
	if hasUpper {
		complexityCount++
	}
	if hasLower {
		complexityCount++
	}
	if hasNumber {
		complexityCount++
	}
	if hasSpecial {
		complexityCount++
	}

	if complexityCount < 3 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "密码必须包含大写字母、小写字母、数字、特殊字符中的至少3种")
	}

	// 检查常见弱密码
	weakPasswords := []string{"123456", "password", "123456789", "qwerty", "abc123"}
	for _, weak := range weakPasswords {
		if strings.Contains(strings.ToLower(password), weak) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "密码过于简单，请使用更复杂的密码")
		}
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func (s *sUser) ValidateEmail(email string) error {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return gerror.NewCode(gcode.CodeInvalidParameter, "邮箱格式不正确")
	}
	return nil
}

// ValidatePhone 验证手机号格式
func (s *sUser) ValidatePhone(phone string) error {
	// 简单的中国手机号验证
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, phone)
	if !matched {
		return gerror.NewCode(gcode.CodeInvalidParameter, "手机号格式不正确")
	}
	return nil
}

// ===============================
// 用户会话管理相关业务逻辑
// ===============================

// CreateUserSession 创建用户会话
func (s *sUser) CreateUserSession(ctx context.Context, userId string, loginData map[string]interface{}) (map[string]interface{}, error) {
	sessionData := make(map[string]interface{})

	// 生成会话ID
	sessionData["session_id"] = s.generateSessionId()
	sessionData["user_id"] = userId
	sessionData["created_at"] = gtime.Now()
	sessionData["last_activity"] = gtime.Now()

	// 设置会话过期时间（默认2小时）
	sessionData["expires_at"] = gtime.New(time.Now().Add(2 * time.Hour))

	// 记录登录信息
	if ipAddress, ok := loginData["ip_address"].(string); ok {
		sessionData["ip_address"] = ipAddress
	}
	if userAgent, ok := loginData["user_agent"].(string); ok {
		sessionData["user_agent"] = userAgent
	}
	if deviceInfo, ok := loginData["device_info"].(string); ok {
		sessionData["device_info"] = deviceInfo
	}

	// 设置状态
	sessionData["status"] = "active"

	g.Log().Info(ctx, "创建用户会话", g.Map{
		"user_id":    userId,
		"session_id": sessionData["session_id"],
		"ip_address": sessionData["ip_address"],
	})

	return sessionData, nil
}

// ValidateUserSession 验证用户会话
func (s *sUser) ValidateUserSession(ctx context.Context, sessionData map[string]interface{}) error {
	// 检查会话是否过期
	if expiresAt, ok := sessionData["expires_at"].(*gtime.Time); ok && expiresAt != nil {
		if time.Now().After(expiresAt.Time) {
			return gerror.NewCode(gcode.CodeNotAuthorized, "会话已过期")
		}
	}

	// 检查会话状态
	if status, ok := sessionData["status"].(string); ok && status != "active" {
		return gerror.NewCode(gcode.CodeNotAuthorized, "会话无效")
	}

	// 检查会话活跃时间
	if lastActivity, ok := sessionData["last_activity"].(*gtime.Time); ok && lastActivity != nil {
		if time.Since(lastActivity.Time) > 30*time.Minute {
			return gerror.NewCode(gcode.CodeNotAuthorized, "会话超时")
		}
	}

	return nil
}

// RefreshUserSession 刷新用户会话
func (s *sUser) RefreshUserSession(ctx context.Context, sessionData map[string]interface{}) map[string]interface{} {
	// 更新最后活跃时间
	sessionData["last_activity"] = gtime.Now()

	// 延长过期时间
	sessionData["expires_at"] = gtime.New(time.Now().Add(2 * time.Hour))

	return sessionData
}

// generateSessionId 生成会话ID
func (s *sUser) generateSessionId() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ===============================
// 用户权限管理相关业务逻辑
// ===============================

// CalculateUserPermissions 计算用户权限
func (s *sUser) CalculateUserPermissions(ctx context.Context, userRoles []map[string]interface{}) []string {
	permissionSet := make(map[string]bool)

	for _, role := range userRoles {
		if permissions, ok := role["permissions"].([]string); ok {
			for _, permission := range permissions {
				permissionSet[permission] = true
			}
		}
	}

	// 转换为数组
	var permissions []string
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return permissions
}

// CheckUserPermission 检查用户权限
func (s *sUser) CheckUserPermission(ctx context.Context, userPermissions []string, requiredPermission string, resource string) bool {
	// 检查直接权限
	for _, permission := range userPermissions {
		if permission == requiredPermission {
			return true
		}

		// 检查通配符权限
		if strings.HasSuffix(permission, "*") {
			prefix := strings.TrimSuffix(permission, "*")
			if strings.HasPrefix(requiredPermission, prefix) {
				return true
			}
		}
	}

	// 检查资源特定权限
	resourcePermission := fmt.Sprintf("%s:%s", resource, requiredPermission)
	for _, permission := range userPermissions {
		if permission == resourcePermission {
			return true
		}
	}

	return false
}

// ===============================
// 用户行为分析相关业务逻辑
// ===============================

// AnalyzeUserBehavior 分析用户行为模式
func (s *sUser) AnalyzeUserBehavior(ctx context.Context, activityData []map[string]interface{}) map[string]interface{} {
	if len(activityData) == 0 {
		return map[string]interface{}{
			"status":  "insufficient_data",
			"message": "活动数据不足",
		}
	}

	analysis := make(map[string]interface{})

	// 统计活动类型分布
	actionCounts := make(map[string]int)
	hourCounts := make(map[int]int)
	var totalSessions, totalDuration float64

	for _, activity := range activityData {
		// 统计操作类型
		if action, ok := activity["action_type"].(string); ok {
			actionCounts[action]++
		}

		// 统计活跃时间段
		if timestamp, ok := activity["timestamp"].(*gtime.Time); ok && timestamp != nil {
			hour := timestamp.Time.Hour()
			hourCounts[hour]++
		}

		// 统计会话时长
		if duration, ok := activity["session_duration"].(float64); ok {
			totalDuration += duration
			totalSessions++
		}
	}

	analysis["action_distribution"] = actionCounts
	analysis["active_hours"] = hourCounts

	if totalSessions > 0 {
		analysis["avg_session_duration"] = totalDuration / totalSessions
	}

	// 识别活跃时间段
	analysis["most_active_hour"] = s.findMostActiveHour(hourCounts)

	// 识别行为模式
	analysis["behavior_pattern"] = s.identifyBehaviorPattern(actionCounts, hourCounts)

	// 计算活跃度评分
	analysis["activity_score"] = s.calculateActivityScore(activityData)

	return analysis
}

// findMostActiveHour 找到最活跃的时间段
func (s *sUser) findMostActiveHour(hourCounts map[int]int) int {
	maxCount := 0
	mostActiveHour := 0

	for hour, count := range hourCounts {
		if count > maxCount {
			maxCount = count
			mostActiveHour = hour
		}
	}

	return mostActiveHour
}

// identifyBehaviorPattern 识别行为模式
func (s *sUser) identifyBehaviorPattern(actionCounts map[string]int, hourCounts map[int]int) string {
	// 根据活动类型判断用户类型
	readCount := actionCounts["view"] + actionCounts["read"]
	writeCount := actionCounts["create"] + actionCounts["update"] + actionCounts["delete"]

	if writeCount > readCount*2 {
		return "power_user"
	} else if readCount > writeCount*5 {
		return "viewer"
	} else {
		return "normal_user"
	}
}

// calculateActivityScore 计算活跃度评分
func (s *sUser) calculateActivityScore(activityData []map[string]interface{}) float64 {
	if len(activityData) == 0 {
		return 0
	}

	score := 0.0

	// 基于活动频率
	activityPerDay := float64(len(activityData)) / 30.0 // 假设30天的数据
	if activityPerDay > 10 {
		score += 50
	} else {
		score += activityPerDay * 5
	}

	// 基于活动多样性
	actionTypes := make(map[string]bool)
	for _, activity := range activityData {
		if action, ok := activity["action_type"].(string); ok {
			actionTypes[action] = true
		}
	}
	score += float64(len(actionTypes)) * 5

	return math.Min(score, 100)
}

// ===============================
// 辅助函数
// ===============================

// isSameNetwork 检查两个IP是否在同一网络
func (s *sUser) isSameNetwork(ip1, ip2 string) bool {
	// 简单实现：检查前三段是否相同
	parts1 := strings.Split(ip1, ".")
	parts2 := strings.Split(ip2, ".")

	if len(parts1) != 4 || len(parts2) != 4 {
		return false
	}

	return parts1[0] == parts2[0] && parts1[1] == parts2[1] && parts1[2] == parts2[2]
}

// isKnownMaliciousIP 检查是否为已知恶意IP
func (s *sUser) isKnownMaliciousIP(ip string) bool {
	// 实际实现中应该查询恶意IP数据库
	// 这里只是示例
	maliciousIPs := []string{"192.168.1.100", "10.0.0.100"}
	for _, maliciousIP := range maliciousIPs {
		if ip == maliciousIP {
			return true
		}
	}
	return false
}

// isMobileDevice 检查是否为移动设备
func (s *sUser) isMobileDevice(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	mobileIndicators := []string{"mobile", "android", "iphone", "ipad", "tablet"}

	for _, indicator := range mobileIndicators {
		if strings.Contains(userAgent, indicator) {
			return true
		}
	}
	return false
}

// isAbnormalLoginPattern 检查是否为异常登录模式
func (s *sUser) isAbnormalLoginPattern(userData, loginData map[string]interface{}) bool {
	// 实际实现中需要更复杂的机器学习算法
	// 这里只是简单示例

	// 检查登录频率是否异常
	if recentLogins, ok := userData["recent_login_count"].(int); ok {
		if recentLogins > 10 { // 短时间内登录次数过多
			return true
		}
	}

	return false
}
