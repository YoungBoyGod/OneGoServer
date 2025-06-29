package biz

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers(offset, limit int, filter *UserFilter) ([]*User, int64, error)
}

// UserFilter 用户查询过滤器
type UserFilter struct {
	Username string
	Email    string
	Status   *UserStatus
	Role     *UserRole
}

// UserUsecase 用户用例
type UserUsecase struct {
	repo UserRepo
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(username, email, password string) (*User, error) {
	// 验证输入
	if err := uc.validateUserInput(username, email, password); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	if _, err := uc.repo.GetUserByUsername(username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := uc.repo.GetUserByEmail(email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := uc.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// 创建用户实体
	user := &User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Status:    UserStatusActive,
		Role:      UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库
	if err := uc.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser 获取用户信息
func (uc *UserUsecase) GetUser(id uint) (*User, error) {
	return uc.repo.GetUserByID(id)
}

// GetUserByEmail 通过邮箱获取用户
func (uc *UserUsecase) GetUserByEmail(email string) (*User, error) {
	return uc.repo.GetUserByEmail(email)
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(id uint, updates map[string]interface{}) (*User, error) {
	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if username, ok := updates["username"].(string); ok {
		user.Username = username
	}
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if avatar, ok := updates["avatar"].(string); ok {
		user.Avatar = avatar
	}
	if status, ok := updates["status"].(UserStatus); ok {
		user.Status = status
	}

	user.UpdatedAt = time.Now()

	if err := uc.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func (uc *UserUsecase) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !uc.checkPassword(oldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	// 验证新密码
	if err := uc.validatePassword(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := uc.hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	return uc.repo.UpdateUser(user)
}

// DeleteUser 删除用户
func (uc *UserUsecase) DeleteUser(id uint) error {
	return uc.repo.DeleteUser(id)
}

// ListUsers 用户列表
func (uc *UserUsecase) ListUsers(page, pageSize int, filter *UserFilter) ([]*User, int64, error) {
	offset := (page - 1) * pageSize
	return uc.repo.ListUsers(offset, pageSize, filter)
}

// VerifyPassword 验证密码
func (uc *UserUsecase) VerifyPassword(email, password string) (*User, error) {
	user, err := uc.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !uc.checkPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	uc.repo.UpdateUser(user)

	return user, nil
}

// 私有方法

// validateUserInput 验证用户输入
func (uc *UserUsecase) validateUserInput(username, email, password string) error {
	if len(username) < 3 {
		return errors.New("用户名长度不能少于3个字符")
	}

	if len(email) == 0 {
		return errors.New("邮箱不能为空")
	}

	return uc.validatePassword(password)
}

// validatePassword 验证密码
func (uc *UserUsecase) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不能少于8个字符")
	}
	return nil
}

// hashPassword 加密密码
func (uc *UserUsecase) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkPassword 验证密码
func (uc *UserUsecase) checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
