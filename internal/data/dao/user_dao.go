package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz"
	"github.com/YoungBoyGod/OneGoServer/internal/data/postgres"
	"github.com/YoungBoyGod/OneGoServer/internal/data/redis"
	"gorm.io/gorm"
)

// userDAO 用户数据访问对象
type userDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建用户DAO实例
func NewUserDAO() biz.UserRepo {
	return &userDAO{
		db: postgres.GetDB(),
	}
}

// CreateUser 创建用户
func (d *userDAO) CreateUser(user *biz.User) error {
	ctx := context.Background()
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 检查用户名是否已存在
	exists, err := d.UserExistsByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return fmt.Errorf("username already exists: %s", user.Username)
	}

	// 检查邮箱是否已存在
	exists, err = d.UserExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return fmt.Errorf("email already exists: %s", user.Email)
	}

	// 创建用户
	if err := d.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// 清理相关缓存
	d.clearUserCache(ctx, user.ID, user.Username, user.Email)

	return nil
}

// GetUserByID 根据ID获取用户
func (d *userDAO) GetUserByID(id uint) (*biz.User, error) {
	ctx := context.Background()
	if d.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 先尝试从缓存获取
	cacheKey := fmt.Sprintf("user:id:%d", id)
	var user biz.User
	if err := redis.GetObject(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 从数据库获取
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	// 缓存结果
	redis.SetObject(ctx, cacheKey, &user, 30*time.Minute)

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (d *userDAO) GetUserByUsername(username string) (*biz.User, error) {
	ctx := context.Background()
	if d.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 先尝试从缓存获取
	cacheKey := fmt.Sprintf("user:username:%s", username)
	var user biz.User
	if err := redis.GetObject(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 从数据库获取
	if err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	// 缓存结果
	redis.SetObject(ctx, cacheKey, &user, 30*time.Minute)

	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (d *userDAO) GetUserByEmail(email string) (*biz.User, error) {
	ctx := context.Background()
	if d.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	// 先尝试从缓存获取
	cacheKey := fmt.Sprintf("user:email:%s", email)
	var user biz.User
	if err := redis.GetObject(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	// 从数据库获取
	if err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// 缓存结果
	redis.SetObject(ctx, cacheKey, &user, 30*time.Minute)

	return &user, nil
}

// UpdateUser 更新用户信息
func (d *userDAO) UpdateUser(user *biz.User) error {
	ctx := context.Background()
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 更新数据库
	if err := d.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 清理相关缓存
	d.clearUserCache(ctx, user.ID, user.Username, user.Email)

	return nil
}

// DeleteUser 删除用户
func (d *userDAO) DeleteUser(id uint) error {
	ctx := context.Background()
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// 先获取用户信息用于清理缓存
	user, err := d.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user for deletion: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// 软删除用户
	if err := d.db.WithContext(ctx).Delete(&biz.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 清理相关缓存
	d.clearUserCache(ctx, user.ID, user.Username, user.Email)

	return nil
}

// ListUsers 分页获取用户列表
func (d *userDAO) ListUsers(offset, limit int, filter *biz.UserFilter) ([]*biz.User, int64, error) {
	ctx := context.Background()
	if d.db == nil {
		return nil, 0, fmt.Errorf("database not initialized")
	}

	query := d.db.WithContext(ctx).Model(&biz.User{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "role":
			query = query.Where("role = ?", value)
		case "search":
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
				searchTerm, searchTerm, searchTerm)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 分页查询
	var users []*biz.User
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// UserExistsByUsername 检查用户名是否存在
func (d *userDAO) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	if d.db == nil {
		return false, fmt.Errorf("database not initialized")
	}

	var count int64
	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return count > 0, nil
}

// UserExistsByEmail 检查邮箱是否存在
func (d *userDAO) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	if d.db == nil {
		return false, fmt.Errorf("database not initialized")
	}

	var count int64
	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// UpdateUserLoginInfo 更新用户登录信息
func (d *userDAO) UpdateUserLoginInfo(ctx context.Context, userID uint, loginIP string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	updates := map[string]interface{}{
		"last_login_at": time.Now(),
		"last_login_ip": loginIP,
	}

	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user login info: %w", err)
	}

	// 清理用户缓存
	d.clearUserCacheByID(ctx, userID)

	return nil
}

// UpdateUserPassword 更新用户密码
func (d *userDAO) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error {
	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	updates := map[string]interface{}{
		"password":   hashedPassword,
		"updated_at": time.Now(),
	}

	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	// 清理用户缓存
	d.clearUserCacheByID(ctx, userID)

	return nil
}

// GetUserStats 获取用户统计信息
func (d *userDAO) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	if d.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	stats := make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	if err := d.db.WithContext(ctx).Model(&biz.User{}).Count(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}
	stats["total_users"] = totalUsers

	// 活跃用户数
	var activeUsers int64
	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("status = ?", biz.UserStatusActive).Count(&activeUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count active users: %w", err)
	}
	stats["active_users"] = activeUsers

	// 今日新增用户
	today := time.Now().Truncate(24 * time.Hour)
	var todayUsers int64
	if err := d.db.WithContext(ctx).Model(&biz.User{}).Where("created_at >= ?", today).Count(&todayUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count today users: %w", err)
	}
	stats["today_users"] = todayUsers

	// 按角色统计
	var roleStats []struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	if err := d.db.WithContext(ctx).Model(&biz.User{}).
		Select("role, count(*) as count").
		Group("role").
		Scan(&roleStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get role stats: %w", err)
	}
	stats["role_stats"] = roleStats

	return stats, nil
}

// clearUserCache 清理用户相关缓存
func (d *userDAO) clearUserCache(ctx context.Context, userID uint, username, email string) {
	// 清理所有相关缓存键
	cacheKeys := []string{
		fmt.Sprintf("user:id:%d", userID),
		fmt.Sprintf("user:username:%s", username),
		fmt.Sprintf("user:email:%s", email),
	}

	for _, key := range cacheKeys {
		redis.Del(ctx, key)
	}
}

// clearUserCacheByID 根据ID清理用户缓存
func (d *userDAO) clearUserCacheByID(ctx context.Context, userID uint) {
	// 先获取用户信息
	user, err := d.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return
	}

	d.clearUserCache(ctx, user.ID, user.Username, user.Email)
}
