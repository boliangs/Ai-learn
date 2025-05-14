package repository

import (
	"ai-interview/internal/model"
	"ai-interview/pkg/database"

	"github.com/jinzhu/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建新的用户仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
