package user

import (
	"my_web/backend/internal/logger"

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

func (r *repo) getProfile(db *gorm.DB, id int) (*Profile, error) {
	var profile Profile

	result := db.
		Model(&Profile{}).
		Where("id = ?", id).
		First(&profile)
	if result.Error != nil {
		logger.Error(
			"db get profile failed",
			zap.Int("id", id),
		)
		return nil, result.Error
	}

	return &profile, nil
}
