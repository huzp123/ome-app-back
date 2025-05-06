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
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
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

	// 初始化DAO
	appUserDAO := dao.NewAppUserDAO(db)
	userWeightDAO := dao.NewUserWeightDAO(db)
	userGoalDAO := dao.NewUserGoalDAO(db)
	healthAnalysisDAO := dao.NewHealthAnalysisDAO(db)
	dailyNutritionDAO := dao.NewDailyNutritionDAO(db)
	chatDAO := dao.NewChatDAO(db)
	foodRecognitionDAO := dao.NewFoodRecognitionDAO(db)

	// 初始化Service
	userService := service.NewUserService(appUserDAO, userWeightDAO, userGoalDAO)
	healthAnalysisService := service.NewHealthAnalysisService(appUserDAO, userWeightDAO, userGoalDAO, healthAnalysisDAO)
	nutritionService := service.NewNutritionService(dailyNutritionDAO)
	fileService := service.NewFileService(&cfg.Upload)
	aiService := service.NewAIService(&cfg.AI)
	chatService := service.NewChatService(chatDAO, aiService)
	foodRecognitionService := service.NewFoodRecognitionService(
		foodRecognitionDAO,
		dailyNutritionDAO,
		fileService,
		aiService,
	)

	// 初始化API
	userAPI := v1.NewUserAPI(userService)
	healthAnalysisAPI := v1.NewHealthAnalysisAPI(healthAnalysisService)
	nutritionAPI := v1.NewNutritionAPI(nutritionService)
	chatAPI := v1.NewChatAPI(chatService)
	foodRecognitionAPI := v1.NewFoodRecognitionAPI(foodRecognitionService)

	// 设置路由
	r := gin.Default()
	api.SetupRouter(r, userAPI, healthAnalysisAPI, nutritionAPI, chatAPI, foodRecognitionAPI)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("服务器启动在 %s\n", serverAddr)
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
	log.Println("数据库表结构已自动迁移")
}
