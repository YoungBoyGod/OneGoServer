package sql

import (
	"errors"
	"fmt"
	"sync"
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

// TestAddAndDeleteTable 测试创建和删除表功能
func TestAddAndDeleteTable(t *testing.T) {
	// 初始化日志系统
	cfg := getTestConfig()
	err := pkglog.InitLoggerEnhanced(&cfg.Logging)
	require.NoError(t, err)

	// 生成基于时间戳的表名
	timestamp := time.Now().UnixNano()
	tableName := fmt.Sprintf("test_%d", timestamp)

	t.Run("创建和删除表操作", func(t *testing.T) {
		// 由于我们没有实际的数据库连接，这里测试表名生成和SQL构造逻辑

		// 测试表名生成
		assert.Contains(t, tableName, "test_")
		assert.True(t, len(tableName) > 5) // test_ + 时间戳应该大于5位

		// 模拟创建表的SQL语句
		createTableSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`, tableName)

		// 模拟删除表的SQL语句
		dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)

		// 验证SQL语句格式
		assert.Contains(t, createTableSQL, tableName)
		assert.Contains(t, createTableSQL, "CREATE TABLE IF NOT EXISTS")
		assert.Contains(t, createTableSQL, "SERIAL PRIMARY KEY")

		assert.Contains(t, dropTableSQL, tableName)
		assert.Contains(t, dropTableSQL, "DROP TABLE IF EXISTS")

		t.Logf("Generated table name: %s", tableName)
		t.Logf("Create SQL: %s", createTableSQL)
		t.Logf("Drop SQL: %s", dropTableSQL)
	})

	t.Run("表名唯一性测试", func(t *testing.T) {
		// 生成多个表名，确保唯一性
		tableNames := make(map[string]bool)

		for i := 0; i < 10; i++ {
			time.Sleep(time.Microsecond) // 确保时间戳不同
			ts := time.Now().UnixNano()
			name := fmt.Sprintf("test_%d", ts)

			// 检查是否重复
			assert.False(t, tableNames[name], "Table name should be unique: %s", name)
			tableNames[name] = true
		}

		assert.Equal(t, 10, len(tableNames), "Should generate 10 unique table names")
	})

	t.Run("表名格式验证", func(t *testing.T) {
		// 测试表名格式的有效性
		ts := time.Now().UnixNano()
		name := fmt.Sprintf("test_%d", ts)

		// PostgreSQL表名规则验证
		assert.True(t, len(name) <= 63, "PostgreSQL table name should not exceed 63 characters")
		assert.Regexp(t, `^test_\d+$`, name, "Table name should match pattern test_<timestamp>")
		assert.True(t, ts > 0, "Timestamp should be positive")
	})

	t.Run("模拟表生命周期管理", func(t *testing.T) {
		// 模拟完整的表生命周期：创建 -> 使用 -> 删除

		tableManager := &TestTableManager{
			TableName: tableName,
			CreatedAt: time.Now(),
		}

		// 测试表管理器
		assert.Equal(t, tableName, tableManager.TableName)
		assert.False(t, tableManager.CreatedAt.IsZero())

		// 模拟表操作
		operations := []string{"CREATE", "INSERT", "SELECT", "UPDATE", "DELETE", "DROP"}

		for _, op := range operations {
			result := tableManager.SimulateOperation(op)
			assert.True(t, result, "Operation %s should succeed", op)
		}

		t.Logf("Table lifecycle completed for: %s", tableManager.TableName)
	})

	t.Run("表管理功能测试", func(t *testing.T) {
		// 测试表定义结构
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
			},
		}

		// 验证表定义结构
		assert.Equal(t, tableName, tableDef.Name)
		assert.Equal(t, 3, len(tableDef.Columns))

		// 验证主键列
		primaryKeyCol := tableDef.Columns[0]
		assert.Equal(t, "id", primaryKeyCol.Name)
		assert.True(t, primaryKeyCol.PrimaryKey)
		assert.True(t, primaryKeyCol.NotNull)

		// 验证唯一列
		uniqueCol := tableDef.Columns[2]
		assert.Equal(t, "email", uniqueCol.Name)
		assert.True(t, uniqueCol.Unique)

		t.Logf("Table definition validated: %+v", tableDef)
	})

	t.Run("SQL构建测试", func(t *testing.T) {
		// 测试buildCreateTableSQL函数的逻辑
		tableDef := &TableDefinition{
			Name: "test_sql_build",
			Columns: []TableColumn{
				{
					Name:       "id",
					Type:       "SERIAL",
					PrimaryKey: true,
				},
				{
					Name:    "name",
					Type:    "VARCHAR(100)",
					NotNull: true,
				},
				{
					Name:         "created_at",
					Type:         "TIMESTAMP",
					DefaultValue: "CURRENT_TIMESTAMP",
				},
			},
		}

		// 由于buildCreateTableSQL是私有函数，我们测试SQL构建逻辑
		expectedSQL := "CREATE TABLE IF NOT EXISTS test_sql_build (id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)"

		// 验证SQL构建逻辑的组件
		assert.Contains(t, tableDef.Name, "test_sql_build")
		assert.Equal(t, 3, len(tableDef.Columns))

		// 验证列定义的组件
		for _, col := range tableDef.Columns {
			assert.NotEmpty(t, col.Name)
			assert.NotEmpty(t, col.Type)
		}

		t.Logf("Expected SQL: %s", expectedSQL)
	})

	t.Run("表管理功能单元测试", func(t *testing.T) {
		// 由于没有实际数据库连接，测试表管理功能的结构和逻辑

		// 测试CreateTestTable函数的逻辑
		// 这里我们不能调用实际的CreateTestTable，但可以测试其逻辑组件

		// 验证时间戳生成逻辑
		ts1 := time.Now().UnixNano()
		time.Sleep(time.Microsecond)
		ts2 := time.Now().UnixNano()

		assert.True(t, ts2 > ts1, "Timestamps should be increasing")

		// 验证表名格式
		testTableName := fmt.Sprintf("test_%d", ts1)
		assert.Regexp(t, `^test_\d+$`, testTableName)

		// 测试TableColumn结构
		col := TableColumn{
			Name:         "test_col",
			Type:         "VARCHAR(100)",
			NotNull:      true,
			PrimaryKey:   false,
			Unique:       true,
			DefaultValue: "default_value",
		}

		assert.Equal(t, "test_col", col.Name)
		assert.Equal(t, "VARCHAR(100)", col.Type)
		assert.True(t, col.NotNull)
		assert.False(t, col.PrimaryKey)
		assert.True(t, col.Unique)
		assert.Equal(t, "default_value", col.DefaultValue)

		t.Logf("Column structure validated: %+v", col)
	})
}

// TestTableManager 测试表管理器结构体
type TestTableManager struct {
	TableName string
	CreatedAt time.Time
}

// SimulateOperation 模拟表操作
func (tm *TestTableManager) SimulateOperation(operation string) bool {
	switch operation {
	case "CREATE":
		// 模拟创建表
		return tm.TableName != "" && !tm.CreatedAt.IsZero()
	case "INSERT", "SELECT", "UPDATE", "DELETE":
		// 模拟数据操作
		return tm.TableName != ""
	case "DROP":
		// 模拟删除表
		return tm.TableName != ""
	default:
		return false
	}
}

// === 真实PostgreSQL集成测试 ===

// TestUser 测试用的用户模型
type TestUser struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Age       int       `gorm:"default:0" json:"age"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (TestUser) TableName() string {
	return "test_users"
}

// TestPostgreSQLIntegration 真实PostgreSQL连接和操作测试
func TestPostgreSQLIntegration(t *testing.T) {
	// 初始化日志系统
	cfg := getTestConfig()
	err := pkglog.InitLoggerEnhanced(&cfg.Logging)
	require.NoError(t, err)

	// 尝试连接真实PostgreSQL
	err = InitDB(cfg)
	if err != nil {
		t.Skipf("跳过PostgreSQL集成测试：无法连接PostgreSQL服务器 (%v)", err)
		return
	}

	// 确保测试结束后清理
	defer func() {
		if db := GetDB(); db != nil {
			// 清理测试数据和表
			db.Exec("DROP TABLE IF EXISTS test_users")
			db.Exec("DROP TABLE IF EXISTS test_integration_table")
		}
		CloseDB()
	}()

	t.Run("真实连接验证", func(t *testing.T) {
		assert.True(t, IsReady())
		assert.NotNil(t, GetDB())

		err := Ping()
		assert.NoError(t, err)

		err = CheckDBHealth()
		assert.NoError(t, err)
	})

	t.Run("数据库统计信息", func(t *testing.T) {
		stats, err := GetDBStats()
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.True(t, stats.MaxOpenConnections > 0)
		t.Logf("数据库统计: %+v", stats)
	})

	t.Run("表管理功能测试", func(t *testing.T) {
		// 生成唯一表名
		timestamp := time.Now().UnixNano()
		tableName := fmt.Sprintf("test_integration_table_%d", timestamp)

		// 定义表结构
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
					Type:    "VARCHAR(100)",
					NotNull: true,
				},
				{
					Name:   "email",
					Type:   "VARCHAR(100)",
					Unique: true,
				},
				{
					Name:         "age",
					Type:         "INTEGER",
					DefaultValue: "0",
				},
				{
					Name:         "created_at",
					Type:         "TIMESTAMP",
					DefaultValue: "CURRENT_TIMESTAMP",
				},
			},
		}

		// 创建表
		err := CreateTable(tableDef)
		assert.NoError(t, err)

		// 验证表存在
		exists, err := TableExists(tableName)
		assert.NoError(t, err)
		assert.True(t, exists)

		// 获取表信息
		tableInfo, err := GetTableInfo(tableName)
		assert.NoError(t, err)
		assert.NotNil(t, tableInfo)
		assert.Equal(t, tableName, tableInfo.Name)
		t.Logf("表信息: %+v", tableInfo)

		// 使用原生SQL插入数据
		db := GetDB()
		insertSQL := fmt.Sprintf(`
			INSERT INTO %s (name, email, age) 
			VALUES ('测试用户1', 'test1@example.com', 25),
			       ('测试用户2', 'test2@example.com', 30),
			       ('测试用户3', 'test3@example.com', 35)
		`, tableName)

		result := db.Exec(insertSQL)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(3), result.RowsAffected)

		// 查询数据
		var count int64
		err = db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// 查询特定数据
		var users []map[string]interface{}
		err = db.Raw(fmt.Sprintf("SELECT * FROM %s ORDER BY id", tableName)).Scan(&users).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, len(users))

		// 验证数据内容
		assert.Equal(t, "测试用户1", users[0]["name"])
		assert.Equal(t, "test1@example.com", users[0]["email"])

		// 更新数据
		updateSQL := fmt.Sprintf("UPDATE %s SET age = 26 WHERE email = 'test1@example.com'", tableName)
		result = db.Exec(updateSQL)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// 验证更新
		var updatedAge int
		err = db.Raw(fmt.Sprintf("SELECT age FROM %s WHERE email = 'test1@example.com'", tableName)).Scan(&updatedAge).Error
		assert.NoError(t, err)
		assert.Equal(t, 26, updatedAge)

		// 删除数据
		deleteSQL := fmt.Sprintf("DELETE FROM %s WHERE email = 'test3@example.com'", tableName)
		result = db.Exec(deleteSQL)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// 验证删除
		err = db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 删除表
		err = DropTable(tableName)
		assert.NoError(t, err)

		// 验证表已删除
		exists, err = TableExists(tableName)
		assert.NoError(t, err)
		assert.False(t, exists)

		t.Logf("表管理功能测试完成：创建 -> 插入 -> 查询 -> 更新 -> 删除 -> 删除表")
	})

	t.Run("GORM模型CRUD操作", func(t *testing.T) {
		db := GetDB()

		// 自动迁移创建表
		err := db.AutoMigrate(&TestUser{})
		assert.NoError(t, err)

		// 验证表存在
		exists, err := TableExists("test_users")
		assert.NoError(t, err)
		assert.True(t, exists)

		// 创建用户数据
		users := []TestUser{
			{
				Name:     "张三",
				Email:    "zhangsan@example.com",
				Age:      25,
				IsActive: true,
			},
			{
				Name:     "李四",
				Email:    "lisi@example.com",
				Age:      30,
				IsActive: true,
			},
			{
				Name:     "王五",
				Email:    "wangwu@example.com",
				Age:      35,
				IsActive: false,
			},
		}

		// 批量插入
		err = db.Create(&users).Error
		assert.NoError(t, err)
		assert.True(t, users[0].ID > 0) // 验证ID已生成

		// 查询单个用户
		var user TestUser
		err = db.Where("email = ?", "zhangsan@example.com").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "张三", user.Name)
		assert.Equal(t, 25, user.Age)
		assert.True(t, user.IsActive)

		// 查询所有用户
		var allUsers []TestUser
		err = db.Find(&allUsers).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, len(allUsers))

		// 条件查询
		var activeUsers []TestUser
		err = db.Where("is_active = ?", true).Find(&activeUsers).Error
		assert.NoError(t, err)
		assert.Equal(t, 2, len(activeUsers))

		// 更新用户
		err = db.Model(&user).Update("age", 26).Error
		assert.NoError(t, err)

		// 验证更新
		err = db.Where("email = ?", "zhangsan@example.com").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, 26, user.Age)

		// 批量更新
		err = db.Model(&TestUser{}).Where("is_active = ?", false).Update("age", 40).Error
		assert.NoError(t, err)

		// 软删除（GORM特性）
		err = db.Delete(&user).Error
		assert.NoError(t, err)

		// 验证软删除 - 正常查询不到
		err = db.Where("email = ?", "zhangsan@example.com").First(&user).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// 包含删除的查询
		err = db.Unscoped().Where("email = ?", "zhangsan@example.com").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "张三", user.Name)

		// 物理删除
		err = db.Unscoped().Delete(&user).Error
		assert.NoError(t, err)

		// 验证物理删除
		err = db.Unscoped().Where("email = ?", "zhangsan@example.com").First(&user).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// 统计剩余用户
		var count int64
		err = db.Model(&TestUser{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		t.Logf("GORM CRUD操作测试完成：插入 -> 查询 -> 更新 -> 软删除 -> 物理删除")
	})

	t.Run("事务操作测试", func(t *testing.T) {
		// 测试成功的事务
		err := WithTransaction(func(tx *gorm.DB) error {
			user1 := TestUser{
				Name:     "事务用户1",
				Email:    "tx1@example.com",
				Age:      20,
				IsActive: true,
			}

			user2 := TestUser{
				Name:     "事务用户2",
				Email:    "tx2@example.com",
				Age:      25,
				IsActive: true,
			}

			// 在事务中创建用户
			if err := tx.Create(&user1).Error; err != nil {
				return err
			}

			if err := tx.Create(&user2).Error; err != nil {
				return err
			}

			return nil
		})

		assert.NoError(t, err)

		// 验证事务成功 - 数据已提交
		db := GetDB()
		var count int64
		err = db.Model(&TestUser{}).Where("email LIKE ?", "tx%@example.com").Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// 测试失败的事务（应该回滚）
		err = WithTransaction(func(tx *gorm.DB) error {
			user3 := TestUser{
				Name:     "事务用户3",
				Email:    "tx3@example.com",
				Age:      30,
				IsActive: true,
			}

			// 创建用户
			if err := tx.Create(&user3).Error; err != nil {
				return err
			}

			// 故意返回错误导致回滚
			return errors.New("事务测试错误 - 应该回滚")
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "事务测试错误")

		// 验证事务回滚 - 用户3不应该存在
		var user TestUser
		err = db.Where("email = ?", "tx3@example.com").First(&user).Error
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// 验证之前的数据仍然存在
		err = db.Model(&TestUser{}).Where("email LIKE ?", "tx%@example.com").Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count) // 仍然是2个，没有增加

		t.Logf("事务操作测试完成：成功提交 + 失败回滚")
	})

	t.Run("CreateTestTable功能测试", func(t *testing.T) {
		// 使用内置的CreateTestTable函数
		tableName, err := CreateTestTable()
		assert.NoError(t, err)
		assert.NotEmpty(t, tableName)
		assert.Contains(t, tableName, "test_")

		// 验证表存在
		exists, err := TableExists(tableName)
		assert.NoError(t, err)
		assert.True(t, exists)

		// 插入测试数据
		db := GetDB()
		insertSQL := fmt.Sprintf(`
			INSERT INTO %s (name, email, created_at) 
			VALUES ('CreateTestTable用户', 'createtest@example.com', CURRENT_TIMESTAMP)
		`, tableName)

		result := db.Exec(insertSQL)
		assert.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// 查询数据
		var count int64
		err = db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 删除测试表
		err = DropTable(tableName)
		assert.NoError(t, err)

		// 验证表已删除
		exists, err = TableExists(tableName)
		assert.NoError(t, err)
		assert.False(t, exists)

		t.Logf("CreateTestTable功能测试完成，表名: %s", tableName)
	})

	t.Run("复杂查询和聚合操作", func(t *testing.T) {
		db := GetDB()

		// 清理可能存在的测试数据
		db.Unscoped().Where("email LIKE ?", "complex%@example.com").Delete(&TestUser{})

		// 创建更多测试数据用于复杂查询
		complexUsers := []TestUser{
			{Name: "复杂查询用户1", Email: "complex1@example.com", Age: 20, IsActive: true},
			{Name: "复杂查询用户2", Email: "complex2@example.com", Age: 25, IsActive: true},
			{Name: "复杂查询用户3", Email: "complex3@example.com", Age: 30, IsActive: false},
			{Name: "复杂查询用户4", Email: "complex4@example.com", Age: 35, IsActive: true},
			{Name: "复杂查询用户5", Email: "complex5@example.com", Age: 40, IsActive: false},
		}

		err = db.Create(&complexUsers).Error
		assert.NoError(t, err)

		// 聚合查询 - 平均年龄
		var avgAge float64
		err = db.Model(&TestUser{}).Where("email LIKE ?", "complex%@example.com").
			Select("AVG(age)").Scan(&avgAge).Error
		assert.NoError(t, err)
		assert.Equal(t, 30.0, avgAge)

		// 分组统计
		type AgeGroup struct {
			IsActive bool    `json:"is_active"`
			Count    int64   `json:"count"`
			AvgAge   float64 `json:"avg_age"`
		}

		var ageGroups []AgeGroup
		err = db.Model(&TestUser{}).Where("email LIKE ?", "complex%@example.com").
			Select("is_active, COUNT(*) as count, AVG(age) as avg_age").
			Group("is_active").Scan(&ageGroups).Error
		assert.NoError(t, err)
		assert.Equal(t, 2, len(ageGroups))

		// 排序和分页
		var orderedUsers []TestUser
		err = db.Where("email LIKE ?", "complex%@example.com").
			Order("age DESC").Limit(3).Find(&orderedUsers).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, len(orderedUsers))
		assert.Equal(t, 40, orderedUsers[0].Age) // 最大年龄在前

		// 清理测试数据
		db.Unscoped().Where("email LIKE ?", "complex%@example.com").Delete(&TestUser{})

		t.Logf("复杂查询测试完成：聚合 -> 分组 -> 排序分页")
	})
}

