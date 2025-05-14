package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"ai-interview/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ResumeController 简历控制器
type ResumeController struct {
	resumeService *service.ResumeService
}

// NewResumeController 创建新的简历控制器
func NewResumeController() *ResumeController {
	return &ResumeController{
		resumeService: service.NewResumeService(),
	}
}

// UploadResume 上传简历
func (c *ResumeController) UploadResume(ctx *gin.Context) {
	// 获取用户ID（这里假设已经通过中间件设置了用户信息）
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("resume")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 检查文件大小
	maxSize := viper.GetInt("upload.max_size") * 1024 * 1024 // 转换为字节
	if file.Size > int64(maxSize) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件类型
	allowedTypes := viper.GetStringSlice("upload.allowed_types")
	fileExt := filepath.Ext(file.Filename)
	isAllowed := false
	for _, ext := range allowedTypes {
		if fileExt == ext {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型"})
		return
	}

	// 上传并解析简历
	resume, err := c.resumeService.UploadResume(userID, file, viper.GetString("upload.save_path"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"data":    resume,
	})
}

// GetResume 获取简历信息
func (c *ResumeController) GetResume(ctx *gin.Context) {
	// 获取简历ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的简历ID"})
		return
	}

	// 获取简历信息
	resume, err := c.resumeService.GetResume(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "简历不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resume,
	})
}

// GetUserResumes 获取用户的所有简历
func (c *ResumeController) GetUserResumes(ctx *gin.Context) {
	// 获取用户ID
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 获取用户的所有简历
	resumes, err := c.resumeService.GetUserResumes(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取简历列表失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resumes,
	})
}

// DeleteResume 删除简历
func (c *ResumeController) DeleteResume(ctx *gin.Context) {
	// 获取简历ID
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的简历ID"})
		return
	}

	// 删除简历
	if err := c.resumeService.DeleteResume(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除简历失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
