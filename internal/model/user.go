package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(32);unique_index;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	Role     string `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}
