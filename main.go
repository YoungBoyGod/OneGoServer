package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func LoadConfig(filename string) (*Config, error) {
	config, err := ReadConfig(filename)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	return config, nil
}

func main() {
	config, err := LoadConfig("config/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := setupRouter(config)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    config.Server.Host + ":" + config.Server.Port,
		Handler: router,
	}

	// 在goroutine中启动服务器，以便它不会阻塞
	go func() {
		log.Printf("Server starting on %s:%s", config.Server.Host, config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，超时时间为5秒
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	// 优雅关闭的上下文，超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
