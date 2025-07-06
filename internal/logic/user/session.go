package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"OneGfServer/internal/model/user"
)

// ===============================
// 用户会话管理相关业务逻辑
// ===============================

// CreateUserSession 创建用户会话
func (s *sUser) CreateUserSession(ctx context.Context, input *user.CreateUserSessionInput) (*user.CreateUserSessionOutput, error) {
	sessionData := make(map[string]interface{})

	// 生成会话ID
	sessionIdInput := &user.GenerateSessionIdInput{}
	sessionIdOutput := s.GenerateSessionId(ctx, sessionIdInput)

	sessionData["session_id"] = sessionIdOutput.SessionId
	sessionData["user_id"] = input.UserId
	sessionData["created_at"] = gtime.Now()
	sessionData["last_activity"] = gtime.Now()

	// 设置会话过期时间（默认2小时）
	sessionData["expires_at"] = gtime.New(time.Now().Add(2 * time.Hour))

	// 记录登录信息
	if ipAddress, ok := input.LoginData["ip_address"].(string); ok {
		sessionData["ip_address"] = ipAddress
	}
	if userAgent, ok := input.LoginData["user_agent"].(string); ok {
		sessionData["user_agent"] = userAgent
	}
	if deviceInfo, ok := input.LoginData["device_info"].(string); ok {
		sessionData["device_info"] = deviceInfo
	}

	// 设置状态
	sessionData["status"] = "active"

	g.Log().Info(ctx, "创建用户会话", g.Map{
		"user_id":    input.UserId,
		"session_id": sessionData["session_id"],
		"ip_address": sessionData["ip_address"],
	})

	return &user.CreateUserSessionOutput{
		SessionData: sessionData,
	}, nil
}

// ValidateUserSession 验证用户会话
func (s *sUser) ValidateUserSession(ctx context.Context, input *user.ValidateUserSessionInput) (*user.ValidateUserSessionOutput, error) {
	// 检查会话是否过期
	if expiresAt, ok := input.SessionData["expires_at"].(*gtime.Time); ok && expiresAt != nil {
		if time.Now().After(expiresAt.Time) {
			return &user.ValidateUserSessionOutput{
				IsValid: false,
				Message: "会话已过期",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "会话已过期")
		}
	}

	// 检查会话状态
	if status, ok := input.SessionData["status"].(string); ok && status != "active" {
		return &user.ValidateUserSessionOutput{
			IsValid: false,
			Message: "会话无效",
		}, gerror.NewCode(gcode.CodeNotAuthorized, "会话无效")
	}

	// 检查会话活跃时间
	if lastActivity, ok := input.SessionData["last_activity"].(*gtime.Time); ok && lastActivity != nil {
		if time.Since(lastActivity.Time) > 30*time.Minute {
			return &user.ValidateUserSessionOutput{
				IsValid: false,
				Message: "会话超时",
			}, gerror.NewCode(gcode.CodeNotAuthorized, "会话超时")
		}
	}

	return &user.ValidateUserSessionOutput{
		IsValid: true,
		Message: "会话验证成功",
	}, nil
}

// RefreshUserSession 刷新用户会话
func (s *sUser) RefreshUserSession(ctx context.Context, input *user.RefreshUserSessionInput) (*user.RefreshUserSessionOutput, error) {
	// 更新最后活跃时间
	input.SessionData["last_activity"] = gtime.Now()

	// 延长过期时间
	input.SessionData["expires_at"] = gtime.New(time.Now().Add(2 * time.Hour))

	return &user.RefreshUserSessionOutput{
		UpdatedSessionData: input.SessionData,
	}, nil
}

// GenerateSessionId 生成会话ID
func (s *sUser) GenerateSessionId(ctx context.Context, input *user.GenerateSessionIdInput) (*user.GenerateSessionIdOutput, error) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	sessionId := hex.EncodeToString(bytes)

	return &user.GenerateSessionIdOutput{
		SessionId: sessionId,
	}, nil
}
