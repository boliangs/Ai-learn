package database

import (
	"log"

	"ai-interview/internal/model"
)

// Migrate 执行数据库迁移
func Migrate() {
	// 自动迁移数据库结构
	err := DB.AutoMigrate(
		&model.User{},
		&model.Interview{},
		&model.Question{},
		&model.Resume{},
	).Error

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
}

// Seed 初始化基础数据
func Seed() {
	// 检查是否已存在管理员用户
	var adminCount int64
	DB.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount == 0 {
		// 创建默认管理员用户
		admin := model.User{
			Username: "admin",
			Password: "$2a$10$XOPbrlUPQdwdJUpSrIF6X.LbE14qsMmKGhM1A8W9iqDp0.3tQ5/Ym", // 密码: admin123
			Email:    "admin@example.com",
			Role:     "admin",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
		} else {
			log.Println("Default admin user created successfully")
		}
	}
}
