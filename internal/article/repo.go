package article

import (
	"context"
	"fmt"

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

func (r *repo) getAllArticleIDs(ctx context.Context) ([]int, error) {
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

// getArticlesByPage
func (r *repo) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
	var articles []ArticleWithoutContent
	var total int64

	result := r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&total)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("db count articles failed: %w", result.Error)
	}

	result = r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("db find articles by page failed: %w", result.Error)
	}

	return articles, int(total), nil
}

// getArticleByID
func (r *repo) getArticleByID(ctx context.Context, id int) (*Article, error) {
	var article Article

	result := r.db.WithContext(ctx).First(&article, id)
	if result.Error != nil {
		return &article, fmt.Errorf("db get article by id failed: %w", result.Error)
	}

	return &article, nil
}

func (r *repo) mGetArticleByID(ctx context.Context, ids []int) ([]*Article, error) {
	var articles []*Article

	err := r.db.
		WithContext(ctx).
		Where("id IN ?", ids).
		Find(&articles).
		Error
	if err != nil {
		return nil, fmt.Errorf("db mget articles by ids failed: %w", err)
	}

	return articles, nil
}

// getArticlesByPopular
func (r *repo) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
	var articles []ArticleWithoutContent

	result := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Order("views DESC").
		Select("id, created_at, updated_at, title, author_name, views, tags, cover").
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
