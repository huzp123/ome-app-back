package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ome-app-back/api"
	v1 "ome-app-back/api/v1"
	"ome-app-back/config"
	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
	"ome-app-back/internal/service"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.Println("启动服务...")

	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 检查配置有效性
	issues := cfg.CheckConfiguration()
	if len(issues) > 0 {
		log.Println("配置存在问题，请检查日志")
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 连接数据库
	db, err := setupDB(cfg.DB)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表结构
	autoMigrate(db)

	// 初始化服务组件
	appUserDAO := dao.NewAppUserDAO(db)
	userWeightDAO := dao.NewUserWeightDAO(db)
	userGoalDAO := dao.NewUserGoalDAO(db)
	healthAnalysisDAO := dao.NewHealthAnalysisDAO(db)
	dailyNutritionDAO := dao.NewDailyNutritionDAO(db)
	chatDAO := dao.NewChatDAO(db)
	foodRecognitionDAO := dao.NewFoodRecognitionDAO(db)

	userService := service.NewUserService(appUserDAO, userWeightDAO, userGoalDAO)
	healthAnalysisService := service.NewHealthAnalysisService(appUserDAO, userWeightDAO, userGoalDAO, healthAnalysisDAO)
	nutritionService := service.NewNutritionService(dailyNutritionDAO, healthAnalysisDAO)
	fileService := service.NewFileService(&cfg.Upload)
	aiService := service.NewAIService(&cfg.AI)
	chatService := service.NewChatService(chatDAO, aiService)
	foodRecognitionService := service.NewFoodRecognitionService(
		foodRecognitionDAO,
		dailyNutritionDAO,
		healthAnalysisDAO,
		fileService,
		aiService,
	)

	// testAIConnection(aiService)

	// 初始化API
	userAPI := v1.NewUserAPI(userService)
	healthAnalysisAPI := v1.NewHealthAnalysisAPI(healthAnalysisService)
	nutritionAPI := v1.NewNutritionAPI(nutritionService)
	chatAPI := v1.NewChatAPI(chatService)
	foodRecognitionAPI := v1.NewFoodRecognitionAPI(foodRecognitionService)
	fileAPI := v1.NewFileAPI(fileService)

	// 设置路由
	r := gin.Default()
	api.SetupRouter(r, userAPI, healthAnalysisAPI, nutritionAPI, chatAPI, foodRecognitionAPI, fileAPI)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 设置数据库连接
func setupDB(dbConfig config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch dbConfig.Type {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbConfig.GetDSN()), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbConfig.GetDSN()), &gorm.Config{})
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", dbConfig.Type)
	}

	if err != nil {
		return nil, err
	}

	// 获取底层SQL DB以设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	return db, nil
}

// 自动迁移数据库表结构
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.AppUser{},
		&model.UserGoal{},
		&model.UserWeight{},
		&model.HealthAnalysis{},
		&model.DailyNutrition{},
		&model.ChatSession{},
		&model.ChatMessage{},
		&model.FoodRecognition{},
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
}

// func testAIConnection(aiService *service.AIService) {
// 	// 简单的测试消息
// 	messages := []model.OpenAIMessage{
// 		{
// 			Role:    "system",
// 			Content: "你是一个健康助手，请简短地回答问题。",
// 		},
// 		{
// 			Role:    "user",
// 			Content: "你好，能听到我说话吗？请用一句话回复。",
// 		},
// 	}

// 	// 尝试与AI服务通信
// 	log.Println("测试AI服务连接...")
// 	_, err := aiService.ChatWithAI(messages)

// 	if err != nil {
// 		log.Printf("AI服务连接测试失败，请检查网络和API配置")
// 	} else {
// 		log.Printf("AI服务连接测试成功")
// 	}
// }
