package article

import (
	"fmt"
	"my_web/backend/internal/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func repoGetAllArticleIDs(db *gorm.DB) ([]int, error) {
	ids := []int{}

	result := db.
		Model(&Article{}).
		Select("id").
		Where("is_delete = false AND status = ?", ArticlePublic).
		Pluck("id", &ids)

	if result.Error != nil {
		logger.Error(
			"db get all article ids failed",
			zap.Error(result.Error),
		)
		return nil, fmt.Errorf("db get all article ids failed: %w", result.Error)
	}
	return ids, nil
}

// repoGetArticlesByPage
func repoGetArticlesByPage(
	db *gorm.DB,
	page, pageSize int,
) ([]ArticleWithoutContent, int, error) {

	var articles []ArticleWithoutContent
	var total int64

	result := db.
		Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Count(&total)
	if result.Error != nil {
		logger.Error(
			"db count articles failed",
			zap.Error(result.Error),
		)
		return nil, 0, fmt.Errorf("db count articles failed: %w", result.Error)
	}

	result = db.Model(Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		logger.Error(
			"db find articles by page failed",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(result.Error),
		)
		return nil, 0, fmt.Errorf("db find articles by page failed: %w", result.Error)
	}

	return articles, int(total), nil
}

// repoGetArticleByID
func repoGetArticleByID(db *gorm.DB, id int) (*Article, error) {
	var article Article

	result := db.First(&article, id)
	if result.Error != nil {
		logger.Error(
			"db get article by id failed",
			zap.Int("id", id),
			zap.Error(result.Error),
		)
		return &article, fmt.Errorf("db get article by id failed: %w", result.Error)
	}

	return &article, nil
}

func repoMGetArticleByID(db *gorm.DB, ids []int) ([]*Article, error) {
	var articles []*Article

	err := db.
		Where("id IN ?", ids).
		Find(&articles).
		Error
	if err != nil {
		logger.Error(
			"db mget articles by ids failed",
			zap.Ints("ids", ids),
			zap.Error(err),
		)
		return nil, fmt.Errorf("db mget articles by ids failed: %w", err)
	}

	return articles, nil
}

// repoGetArticlesByPopular
func repoGetArticlesByPopular(db *gorm.DB, limit int) ([]ArticleWithoutContent, error) {
	var articles []ArticleWithoutContent

	result := db.
		Model(&Article{}).
		Where("is_delete = false AND status = ?", ArticlePublic).
		Order("views DESC").
		Select("id, created_at, updated_at, title, author_name, views, tags, cover").
		Limit(limit).
		Find(&articles)
	if result.Error != nil {
		logger.Error(
			"db get articles by popular failed",
			zap.Int("limit", limit),
			zap.Error(result.Error),
		)
		return nil, fmt.Errorf("db get articles by popular failed: %w", result.Error)
	}

	return articles, nil
}

// repoIncrementViews 增加文章的 views
func repoIncrementViews(db *gorm.DB, id int, increment int64) error {
	result := db.
		Model(&Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", increment))
	if result.Error != nil {
		logger.Error(
			"db increment article views failed",
			zap.Int("id", id),
			zap.Int64("increment", increment),
			zap.Error(result.Error),
		)
		return fmt.Errorf("db increment article views failed: %w", result.Error)
	}

	return nil
}

// repoBatchUpdateViews 批量更新文章的 views
func repoBatchUpdateViews(db *gorm.DB, viewsMap map[int]int64) error {
	if len(viewsMap) == 0 {
		return nil
	}

	err := db.Transaction(func(tx *gorm.DB) error {
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
		logger.Error(
			"db batch update article views failed",
			zap.Any("views_map", viewsMap),
			zap.Error(err),
		)
		return fmt.Errorf("db batch update article views failed: %w", err)
	}

	return nil
}
