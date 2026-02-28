package article

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"my_web/backend/internal/logger"
	"strconv"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	rdb *redis.Client
	cfg func() *ArticleConfig
}

func NewCache(rdb *redis.Client, cfg func() *ArticleConfig) *cache {
	return &cache{
		rdb: rdb,
		cfg: cfg,
	}
}

func (c *cache) GetArticlesByPage(
	ctx context.Context,
	page, pageSize int,
) ([]ArticleWithoutContent, int, error) {
	// 获取文章列表
	key := ArticleByPageKey(page, pageSize)
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, 0, ErrCacheMiss
	}
	if err != nil {
		logger.Error(
			"cache get articles by page failed",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("cache get articles by page failed: %w", err)
	}

	var articles []ArticleWithoutContent
	if err = json.Unmarshal([]byte(data), &articles); err != nil {
		logger.Error(
			"unmarshal articles by page failed",
			zap.String("key", key),
			zap.String("data", data),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("unmarshal articles by page failed: %w", err)
	}

	// 获取总数
	totalKey := ArticleTotalKey()
	totalData, err := c.rdb.Get(ctx, totalKey).Result()
	if err == redis.Nil {
		return articles, 0, ErrCacheMiss // 列表有但总数未命中
	}
	if err != nil {
		logger.Error(
			"cache get article total failed",
			zap.String("key", totalKey),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("cache get article total failed: %w", err)
	}

	total, err := strconv.Atoi(totalData)
	if err != nil {
		logger.Error(
			"parse article total failed",
			zap.String("key", totalKey),
			zap.String("data", totalData),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("parse article total failed: %w", err)
	}

	return articles, total, nil
}

func (c *cache) SetArticlesByPage(
	ctx context.Context,
	page, pageSize int,
	articles []ArticleWithoutContent,
	total int,
) error {
	// 序列化文章列表
	data, err := json.Marshal(articles)
	if err != nil {
		logger.Error(
			"marshal articles by page failed",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return fmt.Errorf("marshal articles by page failed: %w", err)
	}

	// 使用 Pipeline 批量设置
	pipe := c.rdb.Pipeline()
	pipe.Set(
		ctx,
		ArticleByPageKey(page, pageSize),
		data,
		c.cfg().cacheBaseTTL,
	)
	pipe.Set(
		ctx,
		ArticleTotalKey(),
		strconv.Itoa(total),
		c.cfg().cacheBaseTTL,
	)

	_, err = pipe.Exec(ctx)
	if err != nil {
		logger.Error(
			"set articles by page cache failed",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err),
		)
		return fmt.Errorf("set articles by page cache failed: %w", err)
	}

	return nil
}

func (c *cache) GetArticleByID(
	ctx context.Context,
	id int,
) (*Article, error) {
	key := ArticleByIDKey(id)
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Info(
			"cache miss",
			zap.String("key", key),
		)
		return nil, ErrCacheMiss
	}
	if err != nil {
		logger.Error(
			"cache get article by id failed",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("cache get article by id failed: %w", err)
	}

	var article Article
	if err := json.Unmarshal([]byte(data), &article); err != nil {
		logger.Warn(
			"unmarshal failed",
			zap.String("key", key),
			zap.String("data", data),
			zap.Error(err),
		)
		return nil, fmt.Errorf("unmarshal article by id failed: %w", err)
	}

	return &article, nil
}

func (c *cache) SetArticleByID(
	ctx context.Context,
	id int,
	article *Article,
) error {
	data, err := json.Marshal(article)
	if err != nil {
		logger.Error(
			"marshal article by id failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("marshal article by id failed: %w", err)
	}

	err = c.rdb.Set(
		ctx,
		ArticleByIDKey(id),
		data,
		c.cfg().cacheBaseTTL,
	).Err()
	if err != nil {
		logger.Error(
			"set article by id cache failed",
			zap.Int("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("set article by id cache failed: %w", err)
	}

	return nil
}

func (c *cache) GetArticlesByPopular(
	ctx context.Context,
	limit int,
) ([]ArticleWithoutContent, error) {
	key := ArticleByPopularKey(limit)

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Info(
			"cache miss",
			zap.String("key", key),
		)
		return nil, ErrCacheMiss
	}
	if err != nil {
		logger.Error(
			"cache get articles by popular failed",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("cache get articles by popular failed: %w", err)
	}

	var articles []ArticleWithoutContent
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		logger.Error(
			"unmarshal articles by popular failed",
			zap.String("key", key),
			zap.String("data", data),
			zap.Error(err),
		)
		return nil, fmt.Errorf("unmarshal articles by popular failed: %w", err)
	}

	return articles, nil
}

func (c *cache) SetArticlesByPopular(
	ctx context.Context,
	limit int,
	articles []ArticleWithoutContent,
) error {
	data, err := json.Marshal(articles)
	if err != nil {
		logger.Error(
			"marshal articles by popular failed",
			zap.Int("limit", limit),
			zap.Error(err),
		)
		return fmt.Errorf("marshal articles by popular failed: %w", err)
	}

	err = c.rdb.Set(
		ctx,
		ArticleByPopularKey(limit),
		data,
		c.cfg().cacheBaseTTL,
	).Err()
	if err != nil {
		logger.Error(
			"set articles by popular cache failed",
			zap.Int("limit", limit),
			zap.Error(err),
		)
		return fmt.Errorf("set articles by popular cache failed: %w", err)
	}

	return nil
}

func (c *cache) AddViewUV(
	ctx context.Context,
	id int,
	userID string,
) error {
	return c.rdb.PFAdd(ctx, ArticleViewKey(id), userID).Err()
}

func (c *cache) GetViewUV(
	ctx context.Context,
	id int,
) (int64, error) {
	return c.rdb.PFCount(ctx, ArticleViewKey(id)).Result()
}

func (c *cache) DelViewUV(
	ctx context.Context,
	id int,
) error {
	return c.rdb.Del(ctx, ArticleViewKey(id)).Err()
}
