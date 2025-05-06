package model

import (
	"time"
)

// DailyNutrition 用户每日营养摄入记录
type DailyNutrition struct {
	ID     int64     `json:"id" gorm:"primaryKey"`
	UserID int64     `json:"user_id" gorm:"index;not null"`
	Date   time.Time `json:"date" gorm:"type:date;not null;index:idx_user_date,unique:true,priority:1"`

	// 实际摄入量
	CaloriesIntake float64 `json:"calories_intake" gorm:"type:numeric(6,2);default:0"` // 当日摄入总热量(千卡)
	ProteinIntakeG float64 `json:"protein_intake_g" gorm:"type:numeric(6,2);default:0"`
	CarbIntakeG    float64 `json:"carb_intake_g" gorm:"type:numeric(6,2);default:0"`
	FatIntakeG     float64 `json:"fat_intake_g" gorm:"type:numeric(6,2);default:0"`

	// 目标摄入量（根据健康分析生成）
	TargetCalories float64 `json:"target_calories" gorm:"type:numeric(6,2);default:0"`
	TargetProteinG float64 `json:"target_protein_g" gorm:"type:numeric(6,2);default:0"`
	TargetCarbG    float64 `json:"target_carb_g" gorm:"type:numeric(6,2);default:0"`
	TargetFatG     float64 `json:"target_fat_g" gorm:"type:numeric(6,2);default:0"`

	// 完成率与统计
	CaloriesCompletionRate float64 `json:"calories_completion_rate" gorm:"type:numeric(5,2)"` // 热量目标完成率(%)

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (DailyNutrition) TableName() string {
	return "daily_nutrition"
}
