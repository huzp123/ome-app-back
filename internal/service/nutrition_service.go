package service

import (
	"time"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// NutritionService 处理用户营养服务
type NutritionService struct {
	nutritionDAO *dao.DailyNutritionDAO
}

// NewNutritionService 创建营养服务实例
func NewNutritionService(nutritionDAO *dao.DailyNutritionDAO) *NutritionService {
	return &NutritionService{
		nutritionDAO: nutritionDAO,
	}
}

// GetTodayNutrition 获取用户今日营养数据
func (s *NutritionService) GetTodayNutrition(userID int64) (*model.DailyNutrition, error) {
	return s.nutritionDAO.GetOrCreate(userID, time.Now())
}

// UpdateTodayNutrition 更新今日营养摄入数据
func (s *NutritionService) UpdateTodayNutrition(userID int64, caloriesIntake, proteinIntakeG, carbIntakeG, fatIntakeG float64) (*model.DailyNutrition, error) {
	// 获取或创建今日记录
	nutrition, err := s.nutritionDAO.GetOrCreate(userID, time.Now())
	if err != nil {
		return nil, err
	}

	// 更新数据
	nutrition.CaloriesIntake = caloriesIntake
	nutrition.ProteinIntakeG = proteinIntakeG
	nutrition.CarbIntakeG = carbIntakeG
	nutrition.FatIntakeG = fatIntakeG

	// 保存更新
	if err := s.nutritionDAO.Update(nutrition); err != nil {
		return nil, err
	}

	return nutrition, nil
}

// GetNutritionHistory 获取营养历史记录
func (s *NutritionService) GetNutritionHistory(userID int64, startDate, endDate time.Time) ([]model.DailyNutrition, error) {
	return s.nutritionDAO.GetHistory(userID, startDate, endDate)
}

// GetWeekSummary 获取一周营养摄入统计
func (s *NutritionService) GetWeekSummary(userID int64) (map[string]float64, error) {
	return s.nutritionDAO.GetWeekSummary(userID)
}
