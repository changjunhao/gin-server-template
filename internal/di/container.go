package di

import (
	"gin-server/internal/config"
	"gin-server/internal/database"
	"gin-server/internal/repository"
	"gin-server/internal/repository/mongodb"
	"gin-server/internal/repository/mysql"
	"gin-server/internal/service"
	"gin-server/internal/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	Config         *config.Config
	DB             *gorm.DB
	UserRepository repository.UserRepository
	UserService    *service.UserService
	UserController *controller.UserController
}

// NewContainer 创建依赖注入容器
func NewContainer(cfg *config.Config) (*Container, error) {
	// 创建容器
	container := &Container{
		Config: cfg,
	}

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		return nil, err
	}

	// 获取数据库连接
	container.DB = database.GetDB()

	// 初始化仓库
	if err := container.initRepositories(cfg); err != nil {
		return nil, err
	}

	// 初始化服务
	container.initServices()

	// 初始化控制器
	container.initControllers()

	return container, nil
}

// 初始化仓库
func (c *Container) initRepositories(cfg *config.Config) error {
	// 根据配置选择仓库实现
	switch cfg.Database.Driver {
	case "mysql":
		c.UserRepository = mysql.NewUserRepository()
	case "mongodb":
		c.UserRepository = mongodb.NewUserRepository()
	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Driver)
	}
	return nil
}

// 初始化服务
func (c *Container) initServices() {
	c.UserService = service.NewUserService(c.UserRepository)
}

// 初始化控制器
func (c *Container) initControllers() {
	c.UserController = controller.NewUserController(c.UserService, &c.Config.JWT)
}

// SetupRouter 设置路由
func (c *Container) SetupRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(c.Config.Server.Mode)

	// 创建Gin引擎
	router := gin.Default()

	// 设置受信任的代理
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// 公共路由组
	public := router.Group("/api/v1")
	{
		// 用户相关路由
		userGroup := public.Group("/users")
		{
			userGroup.POST("/register", c.UserController.Register)
			userGroup.POST("/login", c.UserController.Login)
		}
	}

	// 需要认证的路由组
	authorized := router.Group("/api/v1")
	// authorized.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		userGroup := authorized.Group("/users")
		{
			userGroup.GET("/profile", c.UserController.GetProfile)
			userGroup.PUT("/profile", c.UserController.UpdateProfile)
		}
	}

	return router
}