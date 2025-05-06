package dao

import (
	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// ChatDAO 处理聊天相关的数据访问
type ChatDAO struct {
	db *gorm.DB
}

// NewChatDAO 创建聊天DAO实例
func NewChatDAO(db *gorm.DB) *ChatDAO {
	return &ChatDAO{db: db}
}

// CreateSession 创建新的聊天会话
func (d *ChatDAO) CreateSession(userID int64, title string) (*model.ChatSession, error) {
	session := model.ChatSession{
		ID:     generateSessionID(),
		UserID: userID,
		Title:  title,
	}

	if err := d.db.Create(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSession 获取会话信息
func (d *ChatDAO) GetSession(sessionID string) (*model.ChatSession, error) {
	var session model.ChatSession
	if err := d.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListUserSessions 获取用户的所有会话列表
func (d *ChatDAO) ListUserSessions(userID int64) ([]model.ChatSession, error) {
	var sessions []model.ChatSession
	err := d.db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// UpdateSessionTitle 更新会话标题
func (d *ChatDAO) UpdateSessionTitle(sessionID string, title string) error {
	return d.db.Model(&model.ChatSession{}).
		Where("id = ?", sessionID).
		Update("title", title).Error
}

// DeleteSession 删除会话及其消息
func (d *ChatDAO) DeleteSession(sessionID string) error {
	// 在事务中删除会话和相关消息
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 先删除会话相关的所有消息
		if err := tx.Where("session_id = ?", sessionID).Delete(&model.ChatMessage{}).Error; err != nil {
			return err
		}
		// 删除会话
		if err := tx.Where("id = ?", sessionID).Delete(&model.ChatSession{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// AddMessage 添加新的聊天消息
func (d *ChatDAO) AddMessage(message *model.ChatMessage) error {
	return d.db.Create(message).Error
}

// GetMessages 获取会话的消息列表
func (d *ChatDAO) GetMessages(sessionID string) ([]model.ChatMessage, error) {
	var messages []model.ChatMessage
	err := d.db.Where("session_id = ?", sessionID).
		Order("id ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetLastNMessages 获取会话的最近N条消息
func (d *ChatDAO) GetLastNMessages(sessionID string, n int) ([]model.ChatMessage, error) {
	var count int64
	if err := d.db.Model(&model.ChatMessage{}).
		Where("session_id = ?", sessionID).
		Count(&count).Error; err != nil {
		return nil, err
	}

	var messages []model.ChatMessage
	err := d.db.Where("session_id = ?", sessionID).
		Order("id ASC").
		Limit(n).
		Offset(int(count) - n).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// 生成唯一的会话ID
func generateSessionID() string {
	// 实际应用中应使用更安全的方法生成UUID
	return "sess_" + randomString(16)
}

// 生成随机字符串
func randomString(length int) string {
	// 简化实现，实际应使用crypto/rand
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
