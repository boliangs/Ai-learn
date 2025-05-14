package service

import (
	"ai-interview/internal/model"
	"ai-interview/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建新的用户服务
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 创建用户
	return s.userRepo.Create(user)
}

// VerifyUser 验证用户
func (s *UserService) VerifyUser(username, password string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return s.userRepo.Update(id, updates)
}
