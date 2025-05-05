package model

import (
	"gorm.io/gorm"
)

// InitModels 初始化所有模型关系
func InitModels(db *gorm.DB) error {
	// 在这里可以添加数据库初始化逻辑
	// 比如设置外键关系、创建索引等

	return nil
}
