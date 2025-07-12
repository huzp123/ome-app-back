package services

import (
	"errors"
	"time"

	"ome-app-back/repositories"
	"ome-app-back/models"
	"ome-app-back/models/constant"
)

// MoodService 心情服务
type MoodService struct {
	moodDAO *repositories.MoodRecordDAO
}

// NewMoodService 创建心情服务实例
func NewMoodService(moodDAO *repositories.MoodRecordDAO) *MoodService {
	return &MoodService{
		moodDAO: moodDAO,
	}
}

// 请求/响应结构

// CreateMoodRequest 创建心情记录请求
type CreateMoodRequest struct {
	TimeContext string   `json:"time_context" binding:"required,oneof=now today"` // "now" 或 "today"
	MoodLevel   int      `json:"mood_level" binding:"required,min=1,max=7"`       // 1-7级
	MoodTags    []string `json:"mood_tags,omitempty"`                             // 情绪标签数组，可选
	Influences  []string `json:"influences,omitempty"`                            // 影响因素数组，可选
}

// MoodHistoryRequest 获取心情历史请求
type MoodHistoryRequest struct {
	StartDate string `form:"start_date" binding:"required"` // 格式: "2023-12-01"
	EndDate   string `form:"end_date" binding:"required"`   // 格式: "2023-12-31"
	Limit     int    `form:"limit"`
}

// 服务方法

// CreateMood 创建心情记录
func (s *MoodService) CreateMood(userID int64, req *CreateMoodRequest) (*models.MoodRecord, error) {
	// 验证情绪等级
	if req.MoodLevel < constant.MoodLevelMin || req.MoodLevel > constant.MoodLevelMax {
		return nil, errors.New("情绪等级必须在1-7之间")
	}

	// 验证时间上下文
	if req.TimeContext != constant.TimeContextNow && req.TimeContext != constant.TimeContextToday {
		return nil, errors.New("时间上下文必须是 'now' 或 'today'")
	}

	// 验证情绪标签（可选，如果提供则验证）
	for _, tag := range req.MoodTags {
		if tag == "" {
			return nil, errors.New("情绪标签不能为空")
		}
	}

	// 验证影响因素（可选，如果提供则验证）
	for _, influence := range req.Influences {
		if !constant.ValidInfluencesMap[influence] {
			return nil, errors.New("无效的影响因素: " + influence)
		}
	}

	mood := &models.MoodRecord{
		UserID:      userID,
		TimeContext: req.TimeContext,
		MoodLevel:   req.MoodLevel,
		MoodTags:    req.MoodTags,
		Influences:  req.Influences,
		RecordTime:  time.Now(),
	}

	err := s.moodDAO.Create(mood)
	if err != nil {
		return nil, err
	}

	return mood, nil
}

// GetMood 获取单个心情记录
func (s *MoodService) GetMood(userID, moodID int64) (*models.MoodRecord, error) {
	return s.moodDAO.GetByID(userID, moodID)
}

// GetMoodHistory 获取心情历史记录
func (s *MoodService) GetMoodHistory(userID int64, req *MoodHistoryRequest) ([]models.MoodRecord, error) {
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

	return s.moodDAO.GetHistory(userID, startDate, endDate, req.Limit)
}

// GetTodayMoods 获取今日心情记录
func (s *MoodService) GetTodayMoods(userID int64) ([]models.MoodRecord, error) {
	return s.moodDAO.GetTodayMoods(userID)
}

// DeleteMood 删除心情记录
func (s *MoodService) DeleteMood(userID, moodID int64) error {
	return s.moodDAO.Delete(userID, moodID)
}

// GetMoodStatistics 获取心情统计数据
func (s *MoodService) GetMoodStatistics(userID int64, startDate, endDate string) (map[string]interface{}, error) {
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

	return s.moodDAO.GetMoodStatistics(userID, start, end)
}

// GetMoodOptions 获取心情选项（用于前端显示）
func (s *MoodService) GetMoodOptions() map[string]interface{} {
	return map[string]interface{}{
		"time_contexts":    constant.TimeContexts,
		"mood_levels":      constant.MoodLevelDescriptions,
		"influences":       constant.Influences,
		"common_mood_tags": constant.CommonMoodTags,
	}
}
