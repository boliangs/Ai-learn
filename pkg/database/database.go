package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.name"),
	)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 设置连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	// 启用日志
	db.LogMode(viper.GetBool("database.log_mode"))
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
