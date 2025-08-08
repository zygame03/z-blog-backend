package site

import (
	"context"
	"encoding/json"
	"errors"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	rdb  *redis.Client
	conf func() *Config
}

func newCache(rdb *redis.Client, conf func() *Config) *cache {
	return &cache{
		rdb:  rdb,
		conf: conf,
	}
}

// get intro from cache
func (c *cache) getIntro(ctx context.Context) (string, error) {
	key := getIntroKey()

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return data, ErrCacheMiss
	}

	if err != nil {
		return "", err
	}

	return data, nil
}

// set intro to cache
func (c *cache) setIntro(ctx context.Context, intro string) error {
	key := getIntroKey()

	err := c.rdb.Set(ctx, key, intro, c.conf().CacheBaseTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) getAnnouncement(ctx context.Context) ([]*announcementBO, error) {
	key := getAllAnnouncementKey()

	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, err
	}

	var announcement []*announcementBO
	err = json.Unmarshal([]byte(data), &announcement)
	if err != nil {
		return nil, err
	}

	return announcement, nil
}

func (c *cache) setAnnouncement(ctx context.Context, data []*announcementBO) error {
	pipe := c.rdb.Pipeline()

	for _, v := range data {
		err := pipe.Set(ctx, getAnnouncementKey(v.Id), v.Text, c.conf().CacheBaseTTL).Err()
		if err != nil {
			continue
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error(
			"cache execute pipeline commands failed",
			zap.Error(err),
		)
		return err
	}

	return nil
}
