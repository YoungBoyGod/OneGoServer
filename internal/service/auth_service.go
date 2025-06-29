package service

import (
	"errors"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz"
	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// AuthService 认证服务
type AuthService struct {
	userUsecase *biz.UserUsecase
	config      *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(userUsecase *biz.UserUsecase, cfg *config.Config) *AuthService {
	return &AuthService{
		userUsecase: userUsecase,
		config:      cfg,
	}
}

// Claims JWT Claims结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login 用户登录
func (s *AuthService) Login(email, password string) (string, *biz.User, error) {
	// 验证用户密码
	user, err := s.userUsecase.VerifyPassword(email, password)
	if err != nil {
		return "", nil, err
	}

	// 检查用户状态
	if user.Status != biz.UserStatusActive {
		return "", nil, errors.New("账户已被禁用")
	}

	// 生成JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Register 用户注册
func (s *AuthService) Register(username, email, password string) (*biz.User, error) {
	// 创建用户
	user, err := s.userUsecase.CreateUser(username, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateToken 生成JWT token
func (s *AuthService) GenerateToken(user *biz.User) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(s.config.JWT.ExpireHours) * time.Hour)

	// 创建Claims
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Issuer,
			Subject:   user.Username,
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名方法错误")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token有效性
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	// 验证原token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 获取用户信息
	user, err := s.userUsecase.GetUser(claims.UserID)
	if err != nil {
		return "", err
	}

	// 检查用户状态
	if user.Status != biz.UserStatusActive {
		return "", errors.New("账户已被禁用")
	}

	// 生成新token
	return s.GenerateToken(user)
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID uint) (*biz.User, error) {
	user, err := s.userUsecase.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if user.Status != biz.UserStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	return s.userUsecase.ChangePassword(userID, oldPassword, newPassword)
}

// ValidateJWTMiddleware JWT中间件验证函数
func (s *AuthService) ValidateJWTMiddleware(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token不能为空")
	}

	// 移除Bearer前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	return s.ValidateToken(tokenString)
}

// HasPermission 检查用户权限
func (s *AuthService) HasPermission(userRole biz.UserRole, requiredRole biz.UserRole) bool {
	// 定义角色权限层级
	roleLevel := map[biz.UserRole]int{
		biz.UserRoleUser:  1,
		biz.UserRoleAdmin: 2,
		biz.UserRoleRoot:  3,
	}

	userLevel, exists := roleLevel[userRole]
	if !exists {
		return false
	}

	requiredLevel, exists := roleLevel[requiredRole]
	if !exists {
		return false
	}

	return userLevel >= requiredLevel
}

// CheckUserAccess 检查用户访问权限
func (s *AuthService) CheckUserAccess(currentUserID, targetUserID uint, currentUserRole biz.UserRole) bool {
	// 用户可以访问自己的信息
	if currentUserID == targetUserID {
		return true
	}

	// 管理员和超级管理员可以访问其他用户信息
	return s.HasPermission(currentUserRole, biz.UserRoleAdmin)
}

// Logout 用户登出（在使用Redis等存储的黑名单机制时使用）
func (s *AuthService) Logout(tokenString string) error {
	// 这里可以将token加入黑名单
	// 目前JWT是无状态的，所以客户端删除token即可
	// 如果需要服务端控制，可以将token加入Redis黑名单

	// TODO: 实现token黑名单机制
	// redis.Set(ctx, "blacklist:"+tokenString, "1", time.Duration(s.config.JWT.ExpireHours)*time.Hour)

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (s *AuthService) IsTokenBlacklisted(tokenString string) bool {
	// TODO: 检查Redis黑名单
	// exists, _ := redis.Exists(ctx, "blacklist:"+tokenString)
	// return exists

	return false
}
