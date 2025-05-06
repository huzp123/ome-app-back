package dao

import (
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// DailyNutritionDAO 处理用户每日营养数据访问
type DailyNutritionDAO struct {
	db *gorm.DB
}

// NewDailyNutritionDAO 创建每日营养DAO实例
func NewDailyNutritionDAO(db *gorm.DB) *DailyNutritionDAO {
	return &DailyNutritionDAO{db: db}
}

// GetOrCreate 获取或创建指定日期的营养记录
func (d *DailyNutritionDAO) GetOrCreate(userID int64, date time.Time) (*model.DailyNutrition, error) {
	// 将时间设置为当天的0点
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var nutrition model.DailyNutrition
	result := d.db.Where("user_id = ? AND date = ?", userID, dateOnly).First(&nutrition)

	if result.Error == nil {
		return &nutrition, nil
	}

	// 不存在则创建新记录，同时获取最新的健康分析设置目标值
	var analysis model.HealthAnalysis
	d.db.Where("user_id = ?", userID).Order("created_at DESC").First(&analysis)

	// 创建新记录
	nutrition = model.DailyNutrition{
		UserID:         userID,
		Date:           dateOnly,
		TargetCalories: analysis.RecommendedCalories,
		TargetProteinG: analysis.ProteinNeedG,
		TargetCarbG:    analysis.CarbNeedG,
		TargetFatG:     analysis.FatNeedG,
	}

	if err := d.db.Create(&nutrition).Error; err != nil {
		return nil, err
	}

	return &nutrition, nil
}

// Update 更新当日营养摄入数据
func (d *DailyNutritionDAO) Update(nutrition *model.DailyNutrition) error {
	// 计算热量完成率
	if nutrition.TargetCalories > 0 {
		nutrition.CaloriesCompletionRate = (nutrition.CaloriesIntake / nutrition.TargetCalories) * 100
	}

	return d.db.Save(nutrition).Error
}

// GetHistory 获取用户营养历史记录
func (d *DailyNutritionDAO) GetHistory(userID int64, startDate, endDate time.Time) ([]model.DailyNutrition, error) {
	var records []model.DailyNutrition
	err := d.db.Where("user_id = ? AND date BETWEEN ? AND ?",
		userID, startDate, endDate).
		Order("date DESC").
		Find(&records).Error
	return records, err
}

// GetWeekSummary 获取用户一周营养摄入统计
func (d *DailyNutritionDAO) GetWeekSummary(userID int64) (map[string]float64, error) {
	// 获取最近7天数据
	now := time.Now()
	startDate := now.AddDate(0, 0, -6) // 6天前

	var records []model.DailyNutrition
	err := d.db.Where("user_id = ? AND date BETWEEN ? AND ?",
		userID, startDate, now).Find(&records).Error

	if err != nil {
		return nil, err
	}

	// 计算平均值
	summary := map[string]float64{
		"avg_calories":        0,
		"avg_protein":         0,
		"avg_carb":            0,
		"avg_fat":             0,
		"avg_completion_rate": 0,
	}

	if len(records) > 0 {
		for _, r := range records {
			summary["avg_calories"] += r.CaloriesIntake
			summary["avg_protein"] += r.ProteinIntakeG
			summary["avg_carb"] += r.CarbIntakeG
			summary["avg_fat"] += r.FatIntakeG
			summary["avg_completion_rate"] += r.CaloriesCompletionRate
		}

		days := float64(len(records))
		summary["avg_calories"] /= days
		summary["avg_protein"] /= days
		summary["avg_carb"] /= days
		summary["avg_fat"] /= days
		summary["avg_completion_rate"] /= days
	}

	return summary, nil
}
