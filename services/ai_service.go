package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	"ome-app-back/config"
	"ome-app-back/models"
)

// AIService 处理AI相关服务
type AIService struct {
	apiKey       string
	apiURL       string
	maxTokens    int
	temperature  float64
	defaultModel string
	client       *http.Client
	testMode     bool // 是否为测试模式
}

// NewAIService 创建AI服务实例
func NewAIService(cfg *config.AIConfig) *AIService {
	// 创建具有适当超时设置的HTTP客户端
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			// 如果配置了代理URL，使用它
			if cfg.ProxyURL != "" {
				log.Printf("[AI服务] 使用代理: %s", cfg.ProxyURL)
				return url.Parse(cfg.ProxyURL)
			}
			// 否则使用环境变量中的代理
			return http.ProxyFromEnvironment(req)
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second, // 总请求超时
	}

	service := &AIService{
		apiKey:       cfg.APIKey,
		apiURL:       cfg.APIURL,
		maxTokens:    cfg.MaxTokens,
		temperature:  cfg.Temperature,
		defaultModel: cfg.Model,
		client:       client,
		testMode:     cfg.TestMode, // 从配置读取测试模式
	}

	if service.testMode {
		log.Printf("[AI服务] 运行在测试模式，将使用预定义响应")
	}

	log.Printf("[AI服务] 初始化完成，使用模型: %s, API URL: %s", service.defaultModel, service.apiURL)

	return service
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Model       string                   `json:"model"`
	Messages    []map[string]interface{} `json:"messages"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float64                  `json:"temperature"`
	Stream      bool                     `json:"stream,omitempty"`
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

// ChatStreamChoice 流式响应中的选择
type ChatStreamChoice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

// ChatStreamResponse AI流式响应结构
type ChatStreamResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []ChatStreamChoice `json:"choices"`
}

// makeAPIRequest 执行API请求，包含重试逻辑
func (s *AIService) makeAPIRequest(jsonData []byte, logPrefix string) ([]byte, error) {
	maxRetries := 2
	retryDelay := 1 * time.Second

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("%s 第%d次重试, 等待%v...", logPrefix, attempt, retryDelay)
			time.Sleep(retryDelay)
			retryDelay *= 2 // 指数级增加重试延迟
		}

		responseBody, err := s.executeRequest(jsonData, logPrefix, attempt)
		if err == nil {
			return responseBody, nil
		}

		lastErr = err
		log.Printf("%s 请求失败(尝试%d/%d): %v", logPrefix, attempt+1, maxRetries+1, err)

		// 根据错误类型决定是否重试
		if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
			continue // 网络超时或临时错误，重试
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			continue // 上下文超时或取消，重试
		} else if s.isConnectionReset(err) {
			continue // 连接重置错误，重试
		}

		// 其他错误不重试
		break
	}

	return nil, fmt.Errorf("在%d次尝试后仍然失败: %w", maxRetries+1, lastErr)
}

// isConnectionReset 检查错误是否为连接重置类型
func (s *AIService) isConnectionReset(err error) bool {
	errStr := err.Error()
	return errStr == "EOF" ||
		errStr == "unexpected EOF" ||
		errStr == "connection reset by peer" ||
		errStr == "broken pipe"
}

// executeRequest 执行单次HTTP请求并返回响应体
func (s *AIService) executeRequest(jsonData []byte, logPrefix string, attempt int) ([]byte, error) {
	// 创建带跟踪的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// 启用HTTP跟踪 (简化跟踪，只保留关键连接信息)
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			// 简化日志，移除DNS查询开始信息
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			if info.Err != nil {
				log.Printf("%s DNS查询出错: %v", logPrefix, info.Err)
			}
		},
		ConnectStart: func(network, addr string) {
			// 简化日志，移除连接开始信息
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("%s 连接失败: %s, 错误: %v", logPrefix, addr, err)
			}
		},
		GotConn: func(info httptrace.GotConnInfo) {
			// 简化日志，移除获得连接信息
		},
	}
	ctx = httptrace.WithClientTrace(ctx, trace)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("User-Agent", "OME-Nutrition-App/1.0")

	// 记录请求开始时间
	startTime := time.Now()
	log.Printf("%s 发送请求(尝试%d)...", logPrefix, attempt+1)

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		netErr, isNetErr := err.(net.Error)
		if isNetErr {
			log.Printf("%s 网络错误: %v (超时=%v, 临时=%v)",
				logPrefix, err, netErr.Timeout(), netErr.Temporary())
		} else {
			log.Printf("%s 请求错误: %v", logPrefix, err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 计算请求耗时
	requestDuration := time.Since(startTime)
	log.Printf("%s 收到响应: HTTP状态=%d, 耗时=%.2f秒",
		logPrefix, resp.StatusCode, requestDuration.Seconds())

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败: HTTP %d, %s", resp.StatusCode, truncateResponse(body, 200))
	}

	return body, nil
}

// truncateResponse 截断响应体，避免日志过长
func truncateResponse(response []byte, maxLen int) string {
	if len(response) <= maxLen {
		return string(response)
	}
	return string(response[:maxLen]) + "...(已截断)"
}

// 预定义的测试响应
const testChatResponse = `有氧运动主要以提高心肺功能为主，但也可以辅助增肌。以下是一些适合的有氧运动：

