package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ome-app-back/internal/service"
	"ome-app-back/pkg/api"
)

// DailyNutritionController 处理用户每日营养相关API
type DailyNutritionController struct {
	nutritionService *service.NutritionService
}

// NewDailyNutritionController 创建每日营养控制器
func NewDailyNutritionController(nutritionService *service.NutritionService) *DailyNutritionController {
	return &DailyNutritionController{
		nutritionService: nutritionService,
	}
}

// GetTodayNutrition 获取今日营养数据
func (c *DailyNutritionController) GetTodayNutrition(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		api.ResponseError(ctx, http.StatusUnauthorized, "未授权")
		return
	}

	nutrition, err := c.nutritionService.GetTodayNutrition(userID)
	if err != nil {
		if err.Error() == "用户尚未生成健康分析报告，请先生成健康分析" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    10001,
				"msg":     "请先生成健康分析",
				"data":    nil,
				"details": []string{"用户尚未生成健康分析报告，无法创建营养记录"},
			})
			return
		}
		api.ResponseError(ctx, http.StatusInternalServerError, "获取营养数据失败")
		return
	}

	api.ResponseSuccess(ctx, nutrition)
}

// UpdateNutritionInput 更新营养摄入的请求参数
type UpdateNutritionInput struct {
	CaloriesIntake float64 `json:"calories_intake"`
	ProteinIntakeG float64 `json:"protein_intake_g"`
	CarbIntakeG    float64 `json:"carb_intake_g"`
	FatIntakeG     float64 `json:"fat_intake_g"`
}

// UpdateTodayNutrition 更新今日营养摄入量
func (c *DailyNutritionController) UpdateTodayNutrition(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		api.ResponseError(ctx, http.StatusUnauthorized, "未授权")
		return
	}

	var input UpdateNutritionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ResponseError(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	nutrition, err := c.nutritionService.UpdateTodayNutrition(
		userID,
		input.CaloriesIntake,
		input.ProteinIntakeG,
		input.CarbIntakeG,
		input.FatIntakeG,
	)
	if err != nil {
		api.ResponseError(ctx, http.StatusInternalServerError, "更新营养数据失败")
		return
	}

	api.ResponseSuccess(ctx, nutrition)
}

// GetNutritionHistoryInput 获取历史记录的请求参数
type GetNutritionHistoryInput struct {
	StartDate string `form:"start_date" binding:"required"` // 格式: 2023-04-01
	EndDate   string `form:"end_date" binding:"required"`   // 格式: 2023-04-07
}

// GetNutritionHistory 获取营养历史记录
func (c *DailyNutritionController) GetNutritionHistory(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		api.ResponseError(ctx, http.StatusUnauthorized, "未授权")
		return
	}

	var input GetNutritionHistoryInput
	if err := ctx.ShouldBindQuery(&input); err != nil {
		api.ResponseError(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		api.ResponseError(ctx, http.StatusBadRequest, "开始日期格式错误")
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		api.ResponseError(ctx, http.StatusBadRequest, "结束日期格式错误")
		return
	}

	// 获取历史记录
	records, err := c.nutritionService.GetNutritionHistory(userID, startDate, endDate)
	if err != nil {
		api.ResponseError(ctx, http.StatusInternalServerError, "获取历史记录失败")
		return
	}

	api.ResponseSuccess(ctx, records)
}

// GetWeekSummary 获取一周营养摄入统计
func (c *DailyNutritionController) GetWeekSummary(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	if userID == 0 {
		api.ResponseError(ctx, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取一周数据统计
	summary, err := c.nutritionService.GetWeekSummary(userID)
	if err != nil {
		api.ResponseError(ctx, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	api.ResponseSuccess(ctx, summary)
}
