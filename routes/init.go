package routes

import (
	v1 "ome-app-back/handlers/v1"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(engine *gin.Engine, handlers *v1.Handlers) {
	SetupRoutes(engine, handlers)
}
