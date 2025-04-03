package data

import (
	"my_web/backend/internal/logger"

	"gorm.io/gorm"
)

// get intro from repository
func repoGetIntro(db *gorm.DB) (string, error) {
	var intro string

	result := db.
		Model(Sitedata{}).
		Where("id = ?", "intro").
		First(&intro)
	if result.Error != nil {
		logger.Error(
			"database get intro failed",
		)
		return "", result.Error
	}

	return intro, nil
}
