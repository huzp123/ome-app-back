package v1

import (
	"ome-app-back/services"
)

// Init 初始化所有处理器
func Init(services *services.Services) *Handlers {
	return &Handlers{
		User:            NewUserAPI(services.UserService),
		HealthAnalysis:  NewHealthAnalysisAPI(services.HealthAnalysisService),
		Nutrition:       NewNutritionAPI(services.NutritionService),
		Chat:            NewChatAPI(services.ChatService),
		FoodRecognition: NewFoodRecognitionAPI(services.FoodRecognitionService),
		File:            NewFileAPI(services.FileService),
		Exercise:        NewExerciseAPI(services.ExerciseService),
		Mood:            NewMoodAPI(services.MoodService),
		Weight:          NewWeightAPI(services.WeightService),
		Height:          NewHeightAPI(services.HeightService),
	}
}
