package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
)

// UserWeightDAO 处理用户体重数据访问
type UserWeightDAO struct {
	db *gorm.DB
}

// NewUserWeightDAO 创建用户体重DAO实例
func NewUserWeightDAO(db *gorm.DB) *UserWeightDAO {
	return &UserWeightDAO{db: db}
}

// Create 记录用户体重
func (d *UserWeightDAO) Create(userID int64, weightKG float64) error {
	weight := models.UserWeight{
		UserID:     userID,
		WeightKG:   weightKG,
		RecordDate: time.Now(),
	}
	return d.db.Create(&weight).Error
}

// GetLatest 获取用户最新体重记录
func (d *UserWeightDAO) GetLatest(userID int64) (*models.UserWeight, error) {
	var weight models.UserWeight
	err := d.db.Where("user_id = ?", userID).
		Order("id DESC").
		First(&weight).Error
	if err != nil {
		return nil, err
	}
	return &weight, nil
}

// GetHistory 获取用户体重历史记录
func (d *UserWeightDAO) GetHistory(userID int64, startDate, endDate time.Time) ([]models.UserWeight, error) {
	var weights []models.UserWeight
	err := d.db.Where("user_id = ? AND record_date BETWEEN ? AND ?",
		userID, startDate, endDate).
		Order("record_date ASC").
		Find(&weights).Error
	if err != nil {
		return nil, err
	}
	return weights, nil
}

// Delete 删除用户体重记录
func (d *UserWeightDAO) Delete(userID, recordID int64) error {
	result := d.db.Where("id = ? AND user_id = ?", recordID, userID).Delete(&models.UserWeight{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("记录不存在或无权限删除")
	}
	return nil
}
