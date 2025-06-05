package service

import (
	"errors"
	"time"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// WeightService 体重服务
type WeightService struct {
	userWeightDAO *dao.UserWeightDAO
}

// NewWeightService 创建体重服务实例
func NewWeightService(userWeightDAO *dao.UserWeightDAO) *WeightService {
	return &WeightService{
		userWeightDAO: userWeightDAO,
	}
}

// CreateWeightRequest 创建体重记录请求
type CreateWeightRequest struct {
	WeightKG float64 `json:"weight_kg" binding:"required,gt=0,lt=500"`
}

// WeightHistoryRequest 体重历史记录请求
type WeightHistoryRequest struct {
	Limit int `form:"limit"`
}

// WeightStatisticsRequest 体重统计请求
type WeightStatisticsRequest struct {
	Days int `form:"days"`
}

// WeightResponse 体重记录响应
type WeightResponse struct {
	ID         int64     `json:"id"`
	WeightKG   float64   `json:"weight_kg"`
	RecordDate time.Time `json:"record_date"`
	CreatedAt  time.Time `json:"created_at"`
}

// CurrentWeightResponse 当前体重响应
type CurrentWeightResponse struct {
	WeightKG   float64   `json:"weight_kg"`
	RecordDate time.Time `json:"record_date"`
	DaysAgo    int       `json:"days_ago"`
}

// WeightStatisticsResponse 体重统计响应
type WeightStatisticsResponse struct {
	CurrentWeight float64            `json:"current_weight"`
	MinWeight     float64            `json:"min_weight"`
	MaxWeight     float64            `json:"max_weight"`
	AvgWeight     float64            `json:"avg_weight"`
	WeightChange  float64            `json:"weight_change"`
	TrendData     []WeightTrendPoint `json:"trend_data"`
}

// WeightTrendPoint 体重趋势数据点
type WeightTrendPoint struct {
	Date     time.Time `json:"date"`
	WeightKG float64   `json:"weight_kg"`
}

// CreateWeight 创建体重记录
func (s *WeightService) CreateWeight(userID int64, req *CreateWeightRequest) error {
	return s.userWeightDAO.Create(userID, req.WeightKG)
}

// GetWeightHistory 获取体重历史记录
func (s *WeightService) GetWeightHistory(userID int64, req *WeightHistoryRequest) ([]WeightResponse, error) {
	// 设置默认限制
	limit := 30
	if req.Limit > 0 && req.Limit <= 365 {
		limit = req.Limit
	}

	// 计算时间范围（最多一年）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -limit)

	weights, err := s.userWeightDAO.GetHistory(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式，按时间倒序排列
	result := make([]WeightResponse, 0, len(weights))
	for i := len(weights) - 1; i >= 0; i-- {
		weight := weights[i]
		result = append(result, WeightResponse{
			ID:         weight.ID,
			WeightKG:   weight.WeightKG,
			RecordDate: weight.RecordDate,
			CreatedAt:  weight.CreatedAt,
		})
	}

	return result, nil
}

// GetCurrentWeight 获取当前体重信息
func (s *WeightService) GetCurrentWeight(userID int64) (*CurrentWeightResponse, error) {
	weight, err := s.userWeightDAO.GetLatest(userID)
	if err != nil {
		return nil, err
	}

	// 计算距离现在多少天
	daysAgo := int(time.Since(weight.RecordDate).Hours() / 24)

	return &CurrentWeightResponse{
		WeightKG:   weight.WeightKG,
		RecordDate: weight.RecordDate,
		DaysAgo:    daysAgo,
	}, nil
}

// DeleteWeight 删除体重记录
func (s *WeightService) DeleteWeight(userID, recordID int64) error {
	return s.userWeightDAO.Delete(userID, recordID)
}

// GetWeightStatistics 获取体重统计分析
func (s *WeightService) GetWeightStatistics(userID int64, req *WeightStatisticsRequest) (*WeightStatisticsResponse, error) {
	// 设置默认天数
	days := 30
	if req.Days > 0 && req.Days <= 365 {
		days = req.Days
	}

	// 计算时间范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	allWeights, err := s.userWeightDAO.GetHistory(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(allWeights) == 0 {
		return nil, errors.New("没有体重记录数据")
	}

	// 过滤每天的记录，只保留每天最后一次记录
	weights := filterDailyLastRecords(allWeights)

	if len(weights) == 0 {
		return nil, errors.New("没有有效的体重记录数据")
	}

	// 计算统计数据
	var totalWeight float64
	minWeight := weights[0].WeightKG
	maxWeight := weights[0].WeightKG

	for _, weight := range weights {
		totalWeight += weight.WeightKG
		if weight.WeightKG < minWeight {
			minWeight = weight.WeightKG
		}
		if weight.WeightKG > maxWeight {
			maxWeight = weight.WeightKG
		}
	}

	avgWeight := totalWeight / float64(len(weights))

	// 计算体重变化（最新 - 最早）
	weightChange := weights[len(weights)-1].WeightKG - weights[0].WeightKG

	// 生成趋势数据（简化版，选择关键数据点）
	trendData := make([]WeightTrendPoint, 0)

	// 如果数据点较少，全部返回
	if len(weights) <= 10 {
		for _, weight := range weights {
			trendData = append(trendData, WeightTrendPoint{
				Date:     weight.RecordDate,
				WeightKG: weight.WeightKG,
			})
		}
	} else {
		// 如果数据点较多，选择均匀分布的点
		step := len(weights) / 10
		for i := 0; i < len(weights); i += step {
			weight := weights[i]
			trendData = append(trendData, WeightTrendPoint{
				Date:     weight.RecordDate,
				WeightKG: weight.WeightKG,
			})
		}
		// 确保包含最后一个数据点
		if len(trendData) == 0 || trendData[len(trendData)-1].Date != weights[len(weights)-1].RecordDate {
			lastWeight := weights[len(weights)-1]
			trendData = append(trendData, WeightTrendPoint{
				Date:     lastWeight.RecordDate,
				WeightKG: lastWeight.WeightKG,
			})
		}
	}

	return &WeightStatisticsResponse{
		CurrentWeight: weights[len(weights)-1].WeightKG,
		MinWeight:     minWeight,
		MaxWeight:     maxWeight,
		AvgWeight:     avgWeight,
		WeightChange:  weightChange,
		TrendData:     trendData,
	}, nil
}

// filterDailyLastRecords 过滤每天的记录，只保留每天最后一次记录
func filterDailyLastRecords(weights []model.UserWeight) []model.UserWeight {
	if len(weights) == 0 {
		return weights
	}

	// 使用map来存储每天的最后一条记录
	dailyLastRecord := make(map[string]model.UserWeight)

	for _, weight := range weights {
		// 获取日期字符串 (YYYY-MM-DD)
		dateKey := weight.RecordDate.Format("2006-01-02")

		// 如果该日期还没有记录，或者当前记录的时间更晚，则更新该日期的记录
		if existing, exists := dailyLastRecord[dateKey]; !exists || weight.RecordDate.After(existing.RecordDate) {
			dailyLastRecord[dateKey] = weight
		}
	}

	// 将map中的记录转换为切片，并按时间排序
	result := make([]model.UserWeight, 0, len(dailyLastRecord))
	for _, weight := range dailyLastRecord {
		result = append(result, weight)
	}

	// 按记录日期排序（从早到晚）
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].RecordDate.After(result[j].RecordDate) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}
