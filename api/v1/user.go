package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ome-app-back/internal/service"
	"ome-app-back/pkg/errcode"
)

type UserAPI struct {
	userService *service.UserService
}

func NewUserAPI(userService *service.UserService) *UserAPI {
	return &UserAPI{userService: userService}
}

// Register 用户注册
func (api *UserAPI) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	resp, err := api.userService.Register(req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": resp,
	})
}

// Login 用户登录
func (api *UserAPI) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	resp, err := api.userService.Login(req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": resp,
	})
}

// UpdateProfile 更新用户档案
func (api *UserAPI) UpdateProfile(c *gin.Context) {
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	// 从JWT中获取用户ID（示例，实际项目中应从认证中间件获取）
	userID := getUserIDFromContext(c)
	req.UserID = userID

	err := api.userService.UpdateProfile(req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": nil,
	})
}

// UpdateGoal 更新用户健康目标
func (api *UserAPI) UpdateGoal(c *gin.Context) {
	var req service.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	// 从JWT中获取用户ID（示例，实际项目中应从认证中间件获取）
	userID := getUserIDFromContext(c)
	req.UserID = userID

	err := api.userService.UpdateGoal(req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": nil,
	})
}

// 从上下文中获取用户ID的辅助函数
func getUserIDFromContext(c *gin.Context) int64 {
	// 实际项目中，这个值应该由认证中间件设置
	// 这里仅作为示例实现
	if value, exists := c.Get("user_id"); exists {
		if userID, ok := value.(int64); ok {
			return userID
		}
	}
	return 0
}
