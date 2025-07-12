package models

import (
	"database/sql"
	"time"
)

// AppUser 用户基础信息及登录凭据
type AppUser struct {
	ID           int64          `json:"id"            gorm:"primaryKey"`               // BIGSERIAL
	UserName     string         `json:"user_name"     gorm:"column:user_name;size:32"` // 新增：昵称 / 用户名
	Phone        sql.NullString `json:"phone"         gorm:"size:20;uniqueIndex:idx_phone,where:phone IS NOT NULL"`
	Email        sql.NullString `json:"email"         gorm:"size:50;uniqueIndex:idx_email,where:email IS NOT NULL"`
	PasswordHash string         `json:"password_hash" gorm:"size:128"` // 微信登录时可为空

	// 微信登录相关字段
	WechatOpenID sql.NullString `json:"wechat_openid" gorm:"column:wechat_openid;size:64;uniqueIndex:idx_wechat_openid,where:wechat_openid IS NOT NULL"` // 微信OpenID
	AvatarURL    sql.NullString `json:"avatar_url"    gorm:"column:avatar_url;size:255"`                                                                 // 头像URL

	HeightCM  sql.NullFloat64 `json:"height_cm"  gorm:"type:decimal(5,2);check:height_cm IS NULL OR (height_cm BETWEEN 50 AND 300)"`
	BirthDate time.Time       `json:"birth_date" gorm:"type:date;default:null"`
	Sex       string          `json:"sex"        gorm:"size:6"` // male / female / other

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AppUser) TableName() string {
	return "app_users"
}
