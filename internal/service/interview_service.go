package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"ai-interview/internal/model"
	"ai-interview/internal/repository"
	"ai-interview/pkg/deepseek"
)

// InterviewService 面试服务
type InterviewService struct {
	resumeRepo    *repository.ResumeRepository
	interviewRepo *repository.InterviewRepository
	deepseek      *deepseek.Client
}

// NewInterviewService 创建新的面试服务
func NewInterviewService() *InterviewService {
	return &InterviewService{
		resumeRepo:    repository.NewResumeRepository(),
		interviewRepo: repository.NewInterviewRepository(),
		deepseek:      deepseek.NewClient(),
	}
}

// GenerateQuestions 生成面试题
func (s *InterviewService) GenerateQuestions(resumeID uint, interviewType string) ([]model.Question, error) {
	// 获取简历信息
	resume, err := s.resumeRepo.GetByID(resumeID)
	if err != nil {
		return nil, fmt.Errorf("获取简历失败: %v", err)
	}

	// 构建提示词
	prompt := fmt.Sprintf(`请根据以下简历内容生成%s面试题：

简历内容：
%s

请生成5个面试题，每个问题应该：
1. 与简历内容相关
2. 考察候选人的专业能力
3. 包含评分标准
4. 难度适中

请以JSON格式返回，格式如下：
{
    "questions": [
        {
            "question": "问题内容",
            "evaluation_criteria": "评分标准",
            "difficulty": "难度等级"
        }
    ]
}`, interviewType, resume.Content)

	// 调用DeepSeek API
	response, err := s.deepseek.GenerateText(prompt)
	if err != nil {
		return nil, fmt.Errorf("生成面试题失败: %v", err)
	}

	// 解析响应
	questions, err := parseQuestions(response)
	if err != nil {
		return nil, fmt.Errorf("解析面试题失败: %v", err)
	}

	// 保存面试记录
	interview := &model.Interview{
		ResumeID: resumeID,
		Type:     interviewType,
		Status:   "completed",
	}
	if err := s.interviewRepo.Create(interview); err != nil {
		return nil, fmt.Errorf("保存面试记录失败: %v", err)
	}

	// 保存问题记录
	for _, q := range questions {
		q.InterviewID = interview.ID
		if err := s.interviewRepo.CreateQuestion(&q); err != nil {
			return nil, fmt.Errorf("保存问题记录失败: %v", err)
		}
	}

	return questions, nil
}

// EvaluateAnswer 评估答案
func (s *InterviewService) EvaluateAnswer(questionID uint, answer string) (*model.Evaluation, error) {
	// 获取问题信息
	question, err := s.interviewRepo.GetQuestionByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("获取问题信息失败: %v", err)
	}

	// 构建提示词
	prompt := fmt.Sprintf(`请评估以下面试答案：

问题：%s
评分标准：%s
答案：%s

请从以下几个方面进行评估：
1. 答案的完整性
2. 答案的准确性
3. 答案的深度
4. 表达的逻辑性
5. 改进建议

请以JSON格式返回，格式如下：
{
    "score": 85,
    "evaluation": "详细评价",
    "suggestions": ["改进建议1", "改进建议2"]
}`, question.Question, question.EvaluationCriteria, answer)

	// 调用DeepSeek API
	response, err := s.deepseek.GenerateText(prompt)
	if err != nil {
		return nil, fmt.Errorf("评估答案失败: %v", err)
	}

	// 解析响应
	evaluation, err := parseEvaluation(response)
	if err != nil {
		return nil, fmt.Errorf("解析评估结果失败: %v", err)
	}

	// 保存答案和评估结果
	answerRecord := &model.Answer{
		QuestionID: questionID,
		Content:    answer,
		Score:      evaluation.Score,
		Feedback:   evaluation.Evaluation,
	}
	if err := s.interviewRepo.CreateAnswer(answerRecord); err != nil {
		return nil, fmt.Errorf("保存答案记录失败: %v", err)
	}

	return evaluation, nil
}

