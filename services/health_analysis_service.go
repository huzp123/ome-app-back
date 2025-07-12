package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"ome-app-back/models"
	"ome-app-back/repositories"
)

// HealthAnalysisService 处理健康分析相关业务逻辑
type HealthAnalysisService struct {
	userDAO           *repositories.AppUserDAO
	userWeightDAO     *repositories.UserWeightDAO
	userHeightDAO     *repositories.UserHeightDAO
	userGoalDAO       *repositories.UserGoalDAO
	healthAnalysisDAO *repositories.HealthAnalysisDAO
}

// NewHealthAnalysisService 创建健康分析服务实例
func NewHealthAnalysisService(
	userDAO *repositories.AppUserDAO,
	userWeightDAO *repositories.UserWeightDAO,
	userHeightDAO *repositories.UserHeightDAO,
	userGoalDAO *repositories.UserGoalDAO,
	healthAnalysisDAO *repositories.HealthAnalysisDAO,
) *HealthAnalysisService {
	return &HealthAnalysisService{
		userDAO:           userDAO,
		userWeightDAO:     userWeightDAO,
		userHeightDAO:     userHeightDAO,
		userGoalDAO:       userGoalDAO,
		healthAnalysisDAO: healthAnalysisDAO,
	}
}

// AnalysisRequest 健康分析请求
type AnalysisRequest struct {
	UserID int64 `json:"user_id"`
}

// AnalysisResponse 健康分析响应
type AnalysisResponse struct {
	BMI                 float64 `json:"bmi"`
	BMICategory         string  `json:"bmi_category"`
	BMR                 float64 `json:"bmr"`
	TDEE                float64 `json:"tdee"`
	RecommendedCalories float64 `json:"recommended_calories"`
	ProteinNeedG        float64 `json:"protein_need_g"`
	CarbNeedG           float64 `json:"carb_need_g"`
	FatNeedG            float64 `json:"fat_need_g"`
	AnalysisContent     string  `json:"analysis_content"`
	CurrentWeightKG     float64 `json:"current_weight_kg"`
	TargetWeightKG      float64 `json:"target_weight_kg"`
	WeeklyChangeKG      float64 `json:"weekly_change_kg"`
	TargetDate          string  `json:"target_date"`
	DaysToTarget        int     `json:"days_to_target"`
}

