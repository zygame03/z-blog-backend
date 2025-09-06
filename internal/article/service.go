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
	listIDs(ctx context.Context) ([]int, error)
	listByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error)
	getByID(ctx context.Context, id int) (*Article, error)
	listPopular(ctx context.Context, limit int) ([]ArticleSummary, error)
	incrementViews(ctx context.Context, id int, inc int64) error
	save(ctx context.Context, article *Article) (int, error)
	delete(ctx context.Context, id int) error
}

type ArticleCache interface {
	getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error)
	setArticlesByPage(ctx context.Context, page, pageSize int, articles []ArticleSummary, total int) error
	getArticleByID(ctx context.Context, id int) (*Article, error)
	setArticleByID(ctx context.Context, id int, article *Article) error
	getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error)
	setArticlesByPopular(ctx context.Context, limit int, articles []ArticleSummary) error
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
	ids, err := s.db.listIDs(ctx)
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

func (s *ArticleService) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error) {
	articles, total, err := s.rdb.getArticlesByPage(ctx, page, pageSize)
	if err == nil {
		logger.Info(
			"get articles by page from cache",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
		return articles, total, nil
	}

	if err != ErrCacheMiss {
		logger.Error(
			"get articles by page from cache failed",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
	} else {
		logger.Info(
			"cache miss for articles by page",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
	}

	articles, total, err = s.db.listByPage(ctx, page, pageSize)
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

func (s *ArticleService) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
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
	} else {
		logger.Error(
			"get articles by popular from cache failed",
			zap.Int("limit", limit),
			zap.Error(err),
		)
	}

	articles, err = s.db.listPopular(ctx, limit)
	if err != nil {
		logger.Error(
			"repo get articles by popular failed",
			zap.Error(err),
		)
		return nil, err
	}

	s.rdb.setArticlesByPopular(ctx, limit, articles)
	return articles, nil
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

	article, err = s.db.getByID(ctx, id)
	if err != nil {
		logger.Error(
			"repo get article by id failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	s.rdb.addViewUV(ctx, id, userID)
	s.rdb.setArticleByID(ctx, id, article)
	return article, nil
}

func (s *ArticleService) save(ctx context.Context, article *Article) (int, error) {
	return s.db.save(ctx, article)
}

func (s *ArticleService) delete(ctx context.Context, id int) error {
	return s.db.delete(ctx, id)
}
