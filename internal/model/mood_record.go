package model

import (
	"time"
)

// MoodRecord 情绪记录表
type MoodRecord struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	UserID      int64     `json:"user_id" gorm:"index;not null"`
	TimeContext string    `json:"time_context" gorm:"size:8;not null"`     // "now" 或 "today"
	MoodLevel   int       `json:"mood_level" gorm:"type:tinyint;not null"` // 1-7级
	MoodTags    []string  `json:"mood_tags" gorm:"type:text"`              // 情绪标签数组，可选
	Influences  []string  `json:"influences" gorm:"type:text"`             // 影响因素数组，可选
	RecordTime  time.Time `json:"record_time" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (MoodRecord) TableName() string {
	return "mood_records"
}
