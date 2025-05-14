package model

import (
	"time"
)

// Question 面试问题记录
type Question struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	InterviewID uint      `gorm:"not null" json:"interview_id"`
	Content     string    `gorm:"type:text;not null" json:"content"` // AI生成的问题内容
	UserAnswer  string    `gorm:"type:text" json:"user_answer"`      // 用户的回答
	AIFeedback  string    `gorm:"type:text" json:"ai_feedback"`      // AI的反馈
	Score       float64   `gorm:"type:decimal(5,2)" json:"score"`    // 评分
	Category    string    `gorm:"size:50" json:"category"`           // 问题类别（由AI生成）
	CreatedAt   time.Time `json:"created_at"`
	Interview   Interview `gorm:"foreignKey:InterviewID" json:"interview,omitempty"`
}
