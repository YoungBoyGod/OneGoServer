# PostgreSQL 数据库模块运行流程详解

## 📋 概述

本文档详细描述了 `pkg/sql/postgresql.go` 模块的完整运行流程，包括初始化、运行时操作、错误处理和生命周期管理。

## 🚀 完整运行流程

### 阶段1: 应用启动和初始化

#### 1.1 应用启动入口
```go
// main.go
func main() {
    // 1. 加载配置文件
    cfg, err := config.LoadConfig("config/config.yaml")
    
    // 2. 初始化日志系统
    err := pkglog.InitLoggerEnhanced(&cfg.Logging)
    
    // 3. 初始化数据库 ⭐ 关键步骤
    if err := sql.InitDB(cfg); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
}
```

#### 1.2 数据库初始化流程 (`InitDB`)

```go
func InitDB(cfg *config.Config) error {
    var initErr error
    
    // 🔒 使用 sync.Once 确保只初始化一次
    once.Do(func() {
        // Step 1: 配置验证
        if err := validateDBConfig(&cfg.Database); err != nil {
            initErr = fmt.Errorf("invalid database config: %w", err)
            return
        }
        
        // Step 2: 创建数据库管理器
        manager = &DBManager{
            config: &cfg.Database,
        }
        
        // Step 3: 建立数据库连接
        if err := manager.connect(); err != nil {
            initErr = fmt.Errorf("failed to connect to database: %w", err)
            return
        }
        
        // Step 4: 设置全局实例
        DB = manager.db
        
        // Step 5: 自动迁移 (可选)
        if cfg.Database.AutoMigrate {
            if err := manager.autoMigrate(); err != nil {
                pkglog.LogWarn("Auto migration failed", zap.Error(err))
            }
        }
        
        // Step 6: 记录成功日志
        pkglog.LogInfo("Database initialized successfully", 
            zap.String("host", cfg.Database.Host),
            zap.Int("port", cfg.Database.Port))
    })
    
    return initErr
}
```

### 阶段2: 配置验证详细流程

#### 2.1 配置验证函数 (`validateDBConfig`)

```go
func validateDBConfig(cfg *config.DatabaseConfig) error {
    // 验证项目清单:
    ✅ nil 检查
    ✅ 主机地址验证 (非空)
    ✅ 端口范围验证 (1-65535)
    ✅ 用户名验证 (必填)
    ✅ 数据库名验证 (必填)
    ✅ 连接池参数验证
    ✅ 负数检查
    ✅ 参数关系验证 (MaxIdleConns ≤ MaxOpenConns)
}
```

#### 2.2 验证流程时序

```
配置传入 → 基础验证 → 参数范围检查 → 逻辑关系验证 → 返回结果
    ↓           ↓           ↓              ↓            ↓
  非空检查    端口范围    连接池参数      参数一致性    成功/错误
```

### 阶段3: 数据库连接建立

#### 3.1 连接建立流程 (`manager.connect`)

```go
func (m *DBManager) connect() error {
    // Step 1: 构建 DSN 连接字符串
    dsn := m.config.GetDsn()
    // 示例: "host=localhost user=admin password=xxx dbname=onegoserver port=5432"
    
    // Step 2: 配置 GORM 参数
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
    }
    
    // Step 3: 建立数据库连接
    db, err := gorm.Open(postgres.Open(dsn), gormConfig)
    
    // Step 4: 配置连接池
    if err := m.configureConnectionPool(sqlDB); err != nil {
        return fmt.Errorf("failed to configure connection pool: %w", err)
    }
    
    // Step 5: 更新管理器状态
    m.mu.Lock()
    m.db = db
    m.isReady = true
    m.lastPing = time.Now()
    m.mu.Unlock()
    
    // Step 6: 执行初始健康检查
    if err := m.ping(); err != nil {
        return fmt.Errorf("initial health check failed: %w", err)
    }
    
    return nil
}
```

#### 3.2 连接池配置 (`configureConnectionPool`)

```go
func (m *DBManager) configureConnectionPool(sqlDB *sql.DB) error {
    // 🔧 连接池参数配置
    sqlDB.SetMaxIdleConns(m.config.MaxIdleConns)     // 最大空闲连接数
    sqlDB.SetMaxOpenConns(m.config.MaxOpenConns)     // 最大开放连接数
    
    // ⏰ 连接生命周期配置
    if m.config.ConnMaxLifetime != "" {
        lifetime, err := time.ParseDuration(m.config.ConnMaxLifetime)
        if err != nil {
            return fmt.Errorf("invalid connection max lifetime: %w", err)
        }
        sqlDB.SetConnMaxLifetime(lifetime)
    }
    
    return nil
}
```

### 阶段4: 运行时操作流程

#### 4.1 数据库实例获取 (`GetDB`)

