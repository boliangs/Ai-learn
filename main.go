package main

import (
	"log"

	"ai-interview/internal/database"
	"ai-interview/internal/router"

	"github.com/spf13/viper"
)

func main() {
	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
