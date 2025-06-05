package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/internal/service"
	"ome-app-back/pkg/errcode"
)

// MoodAPI 心情API
type MoodAPI struct {
	moodService *service.MoodService
}

// NewMoodAPI 创建心情API实例
func NewMoodAPI(moodService *service.MoodService) *MoodAPI {
	return &MoodAPI{
		moodService: moodService,
	}
}

// CreateMood 创建心情记录
func (api *MoodAPI) CreateMood(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req service.CreateMoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	mood, err := api.moodService.CreateMood(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": mood,
	})
}

// GetMood 获取单个心情记录
func (api *MoodAPI) GetMood(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	moodIDStr := c.Param("id")
	moodID, err := strconv.ParseInt(moodIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("心情记录ID格式错误").Response(c)
		return
	}

	mood, err := api.moodService.GetMood(userID, moodID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": mood,
	})
}

// GetMoodHistory 获取心情历史记录
func (api *MoodAPI) GetMoodHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req service.MoodHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	moods, err := api.moodService.GetMoodHistory(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": moods,
	})
}

// GetTodayMoods 获取今日心情记录
func (api *MoodAPI) GetTodayMoods(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	moods, err := api.moodService.GetTodayMoods(userID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": moods,
	})
}

// DeleteMood 删除心情记录
func (api *MoodAPI) DeleteMood(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	moodIDStr := c.Param("id")
	moodID, err := strconv.ParseInt(moodIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("心情记录ID格式错误").Response(c)
		return
	}

	err = api.moodService.DeleteMood(userID, moodID)
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

// GetMoodStatistics 获取心情统计数据
func (api *MoodAPI) GetMoodStatistics(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		errcode.InvalidParams.WithDetails("请提供开始日期和结束日期").Response(c)
		return
	}

	stats, err := api.moodService.GetMoodStatistics(userID, startDate, endDate)
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

// GetMoodOptions 获取心情选项
func (api *MoodAPI) GetMoodOptions(c *gin.Context) {
	options := api.moodService.GetMoodOptions()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": options,
	})
}
