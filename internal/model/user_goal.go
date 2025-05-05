package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// StringSlice 是字符串切片的自定义类型，实现了Scanner和Valuer接口
type StringSlice []string

// Scan 实现了sql.Scanner接口
func (ss *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言为[]byte失败")
	}

	return json.Unmarshal(bytes, &ss)
}

// Value 实现了driver.Valuer接口
func (ss StringSlice) Value() (driver.Value, error) {
	if len(ss) == 0 {
		return nil, nil
	}
	return json.Marshal(ss)
}

// UserGoal 用户健康目标与饮食偏好
type UserGoal struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"user_id" gorm:"index;not null"`

	GoalType         string      `json:"goal_type"        gorm:"type:varchar(16);not null"` // lose_fat / keep_fit / gain_muscle
	TargetWeightKG   float64     `json:"target_weight_kg" gorm:"type:decimal(5,2);not null"`
	WeeklyChangeKG   float64     `json:"weekly_change_kg" gorm:"type:decimal(4,2);not null"`
	TargetDate       time.Time   `json:"target_date"      gorm:"type:date;not null"`
	DietType         string      `json:"diet_type"        gorm:"type:varchar(16);not null"` // normal / vegetarian / low_carb ...
	TastePreferences StringSlice `json:"taste_preferences" gorm:"type:json"`
	FoodIntolerances StringSlice `json:"food_intolerances" gorm:"type:json"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (UserGoal) TableName() string {
	return "user_goals"
}
