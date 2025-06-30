package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// === 类型定义 ===

// DBStats 数据库连接统计信息
type DBStats struct {
	MaxOpenConnections int           `json:"max_open_connections"`
	OpenConnections    int           `json:"open_connections"`
	InUseConnections   int           `json:"in_use_connections"`
	IdleConnections    int           `json:"idle_connections"`
	WaitCount          int64         `json:"wait_count"`
	WaitDuration       time.Duration `json:"wait_duration"`
	MaxIdleClosed      int64         `json:"max_idle_closed"`
	MaxLifetimeClosed  int64         `json:"max_lifetime_closed"`
}

// TransactionFunc 事务函数类型
type TransactionFunc func(tx *gorm.DB) error

// MigrationFunc 迁移函数类型
type MigrationFunc func(db *gorm.DB) error

// DBManager 数据库管理器
type DBManager struct {
	db       *gorm.DB
	config   *config.DatabaseConfig
	mu       sync.RWMutex
	isReady  bool
	lastPing time.Time
}

// === 全局变量 ===

var (
	// 全局数据库实例
	DB      *gorm.DB
	manager *DBManager
	once    sync.Once
)

// === 配置验证 ===

// validateDBConfig 验证数据库配置
func validateDBConfig(cfg *config.DatabaseConfig) error {
	if cfg == nil {
		return errors.New("database config cannot be nil")
	}

	if cfg.Host == "" {
		return errors.New("database host is required")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return errors.New("database port must be between 1 and 65535")
	}

	if cfg.Username == "" {
		return errors.New("database username is required")
	}

	if cfg.DBName == "" {
		return errors.New("database name is required")
	}

	if cfg.MaxIdleConns < 0 {
		return errors.New("max idle connections cannot be negative")
	}

	if cfg.MaxOpenConns <= 0 {
		return errors.New("max open connections must be positive")
	}

	if cfg.MaxIdleConns > cfg.MaxOpenConns {
		return errors.New("max idle connections cannot exceed max open connections")
	}

	return nil
}

// === 初始化函数 ===

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	var initErr error

	once.Do(func() {
		// 验证配置
		if err := validateDBConfig(&cfg.Database); err != nil {
			initErr = fmt.Errorf("invalid database config: %w", err)
			return
		}

		// 创建数据库管理器
		manager = &DBManager{
			config: &cfg.Database,
		}

		// 初始化连接
		if err := manager.connect(); err != nil {
			initErr = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		// 设置全局实例
		DB = manager.db

		// 自动迁移（如果配置启用）
		if cfg.Database.AutoMigrate {
			if err := manager.autoMigrate(); err != nil {
				pkglog.LogWarn("Auto migration failed", zap.Error(err))
			}
		}

		pkglog.LogInfo("Database initialized successfully",
			zap.String("host", cfg.Database.Host),
			zap.Int("port", cfg.Database.Port),
			zap.String("database", cfg.Database.DBName))
	})

	return initErr
}

// InitDBWithRetry 带重试机制的数据库初始化
func InitDBWithRetry(cfg *config.Config, maxRetries int, retryInterval time.Duration) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := InitDB(cfg); err != nil {
			lastErr = err
			pkglog.LogWarn("Database connection attempt failed",
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Error(err))

			if i < maxRetries-1 {
				time.Sleep(retryInterval)
				// 重置once以允许重试
				once = sync.Once{}
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, lastErr)
}

// === 数据库管理器方法 ===

// connect 建立数据库连接
func (m *DBManager) connect() error {
	// 构建DSN
	dsn := m.config.GetDsn()

	// 配置GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 获取底层sql.DB实例
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 配置连接池
	if err := m.configureConnectionPool(sqlDB); err != nil {
		return fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// 设置数据库实例
	m.mu.Lock()
	m.db = db
	m.isReady = true
	m.lastPing = time.Now()
	m.mu.Unlock()

	// 执行健康检查
	if err := m.ping(); err != nil {
		return fmt.Errorf("initial health check failed: %w", err)
	}

	pkglog.LogDBOperation("CONNECT", "postgresql", time.Since(m.lastPing), nil)

	return nil
}

// configureConnectionPool 配置连接池
func (m *DBManager) configureConnectionPool(sqlDB *sql.DB) error {
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)

	// 设置连接最大生命周期
	if m.config.ConnMaxLifetime != "" {
		lifetime, err := time.ParseDuration(m.config.ConnMaxLifetime)
		if err != nil {
			return fmt.Errorf("invalid connection max lifetime: %w", err)
		}
		sqlDB.SetConnMaxLifetime(lifetime)
	}

	return nil
}

