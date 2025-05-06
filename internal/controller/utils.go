package controller

import (
	"github.com/gin-gonic/gin"
)

// 常用上下文键
const (
	UserIDKey = "user_id" // 用户ID键名
)

// getUserIDFromContext 从上下文中获取当前用户ID
func getUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}

	id, ok := userID.(int64)
	if !ok {
		return 0
	}

	return id
}
