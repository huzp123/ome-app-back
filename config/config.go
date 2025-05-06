package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config 应用配置结构
type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"database"`
	AI     AIConfig     `yaml:"ai"`
	Upload UploadConfig `yaml:"upload"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Type            string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"username"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// AIConfig AI服务配置
type AIConfig struct {
	APIKey      string  `yaml:"api_key"`
	APIURL      string  `yaml:"api_url"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Dir     string `yaml:"dir"`
	MaxSize int64  `yaml:"max_size"`
}

// GetDSN 获取数据库连接字符串
func (db *DBConfig) GetDSN() string {
	switch db.Type {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			db.Host, db.User, db.Password, db.DBName, db.Port)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&sql_mode='ALLOW_INVALID_DATES'",
			db.User, db.Password, db.Host, db.Port, db.DBName)
	default:
		log.Fatalf("不支持的数据库类型: %s", db.Type)
		return ""
	}
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
