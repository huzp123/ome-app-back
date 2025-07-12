package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
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
func (d *UserGoalDAO) Create(goal *models.UserGoal) error {
	return d.db.Create(goal).Error
}

// GetByUserID 获取用户最新的目标设置
func (d *UserGoalDAO) GetByUserID(userID int64) (*models.UserGoal, error) {
	var goal models.UserGoal
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
func (d *UserGoalDAO) Update(goal *models.UserGoal) error {
	return d.db.Save(goal).Error
}

// CreateOrUpdate 创建或更新用户目标
func (d *UserGoalDAO) CreateOrUpdate(userID int64, goalType string, targetWeightKG float64,
	weeklyChangeKG float64, targetDate time.Time, dietType string,
	tastePreferences []string, foodIntolerances []string) error {

	// 先查询是否已存在该用户的目标
	var existingGoal models.UserGoal
	result := d.db.Where("user_id = ?", userID).First(&existingGoal)

	// 准备更新或创建的数据
	goal := models.UserGoal{
		UserID:           userID,
		GoalType:         goalType,
		TargetWeightKG:   targetWeightKG,
		WeeklyChangeKG:   weeklyChangeKG,
		TargetDate:       targetDate,
		DietType:         dietType,
		TastePreferences: tastePreferences,
		FoodIntolerances: foodIntolerances,
	}

	// 如果记录已存在，执行更新
	if result.Error == nil {
		goal.ID = existingGoal.ID // 保留原ID
		return d.db.Save(&goal).Error
	}

	// 如果记录不存在且错误是"记录未找到"，则创建新记录
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return d.db.Create(&goal).Error
	}

	// 其他错误直接返回
	return result.Error
}
