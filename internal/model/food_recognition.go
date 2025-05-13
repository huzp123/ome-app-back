package model

import (
	"time"
)

// FoodRecognition 食物识别记录
type FoodRecognition struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	UserID          int64     `json:"user_id" gorm:"index;not null"`                            // 用户ID
	SessionID       string    `json:"session_id" gorm:"size:50;index:idx_session_recognition;"` // 相关会话ID
	ImageURL        string    `json:"image_url" gorm:"size:255;not null"`                       // 图片URL
	RecognizedFoods string    `json:"recognized_foods" gorm:"type:text"`                        // 识别出的食物列表(JSON格式)
	CaloriesIntake  float64   `json:"calories_intake" gorm:"type:numeric(6,2);default:0"`       // 估算热量(千卡)
	ProteinIntakeG  float64   `json:"protein_intake_g" gorm:"type:numeric(6,2);default:0"`      // 估算蛋白质(克)
	CarbIntakeG     float64   `json:"carb_intake_g" gorm:"type:numeric(6,2);default:0"`         // 估算碳水(克)
	FatIntakeG      float64   `json:"fat_intake_g" gorm:"type:numeric(6,2);default:0"`          // 估算脂肪(克)
	AIResponse      string    `json:"ai_response" gorm:"type:text"`                             // AI返回的完整响应
	IsAdopted       bool      `json:"is_adopted" gorm:"default:false"`                          // 用户是否采用此记录到营养摄入
	RecordDate      time.Time `json:"record_date" gorm:"type:date;not null"`                    // 记录日期
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 表名
func (FoodRecognition) TableName() string {
	return "food_recognitions"
}

// FoodRecognitionResult 食物识别结果(用于前端显示)
type FoodRecognitionResult struct {
	ID               int64                    `json:"id"`                // 识别记录ID
	ImageURL         string                   `json:"image_url"`         // 图片URL
	RecognizedFoods  []RecognizedFoodItem     `json:"recognized_foods"`  // 识别出的食物列表
	NutritionSummary FoodRecognitionNutrition `json:"nutrition_summary"` // 营养摘要
	AIAnalysis       string                   `json:"ai_analysis"`       // AI分析结果
	IsAdopted        bool                     `json:"is_adopted"`        // 是否已保存到营养摄入
	RecordDate       string                   `json:"record_date"`       // 记录日期
}

// RecognizedFoodItem 识别出的食物项
type RecognizedFoodItem struct {
	Name     string  `json:"name"`     // 食物名称
	Quantity string  `json:"quantity"` // 数量描述
	Calories float64 `json:"calories"` // 估算热量
}

// FoodRecognitionNutrition 食物识别的营养摘要
type FoodRecognitionNutrition struct {
	CaloriesIntake float64 `json:"calories_intake"`  // 总热量(千卡)
	ProteinIntakeG float64 `json:"protein_intake_g"` // 蛋白质(克)
	CarbIntakeG    float64 `json:"carb_intake_g"`    // 碳水(克)
	FatIntakeG     float64 `json:"fat_intake_g"`     // 脂肪(克)
}
