package websiteData

import (
	"my_web/backend/internal/logger"
	"time"

	"gorm.io/gorm"
)

type database struct {
	DB *gorm.DB
}

func NewDatabase(db *gorm.DB) *database {
	return &database{
		DB: db,
	}
}

// get intro from repository
func (db *database) getIntro() (string, error) {
	var data WebsiteData

	result := db.DB.
		Where("key = ?", "intro").
		First(&data)
	if result.Error != nil {
		logger.Error(
			"database get intro failed",
		)
		return "", result.Error
	}

	return data.Value, nil
}

func (db *database) getAnnouncement() ([]string, error) {
	var data []string
	now := time.Now()

	result := db.DB.
		Model(&Announcement{}).
		Select("text").
		Where("online_at <= ? AND offline_at >= ?", now, now).
		Find(&data)
	if result.Error != nil {
		logger.Error(
			"database get announcement failed",
		)
		return nil, result.Error
	}

	return data, nil
}