1. **游泳**：全身运动，增强肌肉力量和耐力。
2. **骑自行车**：可以调节强度，有助于腿部肌肉发展。
3. **跳绳**：提高心率，锻炼全身肌肉，便于携带。
4. **慢跑或快走**：增强心肺功能，同时可以保持体重。
5. **舞蹈**：愉悦的同时锻炼肌肉力量和协调性。

结合有氧运动和力量训练，可以更好地实现增肌目标。记得饮食均衡，确保摄入足够的蛋白质！`

// 简化版食物识别测试响应
const testFoodRecognitionResponse = `{
  "foods": [
    {"name": "鸡胸肉", "quantity": "约200克", "calories": 330},
    {"name": "糙米", "quantity": "约150克", "calories": 240},
    {"name": "西兰花", "quantity": "约100克", "calories": 55}
  ],
  "nutrition": {
    "calories_intake": 625,
    "protein_intake_g": 45,
    "carb_intake_g": 60,
    "fat_intake_g": 15
  },
  "analysis": "这是一顿营养均衡的健康餐，蛋白质含量丰富，适合健身增肌人群。碳水化合物以复合碳水为主，提供持久能量。添加更多蔬菜可增加纤维和微量元素摄入。"
}`

// ChatWithAI 发送聊天请求到AI
func (s *AIService) ChatWithAI(messages []models.OpenAIMessage) (string, error) {
	logPrefix := "[AI聊天]"

	// 测试模式直接返回预定义响应
	if s.testMode {
		log.Printf("%s 测试模式，返回预定义响应", logPrefix)
		log.Printf("%s 测试模式原始完整回复内容: %s", logPrefix, testChatResponse)
		// 将换行符替换为\n字符串，以便前端正确处理markdown
		processedResponse := strings.ReplaceAll(testChatResponse, "\n", "\\n")
		log.Printf("%s 测试模式处理后发送给前端的内容: %s", logPrefix, processedResponse)
		return processedResponse, nil
	}

	if s.apiKey == "" {
		log.Println(logPrefix + " 错误: API密钥未配置")
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

	// 精简日志，只记录关键信息
	log.Printf("%s 准备请求: 模型=%s, 消息数=%d",
		logPrefix, requestBody.Model, len(requestBody.Messages))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("%s 错误: 请求序列化失败: %v", logPrefix, err)
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 执行请求(带重试)
	responseBody, err := s.makeAPIRequest(jsonData, logPrefix)
	if err != nil {
		return "", err
	}

	// 解析JSON响应
	var responseData ChatResponse
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		log.Printf("%s 错误: 解析响应失败: %v", logPrefix, err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 验证响应
	if len(responseData.Choices) == 0 {
		log.Printf("%s 错误: API返回的选择项为空", logPrefix)
		return "", errors.New("API返回的选择项为空")
	}

	// 记录成功响应的统计信息 (精简统计信息)
	if responseData.Usage.TotalTokens > 0 {
		log.Printf("%s 响应统计: 总令牌=%d",
			logPrefix, responseData.Usage.TotalTokens)
	}

	// 返回AI回复内容
	originalContent := responseData.Choices[0].Message.Content
	log.Printf("%s 成功获取原始回复, 内容: %s", logPrefix, originalContent)

	// 将换行符替换为\n字符串，以便前端正确处理markdown
	processedContent := strings.ReplaceAll(originalContent, "\n", "\\n")
	log.Printf("%s 处理后发送给前端的内容: %s", logPrefix, processedContent)
	return processedContent, nil
}

// ChatWithAIStream 发送聊天请求到AI并以流式返回
func (s *AIService) ChatWithAIStream(messages []models.OpenAIMessage, out chan<- string) error {
	logPrefix := "[AI聊天-流式]"
	defer close(out)

	// 测试模式直接返回预定义响应
	if s.testMode {
		log.Printf("%s 测试模式，返回预定义响应", logPrefix)
		log.Printf("%s 测试模式原始完整回复内容: %s", logPrefix, testChatResponse)
		// 将换行符替换为\n字符串，以便前端正确处理markdown
		processedResponse := strings.ReplaceAll(testChatResponse, "\n", "\\n")
		log.Printf("%s 测试模式处理后发送给前端的内容: %s", logPrefix, processedResponse)
		// 模拟流式输出
		for _, char := range processedResponse {
			out <- string(char)
			time.Sleep(5 * time.Millisecond)
		}
		return nil
	}

	if s.apiKey == "" {
		log.Println(logPrefix + " 错误: API密钥未配置")
		return errors.New("AI API密钥未配置")
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
		Stream:      true,
	}

	log.Printf("%s 准备请求: 模型=%s, 消息数=%d",
		logPrefix, requestBody.Model, len(requestBody.Messages))

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("%s 错误: 请求序列化失败: %v", logPrefix, err)
		return fmt.Errorf("请求序列化失败: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("User-Agent", "OME-Nutrition-App/1.0")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	startTime := time.Now()
	log.Printf("%s 发送请求...", logPrefix)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求执行失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败: HTTP %d, %s", resp.StatusCode, truncateResponse(body, 200))
	}
	log.Printf("%s 收到响应: HTTP状态=%d, 准备接收流数据...", logPrefix, resp.StatusCode)

	var fullContent strings.Builder // 用于收集完整的流式响应内容
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := line[6:]
			if data == "[DONE]" {
				log.Printf("%s 流结束. 总耗时: %.2f秒", logPrefix, time.Since(startTime).Seconds())
				originalFullContent := fullContent.String()
				processedFullContent := strings.ReplaceAll(originalFullContent, "\n", "\\n")
				log.Printf("%s 原始完整回复内容: %s", logPrefix, originalFullContent)
				log.Printf("%s 处理后发送给前端的完整内容: %s", logPrefix, processedFullContent)
				return nil
			}

			var streamResp ChatStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				log.Printf("%s 解析流数据失败: %v, data: %s", logPrefix, err, data)
				continue
			}

			if len(streamResp.Choices) > 0 {
				content := streamResp.Choices[0].Delta.Content
				if content != "" {
					fullContent.WriteString(content) // 收集完整内容
					// 将换行符替换为\n字符串，以便前端正确处理markdown
					processedContent := strings.ReplaceAll(content, "\n", "\\n")
					out <- processedContent
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流失败: %w", err)
	}

	log.Printf("%s 流处理完成. 总耗时: %.2f秒", logPrefix, time.Since(startTime).Seconds())
	originalFullContent := fullContent.String()
	processedFullContent := strings.ReplaceAll(originalFullContent, "\n", "\\n")
	log.Printf("%s 原始完整回复内容: %s", logPrefix, originalFullContent)
	log.Printf("%s 处理后发送给前端的完整内容: %s", logPrefix, processedFullContent)
	return nil
}

// ConvertToMessages 将聊天消息转换为API请求格式
func (s *AIService) ConvertToMessages(chatMessages []models.ChatMessage) []models.OpenAIMessage {
	messages := make([]models.OpenAIMessage, len(chatMessages))
	for i, msg := range chatMessages {
		messages[i] = models.OpenAIMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	return messages
}

// GetSystemMessageForChat 获取聊天的系统消息
func (s *AIService) GetSystemMessageForChat() models.OpenAIMessage {
	return models.OpenAIMessage{
		Role: "system",
		Content: `你是一个专业的营养健康助手。你可以:
