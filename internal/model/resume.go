package model

import (
	"time"
)

type Resume struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FilePath  string    `gorm:"size:255;not null" json:"file_path"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
