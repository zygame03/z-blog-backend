package article

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleRepo interface {
	getAllArticleIDs(ctx context.Context) ([]int, error)
	getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error)
	getArticleByID(ctx context.Context, id int) (*Article, error)
	getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error)
	incrementViews(ctx context.Context, id int, inc int64) error
}

type ArticleCache interface {
	getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error)
	setArticlesByPage(ctx context.Context, page, pageSize int, articles []ArticleWithoutContent, total int) error
	getArticleByID(ctx context.Context, id int) (*Article, error)
	setArticleByID(ctx context.Context, id int, article *Article) error
	getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error)
	setArticlesByPopular(ctx context.Context, limit int, articles []ArticleWithoutContent) error
	addViewUV(ctx context.Context, id int, userID string) error
	getViewUV(ctx context.Context, id int) (int64, error)
	delViewUV(ctx context.Context, id int) error
}

type ArticleService struct {
	db  ArticleRepo
	rdb ArticleCache

	cfg func() *Config
}

func NewService(db *gorm.DB, rdb *redis.Client, cfg func() *Config) *ArticleService {
	service := &ArticleService{
		cfg: cfg,
	}

	service.db = newRepo(db)
	service.rdb = newCache(rdb, cfg)

	return service
}

func (s *ArticleService) RegisterCron(cron *cron.Cron) {
	_, err := cron.AddFunc("@every "+s.cfg().SyncInterval.String(), s.syncArticleViews)
	if err != nil {
		return
	}

	logger.Info(
		"add func successfully",
		zap.Int("interval", int(s.cfg().SyncInterval)),
		zap.String("func", "syncArticleViews"),
	)
}

func (s *ArticleService) syncArticleViews() {
	ctx := context.Background()
	ids, err := s.db.getAllArticleIDs(ctx)
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

		if err := s.db.incrementViews(ctx, id, num); err != nil {
			logger.Error(
				"increment article views failed",
				zap.Int("id", id),
				zap.Int64("increment", num),
				zap.Error(err),
			)
		}
	}
}

func (s *ArticleService) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
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
			zap.Error(err),
		)

		articles, total, err = s.db.getArticlesByPage(ctx, page, pageSize)
		if err != nil {
			logger.Error(
				"repo get articles by page failed",
				zap.Int("page", page),
				zap.Int("page_size", pageSize),
				zap.Error(err),
			)
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

	return s.db.getArticlesByPage(ctx, page, pageSize)
}

func (s *ArticleService) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
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

		articles, err = s.db.getArticlesByPopular(ctx, limit)
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

	return s.db.getArticlesByPopular(ctx, limit)
}

func (s *ArticleService) getArticleByID(ctx context.Context, id int, userID string) (*Article, error) {
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

	article, err = s.db.getArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.rdb.addViewUV(ctx, id, userID)
	s.rdb.setArticleByID(ctx, id, article)
	return article, nil
}

func (s *ArticleService) getArticlesByTag(limit int) ([]Article, error) {
	return nil, nil
}
