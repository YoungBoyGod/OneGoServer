package user

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/internal/model/user"
)

// ===============================
// 用户验证相关业务逻辑
// ===============================

// ValidateUserRegistration 验证用户注册数据
func (s *sUser) ValidateUserRegistration(ctx context.Context, input *user.ValidateUserRegistrationInput) (*user.ValidateUserRegistrationOutput, error) {
	errors := make(map[string]string)

	// 验证必需字段
	requiredFields := []string{"username", "password", "email"}
	for _, field := range requiredFields {
		if value, ok := input.UserData[field]; !ok || value == nil || value == "" {
			errors[field] = fmt.Sprintf("缺少必需字段: %s", field)
		}
	}

	// 验证用户名格式
	if username, ok := input.UserData["username"].(string); ok {
		usernameInput := &user.ValidateUsernameInput{
			Username: username,
		}
		// TODO: Fix assignment mismatch - ValidateUsername returns 2 values
		usernameOutput, _ := s.ValidateUsername(ctx, usernameInput)
		if usernameOutput != nil && !usernameOutput.IsValid {
			errors["username"] = usernameOutput.Message
		}
	}

	// 验证密码强度
	if password, ok := input.UserData["password"].(string); ok {
		passwordInput := &user.ValidatePasswordStrengthInput{
			Password: password,
		}
		// TODO: Fix assignment mismatch - ValidatePasswordStrength returns 2 values
		passwordOutput, _ := s.ValidatePasswordStrength(ctx, passwordInput)
		if passwordOutput != nil && !passwordOutput.IsValid {
			errors["password"] = strings.Join(passwordOutput.Suggestions, "; ")
		}
	}

	// 验证邮箱格式
	if email, ok := input.UserData["email"].(string); ok {
		emailInput := &user.ValidateEmailInput{
			Email: email,
		}
		// TODO: Fix assignment mismatch - ValidateEmail returns 2 values
		emailOutput, _ := s.ValidateEmail(ctx, emailInput)
		if emailOutput != nil && !emailOutput.IsValid {
			errors["email"] = emailOutput.Message
		}
	}

	// 验证手机号格式（如果提供）
	if phone, ok := input.UserData["phone"].(string); ok && phone != "" {
		phoneInput := &user.ValidatePhoneInput{
			Phone: phone,
		}
		// TODO: Fix assignment mismatch - ValidatePhone returns 2 values
		phoneOutput, _ := s.ValidatePhone(ctx, phoneInput)
		if phoneOutput != nil && !phoneOutput.IsValid {
			errors["phone"] = phoneOutput.Message
		}
	}

	isValid := len(errors) == 0
	return &user.ValidateUserRegistrationOutput{
		IsValid: isValid,
		Errors:  errors,
	}, nil
}

// ValidateUsername 验证用户名格式
func (s *sUser) ValidateUsername(ctx context.Context, input *user.ValidateUsernameInput) (*user.ValidateUsernameOutput, error) {
	if len(input.Username) < 3 || len(input.Username) > 20 {
		return &user.ValidateUsernameOutput{
			IsValid: false,
			Message: "用户名长度必须在3-20字符之间",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "用户名长度必须在3-20字符之间")
	}

	// 用户名只能包含字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", input.Username)
	if !matched {
		return &user.ValidateUsernameOutput{
			IsValid: false,
			Message: "用户名只能包含字母、数字和下划线",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "用户名只能包含字母、数字和下划线")
	}

	// 不能以数字开头
	if input.Username[0] >= '0' && input.Username[0] <= '9' {
		return &user.ValidateUsernameOutput{
			IsValid: false,
			Message: "用户名不能以数字开头",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "用户名不能以数字开头")
	}

	return &user.ValidateUsernameOutput{
		IsValid: true,
		Message: "用户名格式正确",
	}, nil
}

// ValidatePasswordStrength 验证密码强度
func (s *sUser) ValidatePasswordStrength(ctx context.Context, input *user.ValidatePasswordStrengthInput) (*user.ValidatePasswordStrengthOutput, error) {
	var suggestions []string
	score := 0.0

	if len(input.Password) < 6 {
		suggestions = append(suggestions, "密码长度不能少于6个字符")
		score -= 20
	}

	if len(input.Password) > 32 {
		suggestions = append(suggestions, "密码长度不能超过32个字符")
		score -= 10
	}

	// 检查密码复杂度
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(input.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(input.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(input.Password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\?]`).MatchString(input.Password)

	complexityCount := 0
	if hasUpper {
		complexityCount++
		score += 10
	}
	if hasLower {
		complexityCount++
		score += 10
	}
	if hasNumber {
		complexityCount++
		score += 10
	}
	if hasSpecial {
		complexityCount++
		score += 15
	}

	if complexityCount < 3 {
		suggestions = append(suggestions, "密码必须包含大写字母、小写字母、数字、特殊字符中的至少3种")
		score -= 30
	}

	// 检查常见弱密码
	weakPasswords := []string{"123456", "password", "123456789", "qwerty", "abc123"}
	for _, weak := range weakPasswords {
		if strings.Contains(strings.ToLower(input.Password), weak) {
			suggestions = append(suggestions, "密码过于简单，请使用更复杂的密码")
			score -= 40
			break
		}
	}

	// 长度奖励
	if len(input.Password) >= 12 {
		score += 20
	} else if len(input.Password) >= 8 {
		score += 10
	}

	// 确定强度等级
	strength := "weak"
	if score >= 80 {
		strength = "strong"
	} else if score >= 60 {
		strength = "medium"
	} else if score >= 40 {
		strength = "fair"
	}

	isValid := len(suggestions) == 0 && score >= 60

	return &user.ValidatePasswordStrengthOutput{
		IsValid:     isValid,
		Strength:    strength,
		Score:       math.Max(0, math.Min(score, 100)),
		Suggestions: suggestions,
	}, nil
}

// ValidateEmail 验证邮箱格式
func (s *sUser) ValidateEmail(ctx context.Context, input *user.ValidateEmailInput) (*user.ValidateEmailOutput, error) {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, input.Email)
	if !matched {
		return &user.ValidateEmailOutput{
			IsValid: false,
			Message: "邮箱格式不正确",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "邮箱格式不正确")
	}

	return &user.ValidateEmailOutput{
		IsValid: true,
		Message: "邮箱格式正确",
	}, nil
}

// ValidatePhone 验证手机号格式
func (s *sUser) ValidatePhone(ctx context.Context, input *user.ValidatePhoneInput) (*user.ValidatePhoneOutput, error) {
	// 简单的中国手机号验证
	phoneRegex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(phoneRegex, input.Phone)
	if !matched {
		return &user.ValidatePhoneOutput{
			IsValid: false,
			Message: "手机号格式不正确",
		}, gerror.NewCode(gcode.CodeInvalidParameter, "手机号格式不正确")
	}

	return &user.ValidatePhoneOutput{
		IsValid: true,
		Message: "手机号格式正确",
	}, nil
}
