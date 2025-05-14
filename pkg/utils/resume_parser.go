package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/unidoc/unioffice/document"
)

var (
	// 正则表达式模式
	phonePattern = regexp.MustCompile(`1[3-9]\d{9}`)
	emailPattern = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// 关键词
	educationKeywords  = []string{"教育经历", "教育背景", "学历", "学校", "大学", "学院"}
	experienceKeywords = []string{"工作经历", "工作经验", "实习经历", "项目经验", "工作背景"}
)

// ResumeParser 简历解析器
type ResumeParser struct {
	FilePath string
}

// NewResumeParser 创建新的简历解析器
func NewResumeParser(filePath string) *ResumeParser {
	return &ResumeParser{
		FilePath: filePath,
	}
}

// Parse 解析简历文件
func (p *ResumeParser) Parse() (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(p.FilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", p.FilePath)
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(p.FilePath))
	if ext != ".docx" && ext != ".doc" {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	// 打开文档
	doc, err := document.Open(p.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open document: %v", err)
	}
	defer doc.Close()

	// 提取文本内容
	var content strings.Builder
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			content.WriteString(run.Text())
		}
		content.WriteString("\n")
	}

	return content.String(), nil
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(src io.Reader, filename string, uploadDir string) (string, error) {
	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 生成文件路径
	filePath := filepath.Join(uploadDir, filename)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filePath, nil
}

// ExtractResumeInfo 从简历内容中提取关键信息
func ExtractResumeInfo(content string) map[string]string {
	info := make(map[string]string)

	// 提取姓名
	if name := extractName(content); name != "" {
		info["name"] = name
	}

	// 提取电话
	if phone := extractPhone(content); phone != "" {
		info["phone"] = phone
	}

	// 提取邮箱
	if email := extractEmail(content); email != "" {
		info["email"] = email
	}

	// 提取教育经历
	if education := extractEducation(content); education != "" {
		info["education"] = education
	}

	// 提取工作经历
	if experience := extractExperience(content); experience != "" {
		info["experience"] = experience
	}

	return info
}

// 辅助函数：提取姓名
func extractName(content string) string {
	// 假设姓名在文档开头，且长度在2-4个字符之间
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if len(firstLine) >= 2 && len(firstLine) <= 4 {
			return firstLine
		}
	}
	return ""
}

// 辅助函数：提取电话
func extractPhone(content string) string {
	matches := phonePattern.FindString(content)
	return matches
}

// 辅助函数：提取邮箱
func extractEmail(content string) string {
	matches := emailPattern.FindString(content)
	return matches
}

// 辅助函数：提取教育经历
func extractEducation(content string) string {
	lines := strings.Split(content, "\n")
	var education strings.Builder
	var inEducationSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否进入教育经历部分
		if !inEducationSection {
			for _, keyword := range educationKeywords {
				if strings.Contains(line, keyword) {
					inEducationSection = true
					education.WriteString(line + "\n")
					break
				}
			}
			continue
		}

		// 检查是否离开教育经历部分
		for _, keyword := range experienceKeywords {
			if strings.Contains(line, keyword) {
				return education.String()
			}
		}

		// 收集教育经历内容
		education.WriteString(line + "\n")
	}

	return education.String()
}

// 辅助函数：提取工作经历
func extractExperience(content string) string {
	lines := strings.Split(content, "\n")
	var experience strings.Builder
	var inExperienceSection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否进入工作经历部分
		if !inExperienceSection {
			for _, keyword := range experienceKeywords {
				if strings.Contains(line, keyword) {
					inExperienceSection = true
					experience.WriteString(line + "\n")
					break
				}
			}
			continue
		}

		// 检查是否离开工作经历部分
		for _, keyword := range educationKeywords {
			if strings.Contains(line, keyword) {
				return experience.String()
			}
		}

		// 收集工作经历内容
		experience.WriteString(line + "\n")
	}

	return experience.String()
}
