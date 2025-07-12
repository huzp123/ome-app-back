package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
)

// MoodRecordDAO 处理情绪记录数据访问
type MoodRecordDAO struct {
	db *gorm.DB
}

// NewMoodRecordDAO 创建情绪记录DAO实例
func NewMoodRecordDAO(db *gorm.DB) *MoodRecordDAO {
	return &MoodRecordDAO{db: db}
}

// Create 创建情绪记录
func (d *MoodRecordDAO) Create(mood *models.MoodRecord) error {
	return d.db.Create(mood).Error
}

// GetByID 根据ID获取情绪记录
func (d *MoodRecordDAO) GetByID(userID, moodID int64) (*models.MoodRecord, error) {
	var mood models.MoodRecord
	err := d.db.Where("id = ? AND user_id = ?", moodID, userID).First(&mood).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("情绪记录不存在")
		}
		return nil, err
	}
	return &mood, nil
}

// GetHistory 获取用户情绪历史记录
func (d *MoodRecordDAO) GetHistory(userID int64, startDate, endDate time.Time, limit int) ([]models.MoodRecord, error) {
	var moods []models.MoodRecord
	query := d.db.Where("user_id = ? AND record_time BETWEEN ? AND ?", userID, startDate, endDate).
		Order("record_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&moods).Error
	if err != nil {
		return nil, err
	}
	return moods, nil
}

// GetTodayMoods 获取今日情绪记录
func (d *MoodRecordDAO) GetTodayMoods(userID int64) ([]models.MoodRecord, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var moods []models.MoodRecord
	err := d.db.Where("user_id = ? AND record_time BETWEEN ? AND ?", userID, startOfDay, endOfDay).
		Order("record_time DESC").
		Find(&moods).Error
	if err != nil {
		return nil, err
	}
	return moods, nil
}

// Delete 删除情绪记录
func (d *MoodRecordDAO) Delete(userID, moodID int64) error {
	result := d.db.Where("id = ? AND user_id = ?", moodID, userID).Delete(&models.MoodRecord{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("记录不存在或无权限删除")
	}
	return nil
}

// GetMoodStatistics 获取情绪统计数据
func (d *MoodRecordDAO) GetMoodStatistics(userID int64, startDate, endDate time.Time) (map[string]interface{}, error) {
	var result struct {
		TotalRecords int64   `json:"total_records"`
		AvgMoodLevel float64 `json:"avg_mood_level"`
	}

	// 获取基本统计
	err := d.db.Model(&models.MoodRecord{}).
		Select(`
			COUNT(*) as total_records,
			COALESCE(AVG(mood_level), 0) as avg_mood_level
		`).
		Where("user_id = ? AND record_time BETWEEN ? AND ?", userID, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 注意：由于mood_tags和influences现在是JSON数组，
	// 复杂的统计查询需要特殊的SQL处理，这里先返回基本统计
	return map[string]interface{}{
		"total_records":  result.TotalRecords,
		"avg_mood_level": result.AvgMoodLevel,
	}, nil
}
