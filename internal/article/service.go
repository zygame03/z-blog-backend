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
	db  *repo
	rdb *cache

	conf func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, conf func() *Config) *Service {
	service := &Service{
		conf: conf,
	}

	service.db = newRepo(db)
	service.rdb = newCache(rdb, service.conf)

	return service
}

func (s *Service) RegisterCron(cron *cron.Cron) {
	_, err := cron.AddFunc("@every 24h", s.syncArticleViews)
	if err != nil {
		return
	}
	logger.Info(
		"添加定时任务成功",
		zap.Int("间隔", int(s.conf().SyncInterval)),
		zap.String("任务描述", "文章浏览数同步"),
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
		num, err := s.rdb.getViewUV(ctx, id)
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

		err = s.rdb.delViewUV(ctx, id)
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
	articles, total, err := s.rdb.getArticlesByPage(ctx, page, pageSize)
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

		s.rdb.setArticlesByPage(ctx, page, pageSize, articles, total)
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
	articles, err := s.rdb.getArticlesByPopular(ctx, limit)
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

		go s.rdb.setArticlesByPopular(ctx, limit, articles)
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
	article, err := s.rdb.getArticleByID(ctx, id)
	if err == nil {
		logger.Info(
			"get article by id from cache",
			zap.Int("id", id),
			zap.String("user_id", userID),
		)
		s.rdb.addViewUV(ctx, id, userID)
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
	s.rdb.addViewUV(ctx, id, userID)
	s.rdb.setArticleByID(ctx, id, article)
	return article, nil
}

func (s *Service) GetArticlesByTag(limit int) ([]Article, error) {
	return nil, nil
}
