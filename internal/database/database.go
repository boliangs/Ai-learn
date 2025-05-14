package database

import (
	"fmt"
	"log"

	"ai-interview/internal/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GormDB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	// 获取数据库配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
	)

	// 配置GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := GormDB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	// 设置连接最大空闲时间
	sqlDB.SetConnMaxIdleTime(viper.GetDuration("database.conn_max_idle_time"))
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(viper.GetDuration("database.conn_max_lifetime"))

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动迁移数据库表失败: %v", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	return GormDB.AutoMigrate(
		&model.User{},
		&model.Resume{},
		&model.Interview{},
		&model.Question{},
		&model.Answer{},
		&model.Feedback{},
	)
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := GormDB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}
	return sqlDB.Close()
}
