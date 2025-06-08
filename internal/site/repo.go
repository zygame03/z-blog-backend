package site

import (
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func newRepo(db *gorm.DB) *repo {
	return &repo{
		db: db,
	}
}

// get intro from repository
func (db *repo) getIntro() (string, error) {
	var data WebsiteData

	result := db.db.
		Where("key = ?", "intro").
		First(&data)
	if result.Error != nil {
		return "", result.Error
	}

	return data.Value, nil
}

func (db *repo) getAnnouncement() ([]string, error) {
	var data []string
	now := time.Now()

	result := db.db.
		Model(&Announcement{}).
		Select("text").
		Where("online_at <= ? AND offline_at >= ?", now, now).
		Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return data, nil
}
