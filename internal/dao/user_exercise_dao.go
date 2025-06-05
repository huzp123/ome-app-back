package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// UserExerciseDAO 处理用户运动数据访问
type UserExerciseDAO struct {
	db *gorm.DB
}

// NewUserExerciseDAO 创建用户运动DAO实例
func NewUserExerciseDAO(db *gorm.DB) *UserExerciseDAO {
	return &UserExerciseDAO{db: db}
}

// Create 创建运动记录
func (d *UserExerciseDAO) Create(exercise *model.UserExercise) error {
	return d.db.Create(exercise).Error
}

// GetByID 根据ID获取运动记录
func (d *UserExerciseDAO) GetByID(userID, exerciseID int64) (*model.UserExercise, error) {
	var exercise model.UserExercise
	err := d.db.Where("id = ? AND user_id = ?", exerciseID, userID).First(&exercise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("运动记录不存在")
		}
		return nil, err
	}
	return &exercise, nil
}

// GetHistory 获取用户运动历史记录
func (d *UserExerciseDAO) GetHistory(userID int64, startDate, endDate time.Time, limit int) ([]model.UserExercise, error) {
	var exercises []model.UserExercise
	query := d.db.Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, startDate, endDate).
		Order("start_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&exercises).Error
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

// GetTodayExercises 获取今日运动记录
func (d *UserExerciseDAO) GetTodayExercises(userID int64) ([]model.UserExercise, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var exercises []model.UserExercise
	err := d.db.Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, startOfDay, endOfDay).
		Order("start_time DESC").
		Find(&exercises).Error
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

// Update 更新运动记录
func (d *UserExerciseDAO) Update(exercise *model.UserExercise) error {
	return d.db.Save(exercise).Error
}

// Delete 删除运动记录
func (d *UserExerciseDAO) Delete(userID, exerciseID int64) error {
	result := d.db.Where("id = ? AND user_id = ?", exerciseID, userID).Delete(&model.UserExercise{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("记录不存在或无权限删除")
	}
	return nil
}

// GetStatistics 获取运动统计数据
func (d *UserExerciseDAO) GetStatistics(userID int64, startDate, endDate time.Time) (map[string]interface{}, error) {
	var result struct {
		TotalExercises int64   `json:"total_exercises"`
		TotalDuration  float64 `json:"total_duration"`
		TotalCalories  float64 `json:"total_calories"`
		TotalDistance  float64 `json:"total_distance"`
		AvgDuration    float64 `json:"avg_duration"`
		AvgCalories    float64 `json:"avg_calories"`
	}

	err := d.db.Model(&model.UserExercise{}).
		Select(`
			COUNT(*) as total_exercises,
			COALESCE(SUM(duration_min), 0) as total_duration,
			COALESCE(SUM(calories_burned), 0) as total_calories,
			COALESCE(SUM(distance_km), 0) as total_distance,
			COALESCE(AVG(duration_min), 0) as avg_duration,
			COALESCE(AVG(calories_burned), 0) as avg_calories
		`).
		Where("user_id = ? AND start_time BETWEEN ? AND ?", userID, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_exercises": result.TotalExercises,
		"total_duration":  result.TotalDuration,
		"total_calories":  result.TotalCalories,
		"total_distance":  result.TotalDistance,
		"avg_duration":    result.AvgDuration,
		"avg_calories":    result.AvgCalories,
	}, nil
}
