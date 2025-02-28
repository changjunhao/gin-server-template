package app

import (
	"fmt"
	"gin-server/internal/config"
	"gin-server/internal/database"
	"gin-server/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

// Server 表示HTTP服务器及其依赖项
type Server struct {
	config *config.Config
	router *gin.Engine
}

// NewServer 创建并配置一个新的Server实例
func NewServer(cfg *config.Config) *Server {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	if err := database.InitDatabase(&cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建Gin引擎
	router := gin.Default()

	// 设置受信任的代理
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 添加全局中间件
	router.Use(middleware.Logger())

	// 创建服务器实例
	s := &Server{
		config: cfg,
		router: router,
	}

	// 设置路由
	s.setupRoutes()

	return s
}

// Run 启动HTTP服务器
func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Server.Port))
}