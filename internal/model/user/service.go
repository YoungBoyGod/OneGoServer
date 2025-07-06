package user

import (
	"context"
)

// IUserService 用户服务接口
type IUserService interface {
	// 认证授权相关
	ValidateUserLogin(ctx context.Context, input *ValidateUserLoginInput) (*ValidateUserLoginOutput, error)
	HashPassword(ctx context.Context, input *HashPasswordInput) (*HashPasswordOutput, error)
	VerifyPassword(ctx context.Context, input *VerifyPasswordInput) (*VerifyPasswordOutput, error)
	GenerateSalt(ctx context.Context, input *GenerateSaltInput) (*GenerateSaltOutput, error)
	HandleFailedLogin(ctx context.Context, input *HandleFailedLoginInput) (*HandleFailedLoginOutput, error)
	CalculateUserRiskScore(ctx context.Context, input *CalculateUserRiskScoreInput) (*CalculateUserRiskScoreOutput, error)

	// 用户验证相关
	ValidateUserRegistration(ctx context.Context, input *ValidateUserRegistrationInput) (*ValidateUserRegistrationOutput, error)
	ValidateUsername(ctx context.Context, input *ValidateUsernameInput) (*ValidateUsernameOutput, error)
	ValidatePasswordStrength(ctx context.Context, input *ValidatePasswordStrengthInput) (*ValidatePasswordStrengthOutput, error)
	ValidateEmail(ctx context.Context, input *ValidateEmailInput) (*ValidateEmailOutput, error)
	ValidatePhone(ctx context.Context, input *ValidatePhoneInput) (*ValidatePhoneOutput, error)

	// 会话管理相关
	CreateUserSession(ctx context.Context, input *CreateUserSessionInput) (*CreateUserSessionOutput, error)
	ValidateUserSession(ctx context.Context, input *ValidateUserSessionInput) (*ValidateUserSessionOutput, error)
	RefreshUserSession(ctx context.Context, input *RefreshUserSessionInput) (*RefreshUserSessionOutput, error)
	GenerateSessionId(ctx context.Context, input *GenerateSessionIdInput) (*GenerateSessionIdOutput, error)

	// 权限管理相关
	CalculateUserPermissions(ctx context.Context, input *CalculateUserPermissionsInput) (*CalculateUserPermissionsOutput, error)
	CheckUserPermission(ctx context.Context, input *CheckUserPermissionInput) (*CheckUserPermissionOutput, error)

	// 行为分析相关
	AnalyzeUserBehavior(ctx context.Context, input *AnalyzeUserBehaviorInput) (*AnalyzeUserBehaviorOutput, error)
	FindMostActiveHour(ctx context.Context, input *FindMostActiveHourInput) (*FindMostActiveHourOutput, error)
	IdentifyBehaviorPattern(ctx context.Context, input *IdentifyBehaviorPatternInput) (*IdentifyBehaviorPatternOutput, error)
	CalculateActivityScore(ctx context.Context, input *CalculateActivityScoreInput) (*CalculateActivityScoreOutput, error)

	// 风险评分相关
	CalculateLoginRiskScore(ctx context.Context, input *CalculateLoginRiskScoreInput) (*CalculateLoginRiskScoreOutput, error)
	CalculateLocationRiskScore(ctx context.Context, input *CalculateLocationRiskScoreInput) (*CalculateLocationRiskScoreOutput, error)
	CalculateDeviceRiskScore(ctx context.Context, input *CalculateDeviceRiskScoreInput) (*CalculateDeviceRiskScoreOutput, error)
	CalculateTimeRiskScore(ctx context.Context, input *CalculateTimeRiskScoreInput) (*CalculateTimeRiskScoreOutput, error)
	CalculateBehaviorRiskScore(ctx context.Context, input *CalculateBehaviorRiskScoreInput) (*CalculateBehaviorRiskScoreOutput, error)
	CalculateHistoryRiskScore(ctx context.Context, input *CalculateHistoryRiskScoreInput) (*CalculateHistoryRiskScoreOutput, error)

	// 辅助功能相关
	IsSameNetwork(ctx context.Context, input *IsSameNetworkInput) (*IsSameNetworkOutput, error)
	IsKnownMaliciousIP(ctx context.Context, input *IsKnownMaliciousIPInput) (*IsKnownMaliciousIPOutput, error)
	IsMobileDevice(ctx context.Context, input *IsMobileDeviceInput) (*IsMobileDeviceOutput, error)
	IsAbnormalLoginPattern(ctx context.Context, input *IsAbnormalLoginPatternInput) (*IsAbnormalLoginPatternOutput, error)
}
