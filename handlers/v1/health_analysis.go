package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/services"
	"ome-app-back/pkg/errcode"
)

type HealthAnalysisAPI struct {
	healthAnalysisService *services.HealthAnalysisService
}

func NewHealthAnalysisAPI(healthAnalysisService *services.HealthAnalysisService) *HealthAnalysisAPI {
	return &HealthAnalysisAPI{healthAnalysisService: healthAnalysisService}
}

// GenerateAnalysis 生成健康分析报告
func (api *HealthAnalysisAPI) GenerateAnalysis(c *gin.Context) {
	// 从JWT中获取用户ID（示例，实际项目中应从认证中间件获取）
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	req := services.AnalysisRequest{
		UserID: userID,
	}

	resp, err := api.healthAnalysisService.GenerateAnalysis(req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": resp,
	})
}

// GetHistoryAnalysis 获取用户健康分析历史记录
func (api *HealthAnalysisAPI) GetHistoryAnalysis(c *gin.Context) {
	// 从JWT中获取用户ID（示例，实际项目中应从认证中间件获取）
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	limitStr := c.Query("limit")
	limit := 10 // 默认返回10条记录
	if limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	analyses, err := api.healthAnalysisService.GetHistoryAnalysis(userID, limit)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": analyses,
	})
}
