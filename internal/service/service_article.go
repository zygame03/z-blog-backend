package service

import (
	"context"
	"my_web/backend/internal/cache"
	"my_web/backend/internal/models"
	"my_web/backend/internal/repo"
)

type ArticleService struct {
	DB  *repo.ArticleRepo
	RDB *cache.ArticleCache
}

func NewArticleService(db *repo.ArticleRepo, rdb *cache.ArticleCache) *ArticleService {
	return &ArticleService{
		DB:  db,
		RDB: rdb,
	}
}

// 分页查找
func (s *ArticleService) GetArticlesByPage(ctx context.Context, page, pageSize int) ([]repo.ArticleVO, int, error) {
	// 1. 先查缓存
	articles, total, err := s.RDB.GetArticlesByPage(ctx, page, pageSize)
	if err == nil {
		return articles, total, nil
	}

	// 2. 缓存未命中或出错，查数据库
	if err == cache.ErrCacheMiss {
		articles, total, err := s.DB.GetArticlesByPage(page, pageSize)
		if err != nil {
			return nil, 0, err
		}

		go s.RDB.SetArticlesByPage(ctx, page, pageSize, articles, total)

		return articles, total, nil
	}

	return s.DB.GetArticlesByPage(page, pageSize)
}

// 获取热门文章，目前只基于view数，后续增加其他项综合判断
func (s *ArticleService) GetArticlesByPopular(limit int) ([]repo.ArticleVO, error) {
	ctx := context.Background()

	// 1. 先查缓存
	articles, err := s.RDB.GetArticlesByPopular(ctx, limit)
	if err == nil {
		return articles, nil
	}

	// 2. 缓存未命中，查数据库
	if err == cache.ErrCacheMiss {
		articles, err := s.DB.GetArticlesByPopular(limit)
		if err != nil {
			return nil, err
		}

		// 3. 异步回写缓存
		go s.RDB.SetArticlesByPopular(ctx, limit, articles)

		return articles, nil
	}

	// 4. 缓存出错，降级到数据库
	return s.DB.GetArticlesByPopular(limit)
}

// 通过ID获取文章，获取后增加views
func (s *ArticleService) GetArticleByID(id int) (*models.Article, error) {
	ctx := context.Background()

	// 1. 先查缓存
	article, err := s.RDB.GetArticleByID(ctx, id)
	if err == nil {
		// 缓存命中，异步增加 views
		go func() {
			s.DB.SetViewsIncr(id)
			s.RDB.SetViewsIncr(ctx, id)
		}()
		return article, nil
	}

	// 2. 缓存未命中，查数据库
	if err == cache.ErrCacheMiss {
		article, err := s.DB.GetArticleByID(id)
		if err != nil {
			return nil, err
		}

		// 3. 增加 views
		if err := s.DB.SetViewsIncr(id); err != nil {
			// views 增加失败不影响返回结果，只记录错误
		}
		s.RDB.SetViewsIncr(ctx, id)

		// 4. 回写缓存
		go s.RDB.SetArticleByID(ctx, id, article)

		return article, nil
	}

	// 5. 缓存出错，降级到数据库
	article, err = s.DB.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	// 即使缓存出错，也尝试增加 views
	go func() {
		s.DB.SetViewsIncr(id)
		s.RDB.SetViewsIncr(ctx, id)
	}()

	return article, nil
}

// Todo 按tags找文章
func (s *ArticleService) GetArticlesByTag(limit int) ([]models.Article, error) {
	return nil, nil
}
