package dao

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// FoodRecognitionDAO 处理食物识别相关的数据访问
type FoodRecognitionDAO struct {
	db *gorm.DB
}

// NewFoodRecognitionDAO 创建食物识别DAO实例
func NewFoodRecognitionDAO(db *gorm.DB) *FoodRecognitionDAO {
	return &FoodRecognitionDAO{db: db}
}

// CreateRecognition 创建食物识别记录
func (d *FoodRecognitionDAO) CreateRecognition(userID int64, sessionID, imageURL string, foods []model.RecognizedFoodItem, nutrition model.FoodRecognitionNutrition, aiResponse string) (*model.FoodRecognition, error) {
	// 将食物列表转为JSON字符串存储
	foodsJSON, err := json.Marshal(foods)
	if err != nil {
		return nil, err
	}

	recognition := model.FoodRecognition{
		UserID:          userID,
		SessionID:       sessionID,
		ImageURL:        imageURL,
		RecognizedFoods: string(foodsJSON),
		CaloriesIntake:  nutrition.CaloriesIntake,
		ProteinIntakeG:  nutrition.ProteinIntakeG,
		CarbIntakeG:     nutrition.CarbIntakeG,
		FatIntakeG:      nutrition.FatIntakeG,
		AIResponse:      aiResponse,
		RecordDate:      time.Now().Truncate(24 * time.Hour), // 当天日期，去除时分秒
	}

	if err := d.db.Create(&recognition).Error; err != nil {
		return nil, err
	}

	return &recognition, nil
}

// GetRecognitionByID 通过ID获取食物识别记录
func (d *FoodRecognitionDAO) GetRecognitionByID(id int64) (*model.FoodRecognition, error) {
	var recognition model.FoodRecognition
	if err := d.db.First(&recognition, id).Error; err != nil {
		return nil, err
	}
	return &recognition, nil
}

// GetUserRecognitionsByDate 获取用户某天的所有食物识别记录
func (d *FoodRecognitionDAO) GetUserRecognitionsByDate(userID int64, date time.Time) ([]model.FoodRecognition, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nextDay := dateOnly.AddDate(0, 0, 1)

	var recognitions []model.FoodRecognition
	err := d.db.Where("user_id = ? AND record_date >= ? AND record_date < ?",
		userID, dateOnly, nextDay).
		Order("created_at DESC").
		Find(&recognitions).Error
	if err != nil {
		return nil, err
	}
	return recognitions, nil
}

// GetUserTodayRecognitions 获取用户今天的所有食物识别记录
func (d *FoodRecognitionDAO) GetUserTodayRecognitions(userID int64) ([]model.FoodRecognition, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return d.GetUserRecognitionsByDate(userID, today)
}

// GetUserRecentRecognitions 获取用户最近的食物识别记录
func (d *FoodRecognitionDAO) GetUserRecentRecognitions(userID int64, limit int) ([]model.FoodRecognition, error) {
	var recognitions []model.FoodRecognition
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&recognitions).Error
	if err != nil {
		return nil, err
	}
	return recognitions, nil
}

// ConvertToResult 将数据库记录转换为前端可用的结果对象
func (d *FoodRecognitionDAO) ConvertToResult(recognition *model.FoodRecognition) (*model.FoodRecognitionResult, error) {
	var foods []model.RecognizedFoodItem
	if err := json.Unmarshal([]byte(recognition.RecognizedFoods), &foods); err != nil {
		return nil, err
	}

	result := &model.FoodRecognitionResult{
		ImageURL:        recognition.ImageURL,
		RecognizedFoods: foods,
		NutritionSummary: model.FoodRecognitionNutrition{
			CaloriesIntake: recognition.CaloriesIntake,
			ProteinIntakeG: recognition.ProteinIntakeG,
			CarbIntakeG:    recognition.CarbIntakeG,
			FatIntakeG:     recognition.FatIntakeG,
		},
		AIAnalysis: recognition.AIResponse,
	}

	return result, nil
}

// SummarizeTodayNutrition 汇总用户今日所有食物识别的营养摄入
func (d *FoodRecognitionDAO) SummarizeTodayNutrition(userID int64) (model.FoodRecognitionNutrition, error) {
	records, err := d.GetUserTodayRecognitions(userID)
	if err != nil {
		return model.FoodRecognitionNutrition{}, err
	}

	var summary model.FoodRecognitionNutrition
	for _, record := range records {
		summary.CaloriesIntake += record.CaloriesIntake
		summary.ProteinIntakeG += record.ProteinIntakeG
		summary.CarbIntakeG += record.CarbIntakeG
		summary.FatIntakeG += record.FatIntakeG
	}

	return summary, nil
}
