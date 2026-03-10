package stats

import (
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

func (r *repo) updateViews(num int64) error {
	result := r.db.
		Model(&NumStats{}).
		Where("key = ?", "view").
		UpdateColumn("value", gorm.Expr("value + ?", num))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repo) getViews() (int64, error) {
	var num int64
	result := r.db.
		Model(&NumStats{}).
		Where("key = ?", "view").
		Select("value").
		Find(&num)
	if result.Error != nil {
		return -1, result.Error
	}
	return num, nil
}
