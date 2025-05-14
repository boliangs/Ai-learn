package model

import (
	"time"
)

// Interview 面试记录
type Interview struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	UserID     uint       `gorm:"not null" json:"user_id"`
	Title      string     `gorm:"size:100;not null" json:"title"` // 面试标题
	Status     string     `gorm:"type:enum('pending','in_progress','completed');default:'pending'" json:"status"`
	TotalScore float64    `gorm:"type:decimal(5,2)" json:"total_score"` // 总分
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	ResumeID   uint       `gorm:"not null" json:"resume_id"` // 关联的简历ID
	JobTitle   string     `gorm:"size:100" json:"job_title"` // 应聘职位
	Company    string     `gorm:"size:100" json:"company"`   // 公司名称
	CreatedAt  time.Time  `json:"created_at"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Resume     Resume     `gorm:"foreignKey:ResumeID" json:"resume,omitempty"`
	Questions  []Question `gorm:"foreignKey:InterviewID" json:"questions,omitempty"`
}

type InterviewQuestion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	InterviewID uint      `gorm:"not null" json:"interview_id"`
	QuestionID  uint      `gorm:"not null" json:"question_id"`
	UserAnswer  string    `gorm:"type:text" json:"user_answer"`
	AIFeedback  string    `gorm:"type:text" json:"ai_feedback"`
	Score       float64   `gorm:"type:decimal(5,2)" json:"score"`
	CreatedAt   time.Time `json:"created_at"`
	Question    Question  `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}
