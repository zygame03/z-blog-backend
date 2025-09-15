package site

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
		return fmt.Errorf("cache unmarshal failed: %w", err)
	}

	return nil
}

func (c *cache) set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set failed: %w", err)
	}

	return nil
}

// get intro from cache
func (c *cache) getIntro(ctx context.Context) (string, error) {
	key := getIntroKey()

	var intro string
	err := c.get(ctx, key, &intro)
	if err != nil {
		return "", err
	}

	return intro, nil
}

// set intro to cache
func (c *cache) setIntro(ctx context.Context, intro string) error {
	key := getIntroKey()

	return c.set(ctx, key, intro, c.cfg().CacheBaseTTL)
}

func (c *cache) listAnnouncements(ctx context.Context) ([]Announcement, error) {
	key := announcementKey()
	var data []Announcement
	err := c.get(ctx, key, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *cache) setAnnouncement(ctx context.Context, data []Announcement) error {
	key := announcementKey()
	return c.set(ctx, key, data, c.cfg().CacheBaseTTL)
}

func (c *cache) getDanmaku(ctx context.Context) ([]DanmakuVO, error) {
	var list []DanmakuVO
	key := danmakuKey()
	err := c.get(ctx, key, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *cache) setDanmaku(ctx context.Context, danmaku []DanmakuVO) error {
	key := danmakuKey()
	err := c.set(ctx, key, danmaku, c.cfg().CacheBaseTTL)
	if err != nil {
		return err
	}
	return nil
}
