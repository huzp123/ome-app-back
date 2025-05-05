package dao

import (
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/model"
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
	weight := model.UserWeight{
		UserID:     userID,
		WeightKG:   weightKG,
		RecordDate: time.Now(),
	}
	return d.db.Create(&weight).Error
}

// GetLatest 获取用户最新体重记录
func (d *UserWeightDAO) GetLatest(userID int64) (*model.UserWeight, error) {
	var weight model.UserWeight
	err := d.db.Where("user_id = ?", userID).
		Order("record_date DESC").
		First(&weight).Error
	if err != nil {
		return nil, err
	}
	return &weight, nil
}

// GetHistory 获取用户体重历史记录
func (d *UserWeightDAO) GetHistory(userID int64, startDate, endDate time.Time) ([]model.UserWeight, error) {
	var weights []model.UserWeight
	err := d.db.Where("user_id = ? AND record_date BETWEEN ? AND ?",
		userID, startDate, endDate).
		Order("record_date ASC").
		Find(&weights).Error
	if err != nil {
		return nil, err
	}
	return weights, nil
}
