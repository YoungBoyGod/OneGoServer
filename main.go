package main

import (
	"fmt"
	"log"

	"github.com/YoungBoyGod/OneGoServer/internal/config"
)

func main() {
	fmt.Println("Hello, World!")
	// 加载配置文件
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Server started on port %d", cfg.Server.Port)
	// 打印配置文件内容
	log.Printf("Database Config: %+v", cfg.Database)
	log.Printf("Server Config: %+v", cfg.Server)

}
