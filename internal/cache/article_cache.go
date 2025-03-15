package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"my_web/backend/internal/models"
	"my_web/backend/internal/repo"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ArticleCache struct {
	RDB *redis.Client
}

func (c *ArticleCache) GetArticlesByPage(ctx context.Context, page, pageSize int) ([]repo.ArticleVO, int, error) {
	// 获取文章列表
	key := ArticleByPageKey(page, pageSize)
	data, err := c.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, 0, ErrCacheMiss // 缓存未命中
	}
	if err != nil {
		return nil, 0, fmt.Errorf("缓存获取异常 %w", err)
	}

	var articles []repo.ArticleVO
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		return nil, 0, fmt.Errorf("反序列化失败 %w", err)
	}

	// 获取总数
	totalData, err := c.RDB.Get(ctx, ArticleTotal()).Result()
	if err == redis.Nil {
		return articles, 0, ErrCacheMiss // 列表有但总数未命中
	}
	if err != nil {
		return nil, 0, fmt.Errorf("获取总数失败 %w", err)
	}

	total, err := strconv.Atoi(totalData)
	if err != nil {
		return nil, 0, fmt.Errorf("总数解析失败 %w", err)
	}

	return articles, total, nil
}

func (c *ArticleCache) SetArticlesByPage(ctx context.Context, page, pageSize int, articles []repo.ArticleVO, total int) error {
	// 设置过期时间
	expire := 30 * time.Minute

	// 序列化文章列表
	data, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("序列化失败 %w", err)
	}

	// 使用 Pipeline 批量设置
	pipe := c.RDB.Pipeline()
	pipe.Set(ctx, ArticleByPageKey(page, pageSize), data, expire)
	pipe.Set(ctx, ArticleTotal(), strconv.Itoa(total), expire)

	_, err = pipe.Exec(ctx)
	return err
}

func (c *ArticleCache) GetArticleByID(ctx context.Context, id int) (*models.Article, error) {
	key := ArticleByIDKey(id)

	data, err := c.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("缓存获取异常 %w", err)
	}

	var article models.Article
	if err := json.Unmarshal([]byte(data), &article); err != nil {
		return nil, fmt.Errorf("反序列化失败 %w", err)
	}

	return &article, nil
}

func (c *ArticleCache) SetArticleByID(ctx context.Context, id int, article *models.Article) error {
	expire := 1 * time.Hour

	data, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("序列化失败 %w", err)
	}

	return c.RDB.Set(ctx, ArticleByIDKey(id), data, expire).Err()
}

func (c *ArticleCache) GetArticlesByPopular(ctx context.Context, limit int) ([]repo.ArticleVO, error) {
	key := ArticleByPopularKey(limit)

	data, err := c.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("缓存获取异常 %w", err)
	}

	var articles []repo.ArticleVO
	if err := json.Unmarshal([]byte(data), &articles); err != nil {
		return nil, fmt.Errorf("反序列化失败 %w", err)
	}

	return articles, nil
}

func (c *ArticleCache) SetArticlesByPopular(ctx context.Context, limit int, articles []repo.ArticleVO) error {
	expire := 1 * time.Hour

	data, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("序列化失败 %w", err)
	}

	return c.RDB.Set(ctx, ArticleByPopularKey(limit), data, expire).Err()
}

func (c *ArticleCache) SetViewsIncr(ctx context.Context, id int) error {
	// 使用 INCR 命令增加 views
	viewsKey := fmt.Sprintf("Article:Views:%d", id)
	_, err := c.RDB.Incr(ctx, viewsKey).Result()
	if err != nil {
		return fmt.Errorf("增加 views 失败 %w", err)
	}

	// 设置过期时间（24小时）
	c.RDB.Expire(ctx, viewsKey, 24*time.Hour)

	return nil
}
