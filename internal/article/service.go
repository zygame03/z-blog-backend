package article

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db  *database
	rdb *cache

	getCfg func() *ArticleConfig
}

func NewService(db *gorm.DB, rdb *redis.Client, getCfg func() *ArticleConfig) *Service {
	service := &Service{
		getCfg: getCfg,
	}

	service.db = newDatabase(db)
	service.rdb = newCache(rdb, service.getCfg)

	return service
}

func (s *Service) RegisterCron(cron *cron.Cron) {
	_, err := cron.AddFunc("@every 24h", s.syncArticleViews)
	if err != nil {
		return
	}
	logger.Info(
		"add func successfully",
		zap.Int("interval", int(s.getCfg().SyncInterval)),
		zap.String("func name", "syncArticleViews"),
	)
}

func (s *Service) syncArticleViews() {
	ctx := context.Background()
	ids, err := s.db.getAllArticleIDs()
	if err != nil {
		logger.Error(
			"load article ids for view sync failed",
			zap.Error(err),
		)
		return
	}

	for _, id := range ids {
		num, err := s.rdb.GetViewUV(ctx, id)
		if err != nil {
			logger.Error(
				"get view uv from cache failed",
				zap.Int("id", id),
				zap.Error(err),
			)
			continue
		}
		if num == 0 {
			continue
		}

		err = s.rdb.DelViewUV(ctx, id)
		if err != nil {
			logger.Error(
				"delete view uv from cache failed",
				zap.Int("id", id),
				zap.Error(err),
			)
			continue
		}

		if err := s.db.incrementViews(id, num); err != nil {
			logger.Error(
				"increment article views failed",
				zap.Int("id", id),
				zap.Int64("increment", num),
				zap.Error(err),
			)
		}
	}
}

func (s *Service) GetArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
	articles, total, err := s.rdb.GetArticlesByPage(ctx, page, pageSize)
	if err == nil {
		logger.Info(
			"get articles by page from cache",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
		return articles, total, nil
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for articles by page",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)

		articles, total, err = s.db.getArticlesByPage(page, pageSize)
		if err != nil {
			return nil, 0, err
		}

		s.rdb.SetArticlesByPage(ctx, page, pageSize, articles, total)
		return articles, total, nil
	}

	logger.Error(
		"get articles by page from cache failed",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Error(err),
	)

	return s.db.getArticlesByPage(page, pageSize)
}

func (s *Service) GetArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
	articles, err := s.rdb.GetArticlesByPopular(ctx, limit)
	if err == nil {
		logger.Info(
			"get articles by popular from cache",
			zap.Int("limit", limit),
		)
		return articles, nil
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for articles by popular",
			zap.Int("limit", limit),
		)

		articles, err = s.db.getArticlesByPopular(limit)
		if err != nil {
			return nil, err
		}

		go s.rdb.SetArticlesByPopular(ctx, limit, articles)
		return articles, nil
	}

	logger.Error(
		"get articles by popular from cache failed",
		zap.Int("limit", limit),
		zap.Error(err),
	)

	return s.db.getArticlesByPopular(limit)
}

func (s *Service) GetArticleByID(ctx context.Context, id int, userID string) (*Article, error) {
	article, err := s.rdb.GetArticleByID(ctx, id)
	if err == nil {
		logger.Info(
			"get article by id from cache",
			zap.Int("id", id),
			zap.String("user_id", userID),
		)
		s.rdb.AddViewUV(ctx, id, userID)
		return article, nil
	}

	if err != ErrCacheMiss {
		logger.Error(
			"get article by id from cache failed",
			zap.Int("id", id),
			zap.String("user_id", userID),
			zap.Error(err),
		)
	} else {
		logger.Info(
			"cache miss for article by id",
			zap.Int("id", id),
			zap.String("user_id", userID),
		)
	}

	article, err = s.db.getArticleByID(id)
	if err != nil {
		return nil, err
	}
	s.rdb.AddViewUV(ctx, id, userID)
	s.rdb.SetArticleByID(ctx, id, article)
	return article, nil
}

func (s *Service) GetArticlesByTag(limit int) ([]Article, error) {
	return nil, nil
}
