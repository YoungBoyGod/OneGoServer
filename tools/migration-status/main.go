package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
	pkgsql "github.com/YoungBoyGod/OneGoServer/pkg/sql"
	"github.com/YoungBoyGod/OneGoServer/pkg/utils"
)

func main() {
	fmt.Println("=== 数据库迁移状态查看工具 ===")

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

	// 获取迁移状态
	migrations, err := utils.GetMigrationStatus(db)
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Println("🔍 没有找到已执行的迁移记录")
		fmt.Println("   这可能意味着:")
		fmt.Println("   - 还没有执行过任何迁移")
		fmt.Println("   - schema_migrations表还未创建")
		return
	}

	fmt.Printf("\n📊 已执行的迁移列表 (共 %d 个):\n", len(migrations))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("%-10s %-30s %-20s %-10s %-15s\n", "版本号", "迁移名称", "执行时间", "状态", "执行耗时(ms)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, migration := range migrations {
		status := "✅ 成功"
		if !migration.IsSuccess {
			status = "❌ 失败"
		}

		appliedTime := migration.AppliedAt.Format("2006-01-02 15:04:05")
		execTime := getExecutionTime(db, migration.Version)

		fmt.Printf("%-10s %-30s %-20s %-10s %-15s\n",
			migration.Version,
			migration.Name,
			appliedTime,
			status,
			fmt.Sprintf("%d", execTime))
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 显示详细信息
	fmt.Println("\n📋 详细信息:")
	for _, migration := range migrations {
		fmt.Printf("\n🔹 版本 %s (%s):\n", migration.Version, migration.Name)
		fmt.Printf("   执行时间: %s\n", migration.AppliedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   状态: %s\n", getStatusIcon(migration.IsSuccess))
		fmt.Printf("   执行耗时: %d ms\n", getExecutionTime(db, migration.Version))
		fmt.Printf("   校验和: %s\n", getChecksum(db, migration.Version))

		if !migration.IsSuccess {
			errorMsg := getErrorMessage(db, migration.Version)
			if errorMsg != "" {
				fmt.Printf("   ❌ 错误信息: %s\n", errorMsg)
			}
		}
	}

	fmt.Println("\n✨ 查看完成!")
}

func getStatusIcon(isSuccess bool) string {
	if isSuccess {
		return "✅ 成功执行"
	}
	return "❌ 执行失败"
}

func getExecutionTime(db *sql.DB, version string) int64 {
	var execTime sql.NullInt64
	err := db.QueryRow("SELECT execution_time FROM schema_migrations WHERE version = $1", version).Scan(&execTime)
	if err != nil {
		return 0
	}
	if execTime.Valid {
		return execTime.Int64
	}
	return 0
}

func getErrorMessage(db *sql.DB, version string) string {
	var errorMsg sql.NullString
	err := db.QueryRow("SELECT error_message FROM schema_migrations WHERE version = $1", version).Scan(&errorMsg)
	if err != nil {
		return ""
	}
	if errorMsg.Valid {
		return errorMsg.String
	}
	return ""
}

func getChecksum(db *sql.DB, version string) string {
	var checksum sql.NullString
	err := db.QueryRow("SELECT checksum FROM schema_migrations WHERE version = $1", version).Scan(&checksum)
	if err != nil {
		return ""
	}
	if checksum.Valid {
		return checksum.String[:8] + "..." // 只显示前8位
	}
	return ""
}
