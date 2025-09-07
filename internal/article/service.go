package article

import (
	"context"
	"errors"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/zerrors"

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
		"cron job registered",
		zap.String("job", "syncArticleViews"),
		zap.Duration("interval", s.cfg().SyncInterval),
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

	var (
		total   = len(ids)
		success int
		failed  int
	)
	for _, id := range ids {
		num, err := s.rdb.getViewUV(ctx, id)
		if err != nil {
			failed++
			continue
		}
		if num == 0 {
			success++
			continue
		}
		err = s.rdb.delViewUV(ctx, id)
		if err != nil {
			failed++
			continue
		}
		err = s.db.incrementViews(ctx, id, num)
		if err != nil {
			failed++
			continue
		}
		success++
	}

	logger.Info(
		"sync article views done",
		zap.Int("total", total),
		zap.Int("success", success),
		zap.Int("failed", failed),
	)
}

func (s *ArticleService) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error) {
	articles, total, err := s.rdb.getArticlesByPage(ctx, page, pageSize)
	if err == nil {
		return articles, total, nil
	}

	articles, total, err = s.db.listByPage(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	s.rdb.setArticlesByPage(ctx, page, pageSize, articles, total)
	return articles, total, nil
}

func (s *ArticleService) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	articles, err := s.rdb.getArticlesByPopular(ctx, limit)
	if err == nil {
		return articles, nil
	}
	if !errors.Is(err, zerrors.ErrCacheMiss) {
		// 这里可能需要输出缓存异常日志
	}

	articles, err = s.db.listPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	s.rdb.setArticlesByPopular(ctx, limit, articles)
	return articles, nil
}

func (s *ArticleService) getArticleByID(ctx context.Context, id int, userID string) (*Article, error) {
	article, err := s.rdb.getArticleByID(ctx, id)
	if err == nil {
		s.rdb.addViewUV(ctx, id, userID)
		return article, nil
	}
	if !errors.Is(err, zerrors.ErrCacheMiss) {
		// 这里可能需要输出缓存异常日志
	}

	article, err = s.db.getByID(ctx, id)
	if err != nil {
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