// GenerateAnalysis 生成健康分析报告
func (s *HealthAnalysisService) GenerateAnalysis(req AnalysisRequest) (*AnalysisResponse, error) {
	// 获取用户基本信息
	user, err := s.userDAO.GetByID(req.UserID)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}

	// 检查用户信息是否完整
	if user.BirthDate.IsZero() || user.Sex == "" {
		return nil, errors.New("请先完善个人资料")
	}

	// 获取用户当前身高
	heightRecord, err := s.userHeightDAO.GetCurrentHeight(req.UserID)
	if err != nil {
		return nil, errors.New("获取用户身高信息失败")
	}
	if heightRecord == nil {
		return nil, errors.New("请先记录身高信息")
	}

	// 获取用户当前体重
	weightRecord, err := s.userWeightDAO.GetLatest(req.UserID)
	if err != nil {
		return nil, errors.New("获取用户体重信息失败")
	}

	// 获取用户目标
	goal, err := s.userGoalDAO.GetByUserID(req.UserID)
	if err != nil {
		return nil, errors.New("获取用户目标失败")
	}

	// 计算BMI
	bmi := calculateBMI(weightRecord.WeightKG, heightRecord.HeightCM)
	bmiCategory := getBMICategory(bmi)

	// 计算基础代谢率(BMR)
	bmr := calculateBMR(user.Sex, weightRecord.WeightKG, heightRecord.HeightCM, calculateAge(user.BirthDate))

	// 计算每日总能量消耗(TDEE)，假设轻度活动水平系数为1.375
	tdee := bmr * 1.375

	// keep_fit模式下强制每周变化为0
	weeklyChangeKG := goal.WeeklyChangeKG
	if goal.GoalType == "keep_fit" {
		weeklyChangeKG = 0
	}

	// 根据目标计算推荐热量
	recommendedCalories := calculateRecommendedCalories(tdee, goal.GoalType, weeklyChangeKG)

	// 计算营养素建议
	proteinNeedG, carbNeedG, fatNeedG := calculateNutrientNeeds(recommendedCalories, weightRecord.WeightKG, goal.GoalType)

	// 计算距离目标日期天数
	daysToTarget := int(math.Ceil(time.Until(goal.TargetDate).Hours() / 24))
	if daysToTarget < 0 {
		daysToTarget = 0
	}

	// 生成分析文本内容
	analysisContent := generateAnalysisContent(
		bmi, bmiCategory, bmr, tdee, recommendedCalories,
		weightRecord.WeightKG, goal.TargetWeightKG, weeklyChangeKG,
		daysToTarget, goal.GoalType, proteinNeedG, carbNeedG, fatNeedG,
	)

	// 保存分析结果到数据库
	analysis := &models.HealthAnalysis{
		UserID:              req.UserID,
		BMI:                 bmi,
		BMR:                 bmr,
		TDEE:                tdee,
		ProteinNeedG:        proteinNeedG,
		CarbNeedG:           carbNeedG,
		FatNeedG:            fatNeedG,
		RecommendedCalories: recommendedCalories,
		AnalysisContent:     analysisContent,
	}
	if err := s.healthAnalysisDAO.Create(analysis); err != nil {
		// 保存失败不影响返回结果
		fmt.Println("保存健康分析结果失败:", err)
	}

	// 组装返回结果
	return &AnalysisResponse{
		BMI:                 bmi,
		BMICategory:         bmiCategory,
		BMR:                 bmr,
		TDEE:                tdee,
		RecommendedCalories: recommendedCalories,
		ProteinNeedG:        proteinNeedG,
		CarbNeedG:           carbNeedG,
		FatNeedG:            fatNeedG,
		AnalysisContent:     analysisContent,
		CurrentWeightKG:     weightRecord.WeightKG,
		TargetWeightKG:      goal.TargetWeightKG,
		WeeklyChangeKG:      weeklyChangeKG, // 使用调整后的值
		TargetDate:          goal.TargetDate.Format("2006-01-02"),
		DaysToTarget:        daysToTarget,
	}, nil
}

// 计算BMI
func calculateBMI(weightKg float64, heightCm float64) float64 {
	heightM := heightCm / 100.0
	return math.Round((weightKg/(heightM*heightM))*10) / 10
}

// 获取BMI分类
func getBMICategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "偏瘦"
	case bmi < 24:
		return "正常"
	case bmi < 28:
		return "超重"
	default:
		return "肥胖"
	}
}

// 计算年龄
func calculateAge(birthDate time.Time) int {
	if birthDate.IsZero() {
		return 0
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}
	return age
}

// 计算基础代谢率(BMR)
func calculateBMR(sex string, weightKg float64, heightCm float64, age int) float64 {
	if sex == "male" {
		return 88.362 + (13.397 * weightKg) + (4.799 * heightCm) - (5.677 * float64(age))
	}
	// 女性
	return 447.593 + (9.247 * weightKg) + (3.098 * heightCm) - (4.330 * float64(age))
}

// 计算推荐热量
func calculateRecommendedCalories(tdee float64, goalType string, weeklyChangeKg float64) float64 {
	// 1kg脂肪约等于7700千卡
	dailyCalorieAdjust := (weeklyChangeKg * 7700) / 7

	switch goalType {
	case "lose_fat":
		return tdee - dailyCalorieAdjust
	case "gain_muscle":
		return tdee + dailyCalorieAdjust
	default:
		return tdee
	}
}

