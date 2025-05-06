package dao

import (
	"errors"
	"strings"
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

// CreateNutritionParams 创建营养记录的参数
type CreateNutritionParams struct {
	UserID         int64
	Date           time.Time
	TargetCalories float64
	TargetProteinG float64
	TargetCarbG    float64
	TargetFatG     float64
}

// GetOrCreate 获取或创建指定日期的营养记录
func (d *DailyNutritionDAO) GetOrCreate(userID int64, date time.Time, targetParams *CreateNutritionParams) (*model.DailyNutrition, error) {
	// 将时间设置为当天的0点
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var nutrition model.DailyNutrition
	result := d.db.Where("user_id = ? AND date = ?", userID, dateOnly).First(&nutrition)

	if result.Error == nil {
		return &nutrition, nil
	}

	// 如果不是记录未找到的错误，直接返回错误
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 如果没有提供目标参数，返回错误
	if targetParams == nil {
		return nil, errors.New("创建营养记录需要提供目标值参数")
	}

	// 创建新记录
	nutrition = model.DailyNutrition{
		UserID:         userID,
		Date:           dateOnly,
		TargetCalories: targetParams.TargetCalories,
		TargetProteinG: targetParams.TargetProteinG,
		TargetCarbG:    targetParams.TargetCarbG,
		TargetFatG:     targetParams.TargetFatG,
	}

	err := d.db.Create(&nutrition).Error
	if err != nil {
		// 处理唯一键冲突的情况
		// MySQL错误包含 "Duplicate entry" 和 "idx_user_date"
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "idx_user_date") {
			// 如果是唯一键冲突，说明在我们检查和插入之间，有其他请求已经创建了记录
			// 再次尝试获取记录
			retryResult := d.db.Where("user_id = ? AND date = ?", userID, dateOnly).First(&nutrition)
			if retryResult.Error == nil {
				return &nutrition, nil
			}
			return nil, retryResult.Error
		}
		return nil, err
	}

	return &nutrition, nil
}

// GetByDate 通过日期获取特定用户的营养记录
func (d *DailyNutritionDAO) GetByDate(userID int64, date time.Time) (*model.DailyNutrition, error) {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var nutrition model.DailyNutrition
	err := d.db.Where("user_id = ? AND date = ?", userID, dateOnly).First(&nutrition).Error
	if err != nil {
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
