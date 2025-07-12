package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/pkg/errcode"
	"ome-app-back/services"
)

type HeightAPI struct {
	heightService *services.HeightService
}

func NewHeightAPI(heightService *services.HeightService) *HeightAPI {
	return &HeightAPI{heightService: heightService}
}

// CreateHeight 创建身高记录
func (api *HeightAPI) CreateHeight(c *gin.Context) {
	var req services.CreateHeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	err := api.heightService.CreateHeight(userID, req)
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

// GetHeightHistory 获取身高历史记录
func (api *HeightAPI) GetHeightHistory(c *gin.Context) {
	var req services.GetHeightHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	history, err := api.heightService.GetHeightHistory(userID, req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": history,
	})
}

// GetCurrentHeight 获取当前身高
func (api *HeightAPI) GetCurrentHeight(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	height, err := api.heightService.GetCurrentHeight(userID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": height,
	})
}

// DeleteHeight 删除身高记录
func (api *HeightAPI) DeleteHeight(c *gin.Context) {
	heightIDStr := c.Param("id")
	heightID, err := strconv.ParseInt(heightIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("无效的身高记录ID").Response(c)
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	err = api.heightService.DeleteHeight(userID, heightID)
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

// GetHeightStatistics 获取身高统计数据
func (api *HeightAPI) GetHeightStatistics(c *gin.Context) {
	var req services.GetHeightStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	statistics, err := api.heightService.GetHeightStatistics(userID, req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": statistics,
	})
}
