package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"ome-app-back/config"
	"ome-app-back/database"
	v1 "ome-app-back/handlers/v1"
	"ome-app-back/repositories"
	"ome-app-back/routes"
	"ome-app-back/services"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.Println("启动服务...")

	// 初始化配置
	cfg, err := config.Init("config/config.yaml")
	if err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := database.Init(cfg.DB)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化数据访问层
	repos := repositories.Init(db)

	// 初始化服务层
	services := services.Init(repos, cfg)

	// 初始化处理器
	handlers := v1.Init(services)

	// 初始化路由
	r := gin.Default()
	routes.Init(r, handlers)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
