package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ome-app-back/config"
	"ome-app-back/models"
)

// Init 初始化数据库连接和迁移
func Init(dbConfig config.DBConfig) (*gorm.DB, error) {
	// 连接数据库
	db, err := connectDB(dbConfig)
	if err != nil {
		log.Printf("数据库连接失败: %v", err)
		return nil, err
	}

	// 自动迁移数据库表结构
	if err := models.Init(db); err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return nil, err
	}

	return db, nil
}

// connectDB 连接数据库
func connectDB(dbConfig config.DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch dbConfig.Type {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dbConfig.GetDSN()), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dbConfig.GetDSN()), &gorm.Config{})
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", dbConfig.Type)
	}

	if err != nil {
		return nil, err
	}

	// 获取底层SQL DB以设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	return db, nil
}
