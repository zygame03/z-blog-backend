package repo

import (
	"my_web/backend/internal/models"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	DB *gorm.DB
}

type ArticleVO struct {
	models.Model
	Title      string `json:"title"`
	AuthorName string `json:"authorName"` // 作者
	Views      uint   `json:"views"`      // 浏览数
	Tags       string `json:"tags"`       // 标签（逗号分隔形式）
	Cover      string `json:"cover"`      // 封面
}

func (r *ArticleRepo) GetArticlesByPage(page, pageSize int) ([]ArticleVO, int, error) {
	var articles []ArticleVO
	var total int64

	result := r.DB.Model(models.Article{}).
		Where("is_delete = false AND status = 0").
		Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = r.DB.Model(models.Article{}).
		Where("is_delete = false AND status = 0").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, int(total), nil
}

func (r *ArticleRepo) GetArticleByID(id int) (*models.Article, error) {
	var article models.Article

	result := r.DB.First(&article, id)
	if result.Error != nil {
		return &article, result.Error
	}

	return &article, nil
}

func (r *ArticleRepo) GetArticlesByPopular(limit int) ([]ArticleVO, error) {
	var articles []ArticleVO

	result := r.DB.Model(&models.Article{}).
		Where("is_delete = false AND status = 0").
		Order("views DESC").
		Select("id, created_at, updated_at, title, author_name, views, tags, cover").
		Limit(limit).
		Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

func (r *ArticleRepo) SetViewsIncr(id int) error {
	result := r.DB.Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1))

	return result.Error
}
