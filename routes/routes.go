package routes

import (
	v1 "ome-app-back/handlers/v1"
	"ome-app-back/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(engine *gin.Engine, handlers *v1.Handlers) {
	// 全局中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.Cors())

	// API版本前缀
	apiV1 := engine.Group("/api/v1")

	// 无需认证的公共接口
	setupPublicRoutes(apiV1, handlers)

	// 需要认证的接口
	auth := apiV1.Group("")
	auth.Use(middleware.JWT())
	setupAuthRoutes(auth, handlers)
}

// setupPublicRoutes 设置公共路由
func setupPublicRoutes(router *gin.RouterGroup, handlers *v1.Handlers) {
	// 健康检查
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "成功",
			"data": gin.H{"message": "pong"},
		})
	})

	// 用户注册登录
	router.POST("/register", handlers.User.Register)
	router.POST("/login", handlers.User.Login)
	router.POST("/wechat/login", handlers.User.WechatLogin)

	// 文件访问（无需权限验证的公共文件）
	router.GET("/files/*filepath", handlers.File.GetFile)
}

// setupAuthRoutes 设置需要认证的路由
func setupAuthRoutes(router *gin.RouterGroup, handlers *v1.Handlers) {
	// 用户信息与档案
	router.GET("/user/info", handlers.User.GetUserInfo)
	router.PUT("/user/profile", handlers.User.UpdateProfile)
	router.PUT("/user/goal", handlers.User.UpdateGoal)
	router.GET("/user/goal", handlers.User.GetGoal)

	// 文件访问（需要验证权限的用户文件）
	router.GET("/user/files/*filepath", handlers.File.GetUserFile)

	// 健康分析
	router.GET("/health/analysis", handlers.HealthAnalysis.GenerateAnalysis)
	router.GET("/health/history", handlers.HealthAnalysis.GetHistoryAnalysis)

	// 每日营养
	router.GET("/nutrition/today", handlers.Nutrition.GetTodayNutrition)
	router.PUT("/nutrition/today", handlers.Nutrition.UpdateTodayNutrition)
	router.GET("/nutrition/history", handlers.Nutrition.GetNutritionHistory)
	router.GET("/nutrition/weekly-summary", handlers.Nutrition.GetWeekSummary)

	// 聊天会话管理
	router.POST("/chat/sessions", handlers.Chat.CreateSession)
	router.GET("/chat/sessions", handlers.Chat.GetSessions)
	router.PUT("/chat/sessions/:session_id", handlers.Chat.UpdateSessionTitle)
	router.DELETE("/chat/sessions/:session_id", handlers.Chat.DeleteSession)
	router.GET("/chat/sessions/:session_id/messages", handlers.Chat.GetMessages)
	router.POST("/chat/sessions/:session_id/messages", handlers.Chat.SendMessage)

	// 食物识别
	router.POST("/food/recognize", handlers.FoodRecognition.RecognizeFood)
	router.GET("/food/recognition/:id", handlers.FoodRecognition.GetRecognitionByID)
	router.GET("/food/recognition/today", handlers.FoodRecognition.GetTodayRecognitions)
	router.POST("/food/recognition/:id/save", handlers.FoodRecognition.SaveRecognitionToNutrition)
	router.GET("/food/recognition/adopted", handlers.FoodRecognition.GetAdoptedRecognitions)

	// 运动记录
	router.POST("/exercise", handlers.Exercise.CreateExercise)
	router.GET("/exercise/:id", handlers.Exercise.GetExercise)
	router.PUT("/exercise/:id", handlers.Exercise.UpdateExercise)
	router.DELETE("/exercise/:id", handlers.Exercise.DeleteExercise)
	router.GET("/exercise/history", handlers.Exercise.GetExerciseHistory)
	router.GET("/exercise/today", handlers.Exercise.GetTodayExercises)
	router.GET("/exercise/statistics", handlers.Exercise.GetExerciseStatistics)
	router.GET("/exercise/options", handlers.Exercise.GetExerciseOptions)

	// 心情记录
	router.POST("/mood", handlers.Mood.CreateMood)
	router.GET("/mood/:id", handlers.Mood.GetMood)
	router.DELETE("/mood/:id", handlers.Mood.DeleteMood)
	router.GET("/mood/history", handlers.Mood.GetMoodHistory)
	router.GET("/mood/today", handlers.Mood.GetTodayMoods)
	router.GET("/mood/statistics", handlers.Mood.GetMoodStatistics)
	router.GET("/mood/options", handlers.Mood.GetMoodOptions)

	// 体重管理
	router.POST("/user/weight", handlers.Weight.CreateWeight)
	router.GET("/user/weight/history", handlers.Weight.GetWeightHistory)
	router.GET("/user/weight/current", handlers.Weight.GetCurrentWeight)
	router.DELETE("/user/weight/:id", handlers.Weight.DeleteWeight)
	router.GET("/user/weight/statistics", handlers.Weight.GetWeightStatistics)

	// 身高管理
	router.POST("/user/height", handlers.Height.CreateHeight)
	router.GET("/user/height/history", handlers.Height.GetHeightHistory)
	router.GET("/user/height/current", handlers.Height.GetCurrentHeight)
	router.DELETE("/user/height/:id", handlers.Height.DeleteHeight)
	router.GET("/user/height/statistics", handlers.Height.GetHeightStatistics)
}
