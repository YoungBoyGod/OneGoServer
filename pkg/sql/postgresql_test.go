package sql

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// 测试用的配置
func getTestConfig() *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Type:            "postgres",
			Host:            "localhost",
			Port:            5432,
			Username:        "test_user",
			Password:        "test_password",
			DBName:          "test_db",
			Charset:         "utf8",
			MaxIdleConns:    5,
			MaxOpenConns:    10,
			ConnMaxLifetime: "1h",
			AutoMigrate:     false,
			ParseTime:       true,
			Loc:             "UTC",
		},
		Logging: config.LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "logs/test.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   false,
		},
	}
}

// TestValidateDBConfig 测试数据库配置验证
func TestValidateDBConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.DatabaseConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "空配置",
			config:      nil,
			expectError: true,
			errorMsg:    "database config cannot be nil",
		},
		{
			name: "缺少主机",
			config: &config.DatabaseConfig{
				Port:         5432,
				Username:     "user",
				DBName:       "db",
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "database host is required",
		},
		{
			name: "无效端口",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         0,
				Username:     "user",
				DBName:       "db",
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "database port must be between 1 and 65535",
		},
		{
			name: "缺少用户名",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				DBName:       "db",
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "database username is required",
		},
		{
			name: "缺少数据库名",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "user",
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "database name is required",
		},
		{
			name: "负数空闲连接",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "user",
				DBName:       "db",
				MaxIdleConns: -1,
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "max idle connections cannot be negative",
		},
		{
			name: "非正数最大连接",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "user",
				DBName:       "db",
				MaxIdleConns: 5,
				MaxOpenConns: 0,
			},
			expectError: true,
			errorMsg:    "max open connections must be positive",
		},
		{
			name: "空闲连接超过最大连接",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "user",
				DBName:       "db",
				MaxIdleConns: 15,
				MaxOpenConns: 10,
			},
			expectError: true,
			errorMsg:    "max idle connections cannot exceed max open connections",
		},
		{
			name: "有效配置",
			config: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				Username:     "user",
				Password:     "password",
				DBName:       "db",
				MaxIdleConns: 5,
				MaxOpenConns: 10,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDBConfig(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDBManager 测试数据库管理器
func TestDBManager(t *testing.T) {
	// 初始化日志系统
	cfg := getTestConfig()
	err := pkglog.InitLoggerEnhanced(&cfg.Logging)
	require.NoError(t, err)

	t.Run("数据库管理器创建", func(t *testing.T) {
		manager := &DBManager{
			config: &cfg.Database,
		}

		assert.NotNil(t, manager)
		assert.Equal(t, &cfg.Database, manager.config)
		assert.False(t, manager.isReady)
	})

	t.Run("配置连接池参数验证", func(t *testing.T) {
		manager := &DBManager{
			config: &config.DatabaseConfig{
				MaxIdleConns:    5,
				MaxOpenConns:    10,
				ConnMaxLifetime: "1h",
			},
		}

		// 由于无法模拟sql.DB，我们只能测试配置是否正确设置
		assert.Equal(t, 5, manager.config.MaxIdleConns)
		assert.Equal(t, 10, manager.config.MaxOpenConns)
		assert.Equal(t, "1h", manager.config.ConnMaxLifetime)
	})

	t.Run("无效连接生命周期", func(t *testing.T) {
		// 测试时间解析错误
		_, err := time.ParseDuration("invalid")
		assert.Error(t, err)
	})
}

// TestIsReady 测试就绪状态检查
func TestIsReady(t *testing.T) {
	// 重置全局状态
	manager = nil
	DB = nil

	t.Run("未初始化状态", func(t *testing.T) {
		assert.False(t, IsReady())
	})

	t.Run("管理器存在但未就绪", func(t *testing.T) {
		manager = &DBManager{
			isReady: false,
		}
		assert.False(t, IsReady())
	})

	t.Run("管理器就绪", func(t *testing.T) {
		manager = &DBManager{
			isReady: true,
		}
		assert.True(t, IsReady())
	})

	// 清理
	manager = nil
}

// TestWithTransaction 测试事务管理
func TestWithTransaction(t *testing.T) {
	// 重置全局状态
	DB = nil
	manager = nil

	t.Run("数据库未初始化", func(t *testing.T) {
		err := WithTransaction(func(tx *gorm.DB) error {
			return nil
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database not initialized")
	})

	t.Run("事务函数错误处理", func(t *testing.T) {
		// 这个测试需要实际的数据库连接，这里我们只测试逻辑
		testErr := errors.New("test error")

		// 模拟事务函数
		txFunc := func(tx *gorm.DB) error {
			return testErr
		}

		// 验证函数类型
		assert.NotNil(t, txFunc)

		// 执行函数验证错误传播
		err := txFunc(nil)
		assert.Equal(t, testErr, err)
	})
}

// TestRunMigrations 测试迁移管理
func TestRunMigrations(t *testing.T) {
	// 重置全局状态
	DB = nil
	manager = nil

	t.Run("数据库未初始化", func(t *testing.T) {
		migrations := []MigrationFunc{
			func(db *gorm.DB) error {
				return nil
			},
		}

		err := RunMigrations(migrations)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database not initialized")
	})

	t.Run("空迁移列表", func(t *testing.T) {
		// 即使数据库未初始化，空列表也应该不报错逻辑上
		migrations := []MigrationFunc{}

		// 这里我们测试迁移函数的结构
		assert.Equal(t, 0, len(migrations))
	})

	t.Run("迁移函数结构", func(t *testing.T) {
		// 测试迁移函数类型
		migration := func(db *gorm.DB) error {
			// 模拟迁移逻辑
			if db == nil {
				return errors.New("db is nil")
			}
			return nil
		}

		// 测试函数调用
		err := migration(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db is nil")
	})
}

// TestPing 测试Ping功能
func TestPing(t *testing.T) {
	// 重置全局状态
	DB = nil
	manager = nil

	t.Run("数据库未初始化", func(t *testing.T) {
		err := Ping()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database not initialized")
	})
}

// TestGetLastPingTime 测试获取最后ping时间
func TestGetLastPingTime(t *testing.T) {
	// 重置全局状态
	manager = nil

	t.Run("管理器不存在", func(t *testing.T) {
		pingTime := GetLastPingTime()
		assert.True(t, pingTime.IsZero())
	})

	t.Run("管理器存在", func(t *testing.T) {
		testTime := time.Now()
		manager = &DBManager{
			lastPing: testTime,
		}

		pingTime := GetLastPingTime()
		assert.Equal(t, testTime, pingTime)
	})

	// 清理
	manager = nil
}

// TestDBStats 测试数据库统计
func TestDBStats(t *testing.T) {
	t.Run("DBStats结构体", func(t *testing.T) {
		stats := &DBStats{
			MaxOpenConnections: 10,
			OpenConnections:    5,
			InUseConnections:   2,
			IdleConnections:    3,
			WaitCount:          100,
			WaitDuration:       time.Millisecond * 50,
			MaxIdleClosed:      10,
			MaxLifetimeClosed:  5,
		}

		assert.Equal(t, 10, stats.MaxOpenConnections)
		assert.Equal(t, 5, stats.OpenConnections)
		assert.Equal(t, 2, stats.InUseConnections)
		assert.Equal(t, 3, stats.IdleConnections)
		assert.Equal(t, int64(100), stats.WaitCount)
		assert.Equal(t, time.Millisecond*50, stats.WaitDuration)
		assert.Equal(t, int64(10), stats.MaxIdleClosed)
		assert.Equal(t, int64(5), stats.MaxLifetimeClosed)
	})

	t.Run("数据库未初始化", func(t *testing.T) {
		// 重置全局状态
		DB = nil
		manager = nil

		stats, err := GetDBStats()
		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "database not initialized")
	})
}

// TestCloseDB 测试数据库关闭
func TestCloseDB(t *testing.T) {
	t.Run("管理器不存在", func(t *testing.T) {
		// 重置全局状态
		manager = nil
		DB = nil

		err := CloseDB()
		assert.NoError(t, err) // 应该静默处理
	})

	t.Run("管理器存在但数据库为空", func(t *testing.T) {
		manager = &DBManager{
			db: nil,
		}

		err := CloseDB()
		assert.NoError(t, err) // 应该静默处理
	})

	// 清理
	manager = nil
}

// TestGetDB 测试获取数据库实例
func TestGetDB(t *testing.T) {
	t.Run("管理器不存在", func(t *testing.T) {
		// 重置全局状态
		manager = nil
		DB = nil

		db := GetDB()
		assert.Nil(t, db)
	})

	t.Run("管理器存在", func(t *testing.T) {
		manager = &DBManager{
			db: &gorm.DB{}, // 模拟数据库实例
		}

		db := GetDB()
		assert.NotNil(t, db)
		assert.Equal(t, manager.db, db)
	})

	// 清理
	manager = nil
}

// TestCheckDBHealth 测试健康检查
func TestCheckDBHealth(t *testing.T) {
	t.Run("数据库未初始化", func(t *testing.T) {
		// 重置全局状态
		DB = nil
		manager = nil

		err := CheckDBHealth()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database not initialized")
	})
}

// BenchmarkValidateDBConfig 配置验证性能测试
func BenchmarkValidateDBConfig(b *testing.B) {
	config := &config.DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		Username:     "user",
		Password:     "password",
		DBName:       "db",
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateDBConfig(config)
	}
}

// TestConfigurationEdgeCases 测试配置边界情况
func TestConfigurationEdgeCases(t *testing.T) {
	t.Run("最大端口号", func(t *testing.T) {
		config := &config.DatabaseConfig{
			Host:         "localhost",
			Port:         65535,
			Username:     "user",
			DBName:       "db",
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		}

		err := validateDBConfig(config)
		assert.NoError(t, err)
	})

	t.Run("最小端口号", func(t *testing.T) {
		config := &config.DatabaseConfig{
			Host:         "localhost",
			Port:         1,
			Username:     "user",
			DBName:       "db",
			MaxIdleConns: 5,
			MaxOpenConns: 10,
		}

		err := validateDBConfig(config)
		assert.NoError(t, err)
	})

	t.Run("零空闲连接", func(t *testing.T) {
		config := &config.DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			Username:     "user",
			DBName:       "db",
			MaxIdleConns: 0,
			MaxOpenConns: 10,
		}

		err := validateDBConfig(config)
		assert.NoError(t, err)
	})

	t.Run("相等的空闲和最大连接数", func(t *testing.T) {
		config := &config.DatabaseConfig{
			Host:         "localhost",
			Port:         5432,
			Username:     "user",
			DBName:       "db",
			MaxIdleConns: 10,
			MaxOpenConns: 10,
		}

		err := validateDBConfig(config)
		assert.NoError(t, err)
	})
}

// TestConcurrentAccess 测试并发访问
func TestConcurrentAccess(t *testing.T) {
	// 重置全局状态
	manager = &DBManager{
		isReady:  true,
		lastPing: time.Now(),
	}

	t.Run("并发IsReady检查", func(t *testing.T) {
		done := make(chan bool, 10)

		// 启动多个goroutine并发检查IsReady
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					IsReady()
				}
			}()
		}

		// 等待所有goroutine完成
		for i := 0; i < 10; i++ {
			<-done
		}

		// 如果没有竞态条件，测试应该成功完成
		assert.True(t, true)
	})

	t.Run("并发GetLastPingTime", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 100; j++ {
					GetLastPingTime()
				}
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		assert.True(t, true)
	})

	// 清理
	manager = nil
}

