package stats

import (
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

func newDatabase(db *gorm.DB) *database {
	return &database{
		db: db,
	}
}

func (db *database) updateViews(num int64) error {
	result := db.db.
		Model(&NumStats{}).
		Where("key = ?", "view").
		UpdateColumn("value", gorm.Expr("value + ?", num))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *database) getViews() (int64, error) {
	var num int64
	result := db.db.
		Model(&NumStats{}).
		Where("key = ?", "view").
		Select("value").
		Find(&num)
	if result.Error != nil {
		return -1, result.Error
	}
	return num, nil
}