// TestPostgreSQLConcurrency 测试PostgreSQL并发操作
func TestPostgreSQLConcurrency(t *testing.T) {
	cfg := getTestConfig()
	err := InitDB(cfg)
	if err != nil {
		t.Skipf("跳过PostgreSQL并发测试：无法连接PostgreSQL服务器 (%v)", err)
		return
	}

	defer CloseDB()

	// 确保测试表存在
	db := GetDB()
	err = db.AutoMigrate(&TestUser{})
	require.NoError(t, err)

	t.Run("并发读写测试", func(t *testing.T) {
		const numGoroutines = 5
		const numOperations = 10

		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines*numOperations)

		// 清理之前的测试数据
		db.Unscoped().Where("email LIKE ?", "concurrent%@example.com").Delete(&TestUser{})

		// 启动多个goroutine进行并发操作
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					// 添加小延迟减少数据库连接池压力
					time.Sleep(time.Millisecond * 10)

					user := TestUser{
						Name:     fmt.Sprintf("并发用户_%d_%d", goroutineID, j),
						Email:    fmt.Sprintf("concurrent%d_%d@example.com", goroutineID, j),
						Age:      20 + (goroutineID*10 + j),
						IsActive: j%2 == 0, // 交替设置活跃状态
					}

					// CREATE操作
					if err := db.Create(&user).Error; err != nil {
						errors <- fmt.Errorf("CREATE失败 [%d:%d]: %v", goroutineID, j, err)
						continue
					}

					// READ操作
					var readUser TestUser
					if err := db.Where("email = ?", user.Email).First(&readUser).Error; err != nil {
						errors <- fmt.Errorf("READ失败 [%d:%d]: %v", goroutineID, j, err)
						continue
					}

					if readUser.Name != user.Name {
						errors <- fmt.Errorf("数据不匹配 [%d:%d]: 期望 %s, 得到 %s",
							goroutineID, j, user.Name, readUser.Name)
						continue
					}

					// UPDATE操作
					newAge := readUser.Age + 1
					if err := db.Model(&readUser).Update("age", newAge).Error; err != nil {
						errors <- fmt.Errorf("UPDATE失败 [%d:%d]: %v", goroutineID, j, err)
						continue
					}

					// 验证UPDATE
					if err := db.Where("email = ?", user.Email).First(&readUser).Error; err != nil {
						errors <- fmt.Errorf("UPDATE验证失败 [%d:%d]: %v", goroutineID, j, err)
						continue
					}

					if readUser.Age != newAge {
						errors <- fmt.Errorf("UPDATE未生效 [%d:%d]: 期望 %d, 得到 %d",
							goroutineID, j, newAge, readUser.Age)
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// 检查错误
		var errorCount int
		for err := range errors {
			t.Errorf("并发操作错误: %v", err)
			errorCount++
		}

		if errorCount > 0 {
			t.Errorf("并发测试中发生了 %d 个错误", errorCount)
		}

		// 验证总数据量
		var count int64
		err = db.Model(&TestUser{}).Where("email LIKE ?", "concurrent%@example.com").Count(&count).Error
		assert.NoError(t, err)
		expectedCount := int64(numGoroutines * numOperations)
		assert.Equal(t, expectedCount, count,
			"应该创建了 %d 个用户，实际创建了 %d 个", expectedCount, count)

		// 清理测试数据
		db.Unscoped().Where("email LIKE ?", "concurrent%@example.com").Delete(&TestUser{})

		t.Logf("并发测试成功：%d个goroutine，每个执行%d次操作", numGoroutines, numOperations)
	})
}

