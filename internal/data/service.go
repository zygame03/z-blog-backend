package data

import (
	"context"
	"my_web/backend/internal/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type service struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func newService(db *gorm.DB, rdb *redis.Client) *service {
	s := service{
		DB:  db,
		RDB: rdb,
	}
	return &s
}

func (s *service) getIntro(ctx context.Context) (string, error) {
	data, err := cacheGetIntro(ctx, s.RDB)
	if err == nil {
		logger.Info(
			"get intro from cache failed",
		)
		return data, err
	}

	if err == ErrCacheMiss {
		logger.Info(
			"cache miss for intro",
		)
		data, err = repoGetIntro(s.DB)
		if err != nil {
			return data, err
		}

		cacheSetIntro(ctx, s.RDB, data)
		return data, err
	}

	return repoGetIntro(s.DB)
}
