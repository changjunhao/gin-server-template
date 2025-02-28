package service

import (
	"errors"
	"gin-server/internal/entity"
	"gin-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register 注册新用户
func (s *UserService) Register(user *entity.User) error {
	// 检查用户名是否已存在
	exist, err := s.userRepo.ExistsByUsername(user.Username)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if user.Email != "" {
		exist, err = s.userRepo.ExistsByEmail(user.Email)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("邮箱已被使用")
		}
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 创建用户
	return s.userRepo.Create(user)
}

// VerifyCredentials 验证用户凭证
func (s *UserService) VerifyCredentials(username, password string) (*entity.User, error) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*entity.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *entity.User) error {
	return s.userRepo.Update(user)
}