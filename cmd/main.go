package main

import (
	"log"

	"ai-interview/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 执行数据库迁移
	database.Migrate()
	// 初始化基础数据
	database.Seed()

	// 设置gin模式
	gin.SetMode(viper.GetString("server.mode"))

	// 初始化路由
	r := gin.Default()

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	port := viper.GetString("server.port")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *gin.Engine) {
	// TODO: 注册路由
	// 示例路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
