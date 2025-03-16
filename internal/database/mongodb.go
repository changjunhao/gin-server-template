package database

import (
	"context"
	"fmt"
	"gin-server-template/internal/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB 客户端实例
var MongoDB *mongo.Client
var MongoDBName string

// InitMongoDB 初始化MongoDB连接
func InitMongoDB(cfg *config.DatabaseConfig) error {
	// 构建MongoDB连接URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(uri)

	// 设置连接池
	clientOptions.SetMaxPoolSize(uint64(cfg.MaxOpenConns))
	clientOptions.SetMinPoolSize(uint64(cfg.MaxIdleConns))
	clientOptions.SetMaxConnIdleTime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 创建连接上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 验证连接
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	// 保存客户端实例和数据库名称
	MongoDB = client
	MongoDBName = cfg.DBName

	log.Println("MongoDB连接成功")
	return nil
}

// GetMongoDB 获取MongoDB客户端实例
func GetMongoDB() *mongo.Client {
	return MongoDB
}

// GetMongoDBName 获取MongoDB数据库名称
func GetMongoDBName() string {
	return MongoDBName
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() error {
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoDB.Disconnect(ctx)
	}
	return nil
}
