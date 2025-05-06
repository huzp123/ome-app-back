package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"ome-app-back/config"
	"ome-app-back/internal/model"
)

// AIService 处理AI相关服务
type AIService struct {
	apiKey       string
	apiURL       string
	maxTokens    int
	temperature  float64
	defaultModel string
}

// NewAIService 创建AI服务实例
func NewAIService(cfg *config.AIConfig) *AIService {
	return &AIService{
		apiKey:       cfg.APIKey,
		apiURL:       cfg.APIURL,
		maxTokens:    cfg.MaxTokens,
		temperature:  cfg.Temperature,
		defaultModel: cfg.Model,
	}
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Model       string                   `json:"model"`
	Messages    []map[string]interface{} `json:"messages"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float64                  `json:"temperature"`
}

// ChatResponse AI响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatWithAI 发送聊天请求到AI
func (s *AIService) ChatWithAI(messages []model.OpenAIMessage) (string, error) {
	if s.apiKey == "" {
		return "", errors.New("AI API密钥未配置")
	}

	// 转换为API请求格式
	apiMessages := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		apiMessages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	// 准备请求体
	requestBody := ChatRequest{
		Model:       s.defaultModel,
		Messages:    apiMessages,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败: HTTP %d, %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var responseData ChatResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 验证响应
	if len(responseData.Choices) == 0 {
		return "", errors.New("API返回的选择项为空")
	}

	// 返回AI回复内容
	return responseData.Choices[0].Message.Content, nil
}

// ConvertToMessages 将聊天消息转换为API请求格式
func (s *AIService) ConvertToMessages(chatMessages []model.ChatMessage) []model.OpenAIMessage {
	messages := make([]model.OpenAIMessage, len(chatMessages))
	for i, msg := range chatMessages {
		messages[i] = model.OpenAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	return messages
}

// GetSystemMessageForChat 获取聊天的系统消息
func (s *AIService) GetSystemMessageForChat() model.OpenAIMessage {
	return model.OpenAIMessage{
		Role: "system",
		Content: `你是一个专业的营养健康助手。你可以:
1. 提供健康饮食建议
2. 帮助用户了解食物的营养价值
3. 回答与健康、饮食相关的问题
4. 根据用户健康目标给出个性化建议

请用友好、专业的方式回答问题，避免医疗诊断，只提供通用健康信息。`,
	}
}

// GetSystemMessageForFoodRecognition 获取食物识别的系统消息
func (s *AIService) GetSystemMessageForFoodRecognition() model.OpenAIMessage {
	return model.OpenAIMessage{
		Role: "system",
		Content: `你是一个专业的食物识别和营养分析AI。你的任务是:
1. 识别图片中的食物
2. 估算每种食物的大致数量
3. 计算总热量(千卡)和主要营养素含量(蛋白质、碳水化合物、脂肪，单位为克)
4. 简要分析这顿饭的营养价值和健康性

请以JSON格式输出结果:
{
  "foods": [
    {"name": "食物名称", "quantity": "份量描述", "calories": 估计热量}
  ],
  "nutrition": {
    "calories": 总热量,
    "protein": 蛋白质克数,
    "carbs": 碳水克数,
    "fat": 脂肪克数
  },
  "analysis": "对这顿饭的简短营养分析"
}

只返回JSON内容，不要添加其他文字说明。`,
	}
}

// AnalyzeImageWithAI 分析图片内容（使用base64编码）
func (s *AIService) AnalyzeImageWithAI(base64Image string, prompt string) (string, error) {
	// 创建带有图像的消息内容
	systemMessage := map[string]interface{}{
		"role":    "system",
		"content": s.GetSystemMessageForFoodRecognition().Content,
	}

	userMessage := map[string]interface{}{
		"role": "user",
		"content": []interface{}{
			map[string]string{
				"type": "text",
				"text": "图片中的食物是什么？请分析营养成分。",
			},
			map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]string{
					"url": "data:image/jpeg;base64," + base64Image,
				},
			},
		},
	}

	messages := []map[string]interface{}{systemMessage, userMessage}

	// 准备请求体
	requestBody := ChatRequest{
		Model:       s.defaultModel,
		Messages:    messages,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败: HTTP %d, %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var responseData ChatResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 验证响应
	if len(responseData.Choices) == 0 {
		return "", errors.New("API返回的选择项为空")
	}

	// 返回AI回复内容
	return responseData.Choices[0].Message.Content, nil
}
