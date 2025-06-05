package model

import (
	"time"
)

// UserExercise 用户运动记录表
type UserExercise struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	UserID         int64     `json:"user_id" gorm:"index;not null"`
	ExerciseType   string    `json:"exercise_type" gorm:"size:32;not null"`
	DurationMin    float64   `json:"duration_min" gorm:"type:numeric(6,2);not null"`
	CaloriesBurned float64   `json:"calories_burned" gorm:"type:numeric(6,2);not null"`
	DistanceKM     *float64  `json:"distance_km,omitempty" gorm:"type:numeric(6,2)"`
	StartTime      time.Time `json:"start_time" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserExercise) TableName() string {
	return "user_exercises"
}
