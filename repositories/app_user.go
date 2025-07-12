package repositories

import (
	"errors"

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
