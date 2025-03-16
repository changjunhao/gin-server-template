package repository

import (
	"gin-server-template/internal/config"
	"gin-server-template/internal/entity"
	"gin-server-template/internal/repository/mongodb"
	"gin-server-template/internal/repository/mysql"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// Create 创建用户
	Create(user *entity.User) error

	// GetByID 根据ID获取用户
	GetByID(id uint) (*entity.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*entity.User, error)

	// ExistsByUsername 检查用户名是否存在
	ExistsByUsername(username string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(email string) (bool, error)

	// Update 更新用户信息
	Update(user *entity.User) error

	// Delete 删除用户
	Delete(id uint) error
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() UserRepository {
	// 根据配置决定使用哪种实现
	// 如果明确指定使用MySQL或MongoDB，则返回对应实现
	// 否则返回模拟实现用于开发和测试

	// 获取当前配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err == nil {
		switch cfg.Database.Driver {
		case "mysql":
			return mysql.NewUserRepository()
		case "mongodb":
			return mongodb.NewUserRepository()
		}
	}

	// 默认返回模拟实现
	return newMockUserRepository()
}

// 模拟实现，用于开发和测试
type mockUserRepository struct {
	users  map[uint]*entity.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[uint]*entity.User),
		nextID: 1,
	}
}

func (r *mockUserRepository) Create(user *entity.User) error {
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	return nil
}

func (r *mockUserRepository) GetByID(id uint) (*entity.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *mockUserRepository) GetByUsername(username string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, nil
}

func (r *mockUserRepository) ExistsByUsername(username string) (bool, error) {
	for _, user := range r.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockUserRepository) ExistsByEmail(email string) (bool, error) {
	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockUserRepository) Update(user *entity.User) error {
	_, exists := r.users[user.ID]
	if !exists {
		return nil
	}
	r.users[user.ID] = user
	return nil
}

func (r *mockUserRepository) Delete(id uint) error {
	delete(r.users, id)
	return nil
}
