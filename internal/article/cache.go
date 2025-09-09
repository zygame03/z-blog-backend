package article

import (
	"context"
	"encoding/json"
	"fmt"
	"my_web/backend/internal/zerrors"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	rdb *redis.Client
	cfg func() *Config
}

func newCache(rdb *redis.Client, cfg func() *Config) *cache {
	return &cache{
		rdb: rdb,
		cfg: cfg,
	}
}

func (c *cache) getArticlesByPage(ctx context.Context, page, pageSize int) ([]ArticleSummary, int, error) {
	key := articleByPageKey(page, pageSize)
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, 0, zerrors.CacheMiss
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	var articles []ArticleSummary
	if err = json.Unmarshal([]byte(data), &articles); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", zerrors.ErrUnmarshal, err)
	}

	totalKey := articleTotalKey()
	totalData, err := c.rdb.Get(ctx, totalKey).Result()
	if err == redis.Nil {
		return articles, 0, zerrors.CacheMiss
	}
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	total, err := strconv.Atoi(totalData)
	if err != nil {
		return nil, 0, zerrors.ErrParse
	}

	return articles, total, nil
}

func (c *cache) setArticlesByPage(ctx context.Context, page, pageSize int, articles []ArticleSummary, total int) error {
	// 序列化文章列表
	data, err := json.Marshal(articles)
	if err != nil {
		return zerrors.ErrMarshal
	}

	// 使用 Pipeline 批量设置
	pipe := c.rdb.Pipeline()
	pipe.Set(ctx, articleByPageKey(page, pageSize), data, c.cfg().CacheBaseTTL)
	pipe.Set(ctx, articleTotalKey(), strconv.Itoa(total), c.cfg().CacheBaseTTL)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	return nil
}

func (c *cache) getArticleByID(ctx context.Context, id int) (*Article, error) {
	key := articleByIDKey(id)
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, zerrors.CacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	var article Article
	if err := json.Unmarshal([]byte(data), &article); err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrUnmarshal, err)
	}

	return &article, nil
}

func (c *cache) setArticleByID(ctx context.Context, id int, article *Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrMarshal, err)
	}

	err = c.rdb.Set(ctx, articleByIDKey(id), data, c.cfg().CacheBaseTTL).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	return nil
}

func (c *cache) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	key := articleByPopularKey(limit)

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, zerrors.CacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	var articles []ArticleSummary
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		return nil, fmt.Errorf("%w: %v", zerrors.ErrUnmarshal, err)
	}

	return articles, nil
}

func (c *cache) setArticlesByPopular(ctx context.Context, limit int, articles []ArticleSummary) error {
	data, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrMarshal, err)
	}

	err = c.rdb.Set(ctx, articleByPopularKey(limit), data, c.cfg().CacheBaseTTL).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", zerrors.ErrCacheOperation, err)
	}

	return nil
}

func (c *cache) addViewUV(ctx context.Context, id int, userID string) error {
	return c.rdb.PFAdd(ctx, articleViewKey(id), userID).Err()
}

func (c *cache) getViewUV(ctx context.Context, id int) (int64, error) {
	return c.rdb.PFCount(ctx, articleViewKey(id)).Result()
}

func (c *cache) delViewUV(ctx context.Context, id int) error {
	return c.rdb.Del(ctx, articleViewKey(id)).Err()
}
