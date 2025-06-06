package service

import (
	"log"
	"strings"
	"sync"
	"time"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// ChatService 处理聊天相关服务
type ChatService struct {
	chatDAO   *dao.ChatDAO
	aiService *AIService
}

// NewChatService 创建聊天服务实例
func NewChatService(chatDAO *dao.ChatDAO, aiService *AIService) *ChatService {
	return &ChatService{
		chatDAO:   chatDAO,
		aiService: aiService,
	}
}

// CreateSession 创建新的聊天会话
func (s *ChatService) CreateSession(userID int64, title string) (*model.ChatSession, error) {
	return s.chatDAO.CreateSession(userID, title)
}

// GetSession 获取会话信息
func (s *ChatService) GetSession(sessionID string) (*model.ChatSession, error) {
	return s.chatDAO.GetSession(sessionID)
}

// ListUserSessions 获取用户的所有会话列表
func (s *ChatService) ListUserSessions(userID int64) ([]model.ChatSession, error) {
	return s.chatDAO.ListUserSessions(userID)
}

// UpdateSessionTitle 更新会话标题
func (s *ChatService) UpdateSessionTitle(sessionID string, title string) error {
	return s.chatDAO.UpdateSessionTitle(sessionID, title)
}

// DeleteSession 删除会话及其消息
func (s *ChatService) DeleteSession(sessionID string) error {
	return s.chatDAO.DeleteSession(sessionID)
}

// SendMessage 发送消息并获取AI回复
func (s *ChatService) SendMessage(userID int64, sessionID string, content string) (*model.ChatMessage, <-chan string, error) {
	log.Printf("[聊天] 用户(ID:%d)在会话(ID:%s)中发送新消息", userID, sessionID)
	start := time.Now()

	// 创建用户消息
	userMessage := &model.ChatMessage{
		SessionID: sessionID,
		UserID:    userID,
		Role:      model.RoleUser,
		Content:   content,
	}

	// 保存用户消息
	if err := s.chatDAO.AddMessage(userMessage); err != nil {
		log.Printf("[聊天] 错误: 保存用户消息失败: %v", err)
		return nil, nil, err
	}
	log.Printf("[聊天] 用户消息已保存(ID:%d)", userMessage.ID)

	// 创建一个通道，用于将AI的响应流式传输给调用者
	responseChan := make(chan string, 10)

	// 启动一个goroutine来处理AI交互和数据库保存
	go func() {
		defer close(responseChan)

		// 获取会话历史消息
		log.Printf("[聊天] 正在获取会话历史消息...")
		messages, err := s.chatDAO.GetLastNMessages(sessionID, 10) // 获取最近10条消息
		if err != nil {
			log.Printf("[聊天] 错误: 获取会话历史消息失败: %v", err)
			return
		}
		log.Printf("[聊天] 成功获取%d条历史消息", len(messages))

		// 转换为AI服务格式的消息
		aiMessages := s.aiService.ConvertToMessages(messages)

		// 添加系统消息
		systemMessage := s.aiService.GetSystemMessageForChat()
		aiMessages = append([]model.OpenAIMessage{systemMessage}, aiMessages...)
		log.Printf("[聊天] 准备向AI发送%d条消息(含系统消息)", len(aiMessages))

		// 用于从AI服务接收流式响应的内部通道
		aiStreamChan := make(chan string)
		var fullResponse strings.Builder
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			for chunk := range aiStreamChan {
				fullResponse.WriteString(chunk)
				responseChan <- chunk
			}
		}()

		// 调用流式AI服务
		log.Printf("[聊天] 正在调用AI流式服务...")
		err = s.aiService.ChatWithAIStream(aiMessages, aiStreamChan)
		if err != nil {
			log.Printf("[聊天] 错误: 从AI获取流式回复失败: %v", err)
		}

		wg.Wait()
		log.Printf("[聊天] 成功从AI获取流式回复，长度%d字符", fullResponse.Len())

		// 创建并保存AI回复消息
		if fullResponse.Len() > 0 {
			assistantMessage := &model.ChatMessage{
				SessionID: sessionID,
				UserID:    userID,
				Role:      model.RoleAssistant,
				Content:   fullResponse.String(),
			}

			if err := s.chatDAO.AddMessage(assistantMessage); err != nil {
				log.Printf("[聊天] 错误: 保存AI回复失败: %v", err)
			} else {
				log.Printf("[聊天] AI回复已保存(ID:%d)", assistantMessage.ID)
			}
		}

		duration := time.Since(start)
		log.Printf("[聊天] 消息处理完成, 总耗时: %.2f秒", duration.Seconds())
	}()

	return userMessage, responseChan, nil
}

// GetMessages 获取会话的消息列表
func (s *ChatService) GetMessages(sessionID string) ([]model.ChatMessage, error) {
	return s.chatDAO.GetMessages(sessionID)
}
