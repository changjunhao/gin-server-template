package app

import (
	"gin-server-template/internal/controller"
	"gin-server-template/internal/middleware"
)

// setupRoutes 配置所有API路由
func (s *Server) setupRoutes() {
	// 创建控制器实例
	userController := controller.NewUserController()

	// 公共路由组
	public := s.router.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := public.Group("/users")
		{
			userGroup.POST("/register", userController.Register)
			userGroup.POST("/login", userController.Login)
		}
	}

	// 需要认证的路由组
	authorized := s.router.Group("/api/v1")
	authorized.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		userGroup := authorized.Group("/users")
		{
			userGroup.GET("/profile", userController.GetProfile)
			userGroup.PUT("/profile", userController.UpdateProfile)
		}
	}
}
