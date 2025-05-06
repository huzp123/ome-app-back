package model

import (
	"gorm.io/gorm"
)

// InitModels 初始化所有模型关系
func InitModels(db *gorm.DB) error {
	// 自动迁移表结构
	err := db.AutoMigrate(
		&AppUser{},
		&UserWeight{},
		&UserGoal{},
		&HealthAnalysis{},
		&DailyNutrition{},
	)
	if err != nil {
		return err
	}

	return nil
}
