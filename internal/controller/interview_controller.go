package controller

import (
	"net/http"
	"strconv"

	"ai-interview/internal/service"
	"ai-interview/pkg/utils"

	"github.com/gin-gonic/gin"
)

// InterviewController 面试控制器
type InterviewController struct {
	interviewService *service.InterviewService
}

// NewInterviewController 创建新的面试控制器
func NewInterviewController() *InterviewController {
	return &InterviewController{
		interviewService: service.NewInterviewService(),
	}
}

// GenerateQuestions 生成面试题
func (c *InterviewController) GenerateQuestions(ctx *gin.Context) {
	// 获取简历ID
	resumeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的简历ID")
		return
	}

	// 获取面试类型
	interviewType := ctx.Query("type")
	if interviewType == "" {
		interviewType = "technical" // 默认为技术面试
	}

	// 生成面试题
	questions, err := c.interviewService.GenerateQuestions(uint(resumeID), interviewType)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "生成面试题失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"questions": questions,
	})
}

// EvaluateAnswer 评估答案
func (c *InterviewController) EvaluateAnswer(ctx *gin.Context) {
	var req struct {
		QuestionID uint   `json:"question_id" binding:"required"`
		Answer     string `json:"answer" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的请求参数")
		return
	}

	// 评估答案
	evaluation, err := c.interviewService.EvaluateAnswer(req.QuestionID, req.Answer)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "评估答案失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"evaluation": evaluation,
	})
}

// GenerateFeedback 生成面试反馈
func (c *InterviewController) GenerateFeedback(ctx *gin.Context) {
	// 获取简历ID
	resumeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的简历ID")
		return
	}

	// 生成反馈
	feedback, err := c.interviewService.GenerateFeedback(uint(resumeID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "生成反馈失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"feedback": feedback,
	})
}

// GetInterviewHistory 获取面试历史
func (c *InterviewController) GetInterviewHistory(ctx *gin.Context) {
	// 获取简历ID
	resumeID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的简历ID")
		return
	}

	// 获取面试历史
	history, err := c.interviewService.GetInterviewHistory(uint(resumeID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取面试历史失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"history": history,
	})
}
