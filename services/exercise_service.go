package services

import (
	"errors"
	"time"

	"ome-app-back/repositories"
	"ome-app-back/models"
	"ome-app-back/models/constant"
)

// ExerciseService 运动服务
type ExerciseService struct {
	exerciseDAO *repositories.UserExerciseDAO
}

// NewExerciseService 创建运动服务实例
func NewExerciseService(exerciseDAO *repositories.UserExerciseDAO) *ExerciseService {
	return &ExerciseService{
		exerciseDAO: exerciseDAO,
	}
}

// 请求/响应结构

// CreateExerciseRequest 创建运动记录请求
type CreateExerciseRequest struct {
	ExerciseType   string   `json:"exercise_type" binding:"required"`
	DurationMin    float64  `json:"duration_min" binding:"required,gt=0"`
	CaloriesBurned float64  `json:"calories_burned" binding:"gte=0"`
	DistanceKM     *float64 `json:"distance_km,omitempty"`
	StartTime      string   `json:"start_time" binding:"required"` // 格式: "2023-12-01T10:30:00Z"
}

// UpdateExerciseRequest 更新运动记录请求
type UpdateExerciseRequest struct {
	ExerciseType   string   `json:"exercise_type,omitempty"`
	DurationMin    *float64 `json:"duration_min,omitempty"`
	CaloriesBurned *float64 `json:"calories_burned,omitempty"`
	DistanceKM     *float64 `json:"distance_km,omitempty"`
	StartTime      string   `json:"start_time,omitempty"`
}

// ExerciseHistoryRequest 获取运动历史请求
type ExerciseHistoryRequest struct {
	StartDate string `form:"start_date" binding:"required"` // 格式: "2023-12-01"
	EndDate   string `form:"end_date" binding:"required"`   // 格式: "2023-12-31"
	Limit     int    `form:"limit"`
}

// 服务方法

// CreateExercise 创建运动记录
func (s *ExerciseService) CreateExercise(userID int64, req *CreateExerciseRequest) (*models.UserExercise, error) {
	// 验证运动类型
	if !constant.ValidExerciseTypesMap[req.ExerciseType] {
		return nil, errors.New("无效的运动类型: " + req.ExerciseType)
	}

	// 解析时间
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errors.New("时间格式错误，请使用 RFC3339 格式")
	}

	exercise := &models.UserExercise{
		UserID:         userID,
		ExerciseType:   req.ExerciseType,
		DurationMin:    req.DurationMin,
		CaloriesBurned: req.CaloriesBurned,
		DistanceKM:     req.DistanceKM,
		StartTime:      startTime,
	}

	err = s.exerciseDAO.Create(exercise)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}

// GetExercise 获取单个运动记录
func (s *ExerciseService) GetExercise(userID, exerciseID int64) (*models.UserExercise, error) {
	return s.exerciseDAO.GetByID(userID, exerciseID)
}

// GetExerciseHistory 获取运动历史记录
func (s *ExerciseService) GetExerciseHistory(userID int64, req *ExerciseHistoryRequest) ([]models.UserExercise, error) {
	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，请使用 YYYY-MM-DD 格式")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，请使用 YYYY-MM-DD 格式")
	}

	// 将结束日期设置为当天的最后一秒
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return s.exerciseDAO.GetHistory(userID, startDate, endDate, req.Limit)
}

// GetTodayExercises 获取今日运动记录
func (s *ExerciseService) GetTodayExercises(userID int64) ([]models.UserExercise, error) {
	return s.exerciseDAO.GetTodayExercises(userID)
}

// UpdateExercise 更新运动记录
func (s *ExerciseService) UpdateExercise(userID, exerciseID int64, req *UpdateExerciseRequest) (*models.UserExercise, error) {
	// 先获取现有记录
	exercise, err := s.exerciseDAO.GetByID(userID, exerciseID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.ExerciseType != "" {
		// 验证运动类型
		if !constant.ValidExerciseTypesMap[req.ExerciseType] {
			return nil, errors.New("无效的运动类型: " + req.ExerciseType)
		}
		exercise.ExerciseType = req.ExerciseType
	}
	if req.DurationMin != nil {
		exercise.DurationMin = *req.DurationMin
	}
	if req.CaloriesBurned != nil {
		exercise.CaloriesBurned = *req.CaloriesBurned
	}
	if req.DistanceKM != nil {
		exercise.DistanceKM = req.DistanceKM
	}
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return nil, errors.New("时间格式错误，请使用 RFC3339 格式")
		}
		exercise.StartTime = startTime
	}

	err = s.exerciseDAO.Update(exercise)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}

// DeleteExercise 删除运动记录
func (s *ExerciseService) DeleteExercise(userID, exerciseID int64) error {
	return s.exerciseDAO.Delete(userID, exerciseID)
}

// GetExerciseStatistics 获取运动统计数据
func (s *ExerciseService) GetExerciseStatistics(userID int64, startDate, endDate string) (map[string]interface{}, error) {
	// 解析日期
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误，请使用 YYYY-MM-DD 格式")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误，请使用 YYYY-MM-DD 格式")
	}

	// 将结束日期设置为当天的最后一秒
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return s.exerciseDAO.GetStatistics(userID, start, end)
}

// GetExerciseOptions 获取运动选项（用于前端显示）
func (s *ExerciseService) GetExerciseOptions() map[string]interface{} {
	return map[string]interface{}{
		"exercise_types": constant.ExerciseTypes,
	}
}
