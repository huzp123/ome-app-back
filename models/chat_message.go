package models

import (
	"time"
)

// ChatRole 对话角色类型
type ChatRole string

const (
	RoleUser      ChatRole = "user"      // 用户
	RoleAssistant ChatRole = "assistant" // AI助手
	RoleSystem    ChatRole = "system"    // 系统消息
)

// ChatMessage 聊天消息
type ChatMessage struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	SessionID string    `json:"session_id" gorm:"size:50;index:idx_session;not null"` // 会话ID
	UserID    int64     `json:"user_id" gorm:"index;not null"`                        // 用户ID
	Role      ChatRole  `json:"role" gorm:"type:varchar(10);not null"`                // 角色
	Content   string    `json:"content" gorm:"type:text;not null"`                    // 消息内容
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// ChatSession 聊天会话信息
type ChatSession struct {
	ID        string    `json:"id" gorm:"primaryKey;size:50"`
	UserID    int64     `json:"user_id" gorm:"index;not null"` // 用户ID
	Title     string    `json:"title" gorm:"size:100"`         // 会话标题
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// OpenAIMessage OpenAI消息格式
type OpenAIMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}
