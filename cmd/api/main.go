package main

import (
	"gin-server/internal/app"
	"gin-server/internal/config"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化并启动服务器
	server := app.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}