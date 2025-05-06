package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ome-app-back/internal/service"
)

// FoodRecognitionAPI 处理食物识别相关接口
type FoodRecognitionAPI struct {
	recognitionService *service.FoodRecognitionService
}

// NewFoodRecognitionAPI 创建食物识别API处理实例
func NewFoodRecognitionAPI(recognitionService *service.FoodRecognitionService) *FoodRecognitionAPI {
	return &FoodRecognitionAPI{
		recognitionService: recognitionService,
	}
}

// RecognizeFood 分析上传的食物图片
func (a *FoodRecognitionAPI) RecognizeFood(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取会话ID
	sessionID := c.PostForm("session_id")

	// 获取上传的文件
	file, err := c.FormFile("food_image")
	if err != nil {
		responseError(c, http.StatusBadRequest, "获取上传文件失败")
		return
	}

	// 检查文件大小
	if file.Size > 10*1024*1024 { // 限制10MB
		responseError(c, http.StatusBadRequest, "文件大小超过限制")
		return
	}

	// 调用服务处理识别
	result, err := a.recognitionService.RecognizeFood(userID, sessionID, file)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "识别食物失败: "+err.Error())
		return
	}

	responseSuccess(c, result)
}

// GetRecognitionByID 获取食物识别记录详情
func (a *FoodRecognitionAPI) GetRecognitionByID(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取识别记录ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseError(c, http.StatusBadRequest, "无效的ID参数")
		return
	}

	// 获取记录详情
	result, err := a.recognitionService.GetRecognitionByID(id)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取识别记录失败")
		return
	}

	responseSuccess(c, result)
}

// GetTodayRecognitions 获取用户今日的食物识别记录
func (a *FoodRecognitionAPI) GetTodayRecognitions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		responseError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 获取当日记录
	results, err := a.recognitionService.GetUserTodayRecognitions(userID)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "获取识别记录失败")
		return
	}

	responseSuccess(c, results)
}