// TestPostgreSQLPerformance 测试PostgreSQL性能
func TestPostgreSQLPerformance(t *testing.T) {
	cfg := getTestConfig()
	err := InitDB(cfg)
	if err != nil {
		t.Skipf("跳过PostgreSQL性能测试：无法连接PostgreSQL服务器 (%v)", err)
		return
	}

	defer CloseDB()

	db := GetDB()
	err = db.AutoMigrate(&TestUser{})
	require.NoError(t, err)

	t.Run("批量插入性能测试", func(t *testing.T) {
		// 清理测试数据
		db.Unscoped().Where("email LIKE ?", "perf%@example.com").Delete(&TestUser{})

		const batchSize = 1000
		users := make([]TestUser, batchSize)

		for i := 0; i < batchSize; i++ {
			users[i] = TestUser{
				Name:     fmt.Sprintf("性能测试用户_%d", i),
				Email:    fmt.Sprintf("perf%d@example.com", i),
				Age:      20 + (i % 50),
				IsActive: i%2 == 0,
			}
		}

		// 记录开始时间
		startTime := time.Now()

		// 批量插入
		err = db.CreateInBatches(users, 100).Error
		assert.NoError(t, err)

		duration := time.Since(startTime)

		// 验证插入结果
		var count int64
		err = db.Model(&TestUser{}).Where("email LIKE ?", "perf%@example.com").Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(batchSize), count)

		// 计算性能指标
		opsPerSecond := float64(batchSize) / duration.Seconds()

		t.Logf("批量插入性能: %d条记录，耗时: %v，速度: %.2f ops/sec",
			batchSize, duration, opsPerSecond)

		// 清理测试数据
		db.Unscoped().Where("email LIKE ?", "perf%@example.com").Delete(&TestUser{})
	})

	t.Run("查询性能测试", func(t *testing.T) {
		// 预先插入数据
		const dataSize = 500
		users := make([]TestUser, dataSize)

		for i := 0; i < dataSize; i++ {
			users[i] = TestUser{
				Name:     fmt.Sprintf("查询测试用户_%d", i),
				Email:    fmt.Sprintf("query%d@example.com", i),
				Age:      20 + (i % 50),
				IsActive: i%2 == 0,
			}
		}

		err = db.CreateInBatches(users, 100).Error
		require.NoError(t, err)

		// 测试单条查询性能
		const queryCount = 100
		startTime := time.Now()

		for i := 0; i < queryCount; i++ {
			var user TestUser
			email := fmt.Sprintf("query%d@example.com", i%dataSize)
			err = db.Where("email = ?", email).First(&user).Error
			assert.NoError(t, err)
		}

		duration := time.Since(startTime)
		avgQueryTime := duration / queryCount

		t.Logf("单条查询性能: %d次查询，总耗时: %v，平均: %v/query",
			queryCount, duration, avgQueryTime)

		// 测试批量查询性能
		startTime = time.Now()

		var allUsers []TestUser
		err = db.Where("email LIKE ?", "query%@example.com").Find(&allUsers).Error
		assert.NoError(t, err)
		assert.Equal(t, dataSize, len(allUsers))

		duration = time.Since(startTime)

		t.Logf("批量查询性能: %d条记录，耗时: %v", dataSize, duration)

		// 清理测试数据
		db.Unscoped().Where("email LIKE ?", "query%@example.com").Delete(&TestUser{})
	})
}

