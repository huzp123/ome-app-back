package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/services"
	"ome-app-back/pkg/errcode"
)

// WeightAPI 体重API
type WeightAPI struct {
	weightService *services.WeightService
}

// NewWeightAPI 创建体重API实例
func NewWeightAPI(weightService *services.WeightService) *WeightAPI {
	return &WeightAPI{
		weightService: weightService,
	}
}

// CreateWeight 手动记录体重
func (api *WeightAPI) CreateWeight(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req services.CreateWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	err := api.weightService.CreateWeight(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": nil,
	})
}

// GetWeightHistory 获取体重历史记录
func (api *WeightAPI) GetWeightHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req services.WeightHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	weights, err := api.weightService.GetWeightHistory(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": weights,
	})
}

// GetCurrentWeight 获取当前体重信息
func (api *WeightAPI) GetCurrentWeight(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	currentWeight, err := api.weightService.GetCurrentWeight(userID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": currentWeight,
	})
}

// DeleteWeight 删除体重记录
func (api *WeightAPI) DeleteWeight(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("体重记录ID格式错误").Response(c)
		return
	}

	err = api.weightService.DeleteWeight(userID, recordID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": nil,
	})
}

// GetWeightStatistics 获取体重统计分析
func (api *WeightAPI) GetWeightStatistics(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req services.WeightStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	stats, err := api.weightService.GetWeightStatistics(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": stats,
	})
}
