package services

import (
	"ome-app-back/config"
	"ome-app-back/repositories"
)

// Services 包含所有的业务服务
type Services struct {
	UserService            *UserService
	HealthAnalysisService  *HealthAnalysisService
	NutritionService       *NutritionService
	FileService            *FileService
	AIService              *AIService
	ChatService            *ChatService
	FoodRecognitionService *FoodRecognitionService
	ExerciseService        *ExerciseService
	MoodService            *MoodService
	WeightService          *WeightService
	HeightService          *HeightService
}

// Init 初始化所有业务服务
func Init(repos *repositories.Repositories, cfg *config.Config) *Services {
	// 初始化基础服务
	fileService := NewFileService(&cfg.Upload)
	aiService := NewAIService(&cfg.AI)

	// 初始化业务服务
	userService := NewUserService(repos.AppUserDAO, repos.UserWeightDAO, repos.UserGoalDAO)
	healthAnalysisService := NewHealthAnalysisService(repos.AppUserDAO, repos.UserWeightDAO, repos.UserHeightDAO, repos.UserGoalDAO, repos.HealthAnalysisDAO)
	nutritionService := NewNutritionService(repos.DailyNutritionDAO, repos.HealthAnalysisDAO)
	chatService := NewChatService(repos.ChatDAO, aiService)
	foodRecognitionService := NewFoodRecognitionService(
		repos.FoodRecognitionDAO,
		repos.DailyNutritionDAO,
		repos.HealthAnalysisDAO,
		fileService,
		aiService,
	)
	exerciseService := NewExerciseService(repos.UserExerciseDAO)
	moodService := NewMoodService(repos.MoodRecordDAO)
	weightService := NewWeightService(repos.UserWeightDAO)
	heightService := NewHeightService(repos.UserHeightDAO)

	return &Services{
		UserService:            userService,
		HealthAnalysisService:  healthAnalysisService,
		NutritionService:       nutritionService,
		FileService:            fileService,
		AIService:              aiService,
		ChatService:            chatService,
		FoodRecognitionService: foodRecognitionService,
		ExerciseService:        exerciseService,
		MoodService:            moodService,
		WeightService:          weightService,
		HeightService:          heightService,
	}
}
