package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai-interview/internal/controller"
	"ai-interview/internal/database"
	"ai-interview/internal/middleware"

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

	// 确保上传目录存在
	uploadDir := viper.GetString("upload.save_path")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// 获取绝对路径
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	viper.Set("upload.save_path", absUploadDir)

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

	// 注册中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Server starting at http://%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 创建控制器
	resumeController := controller.NewResumeController()

	// 简历相关路由
	resumeGroup := r.Group("/api/resumes")
	{
		// 上传简历
		resumeGroup.POST("/upload", resumeController.UploadResume)
		// 获取简历信息
		resumeGroup.GET("/:id", resumeController.GetResume)
		// 获取用户的所有简历
		resumeGroup.GET("/user", resumeController.GetUserResumes)
		// 删除简历
		resumeGroup.DELETE("/:id", resumeController.DeleteResume)
	}
}
