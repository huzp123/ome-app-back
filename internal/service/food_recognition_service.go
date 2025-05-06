package service

import (
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"ome-app-back/internal/dao"
	"ome-app-back/internal/model"
)

// FoodRecognitionService 处理食物识别相关服务
type FoodRecognitionService struct {
	recognitionDAO    *dao.FoodRecognitionDAO
	nutritionDAO      *dao.DailyNutritionDAO
	healthAnalysisDAO *dao.HealthAnalysisDAO
	fileService       *FileService
	aiService         *AIService
}

// NewFoodRecognitionService 创建食物识别服务实例
func NewFoodRecognitionService(
	recognitionDAO *dao.FoodRecognitionDAO,
	nutritionDAO *dao.DailyNutritionDAO,
	healthAnalysisDAO *dao.HealthAnalysisDAO,
	fileService *FileService,
	aiService *AIService,
) *FoodRecognitionService {
	return &FoodRecognitionService{
		recognitionDAO:    recognitionDAO,
		nutritionDAO:      nutritionDAO,
		healthAnalysisDAO: healthAnalysisDAO,
		fileService:       fileService,
		aiService:         aiService,
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
	log.Printf("[食物识别] 开始处理用户(ID:%d)的识别请求", userID)
	startTime := time.Now()

	// 上传图片
	log.Printf("[食物识别] 上传图片: %s", file.Filename)
	imagePath, err := s.fileService.UploadImage(file, userID)
	if err != nil {
		log.Printf("[食物识别] 错误: 图片上传失败: %v", err)
		return nil, err
	}
	log.Printf("[食物识别] 图片已上传至: %s", imagePath)

	imageBase64, err := s.fileService.GetImageBase64(imagePath)
	if err != nil {
		log.Printf("[食物识别] 错误: 图片转换失败: %v", err)
		return nil, err
	}
	log.Printf("[食物识别] 图片Base64转换完成, 大小: %d字节", len(imageBase64))

	// 调用AI分析
	aiResponse, err := s.aiService.AnalyzeImageWithAI(imageBase64, "请分析这张食物图片的营养成分")
	if err != nil {
		log.Printf("[食物识别] 错误: AI分析失败: %v", err)
		return nil, err
	}

	log.Printf("[食物识别] 解析AI响应...")
	var analysisResult AIAnalysisResult
	if err := json.Unmarshal([]byte(aiResponse), &analysisResult); err != nil {
		log.Printf("[食物识别] 错误: 解析AI响应失败: %v, 原始响应: %s", err, truncateString(aiResponse, 100))
		return nil, errors.New("解析AI响应失败: " + err.Error())
	}

	// 创建识别记录
	log.Printf("[食物识别] 创建数据库记录, 识别到%d种食物", len(analysisResult.Foods))
	recognition, err := s.recognitionDAO.CreateRecognition(
		userID,
		sessionID,
		imagePath,
		analysisResult.Foods,
		analysisResult.Nutrition,
		analysisResult.Analysis,
	)
	if err != nil {
		log.Printf("[食物识别] 错误: 创建识别记录失败: %v", err)
		return nil, err
	}
	log.Printf("[食物识别] 识别记录已保存(ID:%d)", recognition.ID)

	// 转换为前端结果对象
	result, err := s.recognitionDAO.ConvertToResult(recognition)
	if err != nil {
		log.Printf("[食物识别] 错误: 转换结果对象失败: %v", err)
		return nil, err
	}

	// 更新用户当日营养摄入
	log.Printf("[食物识别] 更新用户营养摄入数据...")
	if err := s.updateDailyNutrition(userID, analysisResult.Nutrition); err != nil {
		// 记录错误但不中断流程
		log.Printf("[食物识别] 警告: 更新营养数据失败，但不影响识别结果: %v", err)
	} else {
		log.Printf("[食物识别] 营养数据更新成功")
	}

	// 计算总处理时间
	duration := time.Since(startTime)
	log.Printf("[食物识别] 处理完成, 总耗时: %.2f秒", duration.Seconds())

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
	log.Printf("[食物识别-营养] 开始更新用户(ID:%d)营养摄入", userID)

	// 创建营养服务实例
	nutritionService := NewNutritionService(s.nutritionDAO, s.healthAnalysisDAO)

	// 通过服务层获取今日营养数据
	dailyNutrition, err := nutritionService.GetTodayNutrition(userID)
	if err != nil {
		if err.Error() == "用户尚未生成健康分析报告，请先生成健康分析" {
			log.Printf("[食物识别-营养] 错误: 用户尚未生成健康分析报告")
		} else {
			log.Printf("[食物识别-营养] 错误: 获取今日营养数据失败: %v", err)
		}
		return err
	}

	// 记录原始值
	log.Printf("[食物识别-营养] 更新前: 热量=%.2f, 蛋白质=%.2f, 碳水=%.2f, 脂肪=%.2f",
		dailyNutrition.CaloriesIntake, dailyNutrition.ProteinIntakeG,
		dailyNutrition.CarbIntakeG, dailyNutrition.FatIntakeG)

	// 累加营养数据
	dailyNutrition.CaloriesIntake += nutrition.CaloriesIntake
	dailyNutrition.ProteinIntakeG += nutrition.ProteinIntakeG
	dailyNutrition.CarbIntakeG += nutrition.CarbIntakeG
	dailyNutrition.FatIntakeG += nutrition.FatIntakeG

	// 记录新值
	log.Printf("[食物识别-营养] 更新后: 热量=%.2f, 蛋白质=%.2f, 碳水=%.2f, 脂肪=%.2f",
		dailyNutrition.CaloriesIntake, dailyNutrition.ProteinIntakeG,
		dailyNutrition.CarbIntakeG, dailyNutrition.FatIntakeG)

	// 保存更新
	err = s.nutritionDAO.Update(dailyNutrition)
	if err != nil {
		log.Printf("[食物识别-营养] 错误: 保存营养数据失败: %v", err)
		return err
	}

	log.Printf("[食物识别-营养] 营养数据更新成功(ID:%d)", dailyNutrition.ID)
	return nil
}

// truncateString 截断字符串，用于日志输出
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "...(已截断)"
}
