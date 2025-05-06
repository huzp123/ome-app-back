package service

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"time"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// FoodRecognitionService 处理食物识别相关服务
type FoodRecognitionService struct {
	recognitionDAO *dao.FoodRecognitionDAO
	nutritionDAO   *dao.DailyNutritionDAO
	fileService    *FileService
	aiService      *AIService
}

// NewFoodRecognitionService 创建食物识别服务实例
func NewFoodRecognitionService(
	recognitionDAO *dao.FoodRecognitionDAO,
	nutritionDAO *dao.DailyNutritionDAO,
	fileService *FileService,
	aiService *AIService,
) *FoodRecognitionService {
	return &FoodRecognitionService{
		recognitionDAO: recognitionDAO,
		nutritionDAO:   nutritionDAO,
		fileService:    fileService,
		aiService:      aiService,
	}
}

// AIAnalysisResult AI分析结果结构
type AIAnalysisResult struct {
	Foods     []model.RecognizedFoodItem     `json:"foods"`
	Nutrition model.FoodRecognitionNutrition `json:"nutrition"`
	Analysis  string                         `json:"analysis"`
}

// RecognizeFood 处理食物识别
func (s *FoodRecognitionService) RecognizeFood(userID int64, sessionID string, file *multipart.FileHeader) (*model.FoodRecognitionResult, error) {
	// 上传图片
	imagePath, err := s.fileService.UploadImage(file, userID)
	if err != nil {
		return nil, err
	}

	// 转换图片为Base64
	imageBase64, err := s.fileService.GetImageBase64(imagePath)
	if err != nil {
		return nil, err
	}

	// 调用AI分析
	aiResponse, err := s.aiService.AnalyzeImageWithAI(imageBase64, "请分析这张食物图片的营养成分")
	if err != nil {
		return nil, err
	}

	// 解析AI返回的JSON
	var analysisResult AIAnalysisResult
	if err := json.Unmarshal([]byte(aiResponse), &analysisResult); err != nil {
		return nil, errors.New("解析AI响应失败: " + err.Error())
	}

	// 创建识别记录
	recognition, err := s.recognitionDAO.CreateRecognition(
		userID,
		sessionID,
		imagePath,
		analysisResult.Foods,
		analysisResult.Nutrition,
		analysisResult.Analysis,
	)
	if err != nil {
		return nil, err
	}

	// 转换为前端结果对象
	result, err := s.recognitionDAO.ConvertToResult(recognition)
	if err != nil {
		return nil, err
	}

	// 更新用户当日营养摄入
	if err := s.updateDailyNutrition(userID, analysisResult.Nutrition); err != nil {
		// 记录错误但不中断流程
		// TODO: 添加日志记录
	}

	return result, nil
}

// GetRecognitionByID 获取识别记录详情
func (s *FoodRecognitionService) GetRecognitionByID(id int64) (*model.FoodRecognitionResult, error) {
	recognition, err := s.recognitionDAO.GetRecognitionByID(id)
	if err != nil {
		return nil, err
	}

	return s.recognitionDAO.ConvertToResult(recognition)
}

// GetUserTodayRecognitions 获取用户今日的食物识别记录
func (s *FoodRecognitionService) GetUserTodayRecognitions(userID int64) ([]model.FoodRecognitionResult, error) {
	records, err := s.recognitionDAO.GetUserTodayRecognitions(userID)
	if err != nil {
		return nil, err
	}

	results := make([]model.FoodRecognitionResult, 0, len(records))
	for _, record := range records {
		result, err := s.recognitionDAO.ConvertToResult(&record)
		if err != nil {
			continue // 跳过有错误的记录
		}
		results = append(results, *result)
	}

	return results, nil
}

// updateDailyNutrition 更新用户当日营养摄入数据
func (s *FoodRecognitionService) updateDailyNutrition(userID int64, nutrition model.FoodRecognitionNutrition) error {
	// 获取当日记录
	today := time.Now()
	dailyNutrition, err := s.nutritionDAO.GetOrCreate(userID, today)
	if err != nil {
		return err
	}

	// 累加营养数据
	dailyNutrition.CaloriesIntake += nutrition.CaloriesIntake
	dailyNutrition.ProteinIntakeG += nutrition.ProteinIntakeG
	dailyNutrition.CarbIntakeG += nutrition.CarbIntakeG
	dailyNutrition.FatIntakeG += nutrition.FatIntakeG

	// 保存更新
	return s.nutritionDAO.Update(dailyNutrition)
}
