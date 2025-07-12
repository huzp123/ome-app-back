package services

import (
	"errors"
	"time"

	"ome-app-back/models"
	"ome-app-back/repositories"
)

// HeightService 处理用户身高相关业务逻辑
type HeightService struct {
	heightDAO *repositories.UserHeightDAO
}

// NewHeightService 创建身高服务实例
func NewHeightService(heightDAO *repositories.UserHeightDAO) *HeightService {
	return &HeightService{
		heightDAO: heightDAO,
	}
}

// CreateHeightRequest 创建身高记录请求
type CreateHeightRequest struct {
	HeightCM float64 `json:"height_cm" binding:"required,min=50,max=300"`
}

// CreateHeight 创建身高记录
func (s *HeightService) CreateHeight(userID int64, req CreateHeightRequest) error {
	// 验证身高范围
	if req.HeightCM < 50 || req.HeightCM > 300 {
		return errors.New("身高必须在50-300厘米之间")
	}

	// 检查今天是否已有身高记录
	today := time.Now()
	existingHeight, err := s.heightDAO.GetHeightByDate(userID, today)
	if err != nil {
		return err
	}

	if existingHeight != nil {
		// 更新今天的记录
		existingHeight.HeightCM = req.HeightCM
		return s.heightDAO.Update(existingHeight)
	}

	// 创建新记录
	height := &models.UserHeight{
		UserID:     userID,
		HeightCM:   req.HeightCM,
		RecordDate: today,
	}

	return s.heightDAO.Create(height)
}

// GetHeightHistoryRequest 获取身高历史请求
type GetHeightHistoryRequest struct {
	Limit int `form:"limit" binding:"min=1,max=365"`
}

// GetHeightHistoryResponse 获取身高历史响应
type GetHeightHistoryResponse struct {
	ID         int64     `json:"id"`
	HeightCM   float64   `json:"height_cm"`
	RecordDate time.Time `json:"record_date"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetHeightHistory 获取身高历史记录
func (s *HeightService) GetHeightHistory(userID int64, req GetHeightHistoryRequest) ([]GetHeightHistoryResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 30 // 默认30条
	}

	heights, err := s.heightDAO.GetHeightHistory(userID, req.Limit)
	if err != nil {
		return nil, err
	}

	var response []GetHeightHistoryResponse
	for _, height := range heights {
		response = append(response, GetHeightHistoryResponse{
			ID:         height.ID,
			HeightCM:   height.HeightCM,
			RecordDate: height.RecordDate,
			CreatedAt:  height.CreatedAt,
		})
	}

	return response, nil
}

// GetCurrentHeightResponse 获取当前身高响应
type GetCurrentHeightResponse struct {
	HeightCM   float64   `json:"height_cm"`
	RecordDate time.Time `json:"record_date"`
	DaysAgo    int       `json:"days_ago"`
}

// GetCurrentHeight 获取当前身高
func (s *HeightService) GetCurrentHeight(userID int64) (*GetCurrentHeightResponse, error) {
	height, err := s.heightDAO.GetCurrentHeight(userID)
	if err != nil {
		return nil, err
	}

	if height == nil {
		return nil, errors.New("没有身高记录")
	}

	daysAgo := int(time.Since(height.RecordDate).Hours() / 24)

	return &GetCurrentHeightResponse{
		HeightCM:   height.HeightCM,
		RecordDate: height.RecordDate,
		DaysAgo:    daysAgo,
	}, nil
}

// DeleteHeight 删除身高记录
func (s *HeightService) DeleteHeight(userID int64, heightID int64) error {
	// 验证记录是否属于当前用户
	height, err := s.heightDAO.GetByID(heightID)
	if err != nil {
		return err
	}

	if height.UserID != userID {
		return errors.New("只能删除自己的身高记录")
	}

	return s.heightDAO.Delete(heightID)
}

// GetHeightStatisticsRequest 获取身高统计请求
type GetHeightStatisticsRequest struct {
	Days int `form:"days" binding:"min=1,max=365"`
}

// GetHeightStatisticsResponse 获取身高统计响应
type GetHeightStatisticsResponse struct {
	CurrentHeight float64                  `json:"current_height"`
	MinHeight     float64                  `json:"min_height"`
	MaxHeight     float64                  `json:"max_height"`
	AvgHeight     float64                  `json:"avg_height"`
	HeightChange  float64                  `json:"height_change"`
	TrendData     []map[string]interface{} `json:"trend_data"`
}

// GetHeightStatistics 获取身高统计数据
func (s *HeightService) GetHeightStatistics(userID int64, req GetHeightStatisticsRequest) (*GetHeightStatisticsResponse, error) {
	if req.Days <= 0 {
		req.Days = 30 // 默认30天
	}

	stats, err := s.heightDAO.GetHeightStatistics(userID, req.Days)
	if err != nil {
		return nil, err
	}

	return &GetHeightStatisticsResponse{
		CurrentHeight: stats["current_height"].(float64),
		MinHeight:     stats["min_height"].(float64),
		MaxHeight:     stats["max_height"].(float64),
		AvgHeight:     stats["avg_height"].(float64),
		HeightChange:  stats["height_change"].(float64),
		TrendData:     stats["trend_data"].([]map[string]interface{}),
	}, nil
}
