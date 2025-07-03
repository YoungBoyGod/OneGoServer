package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkgsql "github.com/YoungBoyGod/OneGoServer/pkg/sql"
)

func main() {
	fmt.Println("=== 数据库迁移重置工具 ===")
	fmt.Println("⚠️  警告：此工具将删除所有迁移记录，请谨慎使用！")

	// 加载配置文件
	cfg, err := config.LoadConfig("../../config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := pkgsql.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 获取数据库连接
	gormDB := pkgsql.GetDB()
	if gormDB == nil {
		log.Fatalf("Failed to get database connection: database not initialized")
	}

	db, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// 显示当前迁移记录
	fmt.Println("\n📋 当前迁移记录:")
	showMigrations(db)

	// 清理迁移记录
	fmt.Println("\n🧹 开始清理迁移记录...")

	// 删除所有迁移记录
	result, err := db.Exec("DELETE FROM schema_migrations")
	if err != nil {
		log.Fatalf("Failed to delete migration records: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("✅ 成功删除 %d 条迁移记录\n", rowsAffected)

	// 删除可能创建的表和函数
	fmt.Println("\n🗑️  清理可能存在的表和函数...")

	cleanupSQL := []string{
		"DROP TABLE IF EXISTS task_executions CASCADE",
		"DROP TABLE IF EXISTS device_queues CASCADE",
		"DROP TABLE IF EXISTS tasks CASCADE",
		"DROP TABLE IF EXISTS device_heartbeats CASCADE",
		"DROP TABLE IF EXISTS device_logs CASCADE",
		"DROP TABLE IF EXISTS device_commands CASCADE",
		"DROP TABLE IF EXISTS device_tasks CASCADE",
		"DROP TABLE IF EXISTS devices CASCADE",
		"DROP TABLE IF EXISTS task_assignment_queue CASCADE",
		"DROP TABLE IF EXISTS device_load_monitor CASCADE",
		"DROP TABLE IF EXISTS task_assignment_history CASCADE",
		"DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE",
		"DROP FUNCTION IF EXISTS update_queue_position() CASCADE",
		"DROP FUNCTION IF EXISTS calculate_device_load_score() CASCADE",
	}

	for _, sql := range cleanupSQL {
		_, err := db.Exec(sql)
		if err != nil {
			fmt.Printf("⚠️  Warning: %s - %v\n", sql, err)
		} else {
			fmt.Printf("✅ Executed: %s\n", sql)
		}
	}

	fmt.Println("\n🎉 数据库重置完成！")
	fmt.Println("💡 现在可以重新运行 go run main.go 来执行迁移")
}

func showMigrations(db *sql.DB) {
	rows, err := db.Query("SELECT version, name, is_success, error_message FROM schema_migrations ORDER BY version")
	if err != nil {
		fmt.Printf("无法查询迁移记录: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var version, name string
		var isSuccess bool
		var errorMsg sql.NullString

		err := rows.Scan(&version, &name, &isSuccess, &errorMsg)
		if err != nil {
			continue
		}

		status := "✅ 成功"
		if !isSuccess {
			status = "❌ 失败"
		}

		fmt.Printf("  - 版本 %s (%s): %s", version, name, status)
		if !isSuccess && errorMsg.Valid {
			fmt.Printf(" - %s", errorMsg.String)
		}
		fmt.Println()
		count++
	}

	if count == 0 {
		fmt.Println("  没有找到迁移记录")
	}
}
