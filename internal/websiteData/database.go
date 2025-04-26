package websiteData

import (
	"my_web/backend/internal/logger"

	"gorm.io/gorm"
)

// get intro from repository
func repoGetIntro(db *gorm.DB) (string, error) {
	var data WebsiteData

	result := db.
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
