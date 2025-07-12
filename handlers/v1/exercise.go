package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/services"
	"ome-app-back/pkg/errcode"
)

// ExerciseAPI 运动API
type ExerciseAPI struct {
	exerciseService *services.ExerciseService
}

// NewExerciseAPI 创建运动API实例
func NewExerciseAPI(exerciseService *services.ExerciseService) *ExerciseAPI {
	return &ExerciseAPI{
		exerciseService: exerciseService,
	}
}

// CreateExercise 创建运动记录
func (api *ExerciseAPI) CreateExercise(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req services.CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	exercise, err := api.exerciseService.CreateExercise(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": exercise,
	})
}

// GetExercise 获取单个运动记录
func (api *ExerciseAPI) GetExercise(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("运动记录ID格式错误").Response(c)
		return
	}

	exercise, err := api.exerciseService.GetExercise(userID, exerciseID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": exercise,
	})
}

// GetExerciseHistory 获取运动历史记录
func (api *ExerciseAPI) GetExerciseHistory(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	var req services.ExerciseHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	exercises, err := api.exerciseService.GetExerciseHistory(userID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": exercises,
	})
}

// GetTodayExercises 获取今日运动记录
func (api *ExerciseAPI) GetTodayExercises(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	exercises, err := api.exerciseService.GetTodayExercises(userID)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": exercises,
	})
}

// UpdateExercise 更新运动记录
func (api *ExerciseAPI) UpdateExercise(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("运动记录ID格式错误").Response(c)
		return
	}

	var req services.UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.InvalidParams.WithDetails(err.Error()).Response(c)
		return
	}

	exercise, err := api.exerciseService.UpdateExercise(userID, exerciseID, &req)
	if err != nil {
		errcode.ServerError.WithDetails(err.Error()).Response(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": exercise,
	})
}

// DeleteExercise 删除运动记录
func (api *ExerciseAPI) DeleteExercise(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		errcode.UnauthorizedTokenError.Response(c)
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		errcode.InvalidParams.WithDetails("运动记录ID格式错误").Response(c)
		return
	}

	err = api.exerciseService.DeleteExercise(userID, exerciseID)
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

// GetExerciseStatistics 获取运动统计数据
func (api *ExerciseAPI) GetExerciseStatistics(c *gin.Context) {
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

	stats, err := api.exerciseService.GetExerciseStatistics(userID, startDate, endDate)
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

// GetExerciseOptions 获取运动选项配置
func (api *ExerciseAPI) GetExerciseOptions(c *gin.Context) {
	options := api.exerciseService.GetExerciseOptions()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "成功",
		"data": options,
	})
}