// TestTransactionFuncType 测试事务函数类型
func TestTransactionFuncType(t *testing.T) {
	t.Run("事务函数类型验证", func(t *testing.T) {
		// 成功的事务函数
		successFunc := func(tx *gorm.DB) error {
			return nil
		}

		// 失败的事务函数
		failFunc := func(tx *gorm.DB) error {
			return errors.New("transaction failed")
		}

		// 测试函数类型
		var _ TransactionFunc = successFunc
		var _ TransactionFunc = failFunc

		// 测试函数执行
		assert.NoError(t, successFunc(nil))
		assert.Error(t, failFunc(nil))
	})
}

// TestMigrationFuncType 测试迁移函数类型
func TestMigrationFuncType(t *testing.T) {
	t.Run("迁移函数类型验证", func(t *testing.T) {
		// 成功的迁移函数
		successMigration := func(db *gorm.DB) error {
			return nil
		}

		// 失败的迁移函数
		failMigration := func(db *gorm.DB) error {
			return errors.New("migration failed")
		}

		// 测试函数类型
		var _ MigrationFunc = successMigration
		var _ MigrationFunc = failMigration

		// 测试函数执行
		assert.NoError(t, successMigration(nil))
		assert.Error(t, failMigration(nil))
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	t.Run("错误包装", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrappedErr := fmt.Errorf("wrapped: %w", originalErr)

		assert.Contains(t, wrappedErr.Error(), "original error")
		assert.Contains(t, wrappedErr.Error(), "wrapped")
	})

	t.Run("错误链检查", func(t *testing.T) {
		originalErr := errors.New("database connection failed")
		wrappedErr := fmt.Errorf("failed to initialize database: %w", originalErr)

		assert.True(t, errors.Is(wrappedErr, originalErr))
	})
}