```go
// 🔍 线程安全的实例获取
func GetDB() *gorm.DB {
    if manager != nil {
        manager.mu.RLock()          // 📖 读锁
        defer manager.mu.RUnlock()
        return manager.db
    }
    return DB  // 向后兼容
}
```

#### 4.2 健康检查流程 (`CheckDBHealth`)

```go
func CheckDBHealth() error {
    if manager != nil {
        return manager.ping()
    }
    
    // 向后兼容的健康检查
    if DB == nil {
        return errors.New("database not initialized")
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get sql.DB instance: %w", err)
    }
    
    // ⏱️ 5秒超时的 Ping 检查
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return sqlDB.PingContext(ctx)
}
```

#### 4.3 事务管理流程 (`WithTransaction`)

```go
func WithTransaction(fn TransactionFunc) error {
    // 🏁 执行流程
    startTime := time.Now()
    
    // Step 1: 开始事务
    tx := db.Begin()
    if tx.Error != nil {
        pkglog.LogDBOperation("BEGIN_TX", "postgresql", time.Since(startTime), tx.Error)
        return fmt.Errorf("failed to begin transaction: %w", tx.Error)
    }
    
    // Step 2: Panic 恢复机制
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            pkglog.LogError("Transaction panicked, rolled back",
                zap.Any("panic", r),
                zap.Duration("duration", time.Since(startTime)))
        }
    }()
    
    // Step 3: 执行业务逻辑
    if err := fn(tx); err != nil {
        tx.Rollback()
        pkglog.LogDBOperation("ROLLBACK_TX", "postgresql", time.Since(startTime), err)
        return fmt.Errorf("transaction failed: %w", err)
    }
    
    // Step 4: 提交事务
    if err := tx.Commit().Error; err != nil {
        pkglog.LogDBOperation("COMMIT_TX", "postgresql", time.Since(startTime), err)
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    // Step 5: 记录成功日志
    pkglog.LogDBOperation("COMMIT_TX", "postgresql", time.Since(startTime), nil)
    return nil
}
```

#### 4.4 统计监控流程 (`GetDBStats`)

```go
func GetDBStats() (*DBStats, error) {
    // Step 1: 获取数据库实例
    db := GetDB()
    if db == nil {
        return nil, errors.New("database not initialized")
    }
    
    // Step 2: 获取底层 sql.DB
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB instance: %w", err)
    }
    
    // Step 3: 收集统计信息
    stats := sqlDB.Stats()
    
    // Step 4: 构造返回结构
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
```

#### 4.5 迁移管理流程 (`RunMigrations`)

```go
func RunMigrations(migrations []MigrationFunc) error {
    db := GetDB()
    if db == nil {
        return errors.New("database not initialized")
    }
    
    // 🔄 批量执行迁移
    for i, migration := range migrations {
        startTime := time.Now()
        
        if err := migration(db); err != nil {
            pkglog.LogDBOperation("MIGRATION", fmt.Sprintf("step_%d", i+1), 
                time.Since(startTime), err)
            return fmt.Errorf("migration step %d failed: %w", i+1, err)
        }
        
        pkglog.LogDBOperation("MIGRATION", fmt.Sprintf("step_%d", i+1), 
            time.Since(startTime), nil)
    }
    
    pkglog.LogInfo("All migrations completed successfully", 
        zap.Int("total_migrations", len(migrations)))
    
    return nil
}
```

### 阶段5: 错误处理机制

#### 5.1 错误分类和处理

```go
// 🚨 错误类型分类
1. 配置错误    → validateDBConfig() → 返回详细错误信息
2. 连接错误    → gorm.Open() → 包装连接失败错误  
3. 池配置错误  → configureConnectionPool() → 参数配置错误
4. 健康检查错误 → ping() → 网络或数据库服务错误
5. 事务错误    → WithTransaction() → 业务逻辑或数据库错误
6. 迁移错误    → RunMigrations() → 模式变更错误
```

#### 5.2 错误包装策略

```go
// 📦 统一错误包装格式
return fmt.Errorf("操作描述: %w", originalError)

// 示例:
return fmt.Errorf("failed to connect to database: %w", err)
return fmt.Errorf("invalid database config: %w", err)
return fmt.Errorf("migration step %d failed: %w", i+1, err)
```

### 阶段6: 生命周期管理

#### 6.1 初始化生命周期

```
应用启动 → 配置加载 → 日志初始化 → 数据库初始化 → 应用就绪
    ↓         ↓         ↓           ↓             ↓
  main()   LoadConfig  InitLogger   InitDB()    业务逻辑
```

#### 6.2 运行时生命周期

