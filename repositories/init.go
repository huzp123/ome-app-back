package repositories

import (
	"gorm.io/gorm"
)

// Repositories 包含所有的数据访问对象
type Repositories struct {
	AppUserDAO         *AppUserDAO
	UserWeightDAO      *UserWeightDAO
	UserGoalDAO        *UserGoalDAO
	HealthAnalysisDAO  *HealthAnalysisDAO
	DailyNutritionDAO  *DailyNutritionDAO
	ChatDAO            *ChatDAO
	FoodRecognitionDAO *FoodRecognitionDAO
	UserExerciseDAO    *UserExerciseDAO
	MoodRecordDAO      *MoodRecordDAO
}

// Init 初始化所有数据访问对象
func Init(db *gorm.DB) *Repositories {
	return &Repositories{
		AppUserDAO:         NewAppUserDAO(db),
		UserWeightDAO:      NewUserWeightDAO(db),
		UserGoalDAO:        NewUserGoalDAO(db),
		HealthAnalysisDAO:  NewHealthAnalysisDAO(db),
		DailyNutritionDAO:  NewDailyNutritionDAO(db),
		ChatDAO:            NewChatDAO(db),
		FoodRecognitionDAO: NewFoodRecognitionDAO(db),
		UserExerciseDAO:    NewUserExerciseDAO(db),
		MoodRecordDAO:      NewMoodRecordDAO(db),
	}
}
