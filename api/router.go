package api

import (
	"github.com/gin-gonic/gin"

	"ome-app-back/api/middleware"
	v1 "ome-app-back/api/v1"
)

// SetupRouter 设置API路由
func SetupRouter(engine *gin.Engine, userAPI *v1.UserAPI, healthAnalysisAPI *v1.HealthAnalysisAPI,
	nutritionAPI *v1.NutritionAPI, chatAPI *v1.ChatAPI, foodRecognitionAPI *v1.FoodRecognitionAPI,
	fileAPI *v1.FileAPI, exerciseAPI *v1.ExerciseAPI, moodAPI *v1.MoodAPI, weightAPI *v1.WeightAPI) {
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

		// 文件访问（无需权限验证的公共文件）
		apiV1.GET("/files/*filepath", fileAPI.GetFile)
	}

	// 需要认证的接口
	auth := apiV1.Group("")
	auth.Use(middleware.JWT())
	{
		// 用户信息与档案
		auth.GET("/user/info", userAPI.GetUserInfo)
		auth.PUT("/user/profile", userAPI.UpdateProfile)
		auth.PUT("/user/goal", userAPI.UpdateGoal)
		auth.GET("/user/goal", userAPI.GetGoal)

		// 文件访问（需要验证权限的用户文件）
		auth.GET("/user/files/*filepath", fileAPI.GetUserFile)

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
		auth.POST("/food/recognition/:id/save", foodRecognitionAPI.SaveRecognitionToNutrition)
		auth.GET("/food/recognition/adopted", foodRecognitionAPI.GetAdoptedRecognitions)

		// 运动记录
		auth.POST("/exercise", exerciseAPI.CreateExercise)
		auth.GET("/exercise/:id", exerciseAPI.GetExercise)
		auth.PUT("/exercise/:id", exerciseAPI.UpdateExercise)
		auth.DELETE("/exercise/:id", exerciseAPI.DeleteExercise)
		auth.GET("/exercise/history", exerciseAPI.GetExerciseHistory)
		auth.GET("/exercise/today", exerciseAPI.GetTodayExercises)
		auth.GET("/exercise/statistics", exerciseAPI.GetExerciseStatistics)
		auth.GET("/exercise/options", exerciseAPI.GetExerciseOptions)

		// 心情记录
		auth.POST("/mood", moodAPI.CreateMood)
		auth.GET("/mood/:id", moodAPI.GetMood)
		auth.DELETE("/mood/:id", moodAPI.DeleteMood)
		auth.GET("/mood/history", moodAPI.GetMoodHistory)
		auth.GET("/mood/today", moodAPI.GetTodayMoods)
		auth.GET("/mood/statistics", moodAPI.GetMoodStatistics)
		auth.GET("/mood/options", moodAPI.GetMoodOptions)

		// 体重管理
		auth.POST("/user/weight", weightAPI.CreateWeight)
		auth.GET("/user/weight/history", weightAPI.GetWeightHistory)
		auth.GET("/user/weight/current", weightAPI.GetCurrentWeight)
		auth.DELETE("/user/weight/:id", weightAPI.DeleteWeight)
		auth.GET("/user/weight/statistics", weightAPI.GetWeightStatistics)
	}
}
