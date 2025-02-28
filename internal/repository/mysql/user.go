package mysql

import (
	"errors"
	"gin-server/internal/database"
	"gin-server/internal/entity"

	"gorm.io/gorm"
)

// UserRepository MySQL实现的用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建MySQL用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}