package config

import (
	"log"
)

// Init 初始化配置
func Init(configPath string) (*Config, error) {
	// 加载配置
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Printf("加载配置失败: %v", err)
		return nil, err
	}

	// 检查配置有效性
	issues := cfg.CheckConfiguration()
	if len(issues) > 0 {
		log.Println("配置存在问题，请检查日志")
	}

	return cfg, nil
}
