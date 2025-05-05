package dao

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/internal/model"
)

// UserGoalDAO 处理用户目标数据访问
type UserGoalDAO struct {
	db *gorm.DB
}

// NewUserGoalDAO 创建用户目标DAO实例
func NewUserGoalDAO(db *gorm.DB) *UserGoalDAO {
	return &UserGoalDAO{db: db}
}

// Create 创建用户目标
func (d *UserGoalDAO) Create(goal *model.UserGoal) error {
	return d.db.Create(goal).Error
}

// GetByUserID 获取用户最新的目标设置
func (d *UserGoalDAO) GetByUserID(userID int64) (*model.UserGoal, error) {
	var goal model.UserGoal
	if err := d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&goal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未找到用户目标")
		}
		return nil, err
	}
	return &goal, nil
}

// Update 更新用户目标
func (d *UserGoalDAO) Update(goal *model.UserGoal) error {
	return d.db.Save(goal).Error
}

// CreateOrUpdate 创建或更新用户目标
func (d *UserGoalDAO) CreateOrUpdate(userID int64, goalType string, targetWeightKG float64,
	weeklyChangeKG float64, targetDate time.Time, dietType string,
	tastePreferences []string, foodIntolerances []string) error {

	goal := model.UserGoal{
		UserID:           userID,
		GoalType:         goalType,
		TargetWeightKG:   targetWeightKG,
		WeeklyChangeKG:   weeklyChangeKG,
		TargetDate:       targetDate,
		DietType:         dietType,
		TastePreferences: tastePreferences,
		FoodIntolerances: foodIntolerances,
	}

	return d.db.Create(&goal).Error
}
