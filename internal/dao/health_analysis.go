package dao

import (
	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// HealthAnalysisDAO 处理健康分析数据访问
type HealthAnalysisDAO struct {
	db *gorm.DB
}

// NewHealthAnalysisDAO 创建健康分析DAO实例
func NewHealthAnalysisDAO(db *gorm.DB) *HealthAnalysisDAO {
	return &HealthAnalysisDAO{db: db}
}

// Create 创建健康分析记录
func (d *HealthAnalysisDAO) Create(analysis *model.HealthAnalysis) error {
	return d.db.Create(analysis).Error
}

// GetLatestByUserID 获取用户最新健康分析
func (d *HealthAnalysisDAO) GetLatestByUserID(userID int64) (*model.HealthAnalysis, error) {
	var analysis model.HealthAnalysis
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// GetHistory 获取用户健康分析历史记录
func (d *HealthAnalysisDAO) GetHistory(userID int64, limit int) ([]model.HealthAnalysis, error) {
	var analyses []model.HealthAnalysis
	err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&analyses).Error
	if err != nil {
		return nil, err
	}
	return analyses, nil
}
