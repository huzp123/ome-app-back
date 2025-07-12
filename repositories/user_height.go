package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
)

// UserHeightDAO 处理用户身高数据访问
type UserHeightDAO struct {
	db *gorm.DB
}

// NewUserHeightDAO 创建用户身高DAO实例
func NewUserHeightDAO(db *gorm.DB) *UserHeightDAO {
	return &UserHeightDAO{db: db}
}

// Create 创建身高记录
func (d *UserHeightDAO) Create(height *models.UserHeight) error {
	return d.db.Create(height).Error
}

// GetByID 根据ID获取身高记录
func (d *UserHeightDAO) GetByID(id int64) (*models.UserHeight, error) {
	var height models.UserHeight
	if err := d.db.First(&height, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("身高记录不存在")
		}
		return nil, err
	}
	return &height, nil
}

// GetCurrentHeight 获取用户当前身高
func (d *UserHeightDAO) GetCurrentHeight(userID int64) (*models.UserHeight, error) {
	var height models.UserHeight
	if err := d.db.Where("user_id = ?", userID).
		Order("record_date DESC, created_at DESC").
		First(&height).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回nil表示没有身高记录
		}
		return nil, err
	}
	return &height, nil
}

// GetHeightHistory 获取用户身高历史记录
func (d *UserHeightDAO) GetHeightHistory(userID int64, limit int) ([]models.UserHeight, error) {
	var heights []models.UserHeight
	err := d.db.Where("user_id = ?", userID).
		Order("record_date DESC, created_at DESC").
		Limit(limit).
		Find(&heights).Error
	if err != nil {
		return nil, err
	}
	return heights, nil
}

// GetHeightByDate 根据日期获取身高记录
func (d *UserHeightDAO) GetHeightByDate(userID int64, date time.Time) (*models.UserHeight, error) {
	var height models.UserHeight
	dateStr := date.Format("2006-01-02")
	if err := d.db.Where("user_id = ? AND DATE(record_date) = ?", userID, dateStr).
		First(&height).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &height, nil
}

// Update 更新身高记录
func (d *UserHeightDAO) Update(height *models.UserHeight) error {
	return d.db.Save(height).Error
}

// Delete 删除身高记录
func (d *UserHeightDAO) Delete(id int64) error {
	return d.db.Delete(&models.UserHeight{}, id).Error
}

// GetHeightStatistics 获取身高统计数据
func (d *UserHeightDAO) GetHeightStatistics(userID int64, days int) (map[string]interface{}, error) {
	var result struct {
		CurrentHeight float64 `json:"current_height"`
		MinHeight     float64 `json:"min_height"`
		MaxHeight     float64 `json:"max_height"`
		AvgHeight     float64 `json:"avg_height"`
		HeightChange  float64 `json:"height_change"`
	}

	// 获取当前身高
	currentHeight, err := d.GetCurrentHeight(userID)
	if err != nil {
		return nil, err
	}
	if currentHeight == nil {
		return nil, errors.New("没有身高记录")
	}
	result.CurrentHeight = currentHeight.HeightCM

	// 获取指定天数内的统计数据
	startDate := time.Now().AddDate(0, 0, -days)

	// 获取最早的身高记录
	var firstHeight models.UserHeight
	if err := d.db.Where("user_id = ? AND record_date >= ?", userID, startDate).
		Order("record_date ASC, created_at ASC").
		First(&firstHeight).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 计算身高变化
	if firstHeight.ID != 0 {
		result.HeightChange = currentHeight.HeightCM - firstHeight.HeightCM
	}

	// 获取统计信息
	var stats struct {
		MinHeight float64 `json:"min_height"`
		MaxHeight float64 `json:"max_height"`
		AvgHeight float64 `json:"avg_height"`
	}

	if err := d.db.Model(&models.UserHeight{}).
		Where("user_id = ? AND record_date >= ?", userID, startDate).
		Select("MIN(height_cm) as min_height, MAX(height_cm) as max_height, AVG(height_cm) as avg_height").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	result.MinHeight = stats.MinHeight
	result.MaxHeight = stats.MaxHeight
	result.AvgHeight = stats.AvgHeight

	// 获取趋势数据
	var trendData []map[string]interface{}
	if err := d.db.Model(&models.UserHeight{}).
		Where("user_id = ? AND record_date >= ?", userID, startDate).
		Select("record_date as date, height_cm").
		Order("record_date ASC").
		Scan(&trendData).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"current_height": result.CurrentHeight,
		"min_height":     result.MinHeight,
		"max_height":     result.MaxHeight,
		"avg_height":     result.AvgHeight,
		"height_change":  result.HeightChange,
		"trend_data":     trendData,
	}, nil
}
