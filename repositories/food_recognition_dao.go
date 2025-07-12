package repositories

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
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
func (d *FoodRecognitionDAO) CreateRecognition(userID int64, sessionID, imageURL string, foods []models.RecognizedFoodItem, nutrition models.FoodRecognitionNutrition, aiResponse string) (*models.FoodRecognition, error) {
	// 将食物列表转为JSON字符串存储
	foodsJSON, err := json.Marshal(foods)
	if err != nil {
		return nil, err
	}

	recognition := models.FoodRecognition{
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
func (d *FoodRecognitionDAO) GetRecognitionByID(id int64) (*models.FoodRecognition, error) {
	var recognition models.FoodRecognition
	if err := d.db.First(&recognition, id).Error; err != nil {
		return nil, err
	}
	return &recognition, nil
}

// GetUserRecognitionsByDate 获取用户某天的所有食物识别记录
func (d *FoodRecognitionDAO) GetUserRecognitionsByDate(userID int64, date time.Time) ([]models.FoodRecognition, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nextDay := dateOnly.AddDate(0, 0, 1)

	var recognitions []models.FoodRecognition
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
func (d *FoodRecognitionDAO) GetUserTodayRecognitions(userID int64) ([]models.FoodRecognition, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return d.GetUserRecognitionsByDate(userID, today)
}

// GetUserRecentRecognitions 获取用户最近的食物识别记录
func (d *FoodRecognitionDAO) GetUserRecentRecognitions(userID int64, limit int) ([]models.FoodRecognition, error) {
	var recognitions []models.FoodRecognition
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&recognitions).Error
	if err != nil {
		return nil, err
	}
	return recognitions, nil
}

// UpdateAdoptionStatus 更新识别记录的采用状态
func (d *FoodRecognitionDAO) UpdateAdoptionStatus(id int64, isAdopted bool) error {
	return d.db.Model(&models.FoodRecognition{}).
		Where("id = ?", id).
		Update("is_adopted", isAdopted).Error
}

// ConvertToResult 将数据库记录转换为前端可用的结果对象
func (d *FoodRecognitionDAO) ConvertToResult(recognition *models.FoodRecognition) (*models.FoodRecognitionResult, error) {
	var foods []models.RecognizedFoodItem
	if err := json.Unmarshal([]byte(recognition.RecognizedFoods), &foods); err != nil {
		return nil, err
	}

	result := &models.FoodRecognitionResult{
		ID:              recognition.ID,
		ImageURL:        recognition.ImageURL,
		RecognizedFoods: foods,
		NutritionSummary: models.FoodRecognitionNutrition{
			CaloriesIntake: recognition.CaloriesIntake,
			ProteinIntakeG: recognition.ProteinIntakeG,
			CarbIntakeG:    recognition.CarbIntakeG,
			FatIntakeG:     recognition.FatIntakeG,
		},
		AIAnalysis: recognition.AIResponse,
		IsAdopted:  recognition.IsAdopted,
		RecordDate: recognition.RecordDate.Format("2006-01-02"),
	}

	return result, nil
}

// SummarizeTodayNutrition 汇总用户今日所有食物识别的营养摄入
func (d *FoodRecognitionDAO) SummarizeTodayNutrition(userID int64) (models.FoodRecognitionNutrition, error) {
	records, err := d.GetUserTodayRecognitions(userID)
	if err != nil {
		return models.FoodRecognitionNutrition{}, err
	}

	var summary models.FoodRecognitionNutrition
	for _, record := range records {
		summary.CaloriesIntake += record.CaloriesIntake
		summary.ProteinIntakeG += record.ProteinIntakeG
		summary.CarbIntakeG += record.CarbIntakeG
		summary.FatIntakeG += record.FatIntakeG
	}

	return summary, nil
}

// GetUserAdoptedRecognitions 获取用户已采用的食物识别记录
func (d *FoodRecognitionDAO) GetUserAdoptedRecognitions(userID int64, page, pageSize int, startDate, endDate time.Time) ([]models.FoodRecognition, int64, error) {
	var recognitions []models.FoodRecognition
	var total int64

	// 设置分页参数
	offset := (page - 1) * pageSize

	// 查询总数
	query := d.db.Model(&models.FoodRecognition{}).
		Where("user_id = ? AND is_adopted = ? AND record_date >= ? AND record_date <= ?",
			userID, true, startDate, endDate)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	err = query.Order("record_date DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&recognitions).Error

	if err != nil {
		return nil, 0, err
	}

	return recognitions, total, nil
}

// GetUserAdoptedRecognitionsByDateRange 获取用户已采用的食物识别记录并按日期分组
func (d *FoodRecognitionDAO) GetUserAdoptedRecognitionsByDateRange(userID int64, startDate, endDate time.Time) (map[string][]models.FoodRecognition, error) {
	var recognitions []models.FoodRecognition

	// 查询数据
	err := d.db.Where("user_id = ? AND is_adopted = ? AND record_date >= ? AND record_date <= ?",
		userID, true, startDate, endDate).
		Order("record_date DESC, created_at DESC").
		Find(&recognitions).Error

	if err != nil {
		return nil, err
	}

	// 按日期分组
	result := make(map[string][]models.FoodRecognition)
	for _, recognition := range recognitions {
		dateKey := recognition.RecordDate.Format("2006-01-02")
		result[dateKey] = append(result[dateKey], recognition)
	}

	return result, nil
}
