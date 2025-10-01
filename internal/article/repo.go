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
func (r *repo) listByPage(ctx context.Context, page, pageSize int) (*ArticlePageVO, error) {
	var data = ArticlePageVO{
		Page:     page,
		PageSize: pageSize,
	}
	err := r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&data.Total).
		Error
	if err != nil {
		return nil, fmt.Errorf("repo get total failed: %w", err)
	}

	err = r.db.
		WithContext(ctx).
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&data.List).
		Error
	if err != nil {
		return nil, fmt.Errorf("repo list articles by page failed: %w", err)
	}

	return &data, nil
}

// getArticleByID
func (r *repo) getArticleByID(ctx context.Context, id int) (*Article, error) {
	var article Article

	err := r.db.
		WithContext(ctx).
		Where("id = ? AND is_delete = false AND status = ?", id, ArticlePublic).
		First(&article).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, zerrors.ArticleNotFound
	}
	if err != nil {
		return &article, fmt.Errorf("repo get article by id failed: %w", err)
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

type CommentQuery struct {
	status   *CommentStatus // nil = 不过滤
	id       int
	page     int
	pageSize int
}

func (r *repo) listComment(ctx context.Context, query CommentQuery) (*CommentPageVO, error) {
	var comment = CommentPageVO{
		Page:     query.page,
		PageSize: query.pageSize,
	}
	db := r.db.WithContext(ctx)
	if query.status != nil {
		db = db.Where("status = ? AND article_id = ?", *query.status, query.id)
	}

	err := db.Count(&comment.Total).Error
	if err != nil {
		return nil, fmt.Errorf("repo get comments total faild: %w", zerrors.ErrDBOperation)
	}

	err = db.
		Offset((query.page - 1) * query.pageSize).
		Limit(query.pageSize).
		Find(&comment.Comments).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, zerrors.CommentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("repo list comments failed: %w", err)
	}

	return &comment, nil
}

func (r *repo) saveComment(ctx context.Context, comment *ArticleComment) (int, error) {
	err := r.db.WithContext(ctx).Save(comment).Error
	if err != nil {
		return 0, fmt.Errorf("%w: %v", zerrors.ErrDBOperation, err)
	}

	return comment.ID, nil
}

func (r *repo) updateCommentStatus(ctx context.Context, id int, status CommentStatus) error {
	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		UpdateColumn("status", status).
		Error
	if err != nil {
		return err
	}
	return nil
}
