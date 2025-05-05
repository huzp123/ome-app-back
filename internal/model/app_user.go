package model

import (
	"time"
)

// AppUser 用户基础信息及登录凭据
type AppUser struct {
	ID           int64  `json:"id"            gorm:"primaryKey"`               // BIGSERIAL
	UserName     string `json:"user_name"     gorm:"column:user_name;size:32"` // 新增：昵称 / 用户名
	Phone        string `json:"phone"         gorm:"size:20;unique"`
	Email        string `json:"email"         gorm:"size:50;unique"`
	PasswordHash string `json:"password_hash" gorm:"size:128;not null"`

	HeightCM  float64   `json:"height_cm"  gorm:"type:decimal(5,2);check:height_cm BETWEEN 50 AND 300"`
	BirthDate time.Time `json:"birth_date" gorm:"type:date;default:null"`
	Sex       string    `json:"sex"        gorm:"size:6"` // male / female / other

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AppUser) TableName() string {
	return "app_users"
}
