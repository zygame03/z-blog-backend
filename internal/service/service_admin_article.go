package service

import (
	"my_web/backend/internal/models"

	"gorm.io/gorm"
)

type AdminArticleService struct {
	DB *gorm.DB
}

func NewAdminArticleService(db *gorm.DB) *AdminArticleService {
	return &AdminArticleService{
		DB: db,
	}
}

func (s *AdminArticleService) AdminGetArticlesByPage(page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	result := s.DB.Model(models.Article{}).
		Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = s.DB.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, total, nil
}

// 保存或更新文章
func (s *AdminArticleService) AdminSaveOrUpdateArticle(article *models.Article) error {
	if article.ID == 0 {
		err := article.Create(s.DB)
		if err != nil {
			return err
		}
		return nil
	}

	err := article.Update(s.DB)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminArticleService) AdminDeleteArticle(id int) error {
	result := s.DB.Model(&models.Article{}).
		Where("id = ? AND is_delete = false", id).
		UpdateColumn("is_delete", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// 修改文章状态
func (s *AdminArticleService) AdminChangeStatus(id, status int) error {
	result := s.DB.Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
