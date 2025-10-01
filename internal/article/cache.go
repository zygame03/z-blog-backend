package article

import (
	"context"
	"encoding/json"
	"fmt"
	"my_web/backend/internal/zerrors"
	"time"

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

func (c *cache) get(ctx context.Context, key string, dest any) error {
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return zerrors.CacheMiss
	}
	if err != nil {
		return fmt.Errorf("cache get failed: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	return nil
}

func (c *cache) set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set failed: %w", err)
	}

	return nil
}

func (c *cache) getArticlesByPage(ctx context.Context, page, pageSize int) (*ArticlePageVO, error) {
	var data ArticlePageVO
	err := c.get(ctx, articleByPageKey(page, pageSize), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *cache) setArticlesByPage(ctx context.Context, articles *ArticlePageVO) error {
	err := c.set(ctx, articleByPageKey(articles.Page, articles.PageSize), articles, c.cfg().CacheBaseTTL)
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) getArticleByID(ctx context.Context, id int) (*Article, error) {
	var article Article
	err := c.get(ctx, articleByIDKey(id), &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (c *cache) setArticleByID(ctx context.Context, id int, article *Article) error {
	err := c.set(ctx, articleByIDKey(id), article, c.cfg().CacheBaseTTL)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) getArticleComments(ctx context.Context, id int, page, pageSize int) (*CommentPageVO, error) {
	var comments CommentPageVO
	err := c.get(ctx, commentsByIdKey(id, page, pageSize), &comments)
	if err != nil {
		return nil, fmt.Errorf("cache get comments failed: %w", err)
	}
	return &comments, nil
}

func (c *cache) setArticleComments(ctx context.Context, id int, comment *CommentPageVO) error {
	err := c.set(ctx, commentsByIdKey(id, comment.Page, comment.PageSize), comment, c.cfg().CacheBaseTTL)
	if err != nil {
		return fmt.Errorf("cache set comments failed: %w", err)
	}
	return nil
}

func (c *cache) getArticlesByPopular(ctx context.Context, limit int) ([]ArticleSummary, error) {
	var articles []ArticleSummary
	err := c.get(ctx, articleByPopularKey(limit), &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (c *cache) setArticlesByPopular(ctx context.Context, limit int, articles []ArticleSummary) error {
	err := c.set(ctx, articleByPopularKey(limit), articles, c.cfg().CacheBaseTTL)
	if err != nil {
		return err
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
