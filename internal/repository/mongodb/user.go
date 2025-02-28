package mongodb

import (
	"context"
	"errors"
	"gin-server/internal/database"
	"gin-server/internal/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository MongoDB实现的用户仓库
type UserRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

// NewUserRepository 创建MongoDB用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{
		client:     database.GetMongoDB(),
		database:   database.GetMongoDBName(),
		collection: "users",
	}
}

// getCollection 获取用户集合
func (r *UserRepository) getCollection() *mongo.Collection {
	return r.client.Database(r.database).Collection(r.collection)
}

// Create 创建用户
func (r *UserRepository) Create(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置创建时间和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 插入文档
	result, err := r.getCollection().InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// 获取插入的ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		// 由于MongoDB使用ObjectID，我们需要将其转换为uint
		// 这里简单地使用时间戳作为ID
		user.ID = uint(oid.Timestamp().Unix())
	}

	return nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user entity.User
	err := r.getCollection().FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user entity.User
	err := r.getCollection().FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.getCollection().CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := r.getCollection().CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 更新时间
	user.UpdatedAt = time.Now()

	// 更新文档
	_, err := r.getCollection().UpdateOne(
		ctx,
		bson.M{"id": user.ID},
		bson.M{"$set": user},
	)

	return err
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.getCollection().DeleteOne(ctx, bson.M{"id": id})
	return err
}