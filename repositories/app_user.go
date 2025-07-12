package repositories

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"ome-app-back/models"
)

// AppUserDAO 处理用户数据访问
type AppUserDAO struct {
	db *gorm.DB
}

// NewAppUserDAO 创建用户DAO实例
func NewAppUserDAO(db *gorm.DB) *AppUserDAO {
	return &AppUserDAO{db: db}
}

// Create 创建新用户
func (d *AppUserDAO) Create(user *models.AppUser) error {
	return d.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (d *AppUserDAO) GetByID(id int64) (*models.AppUser, error) {
	var user models.AppUser
	if err := d.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取用户
func (d *AppUserDAO) GetByPhone(phone string) (*models.AppUser, error) {
	var user models.AppUser
	if err := d.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (d *AppUserDAO) GetByEmail(email string) (*models.AppUser, error) {
	var user models.AppUser
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByWechatOpenID 根据微信OpenID获取用户
func (d *AppUserDAO) GetByWechatOpenID(openID string) (*models.AppUser, error) {
	var user models.AppUser
	if err := d.db.Where("wechat_openid = ?", openID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (d *AppUserDAO) Update(user *models.AppUser) error {
	return d.db.Save(user).Error
}

// UpdateBasicInfo 更新用户基本信息（已弃用，请使用Update方法）
func (d *AppUserDAO) UpdateBasicInfo(id int64, heightCM float64, birthDate string, sex string) error {
	// 创建一个更新映射
	updates := map[string]interface{}{}

	if heightCM > 0 {
		updates["height_cm"] = heightCM
	}

	if birthDate != "" {
		date, err := time.Parse("2006-01-02", birthDate)
		if err == nil {
			updates["birth_date"] = date
		}
	}

	if sex != "" {
		updates["sex"] = sex
	}

	return d.db.Model(&models.AppUser{}).
		Where("id = ?", id).
		Updates(updates).Error
}