// 计算营养素需求
func calculateNutrientNeeds(calories float64, weightKg float64, goalType string) (protein, carb, fat float64) {
	switch goalType {
	case "lose_fat":
		// 高蛋白、适中碳水、低脂肪
		protein = weightKg * 2.0 // 每公斤体重2.0克蛋白质
		fat = weightKg * 0.8     // 每公斤体重0.8克脂肪
		// 剩余热量来自碳水
		proteinCalories := protein * 4 // 蛋白质每克4千卡
		fatCalories := fat * 9         // 脂肪每克9千卡
		carbCalories := calories - proteinCalories - fatCalories
		carb = carbCalories / 4 // 碳水每克4千卡
	case "gain_muscle":
		// 高蛋白、高碳水、适中脂肪
		protein = weightKg * 2.2 // 每公斤体重2.2克蛋白质
		fat = weightKg * 1.0     // 每公斤体重1.0克脂肪
		// 剩余热量来自碳水
		proteinCalories := protein * 4 // 蛋白质每克4千卡
		fatCalories := fat * 9         // 脂肪每克9千卡
		carbCalories := calories - proteinCalories - fatCalories
		carb = carbCalories / 4 // 碳水每克4千卡
	default: // keep_fit
		// 平衡蛋白质、碳水和脂肪
		protein = weightKg * 1.8 // 每公斤体重1.8克蛋白质
		fat = weightKg * 1.0     // 每公斤体重1.0克脂肪
		// 剩余热量来自碳水
		proteinCalories := protein * 4 // 蛋白质每克4千卡
		fatCalories := fat * 9         // 脂肪每克9千卡
		carbCalories := calories - proteinCalories - fatCalories
		carb = carbCalories / 4 // 碳水每克4千卡
	}

	// 确保值为正数并四舍五入到整数
	protein = math.Max(0, math.Round(protein))
	carb = math.Max(0, math.Round(carb))
	fat = math.Max(0, math.Round(fat))

	return protein, carb, fat
}

// 生成分析文本内容
func generateAnalysisContent(
	bmi float64, bmiCategory string, bmr float64, tdee float64, recommendedCalories float64,
	currentWeight float64, targetWeight float64, weeklyChange float64,
	daysToTarget int, goalType string, protein float64, carb float64, fat float64,
) string {
	var goalDesc string
	switch goalType {
	case "lose_fat":
		goalDesc = "减脂"
	case "gain_muscle":
		goalDesc = "增肌"
	default:
		goalDesc = "保持体型"
	}

	content := fmt.Sprintf(
		"根据您的身体数据，您的BMI指数为%.1f，属于%s范围。\n\n"+
			"您的基础代谢率(BMR)为%.0f千卡，每日总能量消耗(TDEE)约为%.0f千卡。\n\n"+
			"基于您的%s目标，建议每日摄入%.0f千卡的热量",
		bmi, bmiCategory, bmr, tdee, goalDesc, recommendedCalories,
	)

	// 如果有体重变化计划，则显示
	if weeklyChange != 0 {
		changeWord := "增加"
		if weeklyChange < 0 {
			changeWord = "减少"
			weeklyChange = -weeklyChange // 转为正数显示
		}
		content += fmt.Sprintf("，计划每周%s%.1fkg体重", changeWord, weeklyChange)
	}

	content += fmt.Sprintf("。\n\n营养素建议摄入量：\n- 蛋白质：%.0fg\n- 碳水化合物：%.0fg\n- 脂肪：%.0fg\n\n", protein, carb, fat)

	// keep_fit模式不显示体重变化相关内容，因为weeklyChange已经固定为0
	if goalType != "keep_fit" && targetWeight != currentWeight {
		weightDiff := math.Abs(targetWeight - currentWeight)
		changeWord := "减少"
		if targetWeight > currentWeight {
			changeWord = "增加"
		}
		content += fmt.Sprintf(
			"您的目标是从%.1fkg%s到%.1fkg，总共需要%s%.1fkg。",
			currentWeight, changeWord, targetWeight,
			changeWord, weightDiff,
		)

		if daysToTarget > 0 && weeklyChange != 0 {
			changeWord := "减少"
			if targetWeight > currentWeight {
				changeWord = "增加"
			}

			content += fmt.Sprintf("按照每周%s%.1fkg的速度，还需要约%d天可达成目标。",
				changeWord, weeklyChange, daysToTarget)
		}
	} else if goalType == "keep_fit" {
		content += "您的目标是保持当前体型，建议维持均衡的饮食和规律的运动。"
	}

	return content
}

// GetHistoryAnalysis 获取用户健康分析历史记录
func (s *HealthAnalysisService) GetHistoryAnalysis(userID int64, limit int) ([]models.HealthAnalysis, error) {
	if limit <= 0 {
		limit = 10 // 默认返回10条记录
	}
	return s.healthAnalysisDAO.GetHistory(userID, limit)
}
