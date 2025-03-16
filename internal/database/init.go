package database

import (
	"errors"
	"gin-server-template/internal/config"
	"gin-server-template/internal/entity"
	"log"
)

// InitDatabase 初始化数据库
func InitDatabase(cfg *config.DatabaseConfig) error {
	// 根据配置选择数据库类型
	switch cfg.Driver {
	case "mysql":
		// 初始化MySQL连接
		err := InitMySQL(cfg)
		if err != nil {
			return err
		}

		// 自动迁移数据库模型
		err = AutoMigrate()
		if err != nil {
			return err
		}

	case "mongodb":
		// 初始化MongoDB连接
		err := InitMongoDB(cfg)
		if err != nil {
			return err
		}

	default:
		return errors.New("不支持的数据库类型: " + cfg.Driver)
	}

	log.Println("数据库初始化成功")
	return nil
}

// AutoMigrate 自动迁移数据库模型（仅MySQL使用）
func AutoMigrate() error {
	// 在这里添加需要迁移的模型
	return DB.AutoMigrate(
		&entity.User{},
		// 其他模型...
	)
}
