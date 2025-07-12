package v1

import (
	"ome-app-back/services"
)

// Handlers 包含所有的处理器
type Handlers struct {
	User            *UserAPI
	HealthAnalysis  *HealthAnalysisAPI
	Nutrition       *NutritionAPI
	Chat            *ChatAPI
	FoodRecognition *FoodRecognitionAPI
	File            *FileAPI
	Exercise        *ExerciseAPI
	Mood            *MoodAPI
	Weight          *WeightAPI
	Height          *HeightAPI
}

// NewHandlers 创建新的Handlers实例
func NewHandlers(
	userService *services.UserService,
	healthAnalysisService *services.HealthAnalysisService,
	nutritionService *services.NutritionService,
	chatService *services.ChatService,
	foodRecognitionService *services.FoodRecognitionService,
	fileService *services.FileService,
	exerciseService *services.ExerciseService,
	moodService *services.MoodService,
	weightService *services.WeightService,
	heightService *services.HeightService,
) *Handlers {
	return &Handlers{
		User:            NewUserAPI(userService),
		HealthAnalysis:  NewHealthAnalysisAPI(healthAnalysisService),
		Nutrition:       NewNutritionAPI(nutritionService),
		Chat:            NewChatAPI(chatService),
		FoodRecognition: NewFoodRecognitionAPI(foodRecognitionService),
		File:            NewFileAPI(fileService),
		Exercise:        NewExerciseAPI(exerciseService),
		Mood:            NewMoodAPI(moodService),
		Weight:          NewWeightAPI(weightService),
		Height:          NewHeightAPI(heightService),
	}
}
