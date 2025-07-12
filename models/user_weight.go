package models

import (
	"time"
)

// UserWeight 用户体重记录表
type UserWeight struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	UserID     int64     `json:"user_id" gorm:"index;not null"`
	WeightKG   float64   `json:"weight_kg" gorm:"type:numeric(5,2);not null"`
	RecordDate time.Time `json:"record_date" gorm:"type:date;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserWeight) TableName() string {
	return "user_weights"
}
