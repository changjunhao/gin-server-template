package database

import (
	"fmt"
	"gin-server/internal/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	// 配置GORM日志
	logConfig := logger.Config{
		SlowThreshold: time.Second, // 慢SQL阈值
		LogLevel:      logger.Info, // 日志级别
		Colorful:      true,        // 彩色打印
	}

	// 初始化连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用标准日志库
			logConfig,
		),
	})

	if err != nil {
		return err
	}

	// 获取通用数据库对象，设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour)           // 连接最大生命周期

	DB = db
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// CloseMySQL 关闭MySQL连接
func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}