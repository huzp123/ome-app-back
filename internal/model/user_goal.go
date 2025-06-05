package model

import (
	"time"
)

// UserGoal 用户健康目标与饮食偏好
type UserGoal struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"user_id" gorm:"index;not null"`

	GoalType         string    `json:"goal_type"        gorm:"type:varchar(16);not null"` // lose_fat / keep_fit / gain_muscle
	TargetWeightKG   float64   `json:"target_weight_kg" gorm:"type:decimal(5,2);not null"`
	WeeklyChangeKG   float64   `json:"weekly_change_kg" gorm:"type:decimal(4,2);not null"`
	TargetDate       time.Time `json:"target_date"      gorm:"type:date;not null"`
	DietType         string    `json:"diet_type"        gorm:"type:varchar(16);not null"` // normal / vegetarian / low_carb ...
	TastePreferences []string  `json:"taste_preferences" gorm:"serializer:json;not null"` // 口味偏好，必填
	FoodIntolerances []string  `json:"food_intolerances" gorm:"serializer:json;not null"` // 食物不耐受，必填

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserGoal) TableName() string {
	return "user_goals"
}
