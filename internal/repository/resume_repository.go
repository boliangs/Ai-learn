package repository

import (
	"ai-interview/internal/database"
	"ai-interview/internal/model"
)

// ResumeRepository 简历仓库
type ResumeRepository struct{}

// NewResumeRepository 创建新的简历仓库
func NewResumeRepository() *ResumeRepository {
	return &ResumeRepository{}
}

// Create 创建简历
func (r *ResumeRepository) Create(resume *model.Resume) error {
	return database.DB.Create(resume).Error
}

// GetByID 根据ID获取简历
func (r *ResumeRepository) GetByID(id uint) (*model.Resume, error) {
	var resume model.Resume
	err := database.DB.First(&resume, id).Error
	return &resume, err
}

// GetByUserID 获取用户的所有简历
func (r *ResumeRepository) GetByUserID(userID uint) ([]model.Resume, error) {
	var resumes []model.Resume
	err := database.DB.Where("user_id = ?", userID).Find(&resumes).Error
	return resumes, err
}

// Delete 删除简历
func (r *ResumeRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Resume{}, id).Error
}