// GenerateFeedback 生成面试反馈
func (s *InterviewService) GenerateFeedback(resumeID uint) (*model.Feedback, error) {
	// 获取面试历史
	history, err := s.interviewRepo.GetInterviewHistory(resumeID)
	if err != nil {
		return nil, fmt.Errorf("获取面试历史失败: %v", err)
	}

	// 构建提示词
	var historyText strings.Builder
	for _, h := range history {
		historyText.WriteString(fmt.Sprintf("面试类型：%s\n", h.Type))
		historyText.WriteString("问题及答案：\n")
		for _, q := range h.Questions {
			historyText.WriteString(fmt.Sprintf("问题：%s\n", q.Question))
			if q.Answer != nil {
				historyText.WriteString(fmt.Sprintf("答案：%s\n", q.Answer.Content))
				historyText.WriteString(fmt.Sprintf("评分：%d\n", q.Answer.Score))
				historyText.WriteString(fmt.Sprintf("反馈：%s\n", q.Answer.Feedback))
			}
			historyText.WriteString("---\n")
		}
	}

	prompt := fmt.Sprintf(`请根据以下面试历史生成综合反馈：

%s

请从以下几个方面提供反馈：
1. 整体表现评价
2. 优势分析
3. 不足分析
4. 改进建议
5. 发展建议

请以JSON格式返回，格式如下：
{
    "overall_evaluation": "整体评价",
    "strengths": ["优势1", "优势2"],
    "weaknesses": ["不足1", "不足2"],
    "improvement_suggestions": ["建议1", "建议2"],
    "development_suggestions": ["发展建议1", "发展建议2"]
}`, historyText.String())

	// 调用DeepSeek API
	response, err := s.deepseek.GenerateText(prompt)
	if err != nil {
		return nil, fmt.Errorf("生成反馈失败: %v", err)
	}

	// 解析响应
	feedback, err := parseFeedback(response)
	if err != nil {
		return nil, fmt.Errorf("解析反馈失败: %v", err)
	}

	// 保存反馈记录
	feedbackRecord := &model.Feedback{
		ResumeID:               resumeID,
		OverallEvaluation:      feedback.OverallEvaluation,
		Strengths:              feedback.Strengths,
		Weaknesses:             feedback.Weaknesses,
		ImprovementSuggestions: feedback.ImprovementSuggestions,
		DevelopmentSuggestions: feedback.DevelopmentSuggestions,
	}
	if err := s.interviewRepo.CreateFeedback(feedbackRecord); err != nil {
		return nil, fmt.Errorf("保存反馈记录失败: %v", err)
	}

	return feedback, nil
}

// GetInterviewHistory 获取面试历史
func (s *InterviewService) GetInterviewHistory(resumeID uint) ([]model.Interview, error) {
	return s.interviewRepo.GetInterviewHistory(resumeID)
}

// 解析面试题
func parseQuestions(response string) ([]model.Question, error) {
	var result struct {
		Questions []struct {
			Question           string `json:"question"`
			EvaluationCriteria string `json:"evaluation_criteria"`
			Difficulty         string `json:"difficulty"`
		} `json:"questions"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	questions := make([]model.Question, len(result.Questions))
	for i, q := range result.Questions {
		questions[i] = model.Question{
			Question:           q.Question,
			EvaluationCriteria: q.EvaluationCriteria,
			Difficulty:         q.Difficulty,
		}
	}

	return questions, nil
}

// 解析评估结果
func parseEvaluation(response string) (*model.Evaluation, error) {
	var result struct {
		Score       int      `json:"score"`
		Evaluation  string   `json:"evaluation"`
		Suggestions []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &model.Evaluation{
		Score:       result.Score,
		Evaluation:  result.Evaluation,
		Suggestions: result.Suggestions,
	}, nil
}

// 解析反馈
func parseFeedback(response string) (*model.Feedback, error) {
	var result struct {
		OverallEvaluation      string   `json:"overall_evaluation"`
		Strengths              []string `json:"strengths"`
		Weaknesses             []string `json:"weaknesses"`
		ImprovementSuggestions []string `json:"improvement_suggestions"`
		DevelopmentSuggestions []string `json:"development_suggestions"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return &model.Feedback{
		OverallEvaluation:      result.OverallEvaluation,
		Strengths:              result.Strengths,
		Weaknesses:             result.Weaknesses,
		ImprovementSuggestions: result.ImprovementSuggestions,
		DevelopmentSuggestions: result.DevelopmentSuggestions,
	}, nil
}
