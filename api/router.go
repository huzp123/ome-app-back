package api

import (
	"github.com/gin-gonic/gin"

	"ome-app-back/api/middleware"
	v1 "ome-app-back/api/v1"
)

// SetupRouter 设置API路由
func SetupRouter(engine *gin.Engine, userAPI *v1.UserAPI, healthAnalysisAPI *v1.HealthAnalysisAPI) {
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

		// 健康分析
		auth.GET("/health/analysis", healthAnalysisAPI.GenerateAnalysis)
		auth.GET("/health/history", healthAnalysisAPI.GetHistoryAnalysis)
	}
}
