package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
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
func (d *ChatDAO) CreateSession(userID int64, title string) (*models.ChatSession, error) {
	session := models.ChatSession{
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
func (d *ChatDAO) GetSession(sessionID string) (*models.ChatSession, error) {
	var session models.ChatSession
	if err := d.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListUserSessions 获取用户的所有会话列表
func (d *ChatDAO) ListUserSessions(userID int64) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
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
	return d.db.Model(&models.ChatSession{}).
		Where("id = ?", sessionID).
		Update("title", title).Error
}

// DeleteSession 删除会话及其消息
func (d *ChatDAO) DeleteSession(sessionID string) error {
	// 在事务中删除会话和相关消息
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 先删除会话相关的所有消息
		if err := tx.Where("session_id = ?", sessionID).Delete(&models.ChatMessage{}).Error; err != nil {
			return err
		}
		// 删除会话
		if err := tx.Where("id = ?", sessionID).Delete(&models.ChatSession{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// AddMessage 添加新的聊天消息
func (d *ChatDAO) AddMessage(message *models.ChatMessage) error {
	return d.db.Create(message).Error
}

// GetMessages 获取会话的消息列表
func (d *ChatDAO) GetMessages(sessionID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := d.db.Where("session_id = ?", sessionID).
		Order("id ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// GetLastNMessages 获取会话的最近N条消息
func (d *ChatDAO) GetLastNMessages(sessionID string, n int) ([]models.ChatMessage, error) {
	var count int64
	if err := d.db.Model(&models.ChatMessage{}).
		Where("session_id = ?", sessionID).
		Count(&count).Error; err != nil {
		return nil, err
	}

	var messages []models.ChatMessage
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
	// 生成随机UUID
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		timestamp := time.Now().UnixNano()
		return fmt.Sprintf("sess_%x", timestamp)
	}

	// 设置UUID版本和变体标志位
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // 版本4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 变体RFC4122

	// 将当前时间戳融合到ID中，增加唯一性和安全性
	timestamp := time.Now().UnixNano()

	// 格式: sess_[时间戳前8位]_[UUID]
	return fmt.Sprintf("sess_%x_%s", timestamp%100000000, hex.EncodeToString(uuid))
}