1. 提供健康饮食建议
2. 帮助用户了解食物的营养价值
3. 回答与健康、饮食相关的问题
4. 根据用户健康目标给出个性化建议

请注意以下要求：
- 使用简体中文回答问题
- 回复字数控制在150字以内，保持简洁
- 用友好、专业的方式回答问题
- 避免医疗诊断，只提供通用健康信息`,
	}
}

// GetSystemMessageForFoodRecognition 获取食物识别的系统消息
func (s *AIService) GetSystemMessageForFoodRecognition() models.OpenAIMessage {
	return models.OpenAIMessage{
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
    "calories_intake": 总热量,
    "protein_intake_g": 蛋白质克数,
    "carb_intake_g": 碳水克数,
    "fat_intake_g": 脂肪克数
  },
  "analysis": "对这顿饭的简短营养分析"
}

只返回JSON内容，不要添加其他文字说明。`,
	}
}

// AnalyzeImageWithAI 分析图片内容（使用base64编码）
func (s *AIService) AnalyzeImageWithAI(base64Image string, prompt string) (string, error) {
	logPrefix := "[AI图像分析]"

	// 测试模式直接返回预定义响应
	if s.testMode {
		log.Printf("%s 测试模式，返回预定义响应", logPrefix)
		log.Printf("%s 测试模式原始完整回复内容: %s", logPrefix, testFoodRecognitionResponse)
		// 将换行符替换为\n字符串，以便前端正确处理markdown
		processedResponse := strings.ReplaceAll(testFoodRecognitionResponse, "\n", "\\n")
		log.Printf("%s 测试模式处理后发送给前端的内容: %s", logPrefix, processedResponse)
		return processedResponse, nil
	}

	// 记录请求开始
	log.Printf("%s 开始请求, 提示词: %s", logPrefix, prompt)

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
		Model:       "gpt-4o-mini",
		Messages:    messages,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
	}

	// 精简日志
	log.Printf("%s 准备请求: 模型=%s", logPrefix, requestBody.Model)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("%s 错误: 请求序列化失败: %v", logPrefix, err)
		return "", fmt.Errorf("请求序列化失败: %v", err)
	}

	// 执行请求(带重试)
	responseBody, err := s.makeAPIRequest(jsonData, logPrefix)
	if err != nil {
		return "", err
	}

	// 解析JSON响应
	var responseData ChatResponse
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		log.Printf("%s 错误: 解析响应失败: %v", logPrefix, err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 验证响应
	if len(responseData.Choices) == 0 {
		log.Printf("%s 错误: API返回的选择项为空", logPrefix)
		return "", errors.New("API返回的选择项为空")
	}

	// 精简统计信息
	if responseData.Usage.TotalTokens > 0 {
		log.Printf("%s 响应统计: 总令牌=%d", logPrefix, responseData.Usage.TotalTokens)
	}

	// 返回AI回复内容
	originalContent := responseData.Choices[0].Message.Content
	log.Printf("%s 成功获取原始回复, 完整内容: %s", logPrefix, originalContent)

	// 将换行符替换为\n字符串，以便前端正确处理markdown
	processedContent := strings.ReplaceAll(originalContent, "\n", "\\n")
	log.Printf("%s 处理后发送给前端的内容: %s", logPrefix, processedContent)
	return processedContent, nil
}
