package router

import (
	"ai-interview/internal/controller"
	"ai-interview/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	// 控制器
	userController := controller.NewUserController()
	resumeController := controller.NewResumeController()
	interviewController := controller.NewInterviewController()

	// 公开路由
	public := r.Group("/api")
	{
		// 用户认证
		public.POST("/auth/register", userController.Register)
		public.POST("/auth/login", userController.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.Auth())
	{
		// 用户相关
		authorized.GET("/user/profile", userController.GetProfile)
		authorized.PUT("/user/profile", userController.UpdateProfile)

		// 简历相关
		authorized.POST("/resumes/upload", resumeController.UploadResume)
		authorized.GET("/resumes/:id", resumeController.GetResume)
		authorized.GET("/resumes/user", resumeController.GetUserResumes)
		authorized.DELETE("/resumes/:id", resumeController.DeleteResume)

		// 面试相关
		authorized.POST("/resumes/:id/interview", interviewController.GenerateQuestions)
		authorized.POST("/interview/answer", interviewController.EvaluateAnswer)
		authorized.GET("/resumes/:id/feedback", interviewController.GenerateFeedback)
		authorized.GET("/resumes/:id/history", interviewController.GetInterviewHistory)
	}

	return r
}
