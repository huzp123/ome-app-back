package v1

import (
	"fmt"
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

	// 记录注册请求
	clientIP := c.ClientIP()
	fmt.Printf("[注册请求] IP: %s, 手机号: %s, 邮箱: %s, 用户名: %s\n",
		clientIP, req.Phone, req.Email, req.UserName)

	resp, err := api.userService.Register(req)
	if err != nil {
		// 记录注册失败
		fmt.Printf("[注册失败] IP: %s, 手机号: %s, 邮箱: %s, 错误: %s\n",
			clientIP, req.Phone, req.Email, err.Error())
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	// 记录注册成功
	fmt.Printf("[注册成功] 用户ID: %d, IP: %s, 手机号: %s, 邮箱: %s\n",
		resp.UserID, clientIP, req.Phone, req.Email)

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

// WechatLogin 微信登录
func (api *UserAPI) WechatLogin(c *gin.Context) {
	var req service.WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	// 记录微信登录请求
	clientIP := c.ClientIP()
	fmt.Printf("[微信登录请求] IP: %s, OpenID: %s, 用户名: %s\n",
		clientIP, req.OpenID, req.UserName)

	resp, err := api.userService.WechatLogin(req)
	if err != nil {
		// 记录微信登录失败
		fmt.Printf("[微信登录失败] IP: %s, OpenID: %s, 错误: %s\n",
			clientIP, req.OpenID, err.Error())
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	// 记录微信登录成功
	fmt.Printf("[微信登录成功] 用户ID: %d, IP: %s, OpenID: %s, 是否新用户: %t\n",
		resp.UserID, clientIP, req.OpenID, resp.IsNewUser)

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

// GetGoal 获取用户健康目标
func (api *UserAPI) GetGoal(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	// 调用服务层获取用户目标
	goal, err := api.userService.GetGoal(userID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	// 如果goal为nil，表示用户还没有设置健康目标
	var data interface{}
	if goal == nil {
		data = nil
	} else {
		data = goal
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": data,
	})
}

// GetUserInfo 获取用户信息
func (api *UserAPI) GetUserInfo(c *gin.Context) {
	// 从JWT中获取用户ID
	userID := getUserIDFromContext(c)
	if userID == 0 {
		fmt.Printf("[获取用户信息] 认证失败, IP: %s\n", c.ClientIP())
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	// 记录请求开始
	fmt.Printf("[获取用户信息] 开始处理请求, 用户ID: %d, IP: %s\n", userID, c.ClientIP())

	// 调用服务层获取用户信息
	userInfo, err := api.userService.GetUserInfo(userID)
	if err != nil {
		fmt.Printf("[获取用户信息] 失败, 用户ID: %d, 错误: %s\n", userID, err.Error())
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	// 记录成功，打印出参
	fmt.Printf("[获取用户信息] 成功, 用户ID: %d, 出参: %+v\n", userID, userInfo)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": userInfo,
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
