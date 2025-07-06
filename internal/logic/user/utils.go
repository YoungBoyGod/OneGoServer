package user

import (
	"context"
	"strings"

	"OneGfServer/internal/model/user"
)

// ===============================
// 辅助功能相关业务逻辑
// ===============================

// IsSameNetwork 检查两个IP是否在同一网络
func (s *sUser) IsSameNetwork(ctx context.Context, input *user.IsSameNetworkInput) (*user.IsSameNetworkOutput, error) {
	// 简单实现：检查前三段是否相同
	parts1 := strings.Split(input.IP1, ".")
	parts2 := strings.Split(input.IP2, ".")

	if len(parts1) != 4 || len(parts2) != 4 {
		return &user.IsSameNetworkOutput{
			IsSame: false,
		}, nil
	}

	isSame := parts1[0] == parts2[0] && parts1[1] == parts2[1] && parts1[2] == parts2[2]
	return &user.IsSameNetworkOutput{
		IsSame: isSame,
	}, nil
}

// IsKnownMaliciousIP 检查是否为已知恶意IP
func (s *sUser) IsKnownMaliciousIP(ctx context.Context, input *user.IsKnownMaliciousIPInput) (*user.IsKnownMaliciousIPOutput, error) {
	// 实际实现中应该查询恶意IP数据库
	// 这里只是示例
	maliciousIPs := []string{"192.168.1.100", "10.0.0.100"}
	for _, maliciousIP := range maliciousIPs {
		if input.IP == maliciousIP {
			return &user.IsKnownMaliciousIPOutput{
				IsMalicious: true,
			}, nil
		}
	}

	return &user.IsKnownMaliciousIPOutput{
		IsMalicious: false,
	}, nil
}

// IsMobileDevice 检查是否为移动设备
func (s *sUser) IsMobileDevice(ctx context.Context, input *user.IsMobileDeviceInput) (*user.IsMobileDeviceOutput, error) {
	userAgent := strings.ToLower(input.UserAgent)
	mobileIndicators := []string{"mobile", "android", "iphone", "ipad", "tablet"}

	for _, indicator := range mobileIndicators {
		if strings.Contains(userAgent, indicator) {
			return &user.IsMobileDeviceOutput{
				IsMobile: true,
			}, nil
		}
	}

	return &user.IsMobileDeviceOutput{
		IsMobile: false,
	}, nil
}

// IsAbnormalLoginPattern 检查是否为异常登录模式
func (s *sUser) IsAbnormalLoginPattern(ctx context.Context, input *user.IsAbnormalLoginPatternInput) (*user.IsAbnormalLoginPatternOutput, error) {
	// 实际实现中需要更复杂的机器学习算法
	// 这里只是简单示例

	// 检查登录频率是否异常
	if recentLogins, ok := input.UserData["recent_login_count"].(int); ok {
		if recentLogins > 10 { // 短时间内登录次数过多
			return &user.IsAbnormalLoginPatternOutput{
				IsAbnormal: true,
			}, nil
		}
	}

	return &user.IsAbnormalLoginPatternOutput{
		IsAbnormal: false,
	}, nil
}
