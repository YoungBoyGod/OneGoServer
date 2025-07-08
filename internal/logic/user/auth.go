package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"OneGfServer/internal/model/user"
)

// ===============================
// 用户认证授权相关业务逻辑
// ===============================

// ValidateUserLogin 验证用户登录
func (s *sUser) ValidateUserLogin(ctx context.Context, input *user.ValidateUserLoginInput) (*user.ValidateUserLoginOutput, error) {
	// 检查用户状态
	if status, ok := input.UserData["status"].(string); ok {
		switch status {
		case "inactive":
			return &user.ValidateUserLoginOutput{
				IsValid: false,
				Message: "用户账户已被禁用",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "用户账户已被禁用")
		case "locked":
			return &user.ValidateUserLoginOutput{
				IsValid: false,
				Message: "用户账户已被锁定",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "用户账户已被锁定")
		case "deleted":
			return &user.ValidateUserLoginOutput{
				IsValid: false,
				Message: "用户账户不存在",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "用户账户不存在")
		}
	}

	// 检查账户锁定状态
	if lockedUntil, ok := input.UserData["account_locked_until"].(*gtime.Time); ok && lockedUntil != nil {
		if time.Now().Before(lockedUntil.Time) {
			remainingTime := lockedUntil.Time.Sub(time.Now())
			return &user.ValidateUserLoginOutput{
					IsValid: false,
					Message: fmt.Sprintf("账户已锁定，剩余时间: %d分钟", int(remainingTime.Minutes())),
				}, gerror.NewCode(gcode.CodeNotAuthorized,
					fmt.Sprintf("账户已锁定，剩余时间: %d分钟", int(remainingTime.Minutes())))
		}
	}

	// 验证密码
	if hashedPassword, ok := input.UserData["password"].(string); ok {
		verifyInput := &user.VerifyPasswordInput{
			Password:       input.Password,
			HashedPassword: hashedPassword,
		}
		// TODO: Fix assignment mismatch - VerifyPassword returns 2 values
		// verifyOutput := s.VerifyPassword(ctx, verifyInput)
		verifyOutput, _ := s.VerifyPassword(ctx, verifyInput)
		if verifyOutput == nil || !verifyOutput.IsValid {
			// 增加失败登录次数
			handleInput := &user.HandleFailedLoginInput{
				UserData: input.UserData,
			}
			s.HandleFailedLogin(ctx, handleInput)
			return &user.ValidateUserLoginOutput{
				IsValid: false,
				Message: "用户名或密码错误",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "用户名或密码错误")
		}
	} else {
		return &user.ValidateUserLoginOutput{
			IsValid: false,
			Message: "用户密码数据异常",
		}, gerror.NewCode(gcode.CodeInternalError, "用户密码数据异常")
	}

	// 检查密码是否需要修改
	if forceChange, ok := input.UserData["force_password_change"].(bool); ok && forceChange {
		return &user.ValidateUserLoginOutput{
			IsValid: false,
			Message: "需要修改密码后才能登录",
		}, gerror.NewCode(gcode.CodeValidationFailed, "需要修改密码后才能登录")
	}

	return &user.ValidateUserLoginOutput{
		IsValid: true,
		Message: "登录验证成功",
	}, nil
}

// HashPassword 密码加密
func (s *sUser) HashPassword(ctx context.Context, input *user.HashPasswordInput) (*user.HashPasswordOutput, error) {
	// 生成随机盐值
	saltInput := &user.GenerateSaltInput{}
	// TODO: Fix assignment mismatch - GenerateSalt returns 2 values
	saltOutput, _ := s.GenerateSalt(ctx, saltInput)

	// 使用盐值和密码生成hash
	hash := sha256.Sum256([]byte(input.Password + saltOutput.Salt))
	hashedPassword := saltOutput.Salt + ":" + hex.EncodeToString(hash[:])

	return &user.HashPasswordOutput{
		HashedPassword: hashedPassword,
	}, nil
}

// VerifyPassword 验证密码
func (s *sUser) VerifyPassword(ctx context.Context, input *user.VerifyPasswordInput) (*user.VerifyPasswordOutput, error) {
	parts := strings.Split(input.HashedPassword, ":")
	if len(parts) != 2 {
		return &user.VerifyPasswordOutput{
			IsValid: false,
		}, nil
	}

	salt := parts[0]
	storedHash := parts[1]

	// 使用相同的盐值计算hash
	hash := sha256.Sum256([]byte(input.Password + salt))
	calculatedHash := hex.EncodeToString(hash[:])

	return &user.VerifyPasswordOutput{
		IsValid: calculatedHash == storedHash,
	}, nil
}

// GenerateSalt 生成随机盐值
func (s *sUser) GenerateSalt(ctx context.Context, input *user.GenerateSaltInput) (*user.GenerateSaltOutput, error) {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	salt := hex.EncodeToString(bytes)

	return &user.GenerateSaltOutput{
		Salt: salt,
	}, nil
}

