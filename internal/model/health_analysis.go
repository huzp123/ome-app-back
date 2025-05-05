package model

import (
	"time"
)

// HealthAnalysis 用户健康分析结果
type HealthAnalysis struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	UserID int64 `json:"user_id" gorm:"index;not null"`

	BMI  float64 `json:"bmi" gorm:"type:numeric(5,2)"`
	BMR  float64 `json:"bmr" gorm:"type:numeric(6,2)"`  // 基础代谢率
	TDEE float64 `json:"tdee" gorm:"type:numeric(6,2)"` // 每日总能量消耗

	ProteinNeedG float64 `json:"protein_need_g" gorm:"type:numeric(6,2)"` // 每日蛋白质需求(克)
	CarbNeedG    float64 `json:"carb_need_g" gorm:"type:numeric(6,2)"`    // 每日碳水需求(克)
	FatNeedG     float64 `json:"fat_need_g" gorm:"type:numeric(6,2)"`     // 每日脂肪需求(克)

	RecommendedCalories float64 `json:"recommended_calories" gorm:"type:numeric(6,2)"` // 推荐每日摄入热量

	AnalysisContent string `json:"analysis_content" gorm:"type:text"` // 分析结果文本内容

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (HealthAnalysis) TableName() string {
	return "health_analyses"
}
