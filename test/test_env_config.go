package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
)

func main() {
	fmt.Println("=== 测试.env配置加载 ===")

	// 创建临时.env文件
	envContent := `DB_HOST=test.example.com
DB_PORT=9999
DB_USERNAME=testuser
DB_PASSWORD=testpass123
DB_NAME=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		log.Fatalf("创建.env文件失败: %v", err)
	}
	defer os.Remove(".env") // 测试完成后删除

	fmt.Println("临时.env文件内容:")
	fmt.Println(envContent)
	fmt.Println()

	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	fmt.Println("=== 配置加载结果 ===")
	fmt.Printf("数据库主机: %s (期望: test.example.com)\n", cfg.Database.Host)
	fmt.Printf("数据库端口: %d (期望: 9999)\n", cfg.Database.Port)
	fmt.Printf("数据库用户: %s (期望: testuser)\n", cfg.Database.Username)
	fmt.Printf("数据库密码: %s (期望: testpass123)\n", cfg.Database.Password)
	fmt.Printf("数据库名称: %s (期望: testdb)\n", cfg.Database.DBName)
	fmt.Println()

	// 验证是否从.env文件读取
	if cfg.Database.Host == "test.example.com" &&
		cfg.Database.Port == 9999 &&
		cfg.Database.Username == "testuser" {
		fmt.Println("✅ 成功：配置正确从.env文件读取!")
	} else {
		fmt.Println("❌ 失败：配置未从.env文件读取，可能来自config.yaml")
	}
}
