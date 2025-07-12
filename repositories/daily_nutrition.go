package repositories

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
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
func (d *DailyNutritionDAO) GetOrCreate(userID int64, date time.Time, targetParams *CreateNutritionParams) (*models.DailyNutrition, error) {
	// 将时间设置为当天的0点，使用UTC时区避免时区问题
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	log.Printf("[营养DAO] GetOrCreate 用户(ID:%d), 日期:%s", userID, dateOnly.Format("2006-01-02"))

	if targetParams == nil {
		log.Printf("[营养DAO] 用户(ID:%d)创建记录失败：缺少目标参数", userID)
		return nil, errors.New("创建营养记录需要提供目标值参数")
	}

	var nutritionRecord models.DailyNutrition

	err := d.db.
		Where("user_id = ? AND date = ?", userID, dateOnly).
		Attrs(models.DailyNutrition{
			UserID:         userID,
			Date:           dateOnly,
			TargetCalories: targetParams.TargetCalories,
			TargetProteinG: targetParams.TargetProteinG,
			TargetCarbG:    targetParams.TargetCarbG,
			TargetFatG:     targetParams.TargetFatG,
		}).
		FirstOrCreate(&nutritionRecord).Error

	if err != nil {
		log.Printf("[营养DAO] GetOrCreate 操作失败 用户(ID:%d), 错误: %v", userID, err)
		return nil, err
	}

	log.Printf("[营养DAO] GetOrCreate 操作成功 用户(ID:%d), 记录ID: %d", userID, nutritionRecord.ID)
	return &nutritionRecord, nil
}

// GetByDate 通过日期获取特定用户的营养记录
func (d *DailyNutritionDAO) GetByDate(userID int64, date time.Time) (*models.DailyNutrition, error) {
	dateStr := date.Format("2006-01-02")
	log.Printf("[营养DAO] GetByDate查询 用户(ID:%d), 日期:%s", userID, dateStr)

	var nutrition models.DailyNutrition
	// 使用原生SQL查询，避免时区问题
	err := d.db.Raw("SELECT * FROM daily_nutrition WHERE user_id = ? AND DATE(date) = ?", userID, dateStr).Scan(&nutrition).Error
	if err != nil || nutrition.ID == 0 { // Scan不会在未找到时返回error，需要检查ID
		if err == nil {
			err = gorm.ErrRecordNotFound
		}
		log.Printf("[营养DAO] GetByDate查询失败 用户(ID:%d), 错误:%v", userID, err)
		return nil, err
	}

	log.Printf("[营养DAO] GetByDate查询成功 用户(ID:%d), 记录ID:%d, 记录日期:%s", userID, nutrition.ID, nutrition.Date.Format("2006-01-02 15:04:05"))
	return &nutrition, nil
}

// Update 更新当日营养摄入数据
func (d *DailyNutritionDAO) Update(nutrition *models.DailyNutrition) error {
	// 计算热量完成率
	if nutrition.TargetCalories > 0 {
		nutrition.CaloriesCompletionRate = (nutrition.CaloriesIntake / nutrition.TargetCalories) * 100
	}

	return d.db.Save(nutrition).Error
}

// GetHistory 获取用户营养历史记录
func (d *DailyNutritionDAO) GetHistory(userID int64, startDate, endDate time.Time) ([]models.DailyNutrition, error) {
	var records []models.DailyNutrition
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

	var records []models.DailyNutrition
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
