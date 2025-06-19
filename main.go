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
		log.Printf("Server starting on %s:%s", config.Server.Host, config.Server.Port) // 打印服务器启动信息
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {    // 启动服务器，忽略http.ErrServerClosed（表示正常关闭）
			log.Fatalf("Server startup failed: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，超时时间为5秒
	quit := make(chan os.Signal, 1)                      // 创建一个通道，用于接收中断信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 监听中断信号
	sig := <-quit                                        // 阻塞，直到收到中断信号, 收到中断信号后，会打印服务器关闭信息
	log.Printf("Server is shutting down...:%v", sig)     // 打印服务器关闭信息
	// 优雅关闭的上下文，超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 创建一个上下文，用于优雅关闭服务器
	defer cancel()                                                          // 确保释放资源

	if err := srv.Shutdown(ctx); err != nil { // 优雅关闭服务器
		log.Fatal("Server forced to shutdown:", err) // 如果优雅关闭失败，则强制关闭服务器
	}
	log.Println("Server exited")                                 // 打印服务器退出信息
	log.Println("Server time:", time.Now().Format(time.RFC3339)) // 打印服务器时间
}
