package model

import (
	"gorm.io/gorm"
)

// Interview 面试记录
type Interview struct {
	gorm.Model
	ResumeID  uint       `json:"resume_id" gorm:"index"`
	Type      string     `json:"type"`   // 面试类型：technical, behavioral, etc.
	Status    string     `json:"status"` // 面试状态：pending, completed, etc.
	Questions []Question `json:"questions" gorm:"foreignKey:InterviewID"`
}

// Question 面试问题
type Question struct {
	gorm.Model
	InterviewID        uint    `json:"interview_id" gorm:"index"`
	Question           string  `json:"question"`
	EvaluationCriteria string  `json:"evaluation_criteria"`
	Difficulty         string  `json:"difficulty"`
	Answer             *Answer `json:"answer" gorm:"foreignKey:QuestionID"`
}

// Answer 面试答案
type Answer struct {
	gorm.Model
	QuestionID uint   `json:"question_id" gorm:"index"`
	Content    string `json:"content"`
	Score      int    `json:"score"`
	Feedback   string `json:"feedback"`
}

// Evaluation 评估结果
type Evaluation struct {
	Score       int      `json:"score"`
	Evaluation  string   `json:"evaluation"`
	Suggestions []string `json:"suggestions"`
}

// Feedback 面试反馈
type Feedback struct {
	gorm.Model
	ResumeID               uint     `json:"resume_id" gorm:"index"`
	OverallEvaluation      string   `json:"overall_evaluation"`
	Strengths              []string `json:"strengths" gorm:"type:json"`
	Weaknesses             []string `json:"weaknesses" gorm:"type:json"`
	ImprovementSuggestions []string `json:"improvement_suggestions" gorm:"type:json"`
	DevelopmentSuggestions []string `json:"development_suggestions" gorm:"type:json"`
}
