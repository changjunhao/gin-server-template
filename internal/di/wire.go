//go:build wireinject
// +build wireinject

package di

import (
	"fmt"
	"gin-server/internal/config"
	"gin-server/internal/controller"
	"gin-server/internal/database"
	"gin-server/internal/repository"
	"gin-server/internal/repository/mongodb"
	"gin-server/internal/repository/mysql"
	"gin-server/internal/service"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitializeContainer 使用Wire注入依赖
func InitializeContainer(cfg *config.Config) (*Container, error) {
	wire.Build(
		wire.Struct(new(Container), "*"),
		provideDB,
		provideUserRepository,
		service.NewUserService,
		provideUserController,
	)
	return nil, nil
}

// provideDB 提供数据库连接
func provideDB(cfg *config.Config) (*gorm.DB, error) {
	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		return nil, err
	}

	// 获取数据库连接
	return database.GetDB(), nil
}

// provideUserRepository 根据配置提供用户仓库实现
func provideUserRepository(cfg *config.Config, db *gorm.DB) (repository.UserRepository, error) {
	// 根据配置选择仓库实现
	switch cfg.Database.Driver {
	case "mysql":
		return mysql.NewUserRepository(), nil
	case "mongodb":
		return mongodb.NewUserRepository(), nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Driver)
	}
}

// provideUserController 提供用户控制器
func provideUserController(userService *service.UserService, cfg *config.Config) *controller.UserController {
	return controller.NewUserController(userService, &cfg.JWT)
}