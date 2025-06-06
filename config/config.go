package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	ProxyURL    string  `yaml:"proxy_url"`
	TestMode    bool    `yaml:"test_mode"`
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
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
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

// CheckConfiguration 检查配置的有效性
func (c *Config) CheckConfiguration() []string {
	issues := []string{}

	// 检查AI配置
	if c.AI.APIKey == "" {
		msg := "警告: AI API密钥未设置，AI功能将无法使用"
		log.Println(msg)
		issues = append(issues, msg)
	}

	if c.AI.APIURL == "" {
		msg := "警告: AI API URL未设置，使用默认URL"
		log.Println(msg)
		issues = append(issues, msg)
	}

	if c.AI.Model == "" {
		msg := "警告: AI模型未指定，使用默认模型"
		log.Println(msg)
		issues = append(issues, msg)
	}

	// 检查数据库配置
	if c.DB.User == "" || c.DB.Password == "" {
		msg := "警告: 数据库用户名或密码未设置"
		log.Println(msg)
		issues = append(issues, msg)
	}

	if c.DB.Host == "" {
		msg := "警告: 数据库主机未设置"
		log.Println(msg)
		issues = append(issues, msg)
	}

	// 检查上传目录
	if c.Upload.Dir == "" {
		msg := "警告: 文件上传基本路径未设置"
		log.Println(msg)
		issues = append(issues, msg)
	} else {
		// 检查上传目录是否存在且可写
		if _, err := os.Stat(c.Upload.Dir); os.IsNotExist(err) {
			msg := fmt.Sprintf("警告: 上传目录 %s 不存在", c.Upload.Dir)
			log.Println(msg)
			issues = append(issues, msg)
		} else {
			// 尝试创建测试文件检查写权限
			testPath := fmt.Sprintf("%s/test_write_permission", c.Upload.Dir)
			testFile, err := os.Create(testPath)
			if err != nil {
				msg := fmt.Sprintf("警告: 上传目录 %s 可能没有写权限: %v", c.Upload.Dir, err)
				log.Println(msg)
				issues = append(issues, msg)
			} else {
				testFile.Close()
				os.Remove(testPath) // 清理测试文件
			}
		}
	}

	if len(issues) == 0 {
		log.Println("配置检查完成，未发现问题")
	} else {
		log.Printf("配置检查完成，发现 %d 个问题", len(issues))
	}

	return issues
}