// HandleFailedLogin 处理登录失败
func (s *sUser) HandleFailedLogin(ctx context.Context, input *user.HandleFailedLoginInput) (*user.HandleFailedLoginOutput, error) {
	failedAttempts := 0
	if attempts, ok := input.UserData["failed_login_attempts"].(int); ok {
		failedAttempts = attempts
	}

	failedAttempts++
	input.UserData["failed_login_attempts"] = failedAttempts
	input.UserData["last_failed_login"] = gtime.Now()

	isLocked := false
	var lockDuration string

	// 如果失败次数过多，锁定账户
	if failedAttempts >= 5 {
		lockDuration := time.Duration(math.Pow(2, float64(failedAttempts-5))) * time.Hour
		if lockDuration > 24*time.Hour {
			lockDuration = 24 * time.Hour // 最多锁定24小时
		}
		input.UserData["account_locked_until"] = gtime.New(time.Now().Add(lockDuration))
		isLocked = true

		g.Log().Warning(ctx, "用户账户因多次登录失败被锁定", g.Map{
			"failed_attempts": failedAttempts,
			"lock_duration":   lockDuration.String(),
		})
	}

	return &user.HandleFailedLoginOutput{
		UpdatedUserData: input.UserData,
		IsLocked:        isLocked,
		LockDuration:    lockDuration,
	}, nil
}

// CalculateUserRiskScore 计算用户风险评分
func (s *sUser) CalculateUserRiskScore(ctx context.Context, input *user.CalculateUserRiskScoreInput) (*user.CalculateUserRiskScoreOutput, error) {
	riskScore := 0.0
	riskFactors := make(map[string]float64)

	// 1. 登录频率异常 (权重: 25%)
	loginRiskInput := &user.CalculateLoginRiskScoreInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	// TODO: Fix assignment mismatch - CalculateLoginRiskScore returns 2 values
	loginRiskOutput, _ := s.CalculateLoginRiskScore(ctx, loginRiskInput)
	loginRisk := loginRiskOutput.Score
	riskScore += loginRisk * 0.25
	riskFactors["login_risk"] = loginRisk

	// 2. 地理位置异常 (权重: 20%)
	locationRiskInput := &user.CalculateLocationRiskScoreInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	// TODO: Fix assignment mismatch - CalculateLocationRiskScore returns 2 values
	locationRiskOutput, _ := s.CalculateLocationRiskScore(ctx, locationRiskInput)
	locationRisk := locationRiskOutput.Score
	riskScore += locationRisk * 0.20
	riskFactors["location_risk"] = locationRisk

	// 3. 设备异常 (权重: 20%)
	deviceRiskInput := &user.CalculateDeviceRiskScoreInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	// TODO: Fix assignment mismatch - CalculateDeviceRiskScore returns 2 values
	deviceRiskOutput, _ := s.CalculateDeviceRiskScore(ctx, deviceRiskInput)
	deviceRisk := deviceRiskOutput.Score
	riskScore += deviceRisk * 0.20
	riskFactors["device_risk"] = deviceRisk

	// 4. 时间异常 (权重: 15%)
	timeRiskInput := &user.CalculateTimeRiskScoreInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	// TODO: Fix assignment mismatch - CalculateTimeRiskScore returns 2 values
	timeRiskOutput, _ := s.CalculateTimeRiskScore(ctx, timeRiskInput)
	timeRisk := timeRiskOutput.Score
	riskScore += timeRisk * 0.15
	riskFactors["time_risk"] = timeRisk

	// 5. 行为异常 (权重: 10%)
	behaviorRiskInput := &user.CalculateBehaviorRiskScoreInput{
		UserData:  input.UserData,
		LoginData: input.LoginData,
	}
	// TODO: Fix assignment mismatch - CalculateBehaviorRiskScore returns 2 values
	behaviorRiskOutput, _ := s.CalculateBehaviorRiskScore(ctx, behaviorRiskInput)
	behaviorRisk := behaviorRiskOutput.Score
	riskScore += behaviorRisk * 0.10
	riskFactors["behavior_risk"] = behaviorRisk

	// 6. 历史安全事件 (权重: 10%)
	historyRiskInput := &user.CalculateHistoryRiskScoreInput{
		UserData: input.UserData,
	}
	// TODO: Fix assignment mismatch - CalculateHistoryRiskScore returns 2 values
	historyRiskOutput, _ := s.CalculateHistoryRiskScore(ctx, historyRiskInput)
	historyRisk := historyRiskOutput.Score
	riskScore += historyRisk * 0.10
	riskFactors["history_risk"] = historyRisk

	// 确定风险等级
	riskLevel := "low"
	if riskScore > 80 {
		riskLevel = "critical"
	} else if riskScore > 60 {
		riskLevel = "high"
	} else if riskScore > 40 {
		riskLevel = "medium"
	}

	// 生成建议
	var recommendations []string
	if riskScore > 60 {
		recommendations = append(recommendations, "建议启用双因素认证")
	}
	if locationRisk > 50 {
		recommendations = append(recommendations, "检测到异常登录地点，建议验证身份")
	}
	if deviceRisk > 30 {
		recommendations = append(recommendations, "检测到新设备登录，建议确认设备安全")
	}

	return &user.CalculateUserRiskScoreOutput{
		RiskScore:       math.Min(riskScore, 100),
		RiskFactors:     riskFactors,
		RiskLevel:       riskLevel,
		Recommendations: recommendations,
	}, nil
}
