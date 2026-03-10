package article

import (
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

func (r *repo) getAllArticleIDs() ([]int, error) {
	ids := []int{}

	result := r.db.
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
func (r *repo) getArticlesByPage(page, pageSize int) ([]ArticleWithoutContent, int, error) {
	var articles []ArticleWithoutContent
	var total int64

	result := r.db.
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&total)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("db count articles failed: %w", result.Error)
	}

	result = r.db.
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
func (r *repo) getArticleByID(id int) (*Article, error) {
	var article Article

	result := r.db.First(&article, id)
	if result.Error != nil {
		return &article, fmt.Errorf("db get article by id failed: %w", result.Error)
	}

	return &article, nil
}

func (r *repo) mGetArticleByID(ids []int) ([]*Article, error) {
	var articles []*Article

	err := r.db.
		Where("id IN ?", ids).
		Find(&articles).
		Error
	if err != nil {
		return nil, fmt.Errorf("db mget articles by ids failed: %w", err)
	}

	return articles, nil
}

// getArticlesByPopular
func (r *repo) getArticlesByPopular(limit int) ([]ArticleWithoutContent, error) {
	var articles []ArticleWithoutContent

	result := r.db.
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
func (r *repo) incrementViews(id int, increment int64) error {
	result := r.db.
		Model(&Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", increment))
	if result.Error != nil {
		return fmt.Errorf("db increment article views failed: %w", result.Error)
	}

	return nil
}

// batchUpdateViews 批量更新文章的 views
func (r *repo) batchUpdateViews(viewsMap map[int]int64) error {
	if len(viewsMap) == 0 {
		return nil
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for id, increment := range viewsMap {
			if err := tx.Model(&Article{}).
				Where("id = ?", id).
				UpdateColumn("views", gorm.Expr("views + ?", increment)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("db batch update article views failed: %w", err)
	}

	return nil
}
