package stats

import (
	"context"

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

func (r *repo) updateViews(ctx context.Context, num int64) error {
	result := r.db.
		WithContext(ctx).
		Model(&NumStats{}).
		Where("key = ?", "view").
		UpdateColumn("value", gorm.Expr("value + ?", num))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repo) getViews(ctx context.Context) (int, error) {
	var num int
	result := r.db.
		WithContext(ctx).
		Model(&NumStats{}).
		Where("key = ?", "view").
		Select("value").
		Pluck("value", &num)
	if result.Error != nil {
		return -1, result.Error
	}

	return num, nil
}
