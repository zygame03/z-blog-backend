package site

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
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

	err = c.rdb.Set(ctx, key, b, c.conf().CacheBaseTTL).Err()
	if err != nil {
		return err
	}

	return nil
}
