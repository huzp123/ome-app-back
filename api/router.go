package api

import (
	"github.com/gin-gonic/gin"

	"ome-app-back/api/middleware"
	v1 "ome-app-back/api/v1"
)

// SetupRouter 设置API路由
func SetupRouter(engine *gin.Engine, userAPI *v1.UserAPI, healthAnalysisAPI *v1.HealthAnalysisAPI,
	nutritionAPI *v1.NutritionAPI, chatAPI *v1.ChatAPI, foodRecognitionAPI *v1.FoodRecognitionAPI) {
	// 全局中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.Cors())

	// API版本前缀
	apiV1 := engine.Group("/api/v1")

	// 无需认证的公共接口
	{
		// 健康检查
		apiV1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "成功",
				"data": gin.H{"message": "pong"},
			})
		})

		// 用户注册登录
		apiV1.POST("/register", userAPI.Register)
		apiV1.POST("/login", userAPI.Login)
	}

	// 需要认证的接口
	auth := apiV1.Group("")
	auth.Use(middleware.JWT())
	{
		// 用户信息与档案
		auth.PUT("/user/profile", userAPI.UpdateProfile)
		auth.PUT("/user/goal", userAPI.UpdateGoal)
		auth.GET("/user/goal", userAPI.GetGoal)

		// 健康分析
		auth.GET("/health/analysis", healthAnalysisAPI.GenerateAnalysis)
		auth.GET("/health/history", healthAnalysisAPI.GetHistoryAnalysis)

		// 每日营养
		auth.GET("/nutrition/today", nutritionAPI.GetTodayNutrition)
		auth.PUT("/nutrition/today", nutritionAPI.UpdateTodayNutrition)
		auth.GET("/nutrition/history", nutritionAPI.GetNutritionHistory)
		auth.GET("/nutrition/weekly-summary", nutritionAPI.GetWeekSummary)

		// 聊天会话管理
		auth.POST("/chat/sessions", chatAPI.CreateSession)
		auth.GET("/chat/sessions", chatAPI.GetSessions)
		auth.PUT("/chat/sessions/:session_id", chatAPI.UpdateSessionTitle)
		auth.DELETE("/chat/sessions/:session_id", chatAPI.DeleteSession)
		auth.GET("/chat/sessions/:session_id/messages", chatAPI.GetMessages)
		auth.POST("/chat/sessions/:session_id/messages", chatAPI.SendMessage)

		// 食物识别
		auth.POST("/food/recognize", foodRecognitionAPI.RecognizeFood)
		auth.GET("/food/recognition/:id", foodRecognitionAPI.GetRecognitionByID)
		auth.GET("/food/recognition/today", foodRecognitionAPI.GetTodayRecognitions)
	}
}
