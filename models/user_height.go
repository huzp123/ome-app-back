package models

import (
	"time"
)

// UserHeight 用户身高记录表
type UserHeight struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	UserID     int64     `json:"user_id" gorm:"index;not null"`
	HeightCM   float64   `json:"height_cm" gorm:"type:decimal(5,2);not null;check:height_cm BETWEEN 50 AND 300"`
	RecordDate time.Time `json:"record_date" gorm:"type:date;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserHeight) TableName() string {
	return "user_heights"
}
