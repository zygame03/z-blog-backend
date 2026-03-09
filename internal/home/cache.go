package home

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	rdb *redis.Client
	cfg func()
}

func newCache(rdb *redis.Client, cfg func()) *cache {
	return &cache{
		rdb: rdb,
		cfg: cfg,
	}
}