// ping 执行数据库ping检查
func (m *DBManager) ping() error {
	m.mu.RLock()
	if !m.isReady || m.db == nil {
		m.mu.RUnlock()
		return errors.New("database not initialized")
	}
	db := m.db
	m.mu.RUnlock()

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	m.mu.Lock()
	m.lastPing = time.Now()
	m.mu.Unlock()

	return nil
}

// autoMigrate 执行自动迁移
func (m *DBManager) autoMigrate() error {
	// 这里可以添加您的模型
	// 例如: return m.db.AutoMigrate(&User{}, &Order{})

	pkglog.LogInfo("Auto migration completed")
	return nil
}

// === 公共API函数 ===

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if manager != nil {
		manager.mu.RLock()
		defer manager.mu.RUnlock()
		return manager.db
	}
	return DB
}

// IsReady 检查数据库是否就绪
func IsReady() bool {
	if manager != nil {
		manager.mu.RLock()
		defer manager.mu.RUnlock()
		return manager.isReady
	}
	return DB != nil
}

// CheckDBHealth 检查数据库健康状态
func CheckDBHealth() error {
	if manager != nil {
		return manager.ping()
	}

	// 向后兼容的检查
	if DB == nil {
		return errors.New("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// GetDBStats 获取数据库连接统计
func GetDBStats() (*DBStats, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	stats := sqlDB.Stats()

	return &DBStats{
		MaxOpenConnections: stats.MaxOpenConnections,
		OpenConnections:    stats.OpenConnections,
		InUseConnections:   stats.InUse,
		IdleConnections:    stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
	}, nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if manager != nil {
		manager.mu.Lock()
		defer manager.mu.Unlock()

		if manager.db != nil {
			sqlDB, err := manager.db.DB()
			if err != nil {
				return fmt.Errorf("failed to get sql.DB instance: %w", err)
			}

			if err := sqlDB.Close(); err != nil {
				return fmt.Errorf("failed to close database: %w", err)
			}

			manager.db = nil
			manager.isReady = false

			pkglog.LogInfo("Database connection closed")
			return nil
		}
	}

	// 向后兼容
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB instance: %w", err)
		}
		return sqlDB.Close()
	}

	return nil
}

// === 事务管理 ===

// WithTransaction 执行事务
func WithTransaction(fn TransactionFunc) error {
	db := GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	startTime := time.Now()

	tx := db.Begin()
	if tx.Error != nil {
		pkglog.LogDBOperation("BEGIN_TX", "postgresql", time.Since(startTime), tx.Error)
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			pkglog.LogError("Transaction panicked, rolled back",
				zap.Any("panic", r),
				zap.Duration("duration", time.Since(startTime)))
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		pkglog.LogDBOperation("ROLLBACK_TX", "postgresql", time.Since(startTime), err)
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		pkglog.LogDBOperation("COMMIT_TX", "postgresql", time.Since(startTime), err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	pkglog.LogDBOperation("COMMIT_TX", "postgresql", time.Since(startTime), nil)
	return nil
}

// === 迁移管理 ===

// RunMigrations 运行数据库迁移
func RunMigrations(migrations []MigrationFunc) error {
	db := GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	for i, migration := range migrations {
		startTime := time.Now()

		if err := migration(db); err != nil {
			pkglog.LogDBOperation("MIGRATION", fmt.Sprintf("step_%d", i+1), time.Since(startTime), err)
			return fmt.Errorf("migration step %d failed: %w", i+1, err)
		}

		pkglog.LogDBOperation("MIGRATION", fmt.Sprintf("step_%d", i+1), time.Since(startTime), nil)
	}

	pkglog.LogInfo("All migrations completed successfully",
		zap.Int("total_migrations", len(migrations)))

	return nil
}

// === 表管理功能 ===

// TableColumn 表列定义
type TableColumn struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	NotNull      bool   `json:"not_null"`
	PrimaryKey   bool   `json:"primary_key"`
	Unique       bool   `json:"unique"`
	DefaultValue string `json:"default_value,omitempty"`
}

// TableDefinition 表定义结构
type TableDefinition struct {
	Name    string        `json:"name"`
	Columns []TableColumn `json:"columns"`
}

// CreateTable 创建表
func CreateTable(tableDef *TableDefinition) error {
	db := GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	startTime := time.Now()

	// 构建CREATE TABLE SQL语句
	sql := buildCreateTableSQL(tableDef)

	// 执行SQL
	if err := db.Exec(sql).Error; err != nil {
		pkglog.LogDBOperation("CREATE_TABLE", tableDef.Name, time.Since(startTime), err)
		return fmt.Errorf("failed to create table %s: %w", tableDef.Name, err)
	}

	pkglog.LogDBOperation("CREATE_TABLE", tableDef.Name, time.Since(startTime), nil)
	pkglog.LogInfo("Table created successfully", zap.String("table", tableDef.Name))

	return nil
}

// DropTable 删除表
func DropTable(tableName string) error {
	db := GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	startTime := time.Now()

	// 构建DROP TABLE SQL语句
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

	// 执行SQL
	if err := db.Exec(sql).Error; err != nil {
		pkglog.LogDBOperation("DROP_TABLE", tableName, time.Since(startTime), err)
		return fmt.Errorf("failed to drop table %s: %w", tableName, err)
	}

	pkglog.LogDBOperation("DROP_TABLE", tableName, time.Since(startTime), nil)
	pkglog.LogInfo("Table dropped successfully", zap.String("table", tableName))

	return nil
}

// TableExists 检查表是否存在
func TableExists(tableName string) (bool, error) {
	db := GetDB()
	if db == nil {
		return false, errors.New("database not initialized")
	}

	startTime := time.Now()

	var count int64
	sql := `SELECT COUNT(*) FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = ?`

	if err := db.Raw(sql, tableName).Scan(&count).Error; err != nil {
		pkglog.LogDBOperation("CHECK_TABLE", tableName, time.Since(startTime), err)
		return false, fmt.Errorf("failed to check table existence %s: %w", tableName, err)
	}

	pkglog.LogDBOperation("CHECK_TABLE", tableName, time.Since(startTime), nil)

	return count > 0, nil
}

// CreateTestTable 创建测试表（使用时间戳命名）
func CreateTestTable() (string, error) {
	timestamp := time.Now().UnixNano()
	tableName := fmt.Sprintf("test_%d", timestamp)

	// 定义测试表结构
	tableDef := &TableDefinition{
		Name: tableName,
		Columns: []TableColumn{
			{
				Name:       "id",
				Type:       "SERIAL",
				PrimaryKey: true,
				NotNull:    true,
			},
			{
				Name:    "name",
				Type:    "VARCHAR(255)",
				NotNull: true,
			},
			{
				Name:   "email",
				Type:   "VARCHAR(255)",
				Unique: true,
			},
			{
				Name:         "created_at",
				Type:         "TIMESTAMP",
				DefaultValue: "CURRENT_TIMESTAMP",
			},
			{
				Name:         "updated_at",
				Type:         "TIMESTAMP",
				DefaultValue: "CURRENT_TIMESTAMP",
			},
		},
	}

	if err := CreateTable(tableDef); err != nil {
		return "", fmt.Errorf("failed to create test table: %w", err)
	}

	return tableName, nil
}

// buildCreateTableSQL 构建CREATE TABLE SQL语句
func buildCreateTableSQL(tableDef *TableDefinition) string {
	var columns []string

	for _, col := range tableDef.Columns {
		columnSQL := fmt.Sprintf("%s %s", col.Name, col.Type)

		if col.PrimaryKey {
			columnSQL += " PRIMARY KEY"
		}

		if col.NotNull && !col.PrimaryKey {
			columnSQL += " NOT NULL"
		}

		if col.Unique && !col.PrimaryKey {
			columnSQL += " UNIQUE"
		}

		if col.DefaultValue != "" {
			columnSQL += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
		}

		columns = append(columns, columnSQL)
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)",
		tableDef.Name, strings.Join(columns, ", "))
}

// GetTableInfo 获取表信息
func GetTableInfo(tableName string) (*TableDefinition, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	startTime := time.Now()

	// 查询表结构信息
	sql := `
		SELECT 
			column_name,
			data_type,
			is_nullable,
			column_default
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = ?
		ORDER BY ordinal_position
	`

	rows, err := db.Raw(sql, tableName).Rows()
	if err != nil {
		pkglog.LogDBOperation("GET_TABLE_INFO", tableName, time.Since(startTime), err)
		return nil, fmt.Errorf("failed to get table info %s: %w", tableName, err)
	}
	defer rows.Close()

	tableDef := &TableDefinition{
		Name:    tableName,
		Columns: []TableColumn{},
	}

	for rows.Next() {
		var colName, dataType, isNullable string
		var columnDefault *string

		if err := rows.Scan(&colName, &dataType, &isNullable, &columnDefault); err != nil {
			continue
		}

		column := TableColumn{
			Name:    colName,
			Type:    dataType,
			NotNull: isNullable == "NO",
		}

		if columnDefault != nil {
			column.DefaultValue = *columnDefault
		}

		tableDef.Columns = append(tableDef.Columns, column)
	}

	pkglog.LogDBOperation("GET_TABLE_INFO", tableName, time.Since(startTime), nil)

	return tableDef, nil
}

// === 工具函数 ===

// Ping 手动执行数据库ping
func Ping() error {
	return CheckDBHealth()
}

// GetLastPingTime 获取最后一次ping时间
func GetLastPingTime() time.Time {
	if manager != nil {
		manager.mu.RLock()
		defer manager.mu.RUnlock()
		return manager.lastPing
	}
	return time.Time{}
}