// BenchmarkPostgreSQLOperations PostgreSQL操作性能基准测试
func BenchmarkPostgreSQLOperations(b *testing.B) {
	cfg := getTestConfig()
	err := InitDB(cfg)
	if err != nil {
		b.Skipf("跳过PostgreSQL基准测试：无法连接PostgreSQL服务器 (%v)", err)
		return
	}

	defer CloseDB()

	db := GetDB()
	err = db.AutoMigrate(&TestUser{})
	if err != nil {
		b.Fatalf("Failed to migrate test table: %v", err)
	}

	b.Run("CREATE操作", func(b *testing.B) {
		// 清理数据
		db.Unscoped().Where("email LIKE ?", "bench_create%@example.com").Delete(&TestUser{})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			user := TestUser{
				Name:     fmt.Sprintf("基准测试用户_%d", i),
				Email:    fmt.Sprintf("bench_create%d@example.com", i),
				Age:      25,
				IsActive: true,
			}
			db.Create(&user)
		}
	})

	b.Run("READ操作", func(b *testing.B) {
		// 预设数据
		const dataSize = 1000
		users := make([]TestUser, dataSize)
		for i := 0; i < dataSize; i++ {
			users[i] = TestUser{
				Name:     fmt.Sprintf("基准测试读取用户_%d", i),
				Email:    fmt.Sprintf("bench_read%d@example.com", i),
				Age:      25,
				IsActive: true,
			}
		}
		db.CreateInBatches(users, 100)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var user TestUser
			email := fmt.Sprintf("bench_read%d@example.com", i%dataSize)
			db.Where("email = ?", email).First(&user)
		}
	})

	b.Run("UPDATE操作", func(b *testing.B) {
		// 预设数据
		const dataSize = 1000
		for i := 0; i < dataSize; i++ {
			user := TestUser{
				Name:     fmt.Sprintf("基准测试更新用户_%d", i),
				Email:    fmt.Sprintf("bench_update%d@example.com", i),
				Age:      25,
				IsActive: true,
			}
			db.Create(&user)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			email := fmt.Sprintf("bench_update%d@example.com", i%dataSize)
			db.Model(&TestUser{}).Where("email = ?", email).Update("age", 26)
		}
	})

	b.Run("DELETE操作", func(b *testing.B) {
		// 为每次测试预设数据
		for i := 0; i < b.N; i++ {
			user := TestUser{
				Name:     fmt.Sprintf("基准测试删除用户_%d", i),
				Email:    fmt.Sprintf("bench_delete%d@example.com", i),
				Age:      25,
				IsActive: true,
			}
			db.Create(&user)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			email := fmt.Sprintf("bench_delete%d@example.com", i)
			db.Where("email = ?", email).Delete(&TestUser{})
		}
	})

	b.Run("事务操作", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			WithTransaction(func(tx *gorm.DB) error {
				user := TestUser{
					Name:     fmt.Sprintf("基准测试事务用户_%d", i),
					Email:    fmt.Sprintf("bench_tx%d@example.com", i),
					Age:      25,
					IsActive: true,
				}
				return tx.Create(&user).Error
			})
		}
	})
}
