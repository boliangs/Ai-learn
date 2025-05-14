package model

import (
	"time"

	"gorm.io/gorm"
)

// Resume 简历模型
type Resume struct {
	gorm.Model
	UserID   uint   `json:"user_id" gorm:"index"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`
	Content  string `json:"content" gorm:"type:text"`
}

// TableName 指定表名
func (Resume) TableName() string {
	return "resumes"
}

// BeforeCreate 创建前的钩子
func (r *Resume) BeforeCreate(tx *gorm.DB) error {
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (r *Resume) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}
