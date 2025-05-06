package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ome-app-back/internal/service"
)

// NutritionAPI 处理用户每日营养相关接口
type NutritionAPI struct {
	nutritionService *service.NutritionService
}

// NewNutritionAPI 创建营养API处理实例
func NewNutritionAPI(nutritionService *service.NutritionService) *NutritionAPI {
	return &NutritionAPI{
		nutritionService: nutritionService,
	}
}

// GetTodayNutrition 获取今日营养数据
func (a *NutritionAPI) GetTodayNutrition(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	nutrition, err := a.nutritionService.GetTodayNutrition(userID)
	if err != nil {
		// 特别处理用户没有健康分析报告的情况
		if errors.Is(err, service.ErrNoHealthAnalysis) {
			c.JSON(http.StatusOK, gin.H{
				"code":    10001, // 使用特定错误码标识需要健康分析
				"msg":     "请先生成健康分析",
				"data":    nil,
				"details": []string{"用户尚未生成健康分析报告，无法创建营养记录"},
			})
			return
		}

		// 判断错误类型，给出更具体的错误信息
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseError(c, http.StatusNotFound, "未找到营养数据记录")
			return
		}
		responseError(c, http.StatusInternalServerError, "获取营养数据失败", err.Error())
		return
	}

	responseSuccess(c, nutrition)
}

// UpdateNutritionInput 更新营养摄入的请求参数
type UpdateNutritionInput struct {
	CaloriesIntake float64 `json:"calories_intake"`
	ProteinIntakeG float64 `json:"protein_intake_g"`
	CarbIntakeG    float64 `json:"carb_intake_g"`
	FatIntakeG     float64 `json:"fat_intake_g"`
}

// UpdateTodayNutrition 更新今日营养摄入量
func (a *NutritionAPI) UpdateTodayNutrition(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var input UpdateNutritionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responseError(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	nutrition, err := a.nutritionService.UpdateTodayNutrition(
		userID,
		input.CaloriesIntake,
		input.ProteinIntakeG,
		input.CarbIntakeG,
		input.FatIntakeG,
	)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "更新营养数据失败", err.Error())
		return
	}

	responseSuccess(c, nutrition)
}

// GetNutritionHistoryInput 获取历史记录的请求参数
type GetNutritionHistoryInput struct {
	StartDate string `form:"start_date" binding:"required"` // 格式: 2023-04-01
	EndDate   string `form:"end_date" binding:"required"`   // 格式: 2023-04-07
}

// GetNutritionHistory 获取营养历史记录
func (a *NutritionAPI) GetNutritionHistory(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var input GetNutritionHistoryInput
	if err := c.ShouldBindQuery(&input); err != nil {
		responseError(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		responseError(c, http.StatusBadRequest, "开始日期格式错误")
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		responseError(c, http.StatusBadRequest, "结束日期格式错误")
		return
	}

	// 获取历史记录
	records, err := a.nutritionService.GetNutritionHistory(userID, startDate, endDate)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取历史记录失败", err.Error())
		return
	}

	responseSuccess(c, records)
}

// GetWeekSummary 获取一周营养摄入统计
func (a *NutritionAPI) GetWeekSummary(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取一周数据统计
	summary, err := a.nutritionService.GetWeekSummary(userID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取统计数据失败", err.Error())
		return
	}

	responseSuccess(c, summary)
}

// 辅助函数，获取用户ID
func getUserID(c *gin.Context) int64 {
	value, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	userID, ok := value.(int64)
	if !ok {
		return 0
	}

	return userID
}

// 辅助函数，返回成功响应
func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": data,
	})
}

// 辅助函数，返回错误响应
func responseError(c *gin.Context, code int, msg string, details ...string) {
	resp := gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	}

	if len(details) > 0 {
		resp["details"] = details
	}

	c.JSON(code, resp)
}
