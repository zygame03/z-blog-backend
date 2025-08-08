package site

import (
	"context"
	"my_web/backend/internal/logger"
	"time"

	"go.uber.org/zap"
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
func (r *repo) getIntro() (string, error) {
	var data WebsiteData

	result := r.db.
		Where("key = ?", "intro").
		First(&data)
	if result.Error != nil {
		return "", result.Error
	}

	return data.Value, nil
}

func (r *repo) getAnnouncement(ctx context.Context) ([]*announcementBO, error) {
	var data []*announcementBO
	now := time.Now()

	result := r.db.
		WithContext(ctx).
		Model(&Announcement{}).
		Where("online_at <= ? AND offline_at >= ?", now, now).
		Select("id", "text").
		Find(&data)
	if result.Error != nil {
		logger.Error(
			"repo get announcement failed",
			zap.Error(result.Error),
		)
		return nil, result.Error
	}

	return data, nil
}
