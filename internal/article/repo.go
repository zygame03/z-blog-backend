package article

import (
	"context"
	"fmt"
	"my_web/backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func newRepo(db *gorm.DB) *repo {
	return &repo{
		db,
	}
}

func (r *repo) listIDs(ctx context.Context) ([]int, error) {
	ids := []int{}

	result := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Select("id").
		Where("is_delete = false AND status = ?", ArticlePublic).
		Pluck("id", &ids)

	if result.Error != nil {
		return nil, fmt.Errorf("db get all article ids failed: %w", result.Error)
	}
	return ids, nil
}

// listByPage
func (r *repo) listByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error) {
	var articles []ArticleSummary
	var total int64

	result := r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, int(total), nil
}

// getByID
func (r *repo) getByID(ctx context.Context, id int) (*Article, error) {
	var article Article

	result := r.db.
		WithContext(ctx).
		Where("id = ? AND is_delete = false AND status = ?", id, ArticlePublic).
		First(&article)
	if result.Error != nil {
		return &article, fmt.Errorf("db get article by id failed: %w", result.Error)
	}

	return &article, nil
}

// listPopular
func (r *repo) listPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	var articles []ArticleSummary

	result := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Order("views DESC").
		Limit(limit).
		Find(&articles)
	if result.Error != nil {
		return nil, fmt.Errorf("db get articles by popular failed: %w", result.Error)
	}

	return articles, nil
}

// incrementViews 增加文章的 views
func (r *repo) incrementViews(ctx context.Context, id int, increment int64) error {
	result := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", increment))
	if result.Error != nil {
		return fmt.Errorf("db increment article views failed: %w", result.Error)
	}

	return nil
}

func (r *repo) save(ctx context.Context, article *Article) (int, error) {
	err := r.db.WithContext(ctx).Save(article).Error
	if err != nil {
		logger.Error(
			"save article failed",
			zap.Int("id", article.ID),
			zap.Error(err),
		)
		return 0, err
	}

	return article.ID, nil
}

func (r *repo) delete(ctx context.Context, id int) error {
	err := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("id = ?", id).
		UpdateColumn("is_delete", false).
		Error
	if err != nil {
		return err
	}
	return nil
}
