package stats

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *cache {
	return &cache{
		rdb: rdb,
	}
}

func (c *cache) AddViewUV(ctx context.Context, ip string) error {
	return c.rdb.PFAdd(ctx, viewKey(), ip).Err()
}

func (c *cache) GetViewUV(ctx context.Context) (int64, error) {
	return c.rdb.PFCount(ctx, viewKey()).Result()
}

func (c *cache) DelViewUV(ctx context.Context) error {
	return c.rdb.Del(ctx, viewKey()).Err()
}
