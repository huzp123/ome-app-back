package models

import (
	"log"

	"gorm.io/gorm"
)

// Init 初始化数据库模型（自动迁移）
func Init(db *gorm.DB) error {
	err := db.AutoMigrate(
		&AppUser{},
		&UserGoal{},
		&UserWeight{},
		&UserHeight{},
		&HealthAnalysis{},
		&DailyNutrition{},
		&ChatSession{},
		&ChatMessage{},
		&FoodRecognition{},
		&UserExercise{},
		&MoodRecord{},
	)
	if err != nil {
		log.Printf("数据库自动迁移失败: %v", err)
		return err
	}
	log.Println("数据库自动迁移成功")
	return nil
}
