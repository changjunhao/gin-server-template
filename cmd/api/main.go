package main

import (
	"gin-server-template/internal/app"
	"gin-server-template/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化服务器
	server := app.NewServer(cfg)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在单独的goroutine中启动服务器
	go func() {
		if err := server.Run(); err != nil {
			log.Printf("服务器停止运行: %v", err)
			sigChan <- syscall.SIGTERM // 如果服务器异常停止，发送信号以触发清理
		}
	}()

	// 等待终止信号
	sig := <-sigChan
	log.Printf("接收到信号: %v，准备关闭服务器", sig)

	// 关闭服务器并释放资源
	server.Close()
	log.Println("服务器已安全关闭")
}
