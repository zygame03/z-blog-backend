package stats

import (
	"context"
	"my_web/backend/internal/global"
	"my_web/backend/internal/logger"
	"strconv"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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

func (c *cache) getView(ctx context.Context) (int, error) {
	view, err := c.rdb.Get(ctx, viewKey()).Result()
	if err == redis.Nil {
		logger.Error(
			"cache miss",
			zap.Error(err),
		)
		return -1, global.ErrCacheMiss
	}
	if err != nil {
		logger.Error(
			"cache get view failed",
			zap.Error(err),
		)
		return -1, err
	}

	view_i, err := strconv.Atoi(view)
	if err != nil {
		logger.Error(
			"cache view atoi failed",
			zap.String("value", view),
			zap.Error(err),
		)
		return -1, err
	}

	return view_i, nil
}

func (c *cache) setView(ctx context.Context, num int) error {
	return c.rdb.Set(ctx, viewKey(), num, c.conf().CacheBaseTTL).Err()

}

func (c *cache) addViewUV(ctx context.Context, ip string) error {
	return c.rdb.PFAdd(ctx, viewUVKey(), ip).Err()
}

func (c *cache) getViewUV(ctx context.Context) (int64, error) {
	return c.rdb.PFCount(ctx, viewUVKey()).Result()
}

func (c *cache) delViewUV(ctx context.Context) error {
	return c.rdb.Del(ctx, viewUVKey()).Err()
}
