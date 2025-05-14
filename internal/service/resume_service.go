package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"ai-interview/internal/database"
	"ai-interview/internal/model"
	"ai-interview/internal/repository"
	"ai-interview/pkg/utils"
)

// ResumeService 简历服务
type ResumeService struct {
	resumeRepo *repository.ResumeRepository
}

// NewResumeService 创建新的简历服务
func NewResumeService() *ResumeService {
	return &ResumeService{
		resumeRepo: repository.NewResumeRepository(),
	}
}

// UploadResume 上传并解析简历
func (s *ResumeService) UploadResume(userID uint, file *multipart.FileHeader, uploadDir string) (*model.Resume, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// 生成文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))

	// 保存文件
	filePath, err := utils.SaveUploadedFile(src, filename, uploadDir)
	if err != nil {
		return nil, err
	}

	// 解析简历内容
	parser := utils.NewResumeParser(filePath)
	content, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse resume: %v", err)
	}

	// 创建简历记录
	resume := &model.Resume{
		UserID:   userID,
		FileName: filename,
		FilePath: filePath,
		FileSize: file.Size,
		Content:  content,
	}

	// 保存到数据库
	if err := database.DB.Create(resume).Error; err != nil {
		return nil, fmt.Errorf("failed to save resume: %v", err)
	}

	return resume, nil
}

// CreateResume 创建简历
func (s *ResumeService) CreateResume(resume *model.Resume) error {
	return s.resumeRepo.Create(resume)
}

// GetResume 获取简历
func (s *ResumeService) GetResume(id uint) (*model.Resume, error) {
	return s.resumeRepo.GetByID(id)
}

// GetUserResumes 获取用户的所有简历
func (s *ResumeService) GetUserResumes(userID uint) ([]model.Resume, error) {
	return s.resumeRepo.GetByUserID(userID)
}

// DeleteResume 删除简历
func (s *ResumeService) DeleteResume(id uint) error {
	return s.resumeRepo.Delete(id)
}
