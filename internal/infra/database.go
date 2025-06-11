package infra

import (
	"fmt"
	"my_web/backend/internal/article"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/site"
	"my_web/backend/internal/stats"
	"my_web/backend/internal/user"

	"go.uber.org/zap"
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

	AutoMigrate bool `mapstructure:"auto_migrate"`
}

// InitDatabase 初始化数据库连接
func InitDatabase(conf *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if conf.AutoMigrate {
		logger.Info(
			"start database auto migrate",
		)
		// 自动迁移
		err := db.AutoMigrate(
			&article.Article{},
			&user.Profile{},
			&site.WebsiteData{},
			&stats.NumStats{},
		)
		if err != nil {
			logger.Info(
				"database auto migrate failed",
				zap.Error(err),
			)
			return nil, err
		}
	}

	return db, nil
}
