// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"OneGfServer/internal/model/user"
	"context"
)

type (
	IUser interface {
		// ValidateUserLogin 验证用户登录
		ValidateUserLogin(ctx context.Context, input *user.ValidateUserLoginInput) (*user.ValidateUserLoginOutput, error)
		// HashPassword 密码加密
		HashPassword(ctx context.Context, input *user.HashPasswordInput) (*user.HashPasswordOutput, error)
		// VerifyPassword 验证密码
		VerifyPassword(ctx context.Context, input *user.VerifyPasswordInput) (*user.VerifyPasswordOutput, error)
		// GenerateSalt 生成随机盐值
		GenerateSalt(ctx context.Context, input *user.GenerateSaltInput) (*user.GenerateSaltOutput, error)
		// HandleFailedLogin 处理登录失败
		HandleFailedLogin(ctx context.Context, input *user.HandleFailedLoginInput) (*user.HandleFailedLoginOutput, error)
		// CalculateUserRiskScore 计算用户风险评分
		CalculateUserRiskScore(ctx context.Context, input *user.CalculateUserRiskScoreInput) (*user.CalculateUserRiskScoreOutput, error)
		// AnalyzeUserBehavior 分析用户行为模式
		AnalyzeUserBehavior(ctx context.Context, input *user.AnalyzeUserBehaviorInput) (*user.AnalyzeUserBehaviorOutput, error)
		// FindMostActiveHour 找到最活跃的时间段
		FindMostActiveHour(ctx context.Context, input *user.FindMostActiveHourInput) (*user.FindMostActiveHourOutput, error)
		// IdentifyBehaviorPattern 识别行为模式
		IdentifyBehaviorPattern(ctx context.Context, input *user.IdentifyBehaviorPatternInput) (*user.IdentifyBehaviorPatternOutput, error)
		// CalculateActivityScore 计算活跃度评分
		CalculateActivityScore(ctx context.Context, input *user.CalculateActivityScoreInput) (*user.CalculateActivityScoreOutput, error)
		// CalculateUserPermissions 计算用户权限
		CalculateUserPermissions(ctx context.Context, input *user.CalculateUserPermissionsInput) (*user.CalculateUserPermissionsOutput, error)
		// CheckUserPermission 检查用户权限
		CheckUserPermission(ctx context.Context, input *user.CheckUserPermissionInput) (*user.CheckUserPermissionOutput, error)
		// CalculateLoginRiskScore 计算登录风险评分
		CalculateLoginRiskScore(ctx context.Context, input *user.CalculateLoginRiskScoreInput) (*user.CalculateLoginRiskScoreOutput, error)
		// CalculateLocationRiskScore 计算地理位置风险评分
		CalculateLocationRiskScore(ctx context.Context, input *user.CalculateLocationRiskScoreInput) (*user.CalculateLocationRiskScoreOutput, error)
		// CalculateDeviceRiskScore 计算设备风险评分
		CalculateDeviceRiskScore(ctx context.Context, input *user.CalculateDeviceRiskScoreInput) (*user.CalculateDeviceRiskScoreOutput, error)
		// CalculateTimeRiskScore 计算时间异常风险评分
		CalculateTimeRiskScore(ctx context.Context, input *user.CalculateTimeRiskScoreInput) (*user.CalculateTimeRiskScoreOutput, error)
		// CalculateBehaviorRiskScore 计算行为异常风险评分
		CalculateBehaviorRiskScore(ctx context.Context, input *user.CalculateBehaviorRiskScoreInput) (*user.CalculateBehaviorRiskScoreOutput, error)
		// CalculateHistoryRiskScore 计算历史安全事件风险评分
		CalculateHistoryRiskScore(ctx context.Context, input *user.CalculateHistoryRiskScoreInput) (*user.CalculateHistoryRiskScoreOutput, error)
		// CreateUserSession 创建用户会话
		CreateUserSession(ctx context.Context, input *user.CreateUserSessionInput) (*user.CreateUserSessionOutput, error)
		// ValidateUserSession 验证用户会话
		ValidateUserSession(ctx context.Context, input *user.ValidateUserSessionInput) (*user.ValidateUserSessionOutput, error)
		// RefreshUserSession 刷新用户会话
		RefreshUserSession(ctx context.Context, input *user.RefreshUserSessionInput) (*user.RefreshUserSessionOutput, error)
		// GenerateSessionId 生成会话ID
		GenerateSessionId(ctx context.Context, input *user.GenerateSessionIdInput) (*user.GenerateSessionIdOutput, error)
		// IsSameNetwork 检查两个IP是否在同一网络
		IsSameNetwork(ctx context.Context, input *user.IsSameNetworkInput) (*user.IsSameNetworkOutput, error)
		// IsKnownMaliciousIP 检查是否为已知恶意IP
		IsKnownMaliciousIP(ctx context.Context, input *user.IsKnownMaliciousIPInput) (*user.IsKnownMaliciousIPOutput, error)
		// IsMobileDevice 检查是否为移动设备
		IsMobileDevice(ctx context.Context, input *user.IsMobileDeviceInput) (*user.IsMobileDeviceOutput, error)
		// IsAbnormalLoginPattern 检查是否为异常登录模式
		IsAbnormalLoginPattern(ctx context.Context, input *user.IsAbnormalLoginPatternInput) (*user.IsAbnormalLoginPatternOutput, error)
		// ValidateUserRegistration 验证用户注册数据
		ValidateUserRegistration(ctx context.Context, input *user.ValidateUserRegistrationInput) (*user.ValidateUserRegistrationOutput, error)
		// ValidateUsername 验证用户名格式
		ValidateUsername(ctx context.Context, input *user.ValidateUsernameInput) (*user.ValidateUsernameOutput, error)
		// ValidatePasswordStrength 验证密码强度
		ValidatePasswordStrength(ctx context.Context, input *user.ValidatePasswordStrengthInput) (*user.ValidatePasswordStrengthOutput, error)
		// ValidateEmail 验证邮箱格式
		ValidateEmail(ctx context.Context, input *user.ValidateEmailInput) (*user.ValidateEmailOutput, error)
		// ValidatePhone 验证手机号格式
		ValidatePhone(ctx context.Context, input *user.ValidatePhoneInput) (*user.ValidatePhoneOutput, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
