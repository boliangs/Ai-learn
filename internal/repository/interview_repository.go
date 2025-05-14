package repository

import (
	"ai-interview/internal/database"
	"ai-interview/internal/model"
)

// InterviewRepository 面试仓库
type InterviewRepository struct{}

// NewInterviewRepository 创建新的面试仓库
func NewInterviewRepository() *InterviewRepository {
	return &InterviewRepository{}
}

// Create 创建面试记录
func (r *InterviewRepository) Create(interview *model.Interview) error {
	return database.DB.Create(interview).Error
}

// GetByID 根据ID获取面试记录
func (r *InterviewRepository) GetByID(id uint) (*model.Interview, error) {
	var interview model.Interview
	err := database.DB.Preload("Questions").Preload("Questions.Answer").First(&interview, id).Error
	return &interview, err
}

// GetInterviewHistory 获取面试历史
func (r *InterviewRepository) GetInterviewHistory(resumeID uint) ([]model.Interview, error) {
	var interviews []model.Interview
	err := database.DB.Where("resume_id = ?", resumeID).
		Preload("Questions").
		Preload("Questions.Answer").
		Order("created_at DESC").
		Find(&interviews).Error
	return interviews, err
}

// CreateQuestion 创建面试问题
func (r *InterviewRepository) CreateQuestion(question *model.Question) error {
	return database.DB.Create(question).Error
}

// GetQuestionByID 根据ID获取问题
func (r *InterviewRepository) GetQuestionByID(id uint) (*model.Question, error) {
	var question model.Question
	err := database.DB.First(&question, id).Error
	return &question, err
}

// CreateAnswer 创建答案
func (r *InterviewRepository) CreateAnswer(answer *model.Answer) error {
	return database.DB.Create(answer).Error
}

// CreateFeedback 创建反馈
func (r *InterviewRepository) CreateFeedback(feedback *model.Feedback) error {
	return database.DB.Create(feedback).Error
}

// GetFeedbackByResumeID 根据简历ID获取反馈
func (r *InterviewRepository) GetFeedbackByResumeID(resumeID uint) (*model.Feedback, error) {
	var feedback model.Feedback
	err := database.DB.Where("resume_id = ?", resumeID).First(&feedback).Error
	return &feedback, err
}
