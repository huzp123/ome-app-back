package v1

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"ome-app-back/internal/service"
)

// ChatAPI 处理聊天相关接口
type ChatAPI struct {
	chatService *service.ChatService
}

// NewChatAPI 创建聊天API处理实例
func NewChatAPI(chatService *service.ChatService) *ChatAPI {
	return &ChatAPI{
		chatService: chatService,
	}
}

// 会话创建请求
type CreateSessionRequest struct {
	Title string `json:"title"`
}

// CreateSession 创建新的聊天会话
func (a *ChatAPI) CreateSession(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 设置默认标题
	title := req.Title
	if title == "" {
		title = "新会话"
	}

	session, err := a.chatService.CreateSession(userID, title)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "创建会话失败")
		return
	}

	responseSuccess(c, session)
}

// GetSessions 获取用户的所有会话列表
func (a *ChatAPI) GetSessions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	sessions, err := a.chatService.ListUserSessions(userID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取会话列表失败")
		return
	}

	responseSuccess(c, sessions)
}

// 更新会话标题请求
type UpdateSessionTitleRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateSessionTitle 更新会话标题
func (a *ChatAPI) UpdateSessionTitle(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		responseError(c, http.StatusBadRequest, "会话ID不能为空")
		return
	}

	var req UpdateSessionTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := a.chatService.UpdateSessionTitle(sessionID, req.Title); err != nil {
		responseError(c, http.StatusInternalServerError, "更新会话标题失败")
		return
	}

	responseSuccess(c, nil)
}

// DeleteSession 删除聊天会话
func (a *ChatAPI) DeleteSession(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		responseError(c, http.StatusBadRequest, "会话ID不能为空")
		return
	}

	if err := a.chatService.DeleteSession(sessionID); err != nil {
		responseError(c, http.StatusInternalServerError, "删除会话失败")
		return
	}

	responseSuccess(c, nil)
}

// GetMessages 获取会话的消息列表
func (a *ChatAPI) GetMessages(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		responseError(c, http.StatusBadRequest, "会话ID不能为空")
		return
	}

	messages, err := a.chatService.GetMessages(sessionID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取消息列表失败")
		return
	}

	responseSuccess(c, messages)
}

// 发送消息请求
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// SendMessage 发送消息
func (a *ChatAPI) SendMessage(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		responseError(c, http.StatusBadRequest, "会话ID不能为空")
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	_, responseChan, err := a.chatService.SendMessage(userID, sessionID, req.Content)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "发送消息失败")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-responseChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
