package websiteData

import (
	"context"
	"encoding/json"
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

func (c *cache) getAnnouncement(ctx context.Context) ([]string, error) {
	key := getAnnouncementKey()

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, err
	}

	var announcement []string
	err = json.Unmarshal([]byte(data), &announcement)
	if err != nil {
		logger.Error(
			"announcement unmarshal failed",
		)
		return nil, err
	}

	return announcement, nil
}

func (c *cache) setAnnouncement(ctx context.Context, data []string) error {
	key := getAnnouncementKey()

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = c.rdb.Set(ctx, key, b, c.cfg().CacheBaseTTL).Err()
	if err != nil {
		logger.Error(
			"set announcement failed",
		)
		return err
	}

	return nil
}
