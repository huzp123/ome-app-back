package service

import (
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
func (s *ChatService) SendMessage(userID int64, sessionID string, content string) (*model.ChatMessage, *model.ChatMessage, error) {
	// 创建用户消息
	userMessage := &model.ChatMessage{
		SessionID: sessionID,
		UserID:    userID,
		Role:      model.RoleUser,
		Content:   content,
	}

	// 保存用户消息
	if err := s.chatDAO.AddMessage(userMessage); err != nil {
		return nil, nil, err
	}

	// 获取会话历史消息
	messages, err := s.chatDAO.GetLastNMessages(sessionID, 10) // 获取最近10条消息
	if err != nil {
		return userMessage, nil, err
	}

	// 转换为AI服务格式的消息
	aiMessages := s.aiService.ConvertToMessages(messages)

	// 添加系统消息
	systemMessage := s.aiService.GetSystemMessageForChat()
	aiMessages = append([]model.OpenAIMessage{systemMessage}, aiMessages...)

	// 发送请求给AI
	aiResponse, err := s.aiService.ChatWithAI(aiMessages)
	if err != nil {
		return userMessage, nil, err
	}

	// 创建AI回复消息
	assistantMessage := &model.ChatMessage{
		SessionID: sessionID,
		UserID:    userID,
		Role:      model.RoleAssistant,
		Content:   aiResponse,
	}

	// 保存AI回复
	if err := s.chatDAO.AddMessage(assistantMessage); err != nil {
		return userMessage, nil, err
	}

	return userMessage, assistantMessage, nil
}

// GetMessages 获取会话的消息列表
func (s *ChatService) GetMessages(sessionID string) ([]model.ChatMessage, error) {
	return s.chatDAO.GetMessages(sessionID)
}
