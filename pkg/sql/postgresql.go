package sql

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义全局数据库实例
var DB *gorm.DB

// 1.数据库连接
func InitDB(cfg *config.Config) {
	// 从dsn中获取
	dsn := cfg.Database.GetDsn()
	fmt.Println(dsn)
	// 配置gorm
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	if lifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime); err == nil {
		sqlDB.SetConnMaxLifetime(lifetime)
	}
	// 设置全局数据库实例
	DB = db
	// 健康检查
	if err := CheckDBHealth(); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}
	log.Printf("Database connected successfully")
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// 关闭数据库
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// 检查数据库健康状态
func CheckDBHealth() error {
	if DB == nil {

		return errors.New("database not initialized")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	return nil
}
