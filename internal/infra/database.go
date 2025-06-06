package infra

import (
	"fmt"
	"log"
	"my_web/backend/internal/article"
	"my_web/backend/internal/stats"
	"my_web/backend/internal/user"
	"my_web/backend/internal/websiteData"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Port     uint   `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&article.Article{},
		&user.Profile{},
		&websiteData.WebsiteData{},
		&stats.NumStats{},
	); err != nil {
		return nil, fmt.Errorf("数据库自动迁移失败: %w", err)
	}

	data := stats.NumStats{
		Key:   "view",
		Value: 0,
	}
	db.Model(&stats.NumStats{}).Create(&data)

	log.Println("数据库初始化成功")
	return db, nil
}
