package article

import (
	"context"
	"fmt"
	"my_web/backend/internal/zerrors"

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

	err := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Select("id").
		Where("is_delete = false AND status = ?", ArticlePublic).
		Pluck("id", &ids).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, zerrors.ArticleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return ids, nil
}

// listByPage
func (r *repo) listByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error) {
	var total int64

	// 获取总数
	err := r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	var articles []ArticleSummary
	err = r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return articles, int(total), nil
}

// getByID
func (r *repo) getByID(ctx context.Context, id int) (*Article, error) {
	var article Article

	err := r.db.
		WithContext(ctx).
		Where("id = ? AND is_delete = false AND status = ?", id, ArticlePublic).
		First(&article).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repo get article by id failed: %w", err)
	}

	return &article, nil
}

// listPopular
func (r *repo) listPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	var articles []ArticleSummary

	err := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Order("views DESC").
		Limit(limit).
		Find(&articles).
		Error
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return articles, nil
}

// incrementViews 增加文章的 views
func (r *repo) incrementViews(ctx context.Context, id int, increment int64) error {
	err := r.db.
		WithContext(ctx).
		Model(&Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", increment)).
		Error
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return nil
}

func (r *repo) save(ctx context.Context, article *Article) (int, error) {
	err := r.db.WithContext(ctx).Save(article).Error
	if err != nil {
		return 0, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
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
		return fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}
	return nil
}
