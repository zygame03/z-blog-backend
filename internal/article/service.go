package article

import (
	"context"
	"my_web/backend/internal/logger"
	"my_web/backend/internal/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	DB  *gorm.DB
	RDB *redis.Client

	task utils.TaskRunner
}

func newArticleService(ctx context.Context, db *gorm.DB, rdb *redis.Client) *service {
	service := &service{
		DB:  db,
		RDB: rdb,
	}

	service.task = *utils.NewTaskRunner(
		service,
		utils.WithInterval(GetArticleConfig().syncInterval),
		utils.WithTimeout(GetArticleConfig().syncInterval),
	)

	service.task.Start(ctx)

	return service
}

func (s *service) Run(ctx context.Context) {
	ids, err := repoGetAllArticleIDs(s.DB)
	if err != nil {
		logger.Error(
			"load article ids for view sync failed",
			zap.Error(err),
		)
		return
	}

	for _, id := range ids {
		num, err := cacheGetViewUV(ctx, s.RDB, id)
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

		err = cacheDelViewUV(ctx, s.RDB, id)
		if err != nil {
			logger.Error(
				"delete view uv from cache failed",
				zap.Int("id", id),
				zap.Error(err),
			)
			continue
		}

		if err := repoIncrementViews(s.DB, id, num); err != nil {
			logger.Error(
				"increment article views failed",
				zap.Int("id", id),
				zap.Int64("increment", num),
				zap.Error(err),
			)
		}
	}
}

func (s *service) GetArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleWithoutContent, int, error) {
	articles, total, err := cacheGetArticlesByPage(ctx, s.RDB, page, pageSize)
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

		articles, total, err = repoGetArticlesByPage(s.DB, page, pageSize)
		if err != nil {
			return nil, 0, err
		}

		cacheSetArticlesByPage(ctx, s.RDB, page, pageSize, articles, total)
		return articles, total, nil
	}

	logger.Error(
		"get articles by page from cache failed",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Error(err),
	)

	return repoGetArticlesByPage(s.DB, page, pageSize)
}

func (s *service) GetArticlesByPopular(ctx context.Context, limit int) ([]ArticleWithoutContent, error) {
	articles, err := cacheGetArticlesByPopular(ctx, s.RDB, limit)
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

		articles, err = repoGetArticlesByPopular(s.DB, limit)
		if err != nil {
			return nil, err
		}

		go cacheSetArticlesByPopular(ctx, s.RDB, limit, articles)
		return articles, nil
	}

	logger.Error(
		"get articles by popular from cache failed",
		zap.Int("limit", limit),
		zap.Error(err),
	)

	return repoGetArticlesByPopular(s.DB, limit)
}

func (s *service) GetArticleByID(ctx context.Context, id int, userID string) (*Article, error) {
	article, err := cacheGetArticleByID(ctx, s.RDB, id)
	if err == nil {
		logger.Info(
			"get article by id from cache",
			zap.Int("id", id),
			zap.String("user_id", userID),
		)
		cacheAddViewUV(ctx, s.RDB, id, userID)
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

	article, err = repoGetArticleByID(s.DB, id)
	if err != nil {
		return nil, err
	}

	cacheAddViewUV(ctx, s.RDB, id, userID)
	cacheSetArticleByID(ctx, s.RDB, id, article)
	return article, nil
}

func (s *service) GetArticlesByTag(limit int) ([]Article, error) {
	return nil, nil
}
