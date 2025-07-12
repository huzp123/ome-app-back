package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"ome-app-back/repositories"
	"ome-app-back/models"
)

// FoodRecognitionService 处理食物识别相关服务
type FoodRecognitionService struct {
	recognitionDAO    *repositories.FoodRecognitionDAO
	nutritionDAO      *repositories.DailyNutritionDAO
	healthAnalysisDAO *repositories.HealthAnalysisDAO
	fileService       *FileService
	aiService         *AIService
}

func NewFoodRecognitionService(
	recognitionDAO *repositories.FoodRecognitionDAO,
	nutritionDAO *repositories.DailyNutritionDAO,
	healthAnalysisDAO *repositories.HealthAnalysisDAO,
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
	Foods     []models.RecognizedFoodItem     `json:"foods"`
	Nutrition models.FoodRecognitionNutrition `json:"nutrition"`
	Analysis  string                         `json:"analysis"`
}

// FoodRecognitionHistoryResult 食物识别历史结果
type FoodRecognitionHistoryResult struct {
	Total      int64                                    `json:"total"`       // 总记录数
	Records    []models.FoodRecognitionResult            `json:"records"`     // 记录列表
	DateGroups map[string][]models.FoodRecognitionResult `json:"date_groups"` // 按日期分组的记录
}

// RecognizeFood 处理食物识别
func (s *FoodRecognitionService) RecognizeFood(userID int64, sessionID string, file *multipart.FileHeader) (*models.FoodRecognitionResult, error) {
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

	// 计算总处理时间
	duration := time.Since(startTime)
	log.Printf("[食物识别] 处理完成, 总耗时: %.2f秒", duration.Seconds())

	return result, nil
}

// SaveRecognitionToNutrition 将食物识别结果保存到用户当日营养摄入
func (s *FoodRecognitionService) SaveRecognitionToNutrition(recognitionID int64, userID int64) error {
	log.Printf("[食物识别-保存] 开始将识别结果(ID:%d)保存到用户(ID:%d)的营养摄入", recognitionID, userID)

	// 获取识别记录
	recognition, err := s.recognitionDAO.GetRecognitionByID(recognitionID)
	if err != nil {
		log.Printf("[食物识别-保存] 错误: 获取识别记录失败: %v", err)
		return err
	}

	// 验证识别记录属于当前用户
	if recognition.UserID != userID {
		log.Printf("[食物识别-保存] 错误: 用户(ID:%d)尝试保存不属于自己的识别记录(ID:%d)", userID, recognitionID)
		return errors.New("无权操作此识别记录")
	}

	// 直接从识别记录中获取营养数据
	nutrition := models.FoodRecognitionNutrition{
		CaloriesIntake: recognition.CaloriesIntake,
		ProteinIntakeG: recognition.ProteinIntakeG,
		CarbIntakeG:    recognition.CarbIntakeG,
		FatIntakeG:     recognition.FatIntakeG,
	}

	// 更新用户当日营养摄入
	log.Printf("[食物识别-保存] 更新用户营养摄入数据...")
	if err := s.updateDailyNutrition(userID, nutrition); err != nil {
		log.Printf("[食物识别-保存] 错误: 更新营养数据失败: %v", err)
		return err
	}

	// 更新识别记录的采用状态
	log.Printf("[食物识别-保存] 更新识别记录采用状态...")
	if err := s.recognitionDAO.UpdateAdoptionStatus(recognitionID, true); err != nil {
		log.Printf("[食物识别-保存] 错误: 更新采用状态失败: %v", err)
		// 虽然状态更新失败，但营养数据已成功更新，所以仍然返回成功
		log.Printf("[食物识别-保存] 注意: 营养数据已成功更新，但状态更新失败")
		return nil
	}

	log.Printf("[食物识别-保存] 营养数据更新成功，记录状态已标记为已采用")
	return nil
}

// GetRecognitionByID 获取识别记录详情
func (s *FoodRecognitionService) GetRecognitionByID(id int64) (*models.FoodRecognitionResult, error) {
	recognition, err := s.recognitionDAO.GetRecognitionByID(id)
	if err != nil {
		return nil, err
	}

	return s.recognitionDAO.ConvertToResult(recognition)
}

// GetUserTodayRecognitions 获取用户今日的食物识别记录
func (s *FoodRecognitionService) GetUserTodayRecognitions(userID int64) ([]models.FoodRecognitionResult, error) {
	records, err := s.recognitionDAO.GetUserTodayRecognitions(userID)
	if err != nil {
		return nil, err
	}

	results := make([]models.FoodRecognitionResult, 0, len(records))
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
func (s *FoodRecognitionService) updateDailyNutrition(userID int64, nutrition models.FoodRecognitionNutrition) error {
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

// GetAdoptedRecognitions 获取用户已采用的食物识别记录
func (s *FoodRecognitionService) GetAdoptedRecognitions(userID int64, page, pageSize int, startDate, endDate string) (*FoodRecognitionHistoryResult, error) {
	log.Printf("[食物识别-历史] 查询用户(ID:%d)已采用的食物识别记录, 日期范围: %s - %s", userID, startDate, endDate)

	// 解析日期
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		log.Printf("[食物识别-历史] 错误: 开始日期格式无效: %v", err)
		return nil, fmt.Errorf("开始日期格式无效: %v", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		log.Printf("[食物识别-历史] 错误: 结束日期格式无效: %v", err)
		return nil, fmt.Errorf("结束日期格式无效: %v", err)
	}

	// 将结束日期设置为当天的23:59:59
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())

	// 获取分页记录
	records, total, err := s.recognitionDAO.GetUserAdoptedRecognitions(userID, page, pageSize, start, end)
	if err != nil {
		log.Printf("[食物识别-历史] 错误: 查询记录失败: %v", err)
		return nil, err
	}

	// 转换为前端结果对象
	recordResults := make([]models.FoodRecognitionResult, 0, len(records))
	for _, record := range records {
		result, err := s.recognitionDAO.ConvertToResult(&record)
		if err != nil {
			log.Printf("[食物识别-历史] 警告: 转换记录失败(ID:%d): %v", record.ID, err)
			continue
		}
		recordResults = append(recordResults, *result)
	}

	// 获取按日期分组的记录
	dateGroupRecords, err := s.recognitionDAO.GetUserAdoptedRecognitionsByDateRange(userID, start, end)
	if err != nil {
		log.Printf("[食物识别-历史] 错误: 查询日期分组记录失败: %v", err)
		return nil, err
	}

	// 转换为前端结果对象
	dateGroups := make(map[string][]models.FoodRecognitionResult)
	for date, groupRecords := range dateGroupRecords {
		results := make([]models.FoodRecognitionResult, 0, len(groupRecords))
		for _, record := range groupRecords {
			result, err := s.recognitionDAO.ConvertToResult(&record)
			if err != nil {
				log.Printf("[食物识别-历史] 警告: 转换日期组记录失败(ID:%d): %v", record.ID, err)
				continue
			}
			results = append(results, *result)
		}
		dateGroups[date] = results
	}

	// 构建结果
	historyResult := &FoodRecognitionHistoryResult{
		Total:      total,
		Records:    recordResults,
		DateGroups: dateGroups,
	}

	log.Printf("[食物识别-历史] 查询成功, 共%d条记录, %d个日期组", total, len(dateGroups))
	return historyResult, nil
}
