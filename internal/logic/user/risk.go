package user

import (
	"context"
	"math"
	"time"

	"OneGfServer/internal/model/user"
)

// ===============================
// 用户风险评分相关业务逻辑
// ===============================

// CalculateLoginRiskScore 计算登录风险评分
func (s *sUser) CalculateLoginRiskScore(ctx context.Context, input *user.CalculateLoginRiskScoreInput) (*user.CalculateLoginRiskScoreOutput, error) {
	score := 0.0

	// 检查失败登录次数
	if failedAttempts, ok := input.UserData["failed_login_attempts"].(int); ok {
		if failedAttempts > 3 {
			score += float64(failedAttempts) * 10
		}
	}

	// 检查登录频率
	if loginCount, ok := input.UserData["recent_login_count"].(int); ok {
		if loginCount > 20 { // 最近短时间内登录次数过多
			score += 30
		}
	}

	return &user.CalculateLoginRiskScoreOutput{
		Score: math.Min(score, 100),
	}, nil
}

// CalculateLocationRiskScore 计算地理位置风险评分
func (s *sUser) CalculateLocationRiskScore(ctx context.Context, input *user.CalculateLocationRiskScoreInput) (*user.CalculateLocationRiskScoreOutput, error) {
	score := 0.0

	currentIP, _ := input.LoginData["ip_address"].(string)
	lastIP, _ := input.UserData["last_login_ip"].(string)

	if currentIP != "" && lastIP != "" && currentIP != lastIP {
		// 简单的IP变化检测，实际应该使用地理位置API
		isSameNetworkInput := &user.IsSameNetworkInput{
			IP1: currentIP,
			IP2: lastIP,
		}
		isSameNetworkOutput := s.IsSameNetwork(ctx, isSameNetworkInput)
		if !isSameNetworkOutput.IsSame {
			score += 40
		}
	}

	// 检查已知的恶意IP
	isMaliciousInput := &user.IsKnownMaliciousIPInput{
		IP: currentIP,
	}
	isMaliciousOutput := s.IsKnownMaliciousIP(ctx, isMaliciousInput)
	if isMaliciousOutput.IsMalicious {
		score += 60
	}

	return &user.CalculateLocationRiskScoreOutput{
		Score: math.Min(score, 100),
	}, nil
}

// CalculateDeviceRiskScore 计算设备风险评分
func (s *sUser) CalculateDeviceRiskScore(ctx context.Context, input *user.CalculateDeviceRiskScoreInput) (*user.CalculateDeviceRiskScoreOutput, error) {
	score := 0.0

	currentUA, _ := input.LoginData["user_agent"].(string)
	lastUA, _ := input.UserData["last_user_agent"].(string)

	if currentUA != "" && lastUA != "" && currentUA != lastUA {
		score += 20 // 设备变化
	}

	// 检查是否为移动设备
	isMobileInput := &user.IsMobileDeviceInput{
		UserAgent: currentUA,
	}
	isMobileOutput := s.IsMobileDevice(ctx, isMobileInput)
	if isMobileOutput.IsMobile {
		score += 10 // 移动设备风险稍高
	}

	return &user.CalculateDeviceRiskScoreOutput{
		Score: math.Min(score, 100),
	}, nil
}

// CalculateTimeRiskScore 计算时间异常风险评分
func (s *sUser) CalculateTimeRiskScore(ctx context.Context, input *user.CalculateTimeRiskScoreInput) (*user.CalculateTimeRiskScoreOutput, error) {
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

	return &user.CalculateTimeRiskScoreOutput{
		Score: score,
	}, nil
}

// CalculateBehaviorRiskScore 计算行为异常风险评分
func (s *sUser) CalculateBehaviorRiskScore(ctx context.Context, input *user.CalculateBehaviorRiskScoreInput) (*user.CalculateBehaviorRiskScoreOutput, error) {
	score := 0.0

	// 检查登录模式变化
	isAbnormalInput := &user.IsAbnormalLoginPatternInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	isAbnormalOutput := s.IsAbnormalLoginPattern(ctx, isAbnormalInput)
	if isAbnormalOutput.IsAbnormal {
		score += 25
	}

	return &user.CalculateBehaviorRiskScoreOutput{
		Score: score,
	}, nil
}

// CalculateHistoryRiskScore 计算历史安全事件风险评分
func (s *sUser) CalculateHistoryRiskScore(ctx context.Context, input *user.CalculateHistoryRiskScoreInput) (*user.CalculateHistoryRiskScoreOutput, error) {
	score := 0.0

	// 检查最近的安全事件
	if securityEvents, ok := input.UserData["recent_security_events"].(int); ok {
		score += float64(securityEvents) * 15
	}

	return &user.CalculateHistoryRiskScoreOutput{
		Score: math.Min(score, 100),
	}, nil
}
