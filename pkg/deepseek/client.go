package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Client DeepSeek客户端
type Client struct {
	apiKey     string
	apiURL     string
	model      string
	httpClient *http.Client
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewClient 创建新的DeepSeek客户端
func NewClient() *Client {
	return &Client{
		apiKey: viper.GetString("deepseek.api_key"),
		apiURL: viper.GetString("deepseek.api_url"),
		model:  viper.GetString("deepseek.model"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateText 生成文本
func (c *Client) GenerateText(prompt string) (string, error) {
	// 构建请求
	req := ChatRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的面试官，擅长生成面试题、评估答案和提供反馈。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.apiURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败: %s", string(respBody))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应内容
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return chatResp.Choices[0].Message.Content, nil
}
