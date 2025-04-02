package data

import (
	"context"
	"errors"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

// get intro from cache
func cacheGetIntro(ctx context.Context, rdb *redis.Client) (string, error) {
	key := getIntroKey()

	data, err := rdb.Get(ctx, key).Result()
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
func cacheSetIntro(ctx context.Context, rdb *redis.Client, intro string) error {
	key := getIntroKey()

	err := rdb.Set(ctx, key, intro, GetSitedataConfig().cacheBaseTTL).Err()
	if err != nil {
		logger.Error(
			"set intro failed",
		)
		return err
	}

	return nil
}
