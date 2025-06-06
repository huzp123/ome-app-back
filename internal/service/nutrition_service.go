package service

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// NutritionService 处理用户营养服务
type NutritionService struct {
	nutritionDAO      *dao.DailyNutritionDAO
	healthAnalysisDAO *dao.HealthAnalysisDAO // 添加健康分析DAO依赖
}

// NewNutritionService 创建营养服务实例
func NewNutritionService(nutritionDAO *dao.DailyNutritionDAO, healthAnalysisDAO *dao.HealthAnalysisDAO) *NutritionService {
	return &NutritionService{
		nutritionDAO:      nutritionDAO,
		healthAnalysisDAO: healthAnalysisDAO,
	}
}

// 自定义错误，表示用户尚未生成健康分析报告
var ErrNoHealthAnalysis = errors.New("用户尚未生成健康分析报告，请先生成健康分析")

// GetTodayNutrition 获取用户今日营养数据
func (s *NutritionService) GetTodayNutrition(userID int64) (*model.DailyNutrition, error) {
	log.Printf("[营养服务] 开始获取用户(ID:%d)今日营养数据", userID)

	// 先尝试直接获取今天的记录
	today := time.Now()
	nutrition, err := s.nutritionDAO.GetByDate(userID, today)

	// 如果记录存在，直接返回
	if err == nil {
		log.Printf("[营养服务] 用户(ID:%d)今日营养记录已存在(ID:%d)", userID, nutrition.ID)
		return nutrition, nil
	}

	// 如果错误不是"记录未找到"，直接返回错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[营养服务] 用户(ID:%d)获取营养记录时出现非预期错误: %v", userID, err)
		// 不要使用fmt.Errorf包装错误，直接返回原错误以保持错误类型
		return nil, err
	}

	log.Printf("[营养服务] 用户(ID:%d)今日营养记录不存在，准备创建新记录", userID)

	// 记录不存在，需要创建新记录
	// 先获取用户最新的健康分析数据
	analysis, err := s.healthAnalysisDAO.GetLatestByUserID(userID)
	if err != nil {
		// 如果用户没有健康分析数据，返回特定错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[营养服务] 用户(ID:%d)没有健康分析数据", userID)
			return nil, ErrNoHealthAnalysis
		}
		log.Printf("[营养服务] 用户(ID:%d)获取健康分析数据失败: %v", userID, err)
		// 不要使用fmt.Errorf包装错误，直接返回原错误以保持错误类型
		return nil, err
	}

	log.Printf("[营养服务] 用户(ID:%d)健康分析数据获取成功(ID:%d)，目标热量:%.2f", userID, analysis.ID, analysis.RecommendedCalories)

	// 使用健康分析数据中的目标值创建今日营养记录
	createParams := &dao.CreateNutritionParams{
		UserID:         userID,
		Date:           today,
		TargetCalories: analysis.RecommendedCalories,
		TargetProteinG: analysis.ProteinNeedG,
		TargetCarbG:    analysis.CarbNeedG,
		TargetFatG:     analysis.FatNeedG,
	}

	log.Printf("[营养服务] 用户(ID:%d)开始调用GetOrCreate创建营养记录", userID)
	nutrition, err = s.nutritionDAO.GetOrCreate(userID, today, createParams)
	if err != nil {
		log.Printf("[营养服务] 用户(ID:%d)GetOrCreate失败: %v", userID, err)
		return nil, err
	}

	log.Printf("[营养服务] 用户(ID:%d)营养记录创建成功(ID:%d)", userID, nutrition.ID)
	return nutrition, nil
}

// UpdateTodayNutrition 更新今日营养摄入数据
func (s *NutritionService) UpdateTodayNutrition(userID int64, caloriesIntake, proteinIntakeG, carbIntakeG, fatIntakeG float64) (*model.DailyNutrition, error) {
	// 先获取今日营养记录
	nutrition, err := s.GetTodayNutrition(userID)
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
