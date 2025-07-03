package utils

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration 结构体用于存储迁移信息
type Migration struct {
	Version     string
	Name        string
	SQL         string
	Checksum    string
	AppliedAt   time.Time
	IsSuccess   bool
	RollbackSQL string
}

// RunMigrations 执行数据库迁移
func RunMigrations(db *sql.DB) error {
	// 1. 创建migrations表(如果不存在)
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("创建migrations表失败: %v", err)
	}

	// 2. 获取已执行的migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("获取已执行migrations失败: %v", err)
	}

	// 3. 读取migrations目录下的所有SQL文件
	migrations, err := loadMigrationFiles("internal/data/migrations")
	if err != nil {
		return fmt.Errorf("读取migration文件失败: %v", err)
	}

	// 4. 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// 5. 执行未应用的migrations
	for _, migration := range migrations {
		// 检查是否已执行
		if _, exists := appliedMigrations[migration.Version]; exists {
			continue
		}

		// 开始执行migration
		startTime := time.Now()

		// 开启事务
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("开启事务失败: %v", err)
		}

		// 执行SQL
		_, err = tx.Exec(migration.SQL)
		if err != nil {
			tx.Rollback()
			// 记录失败信息
			if err := recordMigration(db, migration, startTime, false, err.Error()); err != nil {
				return fmt.Errorf("记录失败migration失败: %v", err)
			}
			return fmt.Errorf("执行migration %s 失败: %v", migration.Version, err)
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("提交事务失败: %v", err)
		}

		// 记录成功信息
		if err := recordMigration(db, migration, startTime, true, ""); err != nil {
			return fmt.Errorf("记录成功migration失败: %v", err)
		}
	}

	return nil
}

// createMigrationsTable 创建migrations表
func createMigrationsTable(db *sql.DB) error {
	sql, err := os.ReadFile("internal/data/migrations/000_create_migrations_table.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(sql))
	return err
}

// getAppliedMigrations 获取已执行的migrations
func getAppliedMigrations(db *sql.DB) (map[string]Migration, error) {
	migrations := make(map[string]Migration)

	rows, err := db.Query("SELECT version, name, applied_at, is_success FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Migration
		err := rows.Scan(&m.Version, &m.Name, &m.AppliedAt, &m.IsSuccess)
		if err != nil {
			return nil, err
		}
		migrations[m.Version] = m
	}

	return migrations, rows.Err()
}

// loadMigrationFiles 读取migration文件
func loadMigrationFiles(dir string) ([]Migration, error) {
	var migrations []Migration

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// 跳过migrations表创建文件
		if file.Name() == "000_create_migrations_table.sql" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		// 解析文件名
		version := strings.Split(file.Name(), "_")[0]
		name := strings.TrimSuffix(strings.Join(strings.Split(file.Name(), "_")[1:], "_"), ".sql")

		// 计算checksum
		hash := md5.Sum(content)
		checksum := hex.EncodeToString(hash[:])

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			SQL:      string(content),
			Checksum: checksum,
		})
	}

	return migrations, nil
}

// recordMigration 记录migration执行结果
func recordMigration(db *sql.DB, m Migration, startTime time.Time, success bool, errorMsg string) error {
	execTime := time.Since(startTime).Milliseconds()

	_, err := db.Exec(`
		INSERT INTO schema_migrations (
			version, name, checksum, execution_time, is_success, error_message
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		m.Version, m.Name, m.Checksum, execTime, success, errorMsg)

	return err
}

// RollbackMigration 回滚指定版本的migration
func RollbackMigration(db *sql.DB, version string) error {
	// 获取migration记录
	var rollbackSQL string
	err := db.QueryRow("SELECT rollback_sql FROM schema_migrations WHERE version = $1", version).Scan(&rollbackSQL)
	if err != nil {
		return fmt.Errorf("获取migration记录失败: %v", err)
	}

	if rollbackSQL == "" {
		return fmt.Errorf("migration %s 没有回滚SQL", version)
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	// 执行回滚SQL
	_, err = tx.Exec(rollbackSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("执行回滚SQL失败: %v", err)
	}

	// 删除migration记录
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除migration记录失败: %v", err)
	}

	return tx.Commit()
}

// GetMigrationStatus 获取所有migrations的状态
func GetMigrationStatus(db *sql.DB) ([]Migration, error) {
	rows, err := db.Query(`
		SELECT version, name, applied_at, is_success, error_message 
		FROM schema_migrations 
		ORDER BY version`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		var errMsg sql.NullString
		err := rows.Scan(&m.Version, &m.Name, &m.AppliedAt, &m.IsSuccess, &errMsg)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}

	return migrations, rows.Err()
}
