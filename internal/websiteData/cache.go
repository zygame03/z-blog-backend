package websiteData

import (
	"context"
	"errors"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	rdb *redis.Client
	cfg func() *WebsiteDataConfig
}

func NewCache(rdb *redis.Client, cfg func() *WebsiteDataConfig) *cache {
	return &cache{
		rdb: rdb,
		cfg: cfg,
	}
}

// get intro from cache
func (c *cache) GetIntro(ctx context.Context) (string, error) {
	key := getIntroKey()

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return data, ErrCacheMiss
	}

	if err != nil {
		logger.Error(
			"cache get intro failed",
		)
		return "", err
	}

	return data, nil
}

// set intro to cache
func (c *cache) SetIntro(ctx context.Context, intro string) error {
	key := getIntroKey()

	err := c.rdb.Set(ctx, key, intro, c.cfg().CacheBaseTTL).Err()
	if err != nil {
		logger.Error(
			"set intro failed",
		)
		return err
	}

	return nil
}