```
数据库就绪 → 业务操作循环 → 定期健康检查 → 统计监控 → 应用关闭
    ↓           ↓             ↓            ↓          ↓
  IsReady()   各种DB操作    CheckDBHealth  GetDBStats  CloseDB()
```

#### 6.3 关闭流程 (`CloseDB`)

```go
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
            
            // 🧹 清理状态
            manager.db = nil
            manager.isReady = false
            
            pkglog.LogInfo("Database connection closed")
            return nil
        }
    }
    
    return nil
}
```

## 🔄 并发安全机制

### 读写锁使用

```go
// 📖 读操作 (获取实例)
manager.mu.RLock()
db := manager.db
manager.mu.RUnlock()

// ✏️ 写操作 (状态更新)
manager.mu.Lock()
manager.isReady = true
manager.lastPing = time.Now()
manager.mu.Unlock()
```

### 单例模式保证

```go
// 🔒 确保只初始化一次
var once sync.Once
once.Do(func() {
    // 初始化逻辑
})
```

## 📊 监控和日志

### 结构化日志记录

```go
// 🎯 数据库操作日志
pkglog.LogDBOperation("CONNECT", "postgresql", duration, err)
pkglog.LogDBOperation("COMMIT_TX", "postgresql", duration, nil)

// 📝 系统事件日志
pkglog.LogInfo("Database initialized successfully",
    zap.String("host", cfg.Database.Host),
    zap.Int("port", cfg.Database.Port))
```

### 性能监控指标

```go
// 📈 连接池统计
stats := &DBStats{
    MaxOpenConnections: stats.MaxOpenConnections,  // 最大连接数
    OpenConnections:    stats.OpenConnections,     // 当前连接数
    InUseConnections:   stats.InUse,               // 使用中连接
    IdleConnections:    stats.Idle,                // 空闲连接
    WaitCount:          stats.WaitCount,           // 等待次数
    WaitDuration:       stats.WaitDuration,        // 等待时长
}
```

## 🚀 使用示例流程

### 基础使用流程

```go
// 1️⃣ 应用启动
func main() {
    // 加载配置
    cfg, _ := config.LoadConfig("config/config.yaml")
    
    // 初始化数据库
    if err := sql.InitDB(cfg); err != nil {
        log.Fatal(err)
    }
    
    // 2️⃣ 获取数据库实例
    db := sql.GetDB()
    
    // 3️⃣ 执行业务操作
    var users []User
    db.Find(&users)
    
    // 4️⃣ 健康检查
    if err := sql.CheckDBHealth(); err != nil {
        log.Printf("Health check failed: %v", err)
    }
    
    // 5️⃣ 应用关闭
    defer sql.CloseDB()
}
```

### 高级功能使用

```go
// 🔄 重试初始化
err := sql.InitDBWithRetry(cfg, 3, 2*time.Second)

// 💰 事务操作
err := sql.WithTransaction(func(tx *gorm.DB) error {
    // 业务逻辑
    return tx.Create(&user).Error
})

// 📊 监控统计
stats, err := sql.GetDBStats()
if err == nil {
    log.Printf("Active: %d, Idle: %d", 
        stats.InUseConnections, stats.IdleConnections)
}

// 🔧 数据库迁移
migrations := []sql.MigrationFunc{
    func(db *gorm.DB) error {
        return db.AutoMigrate(&User{})
    },
}
err = sql.RunMigrations(migrations)
```

## 🎯 最佳实践建议

### 生产环境

1. **使用重试初始化**: `InitDBWithRetry(cfg, 3, 2*time.Second)`
2. **定期健康检查**: 每30秒调用 `CheckDBHealth()`
3. **监控连接池**: 每分钟记录 `GetDBStats()`
4. **统一事务管理**: 所有写操作使用 `WithTransaction()`

### 开发环境

1. **快速启动**: 使用 `InitDB(cfg)` 即可
2. **调试日志**: 启用详细的数据库日志
3. **频繁迁移**: 启用 `AutoMigrate: true`

### 错误处理

1. **详细日志**: 记录所有数据库操作的详细信息
2. **优雅降级**: 连接失败时提供备选方案
3. **错误包装**: 使用 `fmt.Errorf` 包装原始错误

## 📋 总结

PostgreSQL 数据库模块提供了完整的企业级数据库管理功能：

- 🔒 **线程安全**: 使用读写锁保护并发访问
- 🔄 **重试机制**: 提高连接可靠性
- 💰 **事务管理**: 自动提交/回滚和panic恢复
- 📊 **监控统计**: 实时连接池性能监控
- 🔧 **迁移管理**: 批量数据库模式变更
- ❤️‍🩹 **健康检查**: 定期检测数据库连接状态
- 📝 **结构化日志**: 详细的操作日志记录

整个流程设计遵循高可用、高性能、易维护的原则，为生产环境提供了坚实的数据层基础。 